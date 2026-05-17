# Question

Which peer-local promise accounting records should PromiseGrid simulations
include so Alice, Bob, Carol, and other peers can make sparse-CAS pull, keep,
advertise, and refusal decisions without a central accounting authority?
Source: `DI-navod`.

Open decision points:

- What observations are necessary before a peer changes future behavior toward
  another peer or site?
- Which decisions belong at L7 group semantics, which belong at L6 CAS, and
  which belong at L5 feed replication?
- How can the record shape stay understandable enough for the 100-year
  layperson mental model?
- What should remain out of scope until a later economics or reputation
  specimen exists?
