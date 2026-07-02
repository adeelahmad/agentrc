---
layout: doc
title: Enforcement (Cedar)
description: "The agentrc Enforcement (Cedar) Profile (0.1.0-draft.6): how a platform compiles typed POLICY requests into Cedar and enforces them deny-by-default."
permalink: /profiles/security/
---
# Enforcement (Cedar) Profile

**Version:** 0.1.0-draft.6 — Working Draft  
**Status:** Working Draft (`# syntax=agentrc.agentfile/v0.1`)  
**Date:** 2026-06-30  
**Audience:** security & compliance reviewers, platform / runner authors

> This profile is the normative home of [§11.2 of the specification](/spec/). It
> defines how a conformant **platform** turns the `org.agentrc.*` request labels
> emitted from an Agentfile into a [Cedar](https://www.cedarpolicy.com/)
> `PolicySet`, and the enforcement properties that decision MUST preserve. The
> keywords **MUST**, **MUST NOT**, **SHOULD**, and **MAY** follow
> [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 1. Cedar is platform-side only

Cedar is **not** an Agentfile author surface and **MUST NOT** appear in an
Agentfile. Authors never write `permit`, `forbid`, or `when`. The only
authorization vocabulary an author writes is the typed `POLICY` request
([spec §8](/spec/)):

```dockerfile
POLICY network dns:api.github.com:443
POLICY agent.sub_agents true
POLICY substrate.device /dev/gpu
```

Cedar is the platform's **enforcement engine and compilation target**. The
platform — never the author — compiles each granted request, together with its
own organization rules, into one Cedar `PolicySet` and evaluates the grant. This
keeps exactly **one author surface** (typed `POLICY` requests) and **one
enforcement engine** (Cedar) beneath it, with a single normative mapping between
them. There is no second, author-facing policy language.

> A `POLICY` line is a **request**, not enforcement. The Agentfile expresses
> *intent*; the platform holds *authority* and is free to **grant, narrow, or
> reject** any request.

## 2. The compilation

The platform reads `org.agentrc.*` labels from the OCI image config — it does
**not** parse the Agentfile — and compiles them into Cedar entities and
policies:

```text
author writes          platform compiles to         platform enforces
─────────────          ────────────────────         ─────────────────
POLICY request   ──►   Cedar entities + policies     deny-by-default,
(org.agentrc.*)        (granted request + org rules)  forbid > permit,
                                                       order-independent,
                                                       monotonic across FROM
```

A conformant platform MUST, for each request it elects to grant, derive Cedar
entities / actions / resources and evaluate them under Cedar's semantics
together with the organization's own Cedar policies.

## 3. Request → Cedar mapping (normative)

The **principal** is the agent identity, taken from
`org.agentrc.identity.name`. The **action** and **resource** are derived from
the request label's namespace:

| Request label | Cedar action | Cedar resource |
|---|---|---|
| `org.agentrc.network.dns.<host>=<port>` | `Action::"NetworkEgress"` | `Host::"<host>:<port>"` |
| `org.agentrc.tool.<name>` (a projected tool) | `Action::"tool.invoke"` | `Tool::"<name>"` |
| `org.agentrc.mcp.<name>` | `Action::"mcp.request"` | `MCPServer::"<name>"` |
| `org.agentrc.agent.sub_agents=true` | `Action::"agent.delegate"` | `Agent::*` (capped by `sub_agents.max`) |
| `org.agentrc.substrate.device=<dev>` | `Action::"device.access"` | `Device::"<dev>"` |

These five rows are the normative mapping a conformant platform MUST implement.
Auto-derived egress (a `network.dns.*` label derived from an `agent.hooks.*` or
`agent.interrupt_endpoint` URL, [spec §8.5](/spec/)) maps through the
`NetworkEgress` row exactly like an explicit `network` request — auto-derivation
is an ergonomic convenience, **not** an implicit grant. The platform MUST still
grant it; an un-granted auto-derived egress is denied.

> Secrets are **deferred** in this draft — there is no `SECRET`/`CRED` keyword
> and no `org.agentrc.secret.*` schema; credential resolution is left entirely to
> the platform and is out of scope for now.

## 4. Enforcement properties a conformant platform MUST preserve

These are Cedar's semantics, and they are why Cedar is the engine.

### 4.1 Deny-by-default

Absence of a grant is a denial. A request that is unrecognised, an auto-derived
egress that was not granted, or any action with no matching `permit` MUST be
**denied**. A platform MUST NOT silently honour a request it did not explicitly
grant.

### 4.2 `forbid` overrides `permit`, order-independently

An organization `forbid` (for example, *"no agent may reach the public
internet"* or *"no agent may write to `/etc`"*) MUST defeat any agent request,
**regardless of evaluation order**. A granted `POLICY` request can never widen
past an org `forbid`. Because the org `forbid` is platform-side Cedar — not in
the Agentfile — the author cannot author around it.

### 4.3 Monotonic composition across `FROM`

When an agent inherits from another agent (`FROM another-agent`,
[spec §2](/spec/)), the effective authorization is the **intersection** of
ceilings:

- A child's granted set MUST NOT exceed its parent's.
- A parent `forbid` is **un-loosenable** by the child.

This is the runtime realization of the spec's *capabilities compose additively,
policy composes by tightening only* rule. Capabilities accumulate; authority
only narrows.

## 5. Two-tier governance, one PolicySet

The agent's `POLICY` requests are the **floor of intent**; the organization's
own Cedar policies are the **ceiling of authority**. The platform compiles both
into a single Cedar `PolicySet` and evaluates the grant:

- **Floor (agent, in the artifact).** What the agent *asks* to do, as typed
  `POLICY` requests compiled to `org.agentrc.*` labels. The author owns this and
  never writes Cedar.
- **Ceiling (organization, out-of-band).** What the org *permits or forbids*,
  authored as Cedar policies by the security team, separate from any Agentfile.
  The author never sees or touches it.

The effective grant is the request **evaluated against** the org ceiling. Where
they disagree, `forbid` over `permit` and deny-by-default resolve it in favour
of the tighter outcome.

## 6. Fail-closed behaviour

A conformant platform MUST fail closed — refuse to grant, refuse to project the
resource, or refuse to boot the agent — when it cannot uphold a required
constraint. Specifically, it MUST fail when:

1. a request label cannot be parsed or recognised against the
   [label catalog](/profiles/oci-package/);
2. the compiled Cedar `PolicySet` cannot be evaluated;
3. a required boundary (network, device, tool/MCP, sub-agent) cannot be
   enforced on the chosen substrate;
4. a `--runtime --fail-if-unavailable` resource cannot be fetched at bootstrap
   ([spec §4.3](/spec/));
5. an org `forbid` and an agent request would otherwise resolve to an unsafe
   grant.

Failing closed means denying access or refusing to boot — never degrading to an
ungoverned execution.

## 7. Platform conformance requirements

To claim conformance to this profile (`agentrc/enforcement-cedar/v0.1`), a
platform MUST:

1. Read authorization from `org.agentrc.*` labels only; it MUST NOT require or
   parse the Agentfile source.
2. Take the principal from `org.agentrc.identity.name` and derive actions /
   resources per the mapping in §3 for every request it grants, including
   auto-derived egress.
3. Compile granted requests together with the organization's Cedar policies into
   a single `PolicySet` and decide under Cedar semantics.
4. Preserve deny-by-default, order-independent `forbid` over `permit`, and
   monotonic intersection across `FROM` (§4).
5. Fail closed on every condition in §6.
6. SHOULD emit an auditable record when it narrows, rejects, or substitutes a
   request (override auditability is
   [open decision #4](/spec/); surfaced, not yet resolved).

The full platform runtime contract — pulling the artifact, fetching `--runtime`
resources, substituting embedded resources via `.origin`, projecting `/mnt`, and
booting `CMD` — is the [Platform Conformance Profile](/profiles/runner-conformance/).
The label namespace this profile reads is catalogued in the
[OCI Labels &amp; Package Profile](/profiles/oci-package/).
