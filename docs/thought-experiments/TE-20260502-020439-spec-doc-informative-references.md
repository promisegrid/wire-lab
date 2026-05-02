# TE-33: Spec-doc Informative References to its workshop, RFC-shaped

**Status:** Draft. Refines TE-31 and TE-32.
**Drafted:** 2026-05-02 02:04 UTC

## Premise

In the conversation drafting TE-31 and TE-32, the bot initially
asserted that a spec doc must not reference its workshop at all.
Steve pushed back: "are you SURE that RFCs can't in any way reference
the background or workshop docs that produced them? how does the
IETF do this?" The bot's first answer was wrong.

The IETF in fact has a well-developed practice for spec-to-workshop
references, captured by:

- [RFC 7322 (RFC Style Guide)](https://www.rfc-editor.org/rfc/rfc7322) -
  RFCs MUST split references into Normative and Informative subsections.
- [RFC 3967](https://datatracker.ietf.org/doc/html/rfc3967) -
  Standards-track RFCs MAY NOT have a normative reference to a
  non-standards-track or lower-maturity document, but MAY have
  informative references to such documents (Internet-Drafts, BCPs,
  workshop output, etc).
- [draft-carpenter-rfc-citation-recs](https://datatracker.ietf.org/doc/html/draft-carpenter-rfc-citation-recs-01) -
  RFCs MAY cite Internet-Drafts informatively, with the draft filename
  and version pinned, marked as "Work in Progress" or "Working Draft".

PromiseGrid spec docs should adopt the same convention, adapted to a
content-addressed world. The IETF setup is actually weaker than what
PromiseGrid can achieve, because Internet-Drafts are version-flaky
(the same `draft-foo-bar` URL points at different content over time),
whereas a PromiseGrid workshop reference by tree-hash or pCID is
forever-stable.

## The three locks

### Lock 1: Two-section bibliography, RFC-shaped

A spec doc MAY have a References section. If it does, it MUST split
into two subsections, in this order:

1. **Normative References.** Documents that must be read or
   implemented for the spec to be implementable. By Lock 2 below,
   normative references are restricted to other content-addressed
   spec doc-CIDs and named external standards (RFC numbers, ISO
   numbers, IEEE numbers, etc.). Normative references MUST NOT
   include workshop pointers.

2. **Informative References.** Background, historical, or
   provenance pointers. May include workshop output, prior art,
   tutorial material. By Lock 3 below, workshop references MUST
   be content-addressed.

This matches RFC 7322 Section 4.8.6 ("References").

A spec doc with no references section is permitted (small specs may
not need one). A spec doc with only one section MUST use the
top-level heading "Normative References" or "Informative References"
to disambiguate, per RFC 7322.

### Lock 2: Normative references are tightly constrained

A normative reference in a spec doc MUST be one of:

- Another spec doc-CID (a fully-qualified pCID), optionally
  accompanied by a human-readable title and date for convenience.
- A named external standard with a stable, registry-backed identifier
  (e.g. "RFC 768", "ISO 8601", "IEEE 754", "Unicode 15.0").
- A frozen `protocols/<slug>-<pcid>.md` doc-CID from this repo or any
  equivalent repo elsewhere.

A normative reference MUST NOT be:

- A workshop tree-hash, branch ref, or commit ID.
- A live `protocols/<slug>.d/` path.
- A draft spec doc that has not yet had a freeze event recorded in
  its A-side CHANGELOG (per TE-32).
- A URL to a moveable resource.

The intent: a normative reference is part of the spec's promise.
Anything load-bearing for implementation MUST be content-stable or
have a registry behind it that promises stability.

### Lock 3: Informative references must be content-addressed

An informative reference in a spec doc MAY be any of:

- Anything permitted as a normative reference (Lock 2).
- A git tree-hash of a workshop directory at a specific moment.
  Format: `git tree-hash sha256:abc123... (protocols/<slug>.d/ at
  freeze of doc-CID Y)`. Git tree-hashes are content-addressed by
  git's design, so they are stable forever.
- The pCID of a specific TE document.
- The pCID of a frozen workshop tree (`protocols/<slug>-<pcid>.d/`).
- A long-lived URL on a domain whose stability is well-known
  (e.g. `rfc-editor.org`, `ietf.org`, `w3.org`, `iana.org`).

An informative reference MUST NOT be:

- A branch ref (`main`, `master`, `ppx/main`, any moveable ref).
- A commit ID alone, without explicit tree-hash anchoring (commits
  carry parent metadata that is not stable across rebases; the tree
  is the stable part).
- A URL to a resource on a domain known to be ephemeral (personal
  blog, paste-bin, social media, transient hosting).
- A non-content-addressed mailing-list archive URL (these are
  permitted only when the archive itself is content-addressed or
  cryptographically signed; vanilla mbox URLs are not stable enough).

Note that the IETF allows non-content-addressed mailing-list URLs and
ephemeral Internet-Draft references; PromiseGrid is stricter because
the inversion (TE-31) and the long-time-horizon constraint (C-2 in
TE-28) require that spec docs remain meaningfully resolvable
decades from now.

## Bidirectional pointers, finally

Combining TE-32 and this TE produces a clean two-way trail between
spec and workshop, both directions content-addressed:

- **Workshop -> spec:** A-side `CHANGELOG.md` `freeze` entry records
  the doc-CID produced. This is machine-checkable: each freeze is
  a row with a pCID. (TE-32, locked.)
- **Spec -> workshop:** Informative References section in the spec
  doc may record a workshop tree-hash or frozen `.d/` pCID. This is
  human-readable and auditable but not load-bearing. (This TE.)

The bidirectionality is symmetric in shape but asymmetric in
authority:

| Direction | Where | Authority | Purpose |
|---|---|---|---|
| Workshop -> spec | A-side CHANGELOG | spec maintainers | machine-checkable freeze record |
| Spec -> workshop | Informative References | spec author | human-readable provenance |

Neither direction is load-bearing for implementation conformance
(B-side CHANGELOG entries reference the spec doc-CID upstream; that
is the only conformance-critical reference per TE-31).

## Why not just one direction?

We considered three options before locking the bidirectional shape:

- **Option A: spec -> workshop only** (the IETF's effective default,
  via Acknowledgments and Informative References). Loses
  machine-checkability of the freeze trail.
- **Option B: workshop -> spec only** (the bot's first wrong answer
  in conversation; was the implicit shape of TE-32 alone). Loses
  the human-readable provenance trail that helps a future
  contributor understand where a spec came from.
- **Option C: bidirectional with different authority levels**
  (this TE). Combines machine-checkability with human-readable
  provenance.

Option C costs one extra section in the spec doc. The cost is
negligible; the benefit is cumulative as the project ages.

## Worked example

A future `udp-binding` v1 spec doc, reaching freeze:

```
# UDP-binding v1

[... spec body, including ABNF, state tables, message formats ...]

## Normative References

- [PromiseGrid Group-Session v1](bafkrei...group-session-v1-cid) -
  This binding carries Group-Session messages. Implementers MUST
  read the Group-Session spec.
- [RFC 768](https://www.rfc-editor.org/rfc/rfc768) - User Datagram
  Protocol. The transport this binding rides on.
- [RFC 5234](https://www.rfc-editor.org/rfc/rfc5234) - ABNF for
  Syntax Specifications. The notation used in section 3.

## Informative References

- [TE-29: Protocols as simulated repos and the L4-binding layer](bafkrei...te-29-cid) -
  The thought experiment that established the binding layer concept.
- Workshop tree at freeze: git tree-hash
  `sha256:abc123def456...` (snapshot of `protocols/udp-binding.d/`
  at the moment the v1 doc-CID was produced; recorded also in the
  A-side CHANGELOG entry for cross-checking).
- [TE-31: Spec-doc as upstream, simrepo as implementation](bafkrei...te-31-cid) -
  The conformance-direction inversion.
- [draft-carpenter-rfc-citation-recs-01](https://datatracker.ietf.org/doc/html/draft-carpenter-rfc-citation-recs-01) -
  IETF practice for citing workshop output, which this spec follows.
```

Note that the workshop tree-hash appears in two places: in the
A-side CHANGELOG (machine-checkable, authoritative for "this is the
freeze record") and in the spec's Informative References
(human-readable, advisory). Cross-checking the two is a smoke-test
the harness can perform.

## Closes / partially closes

- **Closes nothing previously open.** This TE adds capability without
  retracting prior decisions.
- **Refines TE-31 and TE-32.** Adds the spec-side half of the
  bidirectional trail; preserves the A-side CHANGELOG as the
  authoritative freeze record.

## Surfaced questions

- **OQ-33.1: Acknowledgments section.** RFCs almost always have an
  Acknowledgments section thanking contributors. Should PromiseGrid
  spec docs adopt this convention? Lean: yes, as a convention not a
  lock; treat it as informative-references-for-people. Leave the
  format to the spec author.

- **OQ-33.2: Workshop tree-hash format.** When recording a git
  tree-hash in either CHANGELOG or Informative References, use git's
  native SHA-1 (legacy compatibility) or SHA-256 (matches PromiseGrid
  CIDv1 hash family)? Lean: prefer SHA-256 once git tooling supports
  it widely; SHA-1 acceptable as a transitional measure. Worth a
  follow-on note, not a separate TE.

- **OQ-33.3: Mailing-list-equivalent for PromiseGrid.** PromiseGrid
  has no mailing list; the equivalent is `proposals/` and
  (eventually) ppx-dr message exchanges. Should those exchanges be
  citable as Informative References? Lean: yes, via pCID once
  proposals are migrated to ppx-dr message-protocol shape per TODO 016.
  Until that migration, treat `proposals/` paths as not-yet-citable.

- **OQ-33.4: Cross-repo informative references.** A spec doc in this
  repo may want to reference a workshop tree from another repo
  (e.g. an external PromiseGrid implementation that helped shape
  the spec). Tree-hashes are content-addressed and global, so this
  works in principle, but the reader needs to know which repo to
  fetch from. Lean: include the repo URL alongside the tree-hash,
  with the understanding that the URL is a hint and the tree-hash
  is the authority.

- **OQ-33.5: Frozen Informative References.** When a spec doc itself
  is frozen and gets a doc-CID, the Informative References section
  is part of its content. If a referenced workshop later evolves
  beyond the cited tree-hash, the citation remains valid (the
  tree-hash is immutable) but may become harder to fetch if the
  workshop has been pruned. Lean: tolerate this; it is the same
  problem as broken URLs in any RFC. The freeze captures the
  reference; the world's preservation of the referenced content is
  a separate problem (potential future TE on archival policy).

## Verdict

Adopt RFC-shaped two-section bibliography for spec docs. Lock 1 (split
into Normative and Informative). Lock 2 (normative references are
tightly content-addressed). Lock 3 (informative references may include
workshop pointers, but only content-addressed forms). Bidirectional
trail spec<->workshop is now standard, with A-side CHANGELOG as the
machine-checkable side and Informative References as the
human-readable side.
