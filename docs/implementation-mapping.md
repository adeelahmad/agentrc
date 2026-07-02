---
layout: doc
title: Implementation mapping
description: "How current implementation work maps onto the 0.1.0-draft.5 agentrc model: a frontend/compiler, a platform enforcement engine, OCI labels, substrates, and the /mnt projection."
permalink: /docs/implementation-mapping/
---
# Implementation mapping

agentrc is **spec-first**. The [Agentfile specification](/spec/) is the source of
truth; the implementation follows it. This page is an honest map from the current
implementation work to the v0.1 model — and a clear statement of where the code
**lags** the spec.

> **Read this as: "what already exists" vs. "what the spec says it should
> become."** Where the two disagree, the spec wins and the implementation is the
> thing that has to change.

## The v0.1 model in one diagram

```text
Agentfile ──build──►  OCI artifact (labels + layers)  ──run──►  Platform
(human recipe)        org.agentrc.* labels + /mnt resources      reads labels,
                                                                 grants/narrows/
   frontend / compiler      OCI labels & package          rejects, enforces
   (§9 translation)         (registry-portable)           via Cedar, on a
                                                          chosen substrate
```

The implementation splits along the same seams. Each existing component maps to
exactly one role in this pipeline.

## Component map

| Current implementation work | v0.1 role | What it produces / consumes |
|---|---|---|
| Agentfile parser | **Frontend / compiler** — translates the Agentfile into an OCI artifact | Reads the four agentrc keywords (`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY`) plus standard Dockerfile keywords; emits `org.agentrc.*` labels and `/mnt` layers per [spec §9](/spec/). |
| Cedar policy gate | **Platform enforcement engine** — compilation target for typed requests | Consumes the granted `org.agentrc.*` request labels (never the Agentfile), compiles them plus org rules into one Cedar `PolicySet`, evaluates deny-by-default. See [Enforcement profile](/profiles/security/). |
| Credential handling | **Deferred — platform-defined** | Secrets are out of scope for this draft: there is no agentrc secret schema. An agent that needs a credential leaves resolution entirely to the platform (Vault / broker / env / workload identity). |
| OCI image / package work | **OCI labels & package** | Builds the standard OCI artifact: layers carry `/mnt` resources, the image config carries the `org.agentrc.*` labels. See [OCI labels & package profile](/profiles/oci-package/). |
| microVM / runner drivers | **One substrate among many** — execution driver for `CMD` | A substrate executes `CMD`; it is selected at run time (`arc run --backend` / `--isolation`), **not** in the Agentfile. microVM is one substrate, not the product identity. |
| Tool patching / projection | **`/mnt` projection** | Projects `/mnt/tools`, `/mnt/skills`, `/mnt/mcp`, and populates `/mnt/proc`; loads the SOP from `/mnt/SOP`. See [projection profile](/profiles/tool-projection/). |

The implementation may keep its own internal names. None of those names define
the public identity of the project — the **Agentfile** and the `org.agentrc.*`
**labels** do.

## What each component owes the spec

**Frontend / compiler.** Two front doors, one artifact. The BuildKit frontend
(routed by `# syntax=agentrc.agentfile/v0.1`) and the native `arc build` MUST
emit **identical** OCI artifacts — same labels, same layers. The compiler MUST
embed `--cached` resources as layers, record `--runtime` resources as references,
emit both a digest and an `.origin` label for embedded MCP servers and skills,
auto-derive an attributed `network` egress label from hook / interrupt URLs, and
emit the SOP as a **pointer + digest** (`org.agentrc.sop=/mnt/SOP`,
`org.agentrc.sop.sha256=<digest>`) — never the full prompt text in a label.

**Platform enforcement engine.** Cedar is **platform-side only** and MUST NOT
appear in any Agentfile. The engine reads the `org.agentrc.*` labels (not the
Agentfile source), maps each request to a Cedar action/resource with the agent
identity as the principal, and preserves Cedar's properties: deny-by-default,
`forbid` over `permit` order-independently, and monotonic intersection across
`FROM`. The normative mapping lives in the [Enforcement profile](/profiles/security/).

**Credentials (deferred).** Secrets are out of scope for this draft. There is no
`SECRET`/`CRED` keyword and no agentrc secret schema; an agent that needs a
credential leaves resolution entirely to the platform (Vault / broker / env /
workload identity). A credential model may be specified in a later version.

**OCI labels & package.** A built agent is an ordinary OCI artifact: it pushes,
pulls, signs (Sigstore), and mirrors through any OCI-compatible registry. The
package layer carries `/mnt` resources; the config carries the labels.

**Substrate.** A substrate executes `CMD` and nothing more. Whether that is
`local`, `container`, or `microvm` is a deploy-time decision the platform makes —
the artifact is substrate-neutral.

**`/mnt` projection.** The projection layer mounts the embedded and fetched
resources under `/mnt`, makes tools self-describing (`--agentrc-schema` or a
sibling `<tool>.toolspec.json`), and exposes live policy, identity, budgets, and
the audit tail under `/mnt/proc`.

## Honest gap status (spec-first)

The implementation **lags the spec**, and we label that gap rather than hide it.

| Area | Spec status | Implementation status |
|---|---|---|
| Frontend / compiler (Agentfile → labels) | Normative ([§9](/spec/)) | In progress — keyword parsing and label translation are partial; the two build paths are not yet byte-identical. |
| Platform enforcement (Cedar) | Normative ([§11.2](/spec/)) | Prototype gate exists; the full request → Cedar mapping and `FROM` intersection are not complete. |
| Credentials | Deferred — out of scope this draft | No agentrc secret schema; credential resolution is platform-defined. |
| OCI labels & package | Normative | Label emission works; signing / provenance attestation is partial. |
| `/mnt` projection | Normative | Tool projection works; `/mnt/proc` runtime population is incomplete. |
| Substrates | Run-time choice | microVM and local drivers exist; others are adapters yet to be written. |

We do **not** advertise a [conformance profile](/docs/conformance/) we do not yet
pass. The reference frontend and the `arc` CLI are works in progress; the spec is
the contract they are being built against.

## Where to go next

- [Specification](/spec/) — the normative source of truth.
- [Enforcement (Cedar) profile](/profiles/security/) — platform-side enforcement.
- [OCI labels & package profile](/profiles/oci-package/) — the artifact format.
- [Conformance suite](/docs/conformance/) — the adversarial tests an implementation must pass.
