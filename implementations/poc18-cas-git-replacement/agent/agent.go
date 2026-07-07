package agent

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/carbundle"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/economy"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/transport"
)

// Config describes one Docker agent container.
type Config struct {
	Name      string
	Listen    string
	CASRoot   string
	Collector string
	Peers     map[string]string
}

// Runtime owns one agent's local CAS, TCP listener, and observer-only event
// stream. It is intentionally small but shaped like the POC16 supervisor/kernel
// split: local runtime state is private, while exact message copies go one-way
// to the collector.
type Runtime struct {
	cfg       Config
	cas       *store.FileStore
	collector *eventstream.Client
	listener  net.Listener
	headsMu   sync.Mutex
	heads     map[string]cidlib.Cid
	done      chan struct{}
}

// receivedMessage keeps the exact TCP frame bytes and parsed view together so
// parent links and observer artifacts refer to the real message CID, not a
// reconstructed approximation.
type receivedMessage struct {
	View  graph.EnvelopeView
	Bytes []byte
	CID   cidlib.Cid
	Peer  string
	Kind  string
}

// ExchangeID names one promise conversation inside a persistent TCP session.
// It is local diagnostic vocabulary, not a pCID, address, command, or authority.
type ExchangeID string

// retrievalRequest describes one voluntary request for latest or specific
// missing object bytes. The wanted object rows stay inside the pCID-owned
// promise body so pCID never becomes an address or operation code.
type retrievalRequest struct {
	ExchangeID ExchangeID
	Scope      string
	Objects    []carbundle.Block
}

// availabilityOffer is the local parse result for an object_availability
// promise. The signed retrieval token is mandatory; storage payment tokens are
// optional reciprocal economics.
type availabilityOffer struct {
	ExchangeID      ExchangeID
	Objects         []carbundle.Block
	RetrievalToken  []byte
	PaymentToken    []byte
	PaymentTokenCID string
}

// Session owns one persistent TCP connection to a peer and can carry multiple
// named exchanges before closing.
type Session struct {
	runtime   *Runtime
	peer      string
	conn      *transport.Conn
	exchanges int
}

// Run starts one agent, executes its deterministic POC18 role, and shuts down.
func Run(cfg Config) error {
	if cfg.Name == "" || cfg.Listen == "" || cfg.CASRoot == "" || cfg.Collector == "" {
		return fmt.Errorf("agent name, listen address, CAS root, and collector are required")
	}
	cas, casErr := store.Open(cfg.CASRoot)
	if casErr != nil {
		return casErr
	}
	collector, collectorErr := eventstream.Dial("agent:"+cfg.Name, cfg.Collector, 20*time.Second)
	if collectorErr != nil {
		return collectorErr
	}
	runtime := &Runtime{cfg: cfg, cas: cas, collector: collector, heads: map[string]cidlib.Cid{}, done: make(chan struct{})}
	defer runtime.closeCollector()
	if startErr := runtime.startListener(); startErr != nil {
		return startErr
	}
	defer runtime.closeListener()
	if eventErr := runtime.record("agent_started", "kept", "", "listen="+cfg.Listen); eventErr != nil {
		return eventErr
	}
	runErr := runtime.runRole()
	if runErr != nil {
		if recordErr := runtime.record("agent_failed", "broken", "", runErr.Error()); recordErr != nil {
			return recordErr
		}
		return runErr
	}
	if eventErr := runtime.record("agent_completed", "kept", "", "role completed"); eventErr != nil {
		return eventErr
	}
	return runtime.emit(eventstream.Record{Kind: eventstream.KindSupervisorDone, Detail: cfg.Name})
}

func (runtime *Runtime) startListener() error {
	listener, listenErr := net.Listen("tcp", runtime.cfg.Listen)
	if listenErr != nil {
		return listenErr
	}
	runtime.listener = listener
	go runtime.acceptLoop()
	return nil
}

func (runtime *Runtime) acceptLoop() {
	for {
		conn, acceptErr := runtime.listener.Accept()
		if acceptErr != nil {
			select {
			case <-runtime.done:
				return
			default:
				runtime.recordBestEffort("agent_accept_failed", "broken", "", acceptErr.Error())
				return
			}
		}
		go runtime.handleConn(transport.Wrap(conn))
	}
}

func (runtime *Runtime) handleConn(conn *transport.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			runtime.recordBestEffort("agent_connection_close_failed", "broken", "", closeErr.Error())
		}
	}()
	for {
		// Intent: A retrieval exchange is a promise conversation on one TCP stream:
		// interest, availability, redemption, and object bytes. The listener must
		// keep the accepted stream open until the peer closes it, otherwise the
		// redemption promise is forced onto a broken pipe. Source: DI-koriz
		frame, readErr := conn.ReadFrame()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return
			}
			runtime.recordBestEffort("agent_frame_read_failed", "broken", "", readErr.Error())
			return
		}
		message, storeErr := runtime.receiveMessage(frame, "received", "peer")
		if storeErr != nil {
			runtime.recordBestEffort("agent_receive_failed", "broken", "", storeErr.Error())
			return
		}
		switch message.Kind {
		case "sync_interest":
			runtime.handleSyncInterest(conn, message)
		case "object_availability":
			runtime.handleObjectAvailability(conn, message)
		case "object_retrieval_redemption":
			runtime.handleRetrievalRedemption(conn, message)
		case "object_bytes":
			if handleErr := runtime.handleObjectBytes(message, message.Peer); handleErr != nil {
				runtime.recordBestEffort("object_bytes_handle_failed", "broken", message.Peer, handleErr.Error())
			}
		default:
			runtime.recordBestEffort("agent_message_recorded", "kept", message.Peer, "kind="+message.Kind)
		}
	}
}

func (runtime *Runtime) runRole() error {
	switch runtime.cfg.Name {
	case "alice":
		if err := runtime.seedAlice(); err != nil {
			return err
		}
		time.Sleep(9 * time.Second)
	case "bob":
		time.Sleep(1200 * time.Millisecond)
		// Intent: Bob opens one persistent TCP session and runs two independent
		// exchanges on it so POC18 proves multiplexing by message CID and exchange
		// label, not by separate throwaway TCP connections. Source: DI-biruf
		if err := runtime.requestLatestBatch("alice", []ExchangeID{"bob-from-alice", "bob-audit-from-alice"}); err != nil {
			return err
		}
		time.Sleep(7 * time.Second)
	case "carol":
		time.Sleep(2800 * time.Millisecond)
		chosenPeer := runtime.choosePeer("scheduler-continuous-sync", []string{"mallory", "bob"})
		if err := runtime.requestLatest(chosenPeer, "carol-from-"+chosenPeer); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	case "dave":
		time.Sleep(3600 * time.Millisecond)
		if err := runtime.requestLatest("bob", "dave-from-bob"); err != nil {
			return err
		}
		if err := runtime.createLocalPromise("merge_snapshot", "alice", "Dave promises a local merge view derived from Bob's latest advertised snapshot."); err != nil {
			return err
		}
		time.Sleep(4 * time.Second)
	case "ellen":
		time.Sleep(4800 * time.Millisecond)
		if err := runtime.requestLatest("dave", "ellen-from-dave"); err != nil {
			return err
		}
		if err := runtime.createLocalPromise("review_statement", "dave", "Ellen promises she locally reviewed Dave's latest merge view."); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	case "frank":
		time.Sleep(5200 * time.Millisecond)
		if err := runtime.requestLatest("alice", "frank-from-alice"); err != nil {
			return err
		}
		if err := runtime.createLocalPromise("object_retention", "alice", "Frank promises local retention of Alice's latest CID under local storage constraints."); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	case "mallory":
		time.Sleep(8 * time.Second)
	default:
		return fmt.Errorf("unknown agent %s", runtime.cfg.Name)
	}
	return nil
}

func (runtime *Runtime) seedAlice() error {
	chunkContent := []byte("POC18 TCP seed: README content plus POSIX-style graph parents.\n")
	chunkEntry, chunkErr := runtime.cas.Put("chunk", chunkContent)
	if chunkErr != nil {
		return chunkErr
	}
	chunkCID, chunkParseErr := store.ParseCIDText(chunkEntry.CID)
	if chunkParseErr != nil {
		return chunkParseErr
	}
	manifestMessage, manifestErr := graph.StoreMessage(runtime.cas, nil, graph.Payload{
		Promiser:          "alice",
		Promisee:          "bob",
		PromiseKind:       "chunk_manifest",
		PromiseBody:       graph.ChunkManifestBody(chunkCID, int64(len(chunkContent)), "single-chunk-poc18", []any{"deterministic seed chunk"}, []any{graph.ObjectRow("chunk", chunkCID, chunkEntry.Size)}, store.CIDText(chunkCID)),
		ReciprocalPromise: []any{},
		LocalConstraints:  []any{"local seed only; cross-agent copies must move over TCP"},
	})
	if manifestErr != nil {
		return manifestErr
	}
	posixMessage, posixErr := graph.StoreMessage(runtime.cas, []graph.Parent{{Role: "chunk_manifest", CID: manifestMessage.CID}}, graph.Payload{
		Promiser:          "alice",
		Promisee:          "bob",
		PromiseKind:       "posix_node",
		PromiseBody:       graph.PosixNodeBody("node:readme", "regular_file", graph.Target("chunk_manifest", manifestMessage.CID), []any{"mode:0644"}, []any{"materialize as README.md when locally chosen"}),
		ReciprocalPromise: []any{},
		LocalConstraints:  []any{"local materialization constraints apply"},
	})
	if posixErr != nil {
		return posixErr
	}
	refSetMessage, refSetErr := graph.StoreMessage(runtime.cas, []graph.Parent{{Role: "directory_entry", CID: posixMessage.CID}}, graph.Payload{
		Promiser:          "alice",
		Promisee:          "bob",
		PromiseKind:       "reference_set",
		PromiseBody:       graph.ReferenceSetBody("dir:root", "directory", "alice-workspace", []any{graph.ReferenceEntry("README.md", []any{graph.Target("posix_node", posixMessage.CID)}, []any{"default visible file"})}, []any{"directory labels are local names, not file identity"}),
		ReciprocalPromise: []any{},
		LocalConstraints:  []any{"directory promise is local to Alice's workspace view"},
	})
	if refSetErr != nil {
		return refSetErr
	}
	snapshotMessage, snapshotErr := graph.StoreMessage(runtime.cas, []graph.Parent{{Role: "root_directory", CID: refSetMessage.CID}}, graph.Payload{
		Promiser:          "alice",
		Promisee:          "bob",
		PromiseKind:       "snapshot",
		PromiseBody:       graph.SnapshotBody("snapshot:poc18-tcp-seed", refSetMessage.CID, nil, "Alice promises this deterministic local seed snapshot from her own CAS.", []any{"local seed only"}),
		ReciprocalPromise: []any{},
		LocalConstraints:  []any{"cross-agent copies must move over TCP"},
	})
	if snapshotErr != nil {
		return snapshotErr
	}
	runtime.setHead("latest", snapshotMessage.CID)
	if recordErr := runtime.record("workspace_scenario_rich", "kept", "bob", "snapshot="+store.CIDText(snapshotMessage.CID)+" parents="+store.CIDText(refSetMessage.CID)); recordErr != nil {
		return recordErr
	}
	return runtime.record("alice_seeded_local_snapshot", "kept", "bob", "cid="+store.CIDText(snapshotMessage.CID))
}

func (runtime *Runtime) createLocalPromise(kind string, promisee string, text string) error {
	parentRows := []any{}
	runtime.headsMu.Lock()
	if head, ok := runtime.heads["latest"]; ok {
		parentRows = append(parentRows, store.LinkTag(head))
	}
	runtime.headsMu.Unlock()
	payload := graph.Payload{
		Promiser:          runtime.cfg.Name,
		Promisee:          promisee,
		PromiseKind:       kind,
		PromiseBody:       []any{kind + ":" + runtime.cfg.Name, parentRows, text},
		ReciprocalPromise: []any{},
		LocalConstraints:  []any{"local judgment only"},
	}
	message, storeErr := graph.StoreMessage(runtime.cas, nil, payload)
	if storeErr != nil {
		return storeErr
	}
	runtime.setHead("latest", message.CID)
	return runtime.record("local_promise_created", "kept", promisee, "kind="+kind+" cid="+store.CIDText(message.CID))
}

func (runtime *Runtime) choosePeer(goal string, candidates []string) string {
	// Intent: Carol's scheduler path should make an explicit local trust choice
	// before opening TCP, proving that peer selection is not hidden in transport
	// setup or treated as a global routing authority. Source: DI-biruf
	bestPeer := ""
	bestScore := -1 << 30
	details := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		score := runtime.localTrustScore(candidate)
		details = append(details, fmt.Sprintf("%s=%d", candidate, score))
		if bestPeer == "" || score > bestScore {
			bestPeer = candidate
			bestScore = score
		}
	}
	runtime.recordBestEffort("trust_peer_choice", "kept", bestPeer, "goal="+goal+" candidates="+strings.Join(details, ","))
	return bestPeer
}

func (runtime *Runtime) localTrustScore(peer string) int {
	switch peer {
	case "bob":
		return 80
	case "alice":
		return 75
	case "dave":
		return 65
	case "frank":
		return 60
	case "mallory":
		return -80
	default:
		return 10
	}
}

func (runtime *Runtime) requestLatest(peer string, label string) error {
	return runtime.requestLatestBatch(peer, []ExchangeID{ExchangeID(label)})
}

func (runtime *Runtime) requestLatestBatch(peer string, exchangeIDs []ExchangeID) error {
	requests := make([]retrievalRequest, 0, len(exchangeIDs))
	for _, exchangeID := range exchangeIDs {
		requests = append(requests, retrievalRequest{ExchangeID: exchangeID, Scope: "latest"})
	}
	return runtime.requestObjectsBatch(peer, requests)
}

func (runtime *Runtime) requestObjects(peer string, exchangeID ExchangeID, objects []carbundle.Block) error {
	return runtime.requestObjectsBatch(peer, []retrievalRequest{{ExchangeID: exchangeID, Scope: "objects", Objects: objects}})
}

func (runtime *Runtime) requestObjectsBatch(peer string, requests []retrievalRequest) error {
	if len(requests) == 0 {
		return fmt.Errorf("at least one retrieval request is required")
	}
	// Intent: A scheduler-style sync round must produce real TCP promise
	// exchanges. Multiple retrieval requests may share one Session, but each keeps
	// its own ExchangeID and message-CID parent chain. Source: DI-biruf
	session, sessionErr := runtime.openSession(peer)
	if sessionErr != nil {
		return sessionErr
	}
	defer session.Close()
	if recordErr := runtime.record("scheduler_tcp_sync", "kept", peer, fmt.Sprintf("exchanges=%d", len(requests))); recordErr != nil {
		return recordErr
	}
	for _, request := range requests {
		if request.ExchangeID == "" || request.Scope == "" {
			return fmt.Errorf("exchange id and scope are required")
		}
		if err := session.runRetrievalExchange(request); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *Runtime) openSession(peer string) (*Session, error) {
	address, ok := runtime.cfg.Peers[peer]
	if !ok {
		return nil, fmt.Errorf("missing peer address for %s", peer)
	}
	conn, dialErr := transport.Dial(address, 10*time.Second)
	if dialErr != nil {
		return nil, dialErr
	}
	if recordErr := runtime.record("session_opened", "kept", peer, "persistent TCP session"); recordErr != nil {
		if closeErr := conn.Close(); closeErr != nil {
			runtime.recordBestEffort("session_open_cleanup_failed", "broken", peer, closeErr.Error())
		}
		return nil, recordErr
	}
	return &Session{runtime: runtime, peer: peer, conn: conn}, nil
}

func (session *Session) runRetrievalExchange(request retrievalRequest) error {
	session.exchanges++
	if recordErr := session.runtime.record("exchange_started", "kept", session.peer, "exchange="+string(request.ExchangeID)+" scope="+request.Scope); recordErr != nil {
		return recordErr
	}
	interest, interestErr := session.runtime.newMessage(nil, graph.Payload{
		Promiser:          session.runtime.cfg.Name,
		Promisee:          session.peer,
		PromiseKind:       "sync_interest",
		PromiseBody:       []any{string(request.ExchangeID), request.Scope, missingRows(request.Objects), []any{"I promise to receive exact CID-verified object bytes and evaluate them locally."}, []any{"peer may decline if objects are not locally promised"}},
		ReciprocalPromise: []any{"I may reciprocate with local storage or forwarding promises."},
		LocalConstraints:  []any{"no shared CAS; TCP promise exchange required"},
	})
	if interestErr != nil {
		return interestErr
	}
	if sendErr := session.Send(interest); sendErr != nil {
		return sendErr
	}
	availability, receiveErr := session.Read()
	if receiveErr != nil {
		return receiveErr
	}
	offer, parseErr := availabilityOfferFromView(availability.View)
	if parseErr != nil {
		return parseErr
	}
	if offer.ExchangeID != request.ExchangeID {
		return fmt.Errorf("availability exchange=%s, want %s", offer.ExchangeID, request.ExchangeID)
	}
	if _, verifyErr := economy.VerifyRetrievalCapability(offer.RetrievalToken, economy.ExpectedRetrievalCapability{
		Issuer:     session.peer,
		Subject:    session.runtime.cfg.Name,
		Scope:      request.Scope,
		ObjectCIDs: blockCIDTexts(offer.Objects),
	}, time.Now()); verifyErr != nil {
		return verifyErr
	}
	if err := session.runtime.redeemStoragePaymentIfUseful(session.peer, offer); err != nil {
		return err
	}
	redemption, redemptionErr := session.runtime.newMessage([]graph.Parent{{Role: "redeems_availability", CID: availability.CID}}, graph.Payload{
		Promiser:          session.runtime.cfg.Name,
		Promisee:          session.peer,
		PromiseKind:       "object_retrieval_redemption",
		PromiseBody:       []any{string(request.ExchangeID), offer.RetrievalToken, missingRows(offer.Objects), "I promise to evaluate the returned CAR bytes by exact CID before storing."},
		ReciprocalPromise: []any{"local reciprocal economics remain agent-local in this POC18 slice"},
		LocalConstraints:  []any{"collector artifacts do not affect trust or routing"},
	})
	if redemptionErr != nil {
		return redemptionErr
	}
	if sendErr := session.Send(redemption); sendErr != nil {
		return sendErr
	}
	objectBytesMessage, objectReceiveErr := session.Read()
	if objectReceiveErr != nil {
		return objectReceiveErr
	}
	if err := session.runtime.handleObjectBytes(objectBytesMessage, session.peer); err != nil {
		return err
	}
	return session.runtime.record("exchange_completed", "kept", session.peer, "exchange="+string(request.ExchangeID)+" scope="+request.Scope)
}

func (session *Session) Send(message outgoingMessage) error {
	return session.runtime.sendMessage(session.conn, session.peer, message, "sent")
}

func (session *Session) Read() (receivedMessage, error) {
	frame, readErr := session.conn.ReadFrame()
	if readErr != nil {
		return receivedMessage{}, readErr
	}
	return session.runtime.receiveMessage(frame, "received", session.peer)
}

func (session *Session) Close() {
	if closeErr := session.conn.Close(); closeErr != nil {
		session.runtime.recordBestEffort("session_close_failed", "broken", session.peer, closeErr.Error())
		return
	}
	session.runtime.recordBestEffort("session_closed", "kept", session.peer, fmt.Sprintf("exchanges=%d", session.exchanges))
}

func (runtime *Runtime) objectsForRequest(wanted []carbundle.Block) ([]carbundle.Block, error) {
	// Intent: Providers answer from their own sparse CAS only. An explicit wanted
	// list is a request for named CIDs, not permission for the requester to read a
	// shared store. Source: DI-biruf
	if len(wanted) == 0 {
		return runtime.currentHeads(), nil
	}
	objects := make([]carbundle.Block, 0, len(wanted))
	for _, wantedBlock := range wanted {
		objectCID, parseErr := store.ParseCIDText(wantedBlock.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		_, entry, getErr := runtime.cas.Get(objectCID)
		if getErr != nil {
			return nil, getErr
		}
		kind := wantedBlock.Kind
		if kind == "" {
			kind = entry.Kind
		}
		objects = append(objects, carbundle.Block{CID: entry.CID, Kind: kind, Size: entry.Size})
	}
	return objects, nil
}

func (runtime *Runtime) optionalStoragePaymentToken(peer string, objects []carbundle.Block, interestCID cidlib.Cid) ([]byte, string, error) {
	if runtime.cfg.Name != "alice" || peer != "frank" || len(objects) == 0 {
		return []byte{}, "", nil
	}
	// Intent: Alice's storage/forwarding token is reciprocal promise economics
	// carried inside the TCP promise flow, not sideband fixture state or an
	// authorization grant. Source: DI-biruf
	issued, issueErr := economy.IssueBearerToken(runtime.cas, economy.BearerToken{
		Issuer:                 runtime.cfg.Name,
		Scope:                  "poc18-tcp-storage-forwarding",
		ObjectCID:              objects[0].CID,
		Value:                  3,
		Unit:                   "storage_credit",
		ExpiresUnix:            time.Now().Add(1 * time.Hour).Unix(),
		Nonce:                  "poc18-storage-payment:" + peer + ":" + store.CIDText(interestCID),
		Transferable:           true,
		RedeemableCapabilities: []string{"store_forward_object"},
	})
	if issueErr != nil {
		return nil, "", issueErr
	}
	if recordErr := runtime.record("storage_payment_token_issued", "kept", peer, "token_cid="+issued.CID+" object="+objects[0].CID); recordErr != nil {
		return nil, "", recordErr
	}
	return issued.Bytes, issued.CID, nil
}

func (runtime *Runtime) redeemStoragePaymentIfUseful(peer string, offer availabilityOffer) error {
	if runtime.cfg.Name != "frank" || len(offer.PaymentToken) == 0 {
		return nil
	}
	// Intent: Frank decides locally whether Alice's bearer token is worth
	// accepting, then issues Frank's non-transferable service capability as a new
	// local promise. Source: DI-biruf
	if len(offer.Objects) == 0 {
		return fmt.Errorf("storage payment offer has no object")
	}
	ledger := economy.NewLedger()
	report, redeemErr := ledger.RedeemBearerForCapability(runtime.cas, offer.PaymentToken, economy.ExpectedBearerPayment{
		Issuer:     peer,
		Scope:      "poc18-tcp-storage-forwarding",
		ObjectCID:  offer.Objects[0].CID,
		Value:      3,
		Unit:       "storage_credit",
		Capability: "store_forward_object",
	}, runtime.cfg.Name, peer, time.Now())
	if redeemErr != nil {
		return redeemErr
	}
	if recordErr := runtime.record("storage_payment_redemption", "kept", peer, "bearer="+report.BearerTokenCID+" capability="+report.CapabilityTokenCID); recordErr != nil {
		return recordErr
	}
	return runtime.record("storage_capability_token_issued", "kept", peer, "subject="+report.Subject+" capability="+report.Capability)
}

func (runtime *Runtime) handleSyncInterest(conn *transport.Conn, received receivedMessage) {
	peer := received.Peer
	exchangeID, scope, wanted, interestErr := syncInterestFromView(received.View)
	if interestErr != nil {
		runtime.recordBestEffort("sync_interest_parse_failed", "broken", peer, interestErr.Error())
		return
	}
	objects, objectsErr := runtime.objectsForRequest(wanted)
	if objectsErr != nil {
		runtime.recordBestEffort("sync_interest_not_promised", "broken", peer, objectsErr.Error())
		return
	}
	issued, issueErr := economy.IssueRetrievalCapability(runtime.cas, economy.RetrievalCapabilityToken{
		Issuer:       runtime.cfg.Name,
		Subject:      peer,
		Scope:        scope,
		ObjectCIDs:   blockCIDTexts(objects),
		ExpiresUnix:  time.Now().Add(1 * time.Hour).Unix(),
		Nonce:        "poc18:" + runtime.cfg.Name + ":" + peer + ":" + string(exchangeID) + ":" + store.CIDText(received.CID),
		Transferable: false,
	})
	if issueErr != nil {
		runtime.recordBestEffort("retrieval_capability_issue_failed", "broken", peer, issueErr.Error())
		return
	}
	if recordErr := runtime.record("retrieval_capability_issued", "kept", peer, "token_cid="+issued.CID); recordErr != nil {
		runtime.recordBestEffort("retrieval_capability_record_failed", "broken", peer, recordErr.Error())
	}
	paymentToken, paymentCID, paymentErr := runtime.optionalStoragePaymentToken(peer, objects, received.CID)
	if paymentErr != nil {
		runtime.recordBestEffort("storage_payment_issue_failed", "broken", peer, paymentErr.Error())
		return
	}
	body := []any{string(exchangeID), "available", missingRows(objects), issued.Bytes, paymentToken, paymentCID, []any{"redeem this retrieval token over TCP", "CAR bytes are sent only after redemption"}}
	message, messageErr := runtime.newMessage([]graph.Parent{{Role: "responds_to", CID: received.CID}}, graph.Payload{
		Promiser:          runtime.cfg.Name,
		Promisee:          peer,
		PromiseKind:       "object_availability",
		PromiseBody:       body,
		ReciprocalPromise: []any{"redeem retrieval capability before object bytes are sent"},
		LocalConstraints:  []any{"serves only objects in local sparse CAS"},
	})
	if messageErr != nil {
		runtime.recordBestEffort("object_availability_failed", "broken", peer, messageErr.Error())
		return
	}
	if sendErr := runtime.sendMessage(conn, peer, message, "sent"); sendErr != nil {
		runtime.recordBestEffort("object_availability_send_failed", "broken", peer, sendErr.Error())
	}
}

func (runtime *Runtime) handleObjectAvailability(_ *transport.Conn, received receivedMessage) {
	if recordErr := runtime.record("object_availability_recorded", "kept", received.Peer, "stored availability promise"); recordErr != nil {
		runtime.recordBestEffort("object_availability_record_failed", "broken", received.Peer, recordErr.Error())
	}
}

func (runtime *Runtime) handleRetrievalRedemption(conn *transport.Conn, received receivedMessage) {
	peer := received.Peer
	exchangeID, objects, token, parseErr := redemptionObjectsAndToken(received.View)
	if parseErr != nil {
		runtime.recordBestEffort("redemption_parse_failed", "broken", peer, parseErr.Error())
		return
	}
	scope := "latest"
	if len(objects) > 0 && string(exchangeID) != "" && strings.Contains(string(exchangeID), "repair") {
		scope = "objects"
	}
	if _, verifyErr := economy.VerifyRetrievalCapability(token, economy.ExpectedRetrievalCapability{
		Issuer:     runtime.cfg.Name,
		Subject:    peer,
		Scope:      scope,
		ObjectCIDs: blockCIDTexts(objects),
	}, time.Now()); verifyErr != nil {
		runtime.recordBestEffort("redemption_capability_verify_failed", "broken", peer, verifyErr.Error())
		return
	}
	cids := make([]cidlib.Cid, 0, len(objects))
	for _, object := range objects {
		objectCID, cidErr := store.ParseCIDText(object.CID)
		if cidErr != nil {
			runtime.recordBestEffort("redemption_cid_parse_failed", "broken", peer, cidErr.Error())
			return
		}
		cids = append(cids, objectCID)
	}
	carBytes, blocks, carErr := carbundle.Encode(runtime.cas, cids)
	if carErr != nil {
		runtime.recordBestEffort("car_encode_failed", "broken", peer, carErr.Error())
		return
	}
	blockRows := blockRows(blocks)
	body := []any{string(exchangeID), "car_v1", carBytes, blockRows, "I promise these are the original exact bytes I first encountered for each CID."}
	message, messageErr := runtime.newMessage([]graph.Parent{{Role: "responds_to", CID: received.CID}}, graph.Payload{
		Promiser:          runtime.cfg.Name,
		Promisee:          peer,
		PromiseKind:       "object_bytes",
		PromiseBody:       body,
		ReciprocalPromise: []any{"receiver verifies every CID from exact bytes before storing"},
		LocalConstraints:  []any{"CARv1 payload is package format, not authority"},
	})
	if messageErr != nil {
		runtime.recordBestEffort("object_bytes_failed", "broken", peer, messageErr.Error())
		return
	}
	if sendErr := runtime.sendMessage(conn, peer, message, "sent"); sendErr != nil {
		runtime.recordBestEffort("object_bytes_send_failed", "broken", peer, sendErr.Error())
		return
	}
	blockCIDs := make([]string, 0, len(blocks))
	for _, block := range blocks {
		blockCIDs = append(blockCIDs, block.CID)
	}
	artifact := eventstream.CARArtifactFor(runtime.cfg.Name, "sent", peer, store.CIDText(message.CID), carBytes, blockCIDs)
	if emitErr := runtime.emit(eventstream.Record{Kind: eventstream.KindCARArtifact, CARArtifact: &artifact}); emitErr != nil {
		runtime.recordBestEffort("car_artifact_emit_failed", "broken", peer, emitErr.Error())
	}
}

func (runtime *Runtime) handleObjectBytes(received receivedMessage, peer string) error {
	exchangeID, carBytes, kinds, parseErr := objectBytesCAR(received.View)
	if parseErr != nil {
		return parseErr
	}
	blocks, decodeErr := carbundle.DecodeAndStore(runtime.cas, carBytes, kinds)
	if decodeErr != nil {
		return decodeErr
	}
	for _, block := range blocks {
		objectCID, cidErr := store.ParseCIDText(block.CID)
		if cidErr != nil {
			return cidErr
		}
		runtime.setHead("latest", objectCID)
	}
	blockCIDs := make([]string, 0, len(blocks))
	for _, block := range blocks {
		blockCIDs = append(blockCIDs, block.CID)
	}
	messageCID := store.CIDText(received.CID)
	artifact := eventstream.CARArtifactFor(runtime.cfg.Name, "received", peer, messageCID, carBytes, blockCIDs)
	if emitErr := runtime.emit(eventstream.Record{Kind: eventstream.KindCARArtifact, CARArtifact: &artifact}); emitErr != nil {
		return emitErr
	}
	if recordErr := runtime.record("object_bytes_stored", "kept", peer, fmt.Sprintf("exchange=%s blocks=%d", exchangeID, len(blocks))); recordErr != nil {
		return recordErr
	}
	return runtime.repairMissingParents(peer, exchangeID, blocks)
}

func (runtime *Runtime) repairMissingParents(peer string, exchangeID ExchangeID, blocks []carbundle.Block) error {
	// Intent: A sparse CAS may receive a child graph object before its parent
	// links. Repair must therefore be another promise exchange over TCP rather
	// than an in-process shared-CAS shortcut. Source: DI-biruf
	missingParents, missingErr := runtime.missingParentBlocks(blocks)
	if missingErr != nil {
		return missingErr
	}
	if len(missingParents) == 0 {
		return nil
	}
	if recordErr := runtime.record("dag_closure_missing", "kept", peer, fmt.Sprintf("exchange=%s missing=%d", exchangeID, len(missingParents))); recordErr != nil {
		return recordErr
	}
	repairExchange := ExchangeID(runtime.cfg.Name + "-repair-from-" + peer)
	if err := runtime.requestObjects(peer, repairExchange, missingParents); err != nil {
		return err
	}
	return runtime.record("dag_repair_completed", "kept", peer, fmt.Sprintf("exchange=%s requested=%d", repairExchange, len(missingParents)))
}

func (runtime *Runtime) missingParentBlocks(blocks []carbundle.Block) ([]carbundle.Block, error) {
	missing := []carbundle.Block{}
	seen := map[string]bool{}
	for _, block := range blocks {
		if block.Kind != "message" {
			continue
		}
		objectCID, parseErr := store.ParseCIDText(block.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		content, _, getErr := runtime.cas.Get(objectCID)
		if getErr != nil {
			return nil, getErr
		}
		view, viewErr := graph.ParseEnvelope(content)
		if viewErr != nil {
			return nil, viewErr
		}
		for _, parent := range view.Parents {
			parentText := store.CIDText(parent.CID)
			if runtime.cas.Has(parent.CID) || seen[parentText] {
				continue
			}
			seen[parentText] = true
			missing = append(missing, carbundle.Block{CID: parentText, Kind: "message"})
		}
	}
	return missing, nil
}

type outgoingMessage struct {
	Bytes []byte
	Kind  string
	CID   cidlib.Cid
}

func (runtime *Runtime) newMessage(parents []graph.Parent, payload graph.Payload) (outgoingMessage, error) {
	message, messageErr := graph.NewMessage(parents, payload)
	if messageErr != nil {
		return outgoingMessage{}, messageErr
	}
	messageBytes, bytesErr := message.Bytes()
	if bytesErr != nil {
		return outgoingMessage{}, bytesErr
	}
	entry, putErr := runtime.cas.Put("message", messageBytes)
	if putErr != nil {
		return outgoingMessage{}, putErr
	}
	messageCID, parseErr := store.ParseCIDText(entry.CID)
	if parseErr != nil {
		return outgoingMessage{}, parseErr
	}
	return outgoingMessage{Bytes: messageBytes, Kind: payload.PromiseKind, CID: messageCID}, nil
}

func (runtime *Runtime) sendMessage(conn *transport.Conn, peer string, message outgoingMessage, direction string) error {
	if err := conn.WriteFrame(message.Bytes); err != nil {
		return err
	}
	artifact := eventstream.MessageArtifactFor(runtime.cfg.Name, direction, peer, graph.VersionControlPCIDText, message.Kind, message.Bytes)
	if err := runtime.emit(eventstream.Record{Kind: eventstream.KindMessageArtifact, MessageArtifact: &artifact}); err != nil {
		return err
	}
	return runtime.record("message_"+direction, "kept", peer, "kind="+message.Kind+" cid="+artifact.ExactCID)
}

func (runtime *Runtime) receiveMessage(frame []byte, direction string, peer string) (receivedMessage, error) {
	entry, putErr := runtime.cas.Put("message", frame)
	if putErr != nil {
		return receivedMessage{}, putErr
	}
	view, parseErr := graph.ParseEnvelope(frame)
	if parseErr != nil {
		return receivedMessage{}, parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return receivedMessage{}, kindErr
	}
	if peer == "peer" {
		peer = peerFromView(view)
	}
	messageCID, parseErr := store.ParseCIDText(entry.CID)
	if parseErr != nil {
		return receivedMessage{}, parseErr
	}
	artifact := eventstream.MessageArtifactFor(runtime.cfg.Name, direction, peer, graph.VersionControlPCIDText, kind, frame)
	if err := runtime.emit(eventstream.Record{Kind: eventstream.KindMessageArtifact, MessageArtifact: &artifact}); err != nil {
		return receivedMessage{}, err
	}
	received := receivedMessage{View: view, Bytes: append([]byte(nil), frame...), CID: messageCID, Peer: peer, Kind: kind}
	return received, runtime.record("message_"+direction, "kept", peer, "kind="+kind+" cid="+artifact.ExactCID)
}

func (runtime *Runtime) currentHeads() []carbundle.Block {
	runtime.headsMu.Lock()
	defer runtime.headsMu.Unlock()
	roles := make([]string, 0, len(runtime.heads))
	for role := range runtime.heads {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	blocks := make([]carbundle.Block, 0, len(roles))
	for _, role := range roles {
		blocks = append(blocks, carbundle.Block{CID: store.CIDText(runtime.heads[role]), Kind: "message"})
	}
	return blocks
}

func (runtime *Runtime) setHead(role string, objectCID cidlib.Cid) {
	runtime.headsMu.Lock()
	defer runtime.headsMu.Unlock()
	runtime.heads[role] = objectCID
}

func (runtime *Runtime) record(eventName, outcome, peer, detail string) error {
	return runtime.emit(eventstream.Record{Kind: eventstream.KindEvent, Event: &eventstream.Event{
		Observer: runtime.cfg.Name,
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	}})
}

func (runtime *Runtime) recordBestEffort(eventName, outcome, peer, detail string) {
	if recordErr := runtime.record(eventName, outcome, peer, detail); recordErr != nil {
		fmt.Fprintf(os.Stderr, "record %s failed: %v\n", eventName, recordErr)
	}
}

func (runtime *Runtime) emit(record eventstream.Record) error {
	return runtime.collector.Emit(record)
}

func (runtime *Runtime) closeListener() {
	close(runtime.done)
	if runtime.listener != nil {
		if closeErr := runtime.listener.Close(); closeErr != nil {
			runtime.recordBestEffort("listener_close_failed", "broken", "", closeErr.Error())
		}
	}
}

func (runtime *Runtime) closeCollector() {
	if runtime.collector != nil {
		if closeErr := runtime.collector.Close(); closeErr != nil {
			fmt.Printf("collector close failed: %v\n", closeErr)
		}
	}
}

func syncInterestFromView(view graph.EnvelopeView) (ExchangeID, string, []carbundle.Block, error) {
	body, ok := view.Payload[3].([]any)
	if !ok || len(body) < 3 {
		return "", "", nil, fmt.Errorf("sync_interest body must have at least three slots")
	}
	exchangeText, ok := body[0].(string)
	if !ok || strings.TrimSpace(exchangeText) == "" {
		return "", "", nil, fmt.Errorf("sync_interest exchange id must be text")
	}
	scope, ok := body[1].(string)
	if !ok || strings.TrimSpace(scope) == "" {
		return "", "", nil, fmt.Errorf("sync_interest scope must be text")
	}
	objects, objectsErr := blocksFromRows(body[2])
	if objectsErr != nil {
		return "", "", nil, objectsErr
	}
	return ExchangeID(exchangeText), scope, objects, nil
}

func availabilityOfferFromView(view graph.EnvelopeView) (availabilityOffer, error) {
	body, ok := view.Payload[3].([]any)
	if !ok || len(body) < 6 {
		return availabilityOffer{}, fmt.Errorf("object_availability body must have at least six slots")
	}
	exchangeText, ok := body[0].(string)
	if !ok || strings.TrimSpace(exchangeText) == "" {
		return availabilityOffer{}, fmt.Errorf("object_availability exchange id must be text")
	}
	status, ok := body[1].(string)
	if !ok || status != "available" {
		return availabilityOffer{}, fmt.Errorf("object_availability status must be available")
	}
	objects, objectsErr := blocksFromRows(body[2])
	if objectsErr != nil {
		return availabilityOffer{}, objectsErr
	}
	token, ok := body[3].([]byte)
	if !ok || len(token) == 0 {
		return availabilityOffer{}, fmt.Errorf("object_availability token must be bytes")
	}
	paymentToken, ok := body[4].([]byte)
	if !ok {
		return availabilityOffer{}, fmt.Errorf("object_availability payment token must be bytes")
	}
	paymentCID, ok := body[5].(string)
	if !ok {
		return availabilityOffer{}, fmt.Errorf("object_availability payment token CID must be text")
	}
	if len(paymentToken) > 0 && paymentCID == "" {
		return availabilityOffer{}, fmt.Errorf("object_availability payment token CID is required when token bytes are present")
	}
	return availabilityOffer{ExchangeID: ExchangeID(exchangeText), Objects: objects, RetrievalToken: token, PaymentToken: paymentToken, PaymentTokenCID: paymentCID}, nil
}

func redemptionObjectsAndToken(view graph.EnvelopeView) (ExchangeID, []carbundle.Block, []byte, error) {
	body, ok := view.Payload[3].([]any)
	if !ok || len(body) < 3 {
		return "", nil, nil, fmt.Errorf("object_retrieval_redemption body must have at least three slots")
	}
	exchangeText, ok := body[0].(string)
	if !ok || strings.TrimSpace(exchangeText) == "" {
		return "", nil, nil, fmt.Errorf("object_retrieval_redemption exchange id must be text")
	}
	token, ok := body[1].([]byte)
	if !ok || len(token) == 0 {
		return "", nil, nil, fmt.Errorf("object_retrieval_redemption token must be bytes")
	}
	objects, objectsErr := blocksFromRows(body[2])
	if objectsErr != nil {
		return "", nil, nil, objectsErr
	}
	return ExchangeID(exchangeText), objects, token, nil
}

func objectBytesCAR(view graph.EnvelopeView) (ExchangeID, []byte, map[string]string, error) {
	body, ok := view.Payload[3].([]any)
	if !ok || len(body) < 4 {
		return "", nil, nil, fmt.Errorf("object_bytes body must have at least four slots")
	}
	exchangeText, ok := body[0].(string)
	if !ok || strings.TrimSpace(exchangeText) == "" {
		return "", nil, nil, fmt.Errorf("object_bytes exchange id must be text")
	}
	profile, ok := body[1].(string)
	if !ok || profile != "car_v1" {
		return "", nil, nil, fmt.Errorf("object_bytes must carry car_v1")
	}
	carBytes, ok := body[2].([]byte)
	if !ok || len(carBytes) == 0 {
		return "", nil, nil, fmt.Errorf("object_bytes CAR payload must be bytes")
	}
	blocks, blocksErr := blocksFromRows(body[3])
	if blocksErr != nil {
		return "", nil, nil, blocksErr
	}
	kinds := map[string]string{}
	for _, block := range blocks {
		kinds[block.CID] = block.Kind
	}
	return ExchangeID(exchangeText), carBytes, kinds, nil
}

func missingRows(blocks []carbundle.Block) []any {
	rows := make([]any, 0, len(blocks))
	for _, block := range blocks {
		objectCID, parseErr := store.ParseCIDText(block.CID)
		if parseErr != nil {
			continue
		}
		rows = append(rows, []any{block.Kind, store.LinkTag(objectCID)})
	}
	return rows
}

func blockRows(blocks []carbundle.Block) []any {
	rows := make([]any, 0, len(blocks))
	for _, block := range blocks {
		objectCID, parseErr := store.ParseCIDText(block.CID)
		if parseErr != nil {
			continue
		}
		rows = append(rows, []any{block.Kind, store.LinkTag(objectCID), block.Size})
	}
	return rows
}

func blockCIDTexts(blocks []carbundle.Block) []string {
	values := make([]string, 0, len(blocks))
	for _, block := range blocks {
		values = append(values, block.CID)
	}
	sort.Strings(values)
	return values
}

func blocksFromRows(value any) ([]carbundle.Block, error) {
	rows, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("block rows must be array")
	}
	blocks := make([]carbundle.Block, 0, len(rows))
	for _, rowValue := range rows {
		row, ok := rowValue.([]any)
		if !ok || len(row) < 2 {
			return nil, fmt.Errorf("block row must have at least kind and CID")
		}
		kind, ok := row[0].(string)
		if !ok || strings.TrimSpace(kind) == "" {
			return nil, fmt.Errorf("block row kind must be text")
		}
		objectCID, cidErr := store.CIDFromLinkTag(row[1])
		if cidErr != nil {
			return nil, cidErr
		}
		blocks = append(blocks, carbundle.Block{CID: store.CIDText(objectCID), Kind: kind})
	}
	return blocks, nil
}

func peerFromView(view graph.EnvelopeView) string {
	if len(view.Payload) > 0 {
		return fmt.Sprint(view.Payload[0])
	}
	return ""
}

func ParsePeers(values []string) (map[string]string, error) {
	peers := map[string]string{}
	for _, value := range values {
		name, address, found := strings.Cut(value, "=")
		if !found || name == "" || address == "" {
			return nil, fmt.Errorf("peer must be name=address")
		}
		peers[name] = address
	}
	return peers, nil
}

// PeerFlags implements flag.Value for repeated peer arguments.
type PeerFlags []string

func (flags *PeerFlags) String() string {
	return strings.Join(*flags, ",")
}

func (flags *PeerFlags) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}
