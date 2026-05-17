# Conditional Release / Geofencing Scenarios

These scenarios are evidence for `TODO-ralud`. They are not a decision and not a
frozen conditional-release protocol. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Onward-restraint chain | Alice sends content to Bob only if Bob promises to forward it only to recipients who make the same promise. Bob wants to forward to Carol. | Whether the recursive promise graph is represented at group-session, conditional-release, transport/feed, or CAS-object level. | If the graph is central to dispatch semantics, group/session ownership gets stronger; if it composes across sessions, a separate family gets stronger. |
| Geofenced group dispatch | Alice permits content only for group members inside a stated region. Carol is a member but is outside the allowed region. | Whether geofence checks are membership checks, per-message dispatch checks, fetch-policy checks, or storage constraints. | The owner layer must explain both refusal and auditability without assuming a central location oracle. |
| Opaque lower-layer carriage | Bob's node stores encrypted content whose condition vocabulary it cannot parse. | Whether lower layers can safely carry opaque condition references while avoiding accidental promise violations. | If opaque carriage is acceptable, the condition object must be verifiable without every layer understanding its semantics. |
| Replay outside conditions | Mallory replays a valid old content reference to Dave outside the allowed audience or geography. | Whether receivers, feeds, or group/session state detect stale or unauthorized reuse. | Replay handling determines whether conditions must bind to recipients, epochs, locations, or session context. |
| Partial compliance evidence | Bob accepts onward-restraint but Carol's node cannot express the same condition vocabulary. | Whether the protocol rejects forwarding, downgrades the condition, records partial compliance, or asks for a translation promise. | Mixed-version behavior exposes whether the design needs explicit condition-version negotiation. |

## Expected Outputs

- Evidence for whether `TODO-ralud` should open a group-session follow-on, a new
  conditional-release protocol family, or a transport/feed-visible constraint
  question.
- A list of condition-evidence fields that a later TE/DR must either require,
  reject, or deliberately defer.
