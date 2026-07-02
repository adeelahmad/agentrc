---
type: plan
story: S1-01
scope: "tests only"
---
# S1-01 Test Plan — Content correctness (RED)

Tests only. Each check is a function in `scripts/verify-sprint1.sh` (created by the RED
worker; GREEN makes the edits pass it). Every bullet states input/action/assertion and the
CURRENT failing (RED) state.

## T1 — ghost tool purge

- [ ] `scripts/verify-sprint1.sh::t1_no_ghost_ping` — action: `grep -rn "tools/ping"`
      excluding `.git` and `docs/agents/`; assert count == 0. RED now: `tools/ping` present on
      6 site surfaces (`index.md:38`, `spec/index.md:636`, `spec/index.md:716`,
      `examples/Agentfile.minimal:23`, `examples/Agentfile.code-reviewer:52`,
      `examples/index.md:45`).
- [ ] `scripts/verify-sprint1.sh::t1_schema_healthcheck_present` — assert
      `grep -c "file_read --agentrc-schema"` across the six surfaces == 6. RED now: 0
      (replacement not yet made).

## T2 — hello + FROM

- [ ] `scripts/verify-sprint1.sh::t2_hello_has_from` — for each file matching
      `IDENTITY name=hello` (excl. `docs/agents/`), assert `FROM python:3.11-slim` present
      within 3 lines above it. RED now: FAILS — inline hello renders lack `FROM`.
- [ ] `scripts/verify-sprint1.sh::t2_spec_from_sentence` — assert the FROM-required MUST
      sentence exists in `spec/index.md`. RED now: FAILS — sentence absent.
- [ ] `scripts/verify-sprint1.sh::t2_hello_diff_clean` — extract each inline hello block and
      diff against `examples/Agentfile.minimal`; assert no drift (M-001). RED now: FAILS —
      inline snippets diverge from source.

## T3 — quickstart honesty callouts

- [ ] `scripts/verify-sprint1.sh::t3_step2_callout` — assert `docs/quickstart.md` contains
      "not yet published to a public registry". RED now: FAILS — callout absent.
- [ ] `scripts/verify-sprint1.sh::t3_step5_callout` — assert `docs/quickstart.md` contains
      "`arc run` is planned". RED now: FAILS — callout absent.
- [ ] `scripts/verify-sprint1.sh::t3_no_present_tense_claims` — assert no present-tense
      "nothing to install"/"no extra tooling" guarantee remains. RED now: FAILS —
      `docs/quickstart.md:71` present-tense claim exists.

## T4 — polish

- [ ] `scripts/verify-sprint1.sh::t4_brand_casing` — assert
      `grep -rn "Minimal AgentRC agent" .` count == 0. RED now: FAILS —
      `examples/Agentfile.minimal:8` says "Minimal AgentRC agent".
- [ ] `scripts/verify-sprint1.sh::t4_quickstart_hello_canonical` — assert the quickstart hello
      carries `version=0.1`, `author=acme`, `model.name claude-sonnet-4`. RED now: FAILS —
      fields not yet aligned to canonical.

## Suite gates (also run here)

- [ ] `scripts/verify-sprint1.sh::v5_example_lint` — assert
      `go run ./cmd/agentrc lint <f>` passes for every `examples/Agentfile.*`. RED now: passes
      today for existing files but MUST stay green after edits (regression guard).
- [ ] `scripts/verify-sprint1.sh::v3_version_coherence` — assert exactly one `draft.N` value
      and it is `draft.5`. RED now: green; guard against accidental version bump (§0.1).
