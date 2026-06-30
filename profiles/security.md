---
layout: doc
title: Security
description: "Security profile."
permalink: /profiles/security/
---
# Security

**Status:** Working Draft  
**Version:** 0.1.0-draft.4

## Purpose

This profile defines how agentrc packages declare security boundaries and how Cedar is used as the default policy language.

agentrc does not provide isolation itself. A compatible runner enforces declared boundaries or fails closed.

## Boundary classes

| Boundary | Agentfile directive(s) | Example action |
|---|---|---|
| Tool access | `TOOL`, `ALLOW`, `DENY`, `POLICY` | `tool.invoke` |
| Function access | `FUNCTION`, `POLICY` | `function.invoke` |
| Skill access | `SKILL`, `POLICY` | `skill.invoke` |
| Credential access | `CRED`, `BROKER`, `POLICY` | `cred.resolve` |
| Network egress | `URL`, `SERVER`, `MCP`, `POLICY` | `network.egress` |
| Filesystem access | `BIND`, `MOUNT`, `POLICY` | `filesystem.read`, `filesystem.write` |
| Memory access | `MEMORY`, `POLICY` | `memory.read`, `memory.write` |
| Rate limits | `RATELIMIT`, `LIMIT` | implementation-defined |
| Audit | `AUDIT` | lifecycle/event emission |

## Deny-by-default

Under this profile, undeclared capability access is denied.

A runner MUST NOT provide:

1. undeclared tool access;
2. undeclared credential access;
3. undeclared filesystem access;
4. undeclared network egress;
5. undeclared MCP server access.

## Fail-closed behavior

A runner MUST fail closed when:

1. policy cannot be parsed;
2. policy cannot be evaluated;
3. a required boundary cannot be enforced;
4. a required audit stream cannot be produced;
5. a declared credential cannot be resolved;
6. a deny rule conflicts with an allow rule.

## Cedar request shape

A Cedar authorization decision evaluates:

```text
Principal · Action · Resource · Context
```

Recommended entities:

```text
AgentRC::Agent::"<agent-name>"
AgentRC::Tool::"<tool-name>"
AgentRC::Function::"<function-name>"
AgentRC::Skill::"<skill-name>"
AgentRC::Credential::"<credential-name>"
AgentRC::Host::"<host>"
AgentRC::Path::"<path>"
AgentRC::Memory::"<memory-name>"
```

Recommended actions:

```text
tool.invoke
function.invoke
skill.invoke
cred.resolve
network.egress
filesystem.read
filesystem.write
memory.read
memory.write
mcp.request
agent.delegate
```

## Example Cedar policy

```cedar
permit(
  principal == AgentRC::Agent::"github-assistant",
  action == AgentRC::Action::"tool.invoke",
  resource == AgentRC::Tool::"http_request"
)
when { context.url.startsWith("https://api.github.com") };

forbid(
  principal,
  action == AgentRC::Action::"filesystem.write",
  resource
)
when { context.path.startsWith("/etc") };
```

## Credential handling

Credential values MUST NOT appear in:

1. Agentfile source;
2. lockfiles;
3. OCI annotations;
4. package config;
5. logs;
6. audit events;
7. error messages.

Only references and credential names may be recorded.

