---
type: validate
story: S1-04
---
# S1-04 Validation — Lockfile investigation report (T8)

## Pre-flight

- [ ] Grep-located the lock implementation under `cmd/agentrc` + `internal/` (CLI is at
      repo-root `./cmd/agentrc`, per standards correction — not `tooling/`).
- [ ] Confirmed T8 is investigate + report only — no site edits, no A/B decision, no
      fabricated format.

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | report exists | `test -f docs/agents/sprint1/lockfile-report.md && echo OK` | `OK` |
| 2 | documents output | `grep -ci "filename\|format\|record\|emit" docs/agents/sprint1/lockfile-report.md` | `>= 1` |
| 3 | documents build consumption | `grep -ci "build" docs/agents/sprint1/lockfile-report.md` | `>= 1` |
| 4 | A/B present | `grep -c "Option A" docs/agents/sprint1/lockfile-report.md` and `grep -c "Option B" docs/agents/sprint1/lockfile-report.md` | each `>= 1` |
| 5 | no decision taken | manual read: report recommends/presents A vs B for the owner, does NOT select one | confirmed |
| 6 | no site edits | `git status --porcelain index.md spec/ cli.md $(ls docs/*.md) \| wc -l` | `0` |
| 7 | no fabricated format | manual read: every format detail cites the actual `arc lock` code, nothing invented | confirmed |
| 8 | version untouched | `grep -rhoE "draft\.[0-9]+" . \| grep -v .git \| sort -u \| wc -l` | `1` |
