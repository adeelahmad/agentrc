---
type: memory
updated_at: 2026-07-03T00:00:00Z
last_retro_sprint: 0
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
