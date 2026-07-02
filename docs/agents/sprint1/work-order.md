# AGENTRC WORK ORDER v2 — SPRINT 1 slice (intake input)

Staged by the supervisor as the authoritative intake document for Sprint 1
("Fixes & Release", tasks T1–T14). Sprint 2 (T15–T27) is out of scope until
Sprint 1's FINAL-GATE + live verification pass and the owner says go.

Repo: /Users/adeelahmad/work/agentrc (`master`) · Live: https://agentrc.ai
(Jekyll site + Go tooling under `tooling/`, `cmd/`, `internal/`).

## §0 GLOBAL INVARIANTS (never violate)

1. Version `0.1.0-draft.5` stays unchanged through all of Sprint 1. Syntax line
   stays `# syntax=agentrc.agentfile/v0.1`.
2. `substrate.*` (POLICY namespace) and "substrate-neutral" (spec concept) are
   NOT renamed in Sprint 1 (that's Sprint 2's `--backend` work).
3. Exactly four keywords (IDENTITY, CAPABILITY, SOP, POLICY). No new keywords.
   Sprint 1 must NOT add POLICY namespace entries the live spec doesn't define.
4. POLICY lines are requests; Cedar is platform-side only; secrets deferred.
   Never write inline Cedar or a secret keyword anywhere.
5. Open decisions (spec §14.2) stay open; follow each task's stated default.
6. ONE hello: byte-identical wherever rendered inline and identical to
   `examples/Agentfile.minimal` (after T1/T2). Propagate, don't fork.
7. Every example file must pass `arc lint` (the built `agentrc`/`arc` CLI's lint).
8. Locate every target by grep FIRST — repo layout is not assumed.
9. Rebuild the site locally and pass the §V verification suite before any commit.

## SPRINT 1 TASKS

### T1 — Kill the ghost healthcheck tool (`/mnt/tools/ping`)  [P0]
Occurrences: `examples/Agentfile.minimal`, `examples/Agentfile.code-reviewer`,
homepage hero (`index.md`), examples index (`examples/index.md`), the spec's
worked example(s). Every affected file already `COPY`s `file_read`;
`Agentfile.secure-workspace` is already correct.
Fix (default Option A) — replace each ghost line, KEEPING each file's existing
HEALTHCHECK options:
`HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema`
`--agentrc-schema` mirrors existing usage and does NOT resolve §14.2 #3 — no prose
claiming it's decided.
Locate: `grep -rn "tools/ping" --include="*.md" --include="*.html" --include="Agentfile*" .`
Verify: `grep -rn "tools/ping" . | grep -v .git | wc -l` → 0

### T2 — One hello, always with `FROM`  [P0]
Homepage, quickstart step 1, and examples index render hello WITHOUT `FROM`
(BuildKit fails as pasted). 1) Add `FROM python:3.11-slim` so every rendering is
byte-identical to `examples/Agentfile.minimal`. 2) Add to spec §2:
"Exactly as in Dockerfile, every Agentfile MUST contain a `FROM` instruction, and
`FROM` must be the first instruction after the `# syntax=` line, comments, and any
`ARG` that `FROM` consumes."
Locate: `grep -rln "IDENTITY name=hello" --include="*.md" --include="*.html" .`
Verify: every located file has the FROM within 3 lines above hello's IDENTITY;
spec sentence present; snippets diff clean against `Agentfile.minimal`.

### T3 — Quickstart honesty callouts  [P0]
Mirror `/cli/` wording exactly; invent no status claims. At step 2:
"Status: the agentrc frontend image is not yet published to a public registry, so
`docker build -f Agentfile .` won't auto-route through it yet. Build the frontend
locally first (see `tooling/README.md`) or pass
`--build-arg BUILDKIT_SYNTAX=<your-built-image>`. Details on the CLI page (/cli/)."
At step 5:
"Status: `arc run` is planned — agentrc declares agents, it does not ship a
runtime. See the CLI status table (/cli/)."
Reword any remaining present-tense "nothing to install"/"no extra tooling" to
future tense.
Locate: `grep -rn "nothing to install\|no extra tooling\|arc run" docs/`
Verify: both callouts present; zero present-tense stock-`docker build` claims.

NOTE (supervisor): the owner's work order says the frontend image is "not yet
published", but this session actually published `ghcr.io/adeelahmad/agentrc-frontend`
(release v0.1.0). The work-order text T3 is authoritative for Sprint 1 wording; the
published-image discrepancy is flagged for the owner to reconcile at planning, NOT
silently changed. Follow the work order's stated wording unless the owner overrides.

### T4 — Polish  [P3]
`Agentfile.minimal` description "Minimal AgentRC agent" → "Minimal agentrc agent";
align quickstart hello to canonical (`version=0.1`, `author=acme`,
`model.name claude-sonnet-4`) per §0.6.

### T5 — Sidebar + docs-index coverage  [P1]
Sidebar order: What is agentrc? · Quickstart · Specification · Agentfile · Security ·
Package · Runners · Conformance · Implementation mapping · CLI · Acknowledgements.
(No "Workflow draft" entry — parked by T12.) Add Examples if the sidebar pattern
allows. Docs index: add Conformance and Implementation-mapping cards (Conformance:
one sentence on profile-based conformance + adversarial suite).
Verify: every built page's sidebar links `/docs/quickstart/` and `/docs/conformance/`;
no sidebar link to workflows; §V link check clean.

### T6 — De-orphan `/tooling/`  [P1]
Zero inbound links today. Link from the CLI page ("Reference implementation
available" + "Try it now") and the docs index.

### T7 — Kill the duplicate URL  [P1]
`/notes/CURRENT_IMPLEMENTATION_MAPPING/` duplicates `/docs/implementation-mapping/`.
Redirect or `noindex` the `/notes/` URL; remove from sitemap.
Verify: `curl -s https://agentrc.ai/sitemap.xml | grep -c CURRENT_IMPLEMENTATION_MAPPING` → 0 post-deploy.

### T8 — Lockfile: INVESTIGATE + REPORT ONLY  [P1]
Homepage slogan claims a lockfile; `arc lock` is implemented; the spec has zero
lockfile content. In Sprint 1: read the `arc lock` implementation in `tooling/`,
document exactly what it emits (filename, format, records, how `build` consumes it),
and write findings + an A/B recommendation to `docs/agents/` as a report artifact.
Options for the owner at Sprint 2 planning: A) spec subsection under §9 documenting
only what the tooling does, marked "Status: informative in this draft; format TODO",
slogan kept; B) cut the slogan sentence, lockfile lives on `/cli/` only.
DO NOT fabricate a lockfile format. DO NOT change site content in Sprint 1.
Locate: `grep -rn -i "lockfile\|agentrc.lock\|arc lock" . | grep -v .git`

### T9 — "Runner" → "Platform" naming  [P1]
Slug `/profiles/runner-conformance/` + index label say "Runner Conformance"; the
page and conformance table say Platform Conformance (`agentrc/platform/v0.1`).
Default: relabel visible link text only; keep the URL.
Verify: `grep -rn "Runner Conformance\|Runner conformance" .` → 0 visible-label hits.

### T10 — Deduplicate canonical tags  [P2]
Two `<link rel="canonical">` per page (two emitters in layout/head includes). Keep
exactly one.
Verify (local build): `grep -c 'rel="canonical"' <page>` → 1, all pages.

### T11 — `registry.agentrc.io` example hosts  [P2]
Default: LEAVE AS-IS (aspirational own-brand). Switch to `registry.example.com` only
on the owner's explicit word. No change in Sprint 1 unless owner says so.

### T12 — PARK the workflow draft  [P1]
Owner decision: workflows come back later; unpublish now.
1. Unpublish `/docs/workflows/` and `/profiles/workflow-draft/` (front-matter
   `published: false` / build-exclude — DO NOT delete sources; move under a `parked/`
   area if the generator needs it, with header comment
   "Parked 2026-07-03 — will return in a future draft").
2. Remove `examples/agent-workflow.yaml` from the served site and its card/link on
   the examples index (keep the file in-repo under `parked/`).
3. Sweep inbound refs: sidebar (T5), docs index, profiles index card, homepage/spec
   prose links → unlinked "workflow orchestration is parked for a future draft"
   where a mention must remain (spec deferred/non-goals prose may stay).
4. Remove both URLs from the sitemap; update `llms.txt` if it lists them.
5. One changelog line under the current draft: "Workflow draft parked (unpublished);
   returns in a future revision."
Locate: `grep -rni "workflow" --include="*.md" --include="*.html" --include="*.yml" --include="*.yaml" . | grep -v .git`
Verify: built site has no `/docs/workflows/` or `/profiles/workflow-draft/` page;
`grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ | wc -l` → 0;
link checker clean; sources still in git.

### T13 — RELEASE Sprint 1  [P0]
Local build; run full §V (checks 1–7, 10, 11); commit with task-tagged message;
deploy (GH Pages). Wait for Pages build + CDN before T14.
NOTE (supervisor): pushing to `master` requires a PR (repo ruleset: squash-merge,
0 approvals). "Deploy" = merge the sprint branch via PR; the owner is informed.

### T14 — LIVE per-task verification (against production)  [P0]
Re-run each task's verify against the live site; record pass/fail per task ID in the
sprint report. (Full curl suite in the work order §14.) Any FAIL → fix within
Sprint 1 before FINAL-GATE. Sprint 1 done = all green + owner sign-off.

## §V VERIFICATION SUITE (Sprint 1 subset — all must pass before commit)
1. ghost tool: `grep -rn "tools/ping" . | grep -v .git | wc -l` → 0
2. stale model markers: `grep -rnE 'utci:|permit\(|forbid\(' --include="*.md" --include="Agentfile*" . | grep -v .git` → review → 0 real
3. version coherence: exactly one `draft.N` value (draft.5 in Sprint 1)
4. hello + FROM: every file with `IDENTITY name=hello` has `FROM python:3.11-slim`
5. lint all examples: `for f in examples/Agentfile.*; do (cd tooling && go run ./cmd/agentrc lint "../$f"); done` all pass
6. one canonical per page: every `_site/**/*.html` has exactly 1 `rel="canonical"`
7. full internal link check: every internal href in `_site/` → 200
10. sitemap hygiene (post-deploy): sitemap has 0 `CURRENT_IMPLEMENTATION_MAPPING|workflow`
11. workflow parked: `grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ | wc -l` → 0

## §H MANUAL / INFRA (owner, not the pipeline)
CDN purge for `/spec/`; `agentrc.io` DNS 301 → agentrc.ai; GitHub/name-collision.
These are owner tasks, out of the pipeline's scope.
