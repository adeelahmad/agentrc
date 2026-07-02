---
type: tasks
story: S1-05
---
# S1-05 Tasks — Release & live verification (T13, T14)

Supervisor-executed, not worker-TDD. Runs only after S1-01…S1-04 are in and local §V is
green. Deploy = merge the sprint branch via PR (master is protected: PR required,
squash-merge, 0 approvals) — no direct push.

## T13 — RELEASE Sprint 1 [P0]

1. Local build: `bundle exec jekyll build` succeeds.
2. Run the full §V suite (checks 1–7, 10*, 11) via `scripts/verify-sprint1.sh` — all green.
   (*Check 10 sitemap-hygiene is the local `_site/sitemap.xml` form pre-deploy; the live
   curl form runs in T14.) Example-lint (§V.5) via `go run ./cmd/agentrc lint <f>` for every
   `examples/Agentfile.*`.
3. Commit with a task-tagged message referencing T1–T12.
4. Open a PR from the sprint branch to `master`; squash-merge it (0 approvals required).
5. Wait for the GH Pages build (`pages.yml`) + CDN to serve the new content before T14.

Verify: local §V green; PR merged to `master`; Pages build succeeded; new content live.

## T14 — LIVE per-task verification (against production) [P0]

Re-run each task's Verify against `https://agentrc.ai` via `scripts/verify-sprint1-live.sh`;
record pass/fail per task ID in the sprint report. Live checks include:
- T1: live pages contain no `tools/ping`.
- T2: live hello renders carry `FROM python:3.11-slim`.
- T3: quickstart honesty callouts present on the live page.
- T5/T6: sidebar order live; `/tooling/` reachable from `/cli/` and docs index.
- T7: `curl -s https://agentrc.ai/sitemap.xml | grep -c CURRENT_IMPLEMENTATION_MAPPING` → 0.
- T9: no "Runner Conformance" visible label on live pages.
- T12: `curl -s https://agentrc.ai/sitemap.xml | grep -c workflow` → 0; `/docs/workflows/`
  and `/profiles/workflow-draft/` return 404 (or noindex per chosen mechanism).

Any FAIL → fix within Sprint 1 before FINAL-GATE. Sprint 1 done = all green + owner sign-off.
