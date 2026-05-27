package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/auditor"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/compute"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/data"
	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/issuer"
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
}

type WireResponse struct {
	Outcome string            `json:"outcome"`
	Detail  string            `json:"detail"`
	Token   token.Token       `json:"token"`
	Payload map[string]string `json:"payload"`
}

// Node groups one local kernel boundary, one app-level relay, and the resource
// apps that run inside a single container. Intent: Make each container's view
// local and promise-based while still producing one bounded, reproducible demo
// trace across five peers. Source: DI-tugih
type Node struct {
	Name     string
	RunID    string
	RunRoot  string
	Issuer   *token.Issuer
	Wallet   *token.Wallet
	Evidence []token.Event
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

func (node *Node) receive(message WireMessage) (WireResponse, error) {
	node.record("kernel_receive", token.OutcomeKept, message.Token.ID, node.Name+" kernel observed "+message.Kind)
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
		return WireResponse{Outcome: event.Outcome, Detail: event.Detail}, nil
	default:
		return WireResponse{}, fmt.Errorf("issuer cannot handle %s", message.Kind)
	}
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
		issued, err := node.Issuer.Issue(message.Payload["new_token_id"], message.FromNode, message.Payload["resource_kind"], message.Payload["resource_id"], token.TransferNonTransferable)
		if err != nil {
			return WireResponse{}, err
		}
		node.record("trade_accepted", token.OutcomeKept, issued.ID, node.Name+" exchanged bearer token for non-transferable access token")
		return WireResponse{Outcome: token.OutcomeKept, Detail: "trade accepted", Token: issued}, nil
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

func (node *Node) runCommand(message WireMessage) (WireResponse, error) {
	switch message.Kind {
	case "redeem_held_token_v1":
		return node.send(message.ToNode, issuer.AppName, "redeem_token_v1", message.Route, map[string]string{}, message.Token)
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
	if redeem, err := post(nodeAddrs["bob"]+"/script", WireMessage{FromNode: "bob", FromApp: trader.AppName, ToNode: "alice", ToApp: issuer.AppName, Kind: "redeem_held_token_v1", Route: []string{"alice"}, Token: bobAccess}); err != nil || redeem.Outcome != token.OutcomeKept {
		return fmt.Errorf("bob redeem outcome=%#v err=%v", redeem, err)
	}
	storageToken, err := node.requestIssue("bob", []string{"bob"}, "bob-storage-carol-1", "carol", storage.Kind, "storage-slot", token.TransferBearer)
	if err != nil {
		return err
	}
	computeToken, err := node.requestIssue("carol", []string{"bob", "carol"}, "carol-compute-bob-1", "bob", compute.Kind, "fib-55", token.TransferBearer)
	if err != nil {
		return err
	}
	if _, err := node.send("carol", trader.AppName, "receive_token_v1", []string{"bob", "carol"}, nil, storageToken); err != nil {
		return err
	}
	if _, err := node.send("bob", trader.AppName, "receive_token_v1", []string{"bob"}, nil, computeToken); err != nil {
		return err
	}
	if _, err := post(nodeAddrs["dave"]+"/script", WireMessage{FromNode: "dave", FromApp: trader.AppName, ToNode: "dave", ToApp: trader.AppName, Kind: "quote_offer_v1", Payload: map[string]string{"offered_issuer": "bob", "wanted_issuer": "carol"}}); err != nil {
		return err
	}
	if trade, err := node.send("alice", trader.AppName, "trade_for_access_v1", []string{"alice"}, map[string]string{"new_token_id": "alice-data-carol-1", "resource_kind": data.Kind, "resource_id": "dataset-private"}, storageToken); err != nil || trade.Outcome != token.OutcomeKept {
		return fmt.Errorf("trade outcome=%#v err=%v", trade, err)
	}
	if _, err := node.receive(WireMessage{FromNode: "alice", FromApp: issuer.AppName, ToNode: "alice", ToApp: issuer.AppName, Kind: "revoke_token_v1", Payload: map[string]string{"token_id": revoked.ID, "reason": "broken promise history changed Alice's local willingness"}}); err != nil {
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
	message := WireMessage{FromNode: node.Name, FromApp: "scenario", ToNode: toNode, ToApp: toApp, Kind: kind, Route: remainingRoute, Payload: payload, Token: tok}
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
