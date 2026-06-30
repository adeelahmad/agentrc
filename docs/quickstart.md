---
layout: doc
title: Quickstart
description: "Author, build, inspect, push, and run a minimal agentrc agent in five steps."
permalink: /docs/quickstart/
---
# Quickstart

This walks you from a blank file to a running agent in five steps: author an
**Agentfile**, build it (two ways), read the `org.agentrc.*` **labels** the build
emits, push the artifact, and run it. The Agentfile is Dockerfile-shaped — if you
know `docker build`, you already know most of this. For the full normative
grammar, see the [specification](/spec/).

## 1. Create an `Agentfile`

Four new keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) sit on top of
standard Dockerfile keywords. Save this as `Agentfile`:

```dockerfile
# syntax=agentrc.io/agentfile:v1

# --- Who the agent is -------------------------------------------------------
IDENTITY name=hello version=1.0 author=you
IDENTITY description="A concise local assistant"

# --- What it supports -------------------------------------------------------
CAPABILITY text
CAPABILITY streaming

# --- System prompt / objective ----------------------------------------------
SOP You are a concise local assistant. Answer in one short paragraph.

# --- Invocation surface (any framework) -------------------------------------
CMD python ./agent.py

# --- A local tool, embedded into the /mnt tree ------------------------------
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# --- Typed requests to the platform -----------------------------------------
POLICY model.name claude-opus-4
POLICY network dns:api.github.com:443
```

Each `POLICY` line is a **request**, not a guarantee — you are asking the platform
for a model and an egress host. The platform decides what to honour (step 5). The
tool is just an executable placed under `/mnt/tools/`; its destination path under
`/mnt` is what makes it a tool.

## 2. Build the artifact

There are two front doors to the same compiler — they produce **identical** OCI
artifacts.

**BuildKit frontend (nothing to install).** The `# syntax=` line routes the file
through the agentrc frontend image, so a stock Docker / BuildKit install needs no
extra tooling:

```bash
docker build -f Agentfile -t ghcr.io/you/hello:1.0 .
```

**Native `agentrc` CLI** (alias `arc`):

```bash
arc build -t ghcr.io/you/hello:1.0 .
```

Either path parses the four agentrc keywords plus the `ADD --remote` extension,
embeds any `--cached` resources as layers, and writes the `org.agentrc.*` labels.
The output is the same content-addressed OCI artifact regardless of which you use.

## 3. Read the labels

The build translates your authored intent into namespaced OCI image labels under
`org.agentrc.*`. **The platform reads these labels — it never parses the
Agentfile.** Inspect them with `docker inspect` or `arc inspect`:

```bash
arc inspect ghcr.io/you/hello:1.0
```

The Agentfile above emits labels like:

```text
org.agentrc.identity.name=hello
org.agentrc.capability.text=true
org.agentrc.model.name=claude-opus-4
org.agentrc.network.dns.api.github.com=443
org.agentrc.tool.file_read=local
org.agentrc.sop=/mnt/SOP
org.agentrc.sop.sha256=<digest>
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

```bash
arc run ghcr.io/you/hello:1.0 --isolation microvm
```

At run time the platform pulls the artifact, reads every `org.agentrc.*` label,
and evaluates each request — model, egress, tools — against its own organization
policy and available resources, then **grants, narrows, or rejects** it.
Enforcement is **Cedar, platform-side**: the platform compiles the granted
requests plus its own org rules into a Cedar `PolicySet` and evaluates with
deny-by-default and `forbid` over `permit`. You never write Cedar in the
Agentfile.

Substrate is a **run-time** choice (`--isolation local|container|microvm`,
`--substrate <driver>`), not an Agentfile directive — the same artifact runs
locally, in a container, or in a microVM unchanged.

## Where to go next

- [What is agentrc?](/docs/what-is-agentrc/) — the model and who it serves.
- [The Agentfile](/docs/agentfile/) — the four keywords and the `/mnt` layout.
- [Specification](/spec/) — the full normative grammar and label tables.
- [Security](/docs/security/) — label-based boundaries and platform Cedar enforcement.
