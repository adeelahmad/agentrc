---
layout: doc
title: OCI Package
description: "OCI Labels & Package profile (0.1.0-draft.5): the org.agentrc.* label namespace, layers, media types, policy encoding, override, and provenance for an agentrc OCI artifact."
permalink: /profiles/oci-package/
---
# OCI Labels &amp; Package Profile

**Version:** 0.1.0-draft.5 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.agentfile/v0.1`)  
**Date:** 2026-06-30  
**Audience:** registry maintainers, platform / runner authors, security &amp; compliance reviewers

> This profile defines the **on-the-wire shape** of a built agentrc agent: the
> `org.agentrc.*` label namespace carried in the OCI image config, the layers
> that carry the `/mnt` resources, the recommended media types and annotations,
> the two policy-encoding modes, and the override / signing / provenance
> contract. It is the registry-side companion to the
> [Specification](/spec/): the compiler emits what is described here; a platform
> reads it back (see the [platform conformance profile](/profiles/runner-conformance/)).

## 1. Purpose

`agentrc build` (or the BuildKit frontend) compiles an Agentfile into an
**ordinary OCI artifact**. There is no bespoke package format: the embedded
resources live in standard layers, and the machine-readable manifest is the set
of `org.agentrc.*` **labels** in the image config. That is the entire contract —
a platform reads the labels, **never the Agentfile source**.

Because the result is a normal OCI artifact, an agent pushes, pulls, signs,
mirrors, and attests through any OCI-compatible registry exactly like a container
image. No proprietary registry, no special storage.

> **Namespace.** The label namespace is **`org.agentrc.*`**. The legacy
> `io.agentrc.*` / `io.agentio.*` namespaces from earlier drafts are **not used
> in 0.1.0-draft.5**; treat any occurrence as stale.

## 2. What a package contains

| Part | Where it lives | Carries |
|---|---|---|
| **Embedded resources** | Standard image **layers** | The `/mnt` tree: `/mnt/tools/*`, `/mnt/skills/*`, `/mnt/mcp/*`, and `/mnt/SOP` for `COPY`'d / `--cached` resources. |
| **Manifest** | Image **config labels** | The full `org.agentrc.*` label set — identity, capabilities, SOP pointer + digest, resource references, and `POLICY`-derived requests. |
| **Resolved policy manifest** *(optional)* | An extra **layer** | Present only under `--policy-mode digest`: a structured request manifest the labels point to by digest (see [§6](#6-policy-encoding-inline-vs-digest)). |

Resources marked `--runtime` are **not** embedded as layers; only their
reference is recorded in labels and they are fetched when the agent bootstraps
([spec §4.3](/spec/)).

## 3. The `org.agentrc.*` label catalog

These are the labels the compiler emits. They are the only thing a platform
reads. All translations follow the spec's build tables ([spec §9](/spec/)).

### 3.1 Identity, capability, SOP

| Source | Label(s) |
|---|---|
| `IDENTITY name=claims-triage version=1.0 author=acme` | `org.agentrc.identity.name=claims-triage`, `org.agentrc.identity.version=1.0`, `org.agentrc.identity.author=acme` |
| `IDENTITY description="…"` | `org.agentrc.identity.description=…` |
| `CAPABILITY streaming` | `org.agentrc.capability.streaming=true` |
| `SOP …` / `COPY ./sop.md /mnt/SOP` | `org.agentrc.sop=/mnt/SOP`, `org.agentrc.sop.sha256=<digest>` |

The SOP label is always a **pointer plus digest**, never the full prompt text —
the prompt is a readable file at `/mnt/SOP` in a layer. (Whether `capability.*`
is one label per value or a single comma-joined label is
[open decision #5](/spec/).)

### 3.2 Resources — tools, skills, MCP servers

The destination path under `/mnt` determines the resource type, and the label
namespace mirrors it (`tool.*`, `skill.*`, `mcp.*`). The value encodes how the
resource was delivered.

| Source | Build behaviour | Label(s) |
|---|---|---|
| `COPY ./tools/x /mnt/tools/x` | Embed as a layer. | `org.agentrc.tool.x=local` |
| `ADD --remote --cached <url> /mnt/tools/x` | Fetch at build, embed as a layer. | `org.agentrc.tool.x=<digest>` + `org.agentrc.tool.x.origin=<url>` |
| `ADD --remote --runtime <url> /mnt/mcp/x` | Reference only; fetch at bootstrap. | `org.agentrc.mcp.x=runtime:<url>` |

The value of a resource label is therefore one of:

- `local` — a `COPY`'d file embedded in a layer;
- `<digest>` (e.g. `sha256:abc123…`) — a remote `--cached` resource fetched at
  build and embedded;
- `runtime:<url>` — a `--runtime` resource recorded by reference only.

Every **embedded** MCP server or skill MUST emit **both** the resolved digest
and an `.origin` reference, so a platform can re-point it without a rebuild:

```text
org.agentrc.mcp.github=sha256:abc123...
org.agentrc.mcp.github.origin=https://registry.agentrc.io/mcp/github:latest
```

Secrets are **deferred** in this draft — there is no `SECRET`/`CRED` keyword and
no `org.agentrc.secret.*` schema; credential resolution is left entirely to the
platform and is out of scope for now.

### 3.3 `POLICY`-derived requests

Each authored `POLICY <key> <value>` becomes `org.agentrc.<key>=<value>`. These
are **requests**, not grants; the platform grants, narrows, or rejects each one
([Enforcement profile](/profiles/security/)).

| Authored | Emitted label |
|---|---|
| `POLICY agent.idle_timeout 5m` | `org.agentrc.agent.idle_timeout=5m` |
| `POLICY agent.context.type autocompressed` | `org.agentrc.agent.context.type=autocompressed` |
| `POLICY agent.sub_agents.max 5` | `org.agentrc.agent.sub_agents.max=5` |
| `POLICY agent.hooks.pre https://hooks.internal/pre-step` | `org.agentrc.agent.hooks.pre=https://hooks.internal/pre-step` (+ auto-derived egress, below) |
| `POLICY substrate.runtime.memory 8gb` | `org.agentrc.substrate.runtime.memory=8gb` |
| `POLICY substrate.ptty true` | `org.agentrc.substrate.ptty=true` |
| `POLICY model.name claude-opus-4` | `org.agentrc.model.name=claude-opus-4` |
| `POLICY model.capability vision` | `org.agentrc.model.capability.vision=true` |
| `POLICY network dns:api.github.com:443` | `org.agentrc.network.dns.api.github.com=443` |

**Auto-derived egress.** A `POLICY` value that is a URL the agent calls out to
(`agent.hooks.*`, `agent.interrupt_endpoint`) auto-derives an explicit,
**attributed** `network` egress label, so it is never a silent network hole:

```text
org.agentrc.network.dns.hooks.internal=443
org.agentrc.network.dns.hooks.internal.source=auto:agent.hooks.pre
```

The `.source` annotation records which request produced the derived egress.
Auto-derivation is convenience, not an implicit grant — the platform must still
grant it.

### 3.4 Namespace summary

| Prefix | Meaning | Value form |
|---|---|---|
| `org.agentrc.identity.*` | Who the agent is. | string (`name`, `version`, `author`, `description`, …) |
| `org.agentrc.capability.*` | Supported modalities / patterns. | `true` |
| `org.agentrc.sop`, `org.agentrc.sop.sha256` | SOP pointer + digest. | `/mnt/SOP`, `<digest>` |
| `org.agentrc.tool.*`, `.skill.*`, `.mcp.*` | Embedded or referenced resources. | `local` \| `<digest>` \| `runtime:<url>` |
| `org.agentrc.tool.*.origin` (etc.) | Origin reference for an embedded resource. | `<url>` |
| `org.agentrc.agent.*` | Agent-side operational requests. | per spec §8.1 |
| `org.agentrc.substrate.*` | Substrate / resource requests. | per spec §8.2 |
| `org.agentrc.model.*` | Model / capability requests. | per spec §8.3 |
| `org.agentrc.network.dns.*` | Egress requests (+ `.source` for auto-derived). | `<port>` |

The namespace is **extensible**: new `agent.*` / `substrate.*` / `model.*` keys
add new labels without changing the grammar.

## 4. Layers

The embedded `/mnt` tree is carried in ordinary layers — there is no special
on-disk format beyond what the resource label already advertises:

```text
layer(s)
└── /mnt
    ├── tools/     # executables for tool.*=local | <digest>
    ├── skills/    # skill bundles (SKILL.md) for skill.*=local | <digest>
    ├── mcp/       # MCP server bundles for mcp.*=local | <digest>
    └── SOP        # the readable system-prompt file (org.agentrc.sop)
```

`runtime:<url>` resources contribute **no** layer; they are fetched at bootstrap.
`/mnt/proc` is **not** packaged — it is populated by the platform at run time
(see the [projection profile](/profiles/tool-projection/)).

## 5. Recommended media types and annotations

A package SHOULD use a standard OCI image manifest and an agentrc config media
type. Layers use standard tar+gzip layer types.

| Component | Media type |
|---|---|
| Manifest | `application/vnd.oci.image.manifest.v1+json` |
| Agent config | `application/vnd.agentrc.agent.config.v1+json` |
| Resource layer (tools / skills / MCP / SOP) | `application/vnd.agentrc.agent.layer.v1.tar+gzip` |
| Resolved policy manifest *(digest mode)* | `application/vnd.agentrc.policy.manifest.v1+json` |

**Annotations.** A package SHOULD carry the standard OCI image annotations
alongside the agentrc identity labels, so generic OCI tooling shows useful
metadata:

```text
org.opencontainers.image.title          # = org.agentrc.identity.name
org.opencontainers.image.version        # = org.agentrc.identity.version
org.opencontainers.image.authors        # = org.agentrc.identity.author
org.opencontainers.image.source
org.opencontainers.image.revision
org.opencontainers.image.created
```

The authoritative agentrc manifest remains the `org.agentrc.*` label set; the
`org.opencontainers.image.*` annotations are conventional metadata, not a
substitute.

## 6. Policy encoding: inline vs digest

The compiled request set MUST be retrievable by a platform in **either** form:

- **Inline** (`--policy-mode inline`): request values live directly in
  `org.agentrc.*` labels, as in the tables above. Best for small request sets;
  fully self-describing from the config alone.
- **By digest** (`--policy-mode digest`): the labels carry a **digest** of a
  structured request manifest embedded as a layer
  (`application/vnd.agentrc.policy.manifest.v1+json`); the platform fetches the
  full request set from that layer. Best for large request sets.

Both MUST be supported. Which is the **default** for `arc build` is
[open decision #1](/spec/); `--policy-mode inline|digest` is the seam. In digest
mode, a platform first reads the manifest layer named by the digest label, then
evaluates the requests exactly as in inline mode — the grant decision is
identical.

## 7. Override without a rebuild

Because every embedded resource carries an `.origin` reference, a platform or
organization MAY rewrite that reference at deploy time — for example, redirecting
a public MCP server to an internal mirror — **without rebuilding the artifact**:

```text
# As built
org.agentrc.mcp.github=sha256:abc123...
org.agentrc.mcp.github.origin=https://registry.agentrc.io/mcp/github:latest

# Platform override at deploy (mirror)
org.agentrc.mcp.github.origin=https://mirror.internal.acme/mcp/github:latest
```

The platform then fetches the mirror in place of the embedded copy. Whether such
narrowing / substitution MUST emit an audit record is
[open decision #4](/spec/).

## 8. Review before run

Because the manifest is labels, a registry maintainer or compliance reviewer can
vet an agent **before it ever runs**, with nothing more than `docker inspect` /
`arc inspect`. The labels reveal:

- **Identity &amp; capabilities** — `org.agentrc.identity.*`, `org.agentrc.capability.*`.
- **Requested model** — `org.agentrc.model.*` (name, min context, capabilities, fallback).
- **Requested network** — `org.agentrc.network.dns.*`, including `.source` for any auto-derived egress.
- **Tools / skills / MCP** — `org.agentrc.tool.*` / `.skill.*` / `.mcp.*`, with digests and `.origin`.
- **Sub-agent &amp; lifecycle limits** — `org.agentrc.agent.sub_agents*`, timeouts, retries.

Everything an agent *requests* is visible in the artifact. The platform remains
the authority that grants, narrows, or rejects each request and enforces via
Cedar ([Enforcement profile](/profiles/security/)).

## 9. Registry operations

The native CLI (`agentrc`, alias `arc`) and any OCI client operate on the
artifact directly:

```bash
arc build  -t ghcr.io/acme/claims-triage:1.0 [--policy-mode inline|digest] .
arc push   ghcr.io/acme/claims-triage:1.0
arc pull   ghcr.io/acme/claims-triage:1.0

# Inspect the labels (the manifest the platform reads)
arc inspect    ghcr.io/acme/claims-triage:1.0
docker inspect ghcr.io/acme/claims-triage:1.0
```

The BuildKit frontend (`# syntax=agentrc.agentfile/v0.1` + `docker build`) and
`arc build` MUST produce **identical** OCI artifacts and labels ([spec §10](/spec/)).

## 10. Signing, provenance, and reproducibility

As a standard OCI artifact, an agent supports the usual supply-chain controls:

- **Signing** — sign the artifact with [Sigstore](https://www.sigstore.dev/) /
  cosign; the signature attaches to the digest like any container image.
- **Provenance** — attach build provenance via [SLSA](https://slsa.dev/)
  attestations, so consumers can verify how the artifact was produced.
- **Reproducibility** — given the Agentfile, the pinned base (`FROM`) digest, and
  the resolved resource digests (in `org.agentrc.tool.*` / `.skill.*` / `.mcp.*`
  and the optional policy manifest), a rebuild SHOULD yield byte-identical layers
  and labels.

Because resources are content-addressed by digest, the artifact is portable
across registries and verifiable end to end.

## 11. Conformance

A package conforms to this profile (`agentrc/oci-labels/v0.1`) when:

1. It is a valid OCI artifact whose embedded `/mnt` resources are carried as
   standard layers.
2. Its image config carries the `org.agentrc.*` labels per [§3](#3-the-orgagentrc-label-catalog),
   matching the spec build tables ([spec §9](/spec/)).
3. SOP appears only as a `/mnt/SOP` pointer + `sop.sha256` digest, never as inline
   label text.
4. Every embedded MCP / skill carries both a digest value and an `.origin` label.
5. The request set is retrievable in the declared `--policy-mode` (inline or
   digest).

See the adversarial [conformance suite](/docs/conformance/) for the
`embedded-override-origin` case.
