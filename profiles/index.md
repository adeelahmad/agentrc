---
layout: doc
title: Profiles
description: "agentrc conformance profiles for the v0.1 Agentfile model: core, enforcement, OCI labels, projection, platform, and the deferred workflow companion."
permalink: /profiles/
---

# Profiles

agentrc uses profiles to keep the core specification small while letting concrete
implementations prove exactly what they conform to. The [specification](/spec/)
defines a Dockerfile-shaped Agentfile, four new keywords, and a single
`org.agentrc.*` label namespace; each profile below pins down one slice of that
contract so compilers, platforms, and registries can be tested independently.

- [Core](/profiles/core/) — the Agentfile core: the four new keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) plus standard Dockerfile keywords, compiled to `org.agentrc.*` labels and OCI layers.
- [Security](/profiles/security/) — Enforcement (Cedar): the platform-side enforcement engine and compilation target for typed `POLICY` requests; deny-by-default, `forbid` over `permit`, tightening-only `FROM`.
- [OCI Package](/profiles/oci-package/) — OCI labels &amp; package: the `org.agentrc.*` label namespace, layers carrying the `/mnt` resources, media types, and `*.origin` overrides.
- [Tool Projection](/profiles/tool-projection/) — the `/mnt` projection: how `tools/`, `skills/`, `mcp/`, `proc/`, and `SOP` are presented to the running agent.
- [Runner Conformance](/profiles/runner-conformance/) — Platform conformance: what a platform MUST do with the labels — read them (never the Agentfile), grant / narrow / reject, and fail closed.
- Workflow orchestration is parked for a future draft — a deferred, non-normative companion for orchestrating packaged agents by digest; distinct from the deferred A2A protocol and not part of the Agentfile core.
