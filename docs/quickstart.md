---
layout: doc
title: Quickstart
description: "Create and publish a minimal AgentRC package."
permalink: /docs/quickstart/
---
# Quickstart

## 1. Create an Agentfile

```dockerfile
# syntax=agentrc.agentfile/v0.1

AGENT hello-local
CMD python ./agent.py

TOOL utcp:file_read
TOOL utcp:current_time
MOUNT /workspace ro
AUDIT basic

POLICY
  permit(
    principal == AgentRC::Agent::"hello-local",
    action == AgentRC::Action::"tool.invoke",
    resource in [AgentRC::Tool::"file_read", AgentRC::Tool::"current_time"]
  );
END
```

## 2. Validate

```bash
agentrc lint Agentfile
```

Validation should answer:

- Is the file syntactically valid?
- Are unknown directives rejected?
- Are security boundaries explicit?
- Are secrets referenced, not embedded?
- Is the policy parseable?

## 3. Build a portable package

```bash
agentrc compile Agentfile
agentrc build --tag ghcr.io/adeelahmad/hello-local:0.1.0
```

The build produces an AgentRC package containing:

- `Agentfile`
- `agentrc.lock`
- package metadata
- policy bundle
- declared tools, skills, functions, and support artifacts

## 4. Publish

```bash
agentrc push ghcr.io/adeelahmad/hello-local:0.1.0
```

## 5. Run with a compatible runner

```bash
agentrc run ghcr.io/adeelahmad/hello-local:0.1.0
```

A compatible runner may use a local process, container, gVisor sandbox, microVM, serverless platform, or managed cloud runtime. The substrate is deliberately outside the core AgentRC specification.
