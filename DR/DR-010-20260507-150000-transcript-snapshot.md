# DR-010 - Switch session logging to transcript-snapshot procedure

DR-ID: DR-010-20260507-150000
Date: 2026-05-07 15:00:00
Asked by: stevegt+ppx@t7a.org (stevegt-via-perplexity)
State: decided
Question: Should the wire-lab Perplexity Computer session logging continue to require per-turn `sessions/<session-id>/<NNN>-turn.md` files, or should it switch to periodic snapshots of the harness-produced `conversation.md` plus a small live-resume `TURNOVER-*.md`?
Why this blocks progress: Two reinforcing failures have shown up under the per-turn discipline. First, the bot does not have reliable in-process access to its own verbatim response bytes, so per-turn writes drift into paraphrase or truncation; the loose-`turns/` drift remediated on `93802eb` was a downstream symptom of the same root cause. Second, the orphan branch on `stevegt/session-logs` accumulates one commit per turn dominated by reconstruction churn rather than primary signal. The harness already retains a verbatim `conversation.md` for the active session, so the simplest and most fidelity-preserving fix is to copy that file periodically rather than reconstruct it turn by turn.
Affects:
- `AGENTS-ppx.md` Session Logging section (replaces "Per-turn discipline" with snapshot discipline).
- `AGENTS-codex.md` (codex needs to know the snapshot layout when reviewing branches that touch session logging).
- `bin/ppx-bootstrap.sh` (comment-only update; script step list unchanged).
- `stevegt/session-logs` `wire-lab` orphan branch:
    - `README.md` — Layout, Conventions, Anti-pattern, Why-this-exists blocks.
    - `tools/check-layout.sh` — new default-reject rule for new `NNN-turn.md` files plus an opt-in `--allow-per-turn` audit-only flag, accept snapshot path `sessions/<session-id>/conversation.md`.
    - New `TURNOVER-20260507-1500PDT-transcript-snapshot-cutover.md` cutover note.
- `protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md` — parent TODO.
- `protocols/wire-lab.d/TODO/TODO.md` — master cross-list gains the TODO-topit row.
Unblocks: future sessions on this and successor sandboxes; clean session-logs orphan branch with low-volume, full-fidelity history; future audit tooling that prefers a single-file transcript over per-turn reconstruction.
Waiting on: DI-033-20260507-150000 (in `protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md`).

## Candidate alternatives considered

- **Alt-A: keep per-turn discipline.** Continue requiring
  `sessions/<session-id>/<NNN>-turn.md` per turn. Strengthen guards
  further. Rejected because the underlying fidelity problem (the bot
  reconstructing its own response from memory) cannot be guarded
  away; only the layout drift it produces.
- **Alt-B: snapshot-only.** Periodically copy `conversation.md`,
  drop both the per-turn files and the TURNOVER live-resume
  mechanism. Rejected because the cold-start workflow has proven its
  value; new bots need a short read at known location to resume
  in-flight context, and `conversation.md` alone does not surface
  the open threads cleanly.
- **Alt-C: snapshot plus TURNOVER, retire per-turn by default.**
  Periodic `conversation.md` snapshots are the canonical record;
  short `TURNOVER-*.md` files at the worktree root carry the live
  resume point; per-turn files are not written by default but remain
  available behind an opt-in audit flag for the rare case where
  reconstruction is justified. Selected.
- **Alt-D: commit one index file per turn, generated from the
  transcript.** Creates the same commit volume as Alt-A while
  decoupling from prompt-byte fidelity. Rejected because the index
  is reproducible from the snapshot at any time and pinning it into
  history adds noise without adding signal.

## Decision

Adopt Alt-C. The full-fidelity record is `sessions/<session-id>/conversation.md`, refreshed by the bot copying the harness-produced file at session checkpoints (turnover, redact-last, end-of-day, or any explicit Steve request) and at minimum once per session. The live resume point is a `TURNOVER-YYYYMMDD-HHMMTZ[-suffix].md` file at the orphan-branch worktree root. Existing `TURNOVER-*.md` files are append-only at the file-system level — they are never edited or moved. New sessions write new TURNOVER files. Per-turn `NNN-turn.md` files are not written by default; the audit-only escape hatch is the `--allow-per-turn` flag on `tools/check-layout.sh`.

Existing per-turn artifacts on the `wire-lab` orphan branch — both `sessions/ea135ce8/NNN-turn.md` and the frozen `turns/turn-NNN.md` files — remain in place as historical record. This DR does not back-edit history.

The Identity gate (snapshot writes only fire when git config user.email is the bot's identity, or Steve's direct identity for testing) and the redact-last force-push escape hatch carry forward unchanged from the prior Session Logging block. Redact-last now operates on the snapshot file: edit the snapshot in place to remove the offending span, commit, and `git push --force-with-lease private wire-lab` with the message `wire-lab session <session-id> redact snapshot`.

## Linked DI

- `DI-033-20260507-150000` in `protocols/wire-lab.d/TODO/TODO-topit-transcript-snapshot-procedure.md`

## Related commits

- `93802eb` — wire-lab `ppx/main` tip carrying the loose-`turns/`
  anti-pattern callout that this DR supersedes for new sessions.
- `92d8a8c` — `wire-lab` orphan branch tip carrying the layout guard
  this DR extends with the `--allow-per-turn` flag.

## Last updated

2026-05-07 15:00:00 UTC. Filed with `State: decided` because the
decision was reached in chat before the file was written, per
AGENTS-ppx.md "DRs that Steve answered in chat before the file was
written may be created with State: decided directly."
