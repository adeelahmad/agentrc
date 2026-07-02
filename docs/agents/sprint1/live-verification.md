---
type: report
sprint: 1
task: T14
---
# Sprint 1 — T14 live verification (https://agentrc.ai)

Per-task live pass/fail against production after the PR #16 deploy, plus the
post-hotfix re-verification. Build-dependent checks that could not run locally
(no Jekyll in this env) were delegated here and to CI's htmlproofer.

## Initial live run (post PR #16 deploy)

| Task | Check | Result |
|---|---|---|
| T1 | live homepage/examples/spec have 0 `tools/ping` | PASS |
| T2 | rendered `hello` carries `FROM python:3.11-slim` (present; literal grep false-negatives due to syntax-highlight spans — verified via `>FROM<` + `python:3.11-slim` tokens + source `v4_hello_from`) | PASS |
| T3 | quickstart shows "`arc run` is planned" + published-image callout | PASS |
| T4 | brand casing / hello fields (source `t4_*`) | PASS |
| T5 | built sidebar links `/docs/conformance/`, no `/docs/workflows/` | PASS |
| T6 | `/tooling/` returns 200 and is linked from `/cli/` | PASS |
| T7 | live sitemap has 0 `CURRENT_IMPLEMENTATION_MAPPING` | PASS |
| T9 | live `/profiles/` has 0 "Runner Conformance" (Platform relabel) | PASS |
| T10 | exactly one `rel="canonical"` per page (`/`, `/spec/`, `/docs/`, `/cli/`, `/docs/quickstart/`) | PASS |
| T12 | `/docs/workflows/` and `/profiles/workflow-draft/` return 404 | PASS |
| T12 | live sitemap has 0 `workflow` URLs | **FAIL → remediated** |

### The one failure and its fix
The sitemap listed 3 `.../s1-03-park-workflow-draft/...` URLs — not the parked
pages, but the pipeline's own **internal planning artifacts under `docs/agents/`**,
which Jekyll published because they carry YAML front matter and `docs/agents`
was not in `_config.yml`'s `exclude:` list (22 junk pages total leaked). Fix:
add `docs/agents` to `exclude:` (this change). Re-verified after redeploy below.

## Post-hotfix re-verification (after excluding docs/agents)

_Filled in after this fix deploys._

| Check | Result |
|---|---|
| live sitemap has 0 `docs/agents` URLs | PENDING |
| live sitemap has 0 `workflow` URLs | PENDING |
| all other T14 checks still green | PENDING |
