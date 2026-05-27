package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/auditor"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/compute"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/data"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/issuer"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/protocol"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/relay"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/storage"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/token"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/trader"
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

// appProtocolCID names the POC7 app-message protocol spec bytes.
// Intent: Make every app message dispatch through a pCID-selected CBOR grid
// envelope rather than through JSON field names. Source: DI-fibok
var appProtocolCID = protocol.NewProtocolCID([]byte("poc7 app-message protocol: cbor grid envelope with signed payload map and optional token bytes"))

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
}

var nodeAddrs = map[string]string{
	"alice":   "http://alice:8080",
	"bob":     "http://bob:8080",
	"carol":   "http://carol:8080",
	"dave":    "http://dave:8080",
	"mallory": "http://mallory:8080",
}

func main() {
	nodeName := flag.String("node", "", "node name")
	flag.Parse()
	if *nodeName == "" {
		log.Fatal("missing -node")
	}
	node := &Node{
		Name:    *nodeName,
		RunID:   getenv("POC7_RUN_ID", "manual"),
		RunRoot: "/run/poc7",
		Issuer:  token.NewIssuer(*nodeName),
		Wallet:  token.NewWallet(*nodeName),
		Storage: make(map[string]string),
		Data:    initialData(*nodeName),
	}
	server := &http.Server{Addr: ":8080", Handler: node.routes()}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("%s server failed: %v", node.Name, err)
		}
	}()
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
}

func (node *Node) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/relay", node.handleRelay)
	mux.HandleFunc("/script", node.handleScript)
	return mux
}

func (node *Node) handleRelay(response http.ResponseWriter, request *http.Request) {
	var message WireMessage
	if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := node.receive(message)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(response, result)
}

func (node *Node) handleScript(response http.ResponseWriter, request *http.Request) {
	var command WireMessage
	if err := json.NewDecoder(request.Body).Decode(&command); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := node.runCommand(command)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(response, result)
}

func (node *Node) receive(rawMessage WireMessage) (WireResponse, error) {
	message, exactBytes, decodeErr := node.decodeWireMessage(rawMessage)
	if decodeErr != nil {
		return WireResponse{}, decodeErr
	}
	node.record("kernel_receive", token.OutcomeKept, message.Token.ID, node.Name+" kernel observed "+message.Kind+" exact_sha256="+protocol.HashExactBytes(exactBytes))
	if message.ToNode != node.Name {
		return node.forward(message)
	}
	switch message.ToApp {
	case issuer.AppName:
		return node.handleIssuer(message)
	case trader.AppName:
		return node.handleTrader(message)
	case auditor.AppName:
		return node.handleAuditor(message)
	default:
		return WireResponse{}, fmt.Errorf("%s has no app %s", node.Name, message.ToApp)
	}
}

func (node *Node) forward(message WireMessage) (WireResponse, error) {
	if len(message.Route) == 0 {
		return WireResponse{}, fmt.Errorf("no route from %s to %s", node.Name, message.ToNode)
	}
	next := message.Route[0]
	message.Route = message.Route[1:]
	node.record("relay_forward", token.OutcomeKept, message.Token.ID, node.Name+" "+relay.AppName+" promised next hop "+next)
	return post(nodeAddrs[next]+"/relay", message)
}

func (node *Node) handleIssuer(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case "issue_token_v1":
		issued, err := node.Issuer.Issue(message.Payload["token_id"], message.Payload["original_peer"], message.Payload["resource_kind"], message.Payload["resource_id"], message.Payload["transfer_rule"])
		if err != nil {
			return WireResponse{}, err
		}
		node.record("token_issued", token.OutcomeKept, issued.ID, node.Name+" issued "+issued.TransferRule+" token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token issued", Token: issued}, nil
	case "revoke_token_v1":
		if err := node.Issuer.Revoke(message.Payload["token_id"], message.Payload["reason"]); err != nil {
			return WireResponse{}, err
		}
		node.record("token_revoked", token.OutcomeKept, message.Payload["token_id"], node.Name+" revoked token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token revoked"}, nil
	case "redeem_token_v1":
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
		return WireResponse{}, fmt.Errorf("issuer cannot handle %s", message.Kind)
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
		return nil, fmt.Errorf("storage redemption has unsupported op %q", message.Payload["op"])
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
	case "receive_token_v1":
		node.Wallet.Add(message.Token, message.FromNode+" transferred token")
		node.record("token_received", token.OutcomeKept, message.Token.ID, node.Name+" trader received token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "token received"}, nil
	case "quote_offer_v1":
		offer := node.Wallet.Quote(message.Payload["offered_issuer"], message.Payload["wanted_issuer"])
		payload := map[string]string{
			"offered_count": fmt.Sprintf("%d", offer.OfferedCount),
			"wanted_count":  fmt.Sprintf("%d", offer.WantedCount),
		}
		node.record("exchange_rate_quoted", token.OutcomeKept, "", node.Name+" quoted peer-local exchange rate")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "peer-local quote", Payload: payload}, nil
	case "trade_for_access_v1":
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
		return WireResponse{}, fmt.Errorf("trader cannot handle %s", message.Kind)
	}
}

func (node *Node) handleAuditor(message WireMessage) (WireResponse, error) {
	if message.Kind != "audit_token_v1" {
		return WireResponse{}, fmt.Errorf("auditor cannot handle %s", message.Kind)
	}
	node.Wallet.Add(message.Token, message.FromNode+" presented token for audit")
	redeem, err := node.send("alice", issuer.AppName, "redeem_token_v1", []string{"mallory", "alice"}, map[string]string{}, message.Token)
	if err != nil {
		return WireResponse{}, err
	}
	node.Wallet.ApplyRedemption(token.Event{Observer: node.Name, Event: "audit_redemption", Outcome: redeem.Outcome, TokenID: message.Token.ID, Detail: redeem.Detail})
	node.record("audit_completed", redeem.Outcome, message.Token.ID, node.Name+" audited token and updated local trust")
	return WireResponse{Outcome: redeem.Outcome, Detail: "audit completed: " + redeem.Detail}, nil
}

func (node *Node) runCommand(rawMessage WireMessage) (WireResponse, error) {
	message, _, decodeErr := node.decodeWireMessage(rawMessage)
	if decodeErr != nil {
		return WireResponse{}, decodeErr
	}
	switch message.Kind {
	case "redeem_held_token_v1":
		if !node.Wallet.Holds(message.Token.ID) {
			return WireResponse{}, fmt.Errorf("%s cannot redeem token %s it does not hold", node.Name, message.Token.ID)
		}
		return node.send(message.ToNode, issuer.AppName, "redeem_token_v1", message.Route, message.Payload, message.Token)
	case "trade_for_access_v1":
		if !node.Wallet.Holds(message.Token.ID) {
			return WireResponse{}, fmt.Errorf("%s cannot trade token %s it does not hold", node.Name, message.Token.ID)
		}
		response, tradeErr := node.send(message.ToNode, trader.AppName, "trade_for_access_v1", message.Route, message.Payload, message.Token)
		if tradeErr != nil {
			return WireResponse{}, tradeErr
		}
		if response.Outcome == token.OutcomeKept && !response.Token.IsZero() {
			node.Wallet.Add(response.Token, "trade returned non-transferable access token")
		}
		return response, nil
	case "quote_offer_v1":
		return node.receive(message)
	default:
		return WireResponse{}, fmt.Errorf("unknown script command %s", message.Kind)
	}
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
	if _, err := node.send("bob", trader.AppName, "receive_token_v1", []string{"bob"}, nil, aliceBearer); err != nil {
		return err
	}
	if _, err := node.send("bob", trader.AppName, "receive_token_v1", []string{"bob"}, nil, bobAccess); err != nil {
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
	if _, err := node.send("carol", trader.AppName, "receive_token_v1", []string{"bob", "carol"}, nil, storageStoreToken); err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, "receive_token_v1", []string{"bob", "carol"}, nil, storageReadToken); err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, "receive_token_v1", []string{"bob", "carol"}, nil, storageTradeToken); err != nil {
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
	if _, err := node.send("bob", trader.AppName, "receive_token_v1", []string{"bob"}, nil, computeToken); err != nil {
		return err
	}
	computeResult, err := node.requestHolderRedeem("bob", []string{"carol"}, computeToken, map[string]string{"n": "10"})
	if err != nil {
		return err
	}
	if computeResult.Outcome != token.OutcomeKept || computeResult.Payload["result"] != "55" {
		return fmt.Errorf("compute redemption outcome=%#v", computeResult)
	}
	daveQuoteCommand, err := node.newWireMessage("scenario", "dave", trader.AppName, "quote_offer_v1", nil, map[string]string{"offered_issuer": "bob", "wanted_issuer": "carol"}, token.Token{})
	if err != nil {
		return err
	}
	if _, err := post(nodeAddrs["dave"]+"/script", daveQuoteCommand); err != nil {
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
	revokeMessage, err := node.newWireMessage(issuer.AppName, "alice", issuer.AppName, "revoke_token_v1", nil, map[string]string{"token_id": revoked.ID, "reason": "broken promise history changed Alice's local willingness"}, token.Token{})
	if err != nil {
		return err
	}
	if _, err := node.receive(revokeMessage); err != nil {
		return err
	}
	if _, err := node.send("mallory", trader.AppName, "receive_token_v1", []string{"mallory"}, nil, revoked); err != nil {
		return err
	}
	audit, err := node.send("dave", auditor.AppName, "audit_token_v1", []string{"mallory", "dave"}, nil, revoked)
	if err != nil {
		return err
	}
	if audit.Outcome != token.OutcomeBroken {
		return fmt.Errorf("expected revoked token audit to break, got %#v", audit)
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
	response, err := node.send(issuerNode, issuer.AppName, "issue_token_v1", route, map[string]string{"token_id": id, "original_peer": originalPeer, "resource_kind": resourceKind, "resource_id": resourceID, "transfer_rule": transferRule}, token.Token{})
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
	command, commandErr := node.newWireMessage("scenario", issuerNode, trader.AppName, "redeem_held_token_v1", route, payload, tok)
	if commandErr != nil {
		return WireResponse{}, commandErr
	}
	return post(nodeAddrs[holderNode]+"/script", command)
}

// requestHolderTrade asks the holder to offer a bearer token to Alice.
// Intent: Fix the first POC7 semantic bug by making Carol offer Bob's bearer
// token and receive the resulting non-transferable Alice data token. Source: DI-fibok
func (node *Node) requestHolderTrade(holderNode string, route []string, tok token.Token, payload map[string]string) (WireResponse, error) {
	command, commandErr := node.newWireMessage("scenario", "alice", trader.AppName, "trade_for_access_v1", route, payload, tok)
	if commandErr != nil {
		return WireResponse{}, commandErr
	}
	return post(nodeAddrs[holderNode]+"/script", command)
}

// newWireMessage packages app fields into a signed CBOR grid envelope carried
// by the HTTP test harness. Intent: Keep HTTP as disposable container plumbing
// while making the protocol message itself be exact signed CBOR bytes.
// Source: DI-fibok
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

// decodeWireMessage verifies and opens the signed CBOR grid envelope carried by
// the HTTP harness. Intent: Treat the envelope bytes as the protocol message and
// the JSON wrapper as transport-only test plumbing. Source: DI-fibok
func (node *Node) decodeWireMessage(message WireMessage) (WireMessage, []byte, error) {
	if message.Envelope == "" {
		return message, nil, nil
	}
	envelopeBytes, hexErr := hex.DecodeString(message.Envelope)
	if hexErr != nil {
		return WireMessage{}, nil, hexErr
	}
	envelope, parseErr := protocol.ParseEnvelope(envelopeBytes)
	if parseErr != nil {
		return WireMessage{}, nil, parseErr
	}
	if !envelope.ProtocolCID.Equal(appProtocolCID) {
		return WireMessage{}, nil, fmt.Errorf("unsupported app envelope pCID %s", envelope.ProtocolCID)
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
	decoded.Payload = fields
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

// send hands the message to the first promised neighbor and leaves only the
// remaining hops inside the envelope. Intent: Make the trace show app-level
// relay promises between peers, not accidental self-forwards caused by keeping
// the already-consumed first hop in the route. Source: DI-tugih
func (node *Node) send(toNode string, toApp string, kind string, route []string, payload map[string]string, tok token.Token) (WireResponse, error) {
	if len(route) == 0 {
		return WireResponse{}, fmt.Errorf("empty route from %s to %s", node.Name, toNode)
	}
	if payload == nil {
		payload = make(map[string]string)
	}
	nextHop := route[0]
	remainingRoute := append([]string(nil), route[1:]...)
	message, messageErr := node.newWireMessage("scenario", toNode, toApp, kind, remainingRoute, payload, tok)
	if messageErr != nil {
		return WireResponse{}, messageErr
	}
	return post(nodeAddrs[nextHop]+"/relay", message)
}

func post(url string, message WireMessage) (WireResponse, error) {
	body, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return WireResponse{}, marshalErr
	}
	client := &http.Client{Timeout: 10 * time.Second}
	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		response, postErr := client.Post(url, "application/json", bytes.NewReader(body))
		if postErr != nil {
			lastErr = postErr
			time.Sleep(250 * time.Millisecond)
			continue
		}
		defer closeBody(response)
		if response.StatusCode >= 300 {
			return WireResponse{}, fmt.Errorf("post %s got %s", url, response.Status)
		}
		var result WireResponse
		if decodeErr := json.NewDecoder(response.Body).Decode(&result); decodeErr != nil {
			return WireResponse{}, decodeErr
		}
		return result, nil
	}
	return WireResponse{}, fmt.Errorf("post %s failed after retries: %w", url, lastErr)
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

func writeJSON(response http.ResponseWriter, value WireResponse) {
	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func closeBody(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Printf("close response body: %v", err)
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
