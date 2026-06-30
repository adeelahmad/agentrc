---
layout: doc
title: Workflow Draft
description: "Workflow Draft profile."
permalink: /profiles/workflow-draft/
---
# Workflow Draft

**Status:** Non-normative sketch  
**Version:** 0.1.0-draft.4

## Purpose

The Agentfile describes one agent. Multi-agent workflows should be specified separately.

This draft sketches an ASL-inspired YAML format for orchestrating multiple packaged agents. It is intentionally not part of the core Agentfile specification.

## Design principles

1. Workflows reference packaged agents; they do not embed Agentfiles.
2. Workflow state machines should handle retry, catch, timeout, branching, and parallelism.
3. Workflow engines are runners/orchestrators, not part of the Agentfile core.
4. Security ceilings should compose across workflow boundaries.

## Example

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
        - errorEquals: [Runner.Timeout, Tool.TemporaryFailure]
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

## Open questions

1. Should this align closely with Amazon States Language or define a smaller portable subset?
2. How should workflow-level policy ceilings compose with each agent package policy?
3. Should workflow engines pass credentials to agents or should agents resolve credentials independently?
4. How should human-in-the-loop tasks be represented?
5. Should A2A be used for live agent delegation or should the workflow engine invoke packages directly?

