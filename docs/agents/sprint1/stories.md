---
type: stories
sprint: 1
---
# Sprint 1 Stories — "Fixes & Release" (T1–T14)

Stage-2 story breakdown authored from `work-order.md` (T1–T14 + §0 invariants + §V
suite), `intake.md` (five-part Intent), `standards.md` (stack + gate matrix), and
`memory.md` (M-001 doc-drift, M-002 stale-status, M-003 grep-first). Every edit
target below was grep-located against the live repo, not assumed (§0.8, M-003).

## Sprint goal

Make the agentrc spec + `agentrc.ai` site internally consistent and honest, then
ship that consistency to production via PR — without changing the spec's meaning,
version (`0.1.0-draft.5`), or keyword surface (IDENTITY, CAPABILITY, SOP, POLICY).
Purge the ghost `/mnt/tools/ping` tool, make every rendered `hello` byte-identical
to `examples/Agentfile.minimal` and buildable (with `FROM`), tell the truth about
what ships, tidy navigation/SEO, park the workflow draft, investigate (not decide)
the lockfile gap, release via PR, and verify every fix against the live site.

## Sprint demo

1. Run the §V suite locally (`scripts/verify-sprint1.sh`) and show every check green:
   0 `tools/ping` hits, every `hello` carries `FROM python:3.11-slim`, one `draft.5`,
   `arc lint` passes all examples, one `rel="canonical"` per built page, internal
   link check clean, and 0 `workflow-draft|docs/workflows|agent-workflow.yaml` in
   `_site/`.
2. Show the built site: sidebar in the T5 order with no "Workflow draft" entry,
   `/tooling/` linked from `/cli/` and the docs index, `/notes/…` de-indexed, and
   "Platform Conformance" visible where "Runner Conformance" was.
3. Show the lockfile investigation report artifact under `docs/agents/`.
4. Show the merged release PR to `master`, then re-run the live per-task curl checks
   (`scripts/verify-sprint1-live.sh`) against `https://agentrc.ai` all green (T14).

## Definition of Done

- [ ] §V checks 1–7, 10, 11 pass locally before commit (via `scripts/verify-sprint1.sh`).
- [ ] Every task T1–T12's own stated Verify passes locally.
- [ ] Version stays `0.1.0-draft.5`; exactly one `draft.N` value tree-wide (§0.1, §V.3).
- [ ] No new keywords / no new POLICY namespace entries / no substrate rename / no
      inline Cedar or secret keyword introduced (§0.2–§0.4).
- [ ] The single canonical `hello` is propagated, never forked; inline snippets diff
      clean against `examples/Agentfile.minimal` (§0.6, M-001).
- [ ] Release lands on `master` via PR (squash-merge, 0 approvals); no direct push (T13).
- [ ] T14 live per-task verification recorded pass/fail per task ID, all green.
- [ ] Owner sign-off recorded.

## Out of scope

- **Sprint 2 entirely** — T15–T27 (incl. `--backend` / substrate rename) until Sprint 1
  FINAL-GATE + live verification pass and the owner says go.
- **T8 lockfile A/B DECISION** — Sprint 1 is investigate + report only; the owner
  decides Option A vs B at Sprint 2 planning. No site change, no fabricated format.
- **T11 `registry.agentrc.io` host rename** — NO-OP this sprint. Rationale: default is
  LEAVE AS-IS (aspirational own-brand host); a switch to `registry.example.com` happens
  only on the owner's explicit word (§work-order T11). No story is created for T11; it
  is deliberately untouched and carried as an owner decision.
- **§H owner/infra tasks** — CDN purge for `/spec/`, `agentrc.io` DNS 301 → agentrc.ai,
  GitHub/name-collision handling.

## User stories

### S1-01 — Content correctness (T1, T2, T3, T4)

- **What's wanted:** Every rendered agent snippet is correct, buildable, and honest.
  T1: replace the ghost `/mnt/tools/ping` healthcheck with
  `CMD /mnt/tools/file_read --agentrc-schema` on every surface, preserving each file's
  existing HEALTHCHECK options. Grep-located surfaces: `index.md:38`,
  `spec/index.md:636`, `spec/index.md:716`, `examples/Agentfile.minimal:23`,
  `examples/Agentfile.code-reviewer:52`, `examples/index.md:45`. T2: add
  `FROM python:3.11-slim` as the first instruction to every rendered `hello` so it
  builds under BuildKit, and add the FROM-required MUST sentence to spec §2. T3: add
  two quickstart honesty callouts (frontend image not-yet-published at step 2;
  `arc run` planned at step 5), mirroring `/cli/` wording, and reword present-tense
  "no extra tooling" (`docs/quickstart.md:71`) to future tense. T4: `Agentfile.minimal`
  description "Minimal AgentRC agent" → "Minimal agentrc agent"; align the quickstart
  hello to canonical fields (`version=0.1`, `author=acme`, `model.name claude-sonnet-4`).
- **Constraints:** `--agentrc-schema` mirrors existing usage and does NOT resolve
  §14.2 #3 — no prose may claim it is decided (§0.5). One canonical `hello`, propagated
  not forked; inline snippets must diff clean against `examples/Agentfile.minimal`
  (§0.6, M-001). Every example must pass `go run ./cmd/agentrc lint <file>` (§0.7).
  Follow work-order T3 wording (open question #1 flags the published-image discrepancy
  to the owner — do not silently reword). No version bump (§0.1).
- **Failure scenarios:** Partial grep sweep leaves `tools/ping` on one surface (M-003;
  §V.1 non-zero). A hello snippet is hand-edited inline and drifts from source (M-001;
  §V.4 diff fails). A status sentence contradicts `/cli/` reality (M-002). An edit
  accidentally changes `0.1.0-draft.5` (§V.3).
- **Success scenarios:** §V.1 `grep -rn "tools/ping" . | grep -v .git | wc -l` → 0;
  §V.4 every `IDENTITY name=hello` file has `FROM python:3.11-slim` within 3 lines above
  and diffs clean vs minimal; §V.5 `arc lint` passes all examples; T3 both callouts
  present with zero present-tense stock-`docker build` claims; T4 description + fields
  aligned.
- **Connections:** Wave 1, parallel with S1-03 and S1-04. Feeds S1-05 release (T13
  cannot run until §V passes). Shares no files with S1-03's parking sweep.

### S1-02 — Information architecture (T5, T6, T7, T9, T10)

- **What's wanted:** Navigation, cross-links, and SEO are coherent. T5: set sidebar
  (`_layouts/doc.html`) to the order What is agentrc? · Quickstart · Specification ·
  Agentfile · Security · Package · Runners · Conformance · Implementation mapping · CLI ·
  Acknowledgements (no "Workflow draft" — parked by T12), and add Conformance +
  Implementation-mapping cards to `docs/index.md`. T6: de-orphan `/tooling/` by linking
  it from `/cli/` and the docs index (zero inbound links today). T7: kill the duplicate
  URL `/notes/CURRENT_IMPLEMENTATION_MAPPING/` (default: `noindex` + remove from sitemap;
  redirect mechanism is open question #4). T9: relabel visible "Runner Conformance" text
  to "Platform Conformance", keeping the `/profiles/runner-conformance/` URL — located at
  `profiles/index.md:20`, `profiles/runner-conformance.md:3-4`, `llms.txt:40`. T10:
  dedupe the two `<link rel="canonical">` emitters — `_includes/head.html:17` (explicit)
  and `{% seo %}` at `_includes/head.html:36` — to exactly one per page.
- **Constraints:** Sidebar must NOT list workflows, so T5 depends on S1-03 (T12) landing
  first. Keep the `/profiles/runner-conformance/` URL (open question #3 confirms
  relabel-only default). For T10, keep the emitter that produces the correct per-page
  canonical URL; do not delete the working one (intake failure scenario). No version bump.
- **Failure scenarios:** Orphan/dangling link after nav edits breaks the internal link
  check (§V.7). Wrong canonical emitter removed leaves pages with a missing/wrong
  canonical (§V.6). Sidebar still lists Workflow draft because T12 had not landed.
- **Success scenarios:** T5 every built page's sidebar links `/docs/quickstart/` and
  `/docs/conformance/` and has no workflows link; docs index has Conformance +
  Implementation-mapping cards; T6 `/tooling/` reachable from `/cli/` and docs index;
  T9 `grep -rn "Runner Conformance\|Runner conformance" .` → 0 visible-label hits;
  §V.6 every `_site/**/*.html` has exactly one `rel="canonical"`; §V.7 internal links 200.
- **Connections:** Wave 2, after S1-03 (parking) so the sidebar cleanly omits workflows.
  On the critical path (S1-03 → S1-02 → S1-05). T7 sitemap hygiene overlaps T12's sweep.

### S1-03 — Park the workflow draft (T12)

- **What's wanted:** Unpublish the workflow draft now (returns in a future draft), keeping
  all sources in git. Unpublish `/docs/workflows/` (`docs/workflows.md`) and
  `/profiles/workflow-draft/` (`profiles/workflow-draft.md`) via front-matter
  `published: false` or a non-built `parked/` move (mechanism is open question #5), with a
  header comment "Parked 2026-07-03 — will return in a future draft". Drop
  `examples/agent-workflow.yaml` from the served site and its card at `examples/index.md:23`
  (keep the file in-repo under `parked/`). Sweep inbound refs: sidebar
  (`_layouts/doc.html:16`), docs index, profiles index card, `llms.txt:31` and
  `llms.txt:41`, `examples/index.md` prose, and spec/non-goals prose links — unlink to
  "workflow orchestration is parked for a future draft" where a mention must remain (spec
  deferred/non-goals prose may stay). Remove both URLs from the sitemap. Add one changelog
  line under the current draft: "Workflow draft parked (unpublished); returns in a future
  revision."
- **Constraints:** DO NOT delete sources (§work-order T12.1). Spec deferred/non-goals prose
  may remain as prose but must not link to unpublished pages that would 404. No version bump.
  Must land before S1-02's T5 sidebar is finalized (the sidebar must not list workflows).
- **Failure scenarios:** A dangling reference survives (sidebar, docs index, profiles card,
  sitemap, `llms.txt`, or examples index), breaking §V.7 / §V.10 / §V.11. Sources get
  deleted instead of parked. A workflow page still builds into `_site/`.
- **Success scenarios:** Built site has no `/docs/workflows/` or `/profiles/workflow-draft/`
  page; §V.11 `grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ | wc -l`
  → 0; §V.10 sitemap has 0 `workflow` entries (post-deploy); internal link check clean;
  sources still tracked in git; one changelog line added.
- **Connections:** Wave 1, parallel with S1-01 and S1-04. Blocks S1-02 (T5 sidebar). First
  node on the critical path (S1-03 → S1-02 → S1-05). Its inbound sweep intersects T5, T6, T7.

### S1-04 — Lockfile investigation report (T8)

- **What's wanted:** Investigate and report only. Read the `arc lock` implementation under
  `cmd/agentrc` + `internal/` (note standards correction: the CLI is at repo-root
  `./cmd/agentrc`, not `tooling/`), document exactly what it emits (filename, format,
  records, how `arc build` consumes it), and write findings + an A/B recommendation to a
  new report artifact under `docs/agents/` (e.g. `docs/agents/sprint1/lockfile-report.md`).
  Option A: spec subsection under §9 documenting only what the tooling does, marked
  "Status: informative in this draft; format TODO", slogan kept. Option B: cut the slogan
  sentence, lockfile lives on `/cli/` only.
- **Constraints:** DO NOT fabricate a lockfile format. DO NOT change site content in
  Sprint 1. DO NOT make the A/B decision (owner decides at Sprint 2 planning). Report is an
  investigation artifact, not a spec/site edit. No version bump.
- **Failure scenarios:** Scope creep — the task fabricates a format or edits site content
  (violates its explicit investigate-only boundary). The report picks a winner instead of
  presenting A/B for the owner.
- **Success scenarios:** A report artifact exists under `docs/agents/` documenting the
  actual `arc lock` output and `build` consumption, ending in an A/B recommendation with no
  decision taken; zero site-content diffs attributable to T8.
- **Connections:** Wave 1, parallel with S1-01 and S1-03. Off the critical path. Its output
  feeds Sprint 2 planning, not this sprint's release.

### S1-05 — Release & live verification (T13, T14)

- **What's wanted:** Ship and verify. T13: local Jekyll build, run full §V (checks 1–7, 10,
  11) green, commit with a task-tagged message, open + squash-merge a PR to `master`
  (master is protected: PR required, 0 approvals), wait for GH Pages build + CDN. T14:
  re-run each task's Verify against production (`https://agentrc.ai`) via
  `scripts/verify-sprint1-live.sh`, record pass/fail per task ID; any FAIL is fixed within
  Sprint 1 before FINAL-GATE.
- **Constraints:** No direct push to `master` — deploy = merge the sprint branch via PR
  (§work-order T13 note). Do not release until all of S1-01…S1-04 are in and local §V
  passes. This story is supervisor-executed, not worker-TDD; its plan.md checks are the §V
  suite + live curl checks, not RED unit tests.
- **Failure scenarios:** A direct push to `master` is attempted and rejected by the ruleset.
  Release proceeds before local §V is green. A live check FAILs post-deploy and is not fixed
  within the sprint.
- **Success scenarios:** Local §V all green; PR merged to `master`; GH Pages + CDN serve the
  new content; §V.10 sitemap-hygiene curl → 0 `CURRENT_IMPLEMENTATION_MAPPING|workflow`; live
  per-task verification recorded all green; owner sign-off.
- **Connections:** Wave 3, after every other story. Depends on S1-01…S1-04 + local §V (T13)
  and on T13's deploy (T14). Last node on the critical path.

## Story dependency graph

```
S1-03 (park) ──┐
S1-01 (content)├─▶ S1-02 (IA) ─▶ S1-05 (release + live verify)
S1-04 (lockfile)┘        ▲
                         └── S1-02 also needs S1-03 landed (sidebar excludes workflows)

Wave 1 (parallel): S1-03, S1-01, S1-04
Wave 2:            S1-02   (after S1-03 parks workflows)
Wave 3:            S1-05   (after S1-01..S1-04 + local §V)
Critical path:     S1-03 → S1-02 → S1-05
```

T11 is a deliberate NO-OP (see Out of scope) and is not a node in this graph.

## Notes / open questions

Carried forward from `intake.md` for the supervisor to surface to the owner — NOT resolved
here:

1. **T3 published-vs-unpublished frontend image (BLOCKING).** Work-order T3 wording says the
   frontend image is "not yet published", but this session published
   `ghcr.io/adeelahmad/agentrc-frontend` (release v0.1.0). Which is authoritative — the
   work-order wording or ground truth? If ground truth, T3's callout text needs owner-approved
   rewording before it ships (else it re-creates the M-002 stale-status drift this sprint kills).
2. **T1 default option confirmation.** T1 defaults to Option A (`file_read --agentrc-schema`);
   confirm it is not an Option B variant (the choice is baked into every affected surface).
3. **T9 relabel-only vs new slug.** Default is relabel visible text only, keeping
   `/profiles/runner-conformance/`; confirm the owner does not want a new `/profiles/platform/`
   slug (which would change the URL and need redirect handling).
4. **T7 redirect mechanism on Jekyll/GitHub Pages.** How to kill
   `/notes/CURRENT_IMPLEMENTATION_MAPPING/` — `jekyll-redirect-from` (301), meta-refresh, or
   `noindex` + sitemap removal? GH Pages has no server-side redirects; the mechanism affects
   T14's live curl verification.
5. **T12 unpublish mechanism on GH Pages.** Front-matter `published: false` vs move sources
   into a non-built `parked/` directory — the generator's behavior determines whether sources
   still build and whether `_site/` cleanly excludes them for §V checks 10/11 while remaining
   in git.
