---
layout: doc
title: Conformance
description: "Conformance profile summary."
permalink: /docs/conformance/
---
# Conformance

agentrc conformance is profile-based.

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

## Conformance suite (v0.1 outline)

A specification without an executable test suite is prose. The conformance suite is what makes a profile claim verifiable: an implementation either passes the suite for a profile or it does not. The suite is intentionally as important as the spec text, and it must include **adversarial** cases — the suite proves a runner does the safe thing under bad input, not just the happy path.

### Positive cases

| ID | Profile | Given | Expect |
|---|---|---|---|
| `core-parse-minimal` | Core | The minimal valid Agentfile | Parses; directive order preserved; structured tree emitted |
| `core-policy-block` | Core | A `POLICY ... END` block | Inner lines captured verbatim, not parsed as directives |
| `oci-roundtrip` | OCI Package | A built package | Push, pull by digest, and inspect reproduce identical content |
| `cedar-permit` | Cedar | A request matching a `permit` | Allowed and (if `AUDIT` requires) recorded |

### Adversarial / fail-closed cases

These are the cases that catch real implementation gaps. Each one has a single correct outcome.

| ID | Profile | Given | MUST |
|---|---|---|---|
| `policy-unparseable-denies` | Cedar | Policy source that does not parse | **Deny** every request (fail closed), never allow |
| `policy-eval-error-denies` | Cedar | Policy that errors during evaluation | **Deny** the request |
| `forbid-overrides-permit` | Cedar | A request matched by both a `permit` and a `forbid` | **Deny** (deny wins) |
| `unknown-required-directive` | Core / Runner | An unknown directive marked required | **Reject** the package, do not silently ignore |
| `cred-value-redacted` | Security | A `CRED` resolves to a secret value | Value **redacted** from logs, audit, lockfile, and package metadata |
| `cred-plaintext-rejected` | Security | A package containing a plaintext secret | **Reject** / fail validation |
| `child-widens-forbid-fails` | Inheritance | A child package that removes or widens an inherited `forbid`/ceiling | **Fail** the build |
| `audit-required-unsupported-fails` | Runner | `AUDIT` required but the runner cannot emit audit events | **Fail closed**, do not run |
| `boundary-unsupported-fails` | Runner | A required boundary the runner cannot enforce | **Fail closed**, do not silently weaken |

A runner that claims a profile but fails any of that profile's adversarial cases is **not conformant** to that profile, regardless of how many positive cases it passes.

## Honest conformance status of the reference implementation

agentrc is the specification; the reference implementation (the `aio-*` packages in this repository) is an implementation and test harness, **not** the definition. Spec-first work means the spec leads the implementation, so the implementation is expected to lag — and that gap must be labeled honestly rather than implied away.

As of this draft, the reference implementation should be described as passing only the profiles it actually passes today (Agentfile parsing and OCI packaging), and **not** yet claiming the Cedar profile or full fail-closed Security/Runner conformance until a real Cedar evaluator and the adversarial cases above are in place. Implementations MUST NOT advertise a profile they do not pass.
