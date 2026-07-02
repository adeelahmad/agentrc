---
type: validate
story: S2-03
---
# S2-03 Validation — CLI `--backend` rename (T21)

## Pre-flight

- [ ] Grep-located EVERY `--substrate` occurrence before editing (§0.9, M-003):
      `cmd/agentrc/run.go:19`, `cli.md:86/103/120`, `docs/quickstart.md`, `tooling/README.md`.
- [ ] Confirmed the rename is CLI-flag-only — no edit touches the `substrate.*` POLICY
      namespace or "substrate-neutral" prose in `spec/` (§0.2).
- [ ] RED test written first and FAILS before implementation (real TDD).
- [ ] Version unchanged: single `draft.5` (§0.1).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | flag present | `go run ./cmd/agentrc run --help \| grep -- '--backend'` | match (§V.9) |
| 2 | old flag gone (CLI/docs) | `rg -l -- '--substrate' cmd/ docs/ cli.md 2>/dev/null \| wc -l` | `0` (§V.8) |
| 3 | namespace intact | `rg -c 'POLICY substrate\.' spec/index.md` | `>= 1` (§V.8) |
| 4 | backend values | `go run ./cmd/agentrc run x --backend bogus 2>&1` | non-zero exit + clear error |
| 5 | scoped build | `go build ./cmd/agentrc/...` | exit 0 (local); full `go build ./...` in CI |
| 6 | scoped vet | `go vet ./cmd/agentrc/...` | exit 0 |
| 7 | unit tests | `go test ./cmd/agentrc/...` | PASS (`-race` full in CI) |
| 8 | mod tidy | `go mod tidy && git diff --exit-code go.mod go.sum` | no diff |
| 9 | GCP/Compose recorded | `grep -c 'GCP' cli.md` and `grep -c 'Compose' cli.md` | `>= 1` each |
| 10 | version unchanged | `grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u \| wc -l` | `1` (`draft.5`) |
</content>
