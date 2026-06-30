---
layout: doc
title: Acknowledgements
description: "Open standards and projects agentrc builds on and credits."
permalink: /acknowledgements/
---
# Acknowledgements

agentrc is deliberately a thin declaration, packaging, and governance layer over proven open standards rather than a reinvention of them. It owes a direct debt to the projects below — listed, where possible, by their source repositories. agentrc declares and governs these; it does not replace any of them.

## Standards agentrc builds on

| Project | What agentrc adopts from it | Source |
|---|---|---|
| **Agent SOP** (Strands Agents) | Natural-language, RFC-2119-constrained operating procedures embedded in the Agentfile via `SOP … END` | [github.com/strands-agents/agent-sop](https://github.com/strands-agents/agent-sop) |
| **microsandbox** | The host-scoped, placeholder-substitution **secret** model — the value never enters the sandbox | [github.com/superradcompany/microsandbox](https://github.com/superradcompany/microsandbox) |
| **Universal Tool Calling Protocol (UTCP)** | Tool description/invocation over native endpoints — the `utcp:` tool namespace | [github.com/universal-tool-calling-protocol/utcp-specification](https://github.com/universal-tool-calling-protocol/utcp-specification) · [utcp.io](https://www.utcp.io/) |
| **Model Context Protocol (MCP)** | The open protocol for model/tool context — the `MCP` directive and `mcp:` namespace | [github.com/modelcontextprotocol](https://github.com/modelcontextprotocol) · [modelcontextprotocol.io](https://modelcontextprotocol.io/) |
| **Cedar** (AWS) | The default policy language and `AgentRC::` authorization vocabulary | [github.com/cedar-policy/cedar](https://github.com/cedar-policy/cedar) · [cedarpolicy.com](https://www.cedarpolicy.com/) |
| **Agent Skills** | The open `SKILL.md` skill-bundle format — the `SKILL` directive | [github.com/agentskills/agentskills](https://github.com/agentskills/agentskills) · [agentskills.io](https://agentskills.io/) |

## Composes with

agentrc references these standards for packaging, signing, provenance, observability, and multi-agent orchestration. They are not required by the core, but conforming runners and registries are expected to lean on them.

| Project | Role in agentrc | Source |
|---|---|---|
| **Open Container Initiative (OCI)** | Content-addressed package and registry distribution; `vnd.agentrc.*` media types | [github.com/opencontainers](https://github.com/opencontainers) · [opencontainers.org](https://opencontainers.org/) |
| **Sigstore** | Package signing and verification | [github.com/sigstore](https://github.com/sigstore) · [sigstore.dev](https://www.sigstore.dev/) |
| **SLSA** | Supply-chain provenance for packages | [github.com/slsa-framework](https://github.com/slsa-framework) · [slsa.dev](https://slsa.dev/) |
| **OpenTelemetry** | Trace emission for the `TRACE` directive and audit | [github.com/open-telemetry](https://github.com/open-telemetry) · [opentelemetry.io](https://opentelemetry.io/) |
| **A2A (Agent2Agent)** | A reference point for the future multi-agent workflow companion | [github.com/a2aproject/A2A](https://github.com/a2aproject/A2A) |

## Naming

**agentrc** is the specification. **AIO** is the reference implementation and test harness. The relationship is the same as ECMAScript to V8, or OCI to runc — the standard is separate from any one implementation of it.
