---
type: validate
story: S1-03
---
# S1-03 Validation — Park the workflow draft (T12)

## Pre-flight

- [ ] Grep-located every workflow inbound ref (M-003): sidebar `_layouts/doc.html:16`,
      `examples/index.md:23,56-64`, `docs/conformance.md:34`, `spec/index.md:823`,
      `docs/non-goals.md:48`, `llms.txt:31,41,45`.
- [ ] Confirmed sources will be KEPT in git (parked, not deleted).
- [ ] `bundle exec jekyll build` succeeds after the parking edits.

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | pages unpublished | `test ! -d _site/docs/workflows && test ! -d _site/profiles/workflow-draft && echo OK` | `OK` |
| 2 | no served refs | `grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ \| wc -l` | `0` |
| 3 | yaml dropped from site | `test ! -f _site/examples/agent-workflow.yaml && echo OK` | `OK` |
| 4 | sitemap clean (local) | `grep -cE 'workflow-draft\|docs/workflows' _site/sitemap.xml` | `0` |
| 5 | llms.txt swept | `grep -c 'docs/workflows/\|profiles/workflow-draft/' llms.txt` | `0` |
| 6 | sidebar unlinked | `grep -c '/docs/workflows/' _layouts/doc.html` | `0` |
| 7 | changelog line | `grep -c "Workflow draft parked" CHANGELOG.md` | `1` |
| 8 | no new version heading | `grep -rhoE "draft\.[0-9]+" . \| grep -v .git \| sort -u \| wc -l` | `1` (still `draft.5`) |
| 9 | sources kept | `git ls-files \| grep -E 'workflows\|workflow-draft\|agent-workflow.yaml' \| wc -l` | `>= 3` |
| 10 | internal link check | `htmlproofer ./_site --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https` | exit 0 |
