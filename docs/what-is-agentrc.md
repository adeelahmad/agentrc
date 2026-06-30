---
layout: doc
title: What is agentrc?
description: "What agentrc is, the problem it solves, and why a Dockerfile-shaped recipe for AI agents is needed."
permalink: /docs/what-is-agentrc/
---
# What is agentrc?

**agentrc** is an open specification for **declaring, packaging, governing, and sharing AI agents** as portable, content-addressed artifacts. At its center is the **Agentfile** — a Dockerfile-shaped recipe for one agent.

You already know the shape: `FROM`, `CMD`, `COPY`, `ADD`, `LABEL`, `HEALTHCHECK`. agentrc adds just **four new keywords** — `IDENTITY`, `CAPABILITY`, `SOP`, and `POLICY` — and otherwise reuses standard Dockerfile keywords. The agentrc BuildKit frontend (or the `agentrc` CLI) compiles the Agentfile into an ordinary **OCI artifact** whose image config carries namespaced `org.agentrc.*` labels. A **platform** reads those labels — never the Agentfile source — and decides what to honour.

agentrc is **not** a runtime, sandbox, cloud platform, model provider, or agent framework. It is the neutral declaration, packaging, and governance layer that sits *above* all of those.

## The problem

AI agents are becoming real software that reads files, calls tools, spends credentials, and reaches the network — but the way they are defined today does not match the risk they carry:

- **Agents are not portable.** An agent built for one framework or cloud usually has to be rewritten to run anywhere else. Its capabilities and limits are scattered across code, config, environment variables, and platform dashboards.
- **Their requests are invisible.** There is rarely a single, machine-readable manifest that says *which tools, hosts, secrets, and model an agent wants*. A security team cannot vet what it cannot see in one place.
- **They are hard to share safely.** There is no common, signable, content-addressed package for an agent the way there is for a container image — so "here is the agent" usually means "here is some code, trust us."
- **Boundaries fail open.** When a platform cannot enforce a control an agent assumed, the agent often runs anyway, quietly less safe than intended.

The result: the same agent gets reimplemented per platform, its real privileges are unknowable, and nobody can sign off on it before it runs.

## What agentrc solves

agentrc gives you one Dockerfile-shaped recipe and one portable package that make an agent's intent explicit and machine-readable:

- **One recipe, any runner.** An `Agentfile` describes a single agent — identity, capabilities, system prompt, tools, skills, MCP servers, model, network, and operational constraints — independent of where it runs. Build it with `docker build` (via the `# syntax=agentrc.agentfile/v0.1` frontend) or with `arc build`; both emit the identical OCI artifact. Run it on a local process, a container, or a microVM — the substrate is a run-time choice, never an Agentfile directive.
- **`POLICY` is a request, not enforcement.** Each `POLICY` line is the developer *asking* the platform for a resource, a model, or an operational constraint. The platform (runtime / operator / organization) is free to **grant, narrow, or reject** any request. The Agentfile expresses *intent*; the platform holds *authority*.
- **Labels are the manifest.** The build translates authored intent into namespaced OCI labels under `org.agentrc.*`. The platform reads those labels at deploy / run time — **without parsing the Agentfile.** Labels are the machine-readable contract; the Agentfile is the human-authored recipe.
- **Secrets are deferred.** This draft defines no `SECRET`/`CRED` keyword and no secret schema; an agent that needs a credential leaves resolution entirely to the platform (Vault / broker / env / workload identity) — out of scope for now.
- **Cedar enforcement is platform-side, deny-by-default, fail closed.** The platform compiles each granted request plus its own organization rules into [Cedar](https://www.cedarpolicy.com/) and evaluates the grant: absence of a grant is a denial, an organization `forbid` defeats any agent request order-independently, and authorization tightens monotonically across `FROM`. A conformant platform that cannot enforce a required boundary refuses to run. Cedar is the enforcement engine and compilation target — **never an Agentfile author surface.**

## Who it is for

| You are… | agentrc gives you… |
|---|---|
| An **agent developer / adopter** | A Dockerfile-shaped recipe — four new keywords over keywords you already know — that builds with `docker build` or `arc build`. |
| A **security / compliance reviewer** | One artifact whose labels state every request: tools, network, model, sub-agents — vetted before it runs. |
| A **platform / runner author** | A labels-only contract: read `org.agentrc.*`, grant / narrow / reject, enforce with Cedar, fail closed. No need to parse the Agentfile. |
| A **registry maintainer** | A standard OCI artifact with digests and `.origin` labels you can mirror, sign, and attest. |

## Standards agentrc builds on

agentrc is deliberately a thin governance layer over proven, open standards rather than a reinvention of them:

| Concern | agentrc uses | v0.1 form |
|---|---|---|
| **MCP servers** | [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) — the open protocol for model/tool context | Projected under `/mnt/mcp/`; added with `COPY` / `ADD --remote`. agentrc declares and governs MCP, it does not replace it. |
| **Skills** | [Agent Skills](https://agentskills.io/) — the open `SKILL.md` format | Skill bundles under `/mnt/skills/`. |
| **Instructions** | [Agent SOP](https://github.com/strands-agents/agent-sop) — natural-language, RFC-2119-constrained operating procedures | The `SOP` keyword; embedded as a readable file at `/mnt/SOP` (the label is a pointer + digest, never the full text). |
| **Authorization** | [Cedar](https://www.cedarpolicy.com/) — the open authorization policy language from AWS | The **platform-side** enforcement engine and compilation target for typed `POLICY` requests. Not an Agentfile author surface. |
| **Packaging** | [OCI](https://opencontainers.org/) — content-addressed, signable artifacts | An OCI artifact: layers carry the `/mnt` resources; the image config carries `org.agentrc.*` labels. |

agentrc declares and governs these; it does not replace any of them.

## What agentrc is not

To stay useful to *every* runtime instead of competing with them, agentrc deliberately does not define a runtime, a sandbox, a model API, an agent framework, a tool-call wire protocol, a proprietary registry, or a second author-facing policy language. See [Non-goals](/docs/non-goals/) for the full list.

## Where to go next

- [Quickstart](/docs/quickstart/) — write, build, and run your first Agentfile.
- [The Agentfile](/docs/agentfile/) — the four new keywords and the `/mnt` projection.
- [Specification](/spec/) — the full 0.1.0-draft.5 working draft.
- [Core profile](/profiles/core/) — the minimal normative compiler behaviour.
- [Enforcement (Cedar) profile](/profiles/security/) — how the platform enforces granted requests.

<div class="callout">
<strong>In one line:</strong> The Agentfile is a Dockerfile-shaped recipe for one agent; the build emits <code>org.agentrc.*</code> OCI labels; the platform reads the labels — not the Agentfile — and grants, narrows, or rejects each request, enforcing the result with Cedar. The agent carries no secret value and writes no policy language.
</div>
