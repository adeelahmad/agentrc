# Current Implementation Mapping

**Input reviewed:** uploaded agentrc source archive containing Go and Python packages, examples, current README/HOWTO, parser code, lockfile code, OCI scaffolding, policy scaffolding, runner/driver experiments, and v2.5 proposal assets.

## What exists today

| Area | Current work observed | Specification treatment |
|---|---|---|
| Agentfile parser | Go parser in `packages/aio-core/pkg/agentfile`; Python parser in `packages/aio-core/src/aio` | Core Agentfile Profile |
| CLI | `agentrc compile`, `agentrc run tool`, `agentrc dmesg`, `agentrc status`, help/version | Implementation tooling, not the spec itself |
| Build/lock | `packages/aio-buildtime` compile and lock generation | Lockfile and OCI Package Profile |
| OCI package work | Manifest/config/layer/push/pull scaffolding | OCI Package Profile |
| Policy | Cedar-like Go stub and Python deny-by-default shim | Cedar Policy Profile target; implementation still needs real Cedar evaluator |
| Audit | Ring buffer, audit events, redaction tests | Security/Audit profile |
| Tools | `TOOL`, `TOOLSET`, tool registry, UTCI-style examples | Capability declaration and optional Tool Projection Profile |
| Functions | `FUNCTION` directive and function examples | Function capability declaration |
| Skills | `SKILL` directive and SKILL.md examples | Skill bundle declaration |
| Secrets | env/keyring/vault broker concepts | Credential declaration; runtime resolution by runner |
| Drivers | local/container/microVM/registry experiments | Runner implementations, not core spec |
| MicroVM/microsandbox | Candidate runner/backing substrate | Example runner profile only |
| Proposal assets | v2.5 branding and docs | Useful background; not normative |

## Main correction made in this draft

The uploaded work contains runtime/driver modules. That is useful implementation work, but it should not define the identity of agentrc.

This draft treats those modules as:

1. proof-of-concept runners;
2. conformance test harnesses;
3. developer tooling;
4. implementation profiles.

They are not the standard itself.

## Naming correction

Current README title:

```text
aio — Agent Isolation Orchestrator
```

Recommended public/spec title:

```text
agentrc Specification
```

or:

```text
agentrc — Agentfile and Agent Package Specification
```

Avoid “isolation orchestrator” in public positioning. It makes agentrc sound like a runtime and drags it into competition with Docker, gVisor, Firecracker, microsandbox, Kubernetes, and cloud runners.

## Directive compatibility notes

The current implementation recognizes a broad directive set:

```text
FROM SHELL CMD TOOL TOOLSET FUNCTION SKILL CRED URL SERVER MCP BIND
POLICY ALLOW DENY RATELIMIT ISOLATION BACKEND TRACE AGENT AUDIT
MEMORY OPTIMIZER BROKER TIMEOUT LIMIT SLICE IMAGE PLUGIN MOUNT
HEALTHCHECK
```

This draft keeps that directive set, but divides it into semantic layers:

| Directive class | Examples | Meaning |
|---|---|---|
| Agent identity/entrypoint | `AGENT`, `FROM`, `SHELL`, `CMD` | What the agent is and how it starts |
| Capability declaration | `TOOL`, `TOOLSET`, `FUNCTION`, `SKILL`, `SERVER`, `MCP` | What the agent requests |
| Boundary declaration | `URL`, `CRED`, `BIND`, `MOUNT`, `LIMIT`, `TIMEOUT`, `RATELIMIT` | What the agent may access and within what limits |
| Governance | `POLICY`, `ALLOW`, `DENY`, `AUDIT`, `TRACE` | How usage is reviewed and recorded |
| Runner hints/profile | `ISOLATION`, `IMAGE`, `SLICE`, `PLUGIN`, `BACKEND`, `BROKER`, `HEALTHCHECK` | What a runner may need to satisfy |
| Experimental | `MEMORY`, `OPTIMIZER` | Useful but should remain optional/profile-based |

## Known implementation/spec gaps

These should be turned into GitHub issues:

1. **Real Cedar evaluation** — current Go implementation contains a Cedar stub. The spec requires actual Cedar semantics for Cedar profile conformance.
2. **Fail-closed translation** — any policy or boundary the runner cannot understand must fail closed, not silently drop.
3. **CRED grammar mismatch** — current examples use `CRED db_password vault://...`, but one parser path treats two-argument `CRED` as URL/auth rather than name/ref. The spec standardizes named credential references.
4. **BIND arity mismatch** — one implementation path expects `BIND <host> <target> <mode>`, another accepts two arguments. The spec recommends three arguments and allows two-argument compatibility normalization.
5. **`ISOLATION` placement** — current work has `ISOLATION` in Agentfile and `--isolation` runtime flags. The draft keeps `ISOLATION` as compatibility but recommends moving placement into runner config later.
6. **OCI config alignment** — current OCI annotations use the legacy `io.agentio.*` namespace; the spec standardizes on `io.agentrc.*` to match the `AgentRC::` policy namespace and `vnd.agentrc.*` media types. Runners should migrate annotations and MAY read the legacy namespace for backward compatibility.
7. **Agent identity** — some examples infer identity from `CMD` or package context. Published packages should declare `AGENT` explicitly.
8. **Spec versioning** — current parser does not support a `SPEC` directive. The draft uses a parser-compatible comment: `# syntax=agentrc.agentfile/v0.1`.

## Recommended next repo changes

1. Rename public README heading to “agentrc Specification”.
2. Move runtime language under “compatible runners” or “reference tooling”.
3. Add `/docs/spec/` with this draft.
4. Add `/docs/profiles/` with profile docs.
5. Add a conformance matrix generated from actual tests.
6. Add issue labels: `spec`, `profile:runner`, `profile:cedar`, `profile:oci`, `compat`, `security`.
7. Split runtime placement from package declaration in a future run manifest.
8. Replace Cedar stub with a real Cedar evaluator before claiming Cedar profile conformance.

## Suggested repository structure

```text
/spec/
  README.md
  SPEC.md
  grammar/AGENTFILE.ebnf
  profiles/
  schemas/
  examples/

/implementation/
  cli/
  runners/
  conformance-tests/
```

If keeping the current monorepo, use:

```text
docs/spec/
docs/profiles/
docs/conformance/
packages/agentrc/                # CLI tool
packages/aio-core/           # parser/types/audit/policy interfaces
packages/aio-buildtime/      # compile/lock/OCI
packages/aio-runtime/        # compatible runner experiments
```
