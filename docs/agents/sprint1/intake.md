---
type: intake
sprint: 1
---
# Sprint 1 Intake — "Fixes & Release" (T1–T14)

Interrogated from `docs/agents/sprint1/work-order.md`, grounded against `README.md`,
`docs/agents/memory.md`, `docs/index.md`, and `cli.md`. This is intake only: no plan,
no code. Blocking ambiguities are listed under Open questions, not resolved.

## What's wanted

Sprint 1 makes the agentrc spec + `agentrc.ai` site internally consistent and honest,
then ships that consistency to production — without touching the spec's meaning, version,
or keyword surface. Concretely: purge a ghost tool, make every rendered `hello` snippet
byte-identical to its canonical source and buildable, tell the truth about what is and
isn't shipped, tidy navigation/SEO, park the workflow draft, investigate (not decide) the
lockfile gap, then release via PR and verify every fix against the live site. Success is
the §V verification suite passing locally, each task's Verify passing, and — after deploy —
the same verifications passing against production (T14). The version string stays
`0.1.0-draft.5` throughout.

- **T1 [P0]** — Replace the ghost healthcheck tool `/mnt/tools/ping` with `file_read --agentrc-schema` on every surface (examples, homepage hero, examples index, spec worked example), preserving each file's existing HEALTHCHECK options.
- **T2 [P0]** — Ensure every rendered `hello` snippet has `FROM python:3.11-slim` as the first instruction (so it builds under BuildKit) and is byte-identical to `examples/Agentfile.minimal`; add the `FROM`-required sentence to spec §2.
- **T3 [P0]** — Add two honesty callouts to the quickstart (frontend image not-yet-published at step 2; `arc run` is planned at step 5), mirroring `/cli/` wording exactly; reword present-tense "nothing to install"/"no extra tooling" to future tense.
- **T4 [P3]** — Polish: `Agentfile.minimal` description "Minimal AgentRC agent" → "Minimal agentrc agent"; align quickstart hello to canonical fields (`version=0.1`, `author=acme`, `model.name claude-sonnet-4`).
- **T5 [P1]** — Fix sidebar order + coverage (no "Workflow draft" entry), add Examples if the pattern allows, and add Conformance + Implementation-mapping cards to the docs index.
- **T6 [P1]** — De-orphan `/tooling/` by linking to it from the CLI page and the docs index (zero inbound links today).
- **T7 [P1]** — Kill the duplicate URL `/notes/CURRENT_IMPLEMENTATION_MAPPING/` (redirect or `noindex`) and remove it from the sitemap.
- **T8 [P1]** — INVESTIGATE + REPORT ONLY: read the `arc lock` implementation, document exactly what it emits and how `build` consumes it, and write findings + an A/B recommendation as a report artifact under `docs/agents/`. No site content changes, no fabricated format.
- **T9 [P1]** — Relabel visible "Runner Conformance" text to "Platform Conformance" (default: label text only, keep the `/profiles/runner-conformance/` URL).
- **T10 [P2]** — Deduplicate the two `<link rel="canonical">` emitters so every page has exactly one.
- **T11 [P2]** — `registry.agentrc.io` example hosts: LEAVE AS-IS unless the owner explicitly says switch to `registry.example.com`.
- **T12 [P1]** — PARK the workflow draft: unpublish `/docs/workflows/` and `/profiles/workflow-draft/`, drop `examples/agent-workflow.yaml` from the served site, sweep inbound refs, remove both URLs from sitemap + `llms.txt`, add one changelog line — keeping all sources in git under a `parked/` area.
- **T13 [P0]** — RELEASE: local build, full §V (checks 1–7, 10, 11), task-tagged commit, deploy to GH Pages via PR (master requires a PR). Wait for Pages build + CDN before T14.
- **T14 [P0]** — LIVE verification: re-run each task's Verify against production, record pass/fail per task ID; any FAIL is fixed within Sprint 1 before FINAL-GATE.

## Constraints

Global invariants from work-order §0 (never violate):

- **Version frozen** — `0.1.0-draft.5` stays unchanged through all of Sprint 1; syntax line stays `# syntax=agentrc.agentfile/v0.1`.
- **Exactly four keywords** — IDENTITY, CAPABILITY, SOP, POLICY. No new keywords; do NOT add POLICY namespace entries the live spec doesn't already define.
- **No new POLICY namespaces / no substrate rename** — `substrate.*` and "substrate-neutral" are NOT renamed (that is Sprint 2's `--backend` work).
- **Secrets deferred, no inline Cedar** — POLICY lines are requests only; Cedar is platform-side; never write inline Cedar or a secret keyword anywhere.
- **Open decisions stay open** — spec §14.2 decisions remain open; follow each task's stated default; `--agentrc-schema` (T1) does NOT resolve §14.2 #3, and no prose may claim it is decided.
- **ONE canonical hello** — byte-identical wherever rendered inline and identical to `examples/Agentfile.minimal` (after T1/T2). Propagate the canonical, never fork it.
- **Every example passes `arc lint`** — the built `agentrc`/`arc` CLI's lint must pass on every example file.
- **Grep-locate first** — locate every edit target by grep before editing; repo layout is not assumed (same string recurs across `.md`, `.html`, and `Agentfile*`).
- **Rebuild + verify before commit** — rebuild the site locally and pass the §V verification suite before any commit.
- **Environment: master is protected** — pushing to `master` requires a PR (repo ruleset: squash-merge, 0 approvals). "Deploy" (T13) means merging the sprint branch via PR, not a direct push.
- **§H is out of the pipeline** — CDN purge, `agentrc.io` DNS 301, GitHub/name-collision are owner/infra tasks, not this pipeline's work.

## Failure scenarios

Concrete ways this sprint goes wrong (each is something the plan must actively prevent):

- **Partial grep sweep** — `tools/ping` (or `IDENTITY name=hello`) is fixed on some surfaces but survives on one (e.g. the spec worked example or examples index), leaving the site internally inconsistent (memory M-003; §V check 1 → non-zero).
- **Inline hello drifts from source** — a hello snippet is hand-edited inline and no longer byte-matches `examples/Agentfile.minimal`; a grep passes but a diff would fail (memory M-001; §V check 4 / T2 diff).
- **Orphan link after parking** — T12 unpublishes a page but leaves a dangling reference in the sidebar, docs index, profiles index card, sitemap, or `llms.txt`, breaking the internal link check (§V checks 7, 10, 11).
- **Stale status claim vs CLI reality** — a doc "status" sentence contradicts the CLI page or the actual published/planned reality (memory M-002; T3), re-introducing the exact drift this sprint is fixing.
- **Wrong canonical emitter removed** — T10 dedupe deletes the emitter that produces the correct per-page canonical URL rather than the duplicate, leaving pages with a wrong or missing canonical (§V check 6).
- **Accidental version bump** — an edit changes `0.1.0-draft.5` (or introduces a second `draft.N`), violating §0.1 and failing §V check 3 (version coherence).
- **T8 scope creep** — the lockfile task fabricates a format or changes site content instead of investigate-and-report-only, violating its explicit boundary.
- **Deploy without PR** — a direct push to `master` is attempted and rejected by the ruleset, or bypasses the required PR/squash-merge path (T13 environment constraint).

## Success scenarios

Acceptance = the §V verification suite passing plus each task's stated Verify, first locally then live (T14):

- **§V check 1** — `grep -rn "tools/ping" . | grep -v .git | wc -l` → 0.
- **§V check 2** — stale model markers (`utci:`, `permit(`, `forbid(`) reviewed → 0 real occurrences.
- **§V check 3** — version coherence: exactly one `draft.N` value (draft.5).
- **§V check 4** — every file with `IDENTITY name=hello` has `FROM python:3.11-slim` within 3 lines above it; snippets diff clean against `Agentfile.minimal`.
- **§V check 5** — `arc lint` passes for every `examples/Agentfile.*`.
- **§V check 6** — every `_site/**/*.html` has exactly one `rel="canonical"`.
- **§V check 7** — every internal href in `_site/` resolves 200 (full internal link check).
- **§V check 10** — post-deploy sitemap has 0 `CURRENT_IMPLEMENTATION_MAPPING` and 0 `workflow` entries.
- **§V check 11** — `grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ | wc -l` → 0.
- **Per-task Verify** — each of T1–T12's own stated Verify passes locally (T5 sidebar links, T6 inbound links to `/tooling/`, T7/T9/T12 greps, T8 report artifact exists).
- **Release + live** — T13 commits + deploys via PR after local §V passes; T14 re-runs every task's Verify against production with a recorded pass/fail per task ID, all green, plus owner sign-off.

## Connections

Inter-task and build dependencies that constrain ordering:

- **T5 depends on T12** — the sidebar must omit the "Workflow draft" entry, which only exists once T12 parks the workflow draft; sidebar coverage (T5) is downstream of the parking sweep (T12).
- **T12 is broad-reach** — it touches the sidebar, docs index, profiles index, sitemap, `llms.txt`, and the changelog simultaneously; its inbound-ref sweep intersects T5 (sidebar), T6 (docs index), and T7 (sitemap hygiene).
- **T13 depends on all of T1–T12 + local §V** — release cannot proceed until every content/nav/SEO fix is in and the local §V suite (checks 1–7, 10, 11) passes.
- **T14 depends on T13's deploy** — live per-task verification can only run after the PR is merged, GH Pages rebuilds, and the CDN serves the new content.
- **T10 canonical lives in the shared layout/head include** — the two emitters are in the shared Jekyll layout/head, so the dedupe is a single shared-include edit that affects every page at once.
- **Toolchain** — verification requires both the Jekyll build (to produce `_site/` for checks 6, 7, 11) and the Go `arc lint` toolchain under `tooling/`/`cmd/` (for check 5); both must be runnable locally before T13.

## In scope

Tasks T1–T14 exactly as described in the work order:

- T1 ghost-tool purge, T2 hello + FROM, T3 quickstart honesty callouts, T4 polish,
  T5 sidebar + docs-index coverage, T6 de-orphan `/tooling/`, T7 kill duplicate URL,
  T8 lockfile investigate + report only, T9 Runner → Platform relabel,
  T10 canonical dedupe, T11 example-hosts (leave as-is), T12 park workflow draft,
  T13 release via PR + local §V, T14 live per-task verification.

## Out of scope

- **Sprint 2 entirely** — T15–T27 (including the `--backend` / substrate rename work) until Sprint 1's FINAL-GATE + live verification pass and the owner says go.
- **The T8 lockfile A/B DECISION** — Sprint 1 is investigate + report only; the owner decides Option A vs B at Sprint 2 planning. Do NOT change site content or fabricate a format.
- **§H owner/infra tasks** — CDN purge for `/spec/`, `agentrc.io` DNS 301 → agentrc.ai, GitHub/name-collision handling.
- **T11 host rename** — leave `registry.agentrc.io` as-is (aspirational own-brand); no change unless the owner explicitly says switch.
- **Resolving the frontend-image "published" wording** — the discrepancy (work order says "not yet published"; the image was actually published this session) is flagged for the owner, not resolved here (see Open questions #1).

## Open questions

Genuinely blocking ambiguities for the owner — not papered over:

1. **T3 published-vs-unpublished frontend image (BLOCKING).** T3's mandated wording says the agentrc frontend image is "not yet published to a public registry", but this session actually published `ghcr.io/adeelahmad/agentrc-frontend` at release v0.1.0. The work order flags T3's text as authoritative for Sprint 1 wording, yet publishing a doc claiming "not yet published" about an image that IS published re-creates exactly the stale-status drift (memory M-002) this sprint exists to kill. Which is authoritative — the work-order wording, or ground truth? If ground truth, T3's callout text needs owner-approved rewording before it ships.
2. **T1 default option confirmation.** T1 defaults to Option A (`file_read --agentrc-schema`). Confirm the intended fix is Option A and not an Option B variant, since the choice is baked into every affected surface.
3. **T9 relabel-only vs new slug.** Default is relabel visible text only, keeping `/profiles/runner-conformance/`. Confirm the owner does not want a new `/profiles/platform/` (or similar) slug — that would change the URL and require redirect handling.
4. **T7 redirect mechanism on Jekyll/GitHub Pages.** How should `/notes/CURRENT_IMPLEMENTATION_MAPPING/` be killed — `jekyll-redirect-from` plugin (301), a meta-refresh page, or `noindex` + sitemap removal? GH Pages has no server-side redirect config, so the mechanism must be one Jekyll/Pages actually supports, and the choice affects T14's live curl verification.
5. **T12 unpublish mechanism on GH Pages.** How should the workflow pages be unpublished — front-matter `published: false`, or move sources into a non-built `parked/` directory? The generator's behavior determines whether the sources still build (and whether `_site/` cleanly excludes them for §V checks 10/11) while remaining in git.
