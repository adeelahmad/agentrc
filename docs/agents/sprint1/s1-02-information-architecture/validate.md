---
type: validate
story: S1-02
---
# S1-02 Validation — Information architecture (T5, T6, T7, T9, T10)

## Pre-flight

- [ ] S1-03 (T12 parking) has landed — the Workflow-draft sidebar link is already removed
      from `_layouts/doc.html` before finalizing the T5 order.
- [ ] Grep-located every target (§0.8): sidebar links, canonical emitters
      (`_includes/head.html:17` + `:36`), Runner-Conformance labels, `/tooling/` inbound gaps.
- [ ] `bundle exec jekyll build` succeeds (produces `_site/` for the canonical + link checks).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | T5 sidebar has Quickstart | `grep -c "/docs/quickstart/" _layouts/doc.html` | `>= 1` |
| 2 | T5 sidebar has Conformance | `grep -c "/docs/conformance/" _layouts/doc.html` | `>= 1` |
| 3 | T5 sidebar has Impl mapping | `grep -c "/docs/implementation-mapping/" _layouts/doc.html` | `>= 1` |
| 4 | T5 no workflows in sidebar | `grep -c "/docs/workflows/" _layouts/doc.html` | `0` |
| 5 | T5 docs-index cards | `grep -c "/docs/conformance/" docs/index.md` and `grep -c "/docs/implementation-mapping/" docs/index.md` | each `>= 1` |
| 6 | T6 /tooling/ inbound | `grep -c "/tooling/" cli.md docs/index.md` | `>= 1` each |
| 7 | T7 sitemap hygiene (local) | `grep -c CURRENT_IMPLEMENTATION_MAPPING _site/sitemap.xml` | `0` |
| 8 | T9 relabel | `grep -rn "Runner Conformance\|Runner conformance" . \| grep -v "docs/agents"` | (no output) |
| 9 | T9 URL kept | `grep -c "permalink: /profiles/runner-conformance/" profiles/runner-conformance.md` | `1` |
| 10 | T10 one canonical | `fail=0; for p in $(find _site -name '*.html'); do [ "$(grep -c 'rel=\"canonical\"' "$p")" -eq 1 ] \|\| { echo "BAD: $p"; fail=1; }; done; [ "$fail" -eq 0 ]` | exit 0, no BAD |
| 11 | internal link check | `htmlproofer ./_site --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https` | exit 0 |
