# AGENTRC WORK ORDER v2 — SPRINT 2 slice (intake input)

Staged by the supervisor as the intake document for Sprint 2 ("Features",
T15–T27). Sprint 1 is complete + live-verified + owner-signed-off ("go").

Repo: /Users/adeelahmad/work/agentrc (`master`). Live: https://agentrc.ai.

## Supervisor-carried resolutions & flags (surface at planning gate)
- **T8 lockfile decision (owner said "go" without picking):** default to **Option A** —
  add an informative spec §9 subsection "Reproducible builds / `agentrc.lock`" documenting
  only what the tooling actually does (`arc lock` writes `agentrc.lock`; `arc build` does
  NOT consume it today — see docs/agents/sprint1/lockfile-report.md), marked
  "Status: informative in this draft; format TODO", keep the homepage slogan. Owner may
  override to Option B (cut slogan; lockfile on /cli/ only). Lands in T20.
- **OPEN QUESTION (blocking for provenance):** `agentrc-draft6-addendum-2026-07-02.md`
  (the Phase-4 "full text" source) is NOT present in the repo. T17–T19 inline specs below
  are detailed enough to author the spec sections, but "commit the addendum for provenance"
  (T-Phase-4 header) cannot be done without the file. Owner: provide it, or accept authoring
  from this work order's inline detail (no separate addendum committed).
- **ENV RISK:** local disk is near-full and Go builds pulling the buildkit dep tree have hit
  ENOSPC. T21–T24 (Go backend translators) real-TDD execution may be blocked locally; may
  need to lean on CI for `go build/vet/test`. Surface feasibility at the planning gate.
- **ENV:** the agentic-agile plugin did not fully load this session (adapted supervisor mode;
  manual gate scripts). A Claude Code restart would restore hook-enforced gates.

## §0 GLOBAL INVARIANTS (never violate)
1. Version `0.1.0-draft.5` stays through Phase 3; bump to `0.1.0-draft.6` **only in T20**
   (sitewide, one commit). Syntax line stays `# syntax=agentrc.agentfile/v0.1` throughout.
2. **Terminology split:** `--backend` is the CLI flag (`local|bedrock|kubernetes`); the old
   `--substrate <driver>` flag + "driver" wording are renamed in T21. `substrate.*` (POLICY
   namespace) and "substrate-neutral" (concept) are NEVER renamed.
3. Exactly four keywords (IDENTITY, CAPABILITY, SOP, POLICY). No new keywords ever. Sprint 2
   adds EXACTLY three POLICY namespace entries: T17 `substrate.<platform>.*`, T18
   `agent.auth.*`, T19 `substrate.runtime.language`. No others.
4. POLICY lines are requests; Cedar platform-side only; secrets deferred. No inline Cedar or
   secret keyword. (agent.auth.* is generic authZ config, fail-closed; NOT a secret.)
5. Open decisions (spec §14.2) stay open; follow each task's stated default.
6. ONE hello: byte-identical wherever rendered inline and identical to examples/Agentfile.minimal.
7. Every example file must pass `arc lint` (`go run ./cmd/agentrc lint <file>` from repo root —
   note: CLI is at repo-root ./cmd/agentrc, NOT tooling/).
8. Backend docs positioning line, VERBATIM: "Reference translators — a proof of concept until
   platforms read `org.agentrc.*` labels natively. Not production runners."
9. Locate every target by grep first. Rebuild site + pass §V before commit. (No local Jekyll in
   this env → build-dependent checks delegated to CI + live, as in Sprint 1.)

## SPRINT 2 TASKS

### Phase 3 — examples expansion
**T15 — `examples/Agentfile.hooked`** [P2] Demonstrates `POLICY agent.hooks.on_tool_call <https
endpoint>` + one `pre`/`post` hook; comments explain the platform auto-derives the hook
endpoint's egress grant and records it with `.source` attribution (requested vs derived), per
/docs/security/. Include one explicit `POLICY network dns:…` line for contrast.
**T16 — `examples/Agentfile.delegator`** [P2] Demonstrates `POLICY agent.sub_agents true`,
`agent.sub_agents.max <n>`, `agent.sub_agent_timeout <dur>`; comment: sub-agent grants are the
platform's call. Spec-defined keys only.
Both files: `# syntax=` line, `FROM`, IDENTITY/CAPABILITY/SOP/CMD, the Sprint-1 healthcheck
pattern (`/mnt/tools/file_read --agentrc-schema`), header comment; add to examples index in the
existing card style. Optional: one index line showing `arc build --policy-mode digest` (no
default-mode prose — §14.2 #1). Verify: `arc lint` passes both; links resolve; §V clean.

### Phase 4 — spec draft.6 (authored from inline detail below; commit addendum iff owner provides)
**T17 — §8.5 `substrate.<platform>.*`** [P1] Tokens `aws|gcp|azure|kubernetes|local`; unknown
tokens MUST parse; foreign-platform keys ignored, never errors (linter MAY warn); labels
`org.agentrc.substrate.<platform>.<key>=<value>`; platform-scoped beats generic on that platform
only; tightening-only across `FROM` per namespace. AWS registry: `roleArn`, `networkMode`,
`securityGroup` (rep.), `subnet` (rep.), `protocol`, `maxLifetime`, `deployment.mode`
(`container` default|`code`), `code.s3.uri`.
**T18 — §8.6 `agent.auth.*` (generic, fail closed)** [P1] `agent.auth.mode` (`platform`
default|`jwt`|`none`), `agent.auth.jwt.discovery_url`, `.allowed_audience` (rep.),
`.allowed_client` (rep.). A platform that can't enforce a requested `jwt` authorizer MUST NOT
expose the invocation endpoint. No secrets.
**T19 — §8.7 `substrate.runtime.language`** [P1] `<language>:<version>`. Optional; container-mode
MAY ignore (base image authoritative); code-mode requires it or resolvable inference, else fail
closed.
**T20 — Supporting edits + version bump + T8 lands** [P1] §14.2 open decision #6 (promotion
candidates: `protocol`, `maxLifetime`); /docs/agentfile/ platform-scoped paragraph + JWT example;
/profiles/core/ accepts unknown `substrate.<token>.*`; Agentfile.code-reviewer gains a commented
`substrate.aws.*` + `agent.auth.*` block; CHANGELOG draft.6 entry; **sitewide bump draft.5 →
draft.6**; implement the T8 A/B choice (default A).

### Phase 5 — CLI: `arc run --backend`
**T21 — Flag rename & surface** [P0* demo] `--substrate` → `--backend` everywhere (CLI code, help,
/cli/, quickstart step 5, tooling README):
  `arc run <ref> --backend local|bedrock|kubernetes`
   local: `[--isolation microvm|container]` (microsandbox MVP exists);
   bedrock: `[--region] [--profile] [--dry-run]`; kubernetes: `[--kubeconfig] [--namespace] [--dry-run]`.
  `--dry-run` prints the translated config and exits. Record in CLI docs: GCP dropped (Agent
  Runtime Python-only managed; GKE via kubernetes backend); Docker Compose dropped (no network.*
  egress enforcement without a bespoke sidecar). Verify: `rg -- '--substrate' -l` → 0;
  `rg 'POLICY substrate\.' spec/` intact.
**T22 — Backend `local`** [P1] Wire the existing microsandbox VMM MVP under `--backend local`
(default). Plumbing only + §0.8 positioning line.
**T23 — Backend `bedrock` (labels → CreateAgentRuntime, 13/13)** [P1] Map org.agentrc.* labels +
image config to Bedrock AgentRuntime fields (agentRuntimeName/description ← IDENTITY;
containerUri ← OCI ref; roleArn ← substrate.aws.roleArn; networkMode ← substrate.aws.networkMode;
securityGroups/subnets ← substrate.aws.securityGroup/subnet; serverProtocol ← substrate.aws.protocol;
env ← image Env; customJWTAuthorizer ← agent.auth.jwt.*; idleRuntimeSessionTimeout ←
agent.idle_timeout; maxLifetime ← substrate.aws.maxLifetime; codeConfiguration ← deployment.mode=code
+ code.s3.uri + substrate.runtime.language). Fail closed: missing roleArn; unenforceable
agent.auth.mode=jwt; code mode without resolvable language. `--dry-run` emits the translated config.
Re-update quickstart step-5 wording (run now real).
**T24 — Backend `kubernetes`** [P1] Emit (dry-run) or apply: Deployment (resources from
substrate.runtime.*, env from image config), Service, deny-by-default NetworkPolicy from `POLICY
network dns:*`, ServiceAccount from substrate.kubernetes.serviceAccount, MCP servers from
/mnt/mcp/* as sidecars. One emission format (manifests OR Helm), not both.
**T25 — CLI docs table** [P1] `run` → `implemented (local, bedrock, kubernetes — reference
translators)`; `sign`/`verify` stay `planned`; §0.8 positioning line above the table.

### Phase 6 — demo & release
**T26 — One agent, three backends** [P2] `arc build -t ghcr.io/agentrc/code-reviewer:1.0 .`;
`arc run ... --backend local --isolation microvm`; `... --backend bedrock --dry-run`;
`... --backend kubernetes --dry-run`. Narrative VERBATIM: "Same artifact, same labels, three
substrates. The translators are the proof of concept; the labels are the standard."
**T27 — RELEASE Sprint 2 + live verification** [P0] Full §V incl checks 8–9; deploy; live re-verify
each Sprint-2 task (version = draft.6 sitewide; --backend in live CLI docs; new examples lint +
render; addendum committed iff provided); record per-task pass/fail; owner sign-off; retrospective.

## §V additions for Sprint 2
- 8 terminology split: `rg -- '--substrate' -l | wc -l` → 0 (CLI/docs); `rg 'POLICY substrate\.'
  spec/` intact.
- 9 backends: `arc run --help | grep -- '--backend'`; `arc run <ref> --backend bedrock --dry-run
  | python3 -m json.tool`; `arc run <ref> --backend kubernetes --dry-run | kubeconform -` (or
  yaml-parse if kubeconform absent).
- 3 version coherence: exactly one `draft.N` value = `draft.6` AFTER T20 (draft.5 before).

## §H owner/infra (out of pipeline): CDN purge for /spec/; agentrc.io DNS 301; name-collision.
