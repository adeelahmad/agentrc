---
layout: doc
title: Compatible runners
description: "Runner conformance without making agentrc a runtime."
permalink: /docs/runners/
---
# Compatible runners

A runner executes an agentrc package. agentrc itself is not the runner.

## Runner examples

- local development runner
- Docker/OCI container runner
- gVisor-backed runner
- Firecracker or microsandbox microVM runner
- Kubernetes job runner
- serverless runner
- managed cloud agent runtime adapter
- framework-native adapter for Strands, LangChain, CrewAI, Langflow, or similar

## Conformance principle

A runner may claim conformance only to the profiles it implements.

| Profile | Runner must do |
|---|---|
| Core | parse Agentfile and reject unknown required directives |
| Security | enforce declared boundaries or fail closed |
| Package | consume agentrc package metadata and lockfile |
| Tool projection | expose declared tools through the specified surface |
| Audit | emit required audit records |

## The key boundary

agentrc defines what must be true. Runners decide how to make it true.

This lets AWS, Google, Docker, gVisor, Firecracker, microsandbox, and local runners implement their own execution layer without owning the portable agent declaration format.
