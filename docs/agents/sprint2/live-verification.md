---
type: report
sprint: 2
task: T27
---
# Sprint 2 — T27 live verification (https://agentrc.ai)

Per-task live pass/fail against production after PR #19 (commit b1bf760) deployed
at 0.1.0-draft.6. Build-dependent checks that can't run locally (no Jekyll) were
delegated here + to CI's htmlproofer (both green on PR #19).

## Live checks (all green)

| Task | Check | Result |
|---|---|---|
| T15/T16 | `/examples/Agentfile.hooked` and `/examples/Agentfile.delegator` reachable (200) | PASS |
| T17 | `/spec/` renders §8.7 `substrate.<platform>.*` | PASS |
| T18 | `/spec/` renders §8.8 `agent.auth.*` | PASS |
| T19 | `/spec/` renders §8.9 `substrate.runtime.language` | PASS |
| T20 | `/spec/` shows `0.1.0-draft.6` (×6) and **0** `0.1.0-draft.5`; homepage draft.6 | PASS |
| T20 (T8-A) | §9.5 `agentrc.lock` informative subsection served; homepage slogan kept | PASS |
| T21 | `/cli/` contains `--backend` (×6) and **0** `--substrate` flag | PASS |
| T25 | `/cli/` `run` row = implemented (reference translators); §0.8 line present | PASS |
| T26 | `/docs/demo/` returns 200 with the verbatim three-backend narrative | PASS |
| — | one `rel="canonical"` per page; no ghost `tools/ping`; `/docs/workflows/` 404 | PASS |
| — | sitemap has **0** `docs/agents` and **0** `workflow` URLs | PASS |

`scripts/verify-sprint2-live.sh https://agentrc.ai` → **2 passed / 0 failed** (t27_live,
t27_no_artifact_leak_live).

## One harness defect found and fixed (not a site defect)

`t27_live` initially reported `spec-not-draft6 no-8.7 no-8.9` while raw `curl` proved
all were live. Root cause: `set -o pipefail` + `echo "$spec" | grep -q PATTERN`. `grep -q`
exits on first match and SIGPIPEs the `echo` of a 97 KB page; `pipefail` turns that
SIGPIPE (141) into a non-zero pipeline → a FALSE "not found". Early-position matches
(draft.6/8.7/8.9) triggered it; the lone late `8.8` match did not (echo finished first) —
exactly the observed pattern. `grep -c` checks were immune (they read to EOF). Fix:
replaced the `echo … | grep -q` pipelines with here-strings (`grep -q PAT <<<"$var"`) —
no pipe, no SIGPIPE. Retrospective memory M-007.

**T27 outcome: all live checks green.** Sprint 2 (T15–T27) shipped via PR #19; the
harness fix + this record via a follow-up PR. Site verified correct throughout.
