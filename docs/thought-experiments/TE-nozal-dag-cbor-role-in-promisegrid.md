# TE-nozal: DAG-CBOR role in PromiseGrid

TE ID: TE-nozal

## Status

needs DF

## Decision Under Test

How much of PromiseGrid should be true DAG-CBOR, and how much should remain
PromiseGrid-specific CBOR that merely uses IPLD-compatible CID links?

This TE tests whether future PromiseGrid messages, payloads, CAS objects,
protocol specs, manifests, chunks, and `grid(...)` envelopes should converge on
true DAG-CBOR, remain custom CBOR, or deliberately use both with precise
boundaries.

## Existing Locked Inputs

- The current formal envelope family is
  `grid([42(pCID), ...protocol-defined-slots])`. Source: `DI-rojij`;
  `DI-punam`.
- Slot `0` is currently `42(pCID)`, where tag `42` is the DAG-CBOR/IPLD CID link
  representation. Source: `DI-kafat`; `DI-sisak`.
- A pCID is the CID of the protocol spec document. It is not a destination
  address, app name, operation, file path, repository name, branch name, or
  payload CID. Source: `DI-009-20260429-173359`; `DI-harih`.
- POC16 and POC18 code currently validate a CIDv1 raw/sha2-256 profile for
  retained exact bytes and pCIDs. Source: `DI-timah`; `DI-jifuj`.
- POC18 uses raw Rabin chunks plus CBOR/grid manifests in the first slice.
  Source: `DI-dofoj`; `DI-jifuj`; `TE-givul`.
- POC6 already proved a true DAG-CBOR path using `go-ipld-prime`, link nodes,
  tag `42`, and CIDv1 `dag-cbor` CIDs. Source: `DI-sagos`.

## Terminology

### CBOR

CBOR is a general binary data format. It has arrays, maps, byte strings, text
strings, tags, integers, floats, indefinite-length forms, and many possible
encoding choices.

### Canonical CBOR

Canonical or deterministic CBOR narrows encoding choices so the same data has
stable bytes. A system can use canonical CBOR without being DAG-CBOR.

### DAG-CBOR

DAG-CBOR is the IPLD data model encoded as a stricter subset of CBOR. The
official IPLD documentation describes DAG-CBOR as a CBOR subset with extra
constraints for hash-consistent representations and an IPLD link type using
CBOR tag `42`. Its spec states the major practical differences from regular CBOR:

- tag `42` is interpreted as CIDs;
- no other tags are supported;
- maps must only be keyed by strings;
- strict deterministic encoding rules apply;
- indefinite-length items are not supported;
- a DAG-CBOR block is one top-level CBOR object;
- special floats such as `NaN` and infinities are not valid DAG-CBOR.

References:

- `https://ipld.io/docs/codecs/known/dag-cbor/`
- `https://ipld.io/specs/codecs/dag-cbor/spec/`

### PromiseGrid-CBOR

PromiseGrid-CBOR is a shorthand used only in this TE for the current
implementation pattern: CBOR messages, usually deterministic, often using
tag-42 CID links, but not necessarily satisfying the DAG-CBOR subset and not
necessarily carrying CIDv1 `dag-cbor` content-type multicodecs.

## Current Mix In This Repo

### True DAG-CBOR evidence: POC6

`implementations/poc6-dag-cbor-interop/` is the cleanest true DAG-CBOR artifact
in the repo:

- It encodes nodes through `go-ipld-prime/codec/dagcbor`.
- It enables IPLD links.
- It derives CIDv1 CIDs with `cid.DagCBOR`.
- It verifies that tag `42` CID links survive the encode/decode path.

This POC proves that real DAG-CBOR is available, understandable, and usable
without an IPFS daemon.

### POC15 and POC16: CBOR with DAG-CBOR-style links

POC15 and POC16 use real CBOR and tag `42` CID links, but their whole messages
are not true DAG-CBOR blocks. Reasons:

- They use the custom PromiseGrid outer tag `grid`, decimal `1735551332`.
- They use pCID-owned raw CBOR slots.
- They use a narrow hand-written CBOR writer/reader rather than an IPLD data
  model encoder.
- They validate CIDs as CIDv1 `raw`/sha2-256 exact-byte identifiers.

The important correction is:

> `42(pCID)` is DAG-CBOR/IPLD link syntax, but using one tag-42 link does not
> make the whole enclosing object a DAG-CBOR block.

### POC18: raw CIDs for all retained bytes

The current POC18 store computes CIDv1 `raw`/sha2-256 CIDs for exact bytes. That
is true for raw chunks and for `.cbor` objects. Therefore:

- A `.bin` chunk is CID-typed as `raw`; this matches its exact-byte identity.
- A `.cbor` object is also CID-typed as `raw`; this means "these exact bytes,"
  not "this is a DAG-CBOR object."
- POC18 uses tag `42` links inside CBOR objects, but those objects are not
  currently stored with CIDv1 `dag-cbor` content-type multicodecs.

This mix was acceptable for a first slice because it made exact-byte CAS easy,
but it is now a source of vocabulary and design ambiguity.

### Protocol spec docs and pCIDs

Current protocol spec documents are Markdown files addressed by raw CIDs. They
are not DAG-CBOR blocks. A pCID currently says "the exact Markdown bytes of this
spec document hash to this CID," not "this spec is a DAG-CBOR IPLD object."

That may remain the right answer: protocol specs are read by humans, developers,
and LLMs. Markdown is a good source format. If we later need a canonical
machine-readable protocol-spec bundle, that bundle could become a DAG-CBOR
object that links to Markdown, schemas, tests, and examples.

## What DAG-CBOR Gains

### IPLD ecosystem compatibility

True DAG-CBOR objects are immediately legible to IPLD libraries, selectors,
CAR files, Graphsync-like tooling, IPFS/IPLD diagnostic tools, and future
ecosystem bridges. This matters because PromiseGrid already wants CID, IPLD,
IPFS, and Bluesky/ATProto-adjacent interoperability.

### Precise object typing through CID content multicodec

If a CID has content multicodec `dag-cbor`, it tells a receiver that the addressed
bytes should be decoded as DAG-CBOR. If a CID has content multicodec `raw`, it
tells a receiver that the addressed bytes are opaque raw bytes. This is a strong
type signal that current POC18 loses by using `raw` for every retained byte.

### Canonical Merkle-DAG semantics

DAG-CBOR's strictness exists so the same IPLD data model value has stable bytes.
That helps with content addressing, parent links, selectors, traversal,
cross-language compatibility, and long-term storage.

### Standard link encoding

Tag `42` links are already the form PromiseGrid wants for pCID, parent, manifest,
and object references. Making some objects true DAG-CBOR would reduce the gap
between "we use IPLD-style links" and "we are valid IPLD data."

### Better transport/archive interop

If POC18 stores manifests, reference sets, snapshots, and review objects as real
DAG-CBOR blocks, those objects can travel in CAR files and be processed by
IPLD-aware tooling without custom PromiseGrid parsers for every layer.

### Cleaner diagnostic claims

Right now `.cbor` means "these bytes decode as CBOR in our code," not "this CID
declares DAG-CBOR." Using CIDv1 `dag-cbor` for true DAG-CBOR objects would make
diagnostics more honest.

## What DAG-CBOR Sacrifices

### Custom CBOR tags, including `grid(...)`

DAG-CBOR permits tag `42` for links and rejects other tags. The PromiseGrid
outer `grid` tag is a non-42 CBOR tag. Therefore a message shaped as:

```text
grid([42(pCID), ...])
```

is not valid DAG-CBOR as a single block.

This is the central conflict.

### Arbitrary CBOR expressiveness

DAG-CBOR excludes CBOR features that might be convenient in protocol-specific
payloads:

- non-string map keys,
- custom tags,
- indefinite-length arrays/maps/strings,
- undefined,
- special floats,
- multiple top-level CBOR items,
- alternate encodings that a plain CBOR decoder might accept.

These exclusions are mostly good for content addressing, but they reduce escape
hatches.

### Easy use of compact array-heavy constrained-device formats

DAG-CBOR allows arrays, so constrained devices can still use compact arrays. The
cost is not array use itself. The cost is that every object must remain inside
the IPLD data model and cannot use non-42 custom tags or odd CBOR conveniences.

### Local implementation freedom

Current POC16/18 code can write exactly the CBOR subset it needs. Moving to
true DAG-CBOR for production objects would push implementation toward IPLD
libraries or very careful in-house conformance. That is probably good for
long-term correctness but increases near-term complexity.

### The pCID-as-spec-doc hash question

If protocol specs remain Markdown, their pCIDs naturally remain raw CIDs. If
specs become DAG-CBOR bundles, the pCID changes because it would name the
encoded bundle bytes, not the Markdown bytes. That may be desirable later, but
it is not a free migration.

## Alternatives

### Alt A: Keep current PromiseGrid-CBOR everywhere

All PromiseGrid messages remain custom CBOR with the `grid` tag. CIDs continue
to use a raw exact-byte profile unless a later object type opts out.

Advantages:

- Minimal implementation churn.
- Keeps the `grid(...)` tag as the visible byte-level family marker.
- Lets pCID specs define arbitrary CBOR slot shapes.
- Works well for small hand-written parsers and constrained runtimes.

Disadvantages:

- `.cbor` files are not necessarily IPLD/DAG-CBOR blocks.
- CIDs do not advertise CBOR object type.
- IPLD selectors/CAR/Graphsync-like tooling cannot understand the custom grid
  layer without adapters.
- The repo will keep needing careful vocabulary: "CBOR with tag-42 links," not
  "DAG-CBOR."

### Alt B: Make every PromiseGrid object true DAG-CBOR

All messages, manifests, reference sets, snapshots, reviews, specs, and metadata
become valid DAG-CBOR blocks with CIDv1 `dag-cbor` content multicodecs. Raw file
chunks may remain `raw`.

Advantages:

- Maximum IPLD interoperability.
- Strong CID-level type signaling.
- Unified object model for CAS traversal.
- Easier future CAR/export/import tooling.

Disadvantages:

- The current `grid(...)` outer tag cannot remain as a CBOR tag.
- Every pCID payload must obey DAG-CBOR restrictions.
- COSE/CWT and other CBOR-family payloads may need to be stored as byte strings
  or linked raw objects instead of inline arbitrary CBOR.
- The promise envelope loses its explicit custom tag unless `grid` becomes data.

### Alt C: Hybrid: DAG-CBOR for data graph objects, PromiseGrid-CBOR for wire envelopes

Keep `grid([42(pCID), ...])` for transport and signing envelopes. Store durable
CAS graph objects such as manifests, reference sets, snapshots, reviews, and
workspace roots as true DAG-CBOR where useful. Raw chunks remain `raw`. Some
payload slots may carry `42(...)` links to true DAG-CBOR objects.

Advantages:

- Preserves the existing PromiseGrid envelope.
- Gains DAG-CBOR for durable graph data where ecosystem tooling matters most.
- Keeps constrained and protocol-specific wire payloads flexible.
- Lets migration happen object family by object family.

Disadvantages:

- Two object classes remain.
- Developers must know whether a CID names a raw object, a DAG-CBOR object, or a
  PromiseGrid-CBOR envelope.
- Requires clear naming, index metadata, diagnostics, and analyzer gates.

### Alt D: Replace `grid` tag with DAG-CBOR data fields

Define a DAG-CBOR object whose top-level value is a map or array with a field or
slot that means "this is a grid message":

```text
[
  "grid",
  42(pCID),
  parents,
  payload,
  proof
]
```

or:

```text
{
  "kind": "grid",
  "pcid": 42(pCID),
  "slots": [...]
}
```

Advantages:

- The whole message can be true DAG-CBOR.
- `grid` remains visible as data rather than a CBOR tag.
- IPLD tooling can traverse the message.

Disadvantages:

- It changes the canonical envelope shape from the already-studied byte-level
  `grid(...)` tag.
- It adds data overhead.
- It weakens the simple "look for the grid tag" parser story.
- It may reopen envelope decisions thought settled by TE-fikoj / TE-lamun.

### Alt E: Treat `grid` as a transport wrapper outside the IPLD block

The wire frame is:

```text
grid(<dag-cbor-message-block-bytes>)
```

or some equivalent transport wrapper, while the content inside is a valid
DAG-CBOR block with CIDv1 `dag-cbor`.

Advantages:

- Keeps a byte-level PromiseGrid carrier marker.
- Lets the durable block be true DAG-CBOR.
- Makes the wrapper more like transport framing than content identity.

Disadvantages:

- The bytes on the wire are not the bytes named by the message CID unless the
  wrapper is excluded from identity.
- Signing, message IDs, and replay protection must say exactly whether they cover
  the wrapper, the inner DAG-CBOR block, or both.
- It risks duplicating the "local storage wrapper versus content identity"
  confusion TE-givul warned about.

## Scenario Analysis

### Scenario 1: Alice sends Bob a compact promise message

Alice sends Bob a small message that Bob only needs to parse by pCID and deliver
to a local handler.

- Alt A is simplest. Bob checks the custom `grid` tag, reads `42(pCID)`, and
  passes raw slots to the pCID parser.
- Alt B requires `grid` to become data, not a tag.
- Alt C keeps the current wire path unchanged.
- Alt D works, but changes the established envelope.
- Alt E adds wrapper-versus-content identity complexity.

Best fit: Alt A or Alt C.

### Scenario 2: Carol stores a version-control snapshot for long-term retrieval

Carol stores a snapshot, directory reference sets, chunk manifests, and review
thread objects. Dave may retrieve them years later with generic IPLD tools.

- Alt A requires PromiseGrid-specific decoders forever.
- Alt B gives maximum generic traversal.
- Alt C gives generic traversal for durable graph objects while preserving wire
  envelope flexibility.
- Alt D also gives generic traversal if all messages are DAG-CBOR.
- Alt E is workable only if the inner block identity is clear.

Best fit: Alt B or Alt C; Alt C is less disruptive.

### Scenario 3: Ellen exports a bundle to an IPFS/IPLD ecosystem peer

Ellen wants to export POC18 objects to a CAR file and have Frank inspect links
using existing tooling.

- Alt A exports opaque raw/CBOR blobs unless custom codecs are registered.
- Alt B works naturally.
- Alt C works for the DAG-CBOR object families, while PromiseGrid-CBOR envelopes
  may remain opaque or be represented by DAG-CBOR summaries.
- Alt D works naturally.
- Alt E works if the durable inner block is what enters the CAR.

Best fit: Alt B, C, or D. Alt C is the pragmatic migration path.

### Scenario 4: Grace runs on a small device

Grace's node has a tiny parser and wants only to check the outer pCID and forward
or reject.

- Alt A keeps the parser small.
- Alt B may require a DAG-CBOR parser and IPLD model awareness even for message
  triage.
- Alt C lets Grace parse the simple PromiseGrid envelope and ignore linked
  durable objects she does not need.
- Alt D can be compact if array-shaped, but requires treating `grid` as a data
  marker.
- Alt E needs wrapper and inner parser rules.

Best fit: Alt A or Alt C.

### Scenario 5: Heidi embeds COSE, CWT, or encrypted payloads

Heidi's protocol uses COSE_Sign1 or CWT bytes, which are themselves CBOR-family
objects with their own tag/profile rules.

- Alt A can carry these as pCID-owned raw CBOR slots or byte strings.
- Alt B may need COSE/CWT to be byte strings or linked raw objects if their tags
  are not valid DAG-CBOR.
- Alt C lets the wire envelope carry opaque COSE/CWT bytes while durable graph
  objects remain DAG-CBOR.
- Alt D faces the same issue as Alt B.
- Alt E can separate transport from durable data, but must define signable views
  carefully.

Best fit: Alt A or Alt C.

### Scenario 6: Ivan maintains a local CAS index

Ivan wants a filesystem CAS where CIDs reveal enough type information for
diagnostics and storage policy.

- Alt A requires local index metadata to distinguish raw chunks, CBOR messages,
  Markdown specs, and other objects, because all may be raw-CID bytes.
- Alt B lets CID multicodec do more work.
- Alt C lets durable DAG-CBOR objects be self-typed by CID, while envelopes and
  raw chunks remain explicit.
- Alt D makes all messages self-typed by CID.
- Alt E depends on whether Ivan indexes wrapper CIDs or inner CIDs.

Best fit: Alt C.

## Where Does The `grid()` Tag Go If Everything Becomes DAG-CBOR?

It cannot remain where it is.

Current PromiseGrid envelopes use a custom CBOR tag:

```text
grid([42(pCID), ...protocol-defined-slots])
```

True DAG-CBOR permits tag `42` for CID links and rejects other tags. Therefore
the custom `grid` CBOR tag has four possible futures:

### Option 1: `grid` remains only on non-DAG-CBOR wire envelopes

Wire messages remain PromiseGrid-CBOR. Durable linked graph objects may be
DAG-CBOR. This is Alt C.

This is the least disruptive path.

### Option 2: `grid` becomes a DAG-CBOR data value

The top-level DAG-CBOR object uses a string or integer marker:

```text
["grid", 42(pCID), parents, payload, proof]
```

This makes the whole message valid DAG-CBOR but changes the byte-level envelope.

### Option 3: `grid` becomes the pCID itself

There is no outer `grid` marker. The message is simply a DAG-CBOR array:

```text
[42(pCID), parents, payload, proof]
```

The receiver recognizes the family by the fact that slot 0 is a tag-42 pCID and
the pCID names a PromiseGrid envelope spec. This is elegant, but it removes the
explicit `grid` byte marker and makes accidental array confusion more plausible.

### Option 4: `grid` becomes transport framing

The transport frame says "this frame carries a PromiseGrid message," but the
content block is DAG-CBOR. This separates transport and content identity, but
creates a signable-view and message-CID question.

## Recommended Direction

The TE recommends Alt C as the working direction:

> Keep `grid([42(pCID), ...])` as the PromiseGrid wire/app/kernel envelope for
> now, while moving durable POC18 graph objects toward true DAG-CBOR where the
> CID content multicodec, IPLD traversal, CAR export, and ecosystem tooling are
> valuable. Keep raw file chunks as `raw` CIDs unless encrypted-chunk work later
> decides otherwise.

This means:

- raw chunks should remain CIDv1 `raw` objects;
- Markdown spec docs may remain CIDv1 `raw` objects;
- durable manifests, reference sets, snapshots, review threads, and mapping
  records are candidates for CIDv1 `dag-cbor`;
- wire envelopes may remain PromiseGrid-CBOR with a custom `grid` tag;
- `42(...)` links should stay the standard link representation;
- diagnostics must stop implying that every `.cbor` object is DAG-CBOR;
- the CAS index should record both local kind and CID content multicodec.

## Rejected Directions

Alt B is rejected as an immediate conversion because it would force the custom
`grid` tag question too early and risks reopening envelope work before POC18 has
peer sync, Git bridge, GC, and encrypted chunks.

Alt D is not rejected forever, but it should require a separate envelope
supersedence TE because it changes the settled `grid(...)` byte shape.

Alt E is risky unless identity and signatures are specified precisely. It should
not be adopted casually.

## Recommended DF Questions

1. **Durable graph objects:** Should POC18 convert durable manifests,
   reference sets, snapshots, reviews, and mapping records to CIDv1 `dag-cbor`
   while keeping wire envelopes as PromiseGrid-CBOR?
   - Recommended: yes.

2. **Wire envelope:** Should the custom `grid` tag remain outside true DAG-CBOR
   for now?
   - Recommended: yes.

3. **CID profiles:** Should the store support multiple CID content multicodecs:
   at least `raw` and `dag-cbor`?
   - Recommended: yes.

4. **Spec docs:** Should Markdown protocol specs remain raw pCID objects until a
   later bundled-spec TE decides otherwise?
   - Recommended: yes.

5. **Terminology:** Should docs distinguish "DAG-CBOR block" from
   "PromiseGrid-CBOR envelope with tag-42 links"?
   - Recommended: yes, immediately.

## Implications For POC18

- Add a DF lock before changing POC18 storage code.
- Add analyzer checks that report CID content multicodec per object.
- Add diagnostic output that says `raw`, `dag-cbor`, or `promisegrid-cbor`.
- Update TE-givul or file a follow-up DF so chunk identity and DAG-CBOR object
  identity are not conflated.
- Consider converting standalone chunk manifests and reference sets before
  converting full grid envelopes.
- Keep POC6 as the executable reference for true DAG-CBOR.

## Decision Status

Needs DF. The recommended surviving path is a hybrid: true DAG-CBOR for durable
IPLD-shaped graph objects, raw CIDs for raw chunks and Markdown spec docs, and
PromiseGrid-CBOR `grid(...)` envelopes for wire/app/kernel messages until a
separate envelope TE proves a better all-DAG-CBOR envelope shape.
