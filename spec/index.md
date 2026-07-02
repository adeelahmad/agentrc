---
layout: doc
title: Specification
description: "The agentrc Agentfile specification (0.1.0-draft.5): a Dockerfile-shaped recipe for portable, governed AI agents."
permalink: /spec/
---
# agentrc Agentfile Specification

**Version:** 0.1.0-draft.5 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.agentfile/v0.1`)  
**Date:** 2026-06-30  
**Editor:** [Adeel Ahmad](https://www.linkedin.com/in/adeelahmadch)  
**Audience:** agent developers, platform / runner authors, security & compliance reviewers, registry maintainers, standards & interop implementers

> This is the single source of truth for the **Agentfile**. The normative
> keywords **MUST**, **MUST NOT**, **SHOULD**, and **MAY** follow
> [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119). Items that are genuinely
> undecided are listed in [§14 Open decisions](#14-normative-checklist-and-open-decisions)
> — they are surfaced, not silently resolved.

## Design principles

Every rule below derives from four principles:

1. **The Agentfile is Dockerfile-shaped.** It reuses standard Dockerfile
   keywords for everything it can and adds the **minimum** number of new
   keywords — just four. The new agentrc keywords are interpreted by the
   **agentrc BuildKit frontend** or the `agentrc` CLI (see [§10](#10-cli-surface)).
   The mental model, file shape, and much of the tooling transfer directly from
   Docker.
2. **`POLICY` declarations are requests, not enforcement.** A `POLICY` line is
   the developer *asking* the platform for a resource, a model, or an
   operational constraint. The platform (runtime / operator / organization) is
   free to **grant, narrow, or reject** any request. The Agentfile expresses
   *intent*; the platform holds *authority*. Enforcement is **Cedar,
   platform-side** — the platform compiles granted requests plus its own
   organization rules into [Cedar](https://www.cedarpolicy.com/) and evaluates
   them (deny-by-default, `forbid` over `permit`, tightening-only across
   `FROM`). Cedar is the enforcement engine and compilation target, **never an
   Agentfile author surface** ([§11.2](#112-cedar--the-platform-enforcement-engine-normative)).
3. **Build emits OCI labels; the platform reads labels.** At build time the
   compiler translates authored intent into namespaced OCI image labels under
   `org.agentrc.*`. The platform reads those labels at deploy / run time and
   decides what to honour, **without parsing the Agentfile source.** Labels are
   the machine-readable manifest; the Agentfile is the human-authored recipe.
4. **Resources can be embedded (cached) or fetched at runtime.** Tools, skills,
   and MCP servers MAY be embedded into the OCI artifact at build time (fast,
   reproducible, air-gappable) or resolved when the agent bootstraps (dynamic,
   fresh, requires network). Both modes MUST be supported, and an embedded
   resource MUST also be visible in labels so the platform can override or
   substitute it.

**Mount point:** the agent's projected filesystem lives under **`/mnt`**. See
[§4.1](#41-the-mnt-projection-layout).

## 1. Keyword summary

There are **four new keywords** (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`).
Everything else is a standard Dockerfile keyword, used as-is or with a
documented extension.

| Keyword | Origin | Role in agentrc |
|---|---|---|
| `FROM` | Dockerfile | Base image **and** agent inheritance (`FROM another-agent`). |
| `CMD` | Dockerfile | The agent's invocation surface — the framework / loop to run. |
| `COPY` | Dockerfile | Add **local** tools, skills, MCP bundles, or an SOP file into the `/mnt` tree. |
| `ADD` | Dockerfile (extended) | Add **remote** tools, skills, MCP servers, or SOP via `--remote` plus delivery flags. |
| `HEALTHCHECK` | Dockerfile | Liveness probe; MAY invoke a projected tool. |
| `LABEL` | Dockerfile | Standard OCI metadata; available for hand-authored `org.agentrc.*` metadata. |
| `ENV` / `ARG` / `WORKDIR` / `USER` / `EXPOSE` / `RUN` | Dockerfile | Standard semantics; available, unchanged. |
| **`IDENTITY`** | **New** | Who / what the agent is: name, version, author, description. |
| **`CAPABILITY`** | **New** | What modalities / patterns the agent supports (text, streaming, multimodal, …). |
| **`SOP`** | **New** | The agent's system prompt / objective / standard operating procedure. |
| **`POLICY`** | **New** | Typed, namespaced **request** for a resource, model, or operational constraint. |

> **There is no `TOOL`, `MCP`, `SERVER`, `FUNC`, `CRED`, `SECRET`, `AUDIT`,
> `MOUNT`, `MEMORY`, or `RATELIMIT` keyword.** Tools / skills / MCP are added with
> `COPY` / `ADD`; memory / context / storage / CPU / model are `POLICY` requests;
> audit rides on `agent.hooks.*`; secrets are [deferred](#121-secrets-deferred).
> Any earlier draft that used those keywords is **stale**.

> **A2A (agent-to-agent: Agent Cards, discovery, delegation) is out of scope for
> this version.** Capability *exposure* is handled via `IDENTITY` / `CAPABILITY`
> / labels; the agent-to-agent *protocol* (how one agent discovers and calls
> another) is [deferred](#143-deferred-to-a-later-version).

## 2. `FROM` — base image and agent inheritance

```dockerfile
FROM <image-or-agent>[:<tag>] [AS <name>]
```

Exactly as in Dockerfile, every Agentfile MUST contain a `FROM` instruction, and
`FROM` must be the first instruction after the `# syntax=` line, comments, and any
`ARG` that `FROM` consumes.

`FROM` behaves like Dockerfile `FROM` for plain base images, and additionally
supports **agent inheritance**: when the source is another agentrc agent, the
child inherits the parent's capabilities and policy under these rules:

- **Capabilities compose additively.** The child receives all of the parent's
  tools, skills, and MCP servers, and MAY add more.
- **Policy composes by tightening only.** The child MAY narrow or extend within
  the parent's declared requests, but MUST NOT widen past a constraint the
  parent locked. The runtime realization of this rule is Cedar's monotonic
  intersection across `FROM` ([§11.2](#112-cedar--the-platform-enforcement-engine-normative)).

```dockerfile
# Plain base image
FROM python:3.11-slim

# Agent inheritance
FROM ghcr.io/acme/pii-redacted-base:1.4
```

## 3. `CMD` — the agent invocation surface

```dockerfile
CMD <command>
```

The command that starts the agent: the framework, loop, or runtime to execute.
agentrc is framework-neutral.

```dockerfile
CMD claude --print
CMD crewai run --agents researcher,writer
CMD python -m my_agent.main
```

## 4. Adding tools, skills, MCP servers, and SOP files

Tools, skills, MCP servers, and a file-backed SOP are **files placed into the
`/mnt` tree**. Local sources use `COPY`; remote sources use `ADD --remote`. The
**destination path** under `/mnt` determines what the resource is — there is no
dedicated keyword per resource type (except `SOP`, which has its own keyword for
the inline / heredoc forms; see [§7](#7-sop--the-agents-system-prompt--objective)).

### 4.1 The `/mnt` projection layout

The compiler and runtime MUST use this layout:

```text
/mnt
├── tools/     # each tool is an executable; argv in, structured output out
├── skills/    # skill bundles (a directory of instructions/scripts/resources)
├── mcp/       # MCP server bundles or configs, projected as tool directories
├── proc/      # runtime-populated: live policy, identity, budgets, audit tail
└── SOP        # the agent's system prompt, as a readable file (see §7)
```

A **tool** is an executable file under `/mnt/tools/`. To be self-describing, the
executable SHOULD expose its schema via `--agentrc-schema` (prints JSON schema
to stdout) **or** ship a sibling `<tool>.toolspec.json`. The projection layer
reads this to render `man` / `cat` output and to enumerate capability. (Which
convention is canonical is [open decision #3](#142-open-decisions-surface-these-do-not-silently-resolve).)

### 4.2 `COPY` — local resources

Standard Dockerfile `COPY`, including `--chmod` and `--chown`.

```dockerfile
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
COPY --chmod=755 ./tools/shell     /mnt/tools/shell
COPY ./skills/code-review          /mnt/skills/code-review
COPY ./mcp/github                  /mnt/mcp/github
COPY ./sop.md                      /mnt/SOP
```

### 4.3 `ADD --remote` — remote resources

`ADD` is extended with a `--remote` flag plus **delivery flags** controlling
*when* the resource is resolved and *what happens if it cannot be.* This follows
existing Dockerfile precedent: `HEALTHCHECK` already takes flags before its
`CMD`, and `RUN` / `COPY` / `ADD` already take flags (`--mount`, `--chmod`,
`--chown`). `ADD --remote` generalizes that.

```dockerfile
ADD --remote [--cached | --runtime] [--fail-if-unavailable | --warn-if-unavailable] \
    [--chmod=<mode>] [--chown=<user>:<group>] \
    <remote-source> <destination-under-/mnt>
```

**Delivery flags (the fstab-style contract):**

| Flag | Meaning | Default |
|---|---|---|
| `--cached` | Fetch **at build time** and embed as an OCI layer. Fast startup, reproducible, air-gappable. Still recorded in labels so the platform can override / substitute. | **Default** if neither `--cached` nor `--runtime` is given. |
| `--runtime` | Do **not** embed. Record the reference in labels and fetch when the agent **bootstraps**. Dynamic, always-fresh; needs network at run time. Use for internal / air-gapped MCP servers resolved on the target network. | — |
| `--fail-if-unavailable` | If the resource cannot be resolved (at build for `--cached`, at bootstrap for `--runtime`), **abort** — fail the build or refuse to boot. | **Default** failure mode. |
| `--warn-if-unavailable` | If the resource cannot be resolved, **log a warning and continue**. | — |

`--chmod` / `--chown` carry their standard Dockerfile meaning.

```dockerfile
# Tool, remote, embedded by default
ADD --remote --chmod=755 https://registry.agentrc.io/tools/http_get:latest /mnt/tools/http_get

# Skill, remote, explicitly cached, hard-fail if missing
ADD --remote --cached --fail-if-unavailable \
    https://registry.agentrc.io/skills/code-review:1.2.3 /mnt/skills/code-review

# MCP server, internal, fetched at runtime, hard-fail if registry unreachable
ADD --remote --runtime --fail-if-unavailable \
    mcp://registry.internal.acme/servers/github:latest /mnt/mcp/github

# MCP server, embedded for speed, but tolerate absence
ADD --remote --cached --warn-if-unavailable \
    https://registry.agentrc.io/mcp/web-search:latest /mnt/mcp/web-search

# SOP, remote file
ADD --remote https://registry.agentrc.io/sops/claims-triage:latest /mnt/SOP
```

### 4.4 Embedding MUST stay visible and overridable

When a tool / skill / MCP server is embedded (`--cached`), the compiler MUST
**also** emit a label recording its origin and resolved digest, so the platform
or organization can override or substitute it at deploy time (e.g., redirect a
public MCP server to an internal mirror) without rebuilding. See [§9.3](#93-resource-delivery--layers--labels).

## 5. `IDENTITY` — who the agent is

```dockerfile
IDENTITY <key>=<value> [<key>=<value> ...]
```

Declares the agent's identity as `key=value` pairs. Repeatable across lines.
Recognised keys: `name`, `version`, `author`, `description`. Additional keys MAY
be added (extensible).

```dockerfile
IDENTITY name=claims-triage version=1.0 author=acme
IDENTITY description="Triages insurance claims by severity and routes escalations"
```

**Build translation:** each key emits a label `org.agentrc.identity.<key>=<value>`
directly (these are short).

```text
org.agentrc.identity.name=claims-triage
org.agentrc.identity.version=1.0
org.agentrc.identity.author=acme
org.agentrc.identity.description=Triages insurance claims by severity and routes escalations
```

## 6. `CAPABILITY` — what the agent supports

```dockerfile
CAPABILITY <capability>
```

Declares one supported modality / pattern per line. Example values: `text`,
`multimodal`, `streaming`, `function-calling`, `agentic`, `vision`, `audio`. The
value set is open / extensible.

```dockerfile
CAPABILITY text
CAPABILITY streaming
CAPABILITY function-calling
```

**Build translation:** each line emits `org.agentrc.capability.<value>=true`.
(Whether a single comma-joined label is used instead is
[open decision #5](#142-open-decisions-surface-these-do-not-silently-resolve);
pick one and be consistent.)

```text
org.agentrc.capability.text=true
org.agentrc.capability.streaming=true
org.agentrc.capability.function-calling=true
```

## 7. `SOP` — the agent's system prompt / objective

`SOP` carries the agent's system prompt, objective, constraints, and voice — the
agent's standard operating procedure. It supports **three forms**:

**Inline (single line):**

```dockerfile
SOP You are a claims-triage specialist. Be thorough but fast; escalate anything ambiguous.
```

**Heredoc (multi-line):**

```dockerfile
SOP <<EOF
You are a claims-triage specialist.
Prioritize by severity. Escalate anything ambiguous to a human.
Never fabricate policy numbers.
EOF
```

**File-backed (reuses `COPY` / `ADD`, no new mechanism):**

```dockerfile
COPY ./sop.md /mnt/SOP
# or
ADD --remote https://registry.agentrc.io/sops/claims-triage:latest /mnt/SOP
```

**Build translation.** Regardless of form, the compiler MUST embed the SOP as a
**readable file at `/mnt/SOP`** (the agent reads it at startup). Because an SOP
can be large, the compiler MUST NOT inline the full text into a label; instead
it emits a **pointer plus digest**:

```text
org.agentrc.sop=/mnt/SOP
org.agentrc.sop.sha256=<digest>
```

so the platform can see that an SOP exists, verify it, or override the file,
without carrying the full prompt in label metadata.

## 8. `POLICY` — typed resource, model, and operational requests

`POLICY` exists as a distinct keyword (rather than raw `LABEL`s) so the request
vocabulary is **standard and discoverable** instead of every team inventing its
own label conventions.

```dockerfile
POLICY <namespaced.key> <value>
```

Each `POLICY` line is a **single request** to the platform. Requests are
**independent** (no block syntax, no `END`), **typed** by a dotted namespace, and
**independently honourable, narrowable, or rejectable**. The developer writes the
**short form** (no `org.agentrc.` prefix); the compiler prepends it when emitting
labels ([§9.1](#91-policy--labels)). Top-level namespaces: `agent.*`,
`substrate.*`, `model.*`, and `network`. All are extensible.

### 8.1 `agent.*` — agent-side operational constraints, state, lifecycle

| Key | Meaning | Example |
|---|---|---|
| `agent.idle_timeout` | Max idle time before the agent is considered done / reaped. | `POLICY agent.idle_timeout 5m` |
| `agent.tool_timeout` | Max wall-clock time for a single tool invocation. | `POLICY agent.tool_timeout 30s` |
| `agent.max_retries` | Max retries the agent loop should attempt on a failed step. | `POLICY agent.max_retries 3` |
| `agent.context` | Requested model context window (tokens); platform honours or caps. | `POLICY agent.context 100k` |
| `agent.context.type` | Context strategy: `autocompressed`, `sliding-window`, `summarized`, `full`. | `POLICY agent.context.type autocompressed` |
| `agent.memory.short` | Agent's short-term / scratch / conversation buffer. | `POLICY agent.memory.short 512mb` |
| `agent.memory.long` | Agent's long-term / persistent knowledge cache. | `POLICY agent.memory.long 2gb` |
| `agent.sub_agents` | Whether the agent may spawn / delegate to sub-agents. | `POLICY agent.sub_agents true` |
| `agent.sub_agents.max` | Max number of sub-agents permitted. | `POLICY agent.sub_agents.max 5` |
| `agent.sub_agent_timeout` | Max wall-clock time for a delegated sub-agent. | `POLICY agent.sub_agent_timeout 60s` |
| `agent.interrupt_endpoint` | Endpoint the agent uses to reach the end user (human-in-the-loop). | `POLICY agent.interrupt_endpoint https://user-service.internal/interrupt` |
| `agent.hooks.pre` | Webhook called before each agent step. | `POLICY agent.hooks.pre https://hooks.internal/pre-step` |
| `agent.hooks.post` | Webhook called after each agent step. | `POLICY agent.hooks.post https://hooks.internal/post-step` |
| `agent.hooks.on_error` | Webhook called on a step error. | `POLICY agent.hooks.on_error https://hooks.internal/error` |
| `agent.hooks.on_tool_call` | Webhook called on each tool invocation (e.g., audit). | `POLICY agent.hooks.on_tool_call https://hooks.internal/audit` |

The `agent.*` namespace is **extensible**: new agent patterns add new keys
(`agent.telemetry_sink`, `agent.feedback_channel`, `agent.audit_log`,
`agent.voice.audio_endpoint`, …) without changing the grammar. Document the
convention, not a closed list.

### 8.2 `substrate.*` — substrate / runtime resource requests

| Key | Meaning | Example |
|---|---|---|
| `substrate.runtime.memory` | Runtime memory the platform should allocate. | `POLICY substrate.runtime.memory 8gb` |
| `substrate.runtime.cpu` | CPU units / cores requested. | `POLICY substrate.runtime.cpu 4` |
| `substrate.device` | A device the agent requests access to. | `POLICY substrate.device /dev/gpu` |
| `substrate.init` | The init / entry process inside the runtime. | `POLICY substrate.init bash` |
| `substrate.shell` | Explicit shell path for the runtime. | `POLICY substrate.shell /bin/bash` |
| `substrate.ptty` | Whether to allocate a pseudo-terminal. | `POLICY substrate.ptty true` |

Extensible (`substrate.storage`, `substrate.runtime.gpu`, …).

### 8.3 `model.*` — model and model-capability requests

| Key | Meaning | Example |
|---|---|---|
| `model.name` | Requested model. | `POLICY model.name claude-opus-4` |
| `model.min_context` | Minimum context window the model must support. | `POLICY model.min_context 200k` |
| `model.capability` | A required model capability (repeatable). | `POLICY model.capability vision` |
| `model.fallback` | Fallback model if the primary is unavailable. | `POLICY model.fallback claude-sonnet-4` |
| `model.temperature` | Requested sampling temperature. | `POLICY model.temperature 0.7` |

The developer says *"I need a model with vision and 200k context, prefer Opus,
fall back to Sonnet"*; the platform honours, substitutes, or rejects based on
what it has.

### 8.4 `network` — egress requests

```dockerfile
POLICY network dns:<host>:<port>
```

Declares an outbound host / port the agent requests. `*` is permitted for
wildcard ports (and MAY be used for host wildcards within a domain). The platform
grants, narrows, or denies.

```dockerfile
POLICY network dns:api.github.com:443
POLICY network dns:internal.acme:*
```

### 8.5 Auto-derived egress from hook / interrupt URLs

Any `POLICY` value that is a URL the agent will call out to — `agent.hooks.*` and
`agent.interrupt_endpoint` — MUST cause the compiler to **auto-derive a
corresponding `network` egress request** for that host, emitted as an explicit,
**attributed** label so it is never a silent network hole:

```text
org.agentrc.network.dns.hooks.internal=443
org.agentrc.network.dns.hooks.internal.source=auto:agent.hooks.pre
```

The platform still must grant the derived egress; auto-derivation is an ergonomic
convenience, not an implicit grant. Deny-by-default still applies (see [§11](#11-runtime-behaviour-what-the-platform-does)).

### 8.6 Why `POLICY` is a keyword and not raw labels

If developers wrote requests as ad-hoc `LABEL`s, every team would invent its own
conventions and nothing would be standard. `POLICY` gives one canonical, typed
surface; the compiler maps it to the `org.agentrc.*` label space so the platform
sees uniform, machine-readable intent.

### 8.7 `substrate.<platform>.*` — platform-scoped substrate requests

Some requests only make sense on a specific execution platform. The
`substrate.<platform>.*` namespace lets an author pin a request to one platform
without leaking that detail into the generic, substrate-neutral `substrate.*`
keys. The platform token is one of `aws | gcp | azure | kubernetes | local`.

Unknown platform tokens MUST parse — an unrecognised token is **never a hard
error**. A key scoped to a platform other than the one currently translating is
**foreign** and is simply ignored (a linter MAY warn, but MUST NOT fail). This
keeps a single Agentfile portable: `substrate.aws.*` and `substrate.gcp.*` can
coexist, and each platform reads only its own keys.

Platform-scoped keys are emitted as labels the same way as any other request:

```text
org.agentrc.substrate.<platform>.<key>=<value>
```

On a given platform, a platform-scoped key **beats** the generic `substrate.*`
key for the same concept — but only on that platform. As with every namespace,
requests are **tightening-only across `FROM`**: a child Agentfile may narrow an
inherited value, never loosen it.

`substrate.<platform>.*` is a **key namespace under the existing `substrate.*`**
request surface — it is not a rename of `substrate.*` and it does not introduce a
new keyword. `substrate.kubernetes.serviceAccount`, for example, is a *key*, not
a separate namespace.

**AWS key registry.** For `substrate.aws.*` the following keys are defined:

| Key | Meaning |
|---|---|
| `roleArn` | Execution role ARN the runtime assumes. |
| `networkMode` | Networking mode for the runtime. |
| `securityGroup` | Security group to attach (repeatable). |
| `subnet` | Subnet to place the runtime in (repeatable). |
| `protocol` | Server protocol the runtime speaks. |
| `maxLifetime` | Maximum lifetime before the runtime is reaped. |
| `deployment.mode` | `container` (default) or `code`. |
| `code.s3.uri` | S3 URI of the code bundle when `deployment.mode` is `code`. |

```dockerfile
POLICY substrate.aws.roleArn arn:aws:iam::123456789012:role/agent-exec
POLICY substrate.aws.networkMode PUBLIC
POLICY substrate.aws.securityGroup sg-0abc123
POLICY substrate.aws.subnet subnet-0def456
POLICY substrate.aws.protocol HTTP
POLICY substrate.aws.maxLifetime 1h
POLICY substrate.aws.deployment.mode code
POLICY substrate.aws.code.s3.uri s3://acme-agents/code-reviewer.zip
```

### 8.8 `agent.auth.*` — invocation authorization

`agent.auth.*` declares how the platform should authorize **calls into** the
running agent's invocation endpoint. It is generic authorization configuration —
**not a secret** (secrets remain deferred; there is no secret keyword) and not
inline policy code.

| Key | Meaning |
|---|---|
| `agent.auth.mode` | Authorizer to enforce: `platform` (default) | `jwt` | `none`. |
| `agent.auth.jwt.discovery_url` | OIDC discovery URL for the JWT authorizer. |
| `agent.auth.jwt.allowed_audience` | Permitted token audience (repeatable). |
| `agent.auth.jwt.allowed_client` | Permitted client identifier (repeatable). |

This namespace is **fail-closed**: a platform that cannot enforce a requested
`jwt` authorizer **MUST NOT expose the invocation endpoint**. It is safer to
refuse invocation than to expose an agent behind an authorizer the platform
cannot honour.

```dockerfile
POLICY agent.auth.mode jwt
POLICY agent.auth.jwt.discovery_url https://auth.acme/.well-known/openid-configuration
POLICY agent.auth.jwt.allowed_audience agentrc://code-reviewer
POLICY agent.auth.jwt.allowed_client acme-ci
```

### 8.9 `substrate.runtime.language` — runtime language / version

`substrate.runtime.language` requests the language runtime the agent's code
expects, written as `<language>:<version>` (for example `python:3.11`). It is
**optional**.

In **container mode** the base image is authoritative, so the platform MAY ignore
this key. In **code mode** the platform requires it — either stated explicitly or
resolvable by inference — otherwise it MUST **fail closed** rather than guess a
runtime.

```dockerfile
POLICY substrate.runtime.language python:3.11
```

## 9. Build-time translation: Agentfile → OCI labels

`agentrc build` MUST translate authored intent into namespaced OCI image labels.
The platform reads **labels**, not the Agentfile.

### 9.1 `POLICY` → labels

Each `POLICY <key> <value>` becomes `org.agentrc.<key>=<value>`.

| Authored | Emitted label |
|---|---|
| `POLICY agent.idle_timeout 5m` | `org.agentrc.agent.idle_timeout=5m` |
| `POLICY agent.context.type autocompressed` | `org.agentrc.agent.context.type=autocompressed` |
| `POLICY agent.hooks.pre https://hooks.internal/pre-step` | `org.agentrc.agent.hooks.pre=https://hooks.internal/pre-step` (+ auto-derived egress, [§8.5](#85-auto-derived-egress-from-hook--interrupt-urls)) |
| `POLICY agent.sub_agents.max 5` | `org.agentrc.agent.sub_agents.max=5` |
| `POLICY substrate.runtime.memory 8gb` | `org.agentrc.substrate.runtime.memory=8gb` |
| `POLICY substrate.ptty true` | `org.agentrc.substrate.ptty=true` |
| `POLICY model.name claude-opus-4` | `org.agentrc.model.name=claude-opus-4` |
| `POLICY model.capability vision` | `org.agentrc.model.capability.vision=true` |
| `POLICY network dns:api.github.com:443` | `org.agentrc.network.dns.api.github.com=443` |

### 9.2 `IDENTITY` / `CAPABILITY` / `SOP` → labels

| Authored | Emitted label(s) |
|---|---|
| `IDENTITY name=claims-triage version=1.0` | `org.agentrc.identity.name=claims-triage`, `org.agentrc.identity.version=1.0` |
| `CAPABILITY streaming` | `org.agentrc.capability.streaming=true` |
| `SOP …` / `COPY ./sop.md /mnt/SOP` | `org.agentrc.sop=/mnt/SOP`, `org.agentrc.sop.sha256=<digest>` (pointer + digest, never full text) |

### 9.3 Resource delivery → layers + labels

| Authored | Build behaviour | Emitted label(s) |
|---|---|---|
| `COPY ./tools/x /mnt/tools/x` | Embed file as a layer. | `org.agentrc.tool.x=local` |
| `ADD --remote --cached <url> /mnt/tools/x` | Fetch at build, embed as a layer. | `org.agentrc.tool.x=<digest>` + `org.agentrc.tool.x.origin=<url>` |
| `ADD --remote --runtime <url> /mnt/mcp/x` | Do not embed; record reference. | `org.agentrc.mcp.x=runtime:<url>` |

For every embedded MCP server or skill, the compiler MUST emit both the
**resolved digest** and the **origin reference**, e.g.:

```text
org.agentrc.mcp.github=sha256:abc123...
org.agentrc.mcp.github.origin=https://registry.agentrc.io/mcp/github:latest
```

so a platform or organization MAY rewrite `org.agentrc.mcp.github.origin` to an
internal mirror at deploy time without rebuilding.

### 9.4 Policy encoding — both forms supported

The compiled request set MUST be retrievable by the platform in **either** form,
and the build MAY support both:

- **Inline**: values live directly in labels (as above). Best for small request
  sets.
- **By digest / reference**: labels carry a digest of a structured manifest
  embedded as a layer (or referenced externally); the full manifest is fetched
  from there. Best for large request sets.

Which is the **default** for `arc build` is
[open decision #1](#142-open-decisions-surface-these-do-not-silently-resolve);
the `--policy-mode inline|digest` flag is the seam.

## 10. CLI surface

agentrc ships **two** build paths. The four agentrc keywords (`IDENTITY`,
`CAPABILITY`, `SOP`, `POLICY`) and the `ADD --remote` extension are interpreted by
the agentrc frontend / CLI; a stock `docker build` **without** the frontend
understands only the standard Dockerfile keywords. See the [CLI page](/cli/) for
detail.

### 10.1 BuildKit frontend (no new tool to install)

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM ...
```

```bash
docker build -f Agentfile -t ghcr.io/acme/my-agent:1.0 .
```

The `# syntax=` directive routes the Agentfile through the agentrc frontend
image, which parses the agentrc keywords, compiles to LLB, and produces the OCI
artifact with all `org.agentrc.*` labels. The frontend image is **published** at
`ghcr.io/adeelahmad/agentrc-frontend`, so a `# syntax=ghcr.io/adeelahmad/agentrc-frontend`
directive routes `docker build` through it directly; see the
[CLI page](/cli/) for current status. The native `arc run`/`sign`/`verify`
commands remain `planned`.

### 10.2 Native `agentrc` CLI

Primary command `agentrc`, short alias `arc`.

```bash
agentrc build  [-f Agentfile] [-t <ref>] [--policy-mode inline|digest] .
agentrc push   <ref>
agentrc pull   <ref>
agentrc run    <ref> --backend local|bedrock|kubernetes [per-backend flags]
```

| Command | Purpose |
|---|---|
| `agentrc build` (`arc build`) | Compile an Agentfile to a signed OCI artifact, emitting `org.agentrc.*` labels and embedding `--cached` resources. |
| `agentrc push` | Push the artifact to any OCI registry. |
| `agentrc pull` | Pull an artifact. |
| `agentrc run` | Run an artifact on a chosen backend. Substrate is a **run-time** choice (`--backend` / `--isolation`); it is **not** an Agentfile directive. |

The build output of 10.1 and 10.2 MUST be **identical** OCI artifacts — the
frontend and the CLI are two front doors to the same compiler.

## 11. Runtime behaviour (what the platform does)

### 11.1 Decision flow

At deploy / run time the platform:

1. Pulls the OCI artifact.
2. Reads all `org.agentrc.*` labels from the image config — **without parsing the
   Agentfile.**
3. Evaluates each request (`POLICY`-derived label, including auto-derived egress)
   against organization / platform policy and available resources, then
   **grants, narrows, or rejects** it. Rejection / narrowing SHOULD be auditable
   ([open decision #4](#142-open-decisions-surface-these-do-not-silently-resolve)).
4. Resolves any credential the agent references via a **platform-defined**
   mechanism (secrets are [deferred](#121-secrets-deferred); the binding is TBD).
5. For `--runtime` resources, fetches them now; applies `--fail-if-unavailable` /
   `--warn-if-unavailable`.
6. May substitute embedded resources by honouring overridden `*.origin` labels.
7. Loads the SOP from `/mnt/SOP`, selects / validates the model from `model.*`,
   and boots the agent (`CMD`) on the chosen substrate with the granted
   constraints, projecting `/mnt/tools`, `/mnt/skills`, `/mnt/mcp`, and
   populating `/mnt/proc`.

> **Reminder:** `POLICY` lines are *requests*. The platform is the authority.
> "Deny by default" applies to the platform's grant decision — an unrecognised or
> disallowed request (or an auto-derived egress) is not silently honoured.

### 11.2 Cedar — the platform enforcement engine (normative)

Cedar is **not** an Agentfile author surface and MUST NOT appear in the
Agentfile. Authors speak only typed `POLICY` requests ([§8](#8-policy--typed-resource-model-and-operational-requests)).
Cedar is the **platform-side enforcement engine and compilation target** for
those requests: the normative answer to *"how does a conformant platform decide
and enforce grants identically?"*

The relationship is a compilation, not an abandonment:

```text
author writes          platform compiles to        platform enforces
─────────────          ────────────────────        ─────────────────
POLICY request   ──►   Cedar entities + policies    deny-by-default,
(org.agentrc.*)        (request + org rules)         forbid > permit,
                                                      order-independent,
                                                      monotonic across FROM
```

A **conformant platform** MUST, for each granted request, derive Cedar entities /
actions and evaluate them under Cedar's semantics. The principal is the agent
identity (`org.agentrc.identity.name`); the action and resource are derived from
the request namespace:

| Request label | Cedar action | Cedar resource |
|---|---|---|
| `org.agentrc.network.dns.<host>=<port>` | `Action::"NetworkEgress"` | `Host::"<host>:<port>"` |
| `org.agentrc.tool.<name>` (a projected tool) | `Action::"tool.invoke"` | `Tool::"<name>"` |
| `org.agentrc.mcp.<name>` | `Action::"mcp.request"` | `MCPServer::"<name>"` |
| `org.agentrc.agent.sub_agents=true` | `Action::"agent.delegate"` | `Agent::*` (capped by `sub_agents.max`) |
| `org.agentrc.substrate.device=<dev>` | `Action::"device.access"` | `Device::"<dev>"` |

**Enforcement properties a conformant platform MUST preserve** (these are
Cedar's, and they are why Cedar is the engine):

- **Deny-by-default.** Absence of a grant is a denial. An unrecognised request, an
  un-granted auto-derived egress, or an action with no matching `permit` is
  denied.
- **`forbid` overrides `permit`, order-independently.** An organization `forbid`
  (e.g., "no agent may reach the public internet") MUST defeat any agent request,
  regardless of evaluation order.
- **Monotonic composition across `FROM`.** When an agent inherits (`FROM
  another-agent`), the effective authorization is the **intersection** of
  ceilings: a child's granted set MUST NOT exceed its parent's. A parent `forbid`
  is un-looseni-able by the child. This is the runtime realization of
  [§2](#2-from--base-image-and-agent-inheritance)'s tightening-only rule.

**Two-tier governance, one author surface.** The agent's `POLICY` requests are
the *floor of intent*; the organization's own Cedar policies are the *ceiling of
authority*. The platform compiles both into one Cedar `PolicySet` and evaluates
the grant. The author never writes Cedar; the organization's security team owns
the org-level Cedar policies out-of-band. This keeps exactly one governance
surface in the Agentfile (typed requests) and one enforcement engine beneath it
(Cedar), with a single normative mapping between them — and no second
author-facing policy language. The [Enforcement profile](/profiles/security/)
defines platform conformance against this mapping.

## 12. Secrets and `HEALTHCHECK`

### 12.1 Secrets (deferred)

Secret resolution is **out of scope for this draft** and intentionally
unspecified. It is a subsystem in its own right (Vault, secret brokers, env
injection, OIDC / workload identity) and will be designed separately — the
Agentfile is **not** pre-committed to any one resolution model (no host-scoped
broker, no fixed label schema) yet.

For now, an agent that needs a credential MAY reference it by name and leave
**all** resolution to the platform / runner; the binding mechanism is TBD. See
[§14.3](#143-deferred-to-a-later-version).

> **Implementer note:** do not add a `SECRET` / `CRED` keyword or an
> `org.agentrc.secret.*` label schema in this pass. Secrets are a future design
> (Vault / broker / workload-identity integration), not a settled part of the
> spec.

### 12.2 `HEALTHCHECK`

Standard Dockerfile `HEALTHCHECK`; the probe MAY invoke a projected tool.

```dockerfile
HEALTHCHECK --interval=60s --timeout=15s --retries=3 CMD /mnt/tools/file_read --agentrc-schema
```

## 13. Complete worked example

```dockerfile
# syntax=agentrc.agentfile/v0.1

# --- Base / inheritance -----------------------------------------------------
FROM ghcr.io/acme/pii-redacted-base:1.4

# --- Identity ---------------------------------------------------------------
IDENTITY name=claims-triage version=1.0 author=acme
IDENTITY description="Triages insurance claims by severity and routes escalations"

# --- Capabilities -----------------------------------------------------------
CAPABILITY text
CAPABILITY streaming
CAPABILITY function-calling

# --- System prompt / objective ----------------------------------------------
SOP <<EOF
You are a claims-triage specialist.
Prioritize by severity. Escalate anything ambiguous to a human.
Never fabricate policy numbers.
EOF

# --- Invocation surface (any framework) -------------------------------------
CMD claude --print

# --- Tools (local, embedded) ------------------------------------------------
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
COPY --chmod=755 ./tools/shell     /mnt/tools/shell

# --- Tools (remote, embedded by default) ------------------------------------
ADD --remote --chmod=755 \
    https://registry.agentrc.io/tools/http_get:latest /mnt/tools/http_get

# --- Skills (remote, cached, hard-fail if missing) --------------------------
ADD --remote --cached --fail-if-unavailable \
    https://registry.agentrc.io/skills/code-review:1.2.3 /mnt/skills/code-review

# --- MCP server (internal, fetched at runtime, hard-fail if unreachable) -----
ADD --remote --runtime --fail-if-unavailable \
    mcp://registry.internal.acme/servers/github:latest /mnt/mcp/github

# --- Secrets: DEFERRED — credential resolution is platform-defined (TBD, see §12.1) ---

# --- Policy: model requests -------------------------------------------------
POLICY model.name        claude-opus-4
POLICY model.min_context 200k
POLICY model.capability  vision
POLICY model.fallback    claude-sonnet-4

# --- Policy: agent-side requests --------------------------------------------
POLICY agent.idle_timeout       5m
POLICY agent.tool_timeout       30s
POLICY agent.max_retries        3
POLICY agent.context            100k
POLICY agent.context.type       autocompressed
POLICY agent.memory.short       512mb
POLICY agent.memory.long        2gb
POLICY agent.sub_agents         true
POLICY agent.sub_agents.max     5
POLICY agent.interrupt_endpoint https://user-service.internal/interrupt
POLICY agent.hooks.pre          https://hooks.internal/pre-step
POLICY agent.hooks.on_tool_call https://hooks.internal/audit

# --- Policy: substrate-side requests ----------------------------------------
POLICY substrate.runtime.memory 8gb
POLICY substrate.runtime.cpu    4
POLICY substrate.device         /dev/gpu
POLICY substrate.init           bash
POLICY substrate.ptty           true

# --- Policy: network egress requests ----------------------------------------
POLICY network dns:api.github.com:443
POLICY network dns:internal.acme:*

# --- Liveness ---------------------------------------------------------------
HEALTHCHECK --interval=60s --timeout=15s --retries=3 CMD /mnt/tools/file_read --agentrc-schema
```

Build and run:

```bash
arc build -t ghcr.io/acme/claims-triage:1.0 .
arc push  ghcr.io/acme/claims-triage:1.0
arc run   ghcr.io/acme/claims-triage:1.0 --isolation microvm
```

## 13a. Packaging, registry, and conformance

The Agentfile compiles to an ordinary **OCI artifact**: standard layers carry the
embedded resources (`/mnt/tools`, `/mnt/skills`, `/mnt/mcp`, `/mnt/SOP`), and the
image config carries the `org.agentrc.*` labels. That means agents push, pull,
sign, and mirror through any OCI-compatible registry exactly like container
images — see the [OCI labels &amp; package profile](/profiles/oci-package/). The
agent's projected filesystem (`/mnt`) is described by the
[projection profile](/profiles/tool-projection/), and what a conformant platform
MUST do with the labels is the [platform conformance profile](/profiles/runner-conformance/),
verified against the adversarial [conformance suite](/docs/conformance/).

## 13b. Who this serves

| You are… | agentrc gives you… |
|---|---|
| An **agent developer / adopter** | A Dockerfile-shaped recipe — four new keywords over keywords you already know — that builds with `docker build` or `arc build`. |
| A **platform / runner author** | A labels-only contract: read `org.agentrc.*`, grant / narrow / reject, enforce with Cedar, fail closed. No need to parse the Agentfile. |
| A **security / compliance reviewer** | One artifact whose labels state every request: tools, network, model, sub-agents, lifecycle — vetted before it runs. |
| A **registry maintainer** | A standard OCI artifact with digests and `.origin` labels you can mirror, sign ([Sigstore](https://www.sigstore.dev/)), and attest ([SLSA](https://slsa.dev/)). |
| A **standards / interop implementer** | A small grammar, a fixed label namespace, and a normative request → Cedar mapping to build against. |

## 14. Normative checklist and open decisions

### 14.1 Do not deviate

- **Four** new keywords: `IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`. Everything
  else is standard Dockerfile.
- No `TOOL` / `MCP` / `SERVER` / `FUNC` / `CRED` / `SECRET` / `AUDIT` / `MOUNT` /
  `MEMORY` / `RATELIMIT` keywords.
- Mount point is **`/mnt`** (`/mnt/tools`, `/mnt/skills`, `/mnt/mcp`,
  `/mnt/proc`, `/mnt/SOP`).
- Tools / skills / MCP / SOP-file added via `COPY` (local) or `ADD --remote`
  (remote); destination path under `/mnt` determines type.
- `ADD --remote` supports `--cached` (default) / `--runtime`,
  `--fail-if-unavailable` (default) / `--warn-if-unavailable`, plus `--chmod` /
  `--chown`.
- `SOP` supports inline, heredoc, and file-backed forms; always embedded at
  `/mnt/SOP`; label is a **pointer + digest**, never the full text.
- `IDENTITY` = `key=value` pairs → `org.agentrc.identity.*`. `CAPABILITY` = one
  per line → `org.agentrc.capability.*`.
- `POLICY` short form authored; compiler prepends `org.agentrc.`. Namespaces:
  `agent.*`, `substrate.*`, `model.*`, `network` — all extensible.
- Hook / interrupt URLs auto-derive an **explicit, attributed** `network` egress
  label; never a silent hole.
- **Secrets are deferred** — no `SECRET` / `CRED` keyword, no
  `org.agentrc.secret.*` schema in this pass
  ([§12.1](#121-secrets-deferred)); credential resolution (Vault / broker /
  workload-identity) is a future design.
- **No `AUDIT` keyword.** Audit rides on `agent.hooks.on_tool_call` (a developer
  sink) or a future `agent.audit` posture request — never its own keyword.
- Embedded MCP / skills emit both digest and `.origin` labels so platforms can
  override.
- `POLICY` is a **request**; the platform grants / narrows / rejects.
- **Cedar is the platform enforcement engine, never an author surface.** No Cedar
  `permit` / `forbid` in the Agentfile; the platform compiles typed requests (+
  org rules) into Cedar and enforces deny-by-default, `forbid` > `permit`, and
  tightening-only `FROM` composition ([§11.2](#112-cedar--the-platform-enforcement-engine-normative)).
- Two build paths (BuildKit frontend + `agentrc` / `arc` CLI) produce identical
  artifacts.
- **A2A is out of scope for this version.** Do not add agent-to-agent protocol
  directives or labels.

### 14.2 Open decisions (surface these; do not silently resolve)

1. **Policy encoding default.** Inline-in-labels vs digest / referenced manifest —
   both supported; which is the **default** for `arc build` is undecided
   (`--policy-mode` flag is the seam).
2. **Cedar-grade conditional authorization — DECIDED** (see [§11.2](#112-cedar--the-platform-enforcement-engine-normative)).
   Conditional / fine-grained authorization is **platform-side only**: Cedar is
   the platform's enforcement engine and the compilation target for typed
   `POLICY` requests. The Agentfile author surface does **not** accept Cedar
   `permit` / `forbid`. One author surface (typed requests), one enforcement
   engine (Cedar), one normative mapping between them.
3. **Tool self-description convention.** `--agentrc-schema` vs sibling
   `<tool>.toolspec.json` (or both).
4. **Override auditability.** Whether platform narrowing / rejecting a request,
   or substituting an embedded resource via `*.origin`, MUST emit an audit
   record.
5. **`CAPABILITY` label shape.** One label per capability
   (`org.agentrc.capability.streaming=true`) vs a single comma-joined label —
   pick one.
6. **`TOOL` sugar.** `TOOL` was removed in favour of `COPY` / `ADD`. A thin
   `TOOL name@source` sugar that desugars to `ADD --remote …/mnt/tools/name` MAY
   be reconsidered later, but is **not** in this spec.

### 14.3 Deferred to a later version

- **Secrets / credential resolution:** Vault, secret brokers, env injection,
  OIDC / workload identity. No keyword and no label schema in this draft; the
  Agentfile leaves credential binding entirely to the platform until this is
  designed ([§12.1](#121-secrets-deferred)).
- **A2A (agent-to-agent protocol):** Agent Cards, agent discovery, cross-agent
  delegation, and the governance algebra across an agent-to-agent call.
  Capability *exposure* (via `IDENTITY` / `CAPABILITY` / labels) is in this
  version; the *protocol* for one agent to find and call another is not.
  External multi-agent orchestration that references packaged agents by digest —
  distinct from the A2A protocol — is parked for a future draft.
