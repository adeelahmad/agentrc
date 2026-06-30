---
layout: doc
title: What is AgentRC?
description: "What AgentRC is, the problem it solves, and why it is needed."
permalink: /docs/what-is-agentrc/
---
# What is AgentRC?

**AgentRC** (Agent Run Config) is an open specification for **declaring, packaging, securing, and sharing AI agents** as portable, content-addressed artifacts.

It defines the *contract* an agent declares — what it is, how it starts, what it may touch, and how those boundaries are governed. Compatible **runners** decide how to execute that contract on their own substrate.

AgentRC is **not** a runtime, sandbox, cloud platform, model provider, or agent framework. It is the neutral declaration, packaging, and governance layer that sits *above* all of those.

## The problem

AI agents are becoming real software that reads files, calls tools, spends credentials, and reaches the network — but the way they are defined today does not match the risk they carry:

- **Agents are not portable.** An agent built for one framework or cloud usually has to be rewritten to run anywhere else. Its capabilities and limits are scattered across code, config, environment variables, and platform dashboards.
- **Their permissions are invisible.** There is rarely a single reviewable artifact that says *which tools, files, hosts, and secrets this agent can use*. A security team cannot vet what it cannot see in one place.
- **They are hard to share safely.** There is no common, signable, content-addressed package for an agent the way there is for a container image — so "here is the agent" usually means "here is some code, trust us."
- **Boundaries fail open.** When a platform cannot enforce a control an agent assumed, the agent often runs anyway, quietly less safe than intended.

The result: the same agent gets reimplemented per platform, its real privileges are unknowable, and nobody can sign off on it before it runs.

## What AgentRC solves

AgentRC introduces one reviewable file and one portable package that make an agent's contract explicit and enforceable:

- **One declaration, any runner.** An `Agentfile` describes a single agent — its entrypoint, tools, mounts, network, credentials, and policy — independent of where it runs. Define once; run on Docker, gVisor, Firecracker, a cloud job, or a local process.
- **Security by declaration.** Every capability and boundary is written down explicitly. There is no ambient authority: undeclared access is denied by default, and policy travels *with* the package.
- **Fail closed, not open.** A conforming runner must enforce the boundaries it claims to support — or refuse to run. Unsupported security controls cause failure, never silent weakening.
- **Portable, signable packages.** Resolved dependencies are pinned in a lockfile and bundled into an OCI-compatible, content-addressed package that can be signed, inspected, and shared like a container image.
- **Reviewable governance.** Boundaries are expressed as machine-evaluable [Cedar](/spec/#11-cedar-policy-profile) policy, so a security reviewer — or a registry — can inspect exactly what an agent may do *before* it executes.

## Who it is for

| You are… | AgentRC lets you… |
|---|---|
| An **agent developer** | Define an agent once and have it run across runtimes without rewriting it. |
| A **security or compliance reviewer** | Read one file to see — and sign off on — exactly what an agent may access. |
| A **platform / runner author** | Consume a neutral spec instead of inventing your own agent format. |
| A **registry maintainer** | Distribute agents, bases, tools, and policies as signed, inspectable artifacts. |

## What AgentRC is not

To stay useful to *every* runtime instead of competing with them, AgentRC deliberately does not define a runtime, a sandbox, a model API, an agent framework, a tool-call wire protocol, or a proprietary registry. See [Non-goals](/docs/non-goals/) for the full list.

## Where to go next

- [Quickstart](/docs/quickstart/) — write and validate your first Agentfile.
- [Specification](/spec/) — the full working draft.
- [Core profile](/profiles/core/) — the minimal normative directive set.
- [Security](/docs/security/) — how declarative boundaries are enforced.

<div class="callout">
<strong>In one line:</strong> The Agentfile declares one agent. The lockfile pins its dependencies. The package makes it portable. The policy makes its boundaries reviewable. The registry makes it shareable. Compatible runners execute it.
</div>
