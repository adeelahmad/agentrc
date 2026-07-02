---
type: validate
story: S2-07
---
# S2-07 Validation — Release & live verification (T27)

## Pre-flight

- [ ] All of S2-01…S2-06 landed; local/CI §V green before opening the PR.
- [ ] `docs/agents/` confirmed in the generator `exclude:` (M-004); `_site` free of pipeline
      artifacts.
- [ ] Git ops serialized — no supervisor git while a worker is active (M-005).
- [ ] Live greps written to tolerate syntax-highlight token splitting (M-006).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | §V local suite | `scripts/verify-sprint2.sh` | all checks exit 0 |
| 2 | version coherence (§V.3) | `grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u` | single `draft.6` |
| 3 | terminology split (§V.8) | `[ "$(rg -l -- '--substrate' cmd/ docs/ cli.md \| wc -l)" -eq 0 ] && rg -q 'POLICY substrate\.' spec/index.md` | pass |
| 4 | backends (§V.9) | `scripts/verify-sprint2.sh::v9_backend_dryruns` (bedrock JSON + k8s YAML) | pass |
| 5 | CI green | GH Actions `go` + `build` jobs on the PR | all green |
| 6 | no artifact leak (M-004) | `grep -rc 'docs/agents' _site/ 2>/dev/null \|\| echo 0` | `0` |
| 7 | PR merged | `gh pr view --json state,baseRefName` | merged to `master` |
| 8 | live version | `scripts/verify-sprint2-live.sh::t27_live` (version = draft.6 sitewide) | pass |
| 9 | live --backend | live `/cli/` shows `--backend`, no `--substrate` (token-split aware, M-006) | pass |
| 10 | live examples render | `https://agentrc.ai/examples/Agentfile.hooked` + `…delegator` reachable + lint | pass |
| 11 | per-task record | per-task pass/fail recorded in a live-verification artifact under `docs/agents/sprint2/` | present, all green |
| 12 | owner sign-off + retro | owner "go" recorded; retrospective written under `docs/agents/` | present |
</content>
