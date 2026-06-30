---
layout: doc
title: OCI Package
description: "OCI Package profile."
permalink: /profiles/oci-package/
---
# OCI Package

**Status:** Working Draft  
**Version:** 0.1.0-draft.4

## Purpose

This profile defines how an Agentfile source tree is packaged for registry distribution.

An agentrc package is a portable agent recipe. It is not a live VM snapshot and not a runtime-specific checkpoint.

## Source tree

A source tree SHOULD contain:

```text
.
├── Agentfile
├── agentrc.lock
├── policy/
│   └── agent.cedar
├── functions/
├── skills/
├── tools/
└── src/
```

## Package contents

A package SHOULD include:

1. Agentfile source;
2. lockfile;
3. package config JSON;
4. policy source or compiled policy profile;
5. bundled skills;
6. bundled functions/source where applicable;
7. tool metadata;
8. provenance metadata;
9. signatures where available.

## Recommended media types

| Component | Media type |
|---|---|
| Manifest | `application/vnd.oci.image.manifest.v1+json` |
| Agent config | `application/vnd.agentrc.agent.config.v1+json` |
| Agentfile | `application/vnd.agentrc.agent.agentfile.v1+text` |
| Lockfile | `application/vnd.agentrc.agent.lock.v1+json` |
| Cedar policy | `application/vnd.agentrc.agent.policy.cedar.v1+text` |
| Source layer | `application/vnd.agentrc.agent.source.layer.v1.tar+gzip` |
| Tool layer | `application/vnd.agentrc.agent.tool.layer.v1.tar+gzip` |
| Skill layer | `application/vnd.agentrc.agent.skill.layer.v1.tar+gzip` |

## Required annotations

A package SHOULD carry:

```text
io.agentrc.agentfile.version
io.agentrc.agent.name
io.agentrc.policy.hash
io.agentrc.base.digest
```

A package MAY also carry:

```text
io.agentrc.risk.tier
io.agentrc.required.runner.profiles
io.agentrc.required.policy.profile
org.opencontainers.image.title
org.opencontainers.image.version
org.opencontainers.image.source
org.opencontainers.image.revision
org.opencontainers.image.created
```

## Registry operations

A registry client SHOULD support:

```bash
agentrc push <ref>
agentrc pull <ref>
agentrc inspect <ref>
agentrc verify <ref>
agentrc policy show <ref>
agentrc diff <ref-a> <ref-b>
```

## Reproducibility

Packages SHOULD be reproducible from:

1. Agentfile;
2. lockfile;
3. source tree;
4. pinned base digest;
5. pinned tool/function/skill/policy digests.

