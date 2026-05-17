# Question

What wire object or promise should mean: "I, peer P, promise to behave as spec
pCID X with open-question answers Q7=yes and Q9=variant-B"? Source:
`DI-pukap`; `TODO-nivus`.

Open decision points:

- Is adoption metadata a signed structured object, an ordinary promise message,
  a spec-manifest entry plus peer answer binding, or a hybrid?
- How are open-question answer keys named so they remain stable across spec
  freezes and later superseding specs?
- How does Alice verify Bob's adoption claim before depending on Bob's protocol
  behavior?
- How does a peer update or revoke an adoption promise without mutating the old
  content-addressed record?
- Does adoption metadata bind to peer identity, device identity, relationship
  context, group/session context, or some combination?
