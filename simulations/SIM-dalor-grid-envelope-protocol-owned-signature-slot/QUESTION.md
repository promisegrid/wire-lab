# Question

Can PromiseGrid use a three-slot outer envelope
`grid([pcid, payload, signature])` where the third slot signs canonical
`[pcid, payload]` bytes and the payload `pcid` itself defines the proof family
and verification rules, without needing a separate outer `sig_pcid` selector?
Source: `DI-kukuk`.

