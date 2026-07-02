---
layout: doc
title: "Demo: one agent, three backends"
description: "Build the code-reviewer agent once, then run the same artifact on three substrates — local microVM, Bedrock, and Kubernetes — via the reference translators."
permalink: /docs/demo/
---
# Demo: one agent, three backends

This walkthrough uses the canonical [code-reviewer Agentfile](/examples/Agentfile.code-reviewer)
to show the core claim of the v0.1 model in one sitting: you build **once**, and
the same OCI artifact — carrying the same `org.agentrc.*` labels — runs on three
different substrates. The substrate is a **run-time** choice (`--backend`), never
an Agentfile directive.

## Build once

```bash
arc build -t ghcr.io/agentrc/code-reviewer:1.0 .
```

This compiles the Agentfile into a standard OCI artifact. The `POLICY` requests
become `org.agentrc.*` labels the platform reads to grant, narrow, or reject each
request. Nothing in the artifact names a substrate.

## Run on three backends

Point `arc run` at the artifact you just built and pick a backend at run time.
The first runs locally under microVM isolation; the other two render the
backend-native deploy form so you can inspect it before applying.

```bash
arc run ghcr.io/agentrc/code-reviewer:1.0 --backend local --isolation microvm
arc run ghcr.io/agentrc/code-reviewer:1.0 --backend bedrock --dry-run
arc run ghcr.io/agentrc/code-reviewer:1.0 --backend kubernetes --dry-run
```

- **`--backend local --isolation microvm`** runs the agent locally, isolated in a
  microVM. `--isolation` is scoped to the local backend.
- **`--backend bedrock --dry-run`** emits the Bedrock deploy form as JSON — pipe it
  to `python3 -m json.tool` to confirm it parses.
- **`--backend kubernetes --dry-run`** emits Kubernetes manifests as YAML — pipe it
  to a YAML parser or `kubeconform` to confirm it parses.

The `--backend` and `--isolation` values change nothing about the artifact: the
labels are identical across all three runs. Each backend has a reference
translator that maps the granted request set onto that substrate's deploy form.

## Why this works

Same artifact, same labels, three substrates. The translators are the proof of concept; the labels are the standard.
