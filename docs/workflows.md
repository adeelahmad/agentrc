---
layout: doc
title: Workflow draft (deferred)
description: "A deferred, non-normative companion: orchestrating packaged agentrc agents externally by digest. Not part of the Agentfile core."
permalink: /docs/workflows/
published: false
---
<!-- Parked 2026-07-03 — will return in a future draft -->
# agentrc Workflow draft (deferred companion)

> **Status: deferred and non-normative.** Nothing on this page is part of the
> [0.1.0-draft.5 Agentfile specification](/spec/). It sketches a *future* companion layer
> and is published only so the boundary is explicit.

An **Agentfile declares one agent.** Multi-agent applications need a separate
layer that references already-packaged agents and orchestrates them. That layer
is **out of scope for this version** and is described here only as direction.

## What is in scope vs. deferred

There are two distinct things people mean by "multi-agent," and v0.1 draws a hard
line between them:

| Concern | Status in v0.1 |
|---|---|
| **Capability *exposure*** — declaring who the agent is and what modalities / patterns it supports, via `IDENTITY`, `CAPABILITY`, and the resulting `org.agentrc.*` labels. | **In scope.** Part of the [spec](/spec/). |
| **The A2A *protocol*** — Agent Cards, agent discovery, cross-agent delegation, and the governance algebra across an agent-to-agent call. | **Deferred.** Not in this version (see the spec's [deferred list](/spec/)). |
| **External workflow orchestration** — a state machine that drives several packaged agents, referenced by digest, from outside any one Agentfile. | **Deferred, non-normative.** Sketched on this page; **not** part of the Agentfile core and **distinct** from the A2A protocol. |

The Agentfile exposes capability through labels; it does **not** define how one
agent finds and calls another. The workflow companion below is also **not** the
A2A protocol — it orchestrates packaged agents externally rather than letting
agents discover and delegate to each other at run time.

## Direction

The companion workflow profile would reference packaged agents by an **immutable
registry reference or digest** — the same OCI artifacts produced by
[`arc build`](/cli/) — and define state-machine orchestration *separately* from
agent packaging. Each referenced agent is still an ordinary agentrc artifact: the
platform reads its `org.agentrc.*` labels and grants, narrows, or rejects its
requests, enforcing them via Cedar (platform-side, [deny-by-default](/profiles/security/)),
exactly as it would for a standalone agent. The workflow layer adds no new
authority of its own; it only sequences agents the platform already governs.

A future workflow file might look like this:

```yaml
# Deferred, non-normative companion — NOT part of the Agentfile core.
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

Each `agent` value is a content-addressed reference to a packaged agentrc
artifact. The workflow engine pulls each one, and the platform applies the same
label-based grant / enforcement flow per agent — there is no implicit trust
between steps and no shared, widened authority.

## Why this is not in the Agentfile

Workflows require state, retries, branching, compensation, timeout behaviour,
failure semantics, and event handling across **multiple** agents. That is a
different problem from packaging **one** agent into a reviewable, governed OCI
artifact. Folding it into the Agentfile would blur the single-agent contract that
makes an Agentfile easy to review and a platform easy to build.

agentrc keeps those layers separate: the Agentfile and its `org.agentrc.*` labels
describe and govern one agent; any cross-agent orchestration lives in a separate,
later companion — and the agent-to-agent *protocol* itself remains deferred.

## Where to go next

- [Specification](/spec/) — the normative v0.1 Agentfile model.
- [What agentrc is not](/docs/non-goals/) — why the workflow layer and the A2A
  protocol are out of scope.
- [Workflow draft profile](/profiles/workflow-draft/) — the (deferred) profile
  notes for this companion.
