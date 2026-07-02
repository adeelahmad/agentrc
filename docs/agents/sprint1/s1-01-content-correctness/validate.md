---
type: validate
story: S1-01
---
# S1-01 Validation — Content correctness (T1–T4)

## Pre-flight

- [ ] Grep-located every target before editing (§0.8, M-003): `tools/ping` and
      `IDENTITY name=hello` surfaces confirmed on the six/five files listed in `tasks.md`.
- [ ] `examples/Agentfile.minimal` confirmed as the single canonical hello source (§0.6).
- [ ] Version unchanged: `grep -rq "draft\.5" .` and no second `draft.N` introduced (§0.1).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | T1 ghost tool | `grep -rn "tools/ping" . \| grep -v .git \| grep -v "docs/agents" \| wc -l` | `0` |
| 2 | T1 schema mirror | `grep -rn "file_read --agentrc-schema" index.md spec/index.md examples/index.md examples/Agentfile.minimal examples/Agentfile.code-reviewer \| wc -l` | `6` |
| 3 | T2 hello+FROM | `for f in $(grep -rln "IDENTITY name=hello" --include="*.md" --include="*.html" . \| grep -v "docs/agents"); do grep -q "FROM python:3.11-slim" "$f" \|\| echo "MISSING: $f"; done` | (no output) |
| 4 | T2 spec sentence | `grep -c "every Agentfile MUST contain a \`FROM\` instruction" spec/index.md` | `>= 1` |
| 5 | T2 diff clean | inline hello blocks in `index.md`/`docs/quickstart.md`/`examples/index.md` diff clean vs `examples/Agentfile.minimal` | no drift (M-001) |
| 6 | T3 callouts | `grep -c "not yet published to a public registry" docs/quickstart.md` and `grep -c "\`arc run\` is planned" docs/quickstart.md` | each `>= 1` |
| 7 | T3 no stale present-tense | `grep -rn "nothing to install\|no extra tooling" docs/ \| grep -v "will need"` | (no present-tense hit) |
| 8 | T4 brand casing | `grep -rn "Minimal AgentRC agent" .` | (no output) |
| 9 | T4 minimal desc | `grep -c "Minimal agentrc agent" examples/Agentfile.minimal` | `1` |
| 10 | example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" \|\| exit 1; done` | all `: ok` |
| 11 | site build | `bundle exec jekyll build` | exit 0 |
