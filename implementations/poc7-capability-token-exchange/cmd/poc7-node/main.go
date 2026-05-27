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

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/compute"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/data"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/issuer"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/protocol"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/relay"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/storage"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/token"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/trader"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/transport"
)

const (
	kindResourcePromiseRequest         = "resource_promise_request_v1"
	kindPromiseRevocationNotice        = "promise_revocation_notice_v1"
	kindPromisePresentedForFulfillment = "promise_presented_for_fulfillment_v1"
	kindPromiseReceived                = "promise_received_v1"
	kindExchangeOfferPromise           = "exchange_offer_promise_v1"
	kindReciprocalExchangePromise      = "reciprocal_exchange_promise_v1"
	kindHeldPromiseFulfillmentRequest  = "held_promise_fulfillment_request_v1"
	kindHeldReciprocalExchangeRequest  = "held_reciprocal_exchange_request_v1"
	kindHeldPromiseTransferRequest     = "held_promise_transfer_request_v1"
)

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

type WireResponse struct {
	Outcome string            `json:"outcome"`
	Detail  string            `json:"detail"`
	Token   token.Token       `json:"token"`
	Payload map[string]string `json:"payload"`
}

// EnvelopeBytes returns the signed CBOR grid bytes carried by a local
// WireMessage value. Intent: Keep the in-process struct as test scaffolding
// while TCP peers exchange exact envelope bytes. Source: DI-tanat
func (message WireMessage) EnvelopeBytes() ([]byte, error) {
	if message.Envelope == "" {
		return nil, fmt.Errorf("wire message missing signed envelope")
	}
	return hex.DecodeString(message.Envelope)
}

// appProtocolCID names the POC7 app-message protocol spec bytes.
// Intent: Make every app message dispatch through a pCID-selected CBOR grid
// envelope rather than through JSON field names. Source: DI-fibok
var appProtocolCID = protocol.NewProtocolCID([]byte("poc7 app-message protocol: cbor grid envelope with signed route, promise-shaped kind, payload map, and optional token bytes"))

// Node groups one local kernel boundary, one app-level relay, and the resource
// apps that run inside a single container. Intent: Make each container's view
// local and promise-based while still producing one bounded, reproducible demo
// trace across five peers. Source: DI-tugih; DI-fibok
type Node struct {
	Name     string
	RunID    string
	RunRoot  string
	Issuer   *token.Issuer
	Wallet   *token.Wallet
	Evidence []token.Event
	Storage  map[string]string
	Data     map[string]string
	// TokenSources records which peer voluntarily circulated a held token to
	// this node. Intent: Dave judges Mallory's stale-token circulation locally
	// instead of receiving Alice-provided out-of-band evidence. Source: DI-pabot
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
	node := &Node{
		Name:         *nodeName,
		RunID:        getenv("POC7_RUN_ID", "manual"),
		RunRoot:      "/run/poc7",
		Issuer:       token.NewIssuer(*nodeName),
		Wallet:       token.NewWallet(*nodeName),
		Storage:      make(map[string]string),
		Data:         initialData(*nodeName),
		TokenSources: make(map[string]string),
	}
	listener, listenErr := net.Listen("tcp", ":8077")
	if listenErr != nil {
		log.Fatalf("%s listen failed: %v", node.Name, listenErr)
	}
	serverErrors := make(chan error, 1)
	go node.serveTCP(listener, serverErrors)
	if node.Name == "alice" {
		if err := node.runScenario(); err != nil {
			log.Printf("scenario failed: %v", err)
			os.Exit(1)
		}
	}
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

// serveTCP accepts bounded length-framed TCP messages from neighbor nodes.
// Intent: Match POC2 through POC5 transport discipline while keeping the signed
// CBOR envelope as the only app-protocol object. Source: DI-tanat
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

// handleTCPConnection reads one signed envelope frame and writes one local
// observation response frame. The response is run-control evidence for this POC,
// not a new global PromiseGrid response authority. Source: DI-tanat
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

func (node *Node) receive(rawMessage WireMessage) (WireResponse, error) {
	envelopeBytes, envelopeErr := rawMessage.EnvelopeBytes()
	if envelopeErr != nil {
		return WireResponse{}, envelopeErr
	}
	return node.receiveFrame(envelopeBytes)
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
	if isHeldPromiseInstruction(message.Kind) {
		return node.runCommandMessage(message)
	}
	switch message.ToApp {
	case issuer.AppName:
		return node.handleIssuer(message)
	case trader.AppName:
		return node.handleTrader(message)
	default:
		return WireResponse{}, fmt.Errorf("%s has no app %s", node.Name, message.ToApp)
	}
}

// isHeldPromiseInstruction identifies scenario messages that must be interpreted
// by the local holder before any issuer or exchange peer is contacted. Source:
// DI-tanat
func isHeldPromiseInstruction(kind string) bool {
	return kind == kindHeldPromiseFulfillmentRequest || kind == kindHeldReciprocalExchangeRequest || kind == kindHeldPromiseTransferRequest
}

func (node *Node) forward(message WireMessage) (WireResponse, error) {
	next, nextErr := node.nextRouteHop(message)
	if nextErr != nil {
		return WireResponse{}, nextErr
	}
	node.record("relay_forward", token.OutcomeKept, message.Token.ID, node.Name+" "+relay.AppName+" promised next hop "+next)
	return transmit(nodeAddrs[next], message)
}

// nextRouteHop finds the next peer after this node inside the signed route.
// Intent: Preserve one immutable route promise in the app message while each
// relay only promises its own next hop. Source: DI-tanat
func (node *Node) nextRouteHop(message WireMessage) (string, error) {
	if len(message.Route) == 0 {
		return "", fmt.Errorf("no route from %s to %s", node.Name, message.ToNode)
	}
	for index, routeNode := range message.Route {
		if routeNode == node.Name && index+1 < len(message.Route) {
			return message.Route[index+1], nil
		}
	}
	return "", fmt.Errorf("route %s has no promised next hop after %s toward %s", strings.Join(message.Route, ","), node.Name, message.ToNode)
}

func (node *Node) handleIssuer(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case kindResourcePromiseRequest:
		issued, err := node.Issuer.Issue(message.Payload["token_id"], message.Payload["original_peer"], message.Payload["resource_kind"], message.Payload["resource_id"], message.Payload["transfer_rule"])
		if err != nil {
			return WireResponse{}, err
		}
		node.record("token_issued", token.OutcomeKept, issued.ID, node.Name+" issued "+issued.TransferRule+" token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token issued", Token: issued}, nil
	case kindPromiseRevocationNotice:
		if err := node.Issuer.Revoke(message.Payload["token_id"], message.Payload["reason"]); err != nil {
			return WireResponse{}, err
		}
		node.record("token_revoked", token.OutcomeKept, message.Payload["token_id"], node.Name+" revoked token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token revoked"}, nil
	case kindPromisePresentedForFulfillment:
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
		}
		return WireResponse{Outcome: event.Outcome, Detail: event.Detail, Payload: payload}, nil
	default:
		return WireResponse{}, fmt.Errorf("issuer has no local promise handler for %s", message.Kind)
	}
}

// fulfillResource performs the app-level work promised by a redeemed token.
// Intent: Make POC7 redemption exercise actual storage, compute, and data
// payload behavior instead of only changing token status. Source: DI-fibok
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
	nText := message.Payload["n"]
	if nText == "" {
		nText = message.Token.ResourceID
	}
	n, parseErr := parseFibonacciN(nText)
	if parseErr != nil {
		return nil, parseErr
	}
	result, resultErr := fibonacci(n)
	if resultErr != nil {
		return nil, resultErr
	}
	return map[string]string{"resource_kind": compute.Kind, "n": fmt.Sprintf("%d", n), "result": fmt.Sprintf("%d", result)}, nil
}

func (node *Node) handleTrader(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case kindPromiseReceived:
		node.Wallet.Add(message.Token, message.FromNode+" transferred token")
		if message.FromNode != "" {
			node.TokenSources[message.Token.ID] = message.FromNode
		}
		node.record("token_received", token.OutcomeKept, message.Token.ID, node.Name+" trader received token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token received"}, nil
	case kindExchangeOfferPromise:
		offer := node.Wallet.Quote(message.Payload["offered_issuer"], message.Payload["wanted_issuer"])
		payload := map[string]string{
			"offered_count": fmt.Sprintf("%d", offer.OfferedCount),
			"wanted_count":  fmt.Sprintf("%d", offer.WantedCount),
		}
		node.record("exchange_rate_quoted", token.OutcomeKept, "", node.Name+" quoted peer-local exchange rate")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "peer-local quote", Payload: payload}, nil
	case kindReciprocalExchangePromise:
		node.Wallet.Add(message.Token, message.FromNode+" offered bearer token for non-transferable access")
		recipientPeer := message.Payload["recipient_peer"]
		if recipientPeer == "" {
			return WireResponse{}, fmt.Errorf("trade missing recipient_peer")
		}
		if message.Token.TransferRule != token.TransferBearer {
			return WireResponse{}, fmt.Errorf("trade offer token %s is not bearer-transferable", message.Token.ID)
		}
		if verifyErr := token.VerifyToken(message.Token); verifyErr != nil {
			return WireResponse{}, verifyErr
		}
		quote := node.Wallet.Quote(message.Token.Issuer, node.Name)
		issued, err := node.Issuer.Issue(message.Payload["new_token_id"], recipientPeer, message.Payload["resource_kind"], message.Payload["resource_id"], token.TransferNonTransferable)
		if err != nil {
			return WireResponse{}, err
		}
		node.record("trade_accepted", token.OutcomeKept, issued.ID, node.Name+" exchanged "+fmt.Sprintf("%d", quote.OfferedCount)+" "+message.Token.Issuer+" bearer token promise for non-transferable access token promised to "+recipientPeer)
		return WireResponse{Outcome: token.OutcomeKept, Detail: "trade accepted", Token: issued, Payload: map[string]string{"offered_count": fmt.Sprintf("%d", quote.OfferedCount), "wanted_count": fmt.Sprintf("%d", quote.WantedCount), "recipient_peer": recipientPeer}}, nil
	default:
		return WireResponse{}, fmt.Errorf("trader has no local promise handler for %s", message.Kind)
	}
}

// runCommand keeps local holder instructions on the same signed-envelope path
// as inter-peer messages. Source: DI-tanat
func (node *Node) runCommand(rawMessage WireMessage) (WireResponse, error) {
	message, _, decodeErr := node.decodeWireMessage(rawMessage)
	if decodeErr != nil {
		return WireResponse{}, decodeErr
	}
	return node.runCommandMessage(message)
}

// runCommandMessage turns a holder-local instruction into the holder's own
// outbound promise presentation or reciprocal exchange promise. Intent: Alice's
// scenario driver may ask Bob or Carol to act, but Bob or Carol still make their
// own local promise presentation. Source: DI-tanat
func (node *Node) runCommandMessage(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case kindHeldPromiseFulfillmentRequest:
		if !node.Wallet.Holds(message.Token.ID) {
			return WireResponse{}, fmt.Errorf("%s has no local held-promise evidence for token %s", node.Name, message.Token.ID)
		}
		issuerNode := message.Payload["promise_issuer_node"]
		if issuerNode == "" {
			return WireResponse{}, fmt.Errorf("held promise fulfillment request has no promise_issuer_node")
		}
		response, redeemErr := node.send(issuerNode, issuer.AppName, kindPromisePresentedForFulfillment, message.Route, message.Payload, message.Token)
		if redeemErr != nil {
			return WireResponse{}, redeemErr
		}
		node.applyRedemptionEvidence(message.Token, response)
		return response, nil
	case kindHeldReciprocalExchangeRequest:
		if !node.Wallet.Holds(message.Token.ID) {
			return WireResponse{}, fmt.Errorf("%s has no local held-promise evidence for exchange token %s", node.Name, message.Token.ID)
		}
		exchangePeer := message.Payload["exchange_peer_node"]
		if exchangePeer == "" {
			return WireResponse{}, fmt.Errorf("held reciprocal exchange request has no exchange_peer_node")
		}
		response, tradeErr := node.send(exchangePeer, trader.AppName, kindReciprocalExchangePromise, message.Route, message.Payload, message.Token)
		if tradeErr != nil {
			return WireResponse{}, tradeErr
		}
		if response.Outcome == token.OutcomeKept && !response.Token.IsZero() {
			node.Wallet.Add(response.Token, "trade returned non-transferable access token")
		}
		return response, nil
	case kindHeldPromiseTransferRequest:
		if !node.Wallet.Holds(message.Token.ID) {
			return WireResponse{}, fmt.Errorf("%s has no local held-promise evidence for transfer token %s", node.Name, message.Token.ID)
		}
		if message.Token.TransferRule != token.TransferBearer {
			return WireResponse{}, fmt.Errorf("held token %s is not bearer-transferable", message.Token.ID)
		}
		recipientNode := message.Payload["recipient_node"]
		if recipientNode == "" {
			return WireResponse{}, fmt.Errorf("held promise transfer request has no recipient_node")
		}
		transferPayload := copyPayload(message.Payload)
		transferPayload["presented_by_peer"] = node.Name
		response, transferErr := node.send(recipientNode, trader.AppName, kindPromiseReceived, message.Route, transferPayload, message.Token)
		if transferErr != nil {
			return WireResponse{}, transferErr
		}
		node.record("token_circulated", response.Outcome, message.Token.ID, node.Name+" voluntarily circulated token to "+recipientNode)
		return response, nil
	case kindExchangeOfferPromise:
		return node.receive(message)
	default:
		return WireResponse{}, fmt.Errorf("no local promise instruction named %s", message.Kind)
	}
}

// applyRedemptionEvidence updates the holder's local trust after the issuer
// keeps, breaks, or refuses a token promise. Intent: Dave judges both Alice's
// issuer promise and Mallory's stale-token circulation without any auditor
// authority or Alice-to-Dave evidence shortcut. Source: DI-pabot
func (node *Node) applyRedemptionEvidence(heldToken token.Token, response WireResponse) {
	event := token.Event{Observer: node.Name, Event: "held_redemption_observed", Outcome: response.Outcome, TokenID: heldToken.ID, Detail: response.Detail}
	node.Wallet.ApplyRedemption(event)
	node.record("holder_trust_updated", response.Outcome, heldToken.ID, node.Name+" updated local trust for issuer "+heldToken.Issuer+" after redemption outcome "+response.Outcome)
	sourcePeer := node.TokenSources[heldToken.ID]
	if sourcePeer == "" || sourcePeer == heldToken.Issuer {
		return
	}
	node.Wallet.ApplyPeerObservation(sourcePeer, event)
	node.record("circulator_trust_updated", response.Outcome, heldToken.ID, node.Name+" updated local trust for circulating peer "+sourcePeer+" after redemption outcome "+response.Outcome)
}

// runScenario is Alice's bounded script for the five-agent exchange. Intent:
// Exercise bearer transfer, non-transferable access, peer-local exchange
// quotes, revoked promises, and local trust updates without a central exchange
// or global authority. Source: DI-tugih
func (node *Node) runScenario() error {
	time.Sleep(2 * time.Second)
	aliceBearer, err := node.localIssue("alice-bearer-1", "bob", data.Kind, "dataset-public", token.TransferBearer)
	if err != nil {
		return err
	}
	bobAccess, err := node.localIssue("alice-data-bob-1", "bob", data.Kind, "dataset-private", token.TransferNonTransferable)
	if err != nil {
		return err
	}
	revoked, err := node.localIssue("alice-revoked-1", "mallory", data.Kind, "dataset-stale", token.TransferBearer)
	if err != nil {
		return err
	}
	if _, err := node.send("bob", trader.AppName, kindPromiseReceived, []string{"bob"}, nil, aliceBearer); err != nil {
		return err
	}
	if _, err := node.send("bob", trader.AppName, kindPromiseReceived, []string{"bob"}, nil, bobAccess); err != nil {
		return err
	}
	if redeem, err := node.requestHolderRedeem("bob", []string{"alice"}, bobAccess, map[string]string{}); err != nil || redeem.Outcome != token.OutcomeKept {
		return fmt.Errorf("bob redeem outcome=%#v err=%v", redeem, err)
	}
	storageStoreToken, err := node.requestIssue("bob", []string{"bob"}, "bob-storage-carol-store-1", "carol", storage.Kind, "storage-slot-store", token.TransferBearer)
	if err != nil {
		return err
	}
	storageReadToken, err := node.requestIssue("bob", []string{"bob"}, "bob-storage-carol-read-1", "carol", storage.Kind, "storage-slot-read", token.TransferBearer)
	if err != nil {
		return err
	}
	storageTradeToken, err := node.requestIssue("bob", []string{"bob"}, "bob-storage-carol-trade-1", "carol", storage.Kind, "storage-slot-trade", token.TransferBearer)
	if err != nil {
		return err
	}
	computeToken, err := node.requestIssue("carol", []string{"bob", "carol"}, "carol-compute-bob-1", "bob", compute.Kind, "fib-55", token.TransferBearer)
	if err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, kindPromiseReceived, []string{"bob", "carol"}, nil, storageStoreToken); err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, kindPromiseReceived, []string{"bob", "carol"}, nil, storageReadToken); err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, kindPromiseReceived, []string{"bob", "carol"}, nil, storageTradeToken); err != nil {
		return err
	}
	storeResult, err := node.requestHolderRedeem("carol", []string{"bob"}, storageStoreToken, map[string]string{"op": "store", "key": "carol-note", "value": "stored through Bob promise token"})
	if err != nil {
		return err
	}
	if storeResult.Outcome != token.OutcomeKept || storeResult.Payload["stored"] != "true" {
		return fmt.Errorf("storage store redemption outcome=%#v", storeResult)
	}
	readResult, err := node.requestHolderRedeem("carol", []string{"bob"}, storageReadToken, map[string]string{"op": "read", "key": "carol-note"})
	if err != nil {
		return err
	}
	if readResult.Outcome != token.OutcomeKept || readResult.Payload["value"] != "stored through Bob promise token" {
		return fmt.Errorf("storage read redemption outcome=%#v", readResult)
	}
	if _, err := node.send("bob", trader.AppName, kindPromiseReceived, []string{"bob"}, nil, computeToken); err != nil {
		return err
	}
	computeResult, err := node.requestHolderRedeem("bob", []string{"carol"}, computeToken, map[string]string{"n": "10"})
	if err != nil {
		return err
	}
	if computeResult.Outcome != token.OutcomeKept || computeResult.Payload["result"] != "55" {
		return fmt.Errorf("compute redemption outcome=%#v", computeResult)
	}
	daveQuoteCommand, err := node.newWireMessage("scenario", "dave", trader.AppName, kindExchangeOfferPromise, nil, map[string]string{"offered_issuer": "bob", "wanted_issuer": "carol"}, token.Token{})
	if err != nil {
		return err
	}
	if _, err := transmit(nodeAddrs["dave"], daveQuoteCommand); err != nil {
		return err
	}
	trade, err := node.requestHolderTrade("carol", []string{"bob", "alice"}, storageTradeToken, map[string]string{"new_token_id": "alice-data-carol-1", "recipient_peer": "carol", "resource_kind": data.Kind, "resource_id": "dataset-private"})
	if err != nil || trade.Outcome != token.OutcomeKept {
		return fmt.Errorf("trade outcome=%#v err=%v", trade, err)
	}
	if trade.Token.OriginalPeer != "carol" {
		return fmt.Errorf("trade issued access token to %s, want carol", trade.Token.OriginalPeer)
	}
	carolDataResult, err := node.requestHolderRedeem("carol", []string{"bob", "alice"}, trade.Token, map[string]string{})
	if err != nil {
		return err
	}
	if carolDataResult.Outcome != token.OutcomeKept || carolDataResult.Payload["value"] == "" {
		return fmt.Errorf("carol data redemption outcome=%#v", carolDataResult)
	}
	revokeMessage, err := node.newWireMessage(issuer.AppName, "alice", issuer.AppName, kindPromiseRevocationNotice, nil, map[string]string{"token_id": revoked.ID, "reason": "broken promise history changed Alice's local willingness"}, token.Token{})
	if err != nil {
		return err
	}
	if _, err := node.receive(revokeMessage); err != nil {
		return err
	}
	if _, err := node.send("mallory", trader.AppName, kindPromiseReceived, []string{"mallory"}, nil, revoked); err != nil {
		return err
	}
	if _, err := node.requestHolderTransfer("mallory", []string{"dave"}, "dave", revoked); err != nil {
		return err
	}
	staleRedemption, err := node.requestHolderRedeem("dave", []string{"mallory", "alice"}, revoked, map[string]string{})
	if err != nil {
		return err
	}
	if staleRedemption.Outcome != token.OutcomeBroken {
		return fmt.Errorf("expected revoked token redemption to break, got %#v", staleRedemption)
	}
	if err := os.MkdirAll(filepath.Join(node.RunRoot, node.RunID), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(node.RunRoot, node.RunID, "done"), []byte("done\n"), 0o644)
}

func (node *Node) localIssue(id string, originalPeer string, resourceKind string, resourceID string, transferRule string) (token.Token, error) {
	issued, err := node.Issuer.Issue(id, originalPeer, resourceKind, resourceID, transferRule)
	if err == nil {
		node.record("token_issued", token.OutcomeKept, id, "alice local issue")
	}
	return issued, err
}

func (node *Node) requestIssue(issuerNode string, route []string, id string, originalPeer string, resourceKind string, resourceID string, transferRule string) (token.Token, error) {
	response, err := node.send(issuerNode, issuer.AppName, kindResourcePromiseRequest, route, map[string]string{"token_id": id, "original_peer": originalPeer, "resource_kind": resourceKind, "resource_id": resourceID, "transfer_rule": transferRule}, token.Token{})
	return response.Token, err
}

// requestHolderRedeem asks the token holder's local process to redeem its token.
// Intent: Keep redemption initiated by the holder agent instead of having Alice
// impersonate Bob or Carol during the scripted scenario. Source: DI-fibok
func (node *Node) requestHolderRedeem(holderNode string, route []string, tok token.Token, payload map[string]string) (WireResponse, error) {
	if len(route) == 0 {
		return WireResponse{}, fmt.Errorf("redeem route for holder %s is empty", holderNode)
	}
	issuerNode := route[len(route)-1]
	commandPayload := copyPayload(payload)
	commandPayload["promise_issuer_node"] = issuerNode
	command, commandErr := node.newWireMessage("scenario", holderNode, trader.AppName, kindHeldPromiseFulfillmentRequest, route, commandPayload, tok)
	if commandErr != nil {
		return WireResponse{}, commandErr
	}
	return transmit(nodeAddrs[holderNode], command)
}

// requestHolderTrade asks the holder to offer a bearer token to Alice.
// Intent: Fix the first POC7 semantic bug by making Carol offer Bob's bearer
// token and receive the resulting non-transferable Alice data token. Source: DI-fibok
func (node *Node) requestHolderTrade(holderNode string, route []string, tok token.Token, payload map[string]string) (WireResponse, error) {
	if len(route) == 0 {
		return WireResponse{}, fmt.Errorf("exchange route for holder %s is empty", holderNode)
	}
	commandPayload := copyPayload(payload)
	commandPayload["exchange_peer_node"] = route[len(route)-1]
	command, commandErr := node.newWireMessage("scenario", holderNode, trader.AppName, kindHeldReciprocalExchangeRequest, route, commandPayload, tok)
	if commandErr != nil {
		return WireResponse{}, commandErr
	}
	return transmit(nodeAddrs[holderNode], command)
}

// requestHolderTransfer asks the current holder to circulate a bearer token to a
// chosen peer. Intent: Mallory must actively present the stale token to Dave;
// Alice no longer sends Dave a second copy out-of-band. Source: DI-pabot
func (node *Node) requestHolderTransfer(holderNode string, route []string, recipientNode string, tok token.Token) (WireResponse, error) {
	if len(route) == 0 {
		return WireResponse{}, fmt.Errorf("transfer route for holder %s is empty", holderNode)
	}
	commandPayload := map[string]string{"recipient_node": recipientNode}
	command, commandErr := node.newWireMessage("scenario", holderNode, trader.AppName, kindHeldPromiseTransferRequest, route, commandPayload, tok)
	if commandErr != nil {
		return WireResponse{}, commandErr
	}
	return transmit(nodeAddrs[holderNode], command)
}

// copyPayload gives local scenario instructions their own map before adding
// holder-only promise routing hints. Source: DI-tanat
func copyPayload(payload map[string]string) map[string]string {
	copiedPayload := make(map[string]string)
	for key, value := range payload {
		copiedPayload[key] = value
	}
	return copiedPayload
}

// encodeRoute stores the demo route in the signed payload as a compact string.
// Node names in this POC are fixed labels without commas. Source: DI-tanat
func encodeRoute(route []string) string {
	return strings.Join(route, ",")
}

// decodeRoute rebuilds the signed route list used for local relay promises.
// Source: DI-tanat
func decodeRoute(routeText string) []string {
	if strings.TrimSpace(routeText) == "" {
		return nil
	}
	return strings.Split(routeText, ",")
}

// newWireMessage packages app fields into a signed CBOR grid envelope. Intent:
// Make route, promise kind, promiser-facing payload, and optional token bytes
// part of the exact pCID-selected message carried over TCP. Source: DI-tanat
func (node *Node) newWireMessage(fromApp string, toNode string, toApp string, kind string, route []string, payload map[string]string, tok token.Token) (WireMessage, error) {
	fields := make(map[string]string)
	for key, value := range payload {
		fields[key] = value
	}
	fields["from_node"] = node.Name
	fields["from_app"] = fromApp
	fields["to_node"] = toNode
	fields["to_app"] = toApp
	fields["kind"] = kind
	fields["route"] = encodeRoute(route)
	if !tok.IsZero() {
		tokenBytes, tokenErr := token.Encode(tok)
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
	copiedRoute := append([]string(nil), route...)
	return WireMessage{
		FromNode: node.Name,
		FromApp:  fromApp,
		ToNode:   toNode,
		ToApp:    toApp,
		Kind:     kind,
		Route:    copiedRoute,
		Payload:  fields,
		Token:    tok,
		Envelope: hex.EncodeToString(envelopeBytes),
	}, nil
}

// decodeWireMessage verifies and opens a locally constructed signed CBOR grid
// envelope. Intent: Let scenario code use the same exact bytes that TCP peers
// exchange instead of a parallel in-process shortcut. Source: DI-tanat
func (node *Node) decodeWireMessage(message WireMessage) (WireMessage, []byte, error) {
	envelopeBytes, envelopeErr := message.EnvelopeBytes()
	if envelopeErr != nil {
		return WireMessage{}, nil, envelopeErr
	}
	return node.decodeEnvelopeBytes(envelopeBytes, message)
}

// decodeEnvelopeBytes verifies and opens the signed CBOR grid bytes received
// from TCP. Intent: Treat the pCID-selected envelope as the app message and TCP
// framing as non-semantic plumbing. Source: DI-tanat
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

// send hands the signed envelope to the first promised neighbor while keeping
// the full route inside the signed payload. Intent: Let each relay make its
// local next-hop promise without rewriting the app message. Source: DI-tanat
func (node *Node) send(toNode string, toApp string, kind string, route []string, payload map[string]string, tok token.Token) (WireResponse, error) {
	if len(route) == 0 {
		return WireResponse{}, fmt.Errorf("empty route from %s to %s", node.Name, toNode)
	}
	if payload == nil {
		payload = make(map[string]string)
	}
	nextHop := route[0]
	message, messageErr := node.newWireMessage("scenario", toNode, toApp, kind, route, payload, tok)
	if messageErr != nil {
		return WireResponse{}, messageErr
	}
	return transmit(nodeAddrs[nextHop], message)
}

// transmit sends one signed envelope over one length-framed TCP connection and
// reads one bounded local observation response for the POC scenario driver.
// Intent: Remove HTTP while avoiding any claim that request/response transport
// is the final PromiseGrid app protocol. Source: DI-tanat
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

func (node *Node) waitForDone() error {
	donePath := filepath.Join(node.RunRoot, node.RunID, "done")
	for deadline := time.Now().Add(45 * time.Second); time.Now().Before(deadline); {
		if _, err := os.Stat(donePath); err == nil {
			node.record("node_complete", token.OutcomeKept, "", node.Name+" observed bounded demo completion")
			marker := filepath.Join(node.RunRoot, node.RunID, node.Name+".complete")
			return os.WriteFile(marker, []byte(node.Name+"\n"), 0o644)
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("%s timed out waiting for %s", node.Name, donePath)
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

// initialData gives Alice actual data resources that redemption can return.
// Intent: Keep data access as Alice's local promise about resources she controls
// instead of a symbolic token label. Source: DI-fibok
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

func parseFibonacciN(value string) (int, error) {
	trimmed := strings.TrimPrefix(value, "fib-")
	n, parseErr := strconv.Atoi(trimmed)
	if parseErr != nil {
		return 0, parseErr
	}
	return n, nil
}

// fibonacci performs bounded compute work for the compute-token redemption.
// Intent: Make compute redemption return a real deterministic result while
// keeping the demo cheap enough for every container run. Source: DI-fibok
func fibonacci(n int) (uint64, error) {
	if n < 0 || n > 93 {
		return 0, fmt.Errorf("fibonacci n %d is outside uint64 demo bounds", n)
	}
	if n == 0 {
		return 0, nil
	}
	var previous uint64
	current := uint64(1)
	for index := 1; index < n; index++ {
		previous, current = current, previous+current
	}
	return current, nil
}
