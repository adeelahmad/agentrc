---
type: plan
story: S1-02
scope: "tests only"
---
# S1-02 Test Plan — Information architecture (RED)

Tests only. Each check is a function in `scripts/verify-sprint1.sh`. Every bullet states
input/action/assertion and the CURRENT failing (RED) state.

## T5 — sidebar + docs-index coverage

- [ ] `scripts/verify-sprint1.sh::t5_sidebar_quickstart` — assert `_layouts/doc.html` links
      `/docs/quickstart/`. RED now: FAILS — sidebar (lines 10–18) has no Quickstart link.
- [ ] `scripts/verify-sprint1.sh::t5_sidebar_conformance` — assert `_layouts/doc.html` links
      `/docs/conformance/`. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint1.sh::t5_sidebar_impl_mapping` — assert `_layouts/doc.html` links
      `/docs/implementation-mapping/`. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint1.sh::t5_sidebar_no_workflows` — assert `_layouts/doc.html` count
      of `/docs/workflows/` == 0. RED now: FAILS — line 16 links Workflow draft.
- [ ] `scripts/verify-sprint1.sh::t5_docsindex_cards` — assert `docs/index.md` has both a
      `/docs/conformance/` card and a `/docs/implementation-mapping/` card. RED now: FAILS —
      card grid ends after Runners (line ~19); neither card exists.

## T6 — de-orphan /tooling/

- [ ] `scripts/verify-sprint1.sh::t6_tooling_inbound` — assert `/tooling/` is linked as an
      anchor from `cli.md` and `docs/index.md`. RED now: FAILS — only prose mentions of
      `tooling/README.md` exist (cli.md:22,34,156); no anchor to the `/tooling/` page.

## T7 — kill duplicate URL

- [ ] `scripts/verify-sprint1.sh::t7_sitemap_no_notes` — action: build then
      `grep -c CURRENT_IMPLEMENTATION_MAPPING _site/sitemap.xml`; assert == 0. RED now: FAILS —
      `notes/CURRENT_IMPLEMENTATION_MAPPING.md` has no front-matter, so jekyll-sitemap includes
      `/notes/CURRENT_IMPLEMENTATION_MAPPING/`.

## T9 — Runner → Platform relabel

- [ ] `scripts/verify-sprint1.sh::t9_no_runner_conformance_label` — action:
      `grep -rn "Runner Conformance\|Runner conformance"` excl. `docs/agents/`; assert count
      == 0. RED now: FAILS — visible-label hits at `profiles/index.md:20`,
      `profiles/runner-conformance.md:3-4`, `llms.txt:40`.
- [ ] `scripts/verify-sprint1.sh::t9_url_preserved` — assert
      `permalink: /profiles/runner-conformance/` still present in
      `profiles/runner-conformance.md`. RED now: green; guard that the URL is NOT changed.

## T10 — canonical dedupe

- [ ] `scripts/verify-sprint1.sh::t10_one_canonical` — action: build then for every
      `_site/**/*.html` assert exactly one `rel="canonical"`. RED now: FAILS — two emitters
      (`_includes/head.html:17` explicit + `{% seo %}` at `:36`) produce two per page.

## Suite gates (also run here)

- [ ] `scripts/verify-sprint1.sh::v7_internal_links` — action: htmlproofer over `_site`;
      assert exit 0 (no dangling internal href). RED now: must stay green after nav edits
      (regression guard against orphan links).
