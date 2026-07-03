---
type: stories
sprint: 2
---
# Sprint 2 Stories — "Features" (T15–T27)

Stage-2 story breakdown authored from `work-order.md` (T15–T27 + §0 invariants +
§V suite + the supervisor-carried resolutions/flags), `intake.md` (five-part
Intent + open questions + the §8.5/8.6 collision note), `standards.md` (Go stack +
gate matrix, ENOSPC→CI note), and `memory.md` (M-001..M-006). Every edit target
below was grep-located against the live repo, not assumed (§0.9, M-003).

## Sprint goal

Turn agentrc from "consistent + honest site" (Sprint 1) into "demonstrably
substrate-neutral": ship two lint-clean examples (`Agentfile.hooked`,
`Agentfile.delegator`), a spec **draft.6** that adds **exactly three** new POLICY
namespaces (`substrate.<platform>.*`, `agent.auth.*`, `substrate.runtime.language`)
and nothing else, the CLI **`--substrate` → `--backend`** rename with three
fail-closed reference translators (`local`, `bedrock`, `kubernetes`), a
one-agent/three-backend demo, and a full release with live re-verification — all
without adding a keyword, renaming the `substrate.*` POLICY namespace or the
"substrate-neutral" concept, or bumping the version before T20.

## Sprint demo

1. Run `scripts/verify-sprint2.sh` (all §V checks incl. new 8/9) green: `--substrate`
   gone from CLI/docs (→0) while `POLICY substrate.` in spec stays intact; exactly
   one `draft.N` value = `draft.6` after T20; both new examples pass
   `go run ./cmd/agentrc lint`; `arc run --help` shows `--backend`.
2. Show the three new spec sections §8.7/§8.8/§8.9 defining exactly the three
   namespaces with their keys and fail-closed semantics, keyword count unchanged.
3. Run the one-agent/three-backend demo (T26): `arc build`, then `--backend local
   --isolation microvm`, `--backend bedrock --dry-run` (pipes through
   `python3 -m json.tool`), `--backend kubernetes --dry-run` (yaml-parses), with the
   verbatim narrative "Same artifact, same labels, three substrates. …".
4. Show the merged release PR to `master`, then re-run the live per-task checks
   (`scripts/verify-sprint2-live.sh`) against `https://agentrc.ai` all green (T27).

## Definition of Done

- [ ] §V checks 1–11 incl. new checks 8 (terminology split) and 9 (backends) pass
      locally where feasible / in CI otherwise (via `scripts/verify-sprint2.sh`).
- [ ] Every task T15–T27's own stated Verify passes (local or CI per gate matrix).
- [ ] Version bumped to `0.1.0-draft.6` in exactly one T20 commit; exactly one
      `draft.N` value tree-wide after T20, `draft.5` before (§0.1, §V.3).
- [ ] Exactly three new POLICY namespaces added; no new keyword; `substrate.*`
      namespace and "substrate-neutral" concept NOT renamed (§0.2, §0.3).
- [ ] All three backend translators are fail-closed (no `roleArn`; unenforceable
      `jwt`; code-mode without resolvable language; k8s deny-by-default NetworkPolicy)
      and each has a failing-first (RED) test proving it (§0-standards, work-order T23/T24).
- [ ] §0.8 positioning line VERBATIM on every backend doc surface; T26 narrative verbatim.
- [ ] Canonical `hello` unchanged, byte-identical to `examples/Agentfile.minimal` (§0.6, M-001).
- [ ] Release lands on `master` via PR; no direct push. T27 live verification recorded
      pass/fail per task ID, all green; owner sign-off; retrospective written.

## Out of scope

- Anything beyond T27.
- Renaming the `substrate.*` POLICY namespace or the "substrate-neutral" concept
  (the rename is CLI-flag-only: `--substrate` → `--backend`).
- Adding a fourth new POLICY namespace or any new keyword.
- Production-grade runners — the three translators are proof-of-concept only.
- Secrets / inline Cedar / a secret keyword (`agent.auth.*` is generic fail-closed
  authZ config, NOT a secret).
- **GCP and Docker Compose backends** — explicitly dropped in T21 (GCP: Agent Runtime
  is Python-only managed, GKE covered via the `kubernetes` backend; Compose: no
  `network.*` egress enforcement without a bespoke sidecar). Recorded in CLI docs.
- §H owner/infra items (CDN purge for `/spec/`, `agentrc.io` DNS 301, name-collision).

## User stories

### S2-01 — Examples expansion (T15, T16)

- **What's wanted:** Two new lint-clean examples using only spec-defined POLICY keys,
  plus their examples-index cards in the existing card style. T15
  `examples/Agentfile.hooked`: demonstrate `POLICY agent.hooks.on_tool_call <https
  endpoint>` + one `pre`/`post` hook; comments explain the platform auto-derives the
  hook endpoint's egress grant and records it with `.source` attribution (requested
  vs derived) per §8.5 / /docs/security/; include one explicit `POLICY network dns:…`
  line for contrast. T16 `examples/Agentfile.delegator`: demonstrate
  `POLICY agent.sub_agents true`, `agent.sub_agents.max <n>`,
  `agent.sub_agent_timeout <dur>`; comment that sub-agent grants are the platform's
  call. Both files: `# syntax=agentrc.agentfile/v0.1` line, `FROM`,
  IDENTITY/CAPABILITY/SOP/CMD, the Sprint-1 healthcheck pattern
  (`/mnt/tools/file_read --agentrc-schema`), header comment; add to examples index in
  the existing card + prose style. Optional: one index line showing
  `arc build --policy-mode digest` (no default-mode prose — §14.2 #1).
- **Constraints:** Spec-defined keys ONLY — no invented POLICY keys (would trip the
  keyword/namespace invariant). Both files pass `go run ./cmd/agentrc lint <file>`
  from repo root (§0.7). Canonical `hello` untouched; if any hello fragment is
  mirrored it must diff clean against `examples/Agentfile.minimal` (§0.6, M-001).
  Version stays `draft.5` here (§0.1). Index-card links must resolve. Grep-locate the
  examples-index card block before editing (M-003).
- **Failure scenarios:** A non-spec POLICY key sneaks in and `arc lint` fails or
  implies a fourth namespace. An index-card link 404s (§V.7). Header/scaffold drifts
  from the existing example style. Accidental `draft.6` in a description string.
- **Success scenarios:** `go run ./cmd/agentrc lint examples/Agentfile.hooked` and
  `… Agentfile.delegator` both pass; each file has the full scaffold + verbatim
  healthcheck; examples index shows two new cards with resolving links; §V still one
  `draft.5`.
- **Connections:** Wave 1, parallel with S2-02's T17–T19. No file overlap with the
  spec sections. Off the critical path; feeds T27 release verification.

### S2-02 — Spec draft.6 (T17, T18, T19, T20)

- **What's wanted:** Author the three new POLICY namespaces and, in a separate late
  step, the supporting edits + sitewide version bump + T8 landing. **§8.x
  renumbering decision (this story's call):** the current spec already uses §8.5
  ("Auto-derived egress") and §8.6 ("Why POLICY is a keyword"), so the three new
  subsections are INSERTED at the next FREE numbers after §8.6 — **§8.7
  `substrate.<platform>.*` (T17), §8.8 `agent.auth.*` (T18), §8.9
  `substrate.runtime.language` (T19)** — NOT at 8.5/8.6 as the work order's shorthand
  suggested. Existing §8.5/§8.6 are left untouched. T17: tokens
  `aws|gcp|azure|kubernetes|local`; unknown tokens MUST parse, foreign-platform keys
  ignored (never error; linter MAY warn); labels
  `ai.agentrc.substrate.<platform>.<key>=<value>`; platform-scoped beats generic on
  that platform only; tightening-only across `FROM` per namespace; AWS registry
  `roleArn`/`networkMode`/`securityGroup`(rep.)/`subnet`(rep.)/`protocol`/`maxLifetime`/`deployment.mode`
  (`container` default|`code`)/`code.s3.uri`. T18: `agent.auth.mode` (`platform`
  default|`jwt`|`none`), `agent.auth.jwt.discovery_url`, `.allowed_audience`(rep.),
  `.allowed_client`(rep.); a platform that can't enforce a requested `jwt` authorizer
  MUST NOT expose the invocation endpoint; not a secret. T19:
  `substrate.runtime.language <language>:<version>`; optional; container-mode MAY
  ignore (base image authoritative); code-mode requires it or resolvable inference,
  else fail closed. **T20 (LATE, single commit):** §14.2 open decision #6 promotion
  candidates (`protocol`, `maxLifetime`); /docs/agentfile/ platform-scoped paragraph
  + JWT example; /profiles/core/ accepts unknown `substrate.<token>.*`;
  `Agentfile.code-reviewer` gains a commented `substrate.aws.*` + `agent.auth.*`
  block; CHANGELOG draft.6 entry; **sitewide bump `0.1.0-draft.5` → `0.1.0-draft.6`**;
  land **T8 Option A** — an informative spec §9 subsection "Reproducible builds /
  `agentrc.lock`" documenting only what the tooling does (`arc lock` writes
  `agentrc.lock`; `arc build` does NOT consume it), marked "Status: informative in
  this draft; format TODO", homepage slogan kept.
- **Constraints:** EXACTLY three new namespaces, no fourth; `substrate.kubernetes.serviceAccount`
  is a KEY under the existing §8.7 `substrate.<platform>.*`, NOT a new namespace
  (intake failure scenario). No new keyword. Syntax line
  `# syntax=agentrc.agentfile/v0.1` frozen. **T17–T19 land at draft.5** (§0.1); the
  bump is T20 ONLY, one commit, sitewide — must not be split across earlier tasks.
  T8 default is Option A unless owner overrides to B (open question #3, carried
  UNRESOLVED). POLICY lines are requests; Cedar platform-side only; secrets deferred
  (§0.4). Grep-locate every version string before the bump (M-003) — incl. schema
  `0.1.0-draft.5` in `internal/**/lock.go`, `_config.yml`/page descriptions, examples
  descriptions, CHANGELOG.
- **Failure scenarios:** A new subsection reuses §8.5/§8.6 (collision) or introduces a
  fourth namespace/keyword. The draft.6 bump misses a surface (homepage, spec, cli.md,
  examples, CHANGELOG, `lock.go`), leaving mixed draft.5/draft.6 (§V.3 finds >1). The
  bump leaks into an earlier task before T20. `substrate.*` POLICY namespace or
  "substrate-neutral" accidentally renamed. `agent.auth.*` documented as a secret.
- **Success scenarios:** §8.7/§8.8/§8.9 present with exactly the documented keys and
  fail-closed semantics; `rg 'POLICY substrate\.' spec/index.md` intact; keyword count
  unchanged; after T20 exactly one `draft.N` value = `draft.6`; T8 Option A subsection
  present and marked informative; slogan kept.
- **Connections:** T17–T19 = Wave 1 (parallel with S2-01, at draft.5). T20 = Wave 4
  (the single sitewide bump, AFTER translators consume the namespaces). On the critical
  path: S2-02(T17–19) → S2-03 → S2-04 → T20 → S2-07. Shares `spec/index.md` with no
  other story except its own T20; shares `examples/Agentfile.code-reviewer` (T20 comment
  block) — sequence T20 after S2-01 to avoid an example-edit collision.

### S2-03 — CLI `--backend` rename (T21)

- **What's wanted:** Rename `--substrate` → `--backend` everywhere (CLI code, help,
  `/cli/`, quickstart step 5, tooling README) and add the backend subcommand surface.
  `arc run <ref> --backend local|bedrock|kubernetes`; local:
  `[--isolation microvm|container]` (microsandbox MVP exists); bedrock:
  `[--region] [--profile] [--dry-run]`; kubernetes:
  `[--kubeconfig] [--namespace] [--dry-run]`. `--dry-run` prints the translated config
  and exits. Record in CLI docs: GCP dropped (Agent Runtime Python-only managed; GKE
  via kubernetes backend); Docker Compose dropped (no `network.*` egress enforcement
  without a bespoke sidecar). Go + docs; **REAL TDD** (RED→GREEN→refactor).
- **Constraints:** The rename is CLI-flag-only — MUST NOT touch the `substrate.*`
  POLICY namespace or the "substrate-neutral" concept (§0.2; §V.8). Grep-first every
  `--substrate` occurrence (M-003): `cmd/agentrc/run.go:19`, `cli.md:86/103/120`, and
  any quickstart/tooling-README mention. Version stays `draft.5` (§0.1). ENV: local
  disk near-full → scope Go gates to `./cmd/agentrc/...`; full `go build ./...` /
  `-race` is CI's job (open question #2, carried UNRESOLVED — feasibility flagged).
- **Failure scenarios:** A blind find-replace of "substrate" renames the POLICY
  namespace or "substrate-neutral" (breaks §V.8 `rg 'POLICY substrate\.' spec/`). A
  `--substrate` mention survives on one surface (§V.8 count > 0). `run` help text left
  stale (M-002). New deps untidy `go.sum` (CI `go build ./...` fails).
- **Success scenarios:** `rg -l -- '--substrate' cmd/ docs/ cli.md` → 0;
  `rg 'POLICY substrate\.' spec/index.md` intact; `arc run --help | grep -- '--backend'`
  present; per-backend flags parse; `TestBackendFlagReplacesSubstrate` green.
- **Connections:** Wave 2. On the critical path (after S2-02's T17–19, before S2-04).
  Shares `cli.md` with S2-05 (T25) — sequence T25 (Wave 4) after T21 (Wave 2). The
  surface T22–T24 hang off.

### S2-04 — Backend translators (T22, T23, T24)

- **What's wanted:** Three fail-closed reference translators mapping `ai.agentrc.*`
  labels + image config → each platform's native config, each asserted through
  `--dry-run` output. T22 `local`: wire the existing microsandbox VMM MVP under
  `--backend local` (default); plumbing only + §0.8 positioning line. T23 `bedrock`:
  map labels + image config → Bedrock `CreateAgentRuntime` fields (13/13:
  agentRuntimeName/description ← IDENTITY; containerUri ← OCI ref; roleArn ←
  substrate.aws.roleArn; networkMode ← substrate.aws.networkMode;
  securityGroups/subnets ← substrate.aws.securityGroup/subnet; serverProtocol ←
  substrate.aws.protocol; env ← image Env; customJWTAuthorizer ← agent.auth.jwt.*;
  idleRuntimeSessionTimeout ← agent.idle_timeout; maxLifetime ← substrate.aws.maxLifetime;
  codeConfiguration ← deployment.mode=code + code.s3.uri + substrate.runtime.language).
  Fail closed: missing `roleArn`; unenforceable `agent.auth.mode=jwt`; code-mode
  without resolvable language. `--dry-run` emits the translated JSON. T24 `kubernetes`:
  emit (dry-run) or apply Deployment (resources from substrate.runtime.*, env from
  image config), Service, deny-by-default NetworkPolicy from `POLICY network dns:*`,
  ServiceAccount from `substrate.kubernetes.serviceAccount`, MCP servers from
  `/mnt/mcp/*` as sidecars. **ONE emission format = manifests** (default; Helm NOT
  emitted — open question #4, carried UNRESOLVED). Go; **REAL TDD**; each translator
  MUST have fail-closed tests.
- **Constraints:** Fail-closed is REQUIRED — a translator MUST NOT emit valid/exposed
  config when a required grant is missing (§0-standards; work-order T23/T24). Reads the
  §8.7/§8.8/§8.9 namespaces (must land first). Prefer plain-struct/YAML emission over
  heavy SDK/k8s clients to avoid enlarging the buildkit-heavy dep tree (ENOSPC ENV
  risk); any added dep keeps `go mod tidy` clean. `substrate.kubernetes.serviceAccount`
  is a KEY under §8.7, not a new namespace. `local` wiring depth is dry-run/plumbing
  acceptable this sprint (open question #5, carried UNRESOLVED). Version stays draft.5
  (§0.1). Scope Go gates to `./cmd/agentrc/...`; full build/-race in CI.
- **Failure scenarios:** A translator emits an invocation endpoint / valid config
  despite missing `roleArn`, an unenforceable `agent.auth.mode=jwt`, or code-mode
  without resolvable `substrate.runtime.language` (not fail-closed). k8s emits BOTH
  manifests and Helm. `--backend bedrock --dry-run` fails `python3 -m json.tool` or
  `--backend kubernetes --dry-run` fails yaml-parse/kubeconform. A translator invents a
  fourth POLICY namespace. §0.8 line paraphrased.
- **Success scenarios:** `--backend bedrock --dry-run | python3 -m json.tool` parses;
  `--backend kubernetes --dry-run` yaml-parses / kubeconform-validates and includes a
  deny-by-default NetworkPolicy; all fail-closed tests green (RED first); §0.8 line
  verbatim; single k8s format.
- **Connections:** Wave 3. On the critical path (after S2-03, before T20). Consumes
  S2-02's §8.7/§8.8/§8.9 as translation inputs. All-Go; no markdown-surface collisions.

### S2-05 — CLI docs table (T25)

- **What's wanted:** Update the `/cli/` status table: `run` →
  `implemented (local, bedrock, kubernetes — reference translators)`; `sign`/`verify`
  stay `planned`; place the §0.8 positioning line VERBATIM above the table. Re-derive
  every status claim from the now-real backends (M-002).
- **Constraints:** §0.8 line VERBATIM, never paraphrased (§0.6-standards). `sign`/`verify`
  MUST stay `planned` (they remain stubs). Grep-locate the table + the pre-rename
  `--substrate`/`--isolation` prose in `cli.md` (lines 103/120) so the docs reflect
  `--backend`. Version stays draft.5 until T20 (§0.1). No spec §10 status claim left
  stale (M-002).
- **Failure scenarios:** §0.8 line paraphrased. `run` still shows `planned` after
  backends land, or `sign`/`verify` flipped to implemented (M-002). A `--substrate`
  mention survives in the same table/prose (breaks §V.8).
- **Success scenarios:** `/cli/` table shows `run` implemented with the reference-translator
  qualifier and the verbatim §0.8 line above it; `sign`/`verify` `planned`;
  `rg -- '--substrate' cli.md` → 0.
- **Connections:** Wave 4 (after S2-04 backends land, alongside T20). Shares `cli.md`
  with S2-03 (T21) — sequenced after it.

### S2-06 — Demo (T26)

- **What's wanted:** A one-agent/three-backend demo:
  `arc build -t ghcr.io/agentrc/code-reviewer:1.0 .`;
  `arc run … --backend local --isolation microvm`;
  `arc run … --backend bedrock --dry-run`;
  `arc run … --backend kubernetes --dry-run`. Narrative VERBATIM: "Same artifact, same
  labels, three substrates. The translators are the proof of concept; the labels are
  the standard."
- **Constraints:** Narrative string VERBATIM (§0.8-standards; work-order T26). Uses the
  real `Agentfile.code-reviewer` (with the T20 commented `substrate.aws.*` +
  `agent.auth.*` block, so demo runs after T20) and all three real backends (needs
  S2-04). `bedrock`/`kubernetes` dry-run outputs must satisfy §V.9. Version = draft.6
  by this point.
- **Failure scenarios:** Narrative paraphrased. A demo command uses the removed
  `--substrate` flag or a dropped backend (gcp/compose). Dry-run output fails §V.9
  parse.
- **Success scenarios:** Demo doc/section runs the four commands end-to-end; narrative
  verbatim; both dry-runs parse (§V.9).
- **Connections:** Wave 5 (after T20 + S2-04 + S2-05). Off the strict critical path but
  gates the demo portion of §V and T27.

### S2-07 — Release & live verification (T27)

- **What's wanted:** Ship and verify. Run full §V incl. checks 8–9 locally where
  feasible / CI otherwise; commit with task-tagged messages; open + squash-merge a PR to
  `master` (protected: PR required); wait for GH Pages build + CDN. Then live re-verify
  each Sprint-2 task against `https://agentrc.ai` (version = draft.6 sitewide; `--backend`
  in live CLI docs; new examples lint + render; addendum committed iff owner-provided —
  open question #1, carried UNRESOLVED); record per-task pass/fail; owner sign-off;
  retrospective.
- **Constraints:** No direct push to `master` — deploy = merge the sprint branch via PR
  (M-005: never run supervisor git in a shared tree while a worker is active — serialize).
  Do not release until all of S2-01…S2-06 are in and local §V (where feasible) + CI is
  green. Live HTML greps must account for syntax-highlight token splitting (M-006) and
  for the pipeline `docs/agents/` artifacts being excluded from the built site (M-004).
  This story is supervisor-executed, not worker-TDD; its plan.md references the §V suite +
  live checks, not RED unit tests.
- **Failure scenarios:** A direct push to `master` rejected by the ruleset (M-005).
  Release before §V/CI green. A live check FAILs post-deploy and is not fixed within the
  sprint. `docs/agents/` artifacts leak into `_site`/sitemap (M-004). A multi-token live
  grep false-negatives (M-006).
- **Success scenarios:** Local/CI §V all green incl. checks 3/8/9; PR merged to `master`;
  GH Pages + CDN serve draft.6 with `--backend` docs and both new examples; live per-task
  verification recorded all green; owner sign-off; retrospective written to `docs/agents/`.
- **Connections:** Wave 6, last node on the critical path. Depends on every other story +
  local/CI §V.

## Story dependency graph

```
S2-01 (examples) ──────────────┐
                               ├─▶ (T27 release + live verify) = S2-07
S2-02 T17–T19 (spec sections) ─┴─▶ S2-03 (--backend rename) ─▶ S2-04 (translators)
                                                                    │
                                        S2-02 T20 (version bump) ◀──┘  ┐
                                        S2-05 (CLI docs table) ────────┼─▶ S2-06 (demo) ─▶ S2-07
                                                                       ┘

Wave 1 (parallel): S2-01 ; S2-02 T17–T19 (spec sections, at draft.5 — NOT the bump)
Wave 2:            S2-03 (--backend rename)
Wave 3:            S2-04 (translators; consume §8.7/§8.8/§8.9)
Wave 4:            S2-02 T20 (single sitewide draft.5→draft.6 commit + T8 A) ; S2-05 (CLI docs)
Wave 5:            S2-06 (demo)
Wave 6:            S2-07 (release + live verify)
Critical path:     S2-02(T17–19) → S2-03 → S2-04 → T20 → S2-07
```

Note: S2-02 spans waves — its spec sections (T17–T19) land in Wave 1 at `draft.5`; its
T20 (the single sitewide version bump + T8) is deliberately deferred to Wave 4 so the
bump is one late commit after the translators consume the new namespaces (§0.1, intake
"Connections").

## Notes / open questions

Carried forward from `intake.md` for the supervisor to surface to the owner — NONE
resolved here (per §0.5 and the planning constraint to carry, not resolve):

1. **Missing addendum (provenance blocker).** `agentrc-draft6-addendum-2026-07-02.md`
   (the Phase-4 "full text" source) is ABSENT from the repo. T17–T19 can be authored
   from the work order's inline detail, but "commit the addendum for provenance" cannot.
   Default: author T17–T19 from inline detail and skip committing a separate addendum.
   Owner: provide the file, or confirm the inline-only default? (Affects T27 live check.)
2. **ENV / Go build feasibility.** Local disk is near-full and Go builds pulling the
   buildkit dep tree have hit ENOSPC. Is real-TDD local execution of T21–T24 feasible, or
   should those tasks be authored + validated via CI only (or deferred)? Plan assumes
   scoped-local (`./cmd/agentrc/...`) + full/`-race` in CI, per the gate matrix.
3. **T8 lockfile A/B.** Proceeding with **Option A** (informative spec §9 subsection
   "Reproducible builds / `agentrc.lock`", marked "Status: informative; format TODO; not
   consumed by `build` today", homepage slogan kept) unless owner picks Option B (cut
   slogan; lockfile on `/cli/` only). Confirm Option A? (Lands in T20.)
4. **k8s emission format.** T24 requires ONE of manifests OR Helm. Default: **manifests**.
   Owner: confirm manifests, or prefer Helm?
5. **`local` backend wiring depth.** Does the owner want the microsandbox `local` backend
   (T22) actually wired even if its VMM MVP isn't runnable in CI/this env, or is
   dry-run / plumbing-only acceptable for this sprint?

Plus the planner's own resolved-for-planning decision (stated, not an owner blocker):

- **§8.x renumbering.** The work order's shorthand (T17→§8.5, T18→§8.6, T19→§8.7)
  collides with the existing §8.5 ("Auto-derived egress") and §8.6 ("Why POLICY is a
  keyword"). Decision: INSERT the three new subsections at the next FREE numbers after
  §8.6 — **§8.7 `substrate.<platform>.*`, §8.8 `agent.auth.*`, §8.9
  `substrate.runtime.language`** — leaving existing §8.5/§8.6 untouched. Applied in S2-02.
</content>
