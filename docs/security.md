---
layout: doc
title: Security model
description: "How agentrc states security boundaries as OCI labels and how the platform enforces them with Cedar."
permalink: /docs/security/
---
# Security model

agentrc separates **what an agent requests** from **what a platform enforces**.
An author writes a Dockerfile-shaped [Agentfile](/docs/agentfile/); the build
emits namespaced `org.agentrc.*` OCI labels; the platform reads **those labels**
— never the Agentfile source — and **grants, narrows, or rejects** each request,
then enforces the result with [Cedar](https://www.cedarpolicy.com/),
platform-side. The Agentfile expresses *intent*; the platform holds *authority*.

> A `POLICY` line or a presence label is a **request**, not a permission.
> Absence of a grant is a denial. The enforcement engine is Cedar, and it lives
> on the platform, not in the artifact.

## Boundaries: declared as a label, enforced by the platform

Every security-relevant boundary is **declared** as a label the compiler emits
from authored intent, and **enforced** by the platform's Cedar engine. Authors
write the short forms; the compiler prepends `org.agentrc.`.

| Boundary | Declared as (authored → label) | Enforced by |
|---|---|---|
| **Network egress** | `POLICY network dns:api.github.com:443` → `org.agentrc.network.dns.api.github.com=443` | Platform Cedar: `Action::"NetworkEgress"` on `Host::"<host>:<port>"` |
| **Tools** | `COPY ./tools/x /mnt/tools/x` → `org.agentrc.tool.x=local` (presence label) | Platform Cedar: `Action::"tool.invoke"` on `Tool::"<name>"` |
| **MCP servers** | `ADD --remote … /mnt/mcp/x` → `org.agentrc.mcp.x=<digest|runtime:url>` (+ `.origin`) | Platform Cedar: `Action::"mcp.request"` on `MCPServer::"<name>"` |
| **Devices** | `POLICY substrate.device /dev/gpu` → `org.agentrc.substrate.device=/dev/gpu` | Platform Cedar: `Action::"device.access"` on `Device::"<dev>"` |
| **Sub-agents** | `POLICY agent.sub_agents true` → `org.agentrc.agent.sub_agents=true` (capped by `sub_agents.max`) | Platform Cedar: `Action::"agent.delegate"` on `Agent::*` |

The platform sees a uniform, machine-readable manifest of *everything the agent
asks for* — tools, network, devices, MCP servers, and sub-agents — and vets it
before the agent runs.

> Secrets are **deferred** in this draft — no `SECRET`/`CRED` keyword and no
> `org.agentrc.secret.*` schema; credential resolution is platform-defined and
> out of scope for now.

## Auto-derived egress is explicit, not implicit

When a `POLICY` value is a URL the agent will call out to — `agent.hooks.*` or
`agent.interrupt_endpoint` — the compiler **auto-derives** a corresponding
`network` egress label and **attributes** it, so a webhook can never open a
silent network hole:

```text
org.agentrc.network.dns.hooks.internal=443
org.agentrc.network.dns.hooks.internal.source=auto:agent.hooks.pre
```

Auto-derivation is an ergonomic convenience, **not** an implicit grant. The
platform must still grant the derived egress; an un-granted one is denied.

## Principles

- **Requests, not enforcement.** A `POLICY` line or a presence label asks the
  platform for something. The platform decides.
- **Deny-by-default.** Absence of a grant is a denial. An unrecognised request,
  an un-granted auto-derived egress, or an action with no matching grant is
  denied.
- **`forbid` overrides `permit`, order-independently.** An organization `forbid`
  (e.g., "no agent may reach the public internet") defeats any agent request,
  regardless of evaluation order.
- **Tightening-only across `FROM`.** When an agent inherits from another agent,
  its effective authorization is the **intersection** of ceilings — a child
  cannot exceed its parent, and a parent `forbid` is un-loosenable.
- **Fail closed.** A platform that cannot understand, evaluate, or enforce a
  required boundary MUST refuse to run the agent rather than run it
  under-protected.
- **No policy language in the Agentfile.** Authors write only typed `POLICY`
  requests. **Cedar is platform-side only** — there is no `permit` / `forbid`
  in an Agentfile, and no second, author-facing policy language to learn.

## Request → Cedar, in brief

Cedar is the platform's **enforcement engine and compilation target**, not an
author surface. For each granted request the platform derives Cedar entities and
evaluates them. The principal is the agent identity
(`org.agentrc.identity.name`); the action and resource come from the request
namespace:

| Request label | Cedar action | Cedar resource |
|---|---|---|
| `org.agentrc.network.dns.<host>=<port>` | `Action::"NetworkEgress"` | `Host::"<host>:<port>"` |
| `org.agentrc.tool.<name>` | `Action::"tool.invoke"` | `Tool::"<name>"` |
| `org.agentrc.mcp.<name>` | `Action::"mcp.request"` | `MCPServer::"<name>"` |
| `org.agentrc.agent.sub_agents=true` | `Action::"agent.delegate"` | `Agent::*` (capped by `sub_agents.max`) |
| `org.agentrc.substrate.device=<dev>` | `Action::"device.access"` | `Device::"<dev>"` |

The agent's `POLICY` requests are the **floor of intent**; the organization's own
Cedar policies (authored out-of-band by the security team) are the **ceiling of
authority**. The platform compiles both into one Cedar `PolicySet` and evaluates
the grant. The [Enforcement (Cedar) profile](/profiles/security/) is the
normative home of this mapping and the conformance requirements.

<div class="callout warning">
<strong>Important:</strong> agentrc declares the security contract; it does not
claim that every platform can enforce every boundary. That is why platform
conformance is profile-based, and why a conformant platform MUST fail closed on
any boundary it cannot enforce.
</div>
