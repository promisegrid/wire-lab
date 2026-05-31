package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/compute"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/data"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/policy"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/protocol"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/storage"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/token"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/transport"
)

const (
	economyAppName             = "promise-economy"
	kindNeedAdvertisement      = "need_advertisement"
	kindOfferPromise           = "offer_promise"
	kindCounterPromise         = "counter_promise"
	kindAcceptancePromise      = "acceptance_promise"
	kindTokenIssuePromise      = "token_issue_promise"
	kindTokenRedemptionPromise = "token_redemption_promise"
	kindOutcomeObservation     = "outcome_observation"
	kindExchangeRateQuote      = "exchange_rate_quote"
)

// WireMessage is the in-process view of one signed POC8 app envelope. Intent:
// Keep the Go struct as local scaffolding while the protocol object exchanged
// across peers remains the exact pCID-selected CBOR grid bytes. Source: DI-sirus
//
// One pCID names the whole POC8 protocol; Kind is a payload variant under that
// protocol, not a separate pCID or hidden RPC method name.
type WireMessage struct {
	FromNode string            `json:"from_node"`
	FromApp  string            `json:"from_app"`
	ToNode   string            `json:"to_node"`
	ToApp    string            `json:"to_app"`
	Kind     string            `json:"kind"`
	Route    []string          `json:"route"`
	Payload  map[string]string `json:"payload"`
	Token    token.Token       `json:"token"`
	Envelope string            `json:"envelope"`
}

// WireResponse is bounded run-control feedback for this POC. It is not a
// PromiseGrid authority; receivers still record their own local observations
// and trust updates from the signed promise they presented. Source: DI-sirus
type WireResponse struct {
	Outcome string            `json:"outcome"`
	Detail  string            `json:"detail"`
	Token   token.Token       `json:"token"`
	Payload map[string]string `json:"payload"`
}

// EnvelopeBytes returns the exact signed CBOR grid bytes carried by the local
// message view. Source: DI-sirus
func (message WireMessage) EnvelopeBytes() ([]byte, error) {
	if message.Envelope == "" {
		return nil, fmt.Errorf("wire message missing signed envelope")
	}
	return hex.DecodeString(message.Envelope)
}

// appProtocolCID names one POC8 promise-economy protocol. Intent: Keep all POC8
// message variants under a single protocol CID so `need_advertisement`,
// `offer_promise`, and the rest are payload kinds defined by this one spec,
// not separate pCIDs. Source: DI-sirus
var appProtocolCID = protocol.NewProtocolCID([]byte("poc8 promise-economy protocol v1: one pCID for autonomous need advertisements, offer promises, counter promises, acceptances, token issue promises, token redemption promises, exchange-rate quotes, and outcome observations"))

// Node groups one local kernel boundary plus the promise-economy app roles in a
// container. Intent: Each process acts from its own local wants, offers, wallet,
// issuer state, resource state, and trust evidence instead of following Alice's
// central transaction script. Source: DI-sirus
type Node struct {
	Name         string
	RunID        string
	RunRoot      string
	Issuer       *token.Issuer
	Wallet       *token.Wallet
	Policy       policy.AgentPolicy
	Evidence     []token.Event
	Storage      map[string]string
	Data         map[string]string
	TokenSources map[string]string
}

var nodeAddrs = map[string]string{
	"alice":   "alice:8077",
	"bob":     "bob:8077",
	"carol":   "carol:8077",
	"dave":    "dave:8077",
	"mallory": "mallory:8077",
}

func main() {
	nodeName := flag.String("node", "", "node name")
	flag.Parse()
	if *nodeName == "" {
		log.Fatal("missing -node")
	}
	node := newNode(*nodeName)
	listener, listenErr := net.Listen("tcp", ":8077")
	if listenErr != nil {
		log.Fatalf("%s listen failed: %v", node.Name, listenErr)
	}
	serverErrors := make(chan error, 1)
	go node.serveTCP(listener, serverErrors)
	go node.runAutonomousPlan()
	if err := node.waitForDone(); err != nil {
		log.Printf("%s wait failed: %v", node.Name, err)
		os.Exit(1)
	}
	if closeErr := listener.Close(); closeErr != nil {
		log.Printf("%s listener close failed: %v", node.Name, closeErr)
	}
	select {
	case serverErr := <-serverErrors:
		if serverErr != nil && !errors.Is(serverErr, net.ErrClosed) {
			log.Printf("%s tcp server failed: %v", node.Name, serverErr)
			os.Exit(1)
		}
	default:
	}
}

// newNode seeds only local facts for one agent. Mallory starts with stale Alice
// bearer-token bytes as historical local evidence; Alice independently starts
// with issuer-local revoked state for those same promise IDs. Source: DI-sirus
func newNode(nodeName string) *Node {
	node := &Node{
		Name:         nodeName,
		RunID:        getenv("POC8_RUN_ID", "manual"),
		RunRoot:      "/run/poc8",
		Issuer:       token.NewIssuer(nodeName),
		Wallet:       token.NewWallet(nodeName),
		Policy:       policy.ForNode(nodeName),
		Storage:      make(map[string]string),
		Data:         initialData(nodeName),
		TokenSources: make(map[string]string),
	}
	if nodeName == "alice" {
		for _, tokenID := range []string{"alice-stale-bearer-1", "alice-stale-bearer-2"} {
			issuedToken, issueErr := node.Issuer.Issue(tokenID, "mallory", data.Kind, "dataset-stale", token.TransferBearer)
			if issueErr != nil {
				log.Fatalf("seed stale token: %v", issueErr)
			}
			if revokeErr := node.Issuer.Revoke(issuedToken.ID, "Alice's local promise history marks this bearer promise stale before POC8 starts"); revokeErr != nil {
				log.Fatalf("seed stale revocation: %v", revokeErr)
			}
		}
	}
	if nodeName == "mallory" {
		for _, tokenID := range []string{"alice-stale-bearer-1", "alice-stale-bearer-2"} {
			staleToken, staleErr := historicalAliceStaleToken(tokenID)
			if staleErr != nil {
				log.Fatalf("seed mallory stale token: %v", staleErr)
			}
			node.Wallet.Add(staleToken, "historical local holding of Alice bearer promise")
			node.TokenSources[staleToken.ID] = "historical"
		}
	}
	return node
}

// serveTCP accepts bounded length-framed TCP messages from neighbor nodes.
// Intent: Keep transport as ordinary plumbing while the pCID-selected grid
// envelope remains the app/kernel and peer/kernel protocol object. Source:
// DI-sirus
func (node *Node) serveTCP(listener net.Listener, serverErrors chan<- error) {
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErrors <- acceptErr
			return
		}
		go node.handleTCPConnection(conn)
	}
}

// handleTCPConnection reads one signed envelope and returns a bounded local
// observation response. The response helps the demo progress; the promise and
// proof remain inside the signed grid envelope. Source: DI-sirus
func (node *Node) handleTCPConnection(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	defer closeFrame(frameConn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		node.record("tcp_frame_read", token.OutcomeBroken, "", readErr.Error())
		return
	}
	result, receiveErr := node.receiveFrame(frameBytes)
	if receiveErr != nil {
		result = WireResponse{Outcome: token.OutcomeBroken, Detail: receiveErr.Error()}
	}
	responseBytes, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		node.record("tcp_response_encode", token.OutcomeBroken, "", marshalErr.Error())
		return
	}
	if writeErr := frameConn.WriteFrame(responseBytes); writeErr != nil {
		node.record("tcp_response_write", token.OutcomeBroken, "", writeErr.Error())
	}
}

func (node *Node) receiveFrame(envelopeBytes []byte) (WireResponse, error) {
	message, exactBytes, decodeErr := node.decodeEnvelopeBytes(envelopeBytes, WireMessage{})
	if decodeErr != nil {
		return WireResponse{}, decodeErr
	}
	node.record("kernel_receive", token.OutcomeKept, message.Token.ID, node.Name+" kernel observed "+message.Kind+" exact_sha256="+protocol.HashExactBytes(exactBytes))
	if message.ToNode != node.Name {
		return node.forward(message)
	}
	if message.ToApp != economyAppName {
		return WireResponse{}, fmt.Errorf("%s has no app %s", node.Name, message.ToApp)
	}
	return node.handleEconomy(message)
}

func (node *Node) forward(message WireMessage) (WireResponse, error) {
	nextHop, nextErr := node.nextRouteHop(message)
	if nextErr != nil {
		return WireResponse{}, nextErr
	}
	node.record("relay_forward", token.OutcomeKept, message.Token.ID, node.Name+" relay app promised next hop "+nextHop)
	return transmit(nodeAddrs[nextHop], message)
}

// nextRouteHop finds the next peer after this node in the signed route. Each
// relay promises only its own next hop; no relay rewrites the app promise.
// Source: DI-sirus
func (node *Node) nextRouteHop(message WireMessage) (string, error) {
	if len(message.Route) == 0 {
		return "", fmt.Errorf("no route from %s to %s", node.Name, message.ToNode)
	}
	for routeIndex, routeNode := range message.Route {
		if routeNode == node.Name && routeIndex+1 < len(message.Route) {
			return message.Route[routeIndex+1], nil
		}
	}
	return "", fmt.Errorf("route %s has no promised next hop after %s toward %s", strings.Join(message.Route, ","), node.Name, message.ToNode)
}

func (node *Node) handleEconomy(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case kindNeedAdvertisement:
		return node.handleNeedAdvertisement(message)
	case kindOfferPromise:
		return node.handleOfferPromise(message)
	case kindCounterPromise:
		return node.handleCounterPromise(message)
	case kindAcceptancePromise:
		return node.handleAcceptancePromise(message)
	case kindTokenIssuePromise:
		return node.handleTokenIssuePromise(message)
	case kindTokenRedemptionPromise:
		return node.handleTokenRedemptionPromise(message)
	case kindOutcomeObservation:
		node.record("outcome_observed", message.Payload["outcome"], message.Token.ID, node.Name+" observed "+message.Payload["detail"])
		return WireResponse{Outcome: token.OutcomeKept, Detail: "outcome recorded"}, nil
	case kindExchangeRateQuote:
		return node.handleExchangeRateQuote(message)
	default:
		return WireResponse{}, fmt.Errorf("promise economy app has no local handler for %s", message.Kind)
	}
}

// runAutonomousPlan gives each node its own local clocked intentions. Intent:
// Startup still bounds the demo, but no Alice-owned script commands Bob, Carol,
// Dave, or Mallory; every inter-peer action after startup is a signed promise
// message from the acting agent. Source: DI-sirus
func (node *Node) runAutonomousPlan() {
	switch node.Name {
	case "alice":
		node.runAlicePlan()
	case "dave":
		node.runDavePlan()
	case "mallory":
		node.runMalloryPlan()
	default:
		node.record("local_plan_ready", token.OutcomeKept, "", node.Name+" waits for peer promises matching local resources")
	}
}

func (node *Node) runAlicePlan() {
	time.Sleep(1500 * time.Millisecond)
	node.advertiseNeed("bob", storage.Kind, "alice-report", map[string]string{
		"need_id":                 "alice-storage-need-1",
		"alice_promises":          "need is current; Alice will not send data until a voluntary storage promise is accepted; Alice will record keep/break evidence locally",
		"requested_peer_promises": "store key alice-report and return the value later",
		"max_price":               "8",
		"stake_requested":         "2",
		"value":                   "private report from Alice",
	})
	time.Sleep(500 * time.Millisecond)
	node.advertiseNeed("carol", compute.Kind, "fib-10", map[string]string{
		"need_id":                 "alice-compute-need-1",
		"alice_promises":          "need is current; Alice will consider compute offers until local expiry; Alice will record the returned result as local evidence",
		"requested_peer_promises": "compute Fibonacci 10 and return the result",
		"max_price":               "6",
		"stake_requested":         "1",
		"n":                       "10",
	})
	time.Sleep(18 * time.Second)
	if doneErr := node.writeDone(); doneErr != nil {
		node.record("demo_done_write", token.OutcomeBroken, "", doneErr.Error())
	}
}

func (node *Node) runDavePlan() {
	time.Sleep(5 * time.Second)
	if _, quoteErr := node.send("bob", kindExchangeRateQuote, map[string]string{"offered_issuer": "alice", "wanted_issuer": "bob"}, token.Token{}); quoteErr != nil {
		node.record("exchange_quote_request", token.OutcomeBroken, "", quoteErr.Error())
	}
}

func (node *Node) runMalloryPlan() {
	time.Sleep(3500 * time.Millisecond)
	firstStale, firstErr := historicalAliceStaleToken("alice-stale-bearer-1")
	if firstErr != nil {
		node.record("stale_offer_prepare", token.OutcomeBroken, "alice-stale-bearer-1", firstErr.Error())
		return
	}
	node.offerTokenToDave(firstStale, "mallory-stale-offer-1")
	time.Sleep(7 * time.Second)
	secondStale, secondErr := historicalAliceStaleToken("alice-stale-bearer-2")
	if secondErr != nil {
		node.record("stale_offer_prepare", token.OutcomeBroken, "alice-stale-bearer-2", secondErr.Error())
		return
	}
	node.offerTokenToDave(secondStale, "mallory-stale-offer-2")
}

// advertiseNeed publishes Alice's own promises and requested reciprocal promise
// shape. It does not command the peer; the peer may ignore, refuse, offer, or
// counter according to its local policy. Source: DI-sirus
func (node *Node) advertiseNeed(peer string, resourceKind string, resourceID string, payload map[string]string) {
	needPayload := copyPayload(payload)
	needPayload["resource_kind"] = resourceKind
	needPayload["resource_id"] = resourceID
	needPayload["promiser"] = node.Name
	needPayload["expires_after"] = "bounded-poc-window"
	node.record("need_advertised", token.OutcomeKept, "", node.Name+" promised local need "+needPayload["need_id"]+" to "+peer)
	if _, sendErr := node.send(peer, kindNeedAdvertisement, needPayload, token.Token{}); sendErr != nil {
		node.record("need_advertise_send", token.OutcomeBroken, "", sendErr.Error())
	}
}

func (node *Node) offerTokenToDave(offeredToken token.Token, offerID string) {
	payload := map[string]string{
		"offer_id":       offerID,
		"resource_kind":  offeredToken.ResourceKind,
		"resource_id":    offeredToken.ResourceID,
		"price":          "1",
		"stake":          "0",
		"offer_promises": "Mallory promises only that she is voluntarily presenting these token bytes; Dave must judge Alice and Mallory locally",
		"issuer":         offeredToken.Issuer,
	}
	if _, offerErr := node.send("dave", kindOfferPromise, payload, offeredToken); offerErr != nil {
		node.record("stale_offer_send", token.OutcomeBroken, offeredToken.ID, offerErr.Error())
	}
}

func (node *Node) handleNeedAdvertisement(message WireMessage) (WireResponse, error) {
	resourceKind := message.Payload["resource_kind"]
	if !node.provides(resourceKind) {
		node.record("need_not_matched", token.OutcomeRefused, "", node.Name+" has no local promise offer for "+resourceKind)
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "resource not offered locally"}, nil
	}
	decision := node.decide(policy.ActionOffer, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	offerPayload := map[string]string{
		"need_id":        message.Payload["need_id"],
		"offer_id":       node.Name + "-offer-for-" + message.Payload["need_id"],
		"resource_kind":  resourceKind,
		"resource_id":    message.Payload["resource_id"],
		"price":          node.initialPrice(resourceKind),
		"max_price":      message.Payload["max_price"],
		"stake":          node.initialStake(resourceKind),
		"offer_promises": node.Name + " promises to consider the reciprocal terms and to issue local resource tokens if voluntary agreement is reached",
		"n":              message.Payload["n"],
		"value":          message.Payload["value"],
	}
	if _, sendErr := node.send(message.FromNode, kindOfferPromise, offerPayload, token.Token{}); sendErr != nil {
		return WireResponse{}, sendErr
	}
	node.record("offer_promised", token.OutcomeKept, "", node.Name+" answered need "+message.Payload["need_id"]+" with offer "+offerPayload["offer_id"])
	return WireResponse{Outcome: token.OutcomeKept, Detail: "offer promised"}, nil
}

func (node *Node) handleOfferPromise(message WireMessage) (WireResponse, error) {
	if node.Name == "alice" && message.Payload["resource_kind"] == storage.Kind && atoi(message.Payload["price"]) > atoi(message.Payload["max_price"]) {
		counterPayload := copyPayload(message.Payload)
		counterPayload["counter_id"] = "alice-counter-for-" + message.Payload["offer_id"]
		counterPayload["price"] = "8"
		counterPayload["stake"] = "2"
		counterPayload["counter_promises"] = "Alice promises to send the storage value only after Bob voluntarily accepts price 8 and stake 2"
		decision := node.decide(policy.ActionCounter, WireMessage{FromNode: message.FromNode, Payload: counterPayload})
		if !decision.Accepted {
			return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
		}
		if _, sendErr := node.send(message.FromNode, kindCounterPromise, counterPayload, token.Token{}); sendErr != nil {
			return WireResponse{}, sendErr
		}
		node.record("counter_promised", token.OutcomeKept, "", "Alice countered "+message.Payload["offer_id"]+" with lower price and larger stake promise")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "counter promised"}, nil
	}
	if node.Name == "alice" && message.Payload["resource_kind"] == compute.Kind {
		decision := node.decide(policy.ActionAccept, message)
		if !decision.Accepted {
			return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
		}
		acceptPayload := copyPayload(message.Payload)
		acceptPayload["acceptance_promises"] = "Alice promises to present a reciprocal token or local payment if Carol issues compute access"
		if _, sendErr := node.send(message.FromNode, kindAcceptancePromise, acceptPayload, token.Token{}); sendErr != nil {
			return WireResponse{}, sendErr
		}
		return WireResponse{Outcome: token.OutcomeKept, Detail: "compute offer accepted"}, nil
	}
	if node.Name == "dave" && !message.Token.IsZero() {
		return node.handleIncomingBearerOffer(message)
	}
	if node.Name == "carol" && !message.Token.IsZero() {
		return node.handleBearerForAccessOffer(message)
	}
	return WireResponse{Outcome: token.OutcomeRefused, Detail: node.Name + " did not match offer locally"}, nil
}

func (node *Node) handleIncomingBearerOffer(message WireMessage) (WireResponse, error) {
	quote := node.Wallet.Quote(message.Token.Issuer, node.Name)
	node.record("exchange_rate_quoted", token.OutcomeKept, message.Token.ID, fmt.Sprintf("%s local quote %d %s bearer for %d %s token", node.Name, quote.OfferedCount, quote.OfferedIssuer, quote.WantedCount, quote.WantedIssuer))
	decision := node.decide(policy.ActionAccept, message)
	if !decision.Accepted {
		node.record("bearer_offer_refused", token.OutcomeRefused, message.Token.ID, decision.Detail())
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	node.Wallet.Add(message.Token, message.FromNode+" offered bearer token")
	node.TokenSources[message.Token.ID] = message.FromNode
	node.record("bearer_offer_accepted", token.OutcomeKept, message.Token.ID, node.Name+" accepted bearer token for local evidence")
	response, redeemErr := node.redeemHeldToken(message.Token, map[string]string{"reason": "Dave asks Alice to keep the bearer promise Mallory circulated"})
	if redeemErr != nil {
		return WireResponse{}, redeemErr
	}
	return response, nil
}

func (node *Node) handleBearerForAccessOffer(message WireMessage) (WireResponse, error) {
	decision := node.decide(policy.ActionAccept, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	node.Wallet.Add(message.Token, message.FromNode+" offered bearer stake token for compute access")
	node.TokenSources[message.Token.ID] = message.FromNode
	issuedToken, issueErr := node.Issuer.Issue("carol-compute-alice-nontransfer-1", message.FromNode, compute.Kind, "fib-10", token.TransferNonTransferable)
	if issueErr != nil {
		return WireResponse{}, issueErr
	}
	payload := map[string]string{"offer_id": message.Payload["offer_id"], "accepted_for": message.Token.ID, "n": "10"}
	if _, sendErr := node.send(message.FromNode, kindTokenIssuePromise, payload, issuedToken); sendErr != nil {
		return WireResponse{}, sendErr
	}
	node.record("bearer_for_nontransferable_trade", token.OutcomeKept, issuedToken.ID, "Carol accepted Bob bearer stake token and issued Alice non-transferable compute access")
	return WireResponse{Outcome: token.OutcomeKept, Detail: "bearer-for-access accepted", Token: issuedToken}, nil
}

func (node *Node) handleCounterPromise(message WireMessage) (WireResponse, error) {
	decision := node.decide(policy.ActionAccept, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	if _, sendErr := node.send(message.FromNode, kindAcceptancePromise, copyPayload(message.Payload), token.Token{}); sendErr != nil {
		return WireResponse{}, sendErr
	}
	return node.issueStorageTokens(message)
}

func (node *Node) handleAcceptancePromise(message WireMessage) (WireResponse, error) {
	decision := node.decide(policy.ActionAccept, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	if node.Name == "carol" && message.Payload["resource_kind"] == compute.Kind {
		issuedToken, issueErr := node.Issuer.Issue("carol-compute-alice-1", message.FromNode, compute.Kind, "fib-10", token.TransferNonTransferable)
		if issueErr != nil {
			return WireResponse{}, issueErr
		}
		if _, sendErr := node.send(message.FromNode, kindTokenIssuePromise, map[string]string{"n": "10", "need_id": message.Payload["need_id"]}, issuedToken); sendErr != nil {
			return WireResponse{}, sendErr
		}
		node.record("compute_token_issued", token.OutcomeKept, issuedToken.ID, "Carol issued compute token after voluntary acceptance")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "compute token issued", Token: issuedToken}, nil
	}
	return WireResponse{Outcome: token.OutcomeKept, Detail: "acceptance recorded"}, nil
}

func (node *Node) issueStorageTokens(message WireMessage) (WireResponse, error) {
	issuedTokens := []struct {
		id           string
		resourceID   string
		transferRule string
	}{
		{id: "bob-storage-alice-store-1", resourceID: "alice-report-store", transferRule: token.TransferNonTransferable},
		{id: "bob-storage-alice-read-1", resourceID: "alice-report-read", transferRule: token.TransferNonTransferable},
		{id: "bob-stake-alice-1", resourceID: "bob-storage-stake", transferRule: token.TransferBearer},
	}
	for _, tokenSpec := range issuedTokens {
		issuedToken, issueErr := node.Issuer.Issue(tokenSpec.id, message.FromNode, storage.Kind, tokenSpec.resourceID, tokenSpec.transferRule)
		if issueErr != nil {
			return WireResponse{}, issueErr
		}
		payload := map[string]string{"need_id": message.Payload["need_id"], "stake": message.Payload["stake"], "price": message.Payload["price"], "max_price": message.Payload["max_price"]}
		if _, sendErr := node.send(message.FromNode, kindTokenIssuePromise, payload, issuedToken); sendErr != nil {
			return WireResponse{}, sendErr
		}
		node.record("storage_token_issued", token.OutcomeKept, issuedToken.ID, "Bob issued "+issuedToken.TransferRule+" storage/stake promise token")
	}
	return WireResponse{Outcome: token.OutcomeKept, Detail: "storage and stake tokens issued"}, nil
}

func (node *Node) handleTokenIssuePromise(message WireMessage) (WireResponse, error) {
	decision := node.decide(policy.ActionAccept, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	node.Wallet.Add(message.Token, message.FromNode+" issued token promise")
	node.TokenSources[message.Token.ID] = message.FromNode
	node.record("token_issue_accepted", token.OutcomeKept, message.Token.ID, node.Name+" accepted issuer promise token from "+message.FromNode)
	if node.Name == "alice" {
		return node.actOnAliceToken(message.Token)
	}
	return WireResponse{Outcome: token.OutcomeKept, Detail: "token accepted"}, nil
}

func (node *Node) actOnAliceToken(issuedToken token.Token) (WireResponse, error) {
	switch issuedToken.ResourceID {
	case "alice-report-store":
		return node.redeemHeldToken(issuedToken, map[string]string{"op": "store", "key": "alice-report", "value": "private report from Alice"})
	case "alice-report-read":
		time.Sleep(750 * time.Millisecond)
		return node.redeemHeldToken(issuedToken, map[string]string{"op": "read", "key": "alice-report"})
	case "bob-storage-stake":
		payload := map[string]string{
			"offer_id":       "alice-offers-bob-stake-for-carol-compute",
			"resource_kind":  compute.Kind,
			"resource_id":    "fib-10",
			"price":          "1",
			"stake":          "0",
			"offer_promises": "Alice promises to transfer Bob's bearer stake token if Carol voluntarily issues non-transferable compute access",
		}
		if _, sendErr := node.send("carol", kindOfferPromise, payload, issuedToken); sendErr != nil {
			return WireResponse{}, sendErr
		}
		return WireResponse{Outcome: token.OutcomeKept, Detail: "stake token offered for compute access"}, nil
	case "fib-10":
		return node.redeemHeldToken(issuedToken, map[string]string{"n": "10"})
	default:
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token retained for later local choice"}, nil
	}
}

func (node *Node) handleTokenRedemptionPromise(message WireMessage) (WireResponse, error) {
	event := node.Issuer.Redeem(message.FromNode, message.Token)
	node.record("token_redeemed", event.Outcome, message.Token.ID, event.Detail)
	payload := map[string]string{}
	if event.Outcome == token.OutcomeKept {
		resourcePayload, resourceErr := node.fulfillResource(message)
		if resourceErr != nil {
			node.record("resource_fulfillment", token.OutcomeBroken, message.Token.ID, resourceErr.Error())
			return WireResponse{Outcome: token.OutcomeBroken, Detail: resourceErr.Error()}, nil
		}
		payload = resourcePayload
		payloadBytes, marshalErr := json.Marshal(resourcePayload)
		if marshalErr != nil {
			node.record("resource_fulfillment", token.OutcomeBroken, message.Token.ID, marshalErr.Error())
			return WireResponse{Outcome: token.OutcomeBroken, Detail: marshalErr.Error()}, nil
		}
		node.record("resource_fulfilled", token.OutcomeKept, message.Token.ID, string(payloadBytes))
	}
	return WireResponse{Outcome: event.Outcome, Detail: event.Detail, Payload: payload}, nil
}

func (node *Node) handleExchangeRateQuote(message WireMessage) (WireResponse, error) {
	quote := node.Wallet.Quote(message.Payload["offered_issuer"], message.Payload["wanted_issuer"])
	payload := map[string]string{"offered_count": fmt.Sprintf("%d", quote.OfferedCount), "wanted_count": fmt.Sprintf("%d", quote.WantedCount)}
	node.record("exchange_rate_quoted", token.OutcomeKept, "", node.Name+" quoted a peer-local exchange rate without a central market")
	return WireResponse{Outcome: token.OutcomeKept, Detail: "peer-local quote", Payload: payload}, nil
}

func (node *Node) redeemHeldToken(heldToken token.Token, payload map[string]string) (WireResponse, error) {
	if !node.Wallet.Holds(heldToken.ID) {
		return WireResponse{}, fmt.Errorf("%s has no local held-promise evidence for token %s", node.Name, heldToken.ID)
	}
	message := WireMessage{FromNode: heldToken.Issuer, Payload: map[string]string{"resource_kind": heldToken.ResourceKind, "resource_id": heldToken.ResourceID}, Token: heldToken}
	decision := node.decide(policy.ActionRedeem, message)
	if !decision.Accepted {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: decision.Detail()}, nil
	}
	response, redeemErr := node.send(heldToken.Issuer, kindTokenRedemptionPromise, payload, heldToken)
	if redeemErr != nil {
		return WireResponse{}, redeemErr
	}
	node.applyRedemptionEvidence(heldToken, response)
	observationPayload := map[string]string{"outcome": response.Outcome, "detail": response.Detail}
	if _, observeErr := node.send(heldToken.Issuer, kindOutcomeObservation, observationPayload, heldToken); observeErr != nil {
		node.record("outcome_observation_send", token.OutcomeBroken, heldToken.ID, observeErr.Error())
	}
	return response, nil
}

// decide records one local economic judgment before an agent acts. Intent: POC8
// opportunities may arrive from peers, but the action is still voluntary and
// judged against local trust, cost, scarcity, stake, and utility. Source:
// DI-sirus
func (node *Node) decide(action string, message WireMessage) policy.Decision {
	decision := node.Policy.Decide(node.policyContext(action, message))
	outcome := token.OutcomeRefused
	if decision.Accepted {
		outcome = token.OutcomeKept
	}
	node.record("local_decision", outcome, message.Token.ID, decision.Detail())
	return decision
}

func (node *Node) policyContext(action string, message WireMessage) policy.ActionContext {
	payload := message.Payload
	if payload == nil {
		payload = map[string]string{}
	}
	resourceKind := payload["resource_kind"]
	if resourceKind == "" {
		resourceKind = message.Token.ResourceKind
	}
	resourceID := payload["resource_id"]
	if resourceID == "" {
		resourceID = message.Token.ResourceID
	}
	issuer := message.Token.Issuer
	if issuer == "" {
		issuer = payload["issuer"]
	}
	if issuer == "" {
		issuer = message.FromNode
	}
	context := policy.ActionContext{
		Agent:          node.Name,
		Action:         action,
		Peer:           message.FromNode,
		SourcePeer:     node.TokenSources[message.Token.ID],
		Issuer:         issuer,
		ResourceKind:   resourceKind,
		ResourceID:     resourceID,
		TransferRule:   message.Token.TransferRule,
		Price:          atoi(payload["price"]),
		MaxPrice:       atoi(payload["max_price"]),
		Stake:          atoi(payload["stake"]),
		Scarcity:       node.scarcity(resourceKind),
		ExpectedRisk:   1,
		ExpectedCost:   1,
		ResourceValue:  atoi(payload["value_score"]),
		Token:          message.Token,
		IssuerTrust:    node.trust(issuer),
		PeerTrust:      node.trust(message.FromNode),
		SourceTrust:    node.trust(node.TokenSources[message.Token.ID]),
		RecipientTrust: node.trust(payload["recipient_node"]),
	}
	if context.SourcePeer == "" {
		context.SourcePeer = message.FromNode
		context.SourceTrust = node.trust(message.FromNode)
	}
	if action == policy.ActionAccept && node.provides(context.ResourceKind) {
		context.ExpectedReturn = context.Price + context.Stake
	}
	return context
}

func (node *Node) trust(peer string) int {
	if peer == "" || peer == "historical" {
		return 0
	}
	return node.Wallet.Trust(peer)
}

func (node *Node) scarcity(resourceKind string) int {
	switch node.Name + ":" + resourceKind {
	case "bob:" + storage.Kind:
		return 1
	case "carol:" + compute.Kind:
		return 2
	default:
		return 1
	}
}

func (node *Node) provides(resourceKind string) bool {
	return (node.Name == "bob" && resourceKind == storage.Kind) || (node.Name == "carol" && resourceKind == compute.Kind)
}

func (node *Node) initialPrice(resourceKind string) string {
	if node.Name == "bob" && resourceKind == storage.Kind {
		return "10"
	}
	if node.Name == "carol" && resourceKind == compute.Kind {
		return "6"
	}
	return "1"
}

func (node *Node) initialStake(resourceKind string) string {
	if node.Name == "bob" && resourceKind == storage.Kind {
		return "1"
	}
	if node.Name == "carol" && resourceKind == compute.Kind {
		return "1"
	}
	return "0"
}

// fulfillResource performs the app-level work promised by a redeemed token.
// Intent: Storage and compute outcomes are real deterministic work, not symbolic
// labels, so keep/break evidence has concrete content. Source: DI-sirus
func (node *Node) fulfillResource(message WireMessage) (map[string]string, error) {
	switch message.Token.ResourceKind {
	case data.Kind:
		value, ok := node.Data[message.Token.ResourceID]
		if !ok {
			return nil, fmt.Errorf("%s has no promised data resource %s", node.Name, message.Token.ResourceID)
		}
		return map[string]string{"resource_kind": data.Kind, "resource_id": message.Token.ResourceID, "value": value}, nil
	case storage.Kind:
		return node.fulfillStorage(message)
	case compute.Kind:
		return node.fulfillCompute(message)
	default:
		return map[string]string{"resource_kind": message.Token.ResourceKind, "resource_id": message.Token.ResourceID}, nil
	}
}

func (node *Node) fulfillStorage(message WireMessage) (map[string]string, error) {
	key := message.Payload["key"]
	if key == "" {
		return nil, fmt.Errorf("storage redemption missing key")
	}
	switch message.Payload["op"] {
	case "store":
		value := message.Payload["value"]
		node.Storage[key] = value
		return map[string]string{"resource_kind": storage.Kind, "op": "store", "key": key, "stored": "true", "value_sha256": protocol.HashExactBytes([]byte(value))}, nil
	case "read":
		value, ok := node.Storage[key]
		if !ok {
			return nil, fmt.Errorf("%s has not stored key %s", node.Name, key)
		}
		return map[string]string{"resource_kind": storage.Kind, "op": "read", "key": key, "value": value}, nil
	default:
		return nil, fmt.Errorf("storage promise has no local behavior for op %q", message.Payload["op"])
	}
}

func (node *Node) fulfillCompute(message WireMessage) (map[string]string, error) {
	inputText := message.Payload["n"]
	if inputText == "" {
		inputText = message.Token.ResourceID
	}
	inputNumber, parseErr := parseFibonacciInput(inputText)
	if parseErr != nil {
		return nil, parseErr
	}
	result, resultErr := fibonacci(inputNumber)
	if resultErr != nil {
		return nil, resultErr
	}
	return map[string]string{"resource_kind": compute.Kind, "n": fmt.Sprintf("%d", inputNumber), "result": fmt.Sprintf("%d", result)}, nil
}

// applyRedemptionEvidence updates the holder's local trust after the issuer
// keeps, breaks, or refuses a token promise. It also judges the peer who
// circulated the token when that peer is different from the issuer. Source:
// DI-sirus
func (node *Node) applyRedemptionEvidence(heldToken token.Token, response WireResponse) {
	event := token.Event{Observer: node.Name, Event: "held_redemption_observed", Outcome: response.Outcome, TokenID: heldToken.ID, Detail: response.Detail}
	node.Wallet.ApplyRedemption(event)
	node.record("holder_trust_updated", response.Outcome, heldToken.ID, node.Name+" updated local trust for issuer "+heldToken.Issuer+" after redemption outcome "+response.Outcome)
	sourcePeer := node.TokenSources[heldToken.ID]
	if sourcePeer == "" || sourcePeer == "historical" || sourcePeer == heldToken.Issuer {
		return
	}
	node.Wallet.ApplyPeerObservation(sourcePeer, event)
	node.record("circulator_trust_updated", response.Outcome, heldToken.ID, node.Name+" updated local trust for circulating peer "+sourcePeer+" after redemption outcome "+response.Outcome)
}

func (node *Node) send(toNode string, kind string, payload map[string]string, messageToken token.Token) (WireResponse, error) {
	route, routeErr := routeBetween(node.Name, toNode)
	if routeErr != nil {
		return WireResponse{}, routeErr
	}
	if payload == nil {
		payload = make(map[string]string)
	}
	nextHop := route[0]
	message, messageErr := node.newWireMessage(toNode, kind, route, payload, messageToken)
	if messageErr != nil {
		return WireResponse{}, messageErr
	}
	return transmit(nodeAddrs[nextHop], message)
}

// newWireMessage packages all app fields into one signed CBOR grid envelope.
// Intent: The pCID names the whole POC8 protocol, and payload `kind` selects
// the variant inside that one protocol vocabulary. Source: DI-sirus
func (node *Node) newWireMessage(toNode string, kind string, route []string, payload map[string]string, messageToken token.Token) (WireMessage, error) {
	fields := make(map[string]string)
	for key, value := range payload {
		fields[key] = value
	}
	fields["from_node"] = node.Name
	fields["from_app"] = economyAppName
	fields["to_node"] = toNode
	fields["to_app"] = economyAppName
	fields["kind"] = kind
	fields["route"] = encodeRoute(route)
	if !messageToken.IsZero() {
		tokenBytes, tokenErr := token.Encode(messageToken)
		if tokenErr != nil {
			return WireMessage{}, tokenErr
		}
		fields["token_bytes"] = hex.EncodeToString(tokenBytes)
	}
	envelope, envelopeErr := protocol.NewEnvelope(appProtocolCID, fields, node.Name)
	if envelopeErr != nil {
		return WireMessage{}, envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return WireMessage{}, bytesErr
	}
	return WireMessage{FromNode: node.Name, FromApp: economyAppName, ToNode: toNode, ToApp: economyAppName, Kind: kind, Route: append([]string(nil), route...), Payload: fields, Token: messageToken, Envelope: hex.EncodeToString(envelopeBytes)}, nil
}

func (node *Node) decodeEnvelopeBytes(envelopeBytes []byte, message WireMessage) (WireMessage, []byte, error) {
	envelope, parseErr := protocol.ParseEnvelope(envelopeBytes)
	if parseErr != nil {
		return WireMessage{}, nil, parseErr
	}
	if !envelope.ProtocolCID.Equal(appProtocolCID) {
		return WireMessage{}, nil, fmt.Errorf("node %s has no local promise to interpret app envelope pCID %s", node.Name, envelope.ProtocolCID)
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return WireMessage{}, nil, verifyErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return WireMessage{}, nil, fieldsErr
	}
	decoded := message
	decoded.FromNode = fields["from_node"]
	decoded.FromApp = fields["from_app"]
	decoded.ToNode = fields["to_node"]
	decoded.ToApp = fields["to_app"]
	decoded.Kind = fields["kind"]
	decoded.Route = decodeRoute(fields["route"])
	decoded.Payload = fields
	decoded.Envelope = hex.EncodeToString(envelopeBytes)
	if tokenHex := fields["token_bytes"]; tokenHex != "" {
		tokenBytes, tokenHexErr := hex.DecodeString(tokenHex)
		if tokenHexErr != nil {
			return WireMessage{}, nil, tokenHexErr
		}
		decodedToken, tokenErr := token.Decode(tokenBytes)
		if tokenErr != nil {
			return WireMessage{}, nil, tokenErr
		}
		decoded.Token = decodedToken
	}
	return decoded, envelopeBytes, nil
}

func transmit(address string, message WireMessage) (WireResponse, error) {
	envelopeBytes, envelopeErr := message.EnvelopeBytes()
	if envelopeErr != nil {
		return WireResponse{}, envelopeErr
	}
	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		frameConn, dialErr := transport.DialFrameConn(address, 10*time.Second)
		if dialErr != nil {
			lastErr = dialErr
			time.Sleep(250 * time.Millisecond)
			continue
		}
		if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
			closeFrame(frameConn)
			return WireResponse{}, writeErr
		}
		responseBytes, readErr := frameConn.ReadFrame()
		closeFrame(frameConn)
		if readErr != nil {
			return WireResponse{}, readErr
		}
		var result WireResponse
		if decodeErr := json.Unmarshal(responseBytes, &result); decodeErr != nil {
			return WireResponse{}, decodeErr
		}
		return result, nil
	}
	return WireResponse{}, fmt.Errorf("tcp frame to %s failed after retries: %w", address, lastErr)
}

func routeBetween(fromNode string, toNode string) ([]string, error) {
	routes := map[string]map[string][]string{
		"alice":   {"bob": {"bob"}, "carol": {"bob", "carol"}, "dave": {"mallory", "dave"}, "mallory": {"mallory"}},
		"bob":     {"alice": {"alice"}, "carol": {"carol"}, "dave": {"carol", "dave"}, "mallory": {"alice", "mallory"}},
		"carol":   {"alice": {"bob", "alice"}, "bob": {"bob"}, "dave": {"dave"}, "mallory": {"dave", "mallory"}},
		"dave":    {"alice": {"mallory", "alice"}, "bob": {"carol", "bob"}, "carol": {"carol"}, "mallory": {"mallory"}},
		"mallory": {"alice": {"alice"}, "bob": {"alice", "bob"}, "carol": {"dave", "carol"}, "dave": {"dave"}},
	}
	if fromNode == toNode {
		return nil, fmt.Errorf("route from %s to itself is not an inter-peer promise", fromNode)
	}
	peerRoutes, ok := routes[fromNode]
	if !ok {
		return nil, fmt.Errorf("unknown route source %s", fromNode)
	}
	route, ok := peerRoutes[toNode]
	if !ok {
		return nil, fmt.Errorf("unknown route from %s to %s", fromNode, toNode)
	}
	return append([]string(nil), route...), nil
}

func encodeRoute(route []string) string {
	return strings.Join(route, ",")
}

func decodeRoute(routeText string) []string {
	if strings.TrimSpace(routeText) == "" {
		return nil
	}
	return strings.Split(routeText, ",")
}

func copyPayload(payload map[string]string) map[string]string {
	copiedPayload := make(map[string]string)
	for key, value := range payload {
		copiedPayload[key] = value
	}
	return copiedPayload
}

func (node *Node) record(event string, outcome string, tokenID string, detail string) {
	record := token.Event{Observer: node.Name, Event: event, Outcome: outcome, TokenID: tokenID, Detail: detail}
	node.Evidence = append(node.Evidence, record)
	bytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		log.Printf("encode evidence record: %v", marshalErr)
		return
	}
	fmt.Println(string(bytes))
}

func (node *Node) writeDone() error {
	if err := os.MkdirAll(filepath.Join(node.RunRoot, node.RunID), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(node.RunRoot, node.RunID, "done"), []byte("done\n"), 0o644)
}

func (node *Node) waitForDone() error {
	donePath := filepath.Join(node.RunRoot, node.RunID, "done")
	for deadline := time.Now().Add(35 * time.Second); time.Now().Before(deadline); {
		if _, err := os.Stat(donePath); err == nil {
			node.record("node_complete", token.OutcomeKept, "", node.Name+" observed bounded demo completion")
			marker := filepath.Join(node.RunRoot, node.RunID, node.Name+".complete")
			return os.WriteFile(marker, []byte(node.Name+"\n"), 0o644)
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("%s timed out waiting for %s", node.Name, donePath)
}

func historicalAliceStaleToken(tokenID string) (token.Token, error) {
	issuer := token.NewIssuer("alice")
	return issuer.Issue(tokenID, "mallory", data.Kind, "dataset-stale", token.TransferBearer)
}

func closeFrame(frameConn transport.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		log.Printf("close tcp frame: %v", closeErr)
	}
}

func getenv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func initialData(nodeName string) map[string]string {
	if nodeName != "alice" {
		return map[string]string{}
	}
	return map[string]string{
		"dataset-public":  "public weather sample",
		"dataset-private": "private Alice dataset shared only after local trust and reciprocal promises",
		"dataset-stale":   "stale revoked dataset",
	}
}

func parseFibonacciInput(value string) (int, error) {
	trimmed := strings.TrimPrefix(value, "fib-")
	inputNumber, parseErr := strconv.Atoi(trimmed)
	if parseErr != nil {
		return 0, parseErr
	}
	return inputNumber, nil
}

// fibonacci performs bounded compute work for compute-token redemption. Intent:
// Keep compute promises concrete and deterministic while staying cheap enough
// for repeated Docker runs. Source: DI-sirus
func fibonacci(inputNumber int) (uint64, error) {
	if inputNumber < 0 || inputNumber > 93 {
		return 0, fmt.Errorf("fibonacci input %d is outside uint64 demo bounds", inputNumber)
	}
	if inputNumber == 0 {
		return 0, nil
	}
	var previous uint64
	current := uint64(1)
	for sequenceIndex := 1; sequenceIndex < inputNumber; sequenceIndex++ {
		previous, current = current, previous+current
	}
	return current, nil
}

func atoi(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	parsed, parseErr := strconv.Atoi(value)
	if parseErr != nil {
		return 0
	}
	return parsed
}
