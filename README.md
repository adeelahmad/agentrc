<div align="center">

# [agentrc](https://agentrc.ai) <sup>ai</sup>

**Agent Run Config** — an open, Dockerfile-shaped specification for declaring, packaging, securing, and sharing portable AI agents.

[**agentrc.ai**](https://agentrc.ai) · [Docs](https://agentrc.ai/docs/) · [Specification](https://agentrc.ai/spec/) · [Discord](https://discord.gg/jWx6Qak5D) · [GitHub](https://github.com/adeelahmad/agentrc)

> ⚠️ **Working Draft (0.1.0-draft.5).** agentrc is an evolving specification draft, not a finished standard. Expect breaking changes.

</div>

---

## What is agentrc?

`agentrc` (Agent Run Config — like `bashrc`/`zshrc`, but for an agent) is the contract an agent declares: what it is, how it starts, what resources it requests, and how those requests are governed. The **Agentfile** is Dockerfile-shaped — four new keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) over the standard Dockerfile keywords you already know. The build emits namespaced `org.agentrc.*` OCI labels; the **platform** reads those labels — never the Agentfile — and decides what to honour.

agentrc is **not** a runtime, sandbox, cloud platform, model provider, or agent framework. It is the neutral declaration, packaging, and governance layer above all of those.

A single reviewable `Agentfile` declares:

- **identity, capabilities, and objective** — who the agent is (`IDENTITY`), what modalities it supports (`CAPABILITY`), and its system prompt / standard operating procedure (`SOP`, embedded at `/mnt/SOP`);
- **tools, skills, and MCP servers** — files placed into the `/mnt` tree with `COPY` (local) or `ADD --remote` (remote); the destination path determines the resource type;
- **typed requests** — `POLICY` lines asking the platform for a model, network egress, or agent / substrate constraints (`model.*`, `network`, `agent.*`, `substrate.*`);
- **secrets** — **deferred**; this draft defines no secret keyword/schema and leaves credential resolution to the platform;
- **packaging** — the build translates intent into `org.agentrc.*` [OCI](https://opencontainers.org/) labels; the platform reads the labels and **grants, narrows, or rejects** each request. Enforcement is **[Cedar](https://www.cedarpolicy.com/), platform-side** (deny-by-default), and Cedar is never written in the Agentfile.

## 📖 Read the docs

Everything lives on **[agentrc.ai](https://agentrc.ai)**:

| | |
|---|---|
| [What is agentrc?](https://agentrc.ai/docs/what-is-agentrc/) | The problem, the need, and what it solves |
| [Quickstart](https://agentrc.ai/docs/quickstart/) | Write, build, read the labels, push, and run your first Agentfile |
| [Specification](https://agentrc.ai/spec/) | The full working draft |
| [Core profile](https://agentrc.ai/profiles/core/) | Parsing the four keywords + Dockerfile keywords into labels |
| [Security](https://agentrc.ai/docs/security/) | Label-based boundaries and platform-side Cedar enforcement |
| [Conformance](https://agentrc.ai/docs/conformance/) | Profiles and the adversarial test-suite outline |
| [CLI](https://agentrc.ai/cli/) | The BuildKit frontend and the `agentrc` / `arc` CLI |
| [Acknowledgements](https://agentrc.ai/acknowledgements/) | The open standards agentrc builds on |

For LLMs: a machine-readable index is published at [`agentrc.ai/llms.txt`](https://agentrc.ai/llms.txt), and every page is available as raw Markdown via its "View Markdown" link.

## Example

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM ghcr.io/acme/pii-redacted-base:1.4

IDENTITY name=code-reviewer version=1.0 author=acme
CAPABILITY text
CAPABILITY streaming

SOP <<EOF
Review pull requests for correctness and security.
Escalate anything ambiguous to a human.
EOF

CMD claude --print

COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
ADD --remote --cached --fail-if-unavailable \
    https://registry.agentrc.io/skills/code-review:1.2.3 /mnt/skills/code-review

POLICY model.name      claude-opus-4
POLICY model.fallback  claude-sonnet-4
POLICY agent.tool_timeout 30s
POLICY network dns:api.github.com:443
```

Build it with the BuildKit frontend (`docker build -f Agentfile -t ghcr.io/acme/code-reviewer:1.0 .`) or the native CLI (`arc build -t ghcr.io/acme/code-reviewer:1.0 .`) — both produce an identical OCI artifact whose `org.agentrc.*` labels carry every request above.

## This repository

This repo is the agentrc specification **and** the source for the [agentrc.ai](https://agentrc.ai) website (a Jekyll site published via GitHub Pages). It contains the spec, profiles, JSON schemas, grammar, examples, and docs, plus a Go reference implementation under [`tooling/`](tooling/README.md): the `agentrc` CLI and a BuildKit frontend for Agentfiles. The reference implementation is a test harness, not the definition — agentrc the standard is independent of any one implementation. See [Conformance](https://agentrc.ai/docs/conformance/) for exactly which profiles it passes today.

Local preview and publishing notes are in [`README_DEV.md`](README_DEV.md).

## Acknowledgements

agentrc builds on [Agent SOP](https://github.com/strands-agents/agent-sop) (an influence on the `SOP` keyword), [microsandbox](https://github.com/superradcompany/microsandbox) (a reference for the deferred secrets design), [UTCP](https://github.com/universal-tool-calling-protocol/utcp-specification), [MCP](https://github.com/modelcontextprotocol) (servers projected under `/mnt/mcp/`), [Cedar](https://github.com/cedar-policy/cedar) (the platform-side enforcement engine and compilation target for typed `POLICY` requests, never an author surface), [Agent Skills](https://github.com/agentskills/agentskills) (`SKILL.md` bundles under `/mnt/skills/`), [OCI](https://github.com/opencontainers) (artifact + `org.agentrc.*` labels), [Sigstore](https://github.com/sigstore), [SLSA](https://github.com/slsa-framework), [OpenTelemetry](https://github.com/open-telemetry), and [A2A](https://github.com/a2aproject/A2A) (a reference point for the deferred multi-agent workflow companion; the agent-to-agent protocol is out of scope this version). See [Acknowledgements](https://agentrc.ai/acknowledgements/).

## License

[Apache License 2.0](LICENSE). Maintained by [Adeel Ahmad](https://www.linkedin.com/in/adeelahmadch).
