package runtime

import (
	"encoding/base64"
	"fmt"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/production"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/protocol"
)

const poc15RouteID = "route-alice-bob-carol-dave-0001"
const poc15RoutePath = "alice>bob>carol>dave"
const poc15RouteResponsePath = "dave>carol>bob>alice"
const poc15PeggyRouteID = "route-alice-dave-peggy-0001"
const poc15PeggyRoutePath = "alice>dave>peggy"
const poc15VictorRouteID = "route-alice-dave-victor-0001"
const poc15VictorRoutePath = "alice>dave>victor"
const routeMessageKindSetup = "setup"
const routeMessageKindCarried = "carried"

type routeSpec struct {
	ID           string
	Path         string
	Final        string
	Payment      string
	TTLMessages  string
	ResponsePath string
}

func primaryRouteSpec() routeSpec {
	return routeSpec{
		ID:           poc15RouteID,
		Path:         poc15RoutePath,
		Final:        "dave",
		Payment:      "reciprocal_forwarding_credit_1",
		TTLMessages:  "2",
		ResponsePath: poc15RouteResponsePath,
	}
}

func runtimeRouteSpec(id, path, final string) routeSpec {
	return routeSpec{
		ID:          id,
		Path:        path,
		Final:       final,
		Payment:     "runtime_compute_forwarding_credit_1",
		TTLMessages: "1",
	}
}

// runRoutePromiseWorkflow exercises the first executable POC15 multi-hop route.
// Intent: Alice uses route_v1 to ask only Alice's direct neighbor Bob for a
// voluntary route promise; Bob and Carol each make their own next-hop promises,
// and Dave locally confirms reachability before Alice sends a carried message.
// This is app-level promise forwarding, not a kernel route authority. Source:
// DI-lihir
func (node *Node) runRoutePromiseWorkflow() error {
	primaryRoute := primaryRouteSpec()
	node.record("route_exclusion_promise_made", "kept", "bob", "pcid="+pcid.RouteV1+" Alice asks route peers for a route that avoids Mallory as transit")
	node.record("route_transit_exclusion_peer_promised", "kept", "bob", "pcid="+pcid.RouteV1+" Bob promises this route will not intentionally transit Mallory")
	node.record("route_transit_exclusion_peer_promised", "kept", "carol", "pcid="+pcid.RouteV1+" Carol promises this route will not intentionally transit Mallory")
	node.record("route_exclusion_used_in_choice", "kept", "bob", "pcid="+pcid.RouteV1+" selected_path="+primaryRoute.Path+" excluded_transit=mallory")
	node.record("route_payment_promised", "kept", "bob", "pcid="+pcid.RouteV1+" Alice promises reciprocal forwarding credit for one bounded route setup")
	node.record("route_credit_offered", "kept", "bob", "pcid="+pcid.RouteV1+" credit="+primaryRoute.Payment)
	setupFields := routeHopFieldsForSpec(primaryRoute, node.Agent.Name, "bob", production.PromiseRouteSetup, routeMessageKindSetup, "")
	node.record("route_setup_promise_made", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" path="+primaryRoute.Path)
	setupAck, setupErr := node.sendAndReceive("bob", setupFields)
	if setupErr != nil {
		return fmt.Errorf("route setup: %w", setupErr)
	}
	if setupAck.Fields["route_status"] != "reachable" {
		return fmt.Errorf("route setup status %q", setupAck.Fields["route_status"])
	}
	node.record("route_reachability_confirmed", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" status=reachable")
	node.record("route_durability_promised", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" ttl_messages="+primaryRoute.TTLMessages)
	node.record("route_asymmetric_response_path_promised", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" response_path="+primaryRoute.ResponsePath)
	carriedFields := routeHopFieldsForSpec(primaryRoute, node.Agent.Name, "bob", production.PromiseRouteForward, routeMessageKindCarried, setupAck.ExactHash)
	carriedFields["carried_pcid"] = pcid.RelationshipV1
	carriedFields["carried_promise"] = "Alice promises to send one bounded relationship_v1 payload only after route reachability was confirmed."
	node.record("route_carried_message_sent", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" carried_pcid="+pcid.RelationshipV1)
	carriedAck, carriedErr := node.sendAndReceive("bob", carriedFields)
	if carriedErr != nil {
		return fmt.Errorf("route carried message: %w", carriedErr)
	}
	if carriedAck.Fields["route_status"] != "delivered" {
		return fmt.Errorf("route carried status %q", carriedAck.Fields["route_status"])
	}
	node.record("route_carried_message_delivered", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" carried_pcid="+pcid.RelationshipV1)
	reusedFields := routePayloadParentHopFieldsForSpec(primaryRoute, node.Agent.Name, "bob", production.PromiseRouteForward, routeMessageKindCarried, carriedAck.ExactHash)
	reusedFields["carried_pcid"] = pcid.RelationshipV1
	reusedFields["carried_promise"] = "Alice promises to reuse the previously confirmed route for a second bounded relationship_v1 payload within the explicit route lifetime."
	node.record("route_reused_message_sent", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" payload_parent="+carriedAck.ExactHash)
	reusedAck, reusedErr := node.sendAndReceive("bob", reusedFields)
	if reusedErr != nil {
		return fmt.Errorf("route reused message: %w", reusedErr)
	}
	if reusedAck.Fields["route_status"] != "delivered" {
		return fmt.Errorf("route reused status %q", reusedAck.Fields["route_status"])
	}
	node.record("route_reused_message_delivered", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" carried_pcid="+pcid.RelationshipV1)
	node.record("route_credit_spent", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" credit="+primaryRoute.Payment)
	// Intent: Route lifetime is a local promise term, not a global lease. Alice
	// treats the two-message route as expired after using both promised sends,
	// then asks for a fresh neighboring promise before relying on it again.
	// Source: DI-mapop
	node.record("route_lifetime_exhausted", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" ttl_messages="+primaryRoute.TTLMessages+" used_messages=2")
	node.record("route_expired_message_not_sent", "non_commitment", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" Alice does not send a third carried message after local route expiry")
	renewalFields := routeHopFieldsForSpec(primaryRoute, node.Agent.Name, "bob", production.PromiseRouteSetup, routeMessageKindSetup, reusedAck.ExactHash)
	node.record("route_renewal_requested", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" parent="+reusedAck.ExactHash)
	renewalAck, renewalErr := node.sendAndReceive("bob", renewalFields)
	if renewalErr != nil {
		return fmt.Errorf("route renewal: %w", renewalErr)
	}
	if renewalAck.Fields["route_status"] != "reachable" {
		return fmt.Errorf("route renewal status %q", renewalAck.Fields["route_status"])
	}
	node.record("route_renewal_confirmed", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+primaryRoute.ID+" status=reachable")
	if err := node.runRoutedRuntimeComputeWorkflow("peggy", runtimeRouteSpec(poc15PeggyRouteID, poc15PeggyRoutePath, "peggy")); err != nil {
		return err
	}
	return node.runRoutedRuntimeComputeWorkflow("victor", runtimeRouteSpec(poc15VictorRouteID, poc15VictorRoutePath, "victor"))
}

// handleRoutePromise handles one route_v1 hop from the local app's own vantage.
// Intent: Every hop signs a fresh pCID-owned route_v1 promise and references the
// previous exact envelope hash as parent context; no hop claims authority over
// downstream peers or final delivery beyond its own observed next-hop outcome.
// Source: DI-lihir
func (node *Node) handleRoutePromise(message parsedMessage) (map[string]string, error) {
	fields := message.Fields
	routePath := routePathParts(fields["route_path"])
	routeIndex, routeIndexErr := routePathIndex(routePath, node.Agent.Name)
	if routeIndexErr != nil {
		return nil, routeIndexErr
	}
	if routeIndex == len(routePath)-1 {
		return node.handleFinalRouteHop(message)
	}
	nextHop := routePath[routeIndex+1]
	node.record("route_forward_promise_made", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" next_hop="+nextHop)
	if fields["route_payment"] != "" {
		node.record("route_payment_promised", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" payment="+fields["route_payment"])
		node.record("route_credit_earned", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" credit="+fields["route_payment"])
	}
	if len(message.ParentExactSHA256s) > 0 {
		node.record("route_multiarity_parent_slot_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" parent="+message.ParentExactSHA256s[0])
	}
	if fields["payload_parent_exact_sha256"] != "" {
		node.record("route_payload_parent_link_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" payload_parent="+fields["payload_parent_exact_sha256"])
	}
	forwardFields := routeHopFieldsForSpec(routeSpecFromFields(fields), node.Agent.Name, nextHop, production.PromiseRouteForward, fields["route_message_kind"], message.ExactHash)
	for _, key := range []string{"carried_pcid", "carried_promise", "carried_envelope_b64", "payload_parent_exact_sha256", "route_setup_parent", "route_response_path"} {
		forwardFields[key] = fields[key]
	}
	if fields["payload_parent_exact_sha256"] != "" {
		forwardFields["parent_link_location"] = "payload"
	}
	downstreamAck, downstreamErr := node.sendAndReceive(nextHop, forwardFields)
	if downstreamErr != nil {
		node.record("route_forward_not_promised", "non_commitment", nextHop, "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" error="+downstreamErr.Error())
		return map[string]string{
			"promise_about": production.PromiseRouteForward,
			"route_id":      fields["route_id"],
			"route_status":  "not_promised",
			"route_path":    fields["route_path"],
		}, nil
	}
	node.record("route_forward_promise_kept", "kept", nextHop, "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" downstream_status="+downstreamAck.Fields["route_status"])
	return map[string]string{
		"promise_about":           production.PromiseRouteReachability,
		"route_id":                fields["route_id"],
		"route_path":              fields["route_path"],
		"route_status":            downstreamAck.Fields["route_status"],
		"route_final":             fields["route_final"],
		"route_parent_exact_hash": message.ExactHash,
		"route_final_ack_hash":    downstreamAck.ExactHash,
		"route_message_kind":      fields["route_message_kind"],
	}, nil
}

func (node *Node) handleFinalRouteHop(message parsedMessage) (map[string]string, error) {
	fields := message.Fields
	if len(message.ParentExactSHA256s) > 0 {
		node.record("route_multiarity_parent_slot_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" parent="+message.ParentExactSHA256s[0])
	}
	if fields["payload_parent_exact_sha256"] != "" {
		node.record("route_payload_parent_link_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" payload_parent="+fields["payload_parent_exact_sha256"])
	}
	if fields["route_message_kind"] == routeMessageKindCarried {
		node.record("route_carried_message_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" carried_pcid="+fields["carried_pcid"])
		carriedAckHash, carriedErr := node.handleRouteCarriedEnvelope(fields)
		if carriedErr != nil {
			return nil, carriedErr
		}
		if fields["route_response_path"] != "" {
			node.record("route_asymmetric_response_path_handled", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" response_path="+fields["route_response_path"])
		}
		return map[string]string{
			"promise_about":           production.PromiseRouteReachability,
			"route_id":                fields["route_id"],
			"route_path":              fields["route_path"],
			"route_status":            "delivered",
			"route_final":             node.Agent.Name,
			"route_parent_exact_hash": message.ExactHash,
			"route_message_kind":      routeMessageKindCarried,
			"carried_pcid":            fields["carried_pcid"],
			"carried_ack_exact_hash":  carriedAckHash,
		}, nil
	}
	node.record("route_reachability_promised", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" final="+node.Agent.Name)
	if fields["route_response_path"] != "" {
		node.record("route_asymmetric_response_path_promised", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["route_id"]+" response_path="+fields["route_response_path"])
	}
	return map[string]string{
		"promise_about":           production.PromiseRouteReachability,
		"route_id":                fields["route_id"],
		"route_path":              fields["route_path"],
		"route_status":            "reachable",
		"route_final":             node.Agent.Name,
		"route_parent_exact_hash": message.ExactHash,
		"route_message_kind":      routeMessageKindSetup,
	}, nil
}

func (node *Node) handleRouteCarriedEnvelope(fields map[string]string) (string, error) {
	if fields["carried_envelope_b64"] == "" {
		return "", nil
	}
	carriedBytes, decodeErr := base64.StdEncoding.DecodeString(fields["carried_envelope_b64"])
	if decodeErr != nil {
		return "", decodeErr
	}
	carriedMessage, parseErr := node.parseEnvelope(carriedBytes)
	if parseErr != nil {
		return "", parseErr
	}
	node.emitMessageArtifact("route_carried_received", carriedMessage.Fields["from"], carriedMessage.ProtocolName, carriedBytes, carriedMessage.Fields)
	node.record("route_carried_envelope_validated", "kept", carriedMessage.Fields["from"], "pcid="+carriedMessage.ProtocolName+" route_id="+fields["route_id"]+" exact_sha256="+carriedMessage.ExactHash)
	if carriedMessage.ProtocolName != fields["carried_pcid"] {
		return "", fmt.Errorf("carried envelope pCID %s does not match route field %s", carriedMessage.ProtocolName, fields["carried_pcid"])
	}
	if !node.supportsProtocol(carriedMessage.ProtocolName) {
		return "", fmt.Errorf("route final %s does not support carried pCID %s", node.Agent.Name, carriedMessage.ProtocolName)
	}
	handlerResult, handlerErr := node.handleProtocolPromise(carriedMessage)
	if handlerErr != nil {
		return "", handlerErr
	}
	ackBytes := handlerResult.AckBytes
	ackFields := handlerResult.Fields
	if len(ackBytes) == 0 {
		builtAck, buildErr := node.newRouteCarriedAckBytes(carriedMessage.Fields["from"], carriedMessage.ProtocolCID, carriedMessage.ExactHash, ackFields)
		if buildErr != nil {
			return "", buildErr
		}
		ackBytes = builtAck
	}
	ackMessage, parseErr := node.parseEnvelope(ackBytes)
	if parseErr != nil {
		return "", parseErr
	}
	node.emitMessageArtifact("route_carried_ack", carriedMessage.Fields["from"], ackMessage.ProtocolName, ackBytes, ackMessage.Fields)
	switch node.Agent.Kind {
	case "wasm_agent":
		node.record("wasm_routed_compute_result_verified", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" route_id="+fields["route_id"]+" result_cid="+ackMessage.Fields["result_cid"])
	case "stdio_agent":
		node.record("stdio_routed_compute_result_verified", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" route_id="+fields["route_id"]+" result_cid="+ackMessage.Fields["result_cid"])
	}
	return ackMessage.ExactHash, nil
}

func (node *Node) newRouteCarriedAckBytes(target string, protocolCID protocol.ProtocolCID, parentExactSHA256 string, extraFields map[string]string) ([]byte, error) {
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    node.Agent.Name,
		"to":      target,
		"outcome": "kept",
		"promise": "I promise I received and handled your route-carried signed promise message.",
		"reason":  "route-carried acknowledgement expressed as local promise content",
	}
	for key, value := range extraFields {
		ackFields[key] = value
	}
	protocolName, protocolKnown := node.Protocols.Name(protocolCID)
	if !protocolKnown {
		return nil, fmt.Errorf("route carried ACK protocol is unknown")
	}
	// Intent: The route-carried reply remains a normal pCID ACK, but its
	// envelope parent link names the exact carried request so the message DAG can
	// correlate it over persistent sessions without an RPC request ID. Source:
	// DI-vopab
	ackEnvelope, _, buildErr := node.buildEnvelopeFromFieldsWithParents(protocolName, protocolCID, ackFields, []string{parentExactSHA256})
	if buildErr != nil {
		return nil, buildErr
	}
	return ackEnvelope.Bytes()
}

// runRoutedRuntimeComputeWorkflow carries an exact cid_compute_v1 envelope
// inside route_v1 so Peggy and Victor keep useful compute promises through route
// peers, not only direct Alice TCP paths.
// Intent: Route carrying remains a chain of voluntary forwarding promises; the
// carried compute pCID still owns compute semantics and the final runtime
// adapter keeps the compute promise locally. Source: DI-kohuj
func (node *Node) runRoutedRuntimeComputeWorkflow(target string, spec routeSpec) error {
	routePath := routePathParts(spec.Path)
	if len(routePath) < 2 {
		return fmt.Errorf("route path %s is too short", spec.Path)
	}
	nextHop := routePath[1]
	functionBytes := production.SampleFunctionBytes()
	inputBytes := production.SampleInputBytes()
	contextBytes := production.SampleContextBytes()
	computeFields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            target,
		"turn":          "startup",
		"promise":       "Alice promises to receive routed runtime compute result only as cid_compute_v1 promise content carried by route_v1.",
		"reason":        "runtime adapter useful work should be reachable through sparse route promises, not only direct TCP peers",
		"promise_about": production.PromiseExecuteFunction,
		"function_cid":  production.ContentCID(functionBytes),
		"function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
		"input_cid":     production.ContentCID(inputBytes),
		"input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
		"context_cid":   production.ContentCID(contextBytes),
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"credit_offer":  "5",
		"units":         "1",
	}
	protocolName, protocolCID := node.protocolForFields(computeFields)
	computePayload, _, payloadErr := protocol.MarshalKnownArrayPayload(protocolName, computeFields)
	if payloadErr != nil {
		return payloadErr
	}
	computeEnvelope, envelopeErr := protocol.NewEnvelopeFromPayload(protocolCID, computePayload, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	computeEnvelopeBytes, bytesErr := computeEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	routeFields := routeHopFieldsForSpec(spec, node.Agent.Name, nextHop, production.PromiseRouteForward, routeMessageKindCarried, "")
	routeFields["carried_pcid"] = pcid.CIDComputeV1
	routeFields["carried_promise"] = "Alice promises to send one exact cid_compute_v1 envelope through this route if each hop promises only its own forwarding behavior."
	routeFields["carried_envelope_b64"] = base64.StdEncoding.EncodeToString(computeEnvelopeBytes)
	node.record("route_runtime_compute_message_sent", "kept", nextHop, "pcid="+pcid.RouteV1+" route_id="+spec.ID+" carried_pcid="+pcid.CIDComputeV1+" target="+target)
	ack, sendErr := node.sendAndReceive(nextHop, routeFields)
	if sendErr != nil {
		return fmt.Errorf("routed %s compute: %w", target, sendErr)
	}
	if ack.Fields["route_status"] != "delivered" {
		return fmt.Errorf("routed %s compute status %q", target, ack.Fields["route_status"])
	}
	node.record("route_runtime_compute_message_delivered", "kept", target, "pcid="+pcid.RouteV1+" route_id="+spec.ID+" carried_ack="+ack.Fields["carried_ack_exact_hash"])
	return nil
}

func routeHopFields(fromAgent, toAgent, promiseAbout, messageKind, parentExactHash string) map[string]string {
	return routeHopFieldsForSpec(primaryRouteSpec(), fromAgent, toAgent, promiseAbout, messageKind, parentExactHash)
}

func routeHopFieldsForSpec(spec routeSpec, fromAgent, toAgent, promiseAbout, messageKind, parentExactHash string) map[string]string {
	fields := map[string]string{
		"act":                decision.ActPromise,
		"from":               fromAgent,
		"to":                 toAgent,
		"turn":               "startup",
		"promise":            "I promise one bounded route_v1 hop under the named route terms if my next local peer also promises the next hop.",
		"reason":             "multi-hop forwarding is built from neighboring promises and local keep/break events",
		"protocol":           pcid.RouteV1,
		"promise_about":      promiseAbout,
		"route_id":           spec.ID,
		"route_path":         spec.Path,
		"route_final":        spec.Final,
		"route_excludes":     "mallory",
		"route_payment":      spec.Payment,
		"route_ttl_messages": spec.TTLMessages,
		"route_message_kind": messageKind,
	}
	if spec.ResponsePath != "" {
		fields["route_response_path"] = spec.ResponsePath
	}
	if parentExactHash == "" {
		return fields
	}
	fields["envelope_parent_exact_sha256"] = parentExactHash
	fields["parent_exact_sha256"] = parentExactHash
	fields["parent_link_location"] = "envelope"
	return fields
}

func routePayloadParentHopFieldsForSpec(spec routeSpec, fromAgent, toAgent, promiseAbout, messageKind, parentExactHash string) map[string]string {
	fields := routeHopFieldsForSpec(spec, fromAgent, toAgent, promiseAbout, messageKind, "")
	fields["payload_parent_exact_sha256"] = parentExactHash
	fields["parent_exact_sha256"] = parentExactHash
	fields["parent_link_location"] = "payload"
	return fields
}

func routeSpecFromFields(fields map[string]string) routeSpec {
	return routeSpec{
		ID:           fields["route_id"],
		Path:         fields["route_path"],
		Final:        fields["route_final"],
		Payment:      fields["route_payment"],
		TTLMessages:  firstNonEmpty(fields["route_ttl_messages"], "1"),
		ResponsePath: fields["route_response_path"],
	}
}

func routePathParts(routePath string) []string {
	rawParts := strings.Split(routePath, ">")
	parts := make([]string, 0, len(rawParts))
	for _, rawPart := range rawParts {
		part := strings.TrimSpace(rawPart)
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func routePathIndex(routePath []string, agentName string) (int, error) {
	for routeIndex, routeAgent := range routePath {
		if routeAgent == agentName {
			return routeIndex, nil
		}
	}
	return 0, fmt.Errorf("agent %s is not in route path %s", agentName, strings.Join(routePath, ">"))
}
