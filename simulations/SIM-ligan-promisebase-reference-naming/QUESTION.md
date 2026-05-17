# Question

How should PromiseGrid handle human-readable references to CAS roots learned
from promisebase without creating competing identities, mutable-name confusion,
or custom non-CID syntax?
Source: `DI-tibis`.

Open decision points:

- Does reference naming belong in L6 CAS objects, L7 group/session metadata, a
  separate reference protocol, or nowhere in the first profile?
- Are references immutable labels, mutable refs with signed history, local-only
  nicknames, or discoverable shared names?
- How does a reader distinguish the reference name, the reference object's CID,
  a pointer-object CID, and the target root CID?
- What happens when Alice and Bob use the same human-readable symbol for
  different roots, or when Mallory publishes a confusingly similar name?
- Which promisebase reference lessons should be adopted, wrapped, rejected, or
  preserved only as prototype history?
