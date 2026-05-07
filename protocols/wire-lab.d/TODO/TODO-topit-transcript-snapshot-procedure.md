# TODO-topit: Switch session logging from per-turn files to transcript-snapshot procedure

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-topit`. No prior
integer or timestamp alias.

For convenience in cross-references during the rollout window, the
master cross-list (`protocols/wire-lab.d/TODO/TODO.md`) carries this
file under the integer alias `TODO-33`.

## Driving observation

The per-turn discipline locked into AGENTS-ppx Session Logging requires
the bot to write `sessions/<session-id>/<NNN>-turn.md` with verbatim
prompt and response, commit, and push at the end of every turn.
Three reinforcing problems have surfaced (see DR-010-20260507-150000):

1. The bot does not have reliable in-process access to its own
   verbatim response bytes. Per-turn writes have repeatedly drifted
   into paraphrase or truncation, and the loose-`turns/` drift
   remediated by the 2026-05-07 anti-pattern callout was downstream
   of the same root cause: bots reconstructing turn contents from
   memory rather than from a captured transcript stream.
2. The harness already retains a complete `conversation.md` for the
   active session, refreshed each turn. That file IS a verbatim
   transcript, with no reconstruction step required.
3. Per-turn commits create high commit volume on the orphan branch
   (`wire-lab` on session-logs) that is dominated by reconstruction
   churn rather than primary signal, while still failing to
   guarantee fidelity.

The transcript-snapshot procedure substitutes a periodic copy of
`conversation.md` for the per-turn discipline, with a small
`TURNOVER-*.md` file as the live resume point and an optional
generated index over the snapshot.

## Locked decisions (from DR-010-20260507-150000)

1. **Canonical artifact:** the full session transcript is the file the
   harness writes as `conversation.md` for the active session.
   Periodic commits of that file to the `wire-lab` orphan branch on
   `stevegt/session-logs` are the durable record.

2. **Snapshot cadence:** the bot copies and commits `conversation.md`
   to `sessions/<session-id>/conversation.md` at session checkpoints
   (turnover, redact-last, end-of-day, or any explicit Steve request)
   and at minimum once per session. Commit message:
   `wire-lab session <session-id>: transcript snapshot @ turn <NNN>`.
   The snapshot is a full file replace, not a diff; the orphan branch
   accumulates the file's history through normal git history.

3. **Live resume point:** a `TURNOVER-YYYYMMDD-HHMMTZ[-suffix].md` file
   at the worktree root captures the in-flight context the next bot
   needs to resume. The TURNOVER file is short (one screen of text)
   and is append-only at the file-system level: existing TURNOVER
   files are NEVER edited or moved. Each new session writes a new
   TURNOVER file rather than updating an old one.

4. **Per-turn files retired by default.** New sessions do NOT write
   `sessions/<session-id>/<NNN>-turn.md`. The audit-only escape
   hatch `--allow-per-turn` on `tools/check-layout.sh` exists for the
   rare case where a specific audit needs reconstructed per-turn
   files; it is opt-in and never default. Existing
   `sessions/ea135ce8/NNN-turn.md` files and the frozen
   `turns/turn-NNN.md` files remain untouched as historical record.

5. **Optional generated index.** A consumer may build a per-turn
   index over a snapshot by parsing turn boundaries inside
   `conversation.md` and emitting an index file. The index is
   generated, not committed: the snapshot itself is canonical and any
   index is reproducible from it. If a future audit pins an index
   into history, that pin lives on a separate non-canonical branch.

6. **Identity gate unchanged.** The Identity gate from the prior
   Session Logging block still applies: snapshot writes only happen
   when git config user.email is the bot's identity (or Steve's
   direct identity for testing). Collaborators who clone wire-lab
   without bot credentials see the procedure documented but never
   execute it.

7. **Force-push policy unchanged.** Force-push to
   `private wire-lab` remains permitted ONLY for the redact-last
   operation, which now operates on the snapshot file rather than on
   per-turn files (`git rm` the snapshot is not the redact-last
   operation; redact-last edits the snapshot in place to remove the
   offending span and force-pushes the resulting commit, with the
   message `wire-lab session <session-id> redact snapshot`).

## Decision Intent Log

### DI-033-20260507-150000

- ID: DI-033-20260507-150000
- Date: 2026-05-07
- Status: locked
- Decision: Replace the per-turn session-logging discipline with a
  transcript-snapshot procedure. The bot copies and commits the
  harness-produced `conversation.md` to
  `sessions/<session-id>/conversation.md` at session checkpoints and
  at minimum once per session. A small `TURNOVER-*.md` at the
  worktree root carries the live resume point for the next bot. The
  bot does NOT reconstruct or commit per-turn `<NNN>-turn.md` files
  by default.
- Intent: drive fidelity by writing the file the harness already
  produces verbatim, rather than asking the bot to reproduce it from
  memory; cut commit volume on the orphan branch; keep the resume-
  on-cold-start workflow that TURNOVER files have proven out.
- Constraints:
    - Pre-existing `sessions/ea135ce8/NNN-turn.md` files and
      `turns/turn-NNN.md` files remain in place; this DI does not
      back-edit history.
    - The Identity gate from the prior Session Logging architecture
      (locked in the 2026-05-07 anti-pattern callout merge `93802eb`)
      still applies.
    - Redact-last continues to be the only force-push case on
      `private wire-lab`; it now operates on the snapshot file rather
      than on per-turn files.
    - `tools/check-layout.sh` in `stevegt/session-logs` is updated to
      reject NEW `sessions/<session-id>/<NNN>-turn.md` files by
      default and to accept the snapshot path
      `sessions/<session-id>/conversation.md`. An opt-in
      `--allow-per-turn` flag preserves the audit-only escape hatch.
- Affects:
    - `AGENTS-ppx.md` Session Logging section.
    - `AGENTS-codex.md` (light edit; codex is the reviewer side and
      needs to know the new layout when reviewing snapshot commits).
    - `bin/ppx-bootstrap.sh` (comment update only; the script's
      step list is unchanged).
    - `stevegt/session-logs` `wire-lab` orphan: `README.md`
      (Layout / Conventions / Anti-pattern blocks) and
      `tools/check-layout.sh` (new default-reject rule plus the
      opt-in flag) and a new
      `TURNOVER-20260507-1500PDT-transcript-snapshot-cutover.md`
      cutover note.
- Author: stevegt@t7a.org (Steve Traugott)

## Cross-references

- DR-010-20260507-150000-transcript-snapshot.md — the DR backing
  this DI.
- DR-009-20260430-204108-group-transport-envelope.md — for the
  prior-art shape of the wire-lab.d DR/DI layout.
- `wire-lab-logs/TURNOVER-20260507-2300PDT.md` and
  `wire-lab-logs/TURNOVER-20260508-0030PDT.md` — the loose-`turns/`
  remediation trail that motivated the cutover.
- `wire-lab-logs/sessions/ea135ce8/NNN-turn.md` — historical per-turn
  files preserved as-is.

## Status

In progress. Twig: `ppx/topit-transcript-snapshot-procedure`.
