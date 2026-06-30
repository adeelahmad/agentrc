---
layout: doc
title: Package model
description: "What an agentrc package contains."
permalink: /docs/package/
---
# Package model

An agentrc package is the portable artifact built from an Agentfile.

## Package contents

| Component | Purpose |
|---|---|
| `Agentfile` | Source declaration |
| `agentrc.lock` | Resolved and pinned dependencies |
| package manifest | package identity, media types, annotations |
| policy bundle | compiled or source policy |
| tool metadata | declared tool refs and schemas |
| function artifacts | code or references for functions |
| skill artifacts | portable procedural bundles |
| runner hints | optional requirements, not runtime identity |
| provenance | signatures, digests, builder metadata |

## Registry model

agentrc packages are intended to be shared through OCI-compatible registries.

```bash
agentrc build --tag ghcr.io/org/code-reviewer:1.0.0
agentrc push ghcr.io/org/code-reviewer:1.0.0
agentrc inspect ghcr.io/org/code-reviewer:1.0.0
```

## Review before run

A package consumer should be able to inspect a package before execution and answer:

- What command will start?
- Which tools can be invoked?
- Which paths can be read or written?
- Which hosts can be contacted?
- Which secrets are required?
- Which policies apply?
- Which digests are pinned?
