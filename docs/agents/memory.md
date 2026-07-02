---
type: memory
updated_at: 2026-07-03T00:00:00Z
last_retro_sprint: 1
---
# Memory — agentrc

Cross-sprint learning, curated at each retrospective. Entries are advisory and
never relax an invariant. First agentic-agile session in this repo (Sprint 1);
seeded from the 2026-07-02/03 docs+examples+spec audit that produced this work order.

## Entries

| id | applies_to | text | origin | recurrence | added |
|---|---|---|---|---|---|
| M-001 | all | Inline doc/HTML snippets drift from the source-of-truth file they mirror (e.g. the "hello" Agentfile rendered on homepage/quickstart/examples-index diverged from `examples/Agentfile.minimal`). When a snippet duplicates a canonical file, its acceptance check MUST diff against that file, not just grep for a substring. | 2026-07-02 audit; T1/T2/T4 | 3 | 2026-07-03 |
| M-002 | all | Status/claim prose goes stale silently ("not yet published", "coming soon", present-tense "nothing to install") after the underlying fact changes. Re-derive every status claim from current ground truth (CLI code, published artifacts) rather than carrying prior wording forward. | 2026-07-02 audit; T3 | 2 | 2026-07-03 |
| M-003 | all | Locate every edit target by grep FIRST — repo layout is not assumed; the same string (e.g. `tools/ping`, `IDENTITY name=hello`) recurs across .md, .html, and Agentfile surfaces, and a partial sweep leaves the site internally inconsistent. | work-order §0.9; T1/T12 | 2 | 2026-07-03 |
| M-004 | all | On a static-site repo (Jekyll), the pipeline's own `docs/agents/` planning artifacts carry YAML front matter and WILL be built into the site (22 junk pages + sitemap URLs leaked to prod). Add the pipeline artifact dir to the generator's `exclude:` at sprint start, and grep `_site`/live sitemap for it in the release checks. | sprint1 T14 defect; PR #17 | 2 | 2026-07-03 |
| M-005 | all | Never run supervisor git / working-tree commands (checkout, pull, append+commit) while a worker is active in the SHARED working tree — their branch/index ops collide (a supervisor `git checkout` failed mid-worker). Serialized-branch execution means: dispatch, WAIT for return, THEN do git. | sprint1 wave-3 | 1 | 2026-07-03 |
| M-006 | all | Live checks that grep rendered HTML for a multi-token literal (e.g. `FROM python:3.11-slim`) false-negative because syntax highlighting splits it across `<span>`s. Verify source for exact strings; for live, strip tags or check tokens (`>FROM<` + `python:3.11-slim`). | sprint1 T14 | 1 | 2026-07-03 |
