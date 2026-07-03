---
type: standards
sprint: 2
---

# Sprint 2 Standards — the active rule digest and gate matrix

Read-only digest produced by the standards activity (The Lawkeeper) for Sprint 2
("Features", T15–T27). Every rule below cites a real source in this repo or in the
Sprint 2 work order; nothing is invented. The gate matrix at the bottom is what the
GREEN gate (per-change, local) and the FINAL gate (pre-commit / pre-deploy, §V incl.
new checks 8–9) run. Baseline carried forward from
`docs/agents/sprint1/standards.md`; this refresh adds the Go-heavy CLI backend work
(T21–T24) and the draft.6 spec/markdown work (T15–T20, T25).

## Detected stack

**Go module for the agentrc / arc CLI** (built at repo root, NOT under `tooling/`):
- `go.mod` at repo root — module `github.com/adeelahmad/agentrc`, `go 1.25.9`.
  Existing deps are BuildKit-centric: `github.com/moby/buildkit v0.31.1`,
  `github.com/google/go-containerregistry v0.21.6`, `github.com/spf13/cobra v1.10.2`,
  `github.com/opencontainers/image-spec v1.1.1` (+ the containerd/otel indirect tree).
- CLI entrypoint is `./cmd/agentrc` (verified: `cmd/agentrc/main.go` registers
  init/lint/lock/build/inspect/push/pull/sign/verify/run). Internal packages:
  `internal/agentfile` (model, extract, validate, labels, security, resources) and
  `internal/llb`.
- **Stubs today (confirmed):** `run` returns "not implemented yet" (`cmd/agentrc/run.go`
  RunE errors out; it already declares `--isolation` and `--substrate` string flags per
  cli.md); `sign` and `verify` are likewise stubs (`cmd/agentrc/sign.go`,
  `verify.go` — "not yet implemented; see cli.md"). Sprint 2 makes `run` real
  (T21–T24); `sign`/`verify` stay `planned` (T25).
- `cli.md` (repo root) documents the run flag surface being renamed:
  `agentrc run <ref> [--isolation local|container|microvm] [--substrate <driver>]`
  (cli.md lines 86, 120, 129). T21 renames `--substrate` → `--backend`.

**New deps Sprint 2 may add:** T23 (bedrock) needs an AWS SDK path
(work-order lines 95–101); T24 (kubernetes) needs either a k8s client or plain-YAML
manifest emission (work-order lines 102–105). **Any added dependency MUST keep
`go build`, `go vet`, and `go test` green and `go.sum` tidy** (`go mod tidy` clean;
the `go` job in `.github/workflows/ci.yml` runs `go build ./...` which will fail on an
untidy module). Prefer plain-YAML/struct emission over heavy clients to avoid
enlarging the already-large dependency tree (see ENV RISK below).

**Jekyll / Ruby static site** (GitHub Pages, deployed to https://agentrc.ai):
- `_config.yml` (`markdown: kramdown`, `permalink: pretty`, `jekyll-seo-tag` +
  `jekyll-sitemap`, `url: https://agentrc.ai`), `Gemfile` (`gem "github-pages"`),
  `_layouts/`, `_includes/head.html` (canonical link emitter), `CNAME` (`agentrc.ai`).
- Spec, docs, examples, profiles carry the markdown Sprint-2 work: examples
  T15–T16 (`examples/Agentfile.hooked`, `examples/Agentfile.delegator` + index card),
  spec draft.6 T17–T20 (§8.5–§8.7, §14.2, CHANGELOG), CLI docs T25.

**GitHub Actions workflows** (`.github/workflows/`):
- `ci.yml` — job `build` (Jekyll build via `bundle exec jekyll build --trace` +
  html-proofer link check) and job `go` (`go build ./...`, `go vet ./...`,
  `go test -race ./...`).
- `pages.yml` — GitHub Pages build + deploy (`actions/jekyll-build-pages@v1`,
  `deploy-pages@v4`) on push to `master`/`main`.
- `release.yml` — release workflow.

Deploy path: GitHub Pages via `pages.yml` on push to `master` (T27 "deploy" = merge
the sprint branch via PR).

## Active rules

- `.github/workflows/ci.yml` `go` job (lines 47–66) + `docs/agents/sprint2/work-order.md` §V — **Go build/vet/test must pass:** `go build ./... && go vet ./... && go test -race ./...` is the authoritative Go gate. **ENV constraint:** local disk is near-full and `go build ./...` links the buildkit-heavy frontend, which has hit ENOSPC locally (work-order lines 20–22, "ENV RISK"). Therefore the *practical local gate* for T21–T24 is per-package — `go test ./cmd/agentrc/... ./internal/...` and `go vet ./cmd/agentrc/...` — with full `go build ./...` / `go test -race ./...` delegated to CI. Do not treat a local inability to run full `go build ./...` as a failure; CI is the source of truth.
- `docs/agents/sprint2/work-order.md` §P ("real-TDD execution", line 22) + T23 (lines 95–101) — **Real TDD for T21–T24:** write a failing test first (RED), implement to GREEN, refactor. Backend translators (local/bedrock/kubernetes) are pure labels→config mapping, unit-testable without a live cloud; each of T23's 13/13 mappings and each fail-closed path needs a table-driven test. Prefer `--dry-run` output as the assertion surface (T21 line 88: "`--dry-run` prints the translated config and exits").
- `docs/agents/sprint2/work-order.md` T23/T24 fail-closed clauses (lines 99–105) + §0.4 (line 36) — **Fail-closed is REQUIRED for backend translators:** bedrock (T23) MUST fail closed on missing `roleArn`, an unenforceable `agent.auth.mode=jwt`, or code-mode without a resolvable language (no endpoint emitted); kubernetes (T24) emits a deny-by-default NetworkPolicy from `POLICY network dns:*`. §8.6 `agent.auth.*` (T18 lines 68–71) is fail-closed generic authZ — a platform that cannot enforce a requested `jwt` authorizer MUST NOT expose the invocation endpoint; §8.7 code-mode language (T19 lines 72–74) fails closed if not resolvable.
- `docs/agents/sprint2/work-order.md` §0.3 (lines 32–34) — **Exactly three new POLICY namespaces, no more:** T17 `substrate.<platform>.*`, T18 `agent.auth.*`, T19 `substrate.runtime.language`. No new keywords ever (the four stay IDENTITY, CAPABILITY, SOP, POLICY); no POLICY entries beyond these three.
- `docs/agents/sprint2/work-order.md` §0.2 (lines 29–31) + §V check 8 (lines 119–120) — **The `--backend` rename must not touch `substrate.*` or "substrate-neutral":** the `--substrate <driver>` CLI flag + "driver" wording are renamed to `--backend local|bedrock|kubernetes` in T21, but the POLICY namespace `substrate.*` and the concept "substrate-neutral" are NEVER renamed.
- `docs/agents/sprint2/work-order.md` §0.8 (lines 42–43) + T26 (lines 111–113) — **Positioning lines, VERBATIM:** wherever backend translators are described (T22, T25 above the CLI docs table) use exactly "Reference translators — a proof of concept until platforms read `ai.agentrc.*` labels natively. Not production runners."; the T26 demo narrative is also verbatim: "Same artifact, same labels, three substrates. The translators are the proof of concept; the labels are the standard."
- `docs/agents/sprint2/work-order.md` §0.1 (lines 27–28) + §V check 3 (line 124) — **draft.6 bump only in T20:** version stays `0.1.0-draft.5` through Phase 3; the sitewide bump to `0.1.0-draft.6` happens in exactly one commit in T20; the syntax line stays `# syntax=agentrc.agentfile/v0.1` throughout. Exactly one `draft.N` value in the tree — `draft.5` before T20, `draft.6` after.
- `docs/agents/sprint2/work-order.md` §0.7 (lines 39–40, 59) — **Every example passes `arc lint`:** each example file must pass `go run ./cmd/agentrc lint <file>` from repo root (CLI is at repo-root `./cmd/agentrc`, NOT `tooling/` — Sprint 1 correction, still valid). Applies to the new T15/T16 files and all existing `examples/Agentfile.*`.
- `docs/agents/sprint2/work-order.md` §0.6 (line 38) — **One canonical hello:** the hello is byte-identical wherever rendered inline and identical to `examples/Agentfile.minimal`; propagate, never fork.
- `docs/agents/sprint2/work-order.md` §0.4 (line 36) — **POLICY lines are requests; Cedar platform-side only; secrets deferred:** no inline Cedar, no secret keyword. `agent.auth.*` is generic authZ config (fail-closed), NOT a secret.
- `docs/agents/sprint2/work-order.md` §0.9 (lines 44–45) — **Grep-locate-first; rebuild site + pass §V before commit:** locate every target by grep before editing; no local Jekyll in this env, so build-dependent checks (site build, link check, one-canonical) delegate to CI + live, as in Sprint 1.
- `docs/agents/sprint2/work-order.md` (with `docs/agents/sprint1/standards.md` as the carried baseline) — no project `CLAUDE.md` exists at repo root; the work order is the authoritative rule source for this sprint.

## Gate matrix

All commands are literal and runnable from the repo root
(`/Users/adeelahmad/work/agentrc`). The **Where** column marks whether the gate is the
practical *local* gate (given the ENOSPC disk constraint) or CI's job.

| Gate | Command | Applies to | Where |
|------|---------|------------|-------|
| go-build | `go build ./cmd/agentrc/...` | Go source under `cmd/agentrc/` (T21–T25); scoped to the CLI pkg to avoid the ENOSPC frontend link. Full `go build ./...` is CI's job (`go` job, `ci.yml`). | local (scoped) + CI (full) |
| go-vet | `go vet ./cmd/agentrc/...` | Go source under `cmd/agentrc/`; full `go vet ./...` in CI. | local (scoped) + CI (full) |
| go-test | `go test ./cmd/agentrc/... ./internal/...` | CLI + translator unit tests (T21–T24 TDD). Full `go test -race ./...` in CI (`ci.yml` line 66). | local (scoped) + CI (-race, full) |
| go-mod-tidy | `go mod tidy && git diff --exit-code go.mod go.sum` | any change adding a dep (T23 AWS SDK, T24 k8s/YAML); keeps `go.sum` tidy so CI `go build ./...` stays green. | local + CI |
| example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" || exit 1; done` | every example Agentfile incl. new T15/T16; §0.7 / §V check 5. | local + CI |
| terminology-split | `[ "$(rg -l -- '--substrate' cmd/ docs/ cli.md 2>/dev/null \| wc -l)" -eq 0 ] && rg -q 'POLICY substrate\.' spec/index.md` | §V.8 (work-order lines 119–120): `--substrate` gone from CLI/docs (→0), `POLICY substrate.` in spec intact. | local + CI |
| backend-dryrun-bedrock | `go run ./cmd/agentrc run <ref> --backend bedrock --dry-run \| python3 -m json.tool` | §V.9 (lines 121–123): bedrock dry-run emits valid JSON. | local + CI |
| backend-dryrun-k8s | `go run ./cmd/agentrc run <ref> --backend kubernetes --dry-run \| kubeconform -` (or yaml-parse / `python3 -c 'import sys,yaml;list(yaml.safe_load_all(sys.stdin))'` if kubeconform absent) | §V.9 (lines 121–123): kubernetes dry-run yaml-parses / kubeconform-validates. | local (yaml-parse) + CI |
| backend-flag-surface | `go run ./cmd/agentrc run --help \| grep -- '--backend'` | §V.9 (line 121): `--backend` present on `run --help`. | local + CI |
| version-coherence | `[ "$(grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u \| wc -l)" -eq 1 ]` — value is `draft.5` before T20, `draft.6` after | §V check 3 / §0.1 (lines 27–28, 124): single `draft.N`; `draft.6` only after T20. | local + CI |
| hello-canonical | `diff <(sed -n '/# syntax=/,$p' examples/Agentfile.minimal)` against each inline render | §0.6 (line 38): one byte-identical hello. | local + CI |
| site-build | `bundle exec jekyll build --trace` | any page/spec/example change; §0.9. No local Jekyll in this env → CI's `build` job (`ci.yml` lines 33–34). | CI |
| internal-link-check | `htmlproofer ./_site --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https` | any change touching pages/links; §V check 7 (`ci.yml` lines 38–45). | CI |
| one-canonical-per-page | `for p in $(find _site -name '*.html'); do [ "$(grep -c 'rel=\"canonical\"' "$p")" -eq 1 ] \|\| exit 1; done` | every built `_site/**/*.html`; §V check 6. | CI / live |

Notes on the ENV disk constraint (work-order lines 20–22): the scoped `go build`/`go vet`/
`go test` commands above are the *practical local GREEN gate* — they exercise only
`./cmd/agentrc/...` and `./internal/...` and skip linking the buildkit frontend that
triggers ENOSPC. The full `go build ./...` and `go test -race ./...` are authoritative
in CI. Site/link/canonical gates have no local Jekyll and run in CI + live per Sprint 1
practice.
