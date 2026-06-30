---
layout: doc
title: Platforms and substrates
description: "How a platform consumes an agentrc artifact by reading its labels, and how a substrate executes the agent — without agentrc becoming a runtime."
permalink: /docs/runners/
---
# Platforms and substrates

A **platform** consumes an agentrc artifact by reading its `org.agentrc.*`
labels — **never the Agentfile**. A **substrate** executes the agent's `CMD`.
agentrc is neither: it is the portable, governed declaration that both read.

The split matters. The platform is the *authority* — it reads the labels,
grants / narrows / rejects each request, resolves secrets, and enforces the
result via Cedar. The substrate is the *execution driver* — local process,
container, microVM, or otherwise — that the platform drives to run `CMD` once
the grant is decided. The substrate is chosen **at run time**
(`arc run --isolation … --substrate …`), not in the Agentfile.

## Substrate examples

A substrate is whatever actually runs `CMD` with the granted constraints. Any
of these can be a substrate, behind one platform contract:

- local development process
- Docker / OCI container
- gVisor-sandboxed container
- Firecracker or microsandbox microVM
- Kubernetes job
- serverless function
- managed cloud agent runtime adapter
- framework-native adapter for Strands, LangChain, CrewAI, Langflow, or similar

```bash
arc run ghcr.io/acme/claims-triage:1.0 --isolation microvm
arc run ghcr.io/acme/claims-triage:1.0 --isolation container --substrate gvisor
```

The same artifact runs on any of them. The platform reads the same labels and
makes the same grant decision regardless of which substrate carries the agent.

## Conformance principle

A platform claims conformance only to the profiles it implements. The contract
is **labels in, governed execution out** — at no point does a conformant
platform parse the Agentfile source.

| A conformant platform MUST… | Detail |
|---|---|
| Read labels, not the Agentfile | Load every `org.agentrc.*` label from the image config; decide entirely from labels. |
| Grant / narrow / reject each request | Evaluate every `POLICY`-derived request — including auto-derived egress — against org / platform policy and available resources. Decisions SHOULD be auditable. |
| Resolve secrets via the broker | Resolve each `org.agentrc.secret.<name>=<scope>` reference through the secret broker; the value never lives in the artifact. |
| Fetch `--runtime` resources | Fetch `runtime:<url>` resources at bootstrap; honour `--fail-if-unavailable` (refuse to boot) and `--warn-if-unavailable` (log and continue). |
| Substitute via `.origin` | MAY redirect an embedded resource by honouring an overridden `org.agentrc.*.origin` label (e.g., a public MCP server to an internal mirror) without rebuilding. |
| Enforce via Cedar | Compile granted requests plus org rules into Cedar and evaluate: **deny-by-default**, **`forbid` over `permit`** (order-independent), **monotonic intersection across `FROM`**. |
| Project `/mnt` | Load the SOP from `/mnt/SOP`, select / validate the model from `model.*`, project `/mnt/tools`, `/mnt/skills`, `/mnt/mcp`, and populate `/mnt/proc`. |
| Fail closed | If a required boundary cannot be enforced, a required secret cannot be resolved, or policy cannot be evaluated, refuse to boot. |

> **`POLICY` lines are requests, not enforcement.** The platform holds the
> authority. Deny-by-default applies to the grant decision: an unrecognised or
> disallowed request — including an un-granted auto-derived egress — is never
> silently honoured.

## The key boundary

agentrc defines **what must be true** from the labels. Platforms decide **how to
make it true**, and substrates carry it out.

This lets AWS, Google, Docker, gVisor, Firecracker, microsandbox, Kubernetes,
and local platforms implement their own execution and enforcement layer without
owning — or even parsing — the portable agent declaration format.

The full normative obligations are in the
[Platform conformance profile](/profiles/runner-conformance/), and the Cedar
enforcement model platforms compile to is the
[Enforcement (Cedar) profile](/profiles/security/).
