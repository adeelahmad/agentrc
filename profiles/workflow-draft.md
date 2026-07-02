---
layout: doc
title: Workflow Draft
description: "A deferred, non-normative sketch for orchestrating multiple packaged agentrc agents by digest — distinct from the (also deferred) A2A protocol."
permalink: /profiles/workflow-draft/
published: false
---
<!-- Parked 2026-07-03 — will return in a future draft -->
# Workflow Draft

**Status:** Deferred, non-normative companion sketch  
**Version:** 0.1.0-draft.6 (Working Draft) — not part of the Agentfile core  
**Date:** 2026-06-30

> This page is **not** part of the [Agentfile specification](/spec/). It is a
> forward-looking sketch, surfaced so the design space is visible, not silently
> resolved. Nothing here is normative, and a conformant platform is not required
> to implement any of it.

## Scope and what is deferred

An Agentfile describes **one** agent: its identity, capabilities, system prompt,
the resources it carries under `/mnt`, and the typed `POLICY` requests it makes
of the platform. Multi-agent orchestration is deliberately **out of the core**.

Two related but distinct things are deferred to a later version:

1. **The A2A (agent-to-agent) protocol** — how one running agent *discovers* and
   *calls* another (Agent Cards, discovery, live cross-agent delegation, and the
   governance algebra across an agent-to-agent call). This is **out of scope for
   0.1.0-draft.6**. Note that capability *exposure* is already in scope: an agent advertises
   what it is and what it does through `IDENTITY`, `CAPABILITY`, and the resulting
   `org.agentrc.*` labels. What is deferred is the *protocol* for one agent to
   find and invoke another.
2. **This workflow companion** — an external orchestrator that runs several
   *packaged* agents in sequence or in parallel. It references each agent by its
   OCI **digest** and drives them from the outside; it does **not** embed
   Agentfiles, and it is **not** the A2A protocol.

The workflow companion sits *above* packaged agents. The A2A protocol would sit
*between* running agents. Both are excluded from the Agentfile core; this page
sketches only the former.

## Design principles

1. **Workflows reference packaged agents by digest.** A workflow step names an
   agent as an OCI reference pinned to a content digest
   (`ghcr.io/bank/claims-triage@sha256:…`). Workflows never embed an Agentfile or
   re-author an agent's `POLICY` requests.
2. **The state machine owns control flow.** Retry, catch, timeout, branching, and
   parallelism live in the workflow definition, not inside any single agent.
3. **The workflow engine is an external orchestrator.** It is a kind of platform
   that invokes packaged agents; it is **not** part of the Agentfile core, and it
   does not change how an individual agent is built or governed.
4. **Authority composes by tightening.** Each agent invoked in a workflow is
   still subject to the platform reading its `org.agentrc.*` labels and
   granting / narrowing / rejecting each request, enforced platform-side via
   Cedar (deny-by-default, `forbid` over `permit`, monotonic across `FROM`). A
   workflow-level ceiling, if present, can only **narrow** what each step's agent
   is allowed to do — never widen it. See the
   [Enforcement (Cedar) profile](/profiles/security/) and the
   [Platform conformance profile](/profiles/runner-conformance/).

## Example

An ASL-inspired YAML state machine that orchestrates three packaged agents, each
referenced by digest. This is illustrative only.

```yaml
apiVersion: agentrc.workflow/v0.1
kind: AgentWorkflow
metadata:
  name: claims-review-flow
spec:
  startAt: Triage
  states:
    Triage:
      type: Task
      agent: ghcr.io/bank/claims-triage@sha256:...
      inputPath: $.claim
      resultPath: $.triage
      timeoutSeconds: 120
      retry:
        - errorEquals: [Platform.Timeout, Tool.TemporaryFailure]
          intervalSeconds: 2
          maxAttempts: 3
          backoffRate: 2.0
      catch:
        - errorEquals: [States.ALL]
          resultPath: $.error
          next: HumanReview
      next: RiskChoice

    RiskChoice:
      type: Choice
      choices:
        - variable: $.triage.risk
          stringEquals: high
          next: HumanReview
        - variable: $.triage.risk
          stringEquals: low
          next: AutoApprove
      default: HumanReview

    HumanReview:
      type: Task
      agent: ghcr.io/bank/human-review-assistant@sha256:...
      end: true

    AutoApprove:
      type: Task
      agent: ghcr.io/bank/approval-agent@sha256:...
      end: true
```

Each `agent:` value is a packaged OCI artifact pinned by digest. When the engine
runs a step, the underlying platform pulls that artifact, reads its
`org.agentrc.*` labels, and grants / narrows / rejects each request before
booting `CMD` — exactly as it would for a standalone run. The workflow adds
control flow around those agents; it does not bypass their governance.

## Open questions

1. Should this align closely with Amazon States Language, or define a smaller,
   portable subset?
2. How should a workflow-level policy ceiling compose with each agent package's
   own `POLICY` requests? (Expected answer: tightening-only — the workflow
   ceiling can narrow but not widen, consistent with monotonic `FROM`
   composition.)
3. How should per-step inputs and outputs be passed between agents — by value in
   the workflow state, or by reference to an external store?
4. How should human-in-the-loop steps be represented as states?
5. When the A2A protocol arrives, where is the boundary between *live agent
   delegation* (A2A, agent-to-agent) and *external orchestration of packages*
   (this workflow companion)?

## Relationship to the core

| This page | The Agentfile core ([spec](/spec/)) |
|---|---|
| Orchestrates **many** packaged agents | Describes **one** agent |
| References agents by **OCI digest** | Is the recipe that **builds** the artifact |
| External engine drives control flow | Platform reads labels, grants/narrows/rejects |
| Deferred, non-normative | Normative |

For the agent-to-agent *protocol* — the other deferred piece, distinct from this
companion — see the deferred-feature note in the
[workflows overview](/docs/workflows/).
