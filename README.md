<div align="center">

# [agentrc](https://agentrc.ai) <sup>ai</sup>

**Agent Run Config** — an open specification for declaring, packaging, securing, and sharing portable AI agents.

[**agentrc.ai**](https://agentrc.ai) · [Docs](https://agentrc.ai/docs/) · [Specification](https://agentrc.ai/spec/) · [Discord](https://discord.gg/jWx6Qak5D) · [GitHub](https://github.com/adeelahmad/agentrc)

> ⚠️ **Working Draft (v0.1).** agentrc is an evolving specification draft, not a finished standard. Expect breaking changes.

</div>

---

## What is agentrc?

`agentrc` (Agent Run Config — like `bashrc`/`zshrc`, but for an agent) is the contract an agent declares: what it is, how it starts, what it may touch, and how those boundaries are governed. Compatible **runners** decide how to execute that contract on their own substrate.

agentrc is **not** a runtime, sandbox, cloud platform, model provider, or agent framework. It is the neutral declaration, packaging, and governance layer above all of those.

A single reviewable `Agentfile` declares:

- **identity & entrypoint** — what the agent is and how it starts;
- **capabilities** — tools ([UTCP](https://www.utcp.io/)), functions, skills ([SKILL.md](https://agentskills.io/)), and [MCP](https://modelcontextprotocol.io/) servers;
- **instructions** — embedded [Agent SOP](https://github.com/strands-agents/agent-sop) operating procedures (`SOP … END`);
- **boundaries** — files, hosts, and host-scoped [secrets](https://docs.microsandbox.dev/sandboxes/secrets);
- **governance** — [Cedar](https://www.cedarpolicy.com/) policy and audit requirements;
- **packaging** — how the package is pinned, signed, shared, and inspected as an [OCI](https://opencontainers.org/) artifact.

## 📖 Read the docs

Everything lives on **[agentrc.ai](https://agentrc.ai)**:

| | |
|---|---|
| [What is agentrc?](https://agentrc.ai/docs/what-is-agentrc/) | The problem, the need, and what it solves |
| [Quickstart](https://agentrc.ai/docs/quickstart/) | Write, validate, build, and publish your first Agentfile |
| [Specification](https://agentrc.ai/spec/) | The full working draft |
| [Core profile](https://agentrc.ai/profiles/core/) | The minimal normative directive set |
| [Security](https://agentrc.ai/docs/security/) | Declarative, fail-closed boundaries |
| [Conformance](https://agentrc.ai/docs/conformance/) | Profiles and the adversarial test-suite outline |
| [CLI](https://agentrc.ai/cli/) | The `agentrc` CLI (coming soon) |
| [Acknowledgements](https://agentrc.ai/acknowledgements/) | The open standards agentrc builds on |

For LLMs: a machine-readable index is published at [`agentrc.ai/llms.txt`](https://agentrc.ai/llms.txt), and every page is available as raw Markdown via its "View Markdown" link.

## Example

```dockerfile
# syntax=agentrc.agentfile/v0.1

AGENT code-reviewer
CMD claude --print

TOOL utcp:file_read
MOUNT /workspace rw
CRED github_token env:GITHUB_TOKEN host:api.github.com
AUDIT all

SOP code-review
  ## Steps
  ### 1. Triage
  - You MUST read the full diff before commenting.
END

POLICY
  permit(
    principal == AgentRC::Agent::"code-reviewer",
    action == AgentRC::Action::"tool.invoke",
    resource == AgentRC::Tool::"file_read"
  );
END
```

## This repository

This repo is the agentrc specification **and** the source for the [agentrc.ai](https://agentrc.ai) website (a Jekyll site published via GitHub Pages). It contains the spec, profiles, JSON schemas, grammar, examples, and docs. The reference implementation (the `aio-*` packages) is a separate test harness — agentrc the standard is independent of any one implementation.

Local preview and publishing notes are in [`README_DEV.md`](README_DEV.md).

## Acknowledgements

agentrc builds on [Agent SOP](https://github.com/strands-agents/agent-sop), [microsandbox](https://github.com/superradcompany/microsandbox), [UTCP](https://github.com/universal-tool-calling-protocol/utcp-specification), [MCP](https://github.com/modelcontextprotocol), [Cedar](https://github.com/cedar-policy/cedar), [Agent Skills](https://github.com/agentskills/agentskills), [OCI](https://github.com/opencontainers), [Sigstore](https://github.com/sigstore), [SLSA](https://github.com/slsa-framework), [OpenTelemetry](https://github.com/open-telemetry), and [A2A](https://github.com/a2aproject/A2A). See [Acknowledgements](https://agentrc.ai/acknowledgements/).

## License

[Apache License 2.0](LICENSE). Maintained by [Adeel Ahmad](https://www.linkedin.com/in/adeelahmadch).
