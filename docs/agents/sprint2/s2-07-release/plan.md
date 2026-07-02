---
type: plan
story: S2-07
scope: "tests only"
---
# S2-07 Test Plan — Release & live verification (RED)

Tests only. This story is supervisor-executed; its "tests" are the full §V suite
(`scripts/verify-sprint2.sh`) run pre-deploy and the live re-verification
(`scripts/verify-sprint2-live.sh`) run post-deploy. Every bullet FAILS until the sprint is
merged + deployed.

## Pre-deploy — full §V suite (local where feasible + CI)

- [ ] `scripts/verify-sprint2.sh::v_all_local` — aggregate: runs every t15…t26 check + guards
      (v3/v8/v9) and asserts all pass. RED now: FAILS — most tasks not yet landed.
- [ ] `scripts/verify-sprint2.sh::v3_version_draft6` — assert single `draft.N` == `draft.6`
      (§V.3). RED now: FAILS pre-T20.
- [ ] `scripts/verify-sprint2.sh::v8_terminology_split` — assert `--substrate` gone from
      CLI/docs (→0) and `POLICY substrate.` intact (§V.8). RED now: FAILS pre-rename.
- [ ] `scripts/verify-sprint2.sh::v9_backend_dryruns` — assert bedrock JSON + k8s YAML dry-runs
      parse (§V.9). RED now: FAILS — translators not landed.
- [ ] `scripts/verify-sprint2.sh::v_no_artifact_leak` — assert `_site` (CI) / repo excludes
      `docs/agents/` from the built site (M-004). RED now: guarded; regression check.

## Post-deploy — live re-verification against https://agentrc.ai

- [ ] `scripts/verify-sprint2-live.sh::t27_live` — aggregate live check: version = draft.6
      sitewide; `/cli/` shows `--backend` and no `--substrate` (token-split aware, M-006); the two
      new examples render + lint; §8.7/§8.8/§8.9 rendered; per-task pass/fail recorded. RED now:
      FAILS — not deployed.
- [ ] `scripts/verify-sprint2-live.sh::t27_no_artifact_leak_live` — assert the live sitemap has 0
      `docs/agents` URLs (M-004). RED now: FAILS — not deployed / guard.
</content>
