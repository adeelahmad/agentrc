---
layout: doc
title: Core
description: "Agentfile Core Profile (0.1.0-draft.5): parse the Dockerfile-shaped Agentfile and compile it to org.agentrc.* labels and layers."
permalink: /profiles/core/
---
# Agentfile Core Profile

**Version:** 0.1.0-draft.5 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.agentfile/v0.1`)  
**Date:** 2026-06-30  
**Audience:** compiler / frontend authors (the `agentrc` BuildKit frontend and the `agentrc` / `arc` CLI)

> The normative keywords **MUST**, **MUST NOT**, **SHOULD**, and **MAY** follow
> [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119). This profile is the
> compile-side companion to the [specification](/spec/); where the two disagree,
> the spec wins.

## Purpose

The Core Profile defines what it means to **read a Dockerfile-shaped Agentfile and
compile it into an OCI artifact** — namely, the `org.agentrc.*` image labels plus
the layers that carry the agent's `/mnt` tree. It is the contract a compiler
implementation claims: recognize the four agentrc keywords (`IDENTITY`,
`CAPABILITY`, `SOP`, `POLICY`) and the standard Dockerfile keywords, and emit the
labels and layers exactly as the spec's translation tables require.

This profile covers **build**, not run time. What a platform does with the labels
(grant / narrow / reject, enforce via Cedar, project `/mnt`) is the
[Platform Conformance Profile](/profiles/runner-conformance/) and the
[Enforcement (Cedar) Profile](/profiles/security/); how the labels and package are
shaped is the [OCI Labels &amp; Package Profile](/profiles/oci-package/).

A small core is deliberate. The Agentfile adds only **four** keywords to a model
developers already know; everything else is a standard Dockerfile keyword. Keeping
the new surface tiny is what lets two independent compilers — the BuildKit frontend
and the native CLI — produce **byte-for-byte identical** artifacts.

## Keyword summary

There are **four new keywords**. Everything else is a standard Dockerfile keyword,
used as-is or with a documented extension. A conformant compiler MUST recognize
all of these.

| Keyword | Origin | Role in agentrc |
|---|---|---|
| **`IDENTITY`** | **New** | Who / what the agent is: name, version, author, description. |
| **`CAPABILITY`** | **New** | What modalities / patterns the agent supports (text, streaming, multimodal, …). |
| **`SOP`** | **New** | The agent's system prompt / objective / standard operating procedure. |
| **`POLICY`** | **New** | Typed, namespaced **request** for a resource, model, or operational constraint. |
| `FROM` | Dockerfile | Base image **and** agent inheritance (`FROM another-agent`). |
| `CMD` | Dockerfile | The agent's invocation surface — the framework / loop to run. |
| `COPY` | Dockerfile | Add **local** tools, skills, MCP bundles, or an SOP file into the `/mnt` tree. |
| `ADD` | Dockerfile (extended) | Add **remote** resources via `--remote` plus delivery flags. |
| `HEALTHCHECK` | Dockerfile | Liveness probe; MAY invoke a projected tool. |
| `LABEL` | Dockerfile | Standard OCI metadata; available for hand-authored `org.agentrc.*` metadata. |
| `ENV` / `ARG` / `WORKDIR` / `USER` / `EXPOSE` / `RUN` | Dockerfile | Standard semantics; available, unchanged. |

> **There is no `TOOL`, `MCP`, `SERVER`, `FUNC`, `CRED`, `SECRET`, `AUDIT`,
> `MOUNT`, `MEMORY`, or `RATELIMIT` keyword.** Tools / skills / MCP are added with
> `COPY` / `ADD`; memory / context / model / CPU are `POLICY` requests; audit
> rides on `agent.hooks.*`; secrets are **deferred**. Any earlier draft that
> recognized those keywords is **stale** and MUST NOT be revived by this profile.

## The `/mnt` projection layout

The destination path under `/mnt` is what tells the compiler *what a resource is*.
There is no per-type keyword; a file copied to `/mnt/tools/x` is a tool, one
copied to `/mnt/mcp/x` is an MCP bundle, and so on. The compiler MUST use this
layout:

```text
/mnt
├── tools/     # each tool is an executable; argv in, structured output out
├── skills/    # skill bundles (a directory of instructions/scripts/resources)
├── mcp/       # MCP server bundles or configs
├── proc/      # runtime-populated (live policy, identity, budgets, audit tail)
└── SOP        # the agent's system prompt, as a readable file
```

`/mnt/proc` is populated at run time by the platform, not by the compiler. See the
[`/mnt` Projection Profile](/profiles/tool-projection/) for the full layout.

## Required behavior

An implementation claiming this profile MUST:

1. read a text file named `Agentfile` or an explicitly supplied path
   (`-f Agentfile`);
2. ignore blank lines and comments (`#`), while honouring the
   `# syntax=agentrc.agentfile/v0.1` parser directive on the first line;
3. parse the four agentrc keywords and the standard Dockerfile keywords listed
   above, case-sensitively;
4. capture an `SOP <<EOF … EOF` **heredoc verbatim**, without interpreting its
   inner lines as Agentfile instructions, and embed it as the readable file
   `/mnt/SOP`;
5. accept the inline and file-backed `SOP` forms and embed each at `/mnt/SOP`,
   emitting a **pointer + digest** label — never the full prompt text;
6. honour the `ADD --remote` delivery flags — `--cached` (default) / `--runtime`
   and `--fail-if-unavailable` (default) / `--warn-if-unavailable`, plus standard
   `--chmod` / `--chown` — fetching `--cached` resources at build and recording
   `--runtime` resources as references;
7. use the **destination path under `/mnt`** to classify each `COPY` / `ADD`
   resource as a tool, skill, MCP bundle, or SOP;
8. translate authored intent into `org.agentrc.*` labels and `/mnt` layers exactly
   per the [spec §9 tables](/spec/) (reproduced below);
9. auto-derive an **explicit, attributed** `network` egress label from any URL in
   `agent.hooks.*` or `agent.interrupt_endpoint`, never a silent egress;
10. report line numbers for parse failures, and reject unknown agentrc keywords
    rather than silently dropping them.

A `POLICY` line is authored in **short form** (no `org.agentrc.` prefix); the
compiler prepends the namespace when it emits the label. The compiler never
evaluates a `POLICY` — it only records the request. Whether a request is granted,
narrowed, or rejected is a platform decision.

## Build translation (compile targets)

These are the tables the compiler MUST implement. They are reproduced from the
[specification](/spec/) so the compile-side contract is self-contained.

### `POLICY` → labels

Each `POLICY <key> <value>` becomes `org.agentrc.<key>=<value>`.

| Authored | Emitted label |
|---|---|
| `POLICY agent.idle_timeout 5m` | `org.agentrc.agent.idle_timeout=5m` |
| `POLICY agent.context.type autocompressed` | `org.agentrc.agent.context.type=autocompressed` |
| `POLICY agent.hooks.pre https://hooks.internal/pre-step` | `org.agentrc.agent.hooks.pre=https://hooks.internal/pre-step` (+ auto-derived egress) |
| `POLICY agent.sub_agents.max 5` | `org.agentrc.agent.sub_agents.max=5` |
| `POLICY substrate.runtime.memory 8gb` | `org.agentrc.substrate.runtime.memory=8gb` |
| `POLICY substrate.ptty true` | `org.agentrc.substrate.ptty=true` |
| `POLICY model.name claude-opus-4` | `org.agentrc.model.name=claude-opus-4` |
| `POLICY model.capability vision` | `org.agentrc.model.capability.vision=true` |
| `POLICY network dns:api.github.com:443` | `org.agentrc.network.dns.api.github.com=443` |

For an auto-derived egress, the compiler MUST emit both the egress label and its
attribution, so it is never a silent network hole:

```text
org.agentrc.network.dns.hooks.internal=443
org.agentrc.network.dns.hooks.internal.source=auto:agent.hooks.pre
```

### `IDENTITY` / `CAPABILITY` / `SOP` → labels

| Authored | Emitted label(s) |
|---|---|
| `IDENTITY name=claims-triage version=1.0` | `org.agentrc.identity.name=claims-triage`, `org.agentrc.identity.version=1.0` |
| `CAPABILITY streaming` | `org.agentrc.capability.streaming=true` |
| `SOP …` / `COPY ./sop.md /mnt/SOP` | `org.agentrc.sop=/mnt/SOP`, `org.agentrc.sop.sha256=<digest>` (pointer + digest, never full text) |

### Resource delivery → layers + labels

| Authored | Build behaviour | Emitted label(s) |
|---|---|---|
| `COPY ./tools/x /mnt/tools/x` | Embed file as a layer. | `org.agentrc.tool.x=local` |
| `ADD --remote --cached <url> /mnt/tools/x` | Fetch at build, embed as a layer. | `org.agentrc.tool.x=<digest>` + `org.agentrc.tool.x.origin=<url>` |
| `ADD --remote --runtime <url> /mnt/mcp/x` | Do not embed; record reference. | `org.agentrc.mcp.x=runtime:<url>` |

For every embedded MCP server or skill, the compiler MUST emit **both** the
resolved digest and the origin reference, so a platform can re-point to a mirror
at deploy time without rebuilding:

```text
org.agentrc.mcp.github=sha256:abc123...
org.agentrc.mcp.github.origin=https://registry.agentrc.io/mcp/github:latest
```

### Policy encoding

The compiled request set MUST be retrievable in **either** form, and the build MAY
support both, selected by `--policy-mode`:

- **inline** — values live directly in labels (best for small request sets);
- **digest** — labels carry a digest of a structured manifest embedded as a layer
  (best for large request sets).

Which form is the **default** for `arc build` is
[open decision #1](/spec/) — the `--policy-mode inline|digest` flag is the seam.

## Identical output from both front doors

agentrc has two build paths — the BuildKit frontend invoked by the
`# syntax=agentrc.agentfile/v0.1` directive, and the native `agentrc` / `arc` CLI
(see the [CLI page](/cli/)). They are two front doors to the same compiler.

> **A conformant implementation MUST produce identical OCI artifacts — same
> layers, same `org.agentrc.*` labels, same digest — from the same Agentfile,
> whether built via the frontend or via `arc build`.** This is the single most
> important conformance property of this profile and is exercised adversarially by
> the [conformance suite](/docs/conformance/) (`build-labels-identical`).

## Worked translation

A minimal Agentfile:

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim
IDENTITY name=claims-triage version=1.0 author=acme
CAPABILITY text
CAPABILITY streaming
SOP You are a claims-triage specialist. Escalate anything ambiguous to a human.
CMD claude --print
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
POLICY model.name        claude-opus-4
POLICY network           dns:api.github.com:443
```

compiles to layers carrying `/mnt/tools/file_read` and `/mnt/SOP`, plus these
labels in the image config:

```text
org.agentrc.identity.name=claims-triage
org.agentrc.identity.version=1.0
org.agentrc.identity.author=acme
org.agentrc.capability.text=true
org.agentrc.capability.streaming=true
org.agentrc.sop=/mnt/SOP
org.agentrc.sop.sha256=<digest>
org.agentrc.tool.file_read=local
org.agentrc.model.name=claude-opus-4
org.agentrc.network.dns.api.github.com=443
```

The platform reads those labels — never the Agentfile source.

## Validation recommendations

A compiler / linter SHOULD warn when:

1. `IDENTITY name=` is missing in an artifact intended for publication;
2. `FROM` uses a mutable tag such as `latest`;
3. an `ADD --remote` uses `--runtime` without a clear failure mode, so the
   default `--fail-if-unavailable` is recorded explicitly;
4. a `COPY` / `ADD` targets a path outside `/mnt` for a resource that is meant to
   be projected (it will not be classified as a tool / skill / mcp);
5. a removed keyword (`TOOL`, `MCP`, `CRED`, `MOUNT`, `RATELIMIT`, …) is
   encountered — it is stale and MUST be rejected, with a pointer to the
   `COPY` / `ADD` / `LABEL` / `POLICY` replacement.

## Compatibility note

The reference compiler — the BuildKit frontend and the `arc` CLI — is in progress.
This profile defines the **target** compile semantics; an implementation SHOULD
publish a support matrix rather than implying it emits every label in the tables
above. Do not advertise a profile you do not pass.
