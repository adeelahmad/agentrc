---
sitemap: false
---
This note is superseded by [Implementation mapping](/docs/implementation-mapping/).

# Current Implementation Mapping

> **Internal implementation note (not a published spec page).** This maps the
> existing reference-implementation code onto the **0.1.0-draft.6** Agentfile model. It is
> spec-first: the specification leads, and the implementation lags it. Where the
> code still reflects the legacy directive family, that gap is called out
> honestly below rather than papered over. The normative source of truth is the
> [Agentfile specification](/spec/); this note only records how today's modules
> line up with it.

## The v0.1 model in one paragraph

An **Agentfile** is a Dockerfile-shaped recipe with exactly **four** new
keywords — `IDENTITY`, `CAPABILITY`, `SOP`, `POLICY` — over standard Dockerfile
keywords (`FROM`, `CMD`, `COPY`, `ADD`, `HEALTHCHECK`, `LABEL`, `ENV`, `ARG`,
`WORKDIR`, `USER`, `EXPOSE`, `RUN`). Tools, skills, and MCP servers are files
placed under **`/mnt`** with `COPY` (local) or `ADD --remote` (remote); the
destination path under `/mnt` determines the resource type. The **compiler /
frontend** translates authored intent into `ai.agentrc.*` OCI labels + layers.
The **platform** reads those labels — never the Agentfile — and grants, narrows,
or rejects each typed `POLICY` request, enforcing the result with **Cedar**
(platform-side only). A **substrate** is the run-time driver that actually
executes `CMD`. Secrets are **deferred** — this draft defines no secret
keyword/schema; credential resolution is platform-defined and out of scope.

## Module → v0.1 role mapping

This is the core of the reframe: every existing module is recast as one piece of
the build → labels → platform → substrate pipeline. None of these modules *is*
the standard — they are a reference implementation of it.

| Current module / work | v0.1 role | Reads / emits |
|---|---|---|
| Agentfile parser (Go in `pkg/agentfile`; Python package) | **Frontend / compiler** — Agentfile → `ai.agentrc.*` labels + layers | parses the four agentrc keywords + standard Dockerfile keywords; emits labels per spec §9 |
| `compile` / `build` tooling | **`agentrc build` / BuildKit frontend** — two front doors, identical OCI artifacts | emits OCI artifact with `ai.agentrc.*` labels; `--policy-mode inline\|digest` |
| OCI manifest/config/layer push/pull scaffolding | **OCI labels & package** | layers carry `/mnt` resources; image config carries `ai.agentrc.*` labels |
| Cedar-like policy stub / deny-by-default shim | **Platform enforcement engine** (Cedar, platform-side) | compiles granted requests + org rules into one Cedar `PolicySet`; deny-by-default |
| Credential-resolution concepts (env / keyring / vault) | **Deferred** — no secret keyword/schema in this draft | credential resolution is platform-defined and out of scope for v0.1 |
| Audit ring buffer / redaction | **Platform audit / telemetry sink** | fed by `agent.hooks.*` / `agent.telemetry_sink`; records grant / narrow / reject decisions |
| Tool registry / projection | **`/mnt` projection** | projects `/mnt/tools`, `/mnt/skills`, `/mnt/mcp`; populates `/mnt/proc` |
| local / container / microVM / registry drivers | **One substrate among many** | chosen at run time via `arc run --isolation \| --substrate`, never in the Agentfile |
| microVM / microsandbox experiments | **One substrate driver** | substrate-neutral; the platform picks it at run time |
| v2.5 proposal / branding assets | Background only | not normative |

## Authoring surface: four keywords, no legacy directives

The implementation historically recognized a broad directive family. In 0.1.0-draft.6 that
family is **removed**. The authoring surface is now:

```text
New agentrc keywords:   IDENTITY  CAPABILITY  SOP  POLICY
Standard Dockerfile:    FROM  CMD  COPY  ADD  HEALTHCHECK  LABEL
                        ENV  ARG  WORKDIR  USER  EXPOSE  RUN
```

The legacy directives the old code parsed —
`AGENT`/`TOOL`/`TOOLSET`/`FUNCTION`/`SKILL`/`SERVER`/`MCP`/`URL`/`CRED`/`BIND`/`MOUNT`/`PLUGIN`/`ALLOW`/`DENY`/`RATELIMIT`/`TIMEOUT`/`LIMIT`/`SLICE`/`IMAGE`/`ISOLATION`/`BACKEND`/`BROKER`/`TRACE`/`MEMORY`/`OPTIMIZER`/`SHELL` —
plus the inline Cedar `POLICY … END` block and the `SOP name … END` block, are
**stale** and must be dropped from the parser. Their concerns map onto the v0.1
surface as follows:

| Legacy directive(s) | v0.1 replacement |
|---|---|
| `AGENT` / identity-from-`CMD` | `IDENTITY name=… version=… author=…` → `ai.agentrc.identity.*` |
| `TOOL` / `TOOLSET` / `FUNCTION` | `COPY` / `ADD --remote … /mnt/tools/<name>` → `ai.agentrc.tool.<name>` |
| `SKILL` | `COPY` / `ADD --remote … /mnt/skills/<name>` → `ai.agentrc.skill.<name>` |
| `SERVER` / `MCP` | `COPY` / `ADD --remote … /mnt/mcp/<name>` → `ai.agentrc.mcp.<name>` (+ `.origin`) |
| `URL` (egress) | `POLICY network dns:<host>:<port>` → `ai.agentrc.network.dns.<host>` |
| `CRED` / `BROKER` | Deferred — no secret keyword/schema in v0.1; credential resolution is platform-defined |
| `BIND` / `MOUNT` | files under `/mnt` via `COPY` / `ADD --remote` |
| `RATELIMIT` / `TIMEOUT` / `LIMIT` | typed `POLICY agent.*` requests (e.g. `agent.tool_timeout`) |
| `MEMORY` / `OPTIMIZER` | `POLICY agent.memory.short` / `agent.memory.long` / `agent.context.type` |
| `ISOLATION` / `IMAGE` / `SLICE` / `BACKEND` | run-time substrate choice (`arc run --isolation \| --substrate`), not authored |
| `ALLOW` / `DENY` / inline Cedar block | **platform-side Cedar only** — never authored in the Agentfile |
| `TRACE` | `POLICY agent.hooks.*` / `agent.telemetry_sink` |
| `SHELL` | `POLICY substrate.init` / `substrate.ptty` |

## Known implementation / spec gaps (spec-first; turn into issues)

The implementation lags the spec. These are honest gaps, not advertised
capabilities:

1. **Frontend keyword set.** The parser must recognize only the four agentrc
   keywords + standard Dockerfile keywords, and **reject / ignore** the legacy
   directive family. SOP heredoc (`SOP <<EOF … EOF`) must be captured verbatim;
   the inline and file-backed (`COPY ./sop.md /mnt/SOP`) forms must all land at
   `/mnt/SOP`.
2. **`ADD --remote` delivery flags.** Implement `--cached` (default) /
   `--runtime` and `--fail-if-unavailable` (default) / `--warn-if-unavailable`,
   plus standard `--chmod` / `--chown`. The `/mnt` destination determines the
   resource type.
3. **Label emission.** Emit `ai.agentrc.*` labels exactly per spec §9 tables —
   `POLICY <key> <value>` → `ai.agentrc.<key>=<value>`; embedded MCP / skills
   emit **both** digest and `.origin`; hook / interrupt URLs auto-derive an
   **explicit, attributed** `network.dns.*` egress label (with `.source`).
4. **Identical artifacts.** The BuildKit frontend and `arc build` must produce
   byte-identical OCI artifacts and labels.
5. **Real Cedar enforcement.** Replace the Cedar stub with a real Cedar
   evaluator before claiming the enforcement profile. The platform must compile
   each granted request + org rules into one `PolicySet` and preserve
   deny-by-default, `forbid` over `permit` (order-independent), and monotonic
   intersection across `FROM`. Cedar MUST NOT appear in any Agentfile.
6. **Fail-closed translation.** Any required boundary the platform cannot
   enforce, or any policy it cannot evaluate, must **fail closed**, never
   silently drop.
7. **Secrets are deferred.** This draft defines no secret keyword/schema;
   credential resolution is platform-defined and out of scope for v0.1.
8. **OCI namespace.** Standardize on the `ai.agentrc.*` label namespace and
   `application/vnd.agentrc.*` media types. The legacy `io.agentio.*` /
   `io.agentrc.*` namespaces are **not** used in 0.1.0-draft.6.
9. **Spec-version comment.** Use the parser-compatible syntax line
   `# syntax=agentrc.agentfile/v0.1` (replacing the old v0.1 syntax line);
   spec-version mentions read "0.1.0-draft.6 (Working Draft)", dated 2026-06-30.
10. **Substrate placement.** Substrate / isolation is a **run-time** choice
    (`arc run --isolation local|container|microvm` / `--substrate <driver>`),
    never an Agentfile directive. Remove any `ISOLATION` parsing.

## Naming

- **agentrc** is the specification (Agentfile + agent package + profiles).
- The reference implementation is *tooling*, not the standard. Avoid
  "isolation orchestrator" positioning — it makes agentrc sound like a runtime
  and drags it into competition with Docker, gVisor, Firecracker, microsandbox,
  Kubernetes, and cloud runners. agentrc is the **build-time + labels contract**;
  any of those can be the substrate.

## Recommended next repo changes

1. Rewrite the parser to the four-keyword + Dockerfile-keyword surface; drop the
   legacy directive family and the inline Cedar / `SOP … END` blocks.
2. Make the frontend and `arc build` emit identical `ai.agentrc.*` labels +
   layers per spec §9.
3. Implement the `ADD --remote` delivery-flag contract and `/mnt` projection.
4. Replace the Cedar stub with a real platform-side Cedar evaluator; wire
   deny-by-default, `forbid` > `permit`, and monotonic `FROM` intersection.
5. Keep secrets deferred — define no secret keyword/schema in v0.1; credential
   resolution stays platform-defined and out of scope.
6. Migrate OCI annotations to `ai.agentrc.*` + `vnd.agentrc.*`; drop legacy
   namespaces.
7. Keep substrate/isolation in `arc run`, not the Agentfile.
8. Generate a conformance matrix from real tests; do not advertise a profile the
   implementation does not pass (see the [conformance suite](/docs/conformance/)).

## Where this maps in the spec

- Authoring surface and label tables: [/spec/](/spec/)
- Compiler behaviour: [/profiles/core/](/profiles/core/)
- Platform Cedar enforcement: [/profiles/security/](/profiles/security/)
- OCI labels & package: [/profiles/oci-package/](/profiles/oci-package/)
- `/mnt` projection: [/profiles/tool-projection/](/profiles/tool-projection/)
- Platform conformance: [/profiles/runner-conformance/](/profiles/runner-conformance/)
