---
type: standards
sprint: 1
---

# Sprint 1 Standards — the active rule digest and gate matrix

Read-only digest produced by the standards activity (The Lawkeeper). Every rule
below cites a real source in this repo; nothing is invented. The gate matrix at
the bottom is what the GREEN gate (per-change, local) and the FINAL gate
(pre-commit / pre-deploy) run.

## Detected stack

**Jekyll / Ruby static site** (GitHub Pages, deployed to https://agentrc.ai):
- `_config.yml` — `markdown: kramdown`, `permalink: pretty`, plugins
  `jekyll-seo-tag` + `jekyll-sitemap`; `url: https://agentrc.ai`.
- `Gemfile` — `gem "github-pages"` plus `jekyll-seo-tag`, `jekyll-sitemap`.
- `_layouts/` (default layout), `_includes/` (incl. `head.html` which emits the
  canonical link), `index.md`, `docs/`, `spec/`, `examples/`, `profiles/`,
  `notes/`, `CNAME` (`agentrc.ai`).
- Build command per `README_DEV.md`: `bundle exec jekyll serve`; production build
  `bundle exec jekyll build`.

**Go module for the agentrc / arc CLI** (built at repo root, NOT under `tooling/`):
- `go.mod` at repo root — module `github.com/adeelahmad/agentrc`, `go 1.25.9`.
- CLI entrypoint is `./cmd/agentrc` (verified: `cmd/agentrc/main.go` registers
  `newLintCmd()` plus build/lock/pull/push/sign/verify/run/init/inspect).
  Internal packages: `internal/agentfile`, `internal/llb`.
- **Correction to the work order:** `tooling/` contains only `README.md` — there
  is no Go code there. The work order's `(cd tooling && go run ./cmd/agentrc …)`
  phrasing is wrong. The verified, runnable command from repo root is
  `go run ./cmd/agentrc lint <file>`. Confirmed working:
  `go run ./cmd/agentrc lint examples/Agentfile.minimal` prints
  `examples/Agentfile.minimal: ok`. The CI `go` job also builds from repo root
  (`go build ./...`), corroborating this.

**GitHub Actions workflows** (`.github/workflows/`):
- `ci.yml` — job `build` (Jekyll build + html-proofer link check) and job `go`
  ("Build, vet & test tooling/": `go build ./...`, `go vet ./...`,
  `go test -race ./...`).
- `pages.yml` — GitHub Pages build + deploy (`actions/jekyll-build-pages@v1`,
  `deploy-pages@v4`) on push to `master`/`main`.
- `release.yml` — release workflow.

Deploy path: GitHub Pages via `pages.yml` on push to `master` (per T13, "deploy"
= merge the sprint branch via PR).

## Active rules

- `docs/agents/sprint1/work-order.md` §0.1–§0.9 — the §0 global invariants are
  binding for all of Sprint 1 and override any "improvement" beyond the slice.
- `docs/agents/sprint1/work-order.md` §0.1 (+ §V check 3): version stays
  `0.1.0-draft.5`, syntax line stays `# syntax=agentrc.agentfile/v0.1`, and
  exactly one `draft.N` value exists in the tree.
- `docs/agents/sprint1/work-order.md` §0.2–§0.4: no renames of `substrate.*` /
  "substrate-neutral"; exactly four keywords (IDENTITY, CAPABILITY, SOP, POLICY),
  no new keywords or POLICY entries the live spec lacks; POLICY lines are
  requests; never write inline Cedar or a secret keyword.
- `docs/agents/sprint1/work-order.md` §0.8–§0.9 (grep-locate-first): locate every
  target by grep before editing — repo layout is not assumed.
- `docs/agents/sprint1/work-order.md` §0.6 (+ §V check 4 / T2): one canonical
  hello, byte-identical everywhere it renders inline and identical to
  `examples/Agentfile.minimal`; propagate not fork; every rendered hello carries
  `FROM python:3.11-slim`.
- `docs/agents/sprint1/work-order.md` §0.7 (+ §V check 5): every example file must
  pass the built CLI's lint, run as `go run ./cmd/agentrc lint <file>` from repo
  root (path corrected — the CLI is not under `tooling/`).
- `docs/agents/sprint1/work-order.md` §0.9: the Jekyll build must succeed locally
  before any commit (`bundle exec jekyll build`); also enforced by the `build`
  job in `.github/workflows/ci.yml`.
- `docs/agents/sprint1/work-order.md` §V (Go changes) via `.github/workflows/ci.yml`
  `go` job: `go build ./...`, `go vet ./...`, `go test -race ./...` must pass for
  any tooling touched.
- `docs/agents/sprint1/work-order.md` T10 / §V check 6: exactly one
  `<link rel="canonical">` per built page — dedupe the two emitters
  (`jekyll-seo-tag`'s `{% seo %}` and the explicit tag in `_includes/head.html`).
- `docs/agents/sprint1/work-order.md` §0.9 and §V checks 1–7, 10, 11: the full §V
  verification suite must pass before any commit.
- No project `CLAUDE.md` exists at repo root (only `.claude/scheduled_tasks.lock`);
  there are no project-local Claude conventions to cite — the authoritative rule
  source for this sprint is `docs/agents/sprint1/work-order.md`.

## Gate matrix

All commands are literal and runnable from the repo root
(`/Users/adeelahmad/work/agentrc`).

| Gate | Command | Applies to |
|------|---------|------------|
| site-build | `bundle exec jekyll build` | any change (site must compile before commit; §0.9) |
| internal-link-check | `gem install html-proofer --no-document && htmlproofer ./_site --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https` (same invocation as the `build` job in `.github/workflows/ci.yml`) | any change that touches pages/links; §V check 7 |
| example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" || exit 1; done` | every example Agentfile; §0.7 / §V check 5 |
| ghost-tool-grep | `[ "$(grep -rn "tools/ping" . \| grep -v .git \| wc -l)" -eq 0 ]` | whole tree; §V check 1 (T1) |
| hello+FROM-diff | `for f in $(grep -rln "IDENTITY name=hello" --include="*.md" --include="*.html" .); do grep -q "FROM python:3.11-slim" "$f" || { echo "MISSING FROM: $f"; exit 1; }; done` | every file rendering hello; §V check 4 (T2) |
| version-coherence | `[ "$(grep -rhoE "draft\.[0-9]+" . \| grep -v .git \| sort -u \| wc -l)" -eq 1 ] && grep -rq "draft\.5" .` | whole tree — exactly one `draft.N`, and it is `draft.5`; §V check 3 / §0.1 |
| one-canonical-per-page | `fail=0; for p in $(find _site -name '*.html'); do [ "$(grep -c 'rel="canonical"' "$p")" -eq 1 ] || { echo "BAD: $p"; fail=1; }; done; [ "$fail" -eq 0 ]` | every built `_site/**/*.html`; §V check 6 (T10) |
| go-build/vet/test | `go build ./... && go vet ./... && go test -race ./...` | only when Go source under `cmd/`, `internal/` changes; `go` job in `.github/workflows/ci.yml` |

Post-deploy / live gates (owner-run against production per T14, not part of the
local GREEN gate):

| Gate | Command | Applies to |
|------|---------|------------|
| sitemap-hygiene | `[ "$(curl -s https://agentrc.ai/sitemap.xml \| grep -cE 'CURRENT_IMPLEMENTATION_MAPPING\|workflow')" -eq 0 ]` | post-deploy; §V check 10 (T7, T12) |
| workflow-parked | `[ "$(grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ \| wc -l)" -eq 0 ]` | built site; §V check 11 (T12) |
