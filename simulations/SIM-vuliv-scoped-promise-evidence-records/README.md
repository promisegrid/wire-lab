# SIM-vuliv-scoped-promise-evidence-records

This simulation is the Promise-Theory-first successor to
`SIM-savak-scoped-claim-card-audit-ledger`. It preserves the useful idea of
scoped, durable evidence while removing claim-card and authority-object
framing. Each durable artifact is either a promise record, a promise reference,
or an observed keep/break record about a promise. Source: `DI-tavaz`.

## Design Under Test

Early apps publish small, scoped promise-evidence records instead of claim
cards:

- a promiser names the promise scope, protocol boundary, and promised behavior;
- a peer records later observations of keep, break, refusal, retry, or timeout;
- durable audit publication cites exact evidence objects rather than using a
  general authority card;
- any capability-like behavior is interpreted as evidence of a promise
  relationship, not as bearer-right authority.

## Why this differs from `savak`

`savak` improved auditability but still centered a scoped claim card as the main
artifact. `vuliv` moves the center of gravity back to agents, promises, and
observed outcomes. Evidence objects support later trust decisions; they do not
grant rights.

## Boundaries

This simulation does not define a final universal promise-record envelope. It
tests whether scoped promise-evidence records are a cleaner pattern than
claim-card lineage for audit-heavy PromiseGrid applications.
