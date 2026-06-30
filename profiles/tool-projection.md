---
layout: doc
title: Tool Projection
description: "The /mnt Projection profile (v1): how a platform projects an agent's tools, skills, MCP servers, SOP, and runtime state under /mnt."
permalink: /profiles/tool-projection/
---
# `/mnt` Projection Profile

**Version:** v1 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.io/agentfile:v1`)  
**Date:** 2026-06-30  
**Audience:** platform / runner authors, substrate implementers

## Purpose

This profile defines the **projected filesystem** an agentrc agent sees at run
time: the `/mnt` tree where its tools, skills, MCP servers, system prompt, and
live runtime state are exposed. It is the run-time companion to the
[Specification](/spec/): the build embeds resources at `/mnt/<...>` and emits
`org.agentrc.*` labels; this profile describes what a platform MUST project so a
`CMD` can find and invoke those resources.

Projection does **not** make agentrc a runtime. It defines the portable surface
a conformant platform exposes on the chosen substrate (`agentrc run --substrate`).
Every projected invocation remains subject to platform policy, enforced by
**Cedar** (see the [Enforcement profile](/profiles/security/)).

## The `/mnt` projection layout

A conformant platform MUST project the agent's filesystem under **`/mnt`** with
this layout:

```text
/mnt
├── tools/     # each tool is an executable; argv in, structured output out
├── skills/    # skill bundles (a directory of instructions/scripts/resources)
├── mcp/       # MCP server bundles or configs, projected as tool directories
├── proc/      # runtime-populated: live policy, identity, budgets, audit tail
└── SOP        # the agent's system prompt, as a readable file
```

The mount point is always `/mnt`. The **destination path** under `/mnt`
determines the resource type — `COPY`/`ADD --remote` into `/mnt/tools/<name>`
makes a tool, into `/mnt/skills/<name>` a skill, into `/mnt/mcp/<name>` an MCP
bundle, and into `/mnt/SOP` the system prompt. A local implementation MAY back
the same model with a real filesystem, a FUSE mount, or a generated directory,
provided the paths and the invocation contract below hold.

## Tools

A **tool** is an executable file under `/mnt/tools/`. The platform projects each
tool the build embedded (`org.agentrc.tool.<name>=local`) or referenced
(`=<digest>` / `=runtime:<url>`) so the `CMD` can invoke it by path.

To be self-describing, a tool SHOULD expose its schema by **one** of:

- `--agentrc-schema` — invoking the executable with this flag prints a JSON
  schema describing its arguments and output to **stdout**, or
- a sibling `<tool>.toolspec.json` next to the executable.

Which convention is canonical is
[open decision #3](/spec/#142-open-decisions-surface-these-do-not-silently-resolve);
a platform SHOULD accept either. The projection layer reads the schema to render
help/`man`-style output and to enumerate the agent's capability surface.

```text
/mnt/tools/
  file_read
  file_read.toolspec.json
  http_get
```

## Skills

Skill bundles are projected under `/mnt/skills/<name>/` as directories of
instructions, scripts, and resources — a `SKILL.md` plus its supporting files.
They are added with `COPY ./skills/<name> /mnt/skills/<name>` or
`ADD --remote … /mnt/skills/<name>` and recorded in `org.agentrc.skill.*` labels.

```text
/mnt/skills/
  code-review/
    SKILL.md
```

## MCP servers

MCP server bundles and configs are projected under `/mnt/mcp/<name>/`. agentrc
declares and governs MCP; it does not replace it. An embedded MCP server carries
both a digest and an `.origin` label, so a platform MAY re-point it to an
internal mirror at deploy time without rebuilding (see the
[OCI labels & package profile](/profiles/oci-package/)). A `--runtime` MCP
server (`org.agentrc.mcp.<name>=runtime:<url>`) is fetched when the agent
bootstraps.

```text
/mnt/mcp/
  github/
    ...
```

## SOP

The system prompt is always a readable file at **`/mnt/SOP`**, regardless of
whether it was authored inline, as a heredoc, or file-backed via `COPY`/`ADD`.
The platform loads `/mnt/SOP` at startup. Labels carry only a **pointer plus
digest** (`org.agentrc.sop=/mnt/SOP`, `org.agentrc.sop.sha256=<digest>`), never
the full text, so the platform can verify or override the file without carrying
the prompt in metadata.

## Invocation contract

A projected tool SHOULD be directly executable, and the invocation contract
SHOULD be:

1. **argv** — arguments passed as command-line arguments;
2. **stdin** — a structured payload (e.g. JSON) on standard input;
3. **stdout** — structured output (e.g. JSON) on standard output;
4. **stderr** — human-readable diagnostics;
5. **exit `0`** — success;
6. **non-zero exit** — failure;
7. a **distinct denied/policy exit code** when the invocation is refused by
   policy, so callers can tell a policy denial apart from an ordinary failure.

```bash
echo '{"path":"./README.md"}' | /mnt/tools/file_read
```

## `/mnt/proc` — runtime state

`/mnt/proc` is **runtime-populated** by the platform (the agent does not author
it). It surfaces the live, granted state so the agent and operators can read
what is actually in force:

```text
/mnt/proc/
  identity     # the resolved org.agentrc.identity.* for this agent
  policy       # the requests as granted/narrowed by the platform (post-Cedar)
  budgets      # live timeouts, memory/CPU, sub-agent and other limits
  audit        # a tail of the audit stream (tool calls, grants, denials)
```

What appears under `/mnt/proc/policy` is the **granted** set, not the raw
request: the platform reads the `org.agentrc.*` labels, evaluates each request
against its own Cedar policies, and projects the result. A request the platform
narrowed or rejected is reflected here, not silently honoured.

## Policy applies to every projected invocation

Projection is a convenience surface, not an authorization. Every invocation of a
projected tool, skill, or MCP server remains subject to platform policy and is
evaluated by **Cedar** (deny-by-default; `forbid` overrides `permit`; monotonic
across `FROM`). The presence of `/mnt/tools/<name>` does **not** by itself grant
`Action::"tool.invoke"` on `Tool::"<name>"`; the platform's grant decision does.
The full request → Cedar mapping is normative in the
[Enforcement (Cedar) profile](/profiles/security/).

## CLI equivalence

A platform MAY offer equivalent commands instead of (or alongside) filesystem
projection. These are conveniences over the same projected surface:

| Projection | CLI equivalent |
|---|---|
| `/mnt/tools/file_read --path x` | `arc run <ref> -- tool file_read --path x` |
| `cat /mnt/proc/policy` | `arc run <ref> -- policy show` |
| `cat /mnt/proc/audit` | `arc run <ref> -- audit tail` |

Substrate and isolation are **run-time** choices (`arc run --isolation` /
`--substrate`), never Agentfile directives; this projection layout is identical
regardless of which substrate runs `CMD`.
