---
layout: doc
title: Runner Conformance
description: "Runner Conformance profile."
permalink: /profiles/runner-conformance/
---
# Runner Conformance

**Status:** Working Draft  
**Version:** 0.1.0-draft.2

## Purpose

This profile describes what a runner must do if it claims compatibility with AgentRC packages.

AgentRC itself is not the runner. This profile exists so independent runtimes, clouds, CLIs, or sandboxes can state what level of AgentRC support they provide.

## Required disclosure

A runner SHOULD publish a support statement containing:

1. supported Agentfile version;
2. supported directives;
3. supported policy profiles;
4. supported credential backends;
5. supported isolation/backing substrate, if any;
6. supported audit/export formats;
7. unsupported directives and failure behavior;
8. known security limitations.

## Core requirements

A runner claiming AgentRC Runner Profile conformance MUST:

1. read an AgentRC package or Agentfile source;
2. validate the Agentfile or fail with a diagnostic;
3. resolve and verify the lockfile where present;
4. execute `CMD` or fail with unsupported entrypoint;
5. enforce supported security boundaries;
6. fail closed on unsupported required security boundaries;
7. resolve credentials only at runtime;
8. redact credential values;
9. emit audit records when required;
10. expose effective support/limits to the operator.

## Runner is not the spec

A runner may use any substrate:

```text
local process
container
Docker
containerd
gVisor-style sandbox
microVM
Kubernetes job
serverless worker
managed cloud agent runtime
SSH remote runner
framework-native adapter
```

The substrate does not change the Agentfile semantics.

## Placement directives

`ISOLATION`, `IMAGE`, `SLICE`, `PLUGIN`, and `BACKEND` are treated as requested runner capabilities. A portable package SHOULD avoid hard-coding placement unless necessary.

A future companion document may define a separate run manifest for placement.

