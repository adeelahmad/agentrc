---
layout: doc
title: Package model
description: "An agentrc agent is a standard OCI artifact: ordinary layers carry the /mnt resources, and the image config carries the ai.agentrc.* labels."
permalink: /docs/package/
---
# Package model

An agentrc agent compiles to an **ordinary OCI artifact**. There is no bespoke
package format, no lockfile-as-a-package, and no proprietary registry: the agent
pushes, pulls, signs, and mirrors through any OCI-compatible registry exactly
like a container image. That is the whole point — registry maintainers already
have the tooling.

Two things make the artifact an *agent*:

1. **Layers carry the `/mnt` resources** — the tools, skills, MCP bundles, and
   the SOP file that were embedded at build time.
2. **The image config carries `ai.agentrc.*` labels** — the machine-readable
   manifest of identity, capabilities, requested model, requested network,
   resource digests, and operational requests.

The platform reads the **labels**, never the Agentfile. See the
[specification](/spec/) for the full label translation and the
[OCI labels &amp; package profile](/profiles/oci-package/) for the normative
namespace catalog.

## Artifact contents

| Component | Where it lives | What it carries |
|---|---|---|
| Embedded tools | layer → `/mnt/tools/` | executables added via `COPY` or `ADD --remote --cached` |
| Embedded skills | layer → `/mnt/skills/` | skill bundles (`SKILL.md` + scripts/resources) |
| Embedded MCP bundles | layer → `/mnt/mcp/` | MCP server bundles / configs |
| SOP file | layer → `/mnt/SOP` | the agent's system prompt / objective, as a readable file |
| `ai.agentrc.*` labels | image config | identity, capabilities, SOP pointer + digest, resource digests + `.origin`, and `POLICY`-derived requests |
| Provenance | registry / attestation | digests, [Sigstore](https://www.sigstore.dev/) signatures, [SLSA](https://slsa.dev/) build attestation |

Resources marked `--runtime` are **not** embedded: they appear only as a
reference label (`ai.agentrc.mcp.<name>=runtime:<url>`) and are fetched when the
agent bootstraps on the target network.

> **Optional:** when the artifact is built with `--policy-mode digest`, the
> compiled request set is embedded as a structured manifest layer and the labels
> carry its digest instead of inline values. `--policy-mode inline` keeps the
> values in the labels directly. Both forms are supported; which is the build
> default is [open](/spec/). See [§9.4 of the spec](/spec/).

## Registry operations

```bash
# Build a signed OCI artifact, emitting ai.agentrc.* labels and embedding
# --cached resources. Two front doors, identical artifacts:
docker build -f Agentfile -t ghcr.io/org/code-reviewer:1.0 .   # BuildKit frontend
arc build      -t ghcr.io/org/code-reviewer:1.0 .              # native CLI

# Push / pull through any OCI registry, like a container image.
arc push  ghcr.io/org/code-reviewer:1.0
arc pull  ghcr.io/org/code-reviewer:1.0
```

The frontend (`# syntax=agentrc.agentfile/v0.1`) and the native `arc` CLI
produce **identical** OCI artifacts. See the [CLI page](/cli/).

## Reading the labels (review before run)

Because the manifest is just labels, you can inspect an agent before it ever
runs — with `docker inspect`, `arc inspect`, or any registry's label viewer —
and answer, from the labels alone:

```bash
docker inspect ghcr.io/org/code-reviewer:1.0   # or: arc inspect ...
```

| Question | Labels that answer it |
|---|---|
| Who / what is this agent? | `ai.agentrc.identity.name`, `…identity.version`, `…identity.author`, `…identity.description` |
| What modalities does it support? | `ai.agentrc.capability.*` (e.g. `…capability.streaming=true`) |
| What is its objective? | `ai.agentrc.sop=/mnt/SOP` + `ai.agentrc.sop.sha256=<digest>` (pointer + digest, never the prompt text) |
| Which model does it request? | `ai.agentrc.model.name`, `…model.min_context`, `…model.capability.*`, `…model.fallback` |
| Which hosts does it want to reach? | `ai.agentrc.network.dns.<host>=<port>` (auto-derived egress also carries `.source`) |
| Which tools / skills / MCP servers? | `ai.agentrc.tool.<name>`, `…skill.<name>`, `…mcp.<name>` (`=local`, `=<digest>`, or `=runtime:<url>`) with `.origin` |
| Can it spawn sub-agents, and how many? | `ai.agentrc.agent.sub_agents`, `…agent.sub_agents.max` |
| What operational limits does it request? | `ai.agentrc.agent.idle_timeout`, `…agent.tool_timeout`, `ai.agentrc.substrate.runtime.memory`, … |

Every value here is a **request**, not a grant. The platform decides what to
honour and enforces the decision with Cedar (deny-by-default, platform-side);
see [docs/security](/docs/security/) and the
[platform conformance profile](/profiles/runner-conformance/).

## Secrets

Secrets are **deferred** in this draft — no `SECRET`/`CRED` keyword and no
`ai.agentrc.secret.*` schema; credential resolution is platform-defined and out
of scope for now.

## Override without rebuilding

Every embedded tool / skill / MCP server records both its resolved **digest**
and an **`.origin`** reference, for example:

```text
ai.agentrc.mcp.github=sha256:abc123...
ai.agentrc.mcp.github.origin=https://registry.example.com/mcp/github:latest
```

This lets a platform or organization re-point a public resource to an internal
mirror — by honouring an overridden `…mcp.github.origin` at deploy time —
**without rebuilding** the agent. The digest preserves verifiability; the
`.origin` is the substitution seam.

## Signing, provenance, and reproducibility

Because the agent is a standard OCI artifact, it slots into the existing supply
chain:

- **Signing** — sign with [Sigstore](https://www.sigstore.dev/) / cosign, exactly
  as you would a container image.
- **Provenance** — attach [SLSA](https://slsa.dev/) build attestation linking the
  artifact digest back to its source and builder.
- **Reproducibility** — `--cached` resources are fetched at build and pinned by
  digest, so the same Agentfile and inputs produce the same content-addressed
  artifact.

## Where to go next

- [OCI labels &amp; package profile](/profiles/oci-package/) — the normative
  `ai.agentrc.*` label namespace, media types, and policy-encoding modes.
- [Specification](/spec/) — the build-time translation tables (§9).
- [Quickstart](/docs/quickstart/) — build, push, and read the labels end to end.
