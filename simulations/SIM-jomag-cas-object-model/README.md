# SIM-jomag: CAS object model

This simulation explores the L6 CAS object-model questions exposed by turn
177. It asks how PromiseGrid should encode pointer objects, raw chunks, and
Merkle nodes so independent sites compute the same CIDs for the same bytes. It
is a standalone design-point simulation, not a frozen L6 CAS spec and not a
shared protocol bundle for other simulations. Source: `DI-navod`; `DI-tibis`.

## Question

Which L6 CAS object model lets peers exchange and verify content-addressed
objects without confusing raw chunks, Merkle nodes, pointer objects, and
application payloads? Source: `DI-navod`.

## Turn 177 pressure

Turn 177 made several CAS details load-bearing rather than cosmetic:

- PromiseGrid messages should use CBOR-shaped bytes, not JSON-shaped bytes.
- Pointer files should be first-class CAS objects, not filesystem symlinks.
- Rabin or FastCDC chunking changes CID stability unless parameters are locked.
- CIDv1 codec fields should distinguish object types instead of filename
  suffixes.
- Promisebase / pitbase prior art should be evaluated without treating it as
  canonical PromiseGrid design.

This simulation keeps those questions visible as specimen pressure while
`TODO-kituj` / TE-43 owns the actual concrete L6 CAS decision path. Source:
`DI-navod`; `DI-tibis`.

## Decision axes

- **CBOR profile:** deterministic CBOR, DAG-CBOR, allowed tags, and
  text-string versus byte-string boundaries.
- **Object typing:** CIDv1 codec / multicodec choices for raw chunks, Merkle
  nodes, pointer objects, and any future application-shaped objects.
- **Pointer shape:** whether a pointer object is a minimal root-CID record, a
  richer CBOR map, or a typed object whose codec carries most of the meaning.
- **Chunking:** Rabin versus FastCDC, target average size, min/max bounds,
  window size, and how those parameters become part of the protocol contract.
- **Prior art:** what to adopt, reject, or adapt from promisebase / pitbase.

## Boundaries

This simulation does not choose a winning grid-envelope variant, does not
rewrite existing `.txt` specimens, and does not assert that all PromiseGrid CAS
implementations must use one specific local directory layout. It exists so
CAS-object design can evolve independently before any frozen pCID spec or guide
prose claims a stable answer. Source: `DI-navod`; `DI-tibis`.
