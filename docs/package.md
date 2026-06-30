---
layout: doc
title: Package model
description: "An agentrc agent is a standard OCI artifact: ordinary layers carry the /mnt resources, and the image config carries the org.agentrc.* labels."
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
2. **The image config carries `org.agentrc.*` labels** — the machine-readable
   manifest of identity, capabilities, requested model, requested network,
   resource digests, secrets-as-references, and operational requests.

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
| `org.agentrc.*` labels | image config | identity, capabilities, SOP pointer + digest, resource digests + `.origin`, secrets-as-references, and `POLICY`-derived requests |
| Provenance | registry / attestation | digests, [Sigstore](https://www.sigstore.dev/) signatures, [SLSA](https://slsa.dev/) build attestation |

Resources marked `--runtime` are **not** embedded: they appear only as a
reference label (`org.agentrc.mcp.<name>=runtime:<url>`) and are fetched when the
agent bootstraps on the target network.

> **Optional:** when the artifact is built with `--policy-mode digest`, the
> compiled request set is embedded as a structured manifest layer and the labels
> carry its digest instead of inline values. `--policy-mode inline` keeps the
> values in the labels directly. Both forms are supported; which is the build
> default is [open](/spec/). See [§9.4 of the spec](/spec/).

## Registry operations

```bash
# Build a signed OCI artifact, emitting org.agentrc.* labels and embedding
# --cached resources. Two front doors, identical artifacts:
docker build -f Agentfile -t ghcr.io/org/code-reviewer:1.0 .   # BuildKit frontend
arc build      -t ghcr.io/org/code-reviewer:1.0 .              # native CLI

# Push / pull through any OCI registry, like a container image.
arc push  ghcr.io/org/code-reviewer:1.0
arc pull  ghcr.io/org/code-reviewer:1.0
```

The frontend (`# syntax=agentrc.io/agentfile:v1`) and the native `arc` CLI
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
| Who / what is this agent? | `org.agentrc.identity.name`, `…identity.version`, `…identity.author`, `…identity.description` |
| What modalities does it support? | `org.agentrc.capability.*` (e.g. `…capability.streaming=true`) |
| What is its objective? | `org.agentrc.sop=/mnt/SOP` + `org.agentrc.sop.sha256=<digest>` (pointer + digest, never the prompt text) |
| Which model does it request? | `org.agentrc.model.name`, `…model.min_context`, `…model.capability.*`, `…model.fallback` |
| Which hosts does it want to reach? | `org.agentrc.network.dns.<host>=<port>` (auto-derived egress also carries `.source`) |
| Which tools / skills / MCP servers? | `org.agentrc.tool.<name>`, `…skill.<name>`, `…mcp.<name>` (`=local`, `=<digest>`, or `=runtime:<url>`) with `.origin` |
| Which secrets are required? | `org.agentrc.secret.<name>=<scope>` — **references only**, never values |
| Can it spawn sub-agents, and how many? | `org.agentrc.agent.sub_agents`, `…agent.sub_agents.max` |
| What operational limits does it request? | `org.agentrc.agent.idle_timeout`, `…agent.tool_timeout`, `org.agentrc.substrate.runtime.memory`, … |

Every value here is a **request**, not a grant. The platform decides what to
honour and enforces the decision with Cedar (deny-by-default, platform-side);
see [docs/security](/docs/security/) and the
[platform conformance profile](/profiles/runner-conformance/).

## Secrets are references, not values

A secret never enters the artifact. The Agentfile names it and its host scope:

```text
org.agentrc.secret.github_token=host:api.github.com
```

The value is resolved by the platform's secret broker at run time and injected,
so the artifact is portable and safe to mirror and sign. A registry maintainer
auditing an agent can confirm there is **no plaintext secret** anywhere in the
layers, labels, or config — only host-scoped references.

## Override without rebuilding

Every embedded tool / skill / MCP server records both its resolved **digest**
and an **`.origin`** reference, for example:

```text
org.agentrc.mcp.github=sha256:abc123...
org.agentrc.mcp.github.origin=https://registry.agentrc.io/mcp/github:latest
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
  `org.agentrc.*` label namespace, media types, and policy-encoding modes.
- [Specification](/spec/) — the build-time translation tables (§9).
- [Quickstart](/docs/quickstart/) — build, push, and read the labels end to end.
