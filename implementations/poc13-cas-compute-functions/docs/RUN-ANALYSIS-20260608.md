# POC13 Superset Repair Analysis — 2026-06-08

This document records the implementation-level analysis for the POC13 superset
repair. It is not a fresh Docker run result. Source: `DI-sinur`.

## What Changed

- POC13 now uses the POC12 app/kernel shape: one local kernel process per
  container and separate app entrypoints for each local role.
- The original POC13 CAS/compute pressure is ported into the repaired runtime as
  pCID-routed promise exchanges instead of a single monolithic supervisor.
- `cas_storage_v1`, `cid_compute_v1`, and `evidence_report_v1` are registered
  pCIDs alongside relationship, shipping, printer-port, and kernel receive
  pCIDs.
- The analyzer now fails if inherited POC11/POC12 evidence or POC13
  storage/compute evidence disappears.
- The repo-level POC rule is now superset-by-default unless a future scoped DI
  explicitly declares and justifies a non-superset exception.

## Validation Status

The local Go tests for the POC13 module pass after the repair. A clean Docker
run remains the operator-level validation because it requires Docker runtime
state and live-provider credentials.

Recommended command:

```sh
cd implementations/poc13-cas-compute-functions
cp config.example.json config.json
printf '%s' "$OPENAI_API_KEY" > openai_api_key.txt
chmod 600 openai_api_key.txt
scripts/run-clean.sh
```
