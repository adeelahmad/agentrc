---
layout: doc
title: Specification
description: "agentrc Agentfile and Package Specification working draft."
permalink: /spec/
---
# agentrc Specification

**Version:** 0.1.0-draft.4  
**Status:** Working Draft  
**Editor:** [Adeel Ahmad](https://www.linkedin.com/in/adeelahmadch)  
**Date:** 2026-06-30  
**Audience:** agent developers, platform engineers, security engineers, registry maintainers, runner implementers, compliance reviewers

## 1. Status of this document

This document is a working draft of the agentrc Agentfile and Package Specification.

This specification defines:

1. a declarative file format for a single AI agent, called an **Agentfile**;
2. a lockfile format for resolved agent dependencies;
3. an OCI-compatible package model for distributing agents;
4. a security-boundary declaration model for filesystems, network egress, credentials, tools, functions, skills, MCP servers, rate limits, and audit;
5. conformance profiles for tools, policy, package distribution, and runners.

This specification does **not** define a runtime. It does **not** require a particular sandbox, model provider, agent framework, cloud provider, or tool-calling transport.

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** are to be interpreted as normative requirement levels.

## 2. Abstract

agentrc defines a portable declaration and packaging format for AI agents.

The **Agentfile** declares what an agent is and what it is allowed to request. The **agent package** carries the Agentfile, lockfile, policy, tools, functions, skills, metadata, and optional source artifacts. Compatible **runners** execute the package and enforce the declared boundaries.

The intended result is separation of concerns:

| Concern | Owned by |
|---|---|
| Agent declaration | Agentfile |
| Dependency pinning | agentrc.lock |
| Package distribution | OCI-compatible registry |
| Security and governance intent | Policy and boundary declarations |
| Execution and isolation | Compatible runner |
| Cloud/runtime-specific enforcement | Runner implementation |
| Multi-agent orchestration | Future workflow companion spec |

agentrc allows developers to define an agent once, security teams to review it before execution, registries to share it as a reusable artifact, and runners to execute it without requiring the agent to be rewritten for a specific runtime or cloud.

## 3. Design goals

The design goals are:

1. **Declarative single-agent format** — one reviewable file describes one agent.
2. **Runtime neutrality** — the package is not bound to Docker, Kubernetes, gVisor, Firecracker, microsandbox, a managed cloud, or any specific local runner.
3. **Security by declaration** — access to tools, credentials, networks, mounts, functions, skills, and sub-systems is explicit.
4. **Governance portability** — policy travels with the package and can be evaluated before execution.
5. **Registry sharing** — agents, bases, skills, tools, toolsets, and policies can be distributed using content-addressed registries.
6. **Reproducibility** — resolved dependencies and digests are captured in a lockfile.
7. **Profile-based conformance** — implementations can support the core Agentfile without implementing every optional execution profile.
8. **Extensibility** — multi-agent workflows, tool projection, and runner-specific enforcement are layered profiles, not hard-coded into the core format.

## 4. Non-goals

agentrc does not define:

1. a new runtime;
2. a new sandboxing technology;
3. a new model API;
4. a new agent framework;
5. a new wire protocol for tool calls;
6. a replacement for MCP or A2A;
7. a replacement for Cedar, OCI, OpenTelemetry, Sigstore, SLSA, Docker, Kubernetes, gVisor, Firecracker, or cloud runtimes;
8. a full workflow language in the single-agent format;
9. cloud infrastructure provisioning;
10. a central proprietary registry.

agentrc packages declare requirements. Compatible runners decide whether and how they can satisfy those requirements.

## 5. Core terminology

### 5.1 Agent

An **agent** is a software unit with an entrypoint, capabilities, dependencies, and declared boundaries.

### 5.2 Agentfile

An **Agentfile** is the source declaration for a single agent package.

### 5.3 Agent package

An **agent package** is the built artifact produced from an Agentfile. It contains the Agentfile, lockfile, policy, metadata, and supporting artifacts required to instantiate the agent.

### 5.4 Runner

A **runner** is an implementation that can execute an agentrc package. Examples include a local CLI runner, a container runner, a gVisor-backed runner, a microVM runner, a serverless runner, or a managed cloud runner.

A runner is outside the core Agentfile format. A runner may claim conformance to one or more profiles defined by this specification.

### 5.5 Capability

A **capability** is any declared thing the agent can use to read, write, call, invoke, or affect the outside world. Tools, functions, MCP servers, mounts, URLs, credentials, skills, memories, and delegated agents are capabilities.

### 5.6 Boundary

A **boundary** is a declared limit on a capability. Boundaries include filesystem access, network egress, credential access, tool permission, rate limits, timeouts, resource limits, and audit requirements.

### 5.7 Policy

A **policy** is a machine-evaluable rule set governing a capability request. The default agentrc policy profile is Cedar.

### 5.8 Registry

A **registry** is a content-addressed distribution system that stores agentrc packages or package components. OCI registries are the default distribution mechanism.

### 5.9 Profile

A **profile** is a named subset of conformance requirements. agentrc uses profiles so the specification can remain neutral while still allowing concrete implementations to prove compatibility.

## 6. Layering model

agentrc is divided into six layers.

| Layer | Purpose | Defined here? |
|---|---|---|
| Core Agentfile | Single-agent declaration grammar and semantics | Yes |
| Security profile | Boundary, policy, credential, audit semantics | Yes |
| Package profile | Lockfile and OCI package layout | Yes |
| Registry profile | Pull/push/inspect/verify expectations | Yes |
| Runner profile | Execution and enforcement obligations for compatible runners | Limited; profile only |
| Workflow profile | Multi-agent state-machine orchestration | Future companion draft |

The core specification deliberately stops at the package boundary. Execution is performed by runners.

## 7. Agentfile syntax model

An Agentfile is a line-oriented text file using Dockerfile-like directives.

Rules:

1. Blank lines are ignored.
2. Lines whose first non-whitespace character is `#` are comments.
3. Inline comments MAY appear after whitespace followed by `#`.
4. Directive names are uppercase ASCII tokens.
5. Arguments are whitespace-separated unless a future profile defines quoting semantics.
6. `POLICY` and `SOP` open multi-line blocks closed by `END`.
7. Unknown directives MUST fail validation unless explicitly placed in an extension namespace profile.

A syntax marker MAY be carried as a comment so older parsers can ignore it:

```dockerfile
# syntax=agentrc.agentfile/v0.1
```

## 8. Minimal valid Agentfile

```dockerfile
# syntax=agentrc.agentfile/v0.1

AGENT hello-local
CMD strands run --agent hello-local
TOOL utcp:file_read
AUDIT basic

POLICY
  permit(
    principal == AgentRC::Agent::"hello-local",
    action in [AgentRC::Action::"tool.invoke"],
    resource == AgentRC::Tool::"file_read"
  );
END
```

For local development, a runner MAY infer an agent identity from `FROM` or the package reference if `AGENT` is absent. Published packages SHOULD declare `AGENT` explicitly.

## 9. Agentfile directives

This section defines the full directive vocabulary. Only a minimal subset is **normative core** for v0.1 — `AGENT`, `FROM`, `CMD`, `TOOL`, `MOUNT`, `CRED`, `URL`, `POLICY`, and `AUDIT`. Every other directive below belongs to an optional profile; see the [Core profile](/profiles/core/) for the full classification table. Implementations MAY parse additional directives under extension profiles, but they MUST NOT silently reinterpret core directives, and an implementation MUST publish which profiles it actually supports rather than implying that every recognized directive is enforced.

The directives in §9.13 and §9.21–§9.25 (`BIND`, `SLICE`, `IMAGE`, `ISOLATION`, `BROKER`, `BACKEND`) plus `PLUGIN` (§9.15) are **placement** concerns, not capability concerns. They describe where and how an agent runs rather than what it is. To preserve content-addressed, signable identity, these SHOULD be carried in a separate **Run Manifest** chosen by the operator or runner and SHOULD NOT contribute to the package digest. See §24 Q1.

### 9.1 `AGENT`

Declares the stable identity of the agent.

```dockerfile
AGENT code-reviewer
```

Published packages SHOULD include `AGENT`. The value SHOULD be a lowercase DNS-label-compatible identifier.

### 9.2 `FROM`

Declares a base package or base environment.

```dockerfile
FROM ghcr.io/org/base-agent:1.4
FROM python:3.11-slim
FROM scratch
```

`FROM` MAY refer to an agentrc agent package, an OCI image, or `scratch`.

If the reference is an agentrc package, inheritance rules apply. If the reference is an ordinary OCI image, it is treated as a base execution environment and no agent inheritance is implied.

A child package MUST NOT weaken inherited security ceilings from a parent agentrc package.

### 9.3 `SHELL`

Declares the agent framework, shell, or adapter expected by the agent entrypoint.

```dockerfile
SHELL strands
SHELL strands:shell@2026.5
SHELL langchain
SHELL crewai
```

`SHELL` is not a runtime. It is a framework or command adapter hint. A runner that does not support the requested shell MUST report an unsupported-shell error.

### 9.4 `CMD`

Declares the command used to start the agent.

```dockerfile
CMD strands run --agent code-reviewer
CMD claude --print
CMD python -m agent.main
```

The command is interpreted inside the execution environment chosen by the runner.

### 9.5 `TOOL`

Declares a tool capability.

```dockerfile
TOOL utcp:file_read
TOOL utcp:shell
TOOL strands:current_time
TOOL mcp:github.get_issue
```

A tool reference SHOULD identify a namespace and tool name. Tool metadata SHOULD include input schema, output schema, version, digest, risk tier, and provenance.

agentrc does not invent a tool-calling transport. Tools described and invoked through the **[Universal Tool Calling Protocol (UTCP)](https://www.utcp.io/)** use the `utcp:` namespace, and SHOULD follow the [UTCP specification](https://github.com/universal-tool-calling-protocol/utcp-specification) for their tool descriptions. Tools provided over the Model Context Protocol use the `mcp:` namespace (see §9.10).

Declaring a tool does not automatically grant use of the tool unless the active security profile allows it. Under deny-by-default security, the tool must also be permitted by policy.

### 9.6 `TOOLSET`

Declares a bundle of tools.

```dockerfile
TOOLSET strands:workspace
TOOLSET registry://tools.example.com/platform-safe@sha256:...
```

A toolset MUST resolve to a deterministic list of tool references during lockfile generation.

### 9.7 `FUNCTION`

Declares a typed function capability.

```dockerfile
FUNCTION ./functions/review.py:summarize_diff
FUNCTION summarize_meeting
```

A function reference MAY be a local module path and symbol, or a symbolic name resolved by a runner/framework adapter.

Function metadata SHOULD include input schema, output schema, version, digest, and postcondition metadata where available.

### 9.8 `SKILL`

Declares a skill bundle.

```dockerfile
SKILL pr-review
SKILL ./skills/incident-sop
SKILL ghcr.io/org/skills/pii-review:1.2
```

A skill SHOULD contain a `SKILL.md` file. Skill metadata SHOULD include name, version, description, license, digest, and provenance.

### 9.9 `SERVER`

Declares an external server dependency, typically an MCP server or API-backed tool server.

```dockerfile
SERVER github https://api.github.com
SERVER internal https://tools.internal.example/mcp
```

Calls to declared servers remain subject to network and policy boundaries.

### 9.10 `MCP`

Declares an MCP dependency by logical name or registry reference.

```dockerfile
MCP github
MCP registry://modelcontextprotocol/filesystem@1.0.2
```

MCP dependencies follow the **[Model Context Protocol (MCP)](https://modelcontextprotocol.io/)** specification. agentrc does not replace or re-specify MCP; it declares the dependency and governs it.

An MCP declaration grants no ambient access. The runner must still apply policy, credentials, and network boundaries to MCP calls.

### 9.11 `URL`

Declares a network destination that the package requests.

```dockerfile
URL https://api.github.com
URL https://registry.npmjs.org
```

Network access is deny-by-default under the Security Declaration Profile. A declared URL is a requested boundary, not an entitlement. The runner MUST fail closed if required egress policy cannot be enforced, unless the operator explicitly overrides the package outside the package itself.

### 9.12 `CRED`

Declares a credential or secret reference.

agentrc adopts the **host-scoped placeholder-substitution** secret model defined by **[microsandbox](https://docs.microsandbox.dev/sandboxes/secrets)**. A secret is bound to a **source** and to one or more **allowed hosts**. The real value MUST NOT enter the agent's execution environment: the agent only ever sees a placeholder, and a conforming runner swaps the placeholder for the real value at the **network layer**, only for outbound requests to the secret's allowed hosts.

```dockerfile
CRED github_token env:GITHUB_TOKEN host:api.github.com
CRED openai_key  env:OPENAI_API_KEY host:api.openai.com
CRED stripe_key  env:STRIPE_KEY     host:*.stripe.com inject:header
CRED db_password vault://secret/data/prod/db#password host:db.internal.example
CRED github_token keyring:github_token host:api.github.com
```

Syntax:

```text
CRED <name> <source> host:<host-or-pattern> [host:<host-or-pattern> ...] [inject:header|query]
```

| Part | Meaning |
|---|---|
| `<name>` | Logical secret name. The agent environment sees a placeholder (e.g. `$<NAME>`), never the value. |
| `env:<VAR>` | Source the value from an environment variable on the host resolving the package. |
| `vault://<path>#<key>` | Source the value from a Vault-like backend. |
| `keyring:<key>` | Source the value from an OS/user keyring. |
| `host:<host>` | An allowed destination host the secret may be substituted into. Patterns such as `*.stripe.com` are permitted. Requests to any other host MUST NOT receive the value. |
| `inject:header` \| `inject:query` | Where the runner injects the secret on an allowed request. Defaults to `header`. |

Rules:

1. A package MUST NOT contain plaintext secret values.
2. The secret value MUST NOT be exposed to agent code, tools, the filesystem, environment dumps, traces, or audit — only a placeholder is visible inside the sandbox.
3. A runner MUST substitute the value only for outbound requests to a declared allowed host; a secret with no `host:` binding MUST NOT be substituted into any request.
4. Runners MUST redact credential values from logs, audit events, lockfiles, and package metadata.
5. If a runner cannot enforce host-scoped substitution for a required secret, it MUST fail closed rather than expose the value to the sandbox.

For backward compatibility a runner MAY accept a host-less `CRED <name> <source>`, but it SHOULD warn, because an unbound secret cannot be safely placeholder-substituted.

### 9.13 `BIND`

Declares host-to-agent filesystem binding intent.

```dockerfile
BIND ./repo /workspace copy
BIND ./out /output direct
BIND ./data /data ro
```

Syntax:

```text
BIND <source> <target> <mode>
```

Recommended modes:

| Mode | Meaning |
|---|---|
| `copy` | Copy source into the execution environment before start |
| `direct` | Expose source directly if supported by runner |
| `ro` | Read-only mount |
| `rw` | Read-write mount |

Compatibility note: early implementations may accept `BIND <source> <target>` and normalize it to `copy`.

### 9.14 `MOUNT`

Declares an internal mount path and mutability.

```dockerfile
MOUNT /workspace rw
MOUNT /data ro
```

`MOUNT` is a sandbox-facing boundary. `BIND` describes source-to-target mapping; `MOUNT` describes the path the agent expects to see. Runners MAY satisfy a `MOUNT` by copy, bind mount, volume, object store mount, memory filesystem, or another substrate-specific mechanism.

### 9.15 `PLUGIN`

Declares an execution capability module requested from the runner.

```dockerfile
PLUGIN network
PLUGIN filesystem
PLUGIN secrets
```

`PLUGIN` is not core execution semantics. It is a runner capability requirement. A runner that cannot provide a required plugin MUST fail with an unsupported-plugin error.

### 9.16 `POLICY` block

Declares inline policy source.

```dockerfile
POLICY
  permit(
    principal == AgentRC::Agent::"code-reviewer",
    action in [AgentRC::Action::"tool.invoke"],
    resource == AgentRC::Tool::"file_read"
  );
END
```

The default policy language is Cedar. A runner that claims Cedar support MUST evaluate with Cedar semantics, not string matching.

### 9.17 `ALLOW`, `DENY`, and `RATELIMIT`

Declare shorthand policy rules.

```dockerfile
ALLOW invoke utcp:file_read
DENY invoke utcp:shell
RATELIMIT file_write 20/hour
```

These directives are convenience syntax. A compiler MAY lower them into the active policy profile. If both an allow and deny match the same request, deny MUST win.

### 9.18 `AUDIT`

Declares audit expectations.

```dockerfile
AUDIT basic
AUDIT all
AUDIT compliance
```

Recommended values:

| Value | Meaning |
|---|---|
| `off` | No package-level audit requirement |
| `basic` | Invocation and lifecycle audit |
| `all` | Tool/function/policy/credential/rate-limit audit |
| `compliance` | Structured, exportable audit suitable for review |

If audit is required and the runner cannot produce audit events, the runner MUST fail closed.

### 9.19 `TIMEOUT`

Declares a default execution timeout in seconds.

```dockerfile
TIMEOUT 30
```

### 9.20 `LIMIT`

Declares a named resource, output, or behavioral limit.

```dockerfile
LIMIT max_output 1048576
LIMIT tool_calls 100
LIMIT network_requests 50/min
```

A runner MUST report unsupported required limits. A runner SHOULD fail closed for security-relevant limits it cannot enforce.

### 9.21 `SLICE`

Declares resource shape requirements.

```dockerfile
SLICE cpu=2 mem=2048
```

`SLICE` is a runner resource hint/profile directive. It does not require agentrc to define resource isolation itself.

### 9.22 `IMAGE`

Declares an image expected by a runner or package builder.

```dockerfile
IMAGE ghcr.io/agentrc/dev-sandbox:1.0
```

`IMAGE` is an execution/package hint. It MUST NOT be interpreted as making agentrc a runtime.

### 9.23 `ISOLATION`

Declares a requested isolation profile.

```dockerfile
ISOLATION local
ISOLATION container
ISOLATION microvm
ISOLATION ssh
```

`ISOLATION` is a compatibility directive for current implementations. In the core specification, execution placement SHOULD be selected by the operator or runner. If `ISOLATION` appears in an Agentfile, it is a requested runner profile, not an implementation mandate.

A future version may move placement into a separate run manifest.

### 9.24 `BROKER`

Declares a secret broker family.

```dockerfile
BROKER vault
BROKER local
```

A broker is a credential-resolution implementation detail. The Agentfile declares the requested broker type, but the runner controls actual resolution.

### 9.25 `BACKEND`

Declares a remote backend or execution endpoint.

```dockerfile
BACKEND ssh://runner.internal.example
```

`BACKEND` is a runner/implementation extension. Published portable packages SHOULD avoid relying on a single backend unless the package is intentionally environment-specific.

### 9.26 `TRACE`

Declares trace emission or trace capture preference.

```dockerfile
TRACE on
TRACE otlp://collector.internal:4317
```

Trace values MUST NOT contain secrets.

### 9.27 `MEMORY`

Declares a bounded memory or state surface.

```dockerfile
MEMORY writing ./schemas/writing.schema.json schema:WritingMemory mode:rw
```

Memory declarations SHOULD include a schema and mode. Runners MUST apply policy and audit semantics to memory read/write operations when supported.

### 9.28 `OPTIMIZER`

Declares an optimizer hook or optimization strategy.

```dockerfile
OPTIMIZER textgrad
```

Optimizers are optional. A runner/compiler that does not support the optimizer MUST report it clearly. Optimizer outputs MUST NOT silently weaken policy or boundaries.

### 9.29 `HEALTHCHECK`

Declares a health probe.

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s CMD utcp:ping
```

A runner MAY use health checks for lifecycle management. A runner that ignores health checks SHOULD report that the directive is unsupported.

### 9.30 `SOP`

Declares a **Standard Operating Procedure** — a reusable, natural-language instruction set that guides the agent through a multi-step task. agentrc adopts the **[Agent SOP](https://github.com/strands-agents/agent-sop)** model from Strands Agents: SOPs are markdown instruction sets with parameterized inputs and RFC 2119 constraints (`MUST` / `SHOULD` / `MAY`) that steer agent behavior without rigid scripting.

A `SOP` may be referenced from a file or registry, or embedded inline in the Agentfile as a block closed by `END`:

```dockerfile
# Reference an external SOP (markdown, or an Agent Skill SKILL.md bundle)
SOP ./sops/code-review.md
SOP ghcr.io/org/sops/incident-response:1.2

# Or embed an SOP directly in the Agentfile
SOP code-review
  ## Overview
  Review a pull request for correctness and security.

  ## Parameters
  - **diff** (required): the unified diff to review

  ## Steps
  ### 1. Triage
  - You MUST read the full diff before commenting.
  - You MUST NOT approve if any required check is failing.

  ### 2. Report
  - You SHOULD group findings by severity.
  - You MAY suggest concrete patches.
END
```

The two forms are distinguished by the argument: an `SOP` line whose argument is a file path or registry/OCI reference (it contains `/`, `://`, or a file extension) is the **single-line reference** form; otherwise the argument is a bare name and the directive **opens a block** closed by a line containing only `END`.

An SOP SHOULD follow the Agent SOP section structure (`Overview`, `Parameters`, `Steps` with constraints, optional `Examples` and `Troubleshooting`). Because Agent SOPs are distributable as the open `SKILL.md` [Agent Skills](https://agentskills.io/) format, an SOP MAY also be packaged and resolved as a `SKILL` (§9.8).

Semantics:

1. Embedded SOP blocks are captured verbatim, exactly like `POLICY` blocks, and MUST NOT be parsed as Agentfile directives.
2. SOP content is **instruction**, not capability or authority: an SOP MUST NOT grant access that policy and boundary declarations do not already permit. Instructions never widen the security ceiling.
3. SOP references SHOULD be pinned by digest in the lockfile so the embedded guidance is reproducible.
4. A runner that injects SOP instructions into the agent's context SHOULD record which SOPs were active when `AUDIT` requires it.

## 10. Security model

agentrc security is declarative and reviewable.

The Agentfile declares requested boundaries; a conforming runner enforces supported boundaries or fails closed.

Principles:

1. **No ambient authority** — tools, credentials, files, networks, functions, skills, and servers are declared explicitly.
2. **Deny by default** — undeclared access is not granted under the Security Declaration Profile.
3. **Fail closed** — unsupported required security controls cause failure, not silent weakening.
4. **Deny wins** — forbids override permits.
5. **Secrets are host-scoped references** — no plaintext secret values appear in source, lockfile, package, logs, or audit. Following [microsandbox](https://docs.microsandbox.dev/sandboxes/secrets), a secret value never enters the sandbox; the agent sees a placeholder and the runner substitutes the real value only at the network layer for the secret's declared allowed hosts.
6. **Policy travels with package** — policy source or compiled policy profile is content-addressed and included in the package.
7. **Inheritance tightens** — child packages cannot weaken inherited ceilings.
8. **Operator override is explicit** — any override is outside the package and auditable.

## 11. Cedar policy profile

The default agentrc policy profile is **[Cedar](https://www.cedarpolicy.com/)**, the open-source authorization policy language created by AWS. agentrc adopts Cedar's syntax and semantics as its policy lingua franca rather than defining a new policy language; policies are written in Cedar and evaluated with a Cedar-compatible evaluator.

A Cedar request evaluates:

| Cedar element | agentrc meaning |
|---|---|
| Principal | agent identity, caller, or delegated agent |
| Action | operation such as `tool.invoke`, `function.invoke`, `cred.resolve`, `network.egress` |
| Resource | tool, function, skill, host, path, credential, memory, server, or agent |
| Context | request attributes such as path, URL, destination host, caller, rate bucket, or mount mode |

Recommended namespace examples:

```text
AgentRC::Agent::"code-reviewer"
AgentRC::Tool::"file_read"
AgentRC::Function::"summarize_meeting"
AgentRC::Skill::"pr-review"
AgentRC::Credential::"github_token"
AgentRC::Host::"api.github.com"
AgentRC::Path::"/workspace"
AgentRC::Action::"tool.invoke"
AgentRC::Action::"function.invoke"
AgentRC::Action::"cred.resolve"
AgentRC::Action::"network.egress"
```

A runner that claims support for the Cedar profile MUST:

1. parse Cedar policy with a Cedar-compatible evaluator;
2. deny on parse or evaluation failure;
3. enforce forbid-overrides-permit behavior;
4. include request context in evaluation;
5. record allow and deny decisions when audit is required.

## 12. Agent package model

An agentrc package is a content-addressed bundle containing:

1. `Agentfile`;
2. `agentrc.lock`;
3. agent metadata/config;
4. policy sources or compiled policy profile;
5. tool metadata or bundled tools;
6. function code or function references;
7. skill bundles;
8. optional source layer;
9. optional framework/runtime adapter metadata;
10. provenance and signature metadata where available.

The package is a recipe for instantiating an agent, not a live VM snapshot and not a specific runtime instance.

## 13. Lockfile model

The lockfile pins resolved artifacts.

A lockfile SHOULD contain:

```json
{
  "version": "1",
  "agentfile_hash": "sha256-or-short-hash",
  "base_image": {
    "registry": "ghcr.io",
    "name": "org/base-agent",
    "tag": "1.0",
    "digest": "sha256:..."
  },
  "tools": [
    {
      "name": "utcp:file_read",
      "version": "1.0.0",
      "schema_hash": "sha256:...",
      "source": "registry://tools/file_read",
      "digest": "sha256:..."
    }
  ],
  "functions": [
    {
      "path": "./functions/review.py:summarize_diff",
      "hash": "sha256:...",
      "type_hint_digest": "sha256:...",
      "source": "local",
      "version": ""
    }
  ],
  "skills": [
    {
      "name": "pr-review",
      "version": "2.1.0",
      "hash": "sha256:...",
      "source": "./skills/pr-review"
    }
  ],
  "sops": [
    {
      "name": "pr-review",
      "version": "1.0.0",
      "hash": "sha256:...",
      "source": "inline"
    }
  ],
  "policy_hash": "sha256:...",
  "credential_names": ["github_token"],
  "timestamp": "2026-06-30T00:00:00Z"
}
```

A lockfile MUST NOT contain credential values.

## 14. OCI package profile

agentrc packages SHOULD be distributable using OCI-compatible registries.

Recommended media types:

| Component | Media type |
|---|---|
| Manifest | `application/vnd.oci.image.manifest.v1+json` or OCI artifact equivalent |
| Agent config | `application/vnd.agentrc.agent.config.v1+json` |
| Agentfile | `application/vnd.agentrc.agent.agentfile.v1+text` |
| Lockfile | `application/vnd.agentrc.agent.lock.v1+json` |
| Cedar policy | `application/vnd.agentrc.agent.policy.cedar.v1+text` |
| Source layer | `application/vnd.agentrc.agent.source.layer.v1.tar+gzip` |
| Tool layer | `application/vnd.agentrc.agent.tool.layer.v1.tar+gzip` |
| Skill layer | `application/vnd.agentrc.agent.skill.layer.v1.tar+gzip` |

Recommended annotations:

```text
io.agentrc.agentfile.version
io.agentrc.agent.name
io.agentrc.policy.hash
io.agentrc.base.digest
io.agentrc.risk.tier
org.opencontainers.image.title
org.opencontainers.image.version
org.opencontainers.image.source
org.opencontainers.image.revision
org.opencontainers.image.created
```

Implementations SHOULD support signing and provenance metadata.

## 15. Registry model

A registry MAY host:

1. complete agent packages;
2. base agent packages;
3. skills;
4. tools;
5. toolsets;
6. policy bundles;
7. verified enterprise templates;
8. metadata/risk reports.

Registry clients SHOULD support:

1. pull by digest;
2. signature verification;
3. provenance inspection;
4. policy inspection before execution;
5. risk-tier display;
6. dependency graph display;
7. offline/cache mode.

agentrc does not require a central registry. Public community registries, private enterprise registries, and air-gapped registries should all be possible.

## 16. Inheritance and monotonic security

If `FROM` references another agentrc package, the child inherits the parent package's declared capabilities and restrictions.

A child package MAY:

1. add implementation detail;
2. add skills;
3. add tools allowed by the parent ceiling;
4. narrow network access;
5. narrow filesystem access;
6. add stricter policy;
7. increase audit requirements.

A child package MUST NOT:

1. remove inherited denies;
2. widen network access beyond the parent ceiling;
3. widen filesystem access beyond the parent ceiling;
4. expose credentials outside the parent ceiling;
5. disable inherited audit requirements;
6. increase delegation authority beyond the parent ceiling;
7. weaken inherited rate limits unless an explicit parent rule permits it.

This allows a security team to publish hardened base packages that downstream teams can reuse safely.

## 17. Tool projection profile

The core Agentfile does not require a single tool implementation model.

A tool MAY be implemented as:

1. a local executable;
2. a framework function;
3. an MCP tool;
4. a remote API adapter;
5. a runner-native capability.

The optional Tool Projection Profile defines a command-oriented surface so tools can be discovered, inspected, and invoked without a provider-specific SDK.

Recommended projection:

```text
/agentrc/tools/<tool-name>
/agentrc/tools/<tool-name>.schema.json
/agentrc/tools/<tool-name>.help
/agentrc/tools/<tool-name>.events
/agentrc/functions/<function-name>
/agentrc/skills/<skill-name>/SKILL.md
/agentrc/proc/audit
/agentrc/proc/policy
/agentrc/proc/limits
```

This profile is optional. It should not be used to imply agentrc is a runtime. It is a portability profile that runners may implement.

## 18. Runner conformance profile

A runner may claim conformance to one or more profiles.

A runner claiming the Core Runner Profile MUST:

1. read the package metadata;
2. identify supported and unsupported directives;
3. reject unknown required directives;
4. execute the declared `CMD` or fail with an unsupported-entrypoint error;
5. enforce declared security boundaries it claims to support;
6. fail closed when it cannot enforce a required boundary;
7. produce audit events if required by `AUDIT`;
8. never leak credential values in logs or audit;
9. report the effective execution profile.

A runner MAY be implemented using local processes, containers, gVisor, microVMs, cloud jobs, Kubernetes, SSH, serverless workers, or managed agent runtimes. The implementation substrate is not specified by agentrc.

## 19. Multi-agent workflows

The Agentfile defines one agent.

Multi-agent workflows should be defined in a companion specification, not inside Agentfile.

A future workflow specification SHOULD reference packaged agents by OCI reference or digest and MAY use state-machine concepts such as:

1. `Task`;
2. `Choice`;
3. `Parallel`;
4. `Map`;
5. `Wait`;
6. `Retry`;
7. `Catch`;
8. `Timeout`.

See `profiles/AGENT_WORKFLOW_DRAFT.md` for a non-normative YAML example.

## 20. Conformance summary

This draft defines these profiles:

| Profile | Purpose |
|---|---|
| Core Agentfile Profile | Parse and validate the Agentfile grammar |
| Security Declaration Profile | Interpret boundary declarations and fail closed |
| Cedar Policy Profile | Evaluate Cedar authorization requests |
| OCI Package Profile | Build, inspect, push, pull, and verify packages |
| Tool Projection Profile | Expose command-oriented tool/function discovery and invocation |
| Runner Conformance Profile | Execute packages while honoring declared boundaries |
| Workflow Draft Profile | Non-normative future multi-agent orchestration |

An implementation should state exactly which profiles it supports.

## 21. Security considerations

Agent packages may request dangerous capabilities. Declarative does not mean safe.

Reviewers and registries should surface:

1. requested network destinations;
2. requested mounts and mutability;
3. requested credentials;
4. high-risk tools such as shell, browser, computer-use, code execution, dynamic tool loading, and broad MCP access;
5. policy denies and permits;
6. inherited restrictions;
7. package signatures;
8. provenance;
9. vulnerability scan results;
10. unsupported runner controls.

A runner MUST NOT silently ignore security controls that the package marks as required.

## 22. Versioning

The specification uses semantic versions:

```text
MAJOR.MINOR.PATCH[-PRERELEASE]
```

Agentfiles SHOULD declare the syntax version using a parser-compatible comment:

```dockerfile
# syntax=agentrc.agentfile/v0.1
```

Packages SHOULD include the Agentfile version as package metadata and OCI annotation.

## 23. Compatibility with existing work

This draft is aligned with the current agentrc repository structure:

- `packages/aio-core` contains Agentfile parsing, types, audit, and Cedar scaffolding.
- `packages/aio-buildtime` contains compile, lock, schema, OCI, and registry package work.
- `packages/aio-runtime` contains candidate runner/driver code and enforcement experiments.
- `packages/aio` contains the CLI.
- `examples/` contains current Agentfile examples.

This draft intentionally frames runtime code as implementation and test harness, not as the definition of agentrc itself.

## 24. Open questions

1. Should `ISOLATION`, `IMAGE`, `SLICE`, `PLUGIN`, and `BACKEND` remain in Agentfile or move to a separate run manifest?
2. Should `BIND` require a mode in v0.1, or should two-argument compatibility be normative?
3. Should Cedar be mandatory for security-conformant runners or only the default policy profile?
4. Should `ALLOW`/`DENY` shorthand remain long-term or compile immediately into Cedar?
5. Should tool projection be required for all runners or only for developer/local profiles?
6. Should an agentrc package be modeled as an OCI image, OCI artifact, or both?
7. What minimal metadata should be required before publishing a package to a public registry?
8. How should multi-agent workflow policy ceilings compose across package boundaries?
9. Should `MODEL` be a directive, or should model selection remain inside `CMD`/runner configuration?
10. What should the first public conformance test suite include?

## 25. Summary

agentrc is the Agentfile and agent-package specification.

The Agentfile declares one agent. The lockfile pins its resolved dependencies. The package makes it portable. The policy makes its boundaries reviewable. The registry makes it shareable. Compatible runners execute it. agentrc itself remains the neutral declaration, packaging, and governance layer.

