---
type: tasks
story: S1-02
---
# S1-02 Tasks — Information architecture (T5, T6, T7, T9, T10)

Runs after S1-03 parks the workflow draft (the sidebar must not list workflows). All
targets grep-located.

## T5 — Sidebar + docs-index coverage [P1]

Edit the sidebar in `_layouts/doc.html` (current links at lines 10–18) to this order,
with NO "Workflow draft" entry (removed by T12 in S1-03):
What is agentrc? · Quickstart · Specification · Agentfile · Security · Package · Runners ·
Conformance · Implementation mapping · CLI · Acknowledgements.
- Add missing links: `/docs/quickstart/`, `/docs/conformance/`,
  `/docs/implementation-mapping/`.
- Add an Examples entry (`/examples/`) if the sidebar pattern allows.
Add two cards to `docs/index.md` (card grid ends after the Runners card at line ~19):
- Conformance card → `/docs/conformance/`, one sentence on profile-based conformance +
  adversarial suite.
- Implementation-mapping card → `/docs/implementation-mapping/`.

Verify: every built page's sidebar links `/docs/quickstart/` and `/docs/conformance/`; no
sidebar link to `/docs/workflows/`; docs index has both new cards; internal link check clean.

## T6 — De-orphan `/tooling/` [P1]

`/tooling/` has zero inbound links today. Add links:
- From `cli.md` — a "Reference implementation available" + "Try it now" link to
  `/tooling/` (the page already references `tooling/README.md` in prose at lines 22, 34, 156;
  add an actual anchor to the tooling page).
- From `docs/index.md` — a card or inline link to `/tooling/`.

Verify: `/tooling/` is reachable from `/cli/` and the docs index; internal link check clean.

## T7 — Kill the duplicate URL [P1]

`/notes/CURRENT_IMPLEMENTATION_MAPPING/` (`notes/CURRENT_IMPLEMENTATION_MAPPING.md`, currently
NO front-matter) duplicates `/docs/implementation-mapping/`. Default mechanism: add
front-matter with `sitemap: false` and a `noindex` robots meta (mechanism is open question #4
— `jekyll-redirect-from` 301 vs meta-refresh vs noindex; confirm with owner if a hard 301 is
wanted). Remove the URL from the sitemap.

Verify (post-deploy, T14):
`curl -s https://agentrc.ai/sitemap.xml | grep -c CURRENT_IMPLEMENTATION_MAPPING` → 0.
Local: the built `_site/sitemap.xml` has 0 `CURRENT_IMPLEMENTATION_MAPPING`.

## T9 — "Runner" → "Platform" naming [P1]

Relabel visible "Runner Conformance" text to "Platform Conformance"; KEEP the
`/profiles/runner-conformance/` URL (open question #3 confirms relabel-only default).
Grep-located visible-label hits:
- `profiles/index.md:20` — `[Runner Conformance](/profiles/runner-conformance/)` link text.
- `profiles/runner-conformance.md:3` — `title: Runner Conformance`.
- `profiles/runner-conformance.md:4` — `description: "Runner Conformance profile."`.
- `llms.txt:40` — `[Runner Conformance](...)` link text.
Do NOT change the permalink `/profiles/runner-conformance/` (line 5).

Verify: `grep -rn "Runner Conformance\|Runner conformance" . | grep -v "docs/agents"` → 0
visible-label hits.

## T10 — Deduplicate canonical tags [P2]

Two emitters produce `<link rel="canonical">` per page:
- `_includes/head.html:17` — explicit `<link rel="canonical" href="{{ page.url | absolute_url }}">`
- `_includes/head.html:36` — `{% seo title=false %}` (jekyll-seo-tag also emits a canonical).
Keep exactly one correct per-page canonical. Default: remove the explicit line 17 and let
`{% seo %}` emit the canonical (it produces the correct absolute per-page URL); OR keep line
17 and suppress seo-tag's canonical — pick the one that yields the correct absolute per-page
URL and delete the duplicate. Do not delete the emitter that produces the correct URL.

Verify (local build): every `_site/**/*.html` has exactly one `rel="canonical"`.
