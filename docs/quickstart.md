---
layout: doc
title: Quickstart
description: "Author, build, inspect, push, and run a minimal agentrc agent in five steps."
permalink: /docs/quickstart/
---
# Quickstart

This walks you from a blank file to a running agent in five steps: author an
**Agentfile**, build it (two ways), read the `ai.agentrc.*` **labels** the build
emits, push the artifact, and run it. The Agentfile is Dockerfile-shaped — if you
know `docker build`, you already know most of this. For the full normative
grammar, see the [specification](/spec/).

<div class="callout">
<strong>Status: steps 1–4 work today; step 5 is planned.</strong> A reference
implementation of the native <code>agentrc</code> CLI (<code>build</code>,
<code>inspect</code>, <code>push</code>, <code>pull</code>, <code>lint</code>,
<code>lock</code>, <code>init</code>) and the BuildKit frontend both exist —
see <a href="/cli/">the live CLI status table</a> and <code>tooling/</code> in
the repository. The frontend image is published at
<code>ghcr.io/adeelahmad/agentrc-frontend</code>, so the <code># syntax=</code>
line routes <code>docker build -f Agentfile .</code> through it (see step 2).
<strong>Step 5 (<code>arc run</code>) is planned</strong> — agentrc declares
agents, it does not ship a runtime; running an artifact is a compatible
platform's job.
</div>

## Installation

You need the native `agentrc` CLI (alias `arc`) for `build`, `inspect`, `push`,
`pull`, and the reference `run` translators. Pick whichever install method suits
you — they all land the same binary. Step 2 also offers a CLI-free path via the
BuildKit frontend, so if you only ever `docker build` you can skip this.

**Quick install (curl).** Installs both the `agentrc` binary and its `arc` alias
to `/usr/local/bin` (or `~/.local/bin` when that is not writable), verified
against the release `checksums.txt`. Covers macOS and Linux on amd64 and arm64:

```bash
curl -fsSL https://agentrc.ai/install.sh | sh
```

As with any piped installer, read the script first if you'd rather vet it:

```bash
curl -fsSL https://agentrc.ai/install.sh | less
```

**Homebrew.** On macOS or Linux with Homebrew:

```bash
brew install adeelahmad/tap/agentrc
```

**Go.** If you have a Go toolchain:

```bash
go install github.com/adeelahmad/agentrc/cmd/agentrc@latest
```

**From source.** Clone and build the binary yourself:

```bash
git clone https://github.com/adeelahmad/agentrc
cd agentrc
go build -o arc ./cmd/agentrc
```

Verify the install:

```bash
arc version
# agentrc <ver> (spec 0.1.0-draft.6, <os>/<arch>)
```

## 1. Create an `Agentfile`

Four new keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) sit on top of
standard Dockerfile keywords. Save this as `Agentfile`:

```dockerfile
# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim

# --- Who the agent is -------------------------------------------------------
IDENTITY name=hello version=0.1 author=acme
IDENTITY description="Minimal agentrc agent"

# --- What it supports -------------------------------------------------------
CAPABILITY text

# --- System prompt / objective ----------------------------------------------
SOP You are a minimal example agent. Read a file when asked; do nothing else.

# --- Invocation surface (any framework) -------------------------------------
CMD python ./agent.py

# --- A local tool, embedded into the /mnt tree ------------------------------
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# --- Typed requests to the platform -----------------------------------------
POLICY model.name         claude-sonnet-4
POLICY agent.tool_timeout 30s

# --- Network egress request -------------------------------------------------
POLICY network dns:api.example.com:443

# --- Liveness probe ---------------------------------------------------------
HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema
```

Each `POLICY` line is a **request**, not a guarantee — you are asking the platform
for a model and an egress host. The platform decides what to honour (step 5). The
tool is just an executable placed under `/mnt/tools/`; its destination path under
`/mnt` is what makes it a tool.

## 2. Build the artifact

There are two front doors to the same compiler — they produce **identical** OCI
artifacts.

**BuildKit frontend.** The agentrc frontend image is published at
`ghcr.io/adeelahmad/agentrc-frontend`. Make it the first line of your Agentfile:

```dockerfile
# syntax=ghcr.io/adeelahmad/agentrc-frontend:latest
```

then a plain build routes the file through it:

```bash
docker build -f Agentfile -t ghcr.io/you/hello:1.0 .
```

To pin the frontend explicitly instead of via the `# syntax=` line, pass it as a
build arg:

```bash
docker build -f Agentfile --build-arg BUILDKIT_SYNTAX=ghcr.io/adeelahmad/agentrc-frontend:latest -t ghcr.io/you/hello:1.0 .
```

**Native `agentrc` CLI** (alias `arc`):

```bash
arc build -t ghcr.io/you/hello:1.0 .
```

Either path parses the four agentrc keywords plus the `ADD --remote` extension,
embeds any `--cached` resources as layers, and writes the `ai.agentrc.*` labels.
The output is the same content-addressed OCI artifact regardless of which you use.

## 3. Read the labels

The build translates your authored intent into namespaced OCI image labels under
`ai.agentrc.*`. **The platform reads these labels — it never parses the
Agentfile.** Inspect them with `docker inspect` or `arc inspect`:

```bash
arc inspect ghcr.io/you/hello:1.0
```

The Agentfile above emits labels like:

```text
ai.agentrc.identity.name=hello
ai.agentrc.capability.text=true
ai.agentrc.model.name=claude-sonnet-4
ai.agentrc.network.dns.api.example.com=443
ai.agentrc.tool.file_read=local
ai.agentrc.sop=/mnt/SOP
ai.agentrc.sop.sha256=<digest>
```

The SOP is embedded as a readable file at `/mnt/SOP`; its label is a
**pointer plus digest**, never the full prompt text. This labels-only manifest is
what a security reviewer reads before the agent ever runs.

## 4. Push to a registry

The artifact is an ordinary OCI image, so it pushes to any OCI-compatible
registry:

```bash
arc push ghcr.io/you/hello:1.0
```

## 5. Run on a substrate

Status: `arc run` is planned — agentrc declares agents, it does not ship a
runtime (see [Non-goals](/docs/non-goals/) and the [CLI status table](/cli/)).
This is the interface a compatible platform provides:

```bash
arc run ghcr.io/you/hello:1.0 --isolation microvm
```

At run time the platform pulls the artifact, reads every `ai.agentrc.*` label,
and evaluates each request — model, egress, tools — against its own organization
policy and available resources, then **grants, narrows, or rejects** it.
Enforcement is **Cedar, platform-side**: the platform compiles the granted
requests plus its own org rules into a Cedar `PolicySet` and evaluates with
deny-by-default and `forbid` over `permit`. You never write Cedar in the
Agentfile.

Substrate is a **run-time** choice (`--backend local|bedrock|kubernetes`, with
`--isolation local|container|microvm` scoped to `--backend local`), not an
Agentfile directive — the same artifact runs locally, on Bedrock, or on
Kubernetes unchanged.

## Where to go next

- [What is agentrc?](/docs/what-is-agentrc/) — the model and who it serves.
- [The Agentfile](/docs/agentfile/) — the four keywords and the `/mnt` layout.
- [Specification](/spec/) — the full normative grammar and label tables.
- [Security](/docs/security/) — label-based boundaries and platform Cedar enforcement.
