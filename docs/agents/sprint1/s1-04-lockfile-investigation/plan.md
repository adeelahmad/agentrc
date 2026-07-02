---
type: plan
story: S1-04
scope: "tests only"
---
# S1-04 Test Plan — Lockfile investigation report (RED)

Tests only. T8 produces a report artifact, not code; each check is a function in
`scripts/verify-sprint1.sh` asserting the report's existence/shape and the no-site-edit
boundary. Every bullet states input/action/assertion and the CURRENT failing (RED) state.

## T8 — investigate + report only

- [ ] `scripts/verify-sprint1.sh::t8_report_exists` — assert
      `docs/agents/sprint1/lockfile-report.md` exists. RED now: FAILS — report not yet written.
- [ ] `scripts/verify-sprint1.sh::t8_report_documents_output` — assert the report contains the
      `arc lock` output facts (filename, format, records) grepped case-insensitively. RED now:
      FAILS — file absent.
- [ ] `scripts/verify-sprint1.sh::t8_report_documents_build` — assert the report documents how
      `arc build` consumes the lock output. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint1.sh::t8_report_ab_recommendation` — assert the report contains
      both "Option A" and "Option B" and takes no decision. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint1.sh::t8_no_site_edits` — action: `git status --porcelain` over
      `index.md spec/ cli.md docs/*.md`; assert 0 changed site files (report-only boundary).
      RED now: green; guard against T8 scope creep into site content.
