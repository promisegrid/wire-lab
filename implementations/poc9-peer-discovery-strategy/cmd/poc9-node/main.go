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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/compute"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/data"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/policy"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/protocol"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/storage"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/token"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/transport"
)

const (
	economyAppName             = "promise-economy-discovery"
	kindNeedAdvertisement      = "need_advertisement"
	kindOfferPromise           = "offer_promise"
	kindCounterPromise         = "counter_promise"
	kindAcceptancePromise      = "acceptance_promise"
	kindTokenIssuePromise      = "token_issue_promise"
	kindTokenRedemptionPromise = "token_redemption_promise"
	kindOutcomeObservation     = "outcome_observation"
	kindReferralPromise        = "referral_promise"
	kindIntroductionPromise    = "introduction_promise"
	kindRoutePromise           = "route_promise"
)

const meshSeed = "poc9 deterministic sparse mesh seed: 2026-06-01 discovery strategy"

const (
	// Intent: POC9 uses deterministic logical time so expired-token evidence does
	// not depend on wall-clock timing or container startup order. Source: DI-vujil
	expiredTokenExpiresAtUnix = 100
	expiredTokenRedeemAtUnix  = 200
)

var allNodeNames = []string{"alice", "bob", "carol", "dave", "ellen", "frank", "mallory"}

var meshEdges = [][2]string{
	{"alice", "bob"},
	{"alice", "ellen"},
	{"bob", "carol"},
	{"bob", "frank"},
	{"carol", "dave"},
	{"carol", "mallory"},
	{"dave", "ellen"},
	{"ellen", "frank"},
	{"frank", "mallory"},
}

var nodeAddrs = map[string]string{
	"alice":   "alice:8077",
	"bob":     "bob:8077",
	"carol":   "carol:8077",
	"dave":    "dave:8077",
	"ellen":   "ellen:8077",
	"frank":   "frank:8077",
	"mallory": "mallory:8077",
}

// appProtocolCID names one POC9 protocol. Intent: Keep discovery, economy, route,
// token, and observation messages as payload variants under one pCID instead of
// minting a new pCID for every message kind. Source: DI-sipuz
var appProtocolCID = protocol.NewProtocolCID([]byte("poc9 peer discovery strategy protocol v1: one pCID for need advertisements, offer promises, counter promises, acceptances, token issue promises, token redemption promises, outcome observations, referrals, introductions, and route promises"))

// WireMessage is the local Go view of one signed POC9 app envelope. Intent: The
// struct is scaffolding; the actual peer-to-peer promise object remains the exact
// pCID-selected CBOR grid envelope carried over framed TCP. Source: DI-sipuz
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

// WireResponse is bounded POC feedback from the final promise observer. Intent:
// It is not a global authority; each receiver records its own local evidence from
// the signed message and the transport behavior it observes. Source: DI-sipuz
type WireResponse struct {
	Outcome string            `json:"outcome"`
	Detail  string            `json:"detail"`
	Token   token.Token       `json:"token"`
	Payload map[string]string `json:"payload"`
}

// EnvelopeBytes returns the exact signed CBOR grid bytes carried by the message.
func (message WireMessage) EnvelopeBytes() ([]byte, error) {
	if message.Envelope == "" {
		return nil, fmt.Errorf("wire message missing signed envelope")
	}
	return hex.DecodeString(message.Envelope)
}

// Mesh is the deterministic sparse physical-neighbor graph for POC9. Intent: The
// graph bounds who can open direct TCP connections while route promises decide
// which non-neighbor paths an agent is locally willing to use. Source: DI-sipuz
type Mesh struct {
	adjacency map[string][]string
}

// Neighbor is one direct TCP neighbor known to the local node.
type Neighbor struct {
	Name    string
	Address string
}

// DiscoveryState is one node's local relationship memory. Intent: Routes,
// referrals, introductions, transport observations, and completed promises are
// local evidence; none of these maps is a registry or global trust table. Source:
// DI-sipuz
type DiscoveryState struct {
	KnownRoutes           map[string][][]string
	Referrals             map[string]string
	Introductions         map[string]bool
	TransportTrust        map[string]int
	Completed             map[string]bool
	TransportObservations []TransportObservation
}

// LocalStrategy owns the deterministic choices one node makes during the POC.
// Intent: Keep the harness reproducible while ensuring each action is still a
// local promise decision by the acting node, not a command from Alice. Source:
// DI-sipuz
type LocalStrategy struct {
	Node *Node
}

// TransportObservation records what one node saw at the transport boundary.
// Intent: TCP behavior feeds local promise accounting, but open TCP is never
// treated as trust or authorization. Source: DI-sipuz
type TransportObservation struct {
	Peer    string `json:"peer"`
	Outcome string `json:"outcome"`
	Detail  string `json:"detail"`
}

// Node groups one local kernel boundary, local app roles, local discovery memory,
// local wallet state, and local transport evidence for one container. Source:
// DI-sipuz
type Node struct {
	Name         string
	RunID        string
	RunRoot      string
	Issuer       *token.Issuer
	Wallet       *token.Wallet
	Policy       policy.AgentPolicy
	Mesh         Mesh
	Discovery    DiscoveryState
	Evidence     []token.Event
	Storage      map[string]string
	Data         map[string]string
	TokenSources map[string]string
	mutex        sync.Mutex
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
	go LocalStrategy{Node: node}.Run()
	if waitErr := node.waitForAllDone(); waitErr != nil {
		log.Printf("%s wait failed: %v", node.Name, waitErr)
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

func newNode(nodeName string) *Node {
	mesh := NewMesh(meshEdges)
	node := &Node{
		Name:    nodeName,
		RunID:   getenv("POC9_RUN_ID", "manual"),
		RunRoot: "/run/poc9",
		Issuer:  token.NewIssuer(nodeName),
		Wallet:  token.NewWallet(nodeName),
		Policy:  policy.ForNode(nodeName),
		Mesh:    mesh,
		Discovery: DiscoveryState{
			KnownRoutes:    make(map[string][][]string),
			Referrals:      make(map[string]string),
			Introductions:  make(map[string]bool),
			TransportTrust: make(map[string]int),
			Completed:      make(map[string]bool),
		},
		Storage:      make(map[string]string),
		Data:         initialData(nodeName),
		TokenSources: make(map[string]string),
	}
	for _, neighborName := range mesh.NeighborNames(nodeName) {
		node.learnRoute(neighborName, []string{neighborName}, "direct mesh neighbor")
	}
	if nodeName == "alice" {
		seedExpiredAliceTokens(node)
	}
	return node
}

func seedExpiredAliceTokens(node *Node) {
	// Intent: Alice's historical bearer tokens are expired by their own signed
	// validity promise, not revoked. Dave should therefore judge Mallory's later
	// circulation, not treat Alice as breaking a promise. Source: DI-vujil
	node.Issuer.SetNowUnix(expiredTokenRedeemAtUnix)
	for _, tokenID := range []string{"alice-expired-bearer-1", "alice-expired-bearer-2"} {
		_, issueErr := node.Issuer.IssueExpiring(tokenID, "mallory", data.Kind, "dataset-expired", token.TransferBearer, expiredTokenExpiresAtUnix)
		if issueErr != nil {
			node.record("expired_seed_issue", token.OutcomeBroken, tokenID, issueErr.Error())
			continue
		}
	}
}

func NewMesh(edges [][2]string) Mesh {
	adjacency := make(map[string][]string)
	for _, nodeName := range allNodeNames {
		adjacency[nodeName] = []string{}
	}
	for _, edge := range edges {
		leftNode := edge[0]
		rightNode := edge[1]
		adjacency[leftNode] = append(adjacency[leftNode], rightNode)
		adjacency[rightNode] = append(adjacency[rightNode], leftNode)
	}
	for nodeName := range adjacency {
		sort.Strings(adjacency[nodeName])
	}
	return Mesh{adjacency: adjacency}
}

func (mesh Mesh) NeighborNames(nodeName string) []string {
	neighbors := append([]string(nil), mesh.adjacency[nodeName]...)
	sort.Strings(neighbors)
	return neighbors
}

func (mesh Mesh) HasEdge(leftNode string, rightNode string) bool {
	for _, neighborName := range mesh.adjacency[leftNode] {
		if neighborName == rightNode {
			return true
		}
	}
	return false
}

func (mesh Mesh) Nodes() []string {
	nodeNames := make([]string, 0, len(mesh.adjacency))
	for nodeName := range mesh.adjacency {
		nodeNames = append(nodeNames, nodeName)
	}
	sort.Strings(nodeNames)
	return nodeNames
}

func (mesh Mesh) Route(fromNode string, toNode string) ([]string, error) {
	if fromNode == toNode {
		return nil, fmt.Errorf("route from %s to itself is not an inter-peer promise", fromNode)
	}
	if _, exists := mesh.adjacency[fromNode]; !exists {
		return nil, fmt.Errorf("unknown mesh source %s", fromNode)
	}
	if _, exists := mesh.adjacency[toNode]; !exists {
		return nil, fmt.Errorf("unknown mesh target %s", toNode)
	}
	queue := []string{fromNode}
	previousNode := map[string]string{fromNode: ""}
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		if currentNode == toNode {
			break
		}
		for _, neighborName := range mesh.adjacency[currentNode] {
			if _, seen := previousNode[neighborName]; seen {
				continue
			}
			previousNode[neighborName] = currentNode
			queue = append(queue, neighborName)
		}
	}
	if _, found := previousNode[toNode]; !found {
		return nil, fmt.Errorf("mesh has no route from %s to %s", fromNode, toNode)
	}
	fullPath := []string{}
	for pathNode := toNode; pathNode != ""; pathNode = previousNode[pathNode] {
		fullPath = append([]string{pathNode}, fullPath...)
	}
	if len(fullPath) < 2 || fullPath[0] != fromNode {
		return nil, fmt.Errorf("invalid reconstructed mesh path from %s to %s", fromNode, toNode)
	}
	return append([]string(nil), fullPath[1:]...), nil
}

func (node *Node) Neighbors() []Neighbor {
	neighborNames := node.Mesh.NeighborNames(node.Name)
	neighbors := make([]Neighbor, 0, len(neighborNames))
	for _, neighborName := range neighborNames {
		neighbors = append(neighbors, Neighbor{Name: neighborName, Address: nodeAddrs[neighborName]})
	}
	return neighbors
}

func (strategy LocalStrategy) Run() {
	time.Sleep(900 * time.Millisecond)
	strategy.Node.publishDiscoveryPromises()
	switch strategy.Node.Name {
	case "alice":
		strategy.runAlice()
	case "bob":
		strategy.runBob()
	case "carol":
		strategy.runCarol()
	case "dave":
		strategy.runDave()
	case "ellen":
		strategy.runRelayNode("ellen")
	case "frank":
		strategy.runRelayNode("frank")
	case "mallory":
		strategy.runMallory()
	}
}

func (strategy LocalStrategy) runAlice() {
	strategy.Node.waitForCompletion("referral:carol:compute", 5*time.Second)
	// Intent: Alice treats Bob's referral and a route promise as separate local
	// evidence before sending a non-neighbor compute need to Carol. Source: DI-sipuz
	strategy.Node.waitForCompletion("route:carol", 8*time.Second)
	strategy.Node.sendNeed("bob", storage.Kind, "alice-public-note", "public", map[string]string{
		"key":        "alice-public-note",
		"value":      "public calibration note",
		"max_price":  "3",
		"stake":      "1",
		"alice_note": "Alice starts with public storage before private data",
	})
	strategy.Node.sendNeed("carol", compute.Kind, "fib-8", "public", map[string]string{
		"n":          "8",
		"max_price":  "3",
		"stake":      "1",
		"alice_note": "Alice starts with known compute before trusting private work",
	})
	if strategy.Node.waitForCompletion("kept:bob:storage:public", 8*time.Second) {
		strategy.Node.sendNeed("bob", storage.Kind, "alice-private-report", "private", map[string]string{
			"key":        "alice-private-report",
			"value":      strategy.Node.Data["dataset-private"],
			"max_price":  "8",
			"stake":      "2",
			"alice_note": "Alice escalates to private storage after Bob keeps a public promise",
		})
	}
	if strategy.Node.waitForCompletion("kept:carol:compute:public", 8*time.Second) {
		strategy.Node.sendNeed("carol", compute.Kind, "fib-10", "private", map[string]string{
			"n":          "10",
			"max_price":  "6",
			"stake":      "2",
			"alice_note": "Alice escalates to higher-value compute after Carol keeps a public promise",
		})
	}
	strategy.Node.waitForCompletion("kept:bob:storage:private", 10*time.Second)
	strategy.Node.waitForCompletion("kept:carol:compute:private", 10*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"kept:bob:storage:public", "kept:bob:storage:private", "kept:carol:compute:public", "kept:carol:compute:private"})
}

func (strategy LocalStrategy) runBob() {
	time.Sleep(700 * time.Millisecond)
	strategy.Node.sendReferral("alice", "carol", compute.Kind, "fib-8", "Bob promises only that Bob has local evidence Carol advertises compute; Alice must judge Carol locally")
	strategy.Node.waitForCompletion("kept:bob:storage:public", 18*time.Second)
	strategy.Node.waitForCompletion("kept:bob:storage:private", 18*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"referral_sent:alice:carol:compute", "kept:bob:storage:public", "kept:bob:storage:private"})
}

func (strategy LocalStrategy) runCarol() {
	strategy.Node.waitForCompletion("transport:malformed:mallory", 12*time.Second)
	strategy.Node.waitForCompletion("route_refused:mallory", 12*time.Second)
	strategy.Node.waitForCompletion("kept:carol:compute:public", 18*time.Second)
	strategy.Node.waitForCompletion("kept:carol:compute:private", 18*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"transport:malformed:mallory", "route_refused:mallory", "kept:carol:compute:public", "kept:carol:compute:private"})
}

func (strategy LocalStrategy) runDave() {
	strategy.Node.waitForCompletion("expired:first:observed", 20*time.Second)
	strategy.Node.waitForCompletion("expired:second:refused", 20*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"expired:first:observed", "expired:second:refused"})
}

func (strategy LocalStrategy) runRelayNode(nodeName string) {
	strategy.Node.waitForCompletion("relayed", 22*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"discovery_published", "relayed"})
	strategy.Node.record("relay_strategy_done", token.OutcomeKept, "", nodeName+" completed local relay evidence")
}

func (strategy LocalStrategy) runMallory() {
	time.Sleep(1400 * time.Millisecond)
	strategy.Node.sendMalformed("carol")
	time.Sleep(1200 * time.Millisecond)
	strategy.Node.sendMalloryRoutePromiseToCarol()
	// Intent: Mallory can only pressure Dave through routes Mallory has locally
	// learned from peers after Carol refuses the direct Carol path. Source: DI-sipuz
	strategy.Node.waitForCompletion("route:dave", 8*time.Second)
	time.Sleep(1800 * time.Millisecond)
	strategy.Node.offerExpiredTokenToDave("alice-expired-bearer-1", "mallory-expired-offer-1")
	time.Sleep(3600 * time.Millisecond)
	strategy.Node.offerExpiredTokenToDave("alice-expired-bearer-2", "mallory-expired-offer-2")
	strategy.Node.waitForCompletion("mallory:expired_second_sent", 5*time.Second)
	strategy.Node.writeNodeDoneIfComplete([]string{"mallory:malformed_sent", "mallory:route_refused_observed", "mallory:expired_first_sent", "mallory:expired_second_sent"})
}

func (node *Node) publishDiscoveryPromises() {
	for _, neighbor := range node.Neighbors() {
		introPayload := map[string]string{
			"introduced_peer": node.Name,
			"promiser":        node.Name,
			"mesh_seed":       meshSeed,
			"promise":         node.Name + " promises only its own direct neighbor presence and willingness to receive this POC9 pCID",
		}
		if _, sendErr := node.send(neighbor.Name, kindIntroductionPromise, introPayload, token.Token{}); sendErr != nil {
			node.record("introduction_send", token.OutcomeBroken, "", sendErr.Error())
		}
		for _, targetNode := range node.Mesh.Nodes() {
			if targetNode == node.Name || targetNode == neighbor.Name {
				continue
			}
			routeFromNode, routeErr := node.Mesh.Route(node.Name, targetNode)
			if routeErr != nil || len(routeFromNode) == 0 || routeFromNode[0] == neighbor.Name {
				continue
			}
			promisedRoute := append([]string{node.Name}, routeFromNode...)
			routePayload := map[string]string{
				"target_node":    targetNode,
				"promised_route": encodeRoute(promisedRoute),
				"promiser":       node.Name,
				"promise":        node.Name + " promises only that this is its local route offer; receivers judge the route locally",
			}
			if _, routeSendErr := node.send(neighbor.Name, kindRoutePromise, routePayload, token.Token{}); routeSendErr != nil {
				node.record("route_promise_send", token.OutcomeBroken, "", routeSendErr.Error())
			}
		}
	}
	node.markCompleted("discovery_published", "published introductions and route promises to direct neighbors")
}

func (node *Node) sendReferral(toNode string, referredPeer string, resourceKind string, resourceID string, promiseText string) {
	payload := map[string]string{
		"referred_peer": referredPeer,
		"resource_kind": resourceKind,
		"resource_id":   resourceID,
		"referrer":      node.Name,
		"promise":       promiseText,
	}
	if _, sendErr := node.send(toNode, kindReferralPromise, payload, token.Token{}); sendErr != nil {
		node.record("referral_send", token.OutcomeBroken, "", sendErr.Error())
		return
	}
	node.markCompleted("referral_sent:"+toNode+":"+referredPeer+":"+resourceKind, "sent referral promise")
}

func (node *Node) sendNeed(toNode string, resourceKind string, resourceID string, stage string, payload map[string]string) {
	needPayload := copyPayload(payload)
	needPayload["resource_kind"] = resourceKind
	needPayload["resource_id"] = resourceID
	needPayload["stage"] = stage
	needPayload["need_id"] = node.Name + "-" + resourceKind + "-" + stage
	needPayload["promiser"] = node.Name
	needPayload["alice_promises"] = "Alice promises this need is current and that she will judge keep/break history locally before escalating risk"
	needPayload["requested_peer_promises"] = "peer promises only its own storage, compute, token, or route behavior"
	node.record("need_advertised", token.OutcomeKept, "", node.Name+" advertised "+stage+" "+resourceKind+" need to "+toNode)
	if _, sendErr := node.send(toNode, kindNeedAdvertisement, needPayload, token.Token{}); sendErr != nil {
		node.record("need_advertise_send", token.OutcomeBroken, "", sendErr.Error())
	}
}

func (node *Node) offerExpiredTokenToDave(tokenID string, offerID string) {
	expiredToken, expiredErr := historicalAliceExpiredToken(tokenID)
	if expiredErr != nil {
		node.record("expired_offer_prepare", token.OutcomeBroken, tokenID, expiredErr.Error())
		return
	}
	payload := map[string]string{
		"offer_id":       offerID,
		"resource_kind":  expiredToken.ResourceKind,
		"resource_id":    expiredToken.ResourceID,
		"stage":          "expired",
		"price":          "1",
		"stake":          "0",
		"offer_promises": "Mallory promises only that she is voluntarily presenting these expired token bytes as useful; Dave must judge Alice and Mallory locally",
		"expires_at":     fmt.Sprintf("%d", expiredToken.ExpiresAtUnix),
		"issuer":         expiredToken.Issuer,
	}
	response, offerErr := node.send("dave", kindOfferPromise, payload, expiredToken)
	if offerErr != nil {
		node.record("expired_offer_send", token.OutcomeBroken, expiredToken.ID, offerErr.Error())
		return
	}
	// Intent: Mallory only counts the expired-token pressure as sent when it
	// reaches Dave; a relay refusal is route evidence, not a completed Dave
	// interaction. Source: DI-sipuz; DI-vujil
	if refusingNode := response.Payload["refusing_node"]; refusingNode != "" {
		node.record("expired_offer_route_refused", token.OutcomeRefused, expiredToken.ID, refusingNode+" refused expired-token route")
		return
	}
	if tokenID == "alice-expired-bearer-1" {
		node.markCompleted("mallory:expired_first_sent", "Mallory circulated first expired Alice token")
	} else {
		node.markCompleted("mallory:expired_second_sent", "Mallory circulated second expired Alice token")
	}
}

func (node *Node) sendMalloryRoutePromiseToCarol() {
	payload := map[string]string{
		"target_node":    "dave",
		"promised_route": "mallory,carol,dave",
		"promiser":       "mallory",
		"promise":        "Mallory asks Carol to treat this route as usable after malformed transport evidence",
	}
	response, sendErr := node.send("carol", kindRoutePromise, payload, token.Token{})
	if sendErr != nil {
		node.record("mallory_route_send", token.OutcomeBroken, "", sendErr.Error())
		return
	}
	if response.Outcome == token.OutcomeRefused {
		node.recordTransportObservation("carol", token.OutcomeRefused, "Carol refused Mallory route after malformed transport evidence")
		node.markCompleted("route_refused_by:carol", "Mallory will avoid Carol as a next hop after Carol's local refusal")
		node.markCompleted("mallory:route_refused_observed", "Mallory observed Carol refusing later route promise")
	}
}

func (node *Node) sendMalformed(toNode string) {
	address := nodeAddrs[toNode]
	frameConn, dialErr := transport.DialFrameConn(address, 5*time.Second)
	if dialErr != nil {
		node.record("malformed_dial", token.OutcomeBroken, "", dialErr.Error())
		return
	}
	defer closeFrame(frameConn)
	malformedBytes := []byte("from_node=mallory; this frame is intentionally not a cbor grid envelope")
	if writeErr := frameConn.WriteFrame(malformedBytes); writeErr != nil {
		node.record("malformed_write", token.OutcomeBroken, "", writeErr.Error())
		return
	}
	if _, readErr := frameConn.ReadFrame(); readErr != nil {
		node.record("malformed_response", token.OutcomeBroken, "", readErr.Error())
		return
	}
	node.markCompleted("mallory:malformed_sent", "Mallory sent intentionally malformed bytes to Carol")
}

func (node *Node) serveTCP(listener net.Listener, serverErrors chan<- error) {
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErrors <- acceptErr
			return
		}
		go node.handleConnection(transport.NewFrameConn(conn))
	}
}

func (node *Node) handleConnection(frameConn transport.FrameConn) {
	defer closeFrame(frameConn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		node.recordTransportObservation("unknown", token.OutcomeBroken, "read frame failed: "+readErr.Error())
		node.writeResponse(frameConn, WireResponse{Outcome: token.OutcomeBroken, Detail: readErr.Error()})
		return
	}
	decodedMessage, _, decodeErr := node.decodeEnvelopeBytes(frameBytes, WireMessage{})
	if decodeErr != nil {
		guessedPeer := guessPeerFromMalformed(frameBytes)
		node.recordTransportObservation(guessedPeer, token.OutcomeBroken, "malformed envelope bytes: "+decodeErr.Error())
		if guessedPeer != "unknown" {
			node.markCompleted("transport:malformed:"+guessedPeer, "observed malformed bytes from "+guessedPeer)
		}
		node.writeResponse(frameConn, WireResponse{Outcome: token.OutcomeBroken, Detail: "malformed envelope observed"})
		return
	}
	node.recordTransportObservation(decodedMessage.FromNode, token.OutcomeKept, "received parseable grid envelope")
	response, handleErr := node.handleWireMessage(decodedMessage)
	if handleErr != nil {
		response = WireResponse{Outcome: token.OutcomeBroken, Detail: handleErr.Error()}
	}
	node.writeResponse(frameConn, response)
}

func (node *Node) writeResponse(frameConn transport.FrameConn, response WireResponse) {
	responseBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		log.Printf("%s response marshal failed: %v", node.Name, marshalErr)
		return
	}
	if writeErr := frameConn.WriteFrame(responseBytes); writeErr != nil {
		log.Printf("%s response write failed: %v", node.Name, writeErr)
	}
}

func (node *Node) handleWireMessage(message WireMessage) (WireResponse, error) {
	if message.ToNode != node.Name {
		nextHop, hopErr := node.nextRouteHop(message)
		if hopErr != nil {
			return WireResponse{}, hopErr
		}
		// Intent: A relay only promises its own forwarding behavior; if its local
		// evidence says the source recently sent malformed bytes, the relay can refuse
		// without appealing to an external authority. Source: DI-sipuz
		if node.hasCompleted("transport:malformed:"+message.FromNode) || node.transportTrust(message.FromNode) < 0 {
			node.record("relay_refused", token.OutcomeRefused, "", node.Name+" refused to relay for low-trust source "+message.FromNode)
			return WireResponse{Outcome: token.OutcomeRefused, Detail: node.Name + " refused low-trust relay source", Payload: map[string]string{"refusing_node": node.Name}}, nil
		}
		response, forwardErr := node.transmit(nextHop, message)
		if forwardErr != nil {
			return WireResponse{}, forwardErr
		}
		node.markCompleted("relayed", node.Name+" relayed "+message.Kind+" toward "+message.ToNode)
		return response, nil
	}
	switch message.Kind {
	case kindIntroductionPromise:
		return node.handleIntroductionPromise(message)
	case kindRoutePromise:
		return node.handleRoutePromise(message)
	case kindReferralPromise:
		return node.handleReferralPromise(message)
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
		return node.handleOutcomeObservation(message)
	default:
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "unknown payload kind " + message.Kind}, nil
	}
}

func (node *Node) handleIntroductionPromise(message WireMessage) (WireResponse, error) {
	introducedPeer := message.Payload["introduced_peer"]
	if introducedPeer == "" {
		introducedPeer = message.FromNode
	}
	node.mutex.Lock()
	node.Discovery.Introductions[introducedPeer] = true
	node.mutex.Unlock()
	node.record("introduction_observed", token.OutcomeKept, "", node.Name+" observed introduction promise for "+introducedPeer+" from "+message.FromNode)
	return WireResponse{Outcome: token.OutcomeKept, Detail: "introduction observed"}, nil
}

func (node *Node) handleRoutePromise(message WireMessage) (WireResponse, error) {
	// Intent: Specific malformed-byte evidence is enough for a receiver to refuse
	// later route promises, even if earlier benign discovery traffic left aggregate
	// transport trust non-negative. Source: DI-sipuz
	if node.hasCompleted("transport:malformed:"+message.FromNode) || node.transportTrust(message.FromNode) < 0 {
		node.record("route_promise_refused", token.OutcomeRefused, "", node.Name+" refused route promise from low-trust peer "+message.FromNode)
		node.markCompleted("route_refused:"+message.FromNode, "refused route promise from low-trust peer")
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "route promise refused after local transport evidence", Payload: map[string]string{"refusing_node": node.Name}}, nil
	}
	targetNode := message.Payload["target_node"]
	// Intent: Route promises use `promised_route` because `route` is already the
	// envelope's current transport path. Keeping those fields separate prevents the
	// transport path from erasing the promise being evaluated. Source: DI-sipuz
	route := decodeRoute(message.Payload["promised_route"])
	if targetNode == "" || len(route) == 0 {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "route promise missing target or route"}, nil
	}
	if route[0] != message.FromNode {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "route promise does not begin with promiser"}, nil
	}
	if !node.validPromisedRoute(route) {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "route promise does not match sparse mesh"}, nil
	}
	node.learnRoute(targetNode, route, message.FromNode)
	node.record("route_promise_observed", token.OutcomeKept, "", node.Name+" learned route to "+targetNode+" via "+message.FromNode+" as "+encodeRoute(route))
	return WireResponse{Outcome: token.OutcomeKept, Detail: "route promise observed"}, nil
}

func (node *Node) handleReferralPromise(message WireMessage) (WireResponse, error) {
	referredPeer := message.Payload["referred_peer"]
	resourceKind := message.Payload["resource_kind"]
	resourceID := message.Payload["resource_id"]
	if referredPeer == "" || resourceKind == "" {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "referral missing peer or resource kind"}, nil
	}
	referralKey := referredPeer + ":" + resourceKind + ":" + resourceID
	node.mutex.Lock()
	node.Discovery.Referrals[referralKey] = message.FromNode
	node.mutex.Unlock()
	node.record("referral_observed", token.OutcomeKept, "", node.Name+" treats "+message.FromNode+" referral to "+referredPeer+" as local evidence, not authority")
	if referredPeer == "carol" && resourceKind == compute.Kind {
		node.markCompleted("referral:carol:compute", "observed Bob referral about Carol compute")
	}
	return WireResponse{Outcome: token.OutcomeKept, Detail: "referral observed"}, nil
}

func (node *Node) handleNeedAdvertisement(message WireMessage) (WireResponse, error) {
	resourceKind := message.Payload["resource_kind"]
	if !node.provides(resourceKind) {
		node.record("need_ignored", token.OutcomeRefused, "", node.Name+" has no local resource promise for "+resourceKind)
		return WireResponse{Outcome: token.OutcomeRefused, Detail: node.Name + " does not promise " + resourceKind}, nil
	}
	stage := message.Payload["stage"]
	if stage == "private" && !node.hasCompleted("kept:"+node.Name+":"+resourceKind+":public") {
		node.record("private_need_without_probe", token.OutcomeRefused, "", node.Name+" refused private need before its public promise evidence existed")
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "private need refused until public promise evidence exists"}, nil
	}
	offerPayload := map[string]string{
		"offer_id":       node.Name + "-" + resourceKind + "-" + stage + "-offer",
		"resource_kind":  resourceKind,
		"resource_id":    message.Payload["resource_id"],
		"stage":          stage,
		"price":          message.Payload["max_price"],
		"stake":          message.Payload["stake"],
		"offer_promises": node.Name + " promises only its own " + resourceKind + " behavior for this bounded POC9 need",
	}
	if _, sendErr := node.send(message.FromNode, kindOfferPromise, offerPayload, token.Token{}); sendErr != nil {
		node.record("offer_send", token.OutcomeBroken, "", sendErr.Error())
		return WireResponse{Outcome: token.OutcomeBroken, Detail: sendErr.Error()}, nil
	}
	node.record("offer_sent", token.OutcomeKept, "", node.Name+" offered "+stage+" "+resourceKind+" promise to "+message.FromNode)
	return WireResponse{Outcome: token.OutcomeKept, Detail: "offer sent"}, nil
}

func (node *Node) handleOfferPromise(message WireMessage) (WireResponse, error) {
	if node.Name == "dave" && !message.Token.IsZero() {
		return node.handleExpiredTokenOffer(message)
	}
	if node.Name != "alice" {
		node.record("offer_observed", token.OutcomeKept, "", node.Name+" observed offer promise from "+message.FromNode)
		return WireResponse{Outcome: token.OutcomeKept, Detail: "offer observed"}, nil
	}
	resourceKind := message.Payload["resource_kind"]
	stage := message.Payload["stage"]
	if !node.canAcceptOffer(message) {
		node.record("offer_refused", token.OutcomeRefused, "", "Alice refused private "+resourceKind+" offer before public evidence")
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "private offer refused before public evidence"}, nil
	}
	acceptPayload := copyPayload(message.Payload)
	acceptPayload["acceptance_promises"] = "Alice promises to present only bounded POC9 data and record keep/break evidence locally"
	if _, sendErr := node.send(message.FromNode, kindAcceptancePromise, acceptPayload, token.Token{}); sendErr != nil {
		node.record("acceptance_send", token.OutcomeBroken, "", sendErr.Error())
		return WireResponse{Outcome: token.OutcomeBroken, Detail: sendErr.Error()}, nil
	}
	node.record("offer_accepted", token.OutcomeKept, "", "Alice accepted "+stage+" "+resourceKind+" promise from "+message.FromNode)
	return WireResponse{Outcome: token.OutcomeKept, Detail: "offer accepted"}, nil
}

func (node *Node) canAcceptOffer(message WireMessage) bool {
	resourceKind := message.Payload["resource_kind"]
	stage := message.Payload["stage"]
	return stage != "private" || node.hasCompleted("kept:"+message.FromNode+":"+resourceKind+":public")
}

func (node *Node) handleCounterPromise(message WireMessage) (WireResponse, error) {
	node.record("counter_observed", token.OutcomeKept, "", node.Name+" observed counter promise from "+message.FromNode)
	return WireResponse{Outcome: token.OutcomeKept, Detail: "counter observed"}, nil
}

func (node *Node) handleAcceptancePromise(message WireMessage) (WireResponse, error) {
	resourceKind := message.Payload["resource_kind"]
	if !node.provides(resourceKind) {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: node.Name + " cannot issue token for " + resourceKind}, nil
	}
	stage := message.Payload["stage"]
	resourceID := message.Payload["resource_id"]
	issuedToken, issueErr := node.Issuer.Issue(node.Name+"-"+resourceKind+"-"+stage+"-"+message.FromNode, message.FromNode, resourceKind, resourceID, token.TransferNonTransferable)
	if issueErr != nil {
		return WireResponse{}, issueErr
	}
	payload := copyPayload(message.Payload)
	payload["token_promises"] = node.Name + " issues a non-transferable promise token for " + resourceKind
	if _, sendErr := node.send(message.FromNode, kindTokenIssuePromise, payload, issuedToken); sendErr != nil {
		node.record("token_issue_send", token.OutcomeBroken, issuedToken.ID, sendErr.Error())
		return WireResponse{Outcome: token.OutcomeBroken, Detail: sendErr.Error()}, nil
	}
	node.record("token_issued", token.OutcomeKept, issuedToken.ID, node.Name+" issued "+stage+" "+resourceKind+" promise token to "+message.FromNode)
	return WireResponse{Outcome: token.OutcomeKept, Detail: "token issued", Token: issuedToken}, nil
}

func (node *Node) handleTokenIssuePromise(message WireMessage) (WireResponse, error) {
	if message.Token.IsZero() {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "token issue carried no token"}, nil
	}
	node.Wallet.Add(message.Token, node.Name+" received token from "+message.FromNode)
	node.TokenSources[message.Token.ID] = message.FromNode
	node.record("token_received", token.OutcomeKept, message.Token.ID, node.Name+" received "+message.Payload["stage"]+" "+message.Token.ResourceKind+" token from "+message.FromNode)
	redemptionPayload := copyPayload(message.Payload)
	redemptionPayload["holder"] = node.Name
	response, redeemErr := node.send(message.Token.Issuer, kindTokenRedemptionPromise, redemptionPayload, message.Token)
	if redeemErr != nil {
		node.record("token_redemption_send", token.OutcomeBroken, message.Token.ID, redeemErr.Error())
		return WireResponse{Outcome: token.OutcomeBroken, Detail: redeemErr.Error()}, nil
	}
	node.applyRedemptionEvidence(message.Token, response, message.Payload["stage"])
	return WireResponse{Outcome: response.Outcome, Detail: response.Detail, Payload: response.Payload, Token: message.Token}, nil
}

func (node *Node) handleTokenRedemptionPromise(message WireMessage) (WireResponse, error) {
	if message.Token.IsZero() {
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "redemption carried no token"}, nil
	}
	redemptionEvent := node.Issuer.Redeem(message.FromNode, message.Token)
	response := WireResponse{Outcome: redemptionEvent.Outcome, Detail: redemptionEvent.Detail, Token: message.Token, Payload: map[string]string{"stage": message.Payload["stage"], "resource_kind": message.Payload["resource_kind"]}}
	if redemptionEvent.Outcome != token.OutcomeKept {
		node.record("token_redemption", redemptionEvent.Outcome, message.Token.ID, redemptionEvent.Detail)
		return response, nil
	}
	switch message.Token.ResourceKind {
	case storage.Kind:
		key := message.Payload["key"]
		value := message.Payload["value"]
		node.Storage[key] = value
		response.Detail = node.Name + " stored and returned " + key
		response.Payload["stored_value"] = node.Storage[key]
	case compute.Kind:
		inputNumber, parseErr := parseFibonacciInput(message.Payload["resource_id"])
		if parseErr != nil {
			return WireResponse{Outcome: token.OutcomeBroken, Detail: parseErr.Error()}, nil
		}
		result, computeErr := fibonacci(inputNumber)
		if computeErr != nil {
			return WireResponse{Outcome: token.OutcomeBroken, Detail: computeErr.Error()}, nil
		}
		response.Detail = fmt.Sprintf("%s computed Fibonacci %d = %d", node.Name, inputNumber, result)
		response.Payload["result"] = fmt.Sprintf("%d", result)
	default:
		response.Outcome = token.OutcomeRefused
		response.Detail = "unknown resource kind " + message.Token.ResourceKind
	}
	node.record("token_redemption", response.Outcome, message.Token.ID, response.Detail)
	if response.Outcome == token.OutcomeKept {
		node.markCompleted("kept:"+node.Name+":"+message.Token.ResourceKind+":"+message.Payload["stage"], response.Detail)
	}
	return response, nil
}

func (node *Node) handleOutcomeObservation(message WireMessage) (WireResponse, error) {
	node.record("outcome_observation", message.Payload["outcome"], message.Payload["token_id"], message.Payload["detail"])
	return WireResponse{Outcome: token.OutcomeKept, Detail: "outcome observation recorded"}, nil
}

func (node *Node) handleExpiredTokenOffer(message WireMessage) (WireResponse, error) {
	// Intent: Dave accepts one expired token only to learn from the issuer's signed
	// expiry evidence. Expiry is neutral for Alice and negative only for Mallory if
	// Mallory presents it as useful. Source: DI-vujil
	if node.Wallet.Trust(message.FromNode) < 0 {
		node.record("expired_offer_refused", token.OutcomeRefused, message.Token.ID, "Dave refused second expired token after local Mallory trust decreased")
		node.markCompleted("expired:second:refused", "Dave refused Mallory after expired-token misrepresentation evidence")
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "expired token refused after local trust decrease"}, nil
	}
	node.Wallet.Add(message.Token, "Dave accepted first Mallory expired token for local evidence")
	node.TokenSources[message.Token.ID] = message.FromNode
	response, redeemErr := node.send(message.Token.Issuer, kindTokenRedemptionPromise, map[string]string{"resource_kind": message.Token.ResourceKind, "resource_id": message.Token.ResourceID, "stage": "expired", "holder": node.Name}, message.Token)
	if redeemErr != nil {
		node.record("expired_redeem_send", token.OutcomeBroken, message.Token.ID, redeemErr.Error())
		return WireResponse{Outcome: token.OutcomeBroken, Detail: redeemErr.Error()}, nil
	}
	node.applyRedemptionEvidence(message.Token, response, "expired")
	if response.Outcome == token.OutcomeExpired {
		node.markCompleted("expired:first:observed", "Dave observed Alice kept the signed expiry promise and judged Mallory locally")
	}
	return WireResponse{Outcome: response.Outcome, Detail: response.Detail, Token: message.Token}, nil
}

func (node *Node) applyRedemptionEvidence(heldToken token.Token, response WireResponse, stage string) {
	event := token.Event{Observer: node.Name, Event: "held_redemption_observed", Outcome: response.Outcome, TokenID: heldToken.ID, Detail: response.Detail}
	node.Wallet.ApplyRedemption(event)
	node.record("holder_trust_updated", response.Outcome, heldToken.ID, node.Name+" recorded local trust evidence for issuer "+heldToken.Issuer+" after redemption outcome "+response.Outcome)
	sourcePeer := node.TokenSources[heldToken.ID]
	if sourcePeer != "" && sourcePeer != "historical" && sourcePeer != heldToken.Issuer {
		node.Wallet.ApplyPeerObservation(sourcePeer, event)
		node.record("circulator_trust_updated", response.Outcome, heldToken.ID, node.Name+" updated local trust for circulating peer "+sourcePeer+" after redemption outcome "+response.Outcome)
	}
	if response.Outcome == token.OutcomeKept {
		node.markCompleted("kept:"+heldToken.Issuer+":"+heldToken.ResourceKind+":"+stage, response.Detail)
	}
}

func (node *Node) provides(resourceKind string) bool {
	return (node.Name == "bob" && resourceKind == storage.Kind) || (node.Name == "carol" && resourceKind == compute.Kind)
}

func (node *Node) send(toNode string, kind string, payload map[string]string, messageToken token.Token) (WireResponse, error) {
	route, routeErr := node.selectRoute(toNode)
	if routeErr != nil {
		return WireResponse{}, routeErr
	}
	if payload == nil {
		payload = make(map[string]string)
	}
	nextHop := route[0]
	if node.transportTrust(nextHop) < 0 {
		node.record("transport_path_refused", token.OutcomeRefused, "", node.Name+" refused to use low-trust TCP next hop "+nextHop+" for "+toNode)
		return WireResponse{Outcome: token.OutcomeRefused, Detail: "local transport trust too low for next hop " + nextHop}, nil
	}
	message, messageErr := node.newWireMessage(toNode, kind, route, payload, messageToken)
	if messageErr != nil {
		return WireResponse{}, messageErr
	}
	return node.transmit(nextHop, message)
}

func (node *Node) transmit(nextHop string, message WireMessage) (WireResponse, error) {
	address := nodeAddrs[nextHop]
	envelopeBytes, envelopeErr := message.EnvelopeBytes()
	if envelopeErr != nil {
		return WireResponse{}, envelopeErr
	}
	var lastErr error
	for attemptNumber := 0; attemptNumber < 20; attemptNumber++ {
		frameConn, dialErr := transport.DialFrameConn(address, 10*time.Second)
		if dialErr != nil {
			lastErr = dialErr
			time.Sleep(250 * time.Millisecond)
			continue
		}
		if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
			closeFrame(frameConn)
			node.recordTransportObservation(nextHop, token.OutcomeBroken, "write frame failed: "+writeErr.Error())
			return WireResponse{}, writeErr
		}
		responseBytes, readErr := frameConn.ReadFrame()
		closeFrame(frameConn)
		if readErr != nil {
			node.recordTransportObservation(nextHop, token.OutcomeBroken, "read response failed: "+readErr.Error())
			return WireResponse{}, readErr
		}
		var response WireResponse
		if decodeErr := json.Unmarshal(responseBytes, &response); decodeErr != nil {
			node.recordTransportObservation(nextHop, token.OutcomeBroken, "decode response failed: "+decodeErr.Error())
			return WireResponse{}, decodeErr
		}
		node.recordTransportObservation(nextHop, token.OutcomeKept, "tcp frame exchange completed for "+message.Kind)
		return response, nil
	}
	node.recordTransportObservation(nextHop, token.OutcomeBroken, "dial failed after retries")
	return WireResponse{}, fmt.Errorf("tcp frame to %s failed after retries: %w", address, lastErr)
}

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

func (node *Node) nextRouteHop(message WireMessage) (string, error) {
	if len(message.Route) == 0 {
		return "", fmt.Errorf("message from %s to %s has no route", message.FromNode, message.ToNode)
	}
	for routeIndex, routeNode := range message.Route {
		if routeNode == node.Name && routeIndex+1 < len(message.Route) {
			return message.Route[routeIndex+1], nil
		}
	}
	return "", fmt.Errorf("route %s has no promised next hop after %s toward %s", encodeRoute(message.Route), node.Name, message.ToNode)
}

func (node *Node) selectRoute(toNode string) ([]string, error) {
	if toNode == node.Name {
		return nil, fmt.Errorf("route from %s to itself is not an inter-peer promise", node.Name)
	}
	node.mutex.Lock()
	candidateRoutes := append([][]string(nil), node.Discovery.KnownRoutes[toNode]...)
	node.mutex.Unlock()
	if len(candidateRoutes) == 0 {
		return nil, fmt.Errorf("%s has no local route promise to %s", node.Name, toNode)
	}
	sort.SliceStable(candidateRoutes, func(leftIndex int, rightIndex int) bool {
		leftRoute := candidateRoutes[leftIndex]
		rightRoute := candidateRoutes[rightIndex]
		if len(leftRoute) != len(rightRoute) {
			return len(leftRoute) < len(rightRoute)
		}
		return encodeRoute(leftRoute) < encodeRoute(rightRoute)
	})
	for _, route := range candidateRoutes {
		if len(route) == 0 {
			continue
		}
		// Intent: A peer's route refusal is local evidence about that next hop for
		// later route selection; it is not global punishment or a command. Source:
		// DI-sipuz
		if !node.hasCompleted("route_refused_by:"+route[0]) && node.transportTrust(route[0]) >= 0 {
			return append([]string(nil), route...), nil
		}
	}
	return nil, fmt.Errorf("%s has only low-trust local routes to %s", node.Name, toNode)
}

func (node *Node) learnRoute(targetNode string, route []string, source string) bool {
	// Intent: Route promises become explicit local evidence that strategies can wait
	// on; direct mesh-neighbor routes stay as startup facts rather than discovery
	// evidence. Source: DI-sipuz
	if targetNode == "" || targetNode == node.Name || len(route) == 0 {
		return false
	}
	routeCopy := append([]string(nil), route...)
	node.mutex.Lock()
	for _, existingRoute := range node.Discovery.KnownRoutes[targetNode] {
		if encodeRoute(existingRoute) == encodeRoute(routeCopy) {
			node.mutex.Unlock()
			return false
		}
	}
	node.Discovery.KnownRoutes[targetNode] = append(node.Discovery.KnownRoutes[targetNode], routeCopy)
	node.mutex.Unlock()
	if source != "direct mesh neighbor" {
		node.markCompleted("route:"+targetNode, "learned local route to "+targetNode+" from "+source+" as "+encodeRoute(routeCopy))
	}
	return true
}

func (node *Node) validPromisedRoute(route []string) bool {
	previousNode := node.Name
	for _, routeNode := range route {
		if !node.Mesh.HasEdge(previousNode, routeNode) {
			return false
		}
		previousNode = routeNode
	}
	return true
}

func (node *Node) recordTransportObservation(peer string, outcome string, detail string) {
	if peer == "" {
		peer = "unknown"
	}
	observation := TransportObservation{Peer: peer, Outcome: outcome, Detail: detail}
	node.mutex.Lock()
	node.Discovery.TransportObservations = append(node.Discovery.TransportObservations, observation)
	switch outcome {
	case token.OutcomeKept:
		node.Discovery.TransportTrust[peer]++
	case token.OutcomeBroken:
		node.Discovery.TransportTrust[peer] -= 2
	case token.OutcomeRefused:
		node.Discovery.TransportTrust[peer]--
	}
	node.mutex.Unlock()
	node.record("transport_observed", outcome, "", node.Name+" observed transport behavior for "+peer+": "+detail)
}

func (node *Node) transportTrust(peer string) int {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	return node.Discovery.TransportTrust[peer]
}

func (node *Node) markCompleted(key string, detail string) {
	node.mutex.Lock()
	alreadyDone := node.Discovery.Completed[key]
	node.Discovery.Completed[key] = true
	node.mutex.Unlock()
	if !alreadyDone {
		node.record("local_condition_met", token.OutcomeKept, "", key+": "+detail)
	}
}

func (node *Node) hasCompleted(key string) bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	return node.Discovery.Completed[key]
}

func (node *Node) waitForCompletion(key string, duration time.Duration) bool {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		if node.hasCompleted(key) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	node.record("local_condition_timeout", token.OutcomeBroken, "", node.Name+" timed out waiting for "+key)
	return false
}

func (node *Node) writeNodeDoneIfComplete(keys []string) {
	missingKeys := []string{}
	for _, key := range keys {
		if !node.hasCompleted(key) {
			missingKeys = append(missingKeys, key)
		}
	}
	if len(missingKeys) > 0 {
		node.record("node_done_refused", token.OutcomeBroken, "", node.Name+" missing local completion evidence: "+strings.Join(missingKeys, ","))
		return
	}
	if doneErr := node.writeNodeDone(); doneErr != nil {
		node.record("node_done_write", token.OutcomeBroken, "", doneErr.Error())
	}
}

func (node *Node) writeNodeDone() error {
	if err := os.MkdirAll(filepath.Join(node.RunRoot, node.RunID), 0o755); err != nil {
		return err
	}
	marker := filepath.Join(node.RunRoot, node.RunID, node.Name+".done")
	if writeErr := os.WriteFile(marker, []byte(node.Name+"\n"), 0o644); writeErr != nil {
		return writeErr
	}
	node.record("node_done", token.OutcomeKept, "", node.Name+" wrote local done marker")
	return nil
}

func (node *Node) waitForAllDone() error {
	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) {
		missingNodes := []string{}
		for _, nodeName := range allNodeNames {
			marker := filepath.Join(node.RunRoot, node.RunID, nodeName+".done")
			if _, statErr := os.Stat(marker); statErr != nil {
				missingNodes = append(missingNodes, nodeName)
			}
		}
		if len(missingNodes) == 0 {
			node.record("all_nodes_done", token.OutcomeKept, "", node.Name+" observed all POC9 node done markers")
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("%s timed out waiting for all node done markers", node.Name)
}

func (node *Node) record(event string, outcome string, tokenID string, detail string) {
	record := token.Event{Observer: node.Name, Event: event, Outcome: outcome, TokenID: tokenID, Detail: detail}
	node.mutex.Lock()
	node.Evidence = append(node.Evidence, record)
	node.mutex.Unlock()
	bytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		log.Printf("encode evidence record: %v", marshalErr)
		return
	}
	fmt.Println(string(bytes))
}

func historicalAliceExpiredToken(tokenID string) (token.Token, error) {
	issuer := token.NewIssuer("alice")
	return issuer.IssueExpiring(tokenID, "mallory", data.Kind, "dataset-expired", token.TransferBearer, expiredTokenExpiresAtUnix)
}

func guessPeerFromMalformed(frameBytes []byte) string {
	frameText := string(frameBytes)
	for _, nodeName := range allNodeNames {
		if strings.Contains(frameText, "from_node="+nodeName) || strings.Contains(frameText, nodeName) {
			return nodeName
		}
	}
	return "unknown"
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
		"dataset-expired": "expired dataset",
	}
}

func parseFibonacciInput(value string) (int, error) {
	trimmedValue := strings.TrimPrefix(value, "fib-")
	inputNumber, parseErr := strconv.Atoi(trimmedValue)
	if parseErr != nil {
		return 0, parseErr
	}
	return inputNumber, nil
}

// fibonacci performs bounded compute work for compute-token redemption. Intent:
// Keep compute promises concrete and deterministic while staying cheap enough for
// repeated Docker runs. Source: DI-sipuz
func fibonacci(inputNumber int) (uint64, error) {
	if inputNumber < 0 || inputNumber > 93 {
		return 0, fmt.Errorf("fibonacci input %d is outside uint64 demo bounds", inputNumber)
	}
	if inputNumber == 0 {
		return 0, nil
	}
	var previousValue uint64
	currentValue := uint64(1)
	for sequenceIndex := 1; sequenceIndex < inputNumber; sequenceIndex++ {
		previousValue, currentValue = currentValue, previousValue+currentValue
	}
	return currentValue, nil
}

func atoi(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	parsedValue, parseErr := strconv.Atoi(value)
	if parseErr != nil {
		return 0
	}
	return parsedValue
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
