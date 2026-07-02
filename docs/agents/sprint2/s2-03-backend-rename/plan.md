---
type: plan
story: S2-03
scope: "tests only"
---
# S2-03 Test Plan — CLI `--backend` rename (RED)

Tests only. REAL Go test handles under `cmd/agentrc/`. Every bullet FAILS now because
`run.go` still declares `--substrate` and errors "not implemented". Scope: `./cmd/agentrc/...`
(ENOSPC; `-race`/full build in CI).

## T21 — flag rename & surface

- [ ] `cmd/agentrc/run_test.go::TestBackendFlagReplacesSubstrate` — given `run --help`, assert
      output contains `--backend` and does NOT contain `--substrate`. FAILS now: `run.go:19`
      declares `--substrate`, no `--backend`.
- [ ] `cmd/agentrc/run_test.go::TestBackendDefaultsToLocal` — given `run <ref>` with no
      `--backend`, assert the resolved backend is `local`. FAILS now: no `--backend` flag exists.
- [ ] `cmd/agentrc/run_test.go::TestBackendRejectsUnknownValue` — given `run <ref> --backend
      bogus`, assert a non-nil error naming valid values. FAILS now: no validation exists.
- [ ] `cmd/agentrc/run_test.go::TestIsolationScopedToLocalBackend` — given `--backend bedrock
      --isolation microvm`, assert an error or that isolation is ignored for non-local. FAILS
      now: `--isolation` is global and unwired.
- [ ] `cmd/agentrc/run_test.go::TestPerBackendFlagsParse` — assert `--region`/`--profile`
      (bedrock) and `--kubeconfig`/`--namespace` (kubernetes) and `--dry-run` parse without
      error. FAILS now: flags absent.
- [ ] `cmd/agentrc/run_test.go::TestDryRunExitsAfterPrint` — given `--backend bedrock --dry-run`
      with a valid ref/labels, assert the command prints and exits 0 without invoking a runtime.
      FAILS now: `run` returns the hard "not implemented" error.

## Suite guards (markdown/repo-level, in the shared harness)

- [ ] `scripts/verify-sprint2.sh::v8_terminology_split` — assert `rg -l -- '--substrate' cmd/
      docs/ cli.md` count == 0 AND `rg 'POLICY substrate\.' spec/index.md` intact (§V.8). FAILS
      now: `--substrate` present in `run.go` + `cli.md`.
- [ ] `scripts/verify-sprint2.sh::v3_version_draft5_guard` — assert single `draft.N` ==
      `draft.5` (no accidental bump in this story; §0.1). RED now: green; regression guard.
</content>
