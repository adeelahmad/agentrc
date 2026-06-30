---
layout: doc
title: Security model
description: "Declarative security boundaries for AgentRC packages."
permalink: /docs/security/
---
# Security model

AgentRC decouples security intent from runtime implementation.

## Boundary types

| Boundary | Declared as | Enforced by |
|---|---|---|
| Tool access | `TOOL`, `POLICY` | runner/tool gateway |
| Filesystem access | `MOUNT`, `POLICY` | runner/sandbox/filesystem layer |
| Network egress | `URL`, `POLICY` | runner/network policy |
| Secrets | `SECRET` / `CRED`, `POLICY` | runner secret broker |
| Rate limits | `RATELIMIT` | runner/gateway |
| Resource limits | `LIMIT` | runner/substrate |
| Audit | `AUDIT` | runner/audit sink |

## Policy profile

Cedar is the default policy profile because it is explicit, machine-evaluable, and well-suited to deny-by-default authorization.

```cedar
permit(
  principal == AgentRC::Agent::"code-reviewer",
  action == AgentRC::Action::"tool.invoke",
  resource == AgentRC::Tool::"file_read"
) when {
  context.path like "/workspace/*"
};
```

## Runner obligation

A runner claiming security-profile conformance must fail closed when it cannot understand, translate, or enforce a declared boundary.

<div class="callout warning">
<strong>Important:</strong> AgentRC declares the security contract. It does not claim that every runner can enforce every boundary. That is why runner conformance is profile-based.
</div>
