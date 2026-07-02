---
type: tasks
story: S2-01
---
# S2-01 Tasks — Examples expansion (T15, T16)

Both files use SPEC-DEFINED POLICY keys ONLY. Full scaffold, verbatim Sprint-1
healthcheck line, header comment. Grep-locate the examples-index card block before
editing (M-003). No version change here — stays `0.1.0-draft.5` (§0.1).

## T15 — `examples/Agentfile.hooked` [P2]

Create `examples/Agentfile.hooked` demonstrating hook policy + auto-derived egress:
- `# syntax=agentrc.agentfile/v0.1` line, then `FROM python:3.11-slim`.
- `IDENTITY name=… version=0.1 author=acme` + `IDENTITY description="…"`; `CAPABILITY`;
  `SOP …`; `CMD python ./agent.py`; `COPY --chmod=755 ./tools/file_read /mnt/tools/file_read`.
- `POLICY agent.hooks.on_tool_call <https endpoint>` plus one `POLICY agent.hooks.pre <https>`
  and one `POLICY agent.hooks.post <https>` (spec-defined `agent.hooks.*` keys only).
- Comment block explaining the platform AUTO-DERIVES the hook endpoint's `network` egress
  grant and records it with `.source` attribution (requested vs derived), per §8.5 /
  `/docs/security/` (e.g. `org.agentrc.network.dns.hooks.internal.source=auto:agent.hooks.pre`).
- One EXPLICIT `POLICY network dns:<host>:<port>` line for contrast (a requested, not
  derived, egress).
- `HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema`
  (verbatim Sprint-1 pattern).

## T16 — `examples/Agentfile.delegator` [P2]

Create `examples/Agentfile.delegator` demonstrating sub-agent policy:
- Same scaffold (syntax line, `FROM python:3.11-slim`, IDENTITY×2, CAPABILITY, SOP, CMD,
  `COPY … file_read`, healthcheck line).
- `POLICY agent.sub_agents true`, `POLICY agent.sub_agents.max <n>`,
  `POLICY agent.sub_agent_timeout <dur>` (spec-defined keys only).
- Comment: sub-agent grants are the platform's call (requests, not guarantees — §0.4).

## Examples index cards (both) — `examples/index.md`

- Add two bullets to the `## Files` list in the existing card style (grep-located at
  `examples/index.md:19-22`), each linking `/examples/Agentfile.hooked` and
  `/examples/Agentfile.delegator` with a one-line description matching the existing tone.
- Optional: one index line showing `arc build --policy-mode digest` — NO default-mode prose
  (§14.2 #1: the default is undecided; do not imply one).

Verify: `go run ./cmd/agentrc lint examples/Agentfile.hooked` and `… Agentfile.delegator`
both pass; index links resolve; `grep -rhoE 'draft\.[0-9]+' .` still single `draft.5`.
</content>
