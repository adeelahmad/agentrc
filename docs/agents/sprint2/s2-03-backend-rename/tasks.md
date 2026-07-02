---
type: tasks
story: S2-03
---
# S2-03 Tasks — CLI `--backend` rename (T21)

REAL TDD (RED→GREEN→refactor). CLI-flag-only rename; NEVER touch the `substrate.*` POLICY
namespace or "substrate-neutral" concept (§0.2). Grep-first every `--substrate` (M-003).
Version stays `draft.5` (§0.1). Scope Go gates to `./cmd/agentrc/...` (ENOSPC; CI runs full).

## T21 — Flag rename & backend surface [P0* demo]

### Go: `cmd/agentrc/run.go`
- Remove the `--substrate` string flag (`run.go:19`); add `--backend` accepting
  `local|bedrock|kubernetes` (default `local`).
- Keep `--isolation` (`local|container|microvm`) but scope it to `--backend local`.
- Add per-backend flags: bedrock `[--region] [--profile] [--dry-run]`; kubernetes
  `[--kubeconfig] [--namespace] [--dry-run]`. `--dry-run` prints the translated config and
  exits (S2-04 fills the translators; here the flag/dispatch surface + validation land).
- Update the command `Short`/help so it no longer says "not yet implemented" for the flag
  surface (M-002); dispatch to a `backend` seam (a `translate(backend, labels) (string,error)`
  function stub that S2-04 implements per backend).
- Reject an unknown `--backend` value with a clear error.

### Docs (grep-located `--substrate` surfaces)
- `cli.md:86` — `agentrc run <ref> [--isolation …] [--substrate <driver>]` →
  `agentrc run <ref> --backend local|bedrock|kubernetes [per-backend flags]`.
- `cli.md:103` — the run table row: replace `--substrate` prose with `--backend`.
- `cli.md:120-129` — the `--isolation`/`--substrate` explanation + the `arc run …
  --isolation microvm` example: reword to `--backend`, keep the "Agentfile never names a
  substrate" sentence and the `POLICY substrate.*` request wording INTACT (§0.2).
- `docs/quickstart.md` step 5 and `tooling/README.md` — any `--substrate` mention → `--backend`.
- Record in CLI docs: **GCP dropped** (Agent Runtime Python-only managed; GKE via kubernetes
  backend); **Docker Compose dropped** (no `network.*` egress enforcement without a bespoke
  sidecar).

Verify: `rg -l -- '--substrate' cmd/ docs/ cli.md` → 0; `rg 'POLICY substrate\.' spec/index.md`
intact; `go run ./cmd/agentrc run --help | grep -- '--backend'` present; scoped
`go build/vet/test ./cmd/agentrc/...` green.
</content>
