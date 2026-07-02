---
type: sprint-plan
sprint: 2
stage: "2"
---
# Sprint 2 Plan — "Features" (T15–T27), Stage 2

Wave/dependency plan for the seven Stage-2 stories (S2-01…S2-07). Grounded in
`stories.md` (dependency graph), `standards.md` (gate matrix, ENOSPC→CI note), and
`work-order.md` §0/§V. This file schedules execution; per-story RED test handles live
in each `s2-NN-*/plan.md`.

## Sprint goal

Ship two lint-clean examples, spec **draft.6** (exactly three new POLICY namespaces),
the `--substrate` → `--backend` rename with three fail-closed reference translators,
a one-agent/three-backend demo, and a released + live-verified site — without adding a
keyword, renaming the `substrate.*` namespace / "substrate-neutral" concept, or bumping
the version before the single T20 commit.

## Critical path

```
S2-02 (T17–T19 spec sections)  →  S2-03 (T21 --backend rename)  →  S2-04 (T22–T24 translators)
   →  S2-02 T20 (single sitewide draft.5→draft.6 commit + T8 A)  →  S2-07 (T27 release + live verify)
```

Rationale: the translators (S2-04) read the §8.7/§8.8/§8.9 namespaces, so the spec
sections gate them; the rename (S2-03) is the surface the translators hang off; the
version bump (T20) is a single late commit AFTER the namespaces are consumed; release
(S2-07) is last. S2-01 (examples) and S2-05/S2-06 (docs/demo) are near-critical but not
on the longest chain.

## Waves

| Wave | Stories / tasks | Runs in parallel? | Gate to advance |
|------|-----------------|-------------------|-----------------|
| 1 | **S2-01** (T15, T16) + **S2-02 T17–T19** (spec sections, authored at `draft.5`) | Yes — different files (`examples/*` vs `spec/index.md` §8.7–8.9) | `example-lint` green for new files; §8.7/§8.8/§8.9 present; `version-coherence` still one `draft.5`; `terminology-split` unaffected |
| 2 | **S2-03** (T21 `--backend` rename) | Single-threaded (touches `cmd/agentrc/run.go`, `cli.md`) | `TestBackendFlagReplacesSubstrate` green; `terminology-split` (§V.8) passes; `backend-flag-surface` shows `--backend`; scoped `go build/vet/test ./cmd/agentrc/...` |
| 3 | **S2-04** (T22 local, T23 bedrock, T24 kubernetes) | Translators are independent packages/files — may fan out, but all in `cmd/agentrc/` so serialize commits | all fail-closed RED tests green; `backend-dryrun-bedrock` (json.tool) + `backend-dryrun-k8s` (yaml-parse) pass; single k8s format; §0.8 line verbatim |
| 4 | **S2-02 T20** (single sitewide bump + T8 A) + **S2-05** (T25 CLI docs) | Yes — T20 spans site/spec/examples/CHANGELOG/`lock.go`; T25 touches `cli.md` table (after T21) | `version-coherence` = exactly one `draft.6`; §0.8 line verbatim in `/cli/`; `run` implemented, `sign`/`verify` `planned` |
| 5 | **S2-06** (T26 demo) | Single | demo commands run; T26 narrative verbatim; §V.9 dry-runs parse |
| 6 | **S2-07** (T27 release + live verify) | Single (supervisor-executed) | full §V incl. 3/8/9 (local where feasible + CI); PR merged to `master`; live per-task pass/fail recorded; owner sign-off; retrospective |

## Dependencies & parallelism

- **S2-01 ∥ S2-02(T17–19)** — Wave 1, no shared files. Examples use only spec-defined
  keys; the spec sections are additive inserts after §8.6.
- **S2-03 after S2-02(T17–19)** — the rename doesn't strictly need the namespaces, but is
  sequenced here to keep `spec/index.md` edits (T17–19) and `cli.md`/`run.go` edits from
  interleaving mid-wave and to front-load the surface the translators need.
- **S2-04 after S2-03** — translators require the `--backend` subcommand surface and read
  §8.7/§8.8/§8.9. The three translators (T22/T23/T24) are internally independent but share
  `cmd/agentrc/`; commit serially (M-005).
- **T20 after S2-04** — the single sitewide `draft.5`→`draft.6` commit lands only after the
  namespaces are consumed; must NOT be split into earlier tasks (§0.1). T20 also edits
  `examples/Agentfile.code-reviewer` — sequence after S2-01 to avoid an examples collision.
- **S2-05 after S2-03 + S2-04** — the `/cli/` table reflects real backend status (M-002) and
  the renamed flag; runs alongside T20 in Wave 4.
- **S2-06 after T20 + S2-04 + S2-05** — demo uses draft.6, the real backends, and the T20
  code-reviewer comment block.
- **S2-07 last** — needs everything + green §V/CI.

## Cross-cutting gates

The `standards.md` gate matrix applies to every wave. **Where** marks the practical local
gate (given the ENOSPC disk constraint) vs CI-authoritative (open question #2, UNRESOLVED).

| Gate | Command | Where |
|------|---------|-------|
| go-build (scoped) | `go build ./cmd/agentrc/...` | local (scoped) + CI (full `go build ./...`) |
| go-vet (scoped) | `go vet ./cmd/agentrc/...` | local (scoped) + CI (full) |
| go-test (scoped) | `go test ./cmd/agentrc/... ./internal/...` | local (scoped) + CI (`-race`, full) |
| go-mod-tidy | `go mod tidy && git diff --exit-code go.mod go.sum` | local + CI |
| example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" \|\| exit 1; done` | local + CI |
| terminology-split (§V.8) | `[ "$(rg -l -- '--substrate' cmd/ docs/ cli.md 2>/dev/null \| wc -l)" -eq 0 ] && rg -q 'POLICY substrate\.' spec/index.md` | local + CI |
| backend-dryrun-bedrock (§V.9) | `go run ./cmd/agentrc run <ref> --backend bedrock --dry-run \| python3 -m json.tool` | local + CI |
| backend-dryrun-k8s (§V.9) | `go run ./cmd/agentrc run <ref> --backend kubernetes --dry-run \| python3 -c 'import sys,yaml;list(yaml.safe_load_all(sys.stdin))'` (kubeconform in CI) | local (yaml-parse) + CI |
| backend-flag-surface (§V.9) | `go run ./cmd/agentrc run --help \| grep -- '--backend'` | local + CI |
| version-coherence (§V.3) | `[ "$(grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u \| wc -l)" -eq 1 ]` — `draft.5` before T20, `draft.6` after | local + CI |
| hello-canonical (§0.6) | diff each inline hello vs `examples/Agentfile.minimal` | local + CI |
| site-build | `bundle exec jekyll build --trace` | CI (no local Jekyll) |
| internal-link-check (§V.7) | `htmlproofer ./_site --disable-external …` | CI |
| one-canonical-per-page (§V.6) | one `rel="canonical"` per `_site/**/*.html` | CI / live |

All markdown-story RED handles live in `scripts/verify-sprint2.sh` (created by the RED
worker); live checks in `scripts/verify-sprint2-live.sh`. Go-story handles are real
`_test.go` functions under `cmd/agentrc/`.

## Cross-story file conflicts (flagged)

- **`cli.md`** — touched by S2-03 (T21 flag rename, lines 86/103/120) AND S2-05 (T25 status
  table + §0.8 line). Mitigation: T21 (Wave 2) before T25 (Wave 4); T25 re-derives the table
  from real status (M-002) and must not reintroduce `--substrate`.
- **`spec/index.md`** — touched by S2-02 T17–T19 (Wave 1, new §8.7–8.9 at draft.5) AND S2-02
  T20 (Wave 4, version bump + §9 T8 subsection + §14.2 #6). Same story, sequenced across
  waves; no other story edits the spec.
- **`examples/Agentfile.code-reviewer`** — touched by S2-02 T20 (commented `substrate.aws.*`
  + `agent.auth.*` block). S2-01 adds NEW example files only (no code-reviewer edit), so no
  collision, but T20 (Wave 4) is safely after S2-01 (Wave 1) regardless.
- **`cmd/agentrc/run.go`** — S2-03 (T21) rewrites the command; S2-04 (T22–T24) add backend
  dispatch. Sequence T21 → translators; serialize translator commits (M-005).
- **Version strings** — the T20 bump must sweep every `draft.5` surface grep-located first
  (M-003): `spec/index.md`, `_config.yml`/page descriptions, `cli.md`, `examples/*`
  descriptions, CHANGELOG, and schema `0.1.0-draft.5` in `internal/**/lock.go`. §V.3 guards
  a single value.
</content>
