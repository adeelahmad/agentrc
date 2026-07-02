---
type: tasks
story: S1-04
---
# S1-04 Tasks — Lockfile investigation report (T8)

Investigate and REPORT ONLY. No site content change, no fabricated format, no A/B decision.

## T8 — Lockfile: INVESTIGATE + REPORT ONLY [P1]

### 1. Locate and read the `arc lock` implementation

Note the standards correction: the CLI is at repo-root `./cmd/agentrc`, NOT `tooling/`
(`tooling/` holds only `README.md`). Grep-locate the lock code:
`grep -rn -i "lockfile\|agentrc.lock\|arc lock\|newLockCmd\|LockCmd" cmd/ internal/`.
Read the lock command in `cmd/agentrc` and the packages it calls under `internal/`
(`internal/agentfile`, `internal/llb`). Also read how `arc build` consumes any lock output.

### 2. Document exactly what `arc lock` emits

Record, from the code (not from assumption):
- The output filename (e.g. `agentrc.lock` or similar).
- The on-disk format (JSON / TOML / text) and its top-level structure.
- The record set (what each entry pins — digests, resources, models, etc.).
- How `arc build` reads it back (if at all) and what changes when it is present/absent.
Cross-reference the homepage slogan that claims a lockfile and the spec's zero lockfile
content (grep: `grep -rn -i "lockfile\|agentrc.lock" index.md spec/index.md cli.md`).

### 3. Write the report artifact

Create `docs/agents/sprint1/lockfile-report.md` containing the findings above plus an
explicit A/B recommendation for the owner (NO decision):
- **Option A** — a spec subsection under §9 documenting only what the tooling does, marked
  "Status: informative in this draft; format TODO"; slogan kept.
- **Option B** — cut the slogan sentence; the lockfile lives on `/cli/` only.
State the trade-offs of each and stop; the owner decides at Sprint 2 planning.

Constraints: DO NOT fabricate a lockfile format. DO NOT change any site content in Sprint 1.
DO NOT pick A or B. No version bump.

Verify: `docs/agents/sprint1/lockfile-report.md` exists, documents the real `arc lock`
output and `build` consumption, and ends with an A/B recommendation with no decision taken;
`git status` shows zero site-content (`index.md`, `spec/`, `docs/*` non-agents, `cli.md`)
changes attributable to T8.
