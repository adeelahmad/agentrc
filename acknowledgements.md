---
layout: doc
title: Acknowledgements
description: "Open standards and projects the agentrc Agentfile (0.1.0-draft.5) builds on and credits."
permalink: /acknowledgements/
---
# Acknowledgements

agentrc is deliberately a thin declaration, packaging, and governance layer over proven open standards rather than a reinvention of them. The Agentfile is Dockerfile-shaped — four new keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) over standard Dockerfile keywords — and it compiles to an OCI artifact carrying `org.agentrc.*` labels. It owes a direct debt to the projects below — listed, where possible, by their source repositories. agentrc declares and governs these; it does not replace any of them.

## Standards agentrc builds on

| Project | What agentrc adopts from it | Source |
|---|---|---|
| **Agent SOP** (Strands Agents) | An influence on the `SOP` keyword: a natural-language, RFC-2119-constrained operating procedure authored in the Agentfile and embedded as a readable file at `/mnt/SOP` (the label is a pointer + digest, never the full text). | [github.com/strands-agents/agent-sop](https://github.com/strands-agents/agent-sop) |
| **microsandbox** | A reference for the **deferred** secrets design (host-scoped credential resolution). Secrets are out of scope for this draft — no `SECRET`/`CRED` keyword and no label schema yet. | [github.com/superradcompany/microsandbox](https://github.com/superradcompany/microsandbox) |
| **Model Context Protocol (MCP)** | The open protocol for model/tool context. MCP servers are projected under `/mnt/mcp/`; agentrc declares and governs MCP, it does not replace it. | [github.com/modelcontextprotocol](https://github.com/modelcontextprotocol) · [modelcontextprotocol.io](https://modelcontextprotocol.io/) |
| **Cedar** (AWS) | The platform-side **enforcement engine and compilation target** for typed `POLICY` requests (never an Agentfile author surface); the `AgentRC::` / Cedar authorization vocabulary. The platform compiles granted requests plus its own org rules into Cedar and enforces deny-by-default, `forbid` over `permit`, and tightening-only `FROM` composition. | [github.com/cedar-policy/cedar](https://github.com/cedar-policy/cedar) · [cedarpolicy.com](https://www.cedarpolicy.com/) |
| **Agent Skills** | The open `SKILL.md` skill-bundle format. Skill bundles are projected under `/mnt/skills/`. | [github.com/agentskills/agentskills](https://github.com/agentskills/agentskills) · [agentskills.io](https://agentskills.io/) |

## Composes with

agentrc references these standards for packaging, signing, provenance, observability, and multi-agent orchestration. They are not required by the core, but conforming platforms and registries are expected to lean on them.

| Project | Role in agentrc | Source |
|---|---|---|
| **Open Container Initiative (OCI)** | The artifact format: standard layers carry the `/mnt` resources and the image config carries the `org.agentrc.*` labels; `application/vnd.agentrc.*` media types. | [github.com/opencontainers](https://github.com/opencontainers) · [opencontainers.org](https://opencontainers.org/) |
| **Sigstore** | Package signing and verification. | [github.com/sigstore](https://github.com/sigstore) · [sigstore.dev](https://www.sigstore.dev/) |
| **SLSA** | Supply-chain provenance for packages. | [github.com/slsa-framework](https://github.com/slsa-framework) · [slsa.dev](https://slsa.dev/) |
| **OpenTelemetry** | Audit and telemetry sinks reached via `POLICY agent.hooks.*` and `agent.telemetry_sink` requests. | [github.com/open-telemetry](https://github.com/open-telemetry) · [opentelemetry.io](https://opentelemetry.io/) |
| **A2A (Agent2Agent)** | A reference point for the **deferred** multi-agent workflow companion. The agent-to-agent protocol is out of scope for this version; capability *exposure* via `IDENTITY` / `CAPABILITY` / labels is in scope. | [github.com/a2aproject/A2A](https://github.com/a2aproject/A2A) |

## Naming

**agentrc** is the specification. **AIO** is the reference implementation and test harness. The relationship is the same as ECMAScript to V8, or OCI to runc — the standard is separate from any one implementation of it.
