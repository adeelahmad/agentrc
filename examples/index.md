---
layout: doc
title: Examples
description: "Example agentrc Agentfiles (0.1.0-draft.5) and the deferred workflow companion draft."
permalink: /examples/
---
# Examples

These are concrete, copy-pasteable Agentfiles in the v0.1 model: a
Dockerfile-shaped recipe with four new keywords (`IDENTITY`, `CAPABILITY`,
`SOP`, `POLICY`) over standard Dockerfile keywords. Each builds with either
`docker build -f Agentfile .` (via the BuildKit frontend) or `arc build .` —
both emit identical OCI artifacts with `org.agentrc.*` labels that the platform
reads to grant, narrow, or reject each request. See the [specification](/spec/)
for the full keyword reference.

## Files

- [Minimal Agentfile](/examples/Agentfile.minimal) — the smallest useful agent.
- [Secure workspace Agentfile](/examples/Agentfile.secure-workspace) — a locked-down agent with tight `POLICY` requests.
- [Code reviewer Agentfile](/examples/Agentfile.code-reviewer) — tools, a skill, and an MCP server projected under `/mnt`.
- [Vault agent Agentfile](/examples/Agentfile.vault-agent) — needs a database credential, but credential resolution is **deferred** (platform-defined).
- [Workflow draft YAML](/examples/agent-workflow.yaml) — a **deferred, non-normative** companion (see below).

## Minimal

```dockerfile
# syntax=agentrc.agentfile/v0.1
IDENTITY name=hello version=0.1 author=acme
IDENTITY description="Minimal AgentRC agent"
CAPABILITY text
SOP You are a minimal example agent. Read a file when asked; do nothing else.
CMD python ./agent.py

# Tool (local, embedded) — projected under /mnt/tools/
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# Model + operational requests (platform grants, narrows, or rejects)
POLICY model.name         claude-sonnet-4
POLICY agent.tool_timeout 30s

# Network egress request
POLICY network dns:api.example.com:443

HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/ping
```

This compiles to labels such as `org.agentrc.identity.name=hello`,
`org.agentrc.capability.text=true`, `org.agentrc.tool.file_read=local`,
`org.agentrc.model.name=claude-sonnet-4`, and
`org.agentrc.network.dns.api.example.com=443`. The `SOP` is embedded as a
readable file at `/mnt/SOP` and recorded as a pointer plus digest
(`org.agentrc.sop=/mnt/SOP`), never inlined into a label. The platform reads
those labels — not the Agentfile — when it decides what to honour.

## The workflow companion is deferred

`agent-workflow.yaml` sketches **external** multi-agent orchestration that
references packaged agents by **digest**. It is a deferred, non-normative
companion: the agent-to-agent (A2A) *protocol* — discovery, delegation, and the
cross-agent governance algebra — is out of scope for this version. Capability
*exposure* via `IDENTITY` / `CAPABILITY` / labels is in scope; the workflow draft
is not part of the Agentfile core. See the
[workflow draft profile](/profiles/workflow-draft/).
