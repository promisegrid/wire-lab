# TE-givul: POC18 chunk storage identity

TE ID: TE-givul

## Status

needs DF

## Decision Under Test

Should POC18 store regular-file chunks as raw `.bin` byte objects under
`chunks/`, or should chunks be stored as CBOR `.cbor` objects under `objects/`?

This is a storage-identity decision, not just a filename decision. The real
question is what a chunk CID names:

- exact raw file bytes,
- a CBOR byte-string wrapper around exact file bytes,
- a typed CBOR chunk object,
- a full `grid(...)` promise message about chunk bytes,
- encrypted chunk bytes, or
- a locally packed/compressed/encrypted storage representation that is not the
  CID identity.

## Existing Locked Inputs

- POC18 uses Rabin content-defined chunking for all regular-file content.
  Source: `DI-dofoj`.
- POC18 represents regular files by chunk manifests and PromiseGrid CAS objects,
  not by Git blobs alone. Source: `DI-dofoj`; `DI-zuruj`.
- Printable object identities are CIDv1 base32; binary CID values remain binary
  on the wire. Source: `DI-jifuj`; `DI-harih`.
- POC18 agents have sparse CAS stores. No agent has complete global storage.
  Source: `TE-vahoj`; `DI-jifuj`.
- POC18 native collaboration is continuous peer DAG synchronization, with
  Git-style push/pull only as bridge behavior. Source: `DI-dibut`; `DI-dofoj`.
- POC18 object semantics must remain promise-shaped and peer-relative, not
  RPC-shaped or authority-shaped. Source: `DI-zuruj`; `DI-harih`.

## Current POC18 Behavior

The first POC18 slice stores chunks as raw files:

```text
/tmp/wire-lab-poc18-run/alice-cas/
  chunks/
    bafkrei....bin
  objects/
    bafkrei....cbor
  index.json
```

A chunk `.bin` file is exactly the byte range produced by the Rabin chunker. It
has no CBOR tag, no CBOR byte-string header, no `grid(...)` tag, no pCID, no
embedded length, no embedded file offset, no signature, and no promise body. Its
CID is currently CIDv1 raw/sha2-256 over those exact bytes.

The references are layered:

```text
regular POSIX node message
  -> 42(chunk_manifest_message_cid)
        -> chunk_manifest promise body
             -> 42(chunk_manifest_object_cid)
             -> chunk rows [offset, length, 42(chunk_cid)]
                   -> raw chunk bytes
```

There are two manifest forms in the first slice:

- A standalone CBOR chunk manifest object.
- A `chunk_manifest` grid message that promises a chunk list and points at the
  standalone manifest object.

That duplication is useful fixture evidence but is not necessarily the final
shape. It should be scrutinized with the chunk-storage decision.

## Prior-Art Questions

The TE must compare systems along four axes:

1. What bytes are hashed?
2. What bytes are stored locally on disk?
3. Is local on-disk wrapping part of object identity?
4. Are objects fetched by content identity, name, placement rule, or some other
   authority?

### IPFS, IPLD, and Kubo

IPFS separates content identity from local datastore layout. A CID identifies
block data together with a multicodec/multihash interpretation. Local Kubo
datastores, including `flatfs`, are implementation choices for storing key/value
block data; local file layout, sharding, and datastore wrapping are not the CID's
semantic content. Reference: Kubo datastore documentation
`https://github.com/ipfs/kubo/blob/master/docs/datastores.md`.

For POC18, the useful lesson is:

- If a chunk is an IPLD block with `raw` codec, the CID names exact raw bytes.
- If a chunk is a DAG-CBOR block, the CID names canonical encoded CBOR bytes.
- If a chunk is a DAG-PB or other structured block, the CID names that encoded
  structured block.
- Local datastore wrappers, filenames, compression, or shard directories should
  not change protocol identity.

IPFS therefore does not answer "always wrap chunks" or "never wrap chunks." It
answers: choose the block's canonical bytes and codec deliberately, then keep
local storage layout separate from identity.

### Git

Git does not hash only raw file bytes for normal blob objects. A Git object ID is
over canonical object bytes shaped as:

```text
<type> SP <size> NUL <content>
```

Loose objects are then compressed for local storage, and packfiles may delta-pack
objects differently, but that local compression/packing is not what the object
ID names. Reference: Git Book, "Git Internals - Git Objects"
`https://git-scm.com/book/en/v2/Git-Internals-Git-Objects`.

For POC18, the useful lesson is Git's type/size domain separation:

- Git object IDs are not over incidental local storage wrappers.
- Git does include minimal typed canonical structure before hashing.
- Git can repack storage without changing logical object IDs.

This supports a possible POC18 alternative where chunk identity is over
canonical typed bytes, but it does not support making the CID depend on a
temporary filesystem extension or local packing format.

### Ceph / RADOS

Ceph's RADOS layer is distributed object storage, but it is not usually a
content-addressed system in the Git/IPFS sense. Objects are named, mapped to
placement groups, replicated or erasure-coded by OSDs, and stored through the
configured OSD backend. Checksums and backend layout protect storage, but a RADOS
object name is not normally the CID of exact object bytes. Reference: Ceph
architecture documentation `https://docs.ceph.com/en/latest/architecture/`.

For POC18, the useful lesson is negative: distributed object storage can be
excellent without making local disk bytes the protocol identity. Ceph-like
placement, replication, and repair ideas may matter later, but POC18's CAS
identity should remain explicit and content-derived.

### OCI Images

OCI images distinguish descriptors, media types, digest algorithms, byte sizes,
manifests, config objects, and layer blobs. A layer digest commonly names the
exact blob bytes being transferred and stored, while separate metadata says how
to interpret the blob. OCI also has the notion of uncompressed layer identity
(`diffID`) distinct from compressed blob digest.

For POC18, the useful lesson is that "same file content" and "same transferred
object" can be different identities if compression or encryption is involved.
If POC18 later stores compressed or encrypted chunks, it must explicitly decide
whether CIDs name plaintext, ciphertext, compressed bytes, or a manifest object
that links between them.

### Nix, Venti, Tahoe-LAFS, Borg, Restic, ZFS, and Btrfs

Other systems reinforce the same split:

- Nix uses content and derivation information to make store paths reproducible,
  but the store path is not just a raw file chunk hash.
- Venti-style archival CAS stores blocks by content hash and treats block
  identity as exact content identity.
- Tahoe-LAFS and backup systems such as Borg/restic typically combine chunking
  with encryption, packing, retention, and repair semantics, so local storage
  bytes and logical content identity may intentionally differ.
- ZFS and Btrfs use checksums to protect blocks inside a storage pool, but those
  checksums are not a portable inter-agent CAS protocol.

For POC18, the recurring lesson is: separate protocol identity from local
storage optimization; do not confuse disk layout with the promise an object makes.

## Alternatives

### Alt A: Raw chunk bytes in `chunks/*.bin`

The chunk CID names the exact raw file byte range produced by the Rabin chunker.
The chunk manifest stores `[offset, length, 42(chunk_cid)]`. Local storage uses
`chunks/<cid>.bin`.

This is the current implementation.

Advantages:

- Maximum deduplication for identical plaintext chunk bytes across files,
  workspaces, and agents.
- Simple reassembly: fetch the bytes named by the chunk CID and concatenate in
  manifest order.
- Small-device friendly: no CBOR decoder is needed to validate the raw chunk
  length after fetching.
- Closest to IPFS `raw` blocks and Venti-like exact-content CAS.
- Local storage can later pack/compress raw chunks without changing identity if
  the store retains a mapping from CID to canonical bytes.

Disadvantages:

- A raw chunk has no embedded type marker or domain separation.
- A chunk viewed alone does not tell the reader whether it is file content,
  encrypted content, a manifest accidentally stored as bytes, or arbitrary data.
- Low-entropy chunks can be guessed by an attacker who can test candidate CIDs.
- Debugging raw chunks is less self-documenting than rendering CBOR.

New obligations:

- The manifest and file-node messages must carry the interpretation.
- The CAS index must retain object kind or equivalent local metadata.
- Privacy-sensitive content needs encryption or access promises, not secrecy of
  raw CIDs.

### Alt B: CBOR byte-string chunk objects in `objects/*.cbor`

The chunk CID names canonical CBOR bytes whose payload is a byte string:

```text
h'<raw chunk bytes>'
```

Local storage uses `objects/<cid>.cbor`.

Advantages:

- Every stored object in `objects/` is syntactically CBOR.
- Diagnostic tooling can uniformly parse the top-level object.
- There is minimal structural domain separation versus raw bytes.

Disadvantages:

- The CID no longer names the raw file byte range; it names CBOR encoding of that
  byte range.
- Any peer or tool wanting raw bytes must know to CBOR-decode first.
- The CBOR byte-string header adds little semantic value because length is
  already in the manifest and the file itself supplies content semantics.
- This is not naturally equivalent to IPFS raw blocks.

New obligations:

- The manifest must say whether chunk links point at raw bytes or CBOR byte
  strings.
- Git bridge code must translate between Git blob content and CBOR chunk content
  carefully.

### Alt C: Typed CBOR chunk records in `objects/*.cbor`

The chunk CID names canonical CBOR such as:

```text
[
  "poc18_chunk",
  chunker_name,
  chunker_parameters,
  raw_chunk_bytes
]
```

Advantages:

- Strong domain separation.
- Self-describing object contents.
- Room for future fields such as compression, encryption profile, or provenance.

Disadvantages:

- Duplicates metadata already in the manifest.
- Expands every chunk, including very small chunks.
- Makes the same raw bytes have different CIDs under different chunker metadata,
  which harms deduplication.
- Risks pushing per-file or per-manifest semantics down into content chunks.

New obligations:

- Define canonical chunk record shape in the pCID spec.
- Decide whether chunker parameters belong on every chunk or only in manifests.
- Define whether identical raw chunks from different chunkers should deduplicate.

### Alt D: Full grid chunk-promise messages

Each chunk is stored as a `grid([42(pCID), parents, payload, proof])` message
whose payload promises something about the bytes.

Advantages:

- Every chunk is a promise-bearing message.
- Per-chunk signatures, parents, and retention/economics hooks are available.

Disadvantages:

- Heavyweight for large files with many chunks.
- Confuses data chunks with promises about data chunks.
- Adds signing and verification work to the hottest storage path.
- Makes raw byte deduplication much worse unless the raw bytes are stored
  separately anyway.

New obligations:

- Define who promises each chunk and why.
- Define whether a file can be reconstructed from unsigned raw bytes or only from
  signed chunk promises.
- Define how per-chunk promise churn interacts with GC and retention economics.

### Alt E: Raw canonical chunks, but stored under `objects/`

The chunk CID still names exact raw bytes, but local storage uses
`objects/<cid>.bin` or `objects/<cid>` rather than a separate `chunks/` directory.

Advantages:

- Preserves raw-byte identity.
- Simplifies CAS directory count.
- Keeps local layout closer to "all CIDs are objects."

Disadvantages:

- Mixed extensions under one directory can make diagnostics and tooling more
  error-prone.
- A single `objects/` directory no longer implies `.cbor` parseability.
- Does not answer the semantic question; it only changes local layout.

New obligations:

- The index or extension must clearly state canonical byte type.
- Diagnostic tooling must not assume every object path is CBOR.

### Alt F: Encrypted chunks with CIDs over ciphertext

The chunk CID names encrypted bytes. Plaintext chunks are never stored or
advertised by CID except inside a local trusted boundary.

Advantages:

- Prevents direct guess-confirm tests against low-entropy plaintext chunks.
- Fits private repositories and cross-legal-entity storage promises.
- A storage peer can retain and serve bytes without seeing plaintext.

Disadvantages:

- Cross-agent deduplication of plaintext is lost unless convergent encryption is
  used, and convergent encryption reintroduces equality leakage.
- Chunk manifests need encryption profiles and key-capability promises.
- Git bridge and checkout paths become more complex.

New obligations:

- Define encryption metadata, key exchange, token promises, and revocation.
- Decide whether parent links and manifests are visible, encrypted, or split.
- Decide whether dedup is by plaintext, ciphertext, or both.

## Scenario Analysis

### Scenario 1: Alice stores a normal source tree

Alice ingests a source tree with text files, symlinks, and generated binary
assets. Bob later materializes the snapshot.

- Alt A is simplest: Bob fetches manifests, fetches raw chunks, checks CIDs and
  lengths, and writes files.
- Alt B requires CBOR decoding for every chunk without adding useful source-tree
  semantics.
- Alt C makes chunks easier to inspect but duplicates manifest data.
- Alt D is unnecessary overhead.
- Alt E is mostly layout style.
- Alt F is unnecessary unless the source tree is private and stored by untrusted
  peers.

Best fit: Alt A or Alt E, with manifests carrying interpretation.

### Scenario 2: Large binary file changes slightly

Carol changes a few bytes in a large file. Rabin chunking preserves most chunk
boundaries, so only a few chunks should change.

- Alt A maximizes reuse because identical raw chunks retain identical CIDs.
- Alt B also reuses chunks if canonical CBOR wrapping is stable, but the stored
  identity is not the raw byte range.
- Alt C may break reuse if metadata changes.
- Alt D adds many signatures and messages for little benefit.
- Alt F deduplicates only under the chosen encryption policy.

Best fit: Alt A.

### Scenario 3: Sparse peer retrieval

Dave has a manifest but lacks one chunk. He asks Ellen whether she promises to
serve the missing CID.

- Alt A makes the request exact: "Do you promise to send bytes matching this raw
  chunk CID?"
- Alt B changes that to "send the CBOR object that decodes to the bytes."
- Alt C and Alt D require the receiver to parse more structure before getting
  file bytes.
- Alt F requires Dave to also have decryption capability promises.

Best fit: Alt A for public chunks; Alt F for private chunks.

### Scenario 4: CAS GC and retention

Frank is low on disk. He decides which objects he still promises to retain.

- Alt A needs manifests and reference sets to discover raw chunk reachability.
- Alt B/C/D make every chunk parseable as CBOR, but reachability still comes from
  manifests and reference sets.
- Alt E has the same semantics as Alt A.
- Alt F adds key-retention and encrypted-object-retention questions.

Best fit: Alt A with clear manifest reachability and local retention promises.

### Scenario 5: Debugging and diagnostics

Grace wants to inspect stored objects after a failed checkout.

- Alt A requires diagnostics to treat chunks as opaque bytes and inspect the
  manifest for meaning.
- Alt B/C/D give more self-description at the chunk object itself.
- Alt E depends on extensions or index metadata.
- Alt F gives little direct visibility without keys.

Best fit: Alt C for human diagnostics, but the operational value probably does
not justify making every chunk self-describing.

### Scenario 6: Low-entropy private content

Heidi stores a file containing a predictable secret, such as a short token or a
boilerplate contract with one guessed clause.

- Alt A leaks equality and supports guess-confirm attacks if an attacker can
  learn or query chunk CIDs.
- Alt B and Alt C do not solve this; deterministic wrapping still lets attackers
  hash guesses if they know the wrapper.
- Alt D does not solve this if the bytes are still deterministic and public.
- Alt F is the relevant alternative: encrypt before advertising or delegating
  storage.

Best fit: Alt F for private/low-entropy data; do not pretend CBOR wrapping is
privacy.

### Scenario 7: Git bridge round-trip

Ivan imports a Git repository, edits in PromiseGrid, and exports back to Git.

- Alt A maps naturally from raw Git blob content to raw file bytes and chunk
  manifests, even though Git object IDs include Git's type/size header.
- Alt B/C require wrapping/unwrapping chunks during bridge conversion.
- Alt D makes Git bridge behavior unnecessarily heavy.
- Alt F requires explicit encrypted-repo semantics and is not a transparent Git
  bridge.

Best fit: Alt A for public Git compatibility.

### Scenario 8: IPFS/IPLD interop

Judy wants to serve POC18 chunks through IPFS-like tooling or compare POC18 chunk
CIDs against raw IPLD blocks.

- Alt A is closest to raw block interop.
- Alt B/C are valid IPLD-style objects only if the chosen multicodec and CBOR
  canonicalization are explicit.
- Alt D is a PromiseGrid message, not a raw file-content block.
- Alt F can interoperate as opaque encrypted bytes, not plaintext chunks.

Best fit: Alt A for raw-block interop; Alt C only if typed chunk objects become
load-bearing.

## Salt, Preimage, and Privacy Analysis

### Preimage attacks

With SHA-256 multihashes, ordinary cryptographic preimage attacks are not the
main concern for POC18 chunk identity. An attacker should not be able to find
bytes matching an arbitrary unknown chunk CID.

### Guess-confirm attacks

The practical risk is different: if content is low entropy or predictable, an
attacker can guess candidate bytes, hash them, and compare the resulting CID to a
known or queried chunk CID. This is not a break of SHA-256; it is an information
leak caused by deterministic public content addressing.

Examples:

- A one-word file.
- A known template with a few possible values.
- A config file whose secret has low entropy.
- A private repository where file equality itself leaks useful information.

### Salting

Salting raw chunks would prevent simple cross-agent guess-confirm checks, but it
would also destroy the main value of public CAS:

- identical chunks would no longer deduplicate across agents,
- peers could not request chunks by stable content identity,
- manifests would need salt distribution and salt-retention semantics,
- Git/IPFS-style interoperability would weaken, and
- chunk CIDs would become local or group-relative rather than universal.

Salt is therefore not a general answer for POC18 public chunk identity. It may be
appropriate for a private, group-scoped, or encrypted profile, but then the TE
must be honest that the object identity is no longer universal raw-content
identity.

### Encryption

For private content, encryption is the stronger PromiseGrid-shaped answer:

- Alice promises storage peers only encrypted chunk bytes.
- Alice issues capability tokens or key promises only to peers she locally
  trusts.
- The CID can name ciphertext bytes for storage and retrieval.
- A manifest can separately promise how authorized agents reconstruct plaintext.

Encryption preserves the distinction between storage promises and read promises.
It also keeps trust local: no global authority decides who may read; Alice
decides what she sends and what key/capability promises she makes.

## Comparison Table

| Alternative | CID names | Stored where | Interop | Dedup | Privacy | Complexity |
|---|---|---|---|---|---|---|
| Alt A raw chunks | raw chunk bytes | `chunks/*.bin` | strong for raw-block/IPFS-like use | strongest | weak for low-entropy public CIDs | low |
| Alt B CBOR byte string | CBOR encoding of bytes | `objects/*.cbor` | valid CBOR, weaker raw interop | strong if wrapper stable | weak | low-medium |
| Alt C typed CBOR chunk | typed CBOR record | `objects/*.cbor` | structured IPLD possible | weaker if metadata varies | weak | medium |
| Alt D grid chunk message | signed promise message | `objects/*.cbor` | PromiseGrid-native only | weak unless raw bytes also stored | depends on payload | high |
| Alt E raw chunks under objects | raw chunk bytes | `objects/*.bin` or `objects/<cid>` | same as Alt A | strongest | weak | low |
| Alt F encrypted chunks | ciphertext bytes | either layout | opaque-byte interop | policy-dependent | strongest | high |

## Conclusions

The TE rejects Alt D as the normal file-content chunk representation. A chunk
should not usually be a full promise message. Promise messages should describe,
link, retain, advertise, price, verify, or authorize access to chunks; they
should not be required for every hot-path byte range.

The TE also rejects CBOR wrapping as a privacy mechanism. CBOR wrapping does not
solve guess-confirm attacks. If privacy matters, POC18 needs encryption and local
capability/key promises.

The strongest default is Alt A:

> Public/plain POC18 chunk CIDs should name exact raw chunk bytes. Manifests and
> grid messages should carry interpretation, chunker parameters, file offsets,
> retention promises, and trust/economics; local disk layout should remain an
> implementation detail.

Alt E remains a surviving layout-only alternative if the team wants one object
directory, but it should not change raw-byte identity. Alt F remains necessary as
a future private/encrypted profile.

## Recommended DF Questions

1. **Public chunk identity:** Should public/plain POC18 chunks use Alt A raw-byte
   identity?
   - Recommended: yes.
   - Reject if: chunk self-description is more important than raw-byte
     deduplication and IPFS-like interop.

2. **Local layout:** Should raw chunks stay under `chunks/*.bin`, or move under
   `objects/` while remaining raw bytes?
   - Recommended for now: keep `chunks/*.bin` until diagnostics and GC mature.
   - Later migration: move to a unified object directory only if the index and
     tooling clearly distinguish raw bytes from CBOR.

3. **Encrypted profile:** Should POC18 add a distinct encrypted-chunk profile
   where CIDs name ciphertext bytes?
   - Recommended: yes, but as a later explicit task, not by changing the public
     raw chunk default.

4. **Manifest duplication:** Should POC18 keep both standalone chunk-manifest
   CBOR objects and `chunk_manifest` grid messages?
   - Recommended: review separately. It may be enough for the grid message to be
     the canonical manifest, or for the grid message to promise a separately
     reusable manifest object.

## Implications For TODO-nahop

- Add a DF task to lock chunk identity before implementing peer retrieval,
  retention, GC, Git bridge, or encrypted chunks.
- Add analyzer gates that distinguish raw chunk CIDs from CBOR object CIDs.
- Add a future encrypted-chunk TE or subtask before claiming private-repository
  security.
- Keep Git bridge design aware that Git object IDs hash typed Git object bytes,
  while POC18 chunk CIDs may hash raw file byte ranges.
- Keep IPFS/IPLD interop language precise: raw chunk blocks and DAG-CBOR objects
  are both valid content-addressing patterns, but they identify different bytes.

## Decision Status

Needs DF. The TE recommends raw public chunk identity, with manifests/grid
messages carrying semantics, but the layout and encrypted-profile questions
remain open.
