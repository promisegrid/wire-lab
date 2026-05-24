# Question

Can PromiseGrid use a three-slot outer envelope
`grid([pCID, payload, signature])` where the third slot signs canonical
`[pCID, payload]` bytes and the protocol specification named by `pCID` defines
the proof family and verification rules, without needing a separate outer
`sig_pcid` selector? Source: `DI-kukuk`; `DI-pozom`.
