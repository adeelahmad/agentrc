---
layout: doc
title: Conformance
description: "Conformance profile summary."
permalink: /docs/conformance/
---
# Conformance

AgentRC conformance is profile-based.

A tool, package builder, registry, runner, or workflow engine should state exactly which profiles it supports.

## Profile names

- `agentrc/core-agentfile/v0.1`
- `agentrc/security-cedar/v0.1`
- `agentrc/oci-package/v0.1`
- `agentrc/tool-projection/v0.1`
- `agentrc/runner/v0.1`
- `agentrc/workflow/v0.1` draft

## Why profiles?

A local validator should not need to implement microVM isolation. A registry should not need to evaluate every runtime boundary. A runner should not need to become a workflow engine.

Profiles keep the spec implementable.
