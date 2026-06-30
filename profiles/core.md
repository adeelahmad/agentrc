---
layout: doc
title: Core
description: "Core Agentfile profile: the minimal normative directive set for v0.1."
permalink: /profiles/core/
---
# Core

**Status:** Working Draft  
**Version:** 0.1.0-draft.4

## Purpose

The Core Agentfile Profile defines the minimum behavior required to parse, validate, and inspect a single-agent Agentfile.

It does not define execution.

A small, stable core is a deliberate design choice. Every directive in the core is something a second, independent implementer must build and a conformance test must cover. The core is therefore kept intentionally narrow so the format can stabilize, while richer capabilities live in clearly separated optional profiles that can evolve without breaking the core.

## Required behavior

An implementation claiming this profile MUST:

1. read a text file named `Agentfile` or an explicitly supplied path;
2. ignore blank lines and comments;
3. parse known directives case-sensitively;
4. preserve directive order for lockfile generation and review;
5. capture `POLICY ... END` and `SOP ... END` blocks without interpreting inner lines as Agentfile directives;
6. reject unknown directives unless an extension profile is enabled;
7. produce a structured parse tree suitable for linting and package generation;
8. report line numbers for parse failures.

## Core directive set

The v0.1 **core** directive set is deliberately minimal — the smallest set that can describe a reviewable, governed, portable single agent:

```text
AGENT  FROM  CMD  TOOL  MOUNT  CRED  URL  POLICY  AUDIT
```

| Directive | Concern | Why it is core |
|---|---|---|
| `AGENT` | Identity | A package needs a stable, signable identity. |
| `FROM` | Inheritance | Establishes the base and the security ceiling that inheritance tightens. |
| `CMD` | Entrypoint | The one thing every runner must know to start the agent. |
| `TOOL` | Capability | The primary capability surface that policy governs. |
| `MOUNT` | Filesystem boundary | A first-class, reviewable security boundary. |
| `CRED` | Credential boundary | Secrets-as-references is a core security guarantee. |
| `URL` | Network boundary | Deny-by-default network egress is a core security guarantee. |
| `POLICY` | Governance | Machine-evaluable boundaries that travel with the package. |
| `AUDIT` | Governance | Fail-closed auditability is a core review guarantee. |

All other recognized directives belong to optional profiles (below). An implementation MAY support any number of them, but it MUST publish a support matrix rather than implying that every recognized directive is enforced.

## Directive classification

Every directive defined by this specification is assigned to exactly one profile. This is what keeps the core small while leaving room to grow.

| Profile | Directives | Status |
|---|---|---|
| **Core Agentfile** | `AGENT` `FROM` `CMD` `TOOL` `MOUNT` `CRED` `URL` `POLICY` `AUDIT` | Normative, v0.1 |
| **Capability Extensions** | `TOOLSET` `FUNCTION` `SKILL` `SERVER` `MCP` `MEMORY` | Optional |
| **Instruction Embedding** | `SOP` | Optional (Agent SOP / Agent Skills) |
| **Security Shorthand & Limits** | `ALLOW` `DENY` `RATELIMIT` `TIMEOUT` `LIMIT` | Optional (lowers into the Security profile) |
| **Placement / Run Manifest** | `ISOLATION` `IMAGE` `SLICE` `BACKEND` `BIND` `BROKER` `PLUGIN` | Optional, non-portable (see below) |
| **Observability** | `TRACE` `HEALTHCHECK` | Optional |
| **Framework & Experimental** | `SHELL` `OPTIMIZER` | Optional / experimental |

## Build vs. placement (Run Manifest)

The directives in the **Placement / Run Manifest** group describe *where and how* an agent runs, not *what it is or what it may do*. Mixing them into the portable declaration breaks the "vet once, run anywhere, content-addressed" promise: two agents identical in capability but different in placement would otherwise produce different artifacts and different digests.

Therefore:

1. The portable, signable artifact (the Agentfile and its package) SHOULD be computed from capability and boundary directives only.
2. Placement directives (`ISOLATION`, `IMAGE`, `SLICE`, `BACKEND`, `BIND`, `BROKER`, `PLUGIN`) SHOULD be expressible in a separate **Run Manifest** chosen by the operator or runner, and SHOULD NOT contribute to the package digest.
3. For backward compatibility, current implementations MAY still accept these directives inline; when they do, a linter SHOULD warn (see below) and the package builder SHOULD exclude them from the content-addressed identity.

This mirrors the Dockerfile (build) vs. Compose/run (placement) separation and is load-bearing for reproducibility and signing.

## Compatibility notes

The current agentrc implementation recognizes the full directive family above, though not every directive has full semantic enforcement. This profile defines target semantics; implementations should publish a support matrix rather than implying all recognized directives are enforced.

## Validation recommendations

A linter SHOULD warn when:

1. `AGENT` is missing in a package intended for publication;
2. `FROM` uses a mutable tag such as `latest`;
3. a `CRED` contains plaintext-like material;
4. `BIND` lacks a mode;
5. a `TOOL` is declared but no policy permits it under deny-by-default mode;
6. `AUDIT` is absent;
7. a placement directive (`ISOLATION`, `IMAGE`, `SLICE`, `BACKEND`, `BIND`, `BROKER`, `PLUGIN`) couples the portable package to one runner, where a Run Manifest should be used instead.
