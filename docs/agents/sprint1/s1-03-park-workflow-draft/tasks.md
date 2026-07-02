---
type: tasks
story: S1-03
---
# S1-03 Tasks — Park the workflow draft (T12)

Unpublish now, keep all sources in git. Must land before S1-02's T5 sidebar is finalized.
All targets grep-located (M-003).

## T12 — PARK the workflow draft [P1]

### 1. Unpublish the two pages (DO NOT delete sources)

- `docs/workflows.md` (permalink `/docs/workflows/`) — front-matter `published: false`
  (or move under a non-built `parked/` area if the generator needs it; mechanism is open
  question #5). Add a header comment "Parked 2026-07-03 — will return in a future draft".
- `profiles/workflow-draft.md` (permalink `/profiles/workflow-draft/`) — same treatment.

### 2. Drop the workflow YAML from the served site

- `examples/agent-workflow.yaml` — remove from the served site (keep the file in-repo under
  `parked/`).
- `examples/index.md:23` — remove the card/link `[Workflow draft YAML](/examples/agent-workflow.yaml)`.

### 3. Sweep inbound refs (unlink where a mention must remain)

- Sidebar: `_layouts/doc.html:16` — remove the `<a href=".../docs/workflows/">Workflow draft</a>`
  link (this also serves T5).
- `examples/index.md` — the "workflow companion is deferred" section (lines ~56–64) and the
  `[workflow draft profile](/profiles/workflow-draft/)` link (line 64): unlink to plain text
  "workflow orchestration is parked for a future draft".
- `docs/conformance.md:34` — the `[Workflow Draft](/profiles/workflow-draft/)` table link:
  unlink (keep the row's prose about the deferred companion; drop the dead link).
- `docs/workflows.md` outbound links from the parked page are irrelevant once unpublished.
- Spec deferred/non-goals prose that LINKS to `/docs/workflows/` or `/profiles/workflow-draft/`
  must be unlinked so it does not 404: `spec/index.md:823` (`[workflow draft](/docs/workflows/)`)
  and `docs/non-goals.md:48` (`[workflow draft](/docs/workflows/)`). Non-linking prose
  mentions of "workflow" may stay (§work-order T12.3).

### 4. Sitemap + llms.txt

- Remove both URLs from the sitemap (achieved by `published: false` / parking — verify
  `_site/sitemap.xml` has no `workflow-draft` or `docs/workflows`).
- `llms.txt:31` (`[Workflows](/docs/workflows/)`), `llms.txt:41`
  (`[Workflow Draft](/profiles/workflow-draft/)`), and `llms.txt:45` ("and a deferred
  workflow YAML") — remove the two link entries and drop the YAML mention.

### 5. Changelog

- `CHANGELOG.md` — add one line under the current `0.1.0-draft.5` section: "Workflow draft
  parked (unpublished); returns in a future revision." Do NOT add a new version heading (§0.1).

Verify: built site has no `/docs/workflows/` or `/profiles/workflow-draft/` page;
`grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' _site/ | wc -l` → 0; internal
link check clean; sources still tracked in git.
