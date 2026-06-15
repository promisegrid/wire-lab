package runtime

import (
	"fmt"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/production"
)

const poc15RouteID = "route-alice-bob-carol-dave-0001"
const poc15RoutePath = "alice>bob>carol>dave"
const routeMessageKindSetup = "setup"
const routeMessageKindCarried = "carried"

// runRoutePromiseWorkflow exercises the first executable POC15 multi-hop route.
// Intent: Alice uses route_v1 to ask only Alice's direct neighbor Bob for a
// voluntary route promise; Bob and Carol each make their own next-hop promises,
// and Dave locally confirms reachability before Alice sends a carried message.
// This is app-level promise forwarding, not a kernel route authority. Source:
// DI-lihir
func (node *Node) runRoutePromiseWorkflow() error {
	node.record("route_exclusion_promise_made", "kept", "bob", "pcid="+pcid.RouteV1+" Alice asks route peers for a route that avoids Mallory as transit")
	node.record("route_exclusion_used_in_choice", "kept", "bob", "pcid="+pcid.RouteV1+" selected_path="+poc15RoutePath+" excluded_transit=mallory")
	node.record("route_payment_promised", "kept", "bob", "pcid="+pcid.RouteV1+" Alice promises reciprocal forwarding credit for one bounded route setup")
	setupFields := routeHopFields(node.Agent.Name, "bob", production.PromiseRouteSetup, routeMessageKindSetup, "")
	node.record("route_setup_promise_made", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+poc15RouteID+" path="+poc15RoutePath)
	setupAck, setupErr := node.sendAndReceive("bob", setupFields)
	if setupErr != nil {
		return fmt.Errorf("route setup: %w", setupErr)
	}
	if setupAck.Fields["field_route_status"] != "reachable" {
		return fmt.Errorf("route setup status %q", setupAck.Fields["field_route_status"])
	}
	node.record("route_reachability_confirmed", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+poc15RouteID+" status=reachable")
	carriedFields := routeHopFields(node.Agent.Name, "bob", production.PromiseRouteForward, routeMessageKindCarried, setupAck.ExactHash)
	carriedFields["field_carried_pcid"] = pcid.RelationshipV1
	carriedFields["field_carried_promise"] = "Alice promises to send one bounded relationship_v1 payload only after route reachability was confirmed."
	node.record("route_carried_message_sent", "kept", "bob", "pcid="+pcid.RouteV1+" route_id="+poc15RouteID+" carried_pcid="+pcid.RelationshipV1)
	carriedAck, carriedErr := node.sendAndReceive("bob", carriedFields)
	if carriedErr != nil {
		return fmt.Errorf("route carried message: %w", carriedErr)
	}
	if carriedAck.Fields["field_route_status"] != "delivered" {
		return fmt.Errorf("route carried status %q", carriedAck.Fields["field_route_status"])
	}
	node.record("route_carried_message_delivered", "kept", "dave", "pcid="+pcid.RouteV1+" route_id="+poc15RouteID+" carried_pcid="+pcid.RelationshipV1)
	return nil
}

// handleRoutePromise handles one route_v1 hop from the local app's own vantage.
// Intent: Every hop signs a fresh pCID-owned route_v1 promise and references the
// previous exact envelope hash as parent context; no hop claims authority over
// downstream peers or final delivery beyond its own observed next-hop outcome.
// Source: DI-lihir
func (node *Node) handleRoutePromise(message parsedMessage) (map[string]string, error) {
	fields := message.Fields
	routePath := routePathParts(fields["field_route_path"])
	routeIndex, routeIndexErr := routePathIndex(routePath, node.Agent.Name)
	if routeIndexErr != nil {
		return nil, routeIndexErr
	}
	if routeIndex == len(routePath)-1 {
		return node.handleFinalRouteHop(message)
	}
	nextHop := routePath[routeIndex+1]
	node.record("route_forward_promise_made", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" next_hop="+nextHop)
	if fields["field_route_payment"] != "" {
		node.record("route_payment_promised", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" payment="+fields["field_route_payment"])
	}
	forwardFields := routeHopFields(node.Agent.Name, nextHop, production.PromiseRouteForward, fields["field_route_message_kind"], message.ExactHash)
	for _, key := range []string{"field_carried_pcid", "field_carried_promise", "field_route_setup_parent"} {
		forwardFields[key] = fields[key]
	}
	downstreamAck, downstreamErr := node.sendAndReceive(nextHop, forwardFields)
	if downstreamErr != nil {
		node.record("route_forward_not_promised", "non_commitment", nextHop, "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" error="+downstreamErr.Error())
		return map[string]string{
			"field_promise_about": production.PromiseRouteForward,
			"field_route_id":      fields["field_route_id"],
			"field_route_status":  "not_promised",
			"field_route_path":    fields["field_route_path"],
		}, nil
	}
	node.record("route_forward_promise_kept", "kept", nextHop, "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" downstream_status="+downstreamAck.Fields["field_route_status"])
	return map[string]string{
		"field_promise_about":           production.PromiseRouteReachability,
		"field_route_id":                fields["field_route_id"],
		"field_route_path":              fields["field_route_path"],
		"field_route_status":            downstreamAck.Fields["field_route_status"],
		"field_route_final":             fields["field_route_final"],
		"field_route_parent_exact_hash": message.ExactHash,
		"field_route_final_ack_hash":    downstreamAck.ExactHash,
		"field_route_message_kind":      fields["field_route_message_kind"],
	}, nil
}

func (node *Node) handleFinalRouteHop(message parsedMessage) (map[string]string, error) {
	fields := message.Fields
	if fields["field_route_message_kind"] == routeMessageKindCarried {
		node.record("route_carried_message_received", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" carried_pcid="+fields["field_carried_pcid"])
		return map[string]string{
			"field_promise_about":           production.PromiseRouteReachability,
			"field_route_id":                fields["field_route_id"],
			"field_route_path":              fields["field_route_path"],
			"field_route_status":            "delivered",
			"field_route_final":             node.Agent.Name,
			"field_route_parent_exact_hash": message.ExactHash,
			"field_route_message_kind":      routeMessageKindCarried,
			"field_carried_pcid":            fields["field_carried_pcid"],
		}, nil
	}
	node.record("route_reachability_promised", "kept", fields["from"], "pcid="+pcid.RouteV1+" route_id="+fields["field_route_id"]+" final="+node.Agent.Name)
	return map[string]string{
		"field_promise_about":           production.PromiseRouteReachability,
		"field_route_id":                fields["field_route_id"],
		"field_route_path":              fields["field_route_path"],
		"field_route_status":            "reachable",
		"field_route_final":             node.Agent.Name,
		"field_route_parent_exact_hash": message.ExactHash,
		"field_route_message_kind":      routeMessageKindSetup,
	}, nil
}

func routeHopFields(fromAgent, toAgent, promiseAbout, messageKind, parentExactHash string) map[string]string {
	fields := map[string]string{
		"act":                       decision.ActPromise,
		"from":                      fromAgent,
		"to":                        toAgent,
		"turn":                      "startup",
		"promise":                   "I promise one bounded route_v1 hop under the named route terms if my next local peer also promises the next hop.",
		"reason":                    "multi-hop forwarding is built from neighboring promises and local keep/break events",
		"field_protocol":            pcid.RouteV1,
		"field_promise_about":       promiseAbout,
		"field_route_id":            poc15RouteID,
		"field_route_path":          poc15RoutePath,
		"field_route_final":         "dave",
		"field_route_excludes":      "mallory",
		"field_route_payment":       "reciprocal_forwarding_credit_1",
		"field_route_ttl_messages":  "2",
		"field_route_message_kind":  messageKind,
		"field_parent_exact_sha256": parentExactHash,
	}
	if parentExactHash == "" {
		delete(fields, "field_parent_exact_sha256")
	}
	return fields
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
