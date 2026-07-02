---
type: tasks
story: S2-07
---
# S2-07 Tasks — Release & live verification (T27)

Supervisor-executed (not worker-TDD). Serialize all git ops — never run supervisor git in a
shared tree while a worker is active (M-005). Version = `draft.6` sitewide by now.

## T27 — RELEASE Sprint 2 + live verification [P0]

1. **Full §V (local where feasible + CI):** run `scripts/verify-sprint2.sh` — all checks incl.
   3 (version coherence = single `draft.6`), 8 (terminology split), 9 (backends). Delegate the
   full `go build ./...` / `go test -race ./...`, site build, link check, and one-canonical
   checks to CI per the gate matrix (ENOSPC; open question #2 UNRESOLVED).
2. **Confirm pipeline artifacts excluded from the site (M-004):** `docs/agents/` must be in
   the generator `exclude:`; grep `_site`/live sitemap for `docs/agents` → 0.
3. **Commit + PR:** task-tagged commits; open + squash-merge a PR to `master` (protected: PR
   required). No direct push.
4. **Deploy:** wait for GH Pages build + CDN.
5. **Live re-verify each Sprint-2 task** against `https://agentrc.ai` via
   `scripts/verify-sprint2-live.sh` — version = draft.6 sitewide; `--backend` in live CLI docs;
   new examples lint + render; §8.7/§8.8/§8.9 rendered; addendum committed iff owner-provided
   (open question #1 UNRESOLVED). Account for syntax-highlight token splitting in HTML greps
   (M-006). Record per-task pass/fail per task ID; any FAIL fixed within Sprint 2 before FINAL-GATE.
6. **Owner sign-off** recorded; **retrospective** written to `docs/agents/` (feeds `memory.md`).

Verify: local/CI §V all green (incl. 3/8/9); PR merged to `master`; live per-task verification
recorded all green; owner sign-off; retrospective present.
</content>
