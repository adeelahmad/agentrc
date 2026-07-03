---
layout: doc
title: Conformance
description: "Profile-based conformance for agentrc 0.1.0-draft.6, with an adversarial, fail-closed test suite that makes every profile claim verifiable."
permalink: /docs/conformance/
---
# Conformance

agentrc conformance is **profile-based**. An implementation — a compiler /
BuildKit frontend, a registry, a platform / runner, or a workflow engine —
states exactly which profiles it supports, and a profile claim is only valid if
the implementation passes that profile's suite, including its **adversarial**
cases.

> A specification without an executable test suite is prose. The conformance
> suite is what makes a profile claim verifiable: an implementation either
> passes a profile's suite or it does not. The adversarial cases prove the
> implementation does the **safe** thing under bad input and fails **closed** —
> not just the happy path.

## Profiles

Each profile corresponds to one of the [conformance profile pages](/profiles/),
which carry the normative requirements. The suite below references them by short
name.

| Profile | Covers | Profile page |
|---|---|---|
| `agentrc/agentfile/v0.1` | **Build conformance** — the compiler / frontend: parse the Dockerfile-shaped Agentfile (four agentrc keywords + standard Dockerfile keywords) and compile Agentfile → `ai.agentrc.*` labels + layers. | [Core](/profiles/core/) |
| `agentrc/enforcement-cedar/v0.1` | Platform-side Cedar enforcement of granted requests (deny-by-default, `forbid` > `permit`, monotonic `FROM`). | [Enforcement (Cedar)](/profiles/security/) |
| `agentrc/oci-labels/v0.1` | The `ai.agentrc.*` label namespace, layers, media types, and package shape. | [OCI Labels &amp; Package](/profiles/oci-package/) |
| `agentrc/projection/v0.1` | The `/mnt` projection (`tools/`, `skills/`, `mcp/`, `proc/`, `SOP`) and the tool invocation contract. | [`/mnt` Projection](/profiles/tool-projection/) |
| `agentrc/platform/v0.1` | **Platform conformance** — runtime behaviour: read labels, grant / narrow / reject, project requests, fetch `--runtime` resources, substitute via `.origin`, enforce, boot `CMD`. | [Platform Conformance](/profiles/runner-conformance/) |
| `agentrc/workflow/v0.1` | The **deferred**, non-normative multi-agent workflow companion (out of scope for 0.1.0-draft.6). | Workflow orchestration is parked for a future draft. |

## Why profiles?

A compiler should not need to implement Cedar evaluation. A registry should not
need to boot a substrate. A platform should not need to become a workflow
engine. Splitting conformance into profiles keeps each surface implementable and
lets an implementation advertise precisely what it does.

The boundary that matters most: **build conformance** (`agentrc/agentfile/v0.1`)
is the compiler — it reads the Agentfile and emits `ai.agentrc.*` labels —
while **platform conformance** (`agentrc/platform/v0.1`) reads those labels and
grants / projects / enforces them. Every profile other than build conformance
consumes **labels**, never the Agentfile source. The suite enforces that
separation.

## Conformance suite (0.1.0-draft.6)

The suite is intentionally as important as the spec text. The adversarial table
is the core of it: each case has a single correct outcome, and a missed case
means the implementation is **not conformant** to that profile no matter how
many positive cases it passes.

### Positive cases

| ID | Profile | Given | Expect |
|---|---|---|---|
| `agentfile-parse-minimal` | agentfile | The minimal valid Agentfile (`# syntax=agentrc.agentfile/v0.1`, `FROM`, `IDENTITY`, `SOP`, `CMD`) | Parses; instruction order preserved; `ai.agentrc.identity.*` and `ai.agentrc.sop` labels emitted |
| `policy-to-label` | agentfile | `POLICY model.name claude-opus-4` | Emits `ai.agentrc.model.name=claude-opus-4` (short form prefixed with `ai.agentrc.`) |
| `add-remote-runtime-recorded` | agentfile | `ADD --remote --runtime <url> /mnt/mcp/github` | No embedded layer; emits `ai.agentrc.mcp.github=runtime:<url>` |
| `oci-roundtrip` | oci-labels | A built artifact | Push, pull by digest, and inspect reproduce identical layers, config, and `ai.agentrc.*` labels |
| `egress-granted-allowed` | enforcement | A granted `ai.agentrc.network.dns.api.github.com=443` request | `Action::"NetworkEgress"` to `Host::"api.github.com:443"` is allowed |
| `mnt-projection-readable` | projection | A built artifact with tools and an SOP | `/mnt/tools/*` are executable; `/mnt/SOP` is a readable file; `/mnt/proc` is populated at boot |

### Adversarial / fail-closed cases

These are the cases that catch real implementation gaps. Each has a single
correct outcome.

| ID | Profile | Given | MUST |
|---|---|---|---|
| `build-labels-identical` | agentfile | The **same** Agentfile built via the BuildKit frontend and via `arc build` | Produce an **identical** OCI artifact and identical `ai.agentrc.*` labels — two front doors, one compiler |
| `sop-pointer-not-inlined` | agentfile | A large `SOP` (heredoc or file-backed) | Embed it at `/mnt/SOP` and emit only `ai.agentrc.sop=/mnt/SOP` + `ai.agentrc.sop.sha256=<digest>` — **never** the full text in a label |
| `unknown-request-denied` | enforcement | A request with no matching grant in the compiled `PolicySet` | **Deny** (deny-by-default); an unrecognised or un-granted request is never silently honoured |
| `forbid-overrides-permit` | enforcement | An org `forbid` against an agent request that a `permit` would otherwise allow | **Deny**, **order-independently** — `forbid` always wins |
| `child-widens-parent-fails` | enforcement | A child agent (`FROM another-agent`) whose requests exceed the parent's ceiling | Effective authorization is the **intersection** of ceilings; the widening request is rejected (a parent `forbid` is un-loosenable) |
| `auto-egress-not-implicit` | enforcement | A `POLICY agent.hooks.pre <url>` that auto-derives `ai.agentrc.network.dns.<host>` with `…source=auto:agent.hooks.pre` | The platform MUST still **grant** the derived egress; if ungranted → **deny**. Auto-derivation is convenience, never an implicit grant |
| `runtime-fetch-failclosed` | platform | An `ADD --remote --runtime --fail-if-unavailable` resource that is unreachable at boot | **Refuse to boot** (fail closed); never start the agent with the required resource missing |
| `embedded-override-origin` | platform | An embedded MCP/skill whose `ai.agentrc.<kind>.<name>.origin` is rewritten to an internal mirror at deploy time | Fetch from the **mirror** and run, **without rebuilding** the artifact |

A platform or compiler that claims a profile but fails any of that profile's
adversarial cases is **not conformant** to that profile, regardless of how many
positive cases it passes.

These cases are the executable counterpart of the spec's normative properties:
the request → Cedar mapping and the deny-by-default / `forbid` > `permit` /
monotonic-`FROM` guarantees are defined in
[§11.2 of the specification](/spec/) and the
[Enforcement (Cedar) profile](/profiles/security/); the label and package shape
in the [OCI Labels &amp; Package profile](/profiles/oci-package/); and the full
platform decision flow in the
[Platform Conformance profile](/profiles/runner-conformance/).

## Honest reference-implementation status

agentrc is the **specification**; the reference work in this repository is an
implementation and test harness, **not** the definition. The project is
spec-first, so the implementation is expected to **lag** the spec — and that gap
is labeled honestly rather than implied away.

As of this Working Draft, a reference implementation of both build paths (the
BuildKit frontend and the native `agentrc` / `arc` CLI — see
[`tooling/`](https://github.com/adeelahmad/agentrc/tree/master/tooling))
covers Agentfile parsing and `ai.agentrc.*` label emission
(`agentrc/agentfile/v0.1`) and produces a standard OCI artifact
(`agentrc/oci-labels/v0.1`). A complete Cedar evaluator and the full adversarial
fail-closed suite for `agentrc/enforcement-cedar/v0.1` and `agentrc/platform/v0.1`
are **not** yet in place, and those profiles are not claimed until they pass.

Implementations **MUST NOT** advertise a profile they do not pass.
