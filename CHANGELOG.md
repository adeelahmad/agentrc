---
layout: doc
title: Changelog
description: "agentrc site and spec changelog."
permalink: /changelog/
---
# Changelog

## 0.1.0-draft.5 — 2026-06-30

The Dockerfile-shaped redesign. The Agentfile is now **Dockerfile-shaped**, the build emits
`org.agentrc.*` OCI labels, and the platform reads those labels — never the
Agentfile — to grant, narrow, or reject each request and enforce it with Cedar.

### Changed (breaking)

- **Redesigned the Agentfile to be Dockerfile-shaped.** There are now exactly
  **four** new keywords — `IDENTITY`, `CAPABILITY`, `SOP`, and `POLICY` — layered
  over standard Dockerfile keywords (`FROM`, `CMD`, `COPY`, `ADD`, `HEALTHCHECK`,
  `LABEL`, `ENV`, `ARG`, `WORKDIR`, `USER`, `EXPOSE`, `RUN`). The mental model,
  file shape, and tooling transfer directly from Docker.
- **Removed the legacy directive family.** `AGENT`, `TOOL`, `TOOLSET`,
  `FUNCTION`, `SKILL`, `SERVER`, `MCP`, `URL`, `CRED`, `BIND`, `MOUNT`, `PLUGIN`,
  `ALLOW`, `DENY`, `RATELIMIT`, `TIMEOUT`, `LIMIT`, `SLICE`, `IMAGE`,
  `ISOLATION`, `BROKER`, `BACKEND`, `TRACE`, `MEMORY`, `OPTIMIZER`, and `SHELL`
  are gone, along with the inline Cedar `POLICY … END` block and the old
  `SOP name … END` block form.
- **Tools, skills, and MCP servers are now files under `/mnt`.** Add local
  resources with `COPY` and remote ones with `ADD --remote` (plus delivery flags
  `--cached`/`--runtime` and `--fail-if-unavailable`/`--warn-if-unavailable`).
  The destination path under `/mnt` (`tools/`, `skills/`, `mcp/`, `SOP`)
  determines the resource type.
- **Secrets are deferred** — removed the `org.agentrc.secret.*` / `LABEL`-secret
  model; no `SECRET`/`CRED` keyword and no secret schema in this draft.
  Credential resolution is platform-defined and out of scope (future design).
- **Removed `AUDIT`** — audit rides on `POLICY agent.hooks.on_tool_call`.
- **Resource, model, network, and lifecycle requests are typed `POLICY` lines.**
  Each `POLICY <namespaced.key> <value>` is a single request in the `agent.*`,
  `substrate.*`, `model.*`, or `network` namespace.
- **`POLICY` is a request, not enforcement.** The build emits `org.agentrc.*`
  OCI labels; the platform reads the labels and **grants, narrows, or rejects**
  each request, with deny-by-default applied to its grant decision.
- **Cedar moved to a platform-side enforcement engine and compilation target.**
  Cedar is no longer an author surface and MUST NOT appear in the Agentfile;
  the platform compiles granted typed requests plus its own organization rules
  into Cedar, with a normative request→Cedar mapping (`NetworkEgress`,
  `tool.invoke`, `mcp.request`, `agent.delegate`, `device.access`) and the
  guarantees `forbid` overrides `permit`, order-independently, and monotonic
  composition across `FROM`.
- **Two build paths, identical artifacts.** The **BuildKit frontend**
  (`# syntax=agentrc.agentfile/v0.1` then `docker build -f Agentfile`) and the
  native **`agentrc` / `arc` CLI** (`build` / `push` / `pull` / `run`) produce
  identical OCI artifacts. Substrate / isolation is a run-time choice
  (`--isolation`, `--substrate`), never an Agentfile directive.

### Deferred

- **A2A (the agent-to-agent protocol)** — Agent Cards, discovery, cross-agent
  delegation, and the governance algebra of an agent-to-agent call — is out of
  scope for this version. Capability *exposure* via `IDENTITY` / `CAPABILITY` /
  labels is in scope; the *protocol* is not.
- Workflow draft parked (unpublished); returns in a future revision.

## Earlier history (superseded by 0.1.0-draft.5)

Before 0.1.0-draft.5, agentrc was published as a series of **0.1.x working drafts** built on
a different, much larger model: roughly thirty directives (`AGENT`, `TOOL`,
`CRED`, `MOUNT`, and many more), an inline Cedar policy block authored directly
in the Agentfile, and an `/agentrc` tool-projection root. **0.1.0-draft.5 replaces that
model entirely** — see the breaking changes above. Those drafts also introduced
the work that carried forward in spirit: embedded operating procedures (now the `SOP`
keyword at `/mnt/SOP`), the standards acknowledgements, OCI-based packaging, and
the standards-style site, brand, and theming.

The full, detailed history of the 0.1.x drafts is preserved in the repository's
git log.
