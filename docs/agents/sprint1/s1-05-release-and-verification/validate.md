---
type: validate
story: S1-05
---
# S1-05 Validation — Release & live verification (T13, T14)

## Pre-flight

- [ ] S1-01, S1-02, S1-03, S1-04 all complete; their story validate.md sign-offs green.
- [ ] Local §V suite green via `scripts/verify-sprint1.sh` (checks 1–7, 10 local, 11).
- [ ] On a sprint branch, NOT `master` (master is protected — deploy is via PR).

## Final sign-off

### T13 — local + release

| # | Check | Command | Expected |
|---|-------|---------|----------|
| 1 | site build | `bundle exec jekyll build` | exit 0 |
| 2 | full §V local | `scripts/verify-sprint1.sh` | all checks PASS, exit 0 |
| 3 | example-lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" \|\| exit 1; done` | all `: ok` |
| 4 | one canonical | `fail=0; for p in $(find _site -name '*.html'); do [ "$(grep -c 'rel=\"canonical\"' "$p")" -eq 1 ] \|\| fail=1; done; [ "$fail" -eq 0 ]` | exit 0 |
| 5 | PR merged | `gh pr list --state merged --base master --limit 1` | shows the sprint PR |
| 6 | no direct push | release landed via squash-merge PR, not a direct push to `master` | confirmed |

### T14 — live per-task (against production)

| # | Task | Command | Expected |
|---|------|---------|----------|
| 7 | T1 live ghost | `curl -s https://agentrc.ai/ https://agentrc.ai/examples/ \| grep -c tools/ping` | `0` |
| 8 | T7 sitemap | `curl -s https://agentrc.ai/sitemap.xml \| grep -c CURRENT_IMPLEMENTATION_MAPPING` | `0` |
| 9 | T12 sitemap | `curl -s https://agentrc.ai/sitemap.xml \| grep -c workflow` | `0` |
| 10 | T12 pages gone | `curl -s -o /dev/null -w '%{http_code}' https://agentrc.ai/docs/workflows/` | `404` (or noindex per chosen mechanism) |
| 11 | T9 live label | `curl -s https://agentrc.ai/profiles/ \| grep -c "Runner Conformance"` | `0` |
| 12 | per-task record | each of T1–T12 re-run live, pass/fail recorded per task ID in the sprint report | all PASS |
| 13 | owner sign-off | owner confirms Sprint 1 done | recorded |
