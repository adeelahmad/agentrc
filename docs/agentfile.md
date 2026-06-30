---
layout: doc
title: Agentfile
description: "The Dockerfile-shaped declaration for one agent: four new keywords over standard Dockerfile keywords."
permalink: /docs/agentfile/
---
# Agentfile

The Agentfile is a **Dockerfile-shaped declaration for one agent.** It reuses
standard Dockerfile keywords for everything it can and adds the **minimum**
number of new keywords — just four. The mental model, file shape, and much of
the tooling transfer directly from Docker.

The Agentfile describes one agent and the contract a platform must satisfy
before running it. It does not describe a cluster, fleet, or workflow.

## The four new keywords

| Keyword | Role in agentrc |
|---|---|
| `IDENTITY` | Who / what the agent is: name, version, author, description (`key=value` pairs). |
| `CAPABILITY` | What modalities / patterns the agent supports (text, streaming, function-calling, multimodal, …), one per line. |
| `SOP` | The agent's system prompt / objective / standard operating procedure (inline, heredoc, or file-backed at `/mnt/SOP`). |
| `POLICY` | A typed, namespaced **request** for a resource, model, or operational constraint. |

> There is **no** `TOOL`, `MCP`, `SERVER`, `FUNCTION`, `CRED`, `MOUNT`,
> `MEMORY`, or `RATELIMIT` keyword. Tools, skills, and MCP servers are added with
> `COPY` / `ADD`; memory, context, CPU, and model are `POLICY` requests.

## Standard Dockerfile keywords agentrc uses

Everything else is a standard Dockerfile keyword, used as-is or with a
documented extension.

| Keyword | Role in agentrc |
|---|---|
| `FROM` | Base image **and** agent inheritance (`FROM another-agent`). |
| `CMD` | The agent's invocation surface — the framework / loop to run. agentrc is framework-neutral. |
| `COPY` | Add **local** tools, skills, MCP bundles, or an SOP file into the `/mnt` tree. |
| `ADD` | Add **remote** resources via `--remote` plus delivery flags (`--cached` / `--runtime`, `--fail-if-unavailable` / `--warn-if-unavailable`). |
| `HEALTHCHECK` | Liveness probe; MAY invoke a projected tool. |
| `LABEL` | Standard OCI metadata; available for hand-authored `org.agentrc.*` metadata. |
| `ENV` / `ARG` / `WORKDIR` / `USER` / `EXPOSE` / `RUN` | Standard Dockerfile semantics, unchanged. |

## The `/mnt` projection layout

The agent's projected filesystem lives under **`/mnt`**. The **destination path**
under `/mnt` determines what a resource is — there is no dedicated keyword per
resource type.

```text
/mnt
├── tools/     # each tool is an executable; argv in, structured output out
├── skills/    # skill bundles (instructions / scripts / resources)
├── mcp/       # MCP server bundles or configs, projected as tool directories
├── proc/      # runtime-populated: live policy, identity, budgets, audit tail
└── SOP        # the agent's system prompt, as a readable file
```

A **tool** is an executable file under `/mnt/tools/`. To be self-describing, it
exposes its schema via `--agentrc-schema` (JSON to stdout) **or** ships a sibling
`<tool>.toolspec.json`.

## Minimal pattern

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim

IDENTITY name=hello version=1.0 author=you
CAPABILITY text
SOP You are a concise local assistant. Answer in one paragraph.
CMD python ./agent.py

# A local tool, projected at /mnt/tools/file_read
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# Typed requests the platform grants, narrows, or rejects
POLICY model.name claude-opus-4
POLICY network dns:api.github.com:443
```

At build time the compiler translates this into namespaced OCI labels under
`org.agentrc.*` — for example `org.agentrc.identity.name=hello`,
`org.agentrc.model.name=claude-opus-4`, and
`org.agentrc.network.dns.api.github.com=443`. The platform reads **labels**, not
the Agentfile.

## Why only four keywords?

- **Dockerfile-shaped.** Authors already know `FROM`, `CMD`, `COPY`, `ADD`,
  `LABEL`, `ENV`, and `HEALTHCHECK`. Reusing them keeps the file familiar and the
  tooling transferable.
- **Minimum new surface.** The four new keywords cover exactly what Dockerfile
  cannot express about an agent: identity, capabilities, the system prompt, and
  typed requests. Resources are just files under `/mnt`.
- **The frontend interprets them.** The four agentrc keywords and the
  `ADD --remote` extension are interpreted by the **agentrc BuildKit frontend**
  (`# syntax=agentrc.agentfile/v0.1`) or the `agentrc` / `arc` CLI; both produce
  identical OCI artifacts. A stock `docker build` without the frontend
  understands only the standard Dockerfile keywords.

## Where to go next

- The [Specification](/spec/) is the single source of truth — full grammar, the
  `POLICY` namespaces, build-time label translation, and runtime behaviour.
- The [Core profile](/profiles/core/) defines how a conformant compiler parses
  the Agentfile and emits `org.agentrc.*` labels.
