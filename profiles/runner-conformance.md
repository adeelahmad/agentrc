---
layout: doc
title: Runner Conformance
description: "Runner Conformance profile."
permalink: /profiles/runner-conformance/
---
# Platform Conformance Profile

**Version:** v1 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.io/agentfile:v1`)  
**Date:** 2026-06-30  
**Audience:** platform / runner authors

> The normative keywords **MUST**, **MUST NOT**, **SHOULD**, and **MAY** follow
> [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119). This profile defines what a
> **platform** (the runtime / operator / organization authority that runs an
> agentrc agent) must do to claim conformance. It is the operational companion to
> the [Specification](/spec/) and the [Enforcement profile](/profiles/security/).

## Purpose

agentrc itself is not a runner. An Agentfile compiles to an ordinary **OCI
artifact** carrying `org.agentrc.*` labels and layers; *something* must read those
labels and run the agent. This profile exists so independent platforms — clouds,
CLIs, sandboxes, Kubernetes operators, framework-native adapters — can state
exactly what level of agentrc support they provide, and so that two conformant
platforms decide and enforce grants **identically**.

The central contract is simple and absolute: **a conformant platform reads
labels, not the Agentfile.** The Agentfile is the human-authored recipe; the
labels are the machine-readable manifest. The platform never parses Agentfile
source at deploy / run time.

## The labels-only contract

A `POLICY` line in an Agentfile is a **request**, not enforcement. The author
expresses *intent*; the platform holds *authority*. For every request — including
the [auto-derived egress](/spec/) the compiler emits from hook / interrupt URLs —
the platform is free to **grant, narrow, or reject**.

```text
author writes          compiler emits              platform decides + enforces
─────────────          ──────────────              ───────────────────────────
POLICY request   ──►   org.agentrc.* labels   ──►  grant / narrow / reject,
(in the Agentfile)     (in the OCI image config)   then enforce with Cedar
```

Everything below is keyed off the labels, never off the source recipe.

## Core requirements

A platform claiming **agentrc Platform Conformance** MUST, at deploy / run time:

1. **Pull** the OCI artifact from any OCI-compatible registry.
2. **Read all `org.agentrc.*` labels** from the image config **without parsing the
   Agentfile.** Source is unavailable and unnecessary; the labels are the
   manifest.
3. **Evaluate each request** — every `POLICY`-derived label, including
   auto-derived `org.agentrc.network.dns.*` egress — against organization /
   platform policy and available resources, and **grant, narrow, or reject** it.
   Narrowing and rejection SHOULD be auditable (this auditability requirement is
   [open decision #4](/spec/)).
4. **Resolve secrets** named in `org.agentrc.secret.<name>=<scope>` via the
   platform's secret broker, at run time only. The value is never in the artifact;
   the platform injects it host-scoped (see [Secrets](#secrets)).
5. **Fetch `--runtime` resources now.** For each `org.agentrc.<kind>.<name>=runtime:<url>`
   label, fetch the resource at bootstrap and apply its failure mode:
   `--fail-if-unavailable` (default — refuse to boot) or `--warn-if-unavailable`
   (log and continue).
6. **Honour `*.origin` overrides.** The platform MAY substitute an embedded
   (`--cached`) resource by reading an overridden `org.agentrc.<kind>.<name>.origin`
   label — e.g. re-point a public MCP server to an internal mirror — without
   rebuilding the artifact.
7. **Enforce via Cedar** (see [Enforcement](#enforcement)): deny-by-default,
   `forbid` over `permit` order-independently, and monotonic intersection across
   `FROM`.
8. **Project the filesystem and boot `CMD`.** Load the SOP from `/mnt/SOP`,
   select / validate the model from `model.*`, project `/mnt/tools`, `/mnt/skills`,
   and `/mnt/mcp`, populate `/mnt/proc` with the live grant / identity / budget /
   audit state, then run `CMD` on the chosen substrate with the granted
   constraints.
9. **Fail closed.** If any *required* constraint cannot be enforced — an
   un-honourable secret scope, a network grant the platform cannot confine, a
   device it cannot isolate — the platform MUST refuse to boot rather than run
   the agent with weaker guarantees than were requested.

## Resolving labels

The label namespaces a conformant platform consumes:

| Label namespace | Source keyword | What the platform does |
|---|---|---|
| `org.agentrc.identity.*` | `IDENTITY` | Establishes the Cedar principal (`org.agentrc.identity.name`). |
| `org.agentrc.capability.*` | `CAPABILITY` | Advertised modalities; informational + match against substrate. |
| `org.agentrc.sop`, `org.agentrc.sop.sha256` | `SOP` | Locates and verifies the prompt file at `/mnt/SOP` (pointer + digest, never full text). |
| `org.agentrc.agent.*` | `POLICY agent.*` | Agent-side operational constraints / lifecycle / hooks. |
| `org.agentrc.substrate.*` | `POLICY substrate.*` | Resource requests (memory, cpu, device, ptty…). |
| `org.agentrc.model.*` | `POLICY model.*` | Model selection / fallback / required capabilities. |
| `org.agentrc.network.dns.*` | `POLICY network` + auto-derived | Egress grants; auto-derived entries carry a `.source` attribution. |
| `org.agentrc.tool.*` / `.skill.*` / `.mcp.*` | `COPY` / `ADD --remote` | Resource delivery: `local`, `<digest>` (+ `.origin`), or `runtime:<url>`. |
| `org.agentrc.secret.*` | `LABEL` | Names + host-scope of secrets the broker must resolve. |

The platform MUST support **both** policy encodings: **inline** values in labels,
and **digest** form where a label carries the digest of a structured manifest
embedded as a layer. The build seam is `--policy-mode inline|digest`; which is the
default is [open decision #1](/spec/).

## Secrets

```text
org.agentrc.secret.github_token=host:api.github.com
```

The secret **value never enters the artifact.** The platform's broker resolves the
named secret, scoped to the declared host, and injects it at run time (a
host-scoped substitution model in the spirit of
[microsandbox](https://docs.microsandbox.dev/sandboxes/secrets)). A conformant
platform MUST resolve secrets only at run time, MUST redact resolved values from
logs and audit records, and MUST fail closed if a required secret cannot be
resolved within its declared scope.

## Enforcement

Enforcement is **Cedar, platform-side.** Cedar MUST NOT appear in the Agentfile;
authors speak only typed `POLICY` requests. The platform compiles each granted
request **plus its own organization Cedar policies** into one Cedar `PolicySet`
and evaluates the grant. A conformant platform MUST derive Cedar entities /
actions from the request labels (principal = `org.agentrc.identity.name`):

| Request label | Cedar action | Cedar resource |
|---|---|---|
| `org.agentrc.network.dns.<host>=<port>` | `Action::"NetworkEgress"` | `Host::"<host>:<port>"` |
| `org.agentrc.tool.<name>` | `Action::"tool.invoke"` | `Tool::"<name>"` |
| `org.agentrc.mcp.<name>` | `Action::"mcp.request"` | `MCPServer::"<name>"` |
| `org.agentrc.agent.sub_agents=true` | `Action::"agent.delegate"` | `Agent::*` (capped by `sub_agents.max`) |
| `org.agentrc.substrate.device=<dev>` | `Action::"device.access"` | `Device::"<dev>"` |
| `org.agentrc.secret.<name>` | `Action::"secret.resolve"` | `Secret::"<name>"` |

The enforcement properties a conformant platform MUST preserve (these are
[Cedar's](https://www.cedarpolicy.com/)):

- **Deny-by-default.** Absence of a grant is a denial. An unrecognised request, an
  un-granted auto-derived egress, or an action with no matching `permit` is denied.
- **`forbid` overrides `permit`, order-independently.** An organization `forbid`
  (e.g. "no agent may reach the public internet") defeats any agent request,
  regardless of evaluation order.
- **Monotonic composition across `FROM`.** When an agent inherits
  (`FROM another-agent`), the effective authorization is the **intersection** of
  ceilings: a child's granted set MUST NOT exceed its parent's, and a parent
  `forbid` is un-loosenable by the child.

The agent's `POLICY` requests are the *floor of intent*; the organization's Cedar
policies are the *ceiling of authority*. One author surface (typed requests), one
engine (Cedar). The full mapping is normative in the
[Specification](/spec/) and the [Enforcement profile](/profiles/security/).

## Substrate-neutral

A conformant platform MAY run the agent on any substrate. The substrate is a
**run-time** choice (`agentrc run --isolation local|container|microvm
--substrate <driver>`), never an Agentfile directive, and it does **not** change
Agentfile semantics:

```text
local process
container (Docker / containerd)
gVisor-style sandbox
microVM
Kubernetes job
serverless worker
managed cloud agent runtime
SSH remote runner
framework-native adapter
```

Whatever the substrate, the `/mnt` projection (`/mnt/tools`, `/mnt/skills`,
`/mnt/mcp`, `/mnt/proc`, `/mnt/SOP`) and the label-derived grants are identical.
Placement is a platform concern expressed at run time; a portable artifact does
**not** hard-code it. A future companion document MAY define a separate run
manifest for placement.

## Required disclosure

A conformant platform SHOULD publish a support statement so authors and security
reviewers know what to expect. It SHOULD declare:

1. supported Agentfile syntax version (e.g. `agentrc.io/agentfile:v1`);
2. which `org.agentrc.*` label namespaces it reads and honours;
3. supported policy encoding(s) — inline, digest, or both;
4. supported secret broker backends and host-scope model;
5. supported substrates / isolation modes, if any;
6. resource delivery support — `--cached`, `--runtime`, and `*.origin` override;
7. audit / export formats for grant / narrow / reject decisions;
8. unsupported labels or constraints, and the failure behaviour for each;
9. known security limitations.

## Verifying conformance

This profile is verified against the adversarial
[conformance suite](/docs/conformance/), which probes the fail-closed and
deny-by-default guarantees above — for example, that an un-granted auto-derived
egress is denied, that a `forbid` defeats a conflicting request, and that a child
agent cannot widen a parent's locked constraint across `FROM`. The packaging side
of the contract is the [OCI labels &amp; package profile](/profiles/oci-package/),
and the `/mnt` layout the platform projects is the
[projection profile](/profiles/tool-projection/).
