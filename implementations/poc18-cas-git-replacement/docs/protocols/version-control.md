# POC18 version-control protocol

## Status

Planned POC18 protocol specification. The exact bytes of this Markdown file
derive the active pCID. Source: `DI-lidaj`; `DI-zuruj`; `DI-dibut`;
`DI-dofoj`; `DI-radaj`; `TE-kopap`.

## Abstract

This protocol lets PromiseGrid agents promise version-control state without a
central repository, global branch authority, global tag authority, or forge
authority. It represents files, directories, POSIX node versions, snapshots,
branches, tags, releases, logical changes, review threads, workspaces, Git
bridge mappings, retention terms, and sparse peer synchronization as signed
PromiseGrid messages over CID-addressed CAS objects.

The native model is PromiseGrid-native: versioned reference sets and sparse CAS
are the source of truth. Conventional Git import, export, push, and pull are
required bridge adapter behavior, not the native synchronization model.

## pCID and envelope

The pCID for this protocol is the CIDv1 raw sha2-256 CID of this exact
Markdown file's bytes. The spec file does not contain its own final pCID because
adding the pCID to the file would change the bytes and therefore change the
pCID.

All protocol messages use this envelope shape:

```text
grid([42(pCID), parents, payload, proof])
```

| Slot | Name | Meaning |
|---|---|---|
| `0` | `42(pCID)` | DAG-CBOR tag 42 containing the protocol CID of this spec. |
| `1` | `parents` | Array of typed parent links to prior exact message CIDs or related object CIDs. |
| `2` | `payload` | Protocol-owned version-control promise payload. |
| `3` | `proof` | Detached proof over the pCID-defined signable view. |

The signable view is:

```text
grid([42(pCID), parents, payload])
```

The proof slot MUST cover the exact bytes of the signable view. The proof
profile for POC18 is a detached single-signer proof containing:

```text
proof = [
  proof_profile,
  signer,
  algorithm,
  public_key_cid,
  signature_bytes,
  created_at
]
```

| Proof slot | Meaning |
|---|---|
| `proof_profile` | Text label for the proof profile, for example `ed25519_detached`. |
| `signer` | Local name, DID, key fingerprint, or other pCID-defined signer reference. |
| `algorithm` | Signature algorithm label. POC18 examples use Ed25519. |
| `public_key_cid` | CID of the public key bytes or key record used to verify the signature. |
| `signature_bytes` | Signature over the exact signable-view bytes. |
| `created_at` | Sender-local timestamp or monotonic event label. Receivers interpret freshness locally. |

Production successors MAY replace this proof profile with COSE or another
pCID-defined proof shape. Such a change creates a different pCID because it
changes the spec bytes.

## Promise Theory model

Every message is a promise by the `promiser`. No agent promises on behalf of
another agent. A reference set, branch, tag, release, review statement, storage
statement, retention statement, Git bridge mapping, or merge statement is only
the promiser's local promise.

Trust is local and relationship-relative. A receiver MAY remember whether prior
promises were kept or broken, but this protocol does not define a global trust
authority, permission service, merge authority, repository authority,
conformance service, role-based access authority, or global monitor.

An agent that does not promise a behavior has not broken a promise by failing to
perform it. A malformed message, unsupported promise kind, missing CAS object,
expired retention term, or locally unsafe materialization target results in
local non-commitment unless the receiver has a prior local outstanding promise
that says otherwise.

## Core terms

| Term | Meaning |
|---|---|
| CAS object | Exact bytes addressed by CID. CAS stores are local and partial. |
| Message CID | CID of the exact CBOR grid message bytes. |
| Object CID | CID of any exact byte object, including chunks, manifests, specs, messages, or mappings. |
| Reference set | Versioned labeled set from labels to one or more target CIDs. |
| POSIX node | Versioned file-system object record: regular file, directory, symlink, hard link, character device, block device, FIFO, or socket. |
| Snapshot | Workspace/root state that names a root directory reference set and parent snapshots. |
| Logical change | Versioned reference set that gives evolving work stable identity across revisions. |
| Review thread | Versioned reference set that groups submissions, comments, tests, acceptance promises, and requested-change promises. |
| Git bridge | Adapter that maps Git refs/objects to and from native PromiseGrid objects. |
| Native sync | Continuous peer DAG synchronization of promises and missing objects, not Git push/pull. |
| Storage-payment token | Exact COSE_Sign1 bytes containing CWT-style CBOR claims. The token is an issuer promise that a bearer may redeem under stated local storage terms. |

## Payload grammar

The payload is a pCID-owned CBOR array:

```text
payload = [
  promiser,
  promisee,
  promise_kind,
  promise_body,
  reciprocal_promise,
  local_constraints
]
```

| Payload slot | Meaning |
|---|---|
| `promiser` | Agent making the promise. |
| `promisee` | Intended promisee, peer, group, local role, or empty string for public advertisement. |
| `promise_kind` | One of the promise kinds defined below. |
| `promise_body` | Kind-specific CBOR array or map defined by this spec. |
| `reciprocal_promise` | Promise requested, offered, or remembered in return; empty array when absent. |
| `local_constraints` | Sender-local constraints such as retention limit, path safety, trust threshold, or bridge loss note. |

`promiser`, `promisee`, and `promise_kind` are protocol fields, not transport
addresses. Routing and delivery remain implementation-local. A pCID selects the
parser and slot grammar; it is not an operation code or destination.

## Parent links

The parent slot is an array of typed links:

```text
parents = [
  [parent_role, 42(parent_cid)],
  ...
]
```

| Parent role | Meaning |
|---|---|
| `previous_message` | Prior message in the same local promise chain. |
| `previous_node` | Prior POSIX node version. |
| `previous_reference_set` | Prior version of the same reference set. |
| `previous_snapshot` | Prior snapshot/change-set version. |
| `supersedes` | Earlier message intentionally replaced by this message. |
| `responds_to` | Message this message answers or narrows. |
| `review_of` | Message or object being reviewed. |
| `git_bridge_source` | Git mapping object or imported Git object that informed this promise. |
| `redeems_token` | Exact storage-payment token bytes consumed by a redemption promise. |
| `paid_by` | Storage-payment redemption promise used as reciprocal economics for a retention promise. |

Receivers MUST treat missing parents as sparse-DAG state, not proof of bad
faith. A receiver MAY request missing parents through `sync_interest`.

## Promise kinds

### `chunk_manifest`

Promises the chunk list for file content. POC18 uses Rabin content-defined
chunking for all regular-file content, including small text files and large
binary files.

```text
promise_body = [
  manifest_cid,
  file_size,
  chunker,
  chunker_parameters,
  chunks,
  content_digest
]

chunks = [
  [offset, length, 42(chunk_cid)],
  ...
]
```

`chunker` is `rabin`. `chunker_parameters` is a pCID-owned map or array
describing average, minimum, and maximum chunk sizes once implementation locks
those values. `content_digest` is optional auxiliary information; CIDs remain
authoritative.

### `posix_node`

Promises one POSIX node version.

```text
promise_body = [
  node_identity,
  node_type,
  content,
  metadata,
  materialization
]
```

| Slot | Meaning |
|---|---|
| `node_identity` | Stable logical node identifier chosen by the promiser, usually a CID of a node-origin record. |
| `node_type` | `regular`, `directory`, `symlink`, `hard_link`, `char_device`, `block_device`, `fifo`, or `socket`. |
| `content` | Type-specific content reference. |
| `metadata` | Mode bits, ownership notes, timestamps, xattrs, or local bridge metadata as permitted by the implementation. |
| `materialization` | Local safety notes and host constraints. |

Content by node type:

| Node type | Content meaning |
|---|---|
| `regular` | `42(chunk_manifest_cid)`. |
| `directory` | `42(reference_set_cid)` whose role is `directory`. |
| `symlink` | Raw target bytes or CID of target-byte record. |
| `hard_link` | Link-group identity or target node identity shared by multiple directory labels. |
| `char_device` | Device metadata record, not live device access. |
| `block_device` | Device metadata record, not live device contents. |
| `fifo` | FIFO metadata record, not stream contents. |
| `socket` | Socket metadata record, not live socket state. |

### `reference_set`

Promises a versioned labeled set. This is the shared mechanism for directories,
filenames, tags, branches, releases, logical changes, review threads, and
workspaces.

```text
promise_body = [
  reference_set_identity,
  role,
  namespace_or_context,
  entries,
  promised_terms
]

entries = [
  [label, targets, entry_terms],
  ...
]

targets = [
  [target_role, 42(target_cid)],
  ...
]
```

| Role | Required meaning |
|---|---|
| `directory` | Labels are filename/path-component dirents. Targets are POSIX node versions or directory reference sets. |
| `branch` | Usually labels `head` to a snapshot CID. |
| `tag` | Labels one or more durable targets. Tags are versioned reference sets, not optional Git side data. |
| `release` | Labels source, binaries, SBOMs, docs, signatures, and review/test artifacts. |
| `logical_change` | Gives evolving work stable identity across revisions. |
| `review_thread` | Groups submissions, comments, test results, requested changes, and adoption promises. |
| `workspace` | Names a root directory plus toolchain, dependency, local override, or materialization receipt CIDs. |

Multiple labels MAY point to the same target CID. One label MAY point to
multiple target CIDs when the role permits it. This makes a multi-target tag and
a directory the same underlying mechanism with different role validation.

### `snapshot`

Promises a root directory state with parent snapshot ancestry.

```text
promise_body = [
  snapshot_identity,
  root_directory,
  parent_snapshots,
  change_summary,
  materialization_terms
]
```

`root_directory` is `42(reference_set_cid)` for a directory role. Parent
snapshots SHOULD also be present in the envelope `parents` slot as
`previous_snapshot` links.

### `object_availability`

Promises that the promiser has, can serve, can forward, or knows how to obtain
selected CIDs.

```text
promise_body = [
  availability_scope,
  objects,
  service_terms
]

objects = [
  [object_role, 42(object_cid), byte_count, availability_status],
  ...
]
```

`availability_status` values include `have`, `can_serve`, `can_forward`,
`missing`, and `not_promised`.

### `object_retention`

Promises retention behavior for selected CIDs.

```text
promise_body = [
  retention_scope,
  objects,
  retention_until,
  pressure_terms,
  compensation_terms
]
```

Retention is local. Garbage collection is promise-based: an agent SHOULD publish
what it is willing to retain or collect instead of silently deleting objects it
previously promised to retain.

### `storage_payment_redemption`

Promises local redemption of a storage-payment bearer token.

```text
promise_body = [
  42(storage_payment_token_cid),
  42(paid_object_cid),
  payment_scope,
  payment_value,
  payment_unit,
  redeemed_at,
  transferable
]
```

The `storage_payment_token_cid` names exact COSE_Sign1 bytes stored in CAS. Those
bytes contain CWT-style CBOR claims, including issuer, subject, expiration time,
token id, capability marker, scope, paid object CID, value, unit, and
transferability. A receiver MUST verify the issuer signature before treating the
token as reciprocal economics, and SHOULD reject token replay using its own local
spent-token ledger. This is not a global balance, currency, authorization, or
settlement authority.

### `review_statement`

Promises a local review, test, adoption, or requested-change statement about
exact target CIDs.

```text
promise_body = [
  review_role,
  targets,
  statement,
  result,
  supporting_objects
]
```

`result` values include `accepted_locally`, `rejected_locally`,
`changes_requested`, `test_kept`, `test_broken`, and `not_promised`.

### `git_bridge_mapping`

Promises how the promiser mapped Git state to or from native PromiseGrid state.

```text
promise_body = [
  bridge_direction,
  git_context,
  mappings,
  loss_records,
  bridge_terms
]
```

| Slot | Meaning |
|---|---|
| `bridge_direction` | `import`, `export`, `pull`, or `push`. |
| `git_context` | Git remote, local path, ref namespace, or bridge run identifier as local text. |
| `mappings` | Array of Git object/ref labels to PromiseGrid target CIDs. |
| `loss_records` | Explicit mapping, loss, refusal, or non-commitment records. |
| `bridge_terms` | Local constraints and verification notes. |

Git bridge behavior MUST NOT make a Git remote, appview, forge, hidden ref, or
role-based permission service the PromiseGrid source of truth.

### `sync_interest`

Promises willingness to receive, reciprocate for, or locally ignore selected
missing objects.

```text
promise_body = [
  interest_scope,
  wanted_objects,
  offer_terms,
  refusal_terms
]
```

`wanted_objects` is an array of typed CIDs. `offer_terms` may include reciprocal
storage, credits, bearer capability tokens, or no reciprocal offer. A peer is
free not to promise service.

## Native object graph

The native POC18 object graph is:

```text
+----------------------+
| reference_set        |
| role: branch         |
| entry: head          |
+----------+-----------+
           |
           v
+----------------------+
| snapshot             |
| root_directory       |
| parent_snapshots     |
+----------+-----------+
           |
           v
+----------------------+
| reference_set        |
| role: directory      |
| label -> node CID    |
+----------+-----------+
           |
           v
+----------------------+
| posix_node           |
| type: regular        |
| content: manifest    |
+----------+-----------+
           |
           v
+----------------------+
| chunk_manifest       |
| rabin chunks         |
+----------+-----------+
           |
           v
+----------------------+
| CAS chunks           |
| exact bytes          |
+----------------------+
```

## Continuous peer DAG sync

Native synchronization is not Git push/pull. It is continuous exchange of
promises between peers that locally choose whether to trust, retain, serve, or
ignore objects.

```text
+-------+       reference_set       +-------+
| Alice | ------------------------> | Bob   |
|       |                           |       |
+---+---+                           +---+---+
    |                                   |
    | object_availability               | local trust check
    v                                   v
+-------+       sync_interest       +-------+
| Carol | <------------------------ | Bob   |
|       |                           |       |
+---+---+                           +---+---+
    |                                   |
    | object_retention                  | object_availability
    v                                   v
+-------+       missing objects     +-------+
| Dave  | ------------------------> | Bob   |
+-------+                           +-------+
```

Bob does not need a complete repository view. Bob may trust Alice's branch
promise, Carol's storage promise, and Dave's retention promise differently.

## Git bridge flow

Git compatibility is adapter behavior around the native model:

```text
+------------------+       read        +------------------+
| Git refs/objects | ----------------> | conversion core  |
+------------------+                   +--------+---------+
                                                |
                                                v
                                      +------------------+
                                      | PromiseGrid CAS  |
                                      | reference sets   |
                                      +--------+---------+
                                                |
                                                v
+------------------+       write       +------------------+
| Git refs/objects | <---------------- | conversion core  |
+------------------+                   +------------------+
```

Import and pull read Git refs/objects into native PromiseGrid objects. Export
and push materialize compatible native state back to Git refs/objects. The same
conversion core MUST be used in both directions.

## Review and logical-change flow

```text
+--------------------+
| logical_change     |
| identity: X        |
+---------+----------+
          |
          v
+--------------------+       review_of       +--------------------+
| submission round 1 | --------------------> | review_thread      |
+---------+----------+                       | comments/tests     |
          |                                  +---------+----------+
          | supersedes                                 |
          v                                            | changes requested
+--------------------+                                 v
| submission round 2 | --------------------> +--------------------+
+---------+----------+       review_of       | review statement   |
          |                                  | accepted locally   |
          v                                  +--------------------+
+--------------------+
| submission round 3 |
+--------------------+
```

Logical identity is a reference-set promise, not a raw change-ID field. Review
rounds are exact CIDs; comments and tests attach to exact targets.

## Materialization flow

```text
+--------------------+       promises root       +--------------------+
| workspace ref set  | ------------------------> | materializer role  |
+---------+----------+                           +---------+----------+
          |                                                |
          | root directory CID                             | local safety check
          v                                                v
+--------------------+                           +--------------------+
| directory ref set  |                           | local filesystem   |
| labels -> nodes    | ------------------------> | files and metadata |
+--------------------+       create locally      +--------------------+
```

The materializer is a local role. It may refuse to create device nodes, sockets,
hard links, symlinks, or paths that violate local policy. Such refusal is local
non-commitment unless the materializer had already promised the unsafe work.

## Parser and receiver behavior

A generic receiver MUST parse only:

1. The outer CBOR grid tag.
2. The array arity.
3. Slot 0 tag 42 pCID bytes.
4. Raw slots 1 through 3 for dispatch to this pCID's parser.

The pCID-specific parser MUST validate:

- Exact arity: four slots.
- Parent slot is an array of two-slot `[role, 42(cid)]` entries.
- Payload is a six-slot array.
- `promise_kind` is one of this spec's defined kinds.
- The `promise_body` shape matches `promise_kind`.
- CIDs are valid CIDv1 base32 when printable and valid binary CID bytes on wire.
- No field treats pCID as a destination, method, repository name, branch name, or
  payload object hash.
- Proof verifies over `grid([42(pCID), parents, payload])`.

Unknown promise kinds, malformed bodies, invalid CIDs, unsupported POSIX node
types, unsafe materialization targets, missing parents, or missing CAS bytes
MUST produce local rejection or non-commitment, not global failure.

## CAS and DAG storage

Agents SHOULD store exact message bytes and exact object bytes by CID. Stores
are sparse and local. No CAS is assumed complete.

Messages form a DAG through the envelope `parents` slot. Payload fields may also
link to object CIDs. Parent links explain version ancestry or response lineage;
payload links name protocol objects.

Agents MAY keep local indexes for:

- Message CID to parent CIDs.
- Reference-set identity to latest locally trusted versions.
- Node identity to locally known node versions.
- Snapshot CID to root directory reference set.
- Git object/ref labels to PromiseGrid CIDs.
- Retention promises and expiration pressure.

Indexes are local derived state. Exact CAS bytes remain the source material.

## Security considerations

All received bytes are untrusted until CID-verified and proof-verified where a
proof is required. A message with a valid proof is still only a promise by the
signer, not evidence of global truth.

Receivers SHOULD defend against:

- CID text that is not canonical CIDv1 base32.
- Binary CID bytes that do not match CIDv1 raw sha2-256 where that profile is
  required.
- Large or cyclic DAG traversals.
- Malformed CBOR, trailing CBOR, wrong arity, wrong tag, or wrong body shape.
- Branch/tag/release promises by agents the receiver does not locally trust.
- Git bridge mappings that silently drop POSIX node information.
- Review statements that imply global acceptance authority.
- Retention promises whose terms have expired or whose reciprocal promise was
  not kept.

## Git interoperability

Git bridge implementations MUST preserve content and DAG semantics where Git can
represent them. For features Git cannot represent, bridge mappings MUST record
explicit local outcomes:

| Native feature | Git bridge requirement |
|---|---|
| Regular file | Export/import as Git blob where possible. |
| Directory | Export/import as Git tree where possible. |
| Symlink | Preserve Git symlink mode and target bytes. |
| Hard link | Record mapping or loss; Git trees do not preserve hard-link identity directly. |
| Char/block device | Record non-commitment, loss, or local materialization note. |
| FIFO/socket | Record non-commitment, loss, or local materialization note. |
| Branch | Map to/from branch-role reference set. |
| Tag | Map to/from tag-role or release-role reference set. |
| Logical change | Map to review/branch artifacts if needed; native identity is a reference set. |
| Review thread | May export as notes, refs, files, or external bridge metadata; native identity remains a reference set. |

Git push/pull are names for bridge interactions with conventional Git remotes.
Native PromiseGrid peers exchange promises continuously and do not need explicit
push/pull.

## Annotated examples

Examples use placeholder CIDs. Printable examples MUST use real CIDv1 base32
strings in executable code and diagnostics.

### Regular file node

```text
grid([42(pCID),
  [
    ["previous_node", 42(bafkreipreviousnodecid...)]
  ],
  [
    "alice",
    "bob",
    "posix_node",
    [
      "node:readme",
      "regular",
      42(bafkreimanifestcid...),
      {"mode": "100644"},
      {"safe_to_materialize": true}
    ],
    [],
    {"trust_required": "local"}
  ],
  proof
])
```

| Item | Meaning |
|---|---|
| `previous_node` | This node version follows an earlier node version. |
| `alice` | Alice promises the node version. |
| `bob` | Bob is the intended promisee. |
| `posix_node` | The payload body follows the POSIX node grammar. |
| `regular` | The node is a regular file. |
| `42(bafkreimanifestcid...)` | File content is named by a Rabin chunk manifest CID. |
| `mode` | POSIX mode metadata retained for materialization. |
| `trust_required` | Bob still judges locally whether to trust or materialize it. |

### Directory reference set

```text
grid([42(pCID),
  [
    ["previous_reference_set", 42(bafkreipreviousdircid...)]
  ],
  [
    "alice",
    "bob",
    "reference_set",
    [
      "refset:src-directory",
      "directory",
      "workspace:demo",
      [
        ["README.md", [["node", 42(bafkreireadmenodecid...)]], []],
        ["cmd", [["directory", 42(bafkreicmddircid...)]], []]
      ],
      ["I promise these labels are my current directory entries."]
    ],
    [],
    []
  ],
  proof
])
```

| Item | Meaning |
|---|---|
| `directory` | This reference set is validated as a directory. |
| `README.md` | Filename label owned by this directory reference-set version. |
| `node` | Target role for a POSIX node version. |
| `cmd` | Directory label pointing at another directory reference set. |

### Branch head reference set

```text
grid([42(pCID),
  [
    ["previous_reference_set", 42(bafkreipreviousbranchcid...)]
  ],
  [
    "alice",
    "",
    "reference_set",
    [
      "refset:main",
      "branch",
      "project:wire-lab-example",
      [
        ["head", [["snapshot", 42(bafkreisnapshotcid...)]], []]
      ],
      ["I promise this is my current main branch head."]
    ],
    [],
    {"audience": "public advertisement"}
  ],
  proof
])
```

The empty promisee means public advertisement, not global truth. Bob may accept,
ignore, or compare Alice's branch promise with other peers' promises.

### Logical change revision

```text
grid([42(pCID),
  [
    ["previous_reference_set", 42(bafkreilogicalchangeold...)]
  ],
  [
    "carol",
    "dave",
    "reference_set",
    [
      "change:portable-cas-checkout",
      "logical_change",
      "project:wire-lab-example",
      [
        ["round-3", [["snapshot", 42(bafkreisnapshotround3...)]], []],
        ["prior-round", [["snapshot", 42(bafkreisnapshotround2...)]], []]
      ],
      ["I promise this is the next revision of the same logical change."]
    ],
    [],
    []
  ],
  proof
])
```

The logical change identity is the reference set itself. It is not a mutable
forge pull request and not a raw change-ID field.

### Review statement

```text
grid([42(pCID),
  [
    ["review_of", 42(bafkreisnapshotround3...)]
  ],
  [
    "dave",
    "carol",
    "review_statement",
    [
      "maintainer_review",
      [["snapshot", 42(bafkreisnapshotround3...)]],
      "I promise I reviewed this exact snapshot locally.",
      "changes_requested",
      [["comment", 42(bafkreicommentcid...)]]
    ],
    [],
    {"authority": "local only"}
  ],
  proof
])
```

Dave promises only his local review result. He does not command Carol or any
other agent.

### Git bridge mapping

```text
grid([42(pCID),
  [
    ["git_bridge_source", 42(bafkreigitcommitbytes...)]
  ],
  [
    "alice",
    "bob",
    "git_bridge_mapping",
    [
      "import",
      {"remote": "origin", "ref": "refs/heads/main"},
      [
        ["git_commit", "abc123...", "snapshot", 42(bafkreisnapshotcid...)],
        ["git_tree", "def456...", "directory", 42(bafkreidirectorycid...)]
      ],
      [
        ["fifo", "path tmp/build.pipe", "not representable in Git tree"]
      ],
      {"verified": "git object bytes and native CIDs checked"}
    ],
    [],
    []
  ],
  proof
])
```

The bridge mapping records what Alice did during import. It does not make Git
the native authority and does not hide loss.

### Continuous sync interest

```text
grid([42(pCID),
  [
    ["responds_to", 42(bafkreibranchadvertisement...)]
  ],
  [
    "bob",
    "alice",
    "sync_interest",
    [
      "missing_objects_for_branch",
      [
        ["snapshot", 42(bafkreisnapshotcid...)],
        ["chunk", 42(bafkreichunkcid...)]
      ],
      ["I promise to receive these objects and reciprocate with storage credit."],
      []
    ],
    ["storage_credit", "2"],
    {"max_bytes": "10485760"}
  ],
  proof
])
```

Bob promises willingness to receive missing objects under his local constraints.
Alice remains free not to serve them.

### Retention promise

```text
grid([42(pCID),
  [
    ["paid_by", 42(bafkreiredemptioncid...)]
  ],
  [
    "frank",
    "alice",
    "object_retention",
    [
      "release-retention",
      [
        ["release", 42(bafkreireleasecid...)],
        ["snapshot", 42(bafkreisnapshotcid...)]
      ],
      "2026-07-31T00:00:00Z",
      ["I may collect unlisted chunks under local disk pressure."],
      [
        ["storage_payment_token", 42(bafkreitokencid...), "redeemed"],
        ["storage_payment_redemption", 42(bafkreiredemptioncid...), "accepted_locally"]
      ]
    ],
    [
      ["storage_payment_token", 42(bafkreitokencid...), "redeemed"],
      ["storage_payment_redemption", 42(bafkreiredemptioncid...), "accepted_locally"]
    ],
    ["store=frank-local-cas"]
  ],
  proof
])
```

Frank promises retention for selected objects and states collection pressure for
unlisted objects. The `paid_by` parent points at Frank's
`storage_payment_redemption` promise, which in turn has a `redeems_token` parent
to Alice's exact signed storage-payment token bytes. Alice judges Frank's future
keep/break history locally.

## Implementation requirements

A conforming POC18 implementation of this protocol MUST:

- Parse the outer grid tag, slot count, tag 42 pCID, parent array, payload array,
  and proof slot exactly as specified.
- Reject malformed CBOR, wrong arity, invalid CID links, unknown promise kinds,
  and promise bodies that do not match their kind.
- Store exact message bytes under their message CID when local retention policy
  promises retention.
- Store exact CAS object bytes under object CIDs using CIDv1 base32 names when
  printable.
- Verify proof bytes before treating the message as a signed promise.
- Keep pCID routing separate from app, branch, path, repository, operation, and
  peer addressing.
- Preserve all POSIX node types in native objects, even when Git bridge output
  must record loss or non-commitment.
- Use Rabin content-defined chunks for all regular-file content.
- Implement Git import/export/push/pull as bridge behavior through one shared
  conversion core.
- Implement native peer synchronization as continuous exchange of
  `reference_set`, `object_availability`, `sync_interest`, and
  `object_retention` promises.

## Non-goals

This protocol does not define:

- A global repository.
- A global branch, tag, release, merge, or review authority.
- A global permission or authorization service.
- A complete global CAS.
- A global monitor.
- A Git replacement command-line user interface.
- A production storage layout.
- A production proof or key-management standard beyond the POC18 detached proof
  profile.

## Source notes

This spec consolidates:

- `DI-zuruj`: versioned reference sets and POC16 baseline.
- `DI-dibut`: Tangled prior-art review and continuous peer DAG sync.
- `DI-dofoj`: Rabin chunking and Git import/export/push/pull bridge.
- `DI-radaj`: all POSIX inode types.
- `DI-lidaj`: one non-versioned implementation-local spec with svgbob-safe
  diagrams and annotated examples.
- `DN-rifir`: versioned reference-set design.
- `DN-dopod`: Tangled comparison.
- `TE-kopap`: hybrid native core plus Git bridge architecture.
