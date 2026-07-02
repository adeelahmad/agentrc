---
type: plan
story: S1-05
scope: "tests only"
---
# S1-05 Test Plan — Release & live verification (checks)

Tests only. This story is supervisor-executed; its "tests" are the full §V suite (local)
plus the T14 live curl checks, expressed as check handles. `scripts/verify-sprint1.sh` (local)
and `scripts/verify-sprint1-live.sh` (live) are created by the RED worker. Every bullet states
input/action/assertion and the CURRENT (RED) state.

## T13 — local §V + release

- [ ] `scripts/verify-sprint1.sh::v_all_local` — action: run the whole local §V suite (checks
      1–7, 10 local, 11); assert every check PASS, exit 0. RED now: FAILS — S1-01/02/03 edits
      not yet applied.
- [ ] `scripts/verify-sprint1.sh::v1_ghost_tool` — assert 0 `tools/ping` (excl `.git`,
      `docs/agents/`). RED now: FAILS (6 surfaces).
- [ ] `scripts/verify-sprint1.sh::v3_version_coherence` — assert exactly one `draft.N`, it is
      `draft.5`. RED now: green; release guard against accidental bump (§0.1).
- [ ] `scripts/verify-sprint1.sh::v4_hello_from` — assert every hello render has
      `FROM python:3.11-slim`. RED now: FAILS.
- [ ] `scripts/verify-sprint1.sh::v5_example_lint` — assert `arc lint` passes every
      `examples/Agentfile.*`. RED now: green; must stay green through release.
- [ ] `scripts/verify-sprint1.sh::v6_one_canonical` — assert one `rel="canonical"` per
      `_site` page. RED now: FAILS (two emitters).
- [ ] `scripts/verify-sprint1.sh::v7_internal_links` — assert htmlproofer exit 0 over `_site`.
      RED now: must be green pre-release (no dangling links after parking + nav edits).
- [ ] `scripts/verify-sprint1.sh::v11_workflow_parked` — assert 0
      `workflow-draft|docs/workflows|agent-workflow.yaml` in `_site/`. RED now: FAILS.
- [ ] `scripts/verify-sprint1.sh::t13_pr_merged` — action: `gh pr list --state merged
      --base master`; assert the sprint PR is merged (deploy via PR, not direct push). RED now:
      FAILS — not yet released.

## T14 — live per-task verification

- [ ] `scripts/verify-sprint1-live.sh::t14_live_ghost` — action: curl live homepage +
      examples; assert 0 `tools/ping`. RED now: FAILS — pre-deploy live still has ghost.
- [ ] `scripts/verify-sprint1-live.sh::t14_live_sitemap_notes` — action:
      `curl -s https://agentrc.ai/sitemap.xml | grep -c CURRENT_IMPLEMENTATION_MAPPING`;
      assert 0. RED now: FAILS.
- [ ] `scripts/verify-sprint1-live.sh::t14_live_sitemap_workflow` — action:
      `curl -s https://agentrc.ai/sitemap.xml | grep -c workflow`; assert 0. RED now: FAILS.
- [ ] `scripts/verify-sprint1-live.sh::t14_live_workflow_404` — action: curl
      `https://agentrc.ai/docs/workflows/`; assert 404 (or noindex per chosen mechanism). RED
      now: FAILS — page live.
- [ ] `scripts/verify-sprint1-live.sh::t14_live_runner_label` — action: curl
      `https://agentrc.ai/profiles/`; assert 0 "Runner Conformance". RED now: FAILS.
- [ ] `scripts/verify-sprint1-live.sh::t14_all_tasks_recorded` — assert every task T1–T12 has
      a recorded live pass/fail in the sprint report. RED now: FAILS — not yet run.
