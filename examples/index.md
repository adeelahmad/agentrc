---
layout: doc
title: Examples
description: "Example agentrc Agentfiles (v1) and the deferred workflow companion draft."
permalink: /examples/
---
# Examples

These are concrete, copy-pasteable Agentfiles in the v1 model: a
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
- [Vault agent Agentfile](/examples/Agentfile.vault-agent) — secrets as references resolved by the platform broker.
- [Workflow draft YAML](/examples/agent-workflow.yaml) — a **deferred, non-normative** companion (see below).

## Minimal

```dockerfile
# syntax=agentrc.io/agentfile:v1
# Minimal agentrc Agentfile: a Dockerfile-shaped recipe for one local agent.
# Build with `docker build -f Agentfile .` or `arc build .`; both emit identical org.agentrc.* labels.

IDENTITY name=hello version=1.0 author=you
CAPABILITY text
SOP You are a concise local assistant. Answer in one paragraph.
CMD python ./agent.py
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
POLICY model.name claude-opus-4
```

This compiles to labels such as `org.agentrc.identity.name=hello`,
`org.agentrc.capability.text=true`, `org.agentrc.tool.file_read=local`, and
`org.agentrc.model.name=claude-opus-4`. The `SOP` is embedded as a readable file
at `/mnt/SOP` and recorded as a pointer plus digest
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
