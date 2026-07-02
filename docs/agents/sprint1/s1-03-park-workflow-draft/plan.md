---
type: plan
story: S1-03
scope: "tests only"
---
# S1-03 Test Plan — Park the workflow draft (RED)

Tests only. Each check is a function in `scripts/verify-sprint1.sh`. Every bullet states
input/action/assertion and the CURRENT failing (RED) state.

## T12 — park workflow draft

- [ ] `scripts/verify-sprint1.sh::t12_no_served_workflow_refs` — action: build then
      `grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ \| wc -l`; assert
      == 0 (§V.11). RED now: FAILS — `/docs/workflows/` and `/profiles/workflow-draft/` build
      and are referenced across the site.
- [ ] `scripts/verify-sprint1.sh::t12_pages_unpublished` — assert `_site/docs/workflows/`
      and `_site/profiles/workflow-draft/` do NOT exist after build. RED now: FAILS — both
      pages currently render.
- [ ] `scripts/verify-sprint1.sh::t12_yaml_dropped` — assert `_site/examples/agent-workflow.yaml`
      does NOT exist and `examples/index.md` has no `/examples/agent-workflow.yaml` link. RED
      now: FAILS — served + linked at `examples/index.md:23`.
- [ ] `scripts/verify-sprint1.sh::t12_sidebar_unlinked` — assert `_layouts/doc.html` has 0
      `/docs/workflows/` links. RED now: FAILS — line 16 links Workflow draft.
- [ ] `scripts/verify-sprint1.sh::t12_llms_swept` — assert `llms.txt` has 0
      `docs/workflows/` and 0 `profiles/workflow-draft/` link entries. RED now: FAILS —
      `llms.txt:31,41` list both.
- [ ] `scripts/verify-sprint1.sh::t12_sitemap_no_workflow` — action: build then
      `grep -cE 'workflow-draft\|docs/workflows' _site/sitemap.xml`; assert == 0. RED now:
      FAILS — both URLs in the generated sitemap.
- [ ] `scripts/verify-sprint1.sh::t12_changelog_line` — assert `CHANGELOG.md` contains
      "Workflow draft parked" under the `0.1.0-draft.5` section. RED now: FAILS — line absent.
- [ ] `scripts/verify-sprint1.sh::t12_sources_kept` — assert the workflow sources are still
      tracked (`git ls-files` shows workflows/workflow-draft/agent-workflow.yaml). RED now:
      green; guard that parking does NOT delete sources.
- [ ] `scripts/verify-sprint1.sh::t12_no_dangling_links` — action: htmlproofer over `_site`;
      assert exit 0 (no unlinked-but-still-referenced 404). RED now: would FAIL after a naive
      unpublish that leaves dangling refs — this guards the sweep is complete (§V.7).
