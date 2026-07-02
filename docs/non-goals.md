---
layout: doc
title: Non-goals
description: "What the agentrc Agentfile specification (0.1.0-draft.5) deliberately does not define."
permalink: /docs/non-goals/
---
# Non-goals

agentrc is a packaging and governance contract: an Agentfile compiles to an
[OCI artifact](/docs/package/) carrying `org.agentrc.*` labels, and a platform
reads those labels to grant, narrow, or reject each request. To keep that
contract small and adoptable, agentrc deliberately does **not** define:

1. **A runtime.** agentrc does not execute the agent. `CMD` names the loop or
   framework; a [platform](/docs/runners/) drives it on a substrate (local,
   container, microVM, …) chosen with `arc run --backend / --isolation` at run
   time, never in the Agentfile.
2. **A sandboxing implementation.** Isolation strength is a platform / substrate
   property. agentrc states *requests* (`POLICY substrate.*`, `POLICY network …`)
   and the platform enforces them; it does not ship a sandbox.
3. **A cloud platform.** agentrc describes the artifact, not where or by whom it
   runs. Any conformant platform may consume it.
4. **A model API.** `POLICY model.*` *requests* a model and its capabilities; the
   platform selects, substitutes, or rejects. agentrc neither hosts models nor
   defines an inference protocol.
5. **An agent framework.** agentrc is framework-neutral. The same four keywords
   wrap any loop the `CMD` starts.
6. **A tool-call wire protocol.** Tools are plain executables under `/mnt/tools/`
   that self-describe (`--agentrc-schema` or a sibling `<tool>.toolspec.json`);
   MCP servers live under `/mnt/mcp/`. agentrc *declares and governs* these, but
   does not define the on-the-wire calling convention — that is MCP's and each
   tool's job.
7. **A proprietary registry.** The built artifact is a standard OCI artifact that
   pushes, pulls, signs, and mirrors through any OCI-compatible registry.
8. **A second, author-facing policy language.** Authors write only typed `POLICY`
   *requests*. [Cedar](https://www.cedarpolicy.com/) is the **platform-side
   enforcement engine and compilation target** — the platform compiles granted
   requests plus its own organization rules into Cedar and evaluates them. Cedar
   `permit` / `forbid` MUST NOT appear in an Agentfile. See
   [§11.2 of the spec](/spec/) and the [Enforcement profile](/profiles/security/).
9. **The A2A (agent-to-agent) protocol.** Capability *exposure* via `IDENTITY` /
   `CAPABILITY` / labels is in scope; the *protocol* by which one agent discovers
   and calls another (Agent Cards, discovery, cross-agent delegation) is
   **deferred to a later version**.
10. **A multi-agent workflow language inside the Agentfile.** An Agentfile
    describes exactly one agent. External orchestration that references packaged
    agents by digest is a separate, non-normative concern — parked for a future
    draft — not part of the Agentfile core.
11. **A secrets manager — credential resolution is deferred and platform-defined.**
    This draft defines no `SECRET`/`CRED` keyword and no secret schema. An agent
    that needs a credential leaves resolution entirely to the platform (Vault /
    broker / env / workload identity); a credential model may come in a later
    version.

agentrc also does **not replace** the standards it builds on — it composes with
them: MCP, Cedar, OCI, OpenTelemetry, Sigstore, SLSA, Docker / BuildKit,
Kubernetes, gVisor, or Firecracker. Where they overlap, agentrc declares and
governs; the underlying standard executes.

The specification should stay narrow enough that many runtimes and clouds can
consume it without treating it as a competitor.
