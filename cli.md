---
layout: doc
title: CLI
description: "Building agentrc agents: the BuildKit frontend (nothing to install) and the native agentrc / arc CLI."
permalink: /cli/
---
# agentrc CLI

An Agentfile compiles to a standard OCI artifact carrying `org.agentrc.*`
labels. There are **two front doors** to that compiler, and they produce
**identical artifacts**: the BuildKit frontend (which needs no new tool) and the
native `agentrc` CLI (alias `arc`). This page is the practical companion to
[§10 of the specification](/spec/).

<div class="callout">
<strong>In progress.</strong> The native <code>agentrc</code> / <code>arc</code>
CLI is being built. The BuildKit frontend path below works with any current
Docker / BuildKit install. The <a href="/spec/">specification</a>, schemas, and
examples on this site are the source of truth today; both build paths implement
them and MUST emit the same labels and layers.
</div>

## BuildKit frontend (nothing to install)

If you already have Docker / BuildKit, you need install nothing. Add the
`# syntax=` directive as the **first line** of your `Agentfile` and build it like
any image:

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim
IDENTITY name=hello version=1.0 author=you
CAPABILITY text
SOP You are a concise local assistant. Answer in one paragraph.
CMD python ./agent.py
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
POLICY model.name claude-opus-4
POLICY network dns:api.github.com:443
```

```bash
docker build -f Agentfile -t ghcr.io/you/hello:1.0 .
```

The `# syntax=` line routes the Agentfile through the agentrc frontend image,
which parses the four agentrc keywords (`IDENTITY`, `CAPABILITY`, `SOP`,
`POLICY`) and the `ADD --remote` extension, compiles them to LLB, embeds
`--cached` resources as layers, and writes the `org.agentrc.*` labels into the
image config. A stock `docker build` *without* the `# syntax=` directive
understands only the standard Dockerfile keywords, so the agentrc keywords would
be unrecognized — the directive is what enables them.

## Native `agentrc` CLI (`arc`)

The primary command is `agentrc`; `arc` is the short alias. It is one tool
covering build, registry transport, and run.

```bash
agentrc build  [-f Agentfile] [-t <ref>] [--policy-mode inline|digest] .
agentrc push   <ref>
agentrc pull   <ref>
agentrc run    <ref> [--isolation local|container|microvm] [--substrate <driver>]
```

The four **core** commands are `build`, `push`, `pull`, and `run` (spec §10);
the rest are tooling around them. All commands are currently `planned`.

| Command | Purpose | Status |
|---|---|---|
| `agentrc init` (`arc init`) | Scaffold a starter Agentfile. | `planned` |
| `agentrc lint` (`arc lint`) | Check an Agentfile for keyword and request errors before building. | `planned` |
| `agentrc lock` (`arc lock`) | Pin `ADD --remote` resources to digests for reproducible builds. | `planned` |
| `agentrc build` (`arc build`) | **Core (§10).** Compile an Agentfile to an OCI artifact, emitting `org.agentrc.*` labels and embedding `--cached` resources as layers. `--policy-mode inline\|digest` selects how the request set is encoded (see below). | `planned` |
| `agentrc inspect` (`arc inspect`) | Read an artifact's `org.agentrc.*` labels to review what an agent requests before it runs. | `planned` |
| `agentrc sign` (`arc sign`) | Sign an artifact (Sigstore). | `planned` |
| `agentrc verify` (`arc verify`) | Verify an artifact's signature and provenance. | `planned` |
| `agentrc push` (`arc push`) | **Core (§10).** Push the artifact to any OCI registry. | `planned` |
| `agentrc pull` (`arc pull`) | **Core (§10).** Pull an artifact from any OCI registry. | `planned` |
| `agentrc run` (`arc run`) | **Core (§10).** Run an artifact on a chosen substrate. `--isolation` / `--substrate` are **run-time** choices, never Agentfile directives. | `planned` |

### `--policy-mode inline | digest`

The compiled request set MUST be retrievable by the platform in either form, and
`build` exposes the choice:

- `inline` — the `POLICY`-derived values live directly in the `org.agentrc.*`
  labels. Best for small request sets.
- `digest` — the labels carry a digest of a structured manifest embedded as a
  layer; the full request set is fetched from there. Best for large request sets.

Which mode is the **default** is [open decision #1](/spec/); `--policy-mode` is
the seam, so you can pin it explicitly today.

### Substrate is a run-time choice

`--isolation` (`local`, `container`, `microvm`) and `--substrate <driver>` are
selected when you run, not when you build. The Agentfile never names a substrate.
A `POLICY substrate.*` line *requests* resources (memory, CPU, a device, a
pseudo-terminal); the platform grants, narrows, or rejects those requests and
realizes them on whatever substrate you point `arc run` at.

```bash
arc build -t ghcr.io/you/hello:1.0 .
arc push  ghcr.io/you/hello:1.0
arc run   ghcr.io/you/hello:1.0 --isolation microvm
```

## Reading the labels

The platform reads the artifact's `org.agentrc.*` labels — not the Agentfile.
The Agentfile above emits labels such as:

```text
org.agentrc.identity.name=hello
org.agentrc.identity.version=1.0
org.agentrc.capability.text=true
org.agentrc.sop=/mnt/SOP
org.agentrc.sop.sha256=<digest>
org.agentrc.tool.file_read=local
org.agentrc.model.name=claude-opus-4
org.agentrc.network.dns.api.github.com=443
```

Inspect them with standard OCI tooling — `docker inspect`, `arc inspect` (in
progress), or any registry client — to review exactly what an agent requests
(identity, capabilities, requested model, requested egress, tools / skills / MCP
with digests, sub-agent limits) **before** it runs.

## Identical artifacts, two paths

The output of `docker build` through the frontend and the output of `arc build`
MUST be **identical** OCI artifacts: same layers, same `org.agentrc.*` labels,
same digest given the same inputs. The frontend and the CLI are two front doors
to one compiler. At run time, the platform reads those labels, grants / narrows /
rejects each request, and enforces the grant with [Cedar, platform-side](/profiles/security/)
(deny-by-default, `forbid` over `permit`, tightening-only across `FROM`).

## In the meantime

- Read [What is agentrc?](/docs/what-is-agentrc/) for the why.
- Walk the [Quickstart](/docs/quickstart/) end to end.
- Read the full build and label translation in the [Specification](/spec/).
- Track or contribute on [GitHub](https://github.com/adeelahmad/agentrc), or say hi on [Discord](https://discord.gg/jWx6Qak5D).
