---
type: validate
story: S2-01
---
# S2-01 Validation — Examples expansion (T15, T16)

## Pre-flight

- [ ] Grep-located the examples-index `## Files` card block before editing
      (`examples/index.md:19-22`) and the existing example scaffold/healthcheck style
      (§0.9, M-003).
- [ ] Confirmed both files will use SPEC-DEFINED POLICY keys only (`agent.hooks.*`,
      `agent.sub_agents*`, `agent.sub_agent_timeout`, `network`) — no invented key /
      no fourth namespace (§0.3).
- [ ] Version unchanged: `grep -rhoE 'draft\.[0-9]+' . | grep -v .git | sort -u` → single
      `draft.5` (§0.1).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | T15 lints | `go run ./cmd/agentrc lint examples/Agentfile.hooked` | exit 0 (`: ok`) |
| 2 | T15 hook keys | `grep -cE 'POLICY agent\.hooks\.(on_tool_call\|pre\|post)' examples/Agentfile.hooked` | `>= 3` |
| 3 | T15 explicit egress | `grep -c 'POLICY network dns:' examples/Agentfile.hooked` | `>= 1` |
| 4 | T15 auto-derive comment | `grep -c '\.source' examples/Agentfile.hooked` | `>= 1` |
| 5 | T15 healthcheck | `grep -c 'CMD /mnt/tools/file_read --agentrc-schema' examples/Agentfile.hooked` | `1` |
| 6 | T16 lints | `go run ./cmd/agentrc lint examples/Agentfile.delegator` | exit 0 (`: ok`) |
| 7 | T16 sub-agent keys | `grep -cE 'POLICY agent\.sub_agents(\.max)? \|POLICY agent\.sub_agent_timeout ' examples/Agentfile.delegator` | `>= 3` |
| 8 | T16 healthcheck | `grep -c 'CMD /mnt/tools/file_read --agentrc-schema' examples/Agentfile.delegator` | `1` |
| 9 | index cards | `grep -cE 'Agentfile\.(hooked\|delegator)' examples/index.md` | `>= 2` |
| 10 | all examples lint | `for f in examples/Agentfile.*; do go run ./cmd/agentrc lint "$f" \|\| exit 1; done` | all `: ok` |
| 11 | version unchanged | `grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u \| wc -l` | `1` (value `draft.5`) |
| 12 | site build (CI) | `bundle exec jekyll build --trace` | exit 0 |
</content>
