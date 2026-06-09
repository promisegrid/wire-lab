# TODO-pazif: POC13 Superset Repair

## Decision Intent Log

ID: DI-sinur
Date: 2026-06-08 17:50:52
Status: active
Decision: Refactor POC13 in place so it is a strict superset of POC11 and POC12, and make future POCs superset-by-default unless a later scoped DI explicitly authorizes a non-superset exception.
Intent: POC13 regressed by preserving CAS/compute evidence while dropping POC11/POC12 architecture and acceptance lessons. Future POCs must retain inherited app/kernel boundaries, pCID routing, local trust semantics, relationship journals, monitor/analyzer gates, and proven workflows unless the repo records a deliberate exception before implementation.
Constraints: Use POC12 as the POC13 repair base; keep the existing `implementations/poc13-cas-compute-functions/` path; require both deterministic analyzer gates and observer-monitor minimum scores; preserve the single top-level `promise` action; keep Promise Theory vocabulary and local-trust semantics.
Affects: AGENTS.md; implementations/poc13-cas-compute-functions/; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md.

ID: DI-punib
Date: 2026-06-08 18:18:58
Status: active
Decision: Fix POC13 clean-run regressions by treating live-agent pCID mismatches as relationship-level promises, preventing negative ACK verdicts from increasing local trust, making bad-proof ACK handling explicit, and discovering the actual run evidence directory shape.
Intent: The first clean Docker run showed that the repaired POC13 architecture worked, but the run failed acceptance because live LLM decisions could select a protocol pCID without providing that pCID's required payload semantics, malformed/non-commitment ACK fields could still increase trust, bad-proof evidence was logged ambiguously, and the clean-run script expected an obsolete `/run` subdirectory. These are implementation hygiene defects, not new protocol features.
Constraints: Preserve the single top-level `promise` action, keep scripted CAS/compute/shipping protocol tests exact, reframe only autonomous live-agent free-form promises, do not turn receiver non-commitment into a broken promise, and keep all run inputs local/ignored.
Affects: implementations/poc13-cas-compute-functions/scripts/run-clean.sh; implementations/poc13-cas-compute-functions/decision/decision.go; implementations/poc13-cas-compute-functions/runtime/node.go; implementations/poc13-cas-compute-functions/runtime/node_test.go.

## Tasks

- [x] pazif.1 Record the POC superset rule in `AGENTS.md` and this DI.
- [x] pazif.2 Rebase POC13 implementation shape on POC12's separate app/kernel process architecture.
- [x] pazif.3 Preserve POC11 structured LLM decisions, local relationship ledgers, economics, direct TCP relationship transitions, and observer monitor.
- [x] pazif.4 Preserve POC12 multi-pCID app routing, shipping/device workflow, printer-port capability tokens, promise journals, non-commitment restraint, duplicate checkpoints, and resource/trust separation.
- [x] pazif.5 Port current POC13 CAS storage, CID compute, replica, token lifecycle, cache, verifier disagreement, bad-proof, and key-rotation evidence into the repaired architecture.
- [x] pazif.6 Expand analyzer gates so POC13 fails if any inherited POC11/POC12 or current POC13 evidence disappears.
- [x] pazif.7 Update README, run narrative, and resource docs to state POC13's superset status and future-POC rule.
- [x] pazif.8 Validate with Go tests, vet, errcheck, and Compose syntax; defer clean Docker/provider run to the operator because the run consumes live provider credentials.
