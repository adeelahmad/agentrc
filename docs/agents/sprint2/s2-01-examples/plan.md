---
type: plan
story: S2-01
scope: "tests only"
---
# S2-01 Test Plan — Examples expansion (RED)

Tests only. Each check is a function in `scripts/verify-sprint2.sh` (created by the RED
worker; GREEN makes the edits pass it). Every bullet states input/action/assertion and the
CURRENT failing (RED) state.

## T15 — `examples/Agentfile.hooked`

- [ ] `scripts/verify-sprint2.sh::t15_hooked_exists` — action: assert `examples/Agentfile.hooked`
      exists. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint2.sh::t15_hooked_lints` — action: `go run ./cmd/agentrc lint
      examples/Agentfile.hooked`; assert exit 0. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint2.sh::t15_hooked_hook_keys` — assert `grep -cE 'POLICY
      agent\.hooks\.(on_tool_call|pre|post)'` >= 3. RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t15_hooked_explicit_egress` — assert one explicit
      `POLICY network dns:` line present (contrast to auto-derived). RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t15_hooked_source_attribution` — assert a `.source`
      auto-derivation comment present (§8.5). RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t15_hooked_healthcheck` — assert verbatim
      `CMD /mnt/tools/file_read --agentrc-schema` present. RED now: FAILS.

## T16 — `examples/Agentfile.delegator`

- [ ] `scripts/verify-sprint2.sh::t16_delegator_exists` — assert `examples/Agentfile.delegator`
      exists. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint2.sh::t16_delegator_lints` — action: `go run ./cmd/agentrc lint
      examples/Agentfile.delegator`; assert exit 0. RED now: FAILS — file absent.
- [ ] `scripts/verify-sprint2.sh::t16_delegator_subagent_keys` — assert `agent.sub_agents`,
      `agent.sub_agents.max`, `agent.sub_agent_timeout` all present (>= 3 matches). RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t16_delegator_healthcheck` — assert verbatim
      `CMD /mnt/tools/file_read --agentrc-schema` present. RED now: FAILS.

## Examples index + suite guards

- [ ] `scripts/verify-sprint2.sh::t15_16_index_cards` — assert `examples/index.md` links both
      `Agentfile.hooked` and `Agentfile.delegator` (>= 2 matches). RED now: FAILS — cards absent.
- [ ] `scripts/verify-sprint2.sh::v5_example_lint` — assert `go run ./cmd/agentrc lint <f>`
      passes for EVERY `examples/Agentfile.*` (regression guard incl. new files). RED now:
      passes for existing files, FAILS once new files are referenced but not yet valid.
- [ ] `scripts/verify-sprint2.sh::v3_version_draft5_guard` — assert exactly one `draft.N`
      value == `draft.5` (guard against an accidental bump in this story; §0.1). RED now: green;
      regression guard.
</content>
