---
layout: doc
title: Workflow draft
description: "Future multi-agent workflow layer."
permalink: /docs/workflows/
---
# AgentRC Workflow draft

The Agentfile declares one agent. Multi-agent applications need a separate layer.

## Direction

The companion workflow profile should reference packaged agents by immutable registry reference or digest and define state-machine orchestration separately from agent packaging.

A future workflow file may look like this:

```yaml
version: agentrc.workflow/v0.1
name: review-and-notify

agents:
  reviewer: ghcr.io/org/code-reviewer@sha256:...
  notifier: ghcr.io/org/slack-notifier@sha256:...

states:
  Review:
    type: task
    agent: reviewer
    next: Notify
    retry:
      max_attempts: 2

  Notify:
    type: task
    agent: notifier
    end: true
```

## Why not put this in Agentfile?

Workflows require state, retries, branching, compensation, timeout behavior, failure semantics, and event handling. That is a different problem from single-agent packaging.

AgentRC keeps those layers separate.
