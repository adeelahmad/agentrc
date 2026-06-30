---
layout: doc
title: Agentfile
description: "The single-agent declaration format."
permalink: /docs/agentfile/
---
# Agentfile

The Agentfile is a Dockerfile-like declaration for one agent.

## Required idea

The Agentfile does not describe a cluster, fleet, workflow, or runtime implementation. It describes one agent and the contract a runner must satisfy before executing it.

## Core directives

| Directive | Purpose |
|---|---|
| `AGENT` | Stable agent identity |
| `FROM` | Optional base agent package or base image reference |
| `CMD` | Entrypoint command |
| `TOOL` | Declared tool capability |
| `FUNCTION` | Typed function capability |
| `SKILL` | Portable skill/procedure bundle |
| `SERVER` | MCP server or compatible tool server |
| `MOUNT` | Filesystem boundary |
| `URL` | Network egress boundary |
| `CRED` | Host-scoped credential reference (placeholder-substitution model) |
| `LIMIT` | Resource, time, or output cap |
| `RATELIMIT` | Invocation or egress rate cap |
| `AUDIT` | Audit capture level |
| `SOP ... END` | Embedded Standard Operating Procedure (instructions) |
| `POLICY ... END` | Inline policy block |

## Minimal pattern

```dockerfile
AGENT <name>
CMD <command>
TOOL <tool-ref>
AUDIT basic
POLICY
  <policy statements>
END
```

## Why single-agent only?

Keeping Agentfile single-agent prevents it from becoming a workflow language. Multi-agent applications should use a companion workflow specification that references packaged agents by digest.
