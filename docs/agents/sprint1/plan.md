---
type: sprint-plan
sprint: 1
stage: "2"
---
# Sprint 1 Plan — Stage 2 (waves, critical path, gates)

Execution plan for the five stories (S1-01…S1-05) covering work-order tasks T1–T14.
T11 is a deliberate NO-OP (see `stories.md` › Out of scope). Every gate command is
literal and runnable from repo root (`/Users/adeelahmad/work/agentrc`).

## Sprint goal

Make the agentrc spec + `agentrc.ai` site internally consistent and honest, then ship
that consistency to production via PR — without changing the spec's meaning, version
(`0.1.0-draft.5`), or keyword surface. Success = §V checks 1–7, 10, 11 green locally,
every task's Verify green, then the same verified live (T14) with owner sign-off.

## Critical path

```
S1-03 (park workflow draft)  →  S1-02 (information architecture)  →  S1-05 (release + live verify)
```

- **S1-03 first** because S1-02's T5 sidebar must not list workflows; the sidebar can
  only be finalized once T12 has parked `/docs/workflows/` and `/profiles/workflow-draft/`.
- **S1-02 second** because navigation/SEO coherence (sidebar, docs-index cards, canonical
  dedupe, `/tooling/` de-orphaning) must be settled before the site is built for release.
- **S1-05 last** because T13 release requires all content/nav/SEO fixes in and local §V
  green, and T14 live verification requires T13's deploy (PR merge + Pages build + CDN).

S1-01 (content) and S1-04 (lockfile report) are off the critical path but are release
blockers for S1-05 (T13 needs §V green, which includes S1-01's ghost-tool and hello checks).

## Waves

### Wave 1 — parallel: S1-03, S1-01, S1-04

- **S1-03 (park workflow draft, T12)** — unpublish workflow pages, drop
  `agent-workflow.yaml` from the served site, sweep inbound refs, sitemap + `llms.txt`
  cleanup, changelog line. Touches `docs/workflows.md`, `profiles/workflow-draft.md`,
  `examples/agent-workflow.yaml`, `examples/index.md`, `_layouts/doc.html`, `llms.txt`,
  `CHANGELOG.md`.
- **S1-01 (content correctness, T1–T4)** — ghost-tool purge, hello+FROM, quickstart
  honesty callouts, polish. Touches `index.md`, `spec/index.md`, `examples/Agentfile.*`,
  `examples/index.md`, `docs/quickstart.md`.
- **S1-04 (lockfile investigation, T8)** — report-only; writes
  `docs/agents/sprint1/lockfile-report.md`. Touches no site content.

These three share no edit targets (verified by grep), so they run fully in parallel.

### Wave 2 — S1-02 (information architecture, T5, T6, T7, T9, T10)

Runs after S1-03 lands so the sidebar (`_layouts/doc.html`) can be set to the T5 order
with no "Workflow draft" entry. Touches `_layouts/doc.html`, `docs/index.md`, `cli.md`,
`_includes/head.html`, `notes/CURRENT_IMPLEMENTATION_MAPPING.md`, `profiles/index.md`,
`profiles/runner-conformance.md`, `llms.txt`.

### Wave 3 — S1-05 (release + live verification, T13, T14)

Runs after S1-01…S1-04 are complete and local §V is green. T13 builds, runs §V, commits,
opens + squash-merges the PR to `master`; T14 re-runs per-task verification live.

## Parallelism

| Wave | Stories (parallel within wave) | Depends on | Shared files / conflict risk |
|------|--------------------------------|------------|------------------------------|
| 1 | S1-03, S1-01, S1-04 | — | None — disjoint edit targets (grep-verified) |
| 2 | S1-02 | S1-03 (sidebar excludes workflows) | `_layouts/doc.html` also touched by S1-03 (sidebar link removal) — sequence S1-03 → S1-02 on that file |
| 3 | S1-05 | S1-01, S1-02, S1-03, S1-04 + local §V | Read-only over the whole tree + `_site/`; PR merge |

Note the one cross-wave file conflict: `_layouts/doc.html` is edited by S1-03 (remove the
Workflow-draft sidebar link) and S1-02 (reorder + add Quickstart/Conformance/Implementation
mapping). Serializing S1-03 → S1-02 resolves it; both land on the same final sidebar.

## Cross-cutting gates

The GREEN gate (per-change, local) and FINAL gate (pre-commit / pre-deploy) run the
`standards.md` gate matrix. All commands are literal, from repo root:

| Gate | Command | Applies to |
|------|---------|------------|
| site-build | `bundle exec jekyll build` | any change (§0.9) |
| internal-link-check | `htmlproofer ./_site --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https` | pages/links; §V.7 |
| example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" || exit 1; done` | every example; §0.7 / §V.5 |
| ghost-tool-grep | `[ "$(grep -rn "tools/ping" . \| grep -v .git \| wc -l)" -eq 0 ]` | whole tree; §V.1 (T1) |
| hello+FROM-diff | `for f in $(grep -rln "IDENTITY name=hello" --include="*.md" --include="*.html" .); do grep -q "FROM python:3.11-slim" "$f" \|\| { echo "MISSING FROM: $f"; exit 1; }; done` | every hello render; §V.4 (T2) |
| version-coherence | `[ "$(grep -rhoE "draft\.[0-9]+" . \| grep -v .git \| sort -u \| wc -l)" -eq 1 ] && grep -rq "draft\.5" .` | whole tree; §V.3 / §0.1 |
| one-canonical-per-page | `fail=0; for p in $(find _site -name '*.html'); do [ "$(grep -c 'rel="canonical"' "$p")" -eq 1 ] \|\| fail=1; done; [ "$fail" -eq 0 ]` | every built page; §V.6 (T10) |
| go-build/vet/test | `go build ./... && go vet ./... && go test -race ./...` | only if Go source under `cmd/`/`internal/` changes |
| workflow-parked | `[ "$(grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ \| wc -l)" -eq 0 ]` | built site; §V.11 (T12) |

Post-deploy / live (owner-run against production per T14, not part of local GREEN):

| Gate | Command | Applies to |
|------|---------|------------|
| sitemap-hygiene | `[ "$(curl -s https://agentrc.ai/sitemap.xml \| grep -cE 'CURRENT_IMPLEMENTATION_MAPPING\|workflow')" -eq 0 ]` | post-deploy; §V.10 (T7, T12) |

The RED worker creates a single `scripts/verify-sprint1.sh` that wraps the local checks
above (one function per assertion, referenced by the story `plan.md` files) and a
`scripts/verify-sprint1-live.sh` for the T14 live curl checks. GREEN makes the edits pass.
