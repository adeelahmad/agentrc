---
layout: doc
title: CLI
description: "Building agentrc agents: the BuildKit frontend (no CLI to install) and the native agentrc / arc CLI."
permalink: /cli/
---
# agentrc CLI

An Agentfile compiles to a standard OCI artifact carrying `org.agentrc.*`
labels. There are **two front doors** to that compiler, and they produce
**identical artifacts**: the BuildKit frontend (which needs no new tool) and the
native `agentrc` CLI (alias `arc`). This page is the practical companion to
[§10 of the specification](/spec/).

<div class="callout">
<strong>Reference implementation available.</strong> A Go implementation of
both build paths lives in <a href="https://github.com/adeelahmad/agentrc/tree/master/tooling">this
repository's <code>tooling/</code> directory</a> — see the table below for
per-command status. The frontend image is published at
<code>ghcr.io/adeelahmad/agentrc-frontend</code>, so a
<code># syntax=ghcr.io/adeelahmad/agentrc-frontend</code> directive makes
<code>docker build -f Agentfile .</code> auto-route through it; you can also pin
it with <code>--build-arg BUILDKIT_SYNTAX=ghcr.io/adeelahmad/agentrc-frontend:latest</code>.
<code>arc run</code> ships as a reference translator (it renders backend
deploy artifacts as a dry run; it does not apply them to a live cloud), while
<code>sign</code> and <code>verify</code> remain planned. The <a href="/spec/">specification</a>
remains the source of truth, not this implementation — see
<a href="/docs/conformance/">Conformance</a> for exactly what it covers.
</div>

## BuildKit frontend

The agentrc frontend image is published at
`ghcr.io/adeelahmad/agentrc-frontend`. A `# syntax=ghcr.io/adeelahmad/agentrc-frontend`
directive on the **first line** of an `Agentfile` makes Docker / BuildKit pull
the frontend automatically, so a plain `docker build` routes through it. You can
also pin it explicitly with `--build-arg BUILDKIT_SYNTAX=<image>`:

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim

IDENTITY name=hello version=0.1 author=acme
IDENTITY description="Minimal agentrc agent"
CAPABILITY text
SOP You are a minimal example agent. Read a file when asked; do nothing else.
CMD python ./agent.py

# Tool (local, embedded) — projected under /mnt/tools/
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# Model + operational requests (platform grants, narrows, or rejects)
POLICY model.name         claude-sonnet-4
POLICY agent.tool_timeout 30s

# Network egress request
POLICY network dns:api.example.com:443

HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema
```

```bash
# Routed via the `# syntax=ghcr.io/adeelahmad/agentrc-frontend:latest` line:
docker build -f Agentfile -t ghcr.io/you/hello:1.0 .

# Or pin the frontend explicitly:
docker build -f Agentfile --build-arg BUILDKIT_SYNTAX=ghcr.io/adeelahmad/agentrc-frontend:latest -t ghcr.io/you/hello:1.0 .
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
agentrc run    <ref> --backend local|bedrock|kubernetes [per-backend flags]
```

The four **core** commands are `build`, `push`, `pull`, and `run` (spec §10);
the rest are tooling around them.

Reference translators — a proof of concept until platforms read `org.agentrc.*` labels natively. Not production runners.

| Command | Purpose | Status |
|---|---|---|
| `agentrc init` (`arc init`) | Scaffold a starter Agentfile. | `implemented` |
| `agentrc lint` (`arc lint`) | Check an Agentfile for keyword and request errors before building. | `implemented` |
| `agentrc lock` (`arc lock`) | Pin `ADD --remote` resources to digests for reproducible builds. | `implemented` |
| `agentrc build` (`arc build`) | **Core (§10).** Compile an Agentfile to an OCI artifact, emitting `org.agentrc.*` labels and embedding `--cached` resources as layers. `--policy-mode inline\|digest` selects how the request set is encoded (see below). | `implemented` |
| `agentrc inspect` (`arc inspect`) | Read an artifact's `org.agentrc.*` labels to review what an agent requests before it runs. | `implemented` |
| `agentrc sign` (`arc sign`) | Sign an artifact (Sigstore). | `planned` |
| `agentrc verify` (`arc verify`) | Verify an artifact's signature and provenance. | `planned` |
| `agentrc push` (`arc push`) | **Core (§10).** Push the artifact to any OCI registry. | `implemented` |
| `agentrc pull` (`arc pull`) | **Core (§10).** Pull an artifact from any OCI registry. | `implemented` |
| `agentrc run` (`arc run`) | **Core (§10).** Translate an artifact into a chosen backend's deploy form (`--backend local\|bedrock\|kubernetes`). `--backend` (and `--isolation`, scoped to `--backend local`) are **run-time** choices, never Agentfile directives. | `implemented (local, bedrock, kubernetes — reference translators)` |

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

`--backend` (`local`, `bedrock`, `kubernetes`) and `--isolation` (`local`,
`container`, `microvm`, scoped to `--backend local`) are selected when you run,
not when you build. The Agentfile never names a substrate. A `POLICY substrate.*`
line *requests* resources (memory, CPU, a device, a pseudo-terminal); the
platform grants, narrows, or rejects those requests and realizes them on
whatever backend you point `arc run` at.

```bash
arc build -t ghcr.io/you/hello:1.0 .
arc push  ghcr.io/you/hello:1.0
arc run   ghcr.io/you/hello:1.0 --backend local --isolation microvm
```

**Dropped backends.** GCP is dropped as a first-class `--backend`: its managed
Agent Runtime is Python-only, and GKE is reached through the `kubernetes` backend
anyway. Docker Compose is dropped: it offers no `network.*` egress enforcement
without a bespoke sidecar, so it cannot honor the granted request set.

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

Inspect them with standard OCI tooling — `docker inspect`, `arc inspect`, or
any registry client — to review exactly what an agent requests
(identity, capabilities, requested model, requested egress, tools / skills / MCP
with digests, sub-agent limits) **before** it runs.

## Identical artifacts, two paths

The output of `docker build` through the frontend and the output of `arc build`
MUST be **identical** OCI artifacts: same layers, same `org.agentrc.*` labels,
same digest given the same inputs. The frontend and the CLI are two front doors
to one compiler. At run time, the platform reads those labels, grants / narrows /
rejects each request, and enforces the grant with [Cedar, platform-side](/profiles/security/)
(deny-by-default, `forbid` over `permit`, tightening-only across `FROM`).

## Try it now

```bash
git clone https://github.com/adeelahmad/agentrc.git
cd agentrc
go build -o bin/agentrc ./cmd/agentrc
./bin/agentrc lint examples/Agentfile.code-reviewer
```

See the [reference implementation](/tooling/) (`tooling/README.md` in the
repository) for the full build/test/frontend walkthrough.

## In the meantime

- Read [What is agentrc?](/docs/what-is-agentrc/) for the why.
- Walk the [Quickstart](/docs/quickstart/) end to end.
- Read the full build and label translation in the [Specification](/spec/).
- Track or contribute on [GitHub](https://github.com/adeelahmad/agentrc), or say hi on [Discord](https://discord.gg/jWx6Qak5D).
