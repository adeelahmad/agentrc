---
type: intake
sprint: 2
---
# Sprint 2 Intake — "Features" (T15–T27)

Interrogated from `docs/agents/sprint2/work-order.md` (T15–T27 + §0 invariants +
the supervisor-carried resolutions/flags), grounded against `docs/agents/memory.md`
(M-001..M-006), `docs/agents/sprint1/lockfile-report.md` (T8), `cli.md`,
`spec/index.md` §8 (POLICY namespaces) and §10 (CLI surface), and the real CLI/backend
code under `cmd/agentrc/` (`run.go`, `build.go`, …) and `internal/`. This is intake
only: no plan, no code. Blocking ambiguities are listed under Open questions, not
resolved.

## What's wanted

Sprint 2 turns agentrc from "consistent + honest site" (Sprint 1) into "demonstrably
substrate-neutral" by shipping five things in order: (1) two new lint-clean examples
(`Agentfile.hooked`, `Agentfile.delegator`) using only spec-defined POLICY keys; (2) a
spec **draft.6** that adds **exactly three** new POLICY namespaces —
`substrate.<platform>.*`, `agent.auth.*`, `substrate.runtime.language` — and nothing
else; (3) the CLI **`--substrate` → `--backend`** rename with three reference
translators (`local`, `bedrock`, `kubernetes`) that map `org.agentrc.*` labels to each
platform's native config, all fail-closed; (4) a one-agent/three-backend demo proving
"same artifact, same labels, three substrates"; and (5) a full release with live
re-verification. The version string stays `0.1.0-draft.5` until T20, which bumps it
sitewide to `0.1.0-draft.6` in one commit. Success is §V (incl. new checks 8 and 9)
passing locally, each task's Verify passing, and — after deploy — the same checks
passing live (T27). This is a single confirmable direction: **new examples + 3 new
POLICY namespaces (draft.6) + `--backend` translators + three-backend demo + release.**

- **T15 [P2]** — `examples/Agentfile.hooked`: demonstrate `POLICY agent.hooks.on_tool_call <https>` + one `pre`/`post` hook; comments explain platform auto-derives the hook endpoint's egress grant with `.source` attribution (requested vs derived); include one explicit `POLICY network dns:…` for contrast.
- **T16 [P2]** — `examples/Agentfile.delegator`: demonstrate `POLICY agent.sub_agents true`, `agent.sub_agents.max <n>`, `agent.sub_agent_timeout <dur>`; comment that sub-agent grants are the platform's call. Spec-defined keys only. Both examples get the standard scaffold + examples-index card + `arc lint` clean.
- **T17 [P1]** — Spec §8.5 `substrate.<platform>.*`: tokens `aws|gcp|azure|kubernetes|local`; unknown tokens MUST parse, foreign-platform keys ignored (never error, linter MAY warn); AWS registry `roleArn`/`networkMode`/`securityGroup`(rep.)/`subnet`(rep.)/`protocol`/`maxLifetime`/`deployment.mode`/`code.s3.uri`; platform-scoped beats generic on that platform only; tightening-only across `FROM`.
- **T18 [P1]** — Spec §8.6 `agent.auth.*` (generic, fail-closed authZ): `agent.auth.mode` (`platform` default|`jwt`|`none`), `agent.auth.jwt.discovery_url`, `.allowed_audience`(rep.), `.allowed_client`(rep.). Platform that can't enforce a requested `jwt` authorizer MUST NOT expose the invocation endpoint. Not a secret.
- **T19 [P1]** — Spec §8.7 `substrate.runtime.language`: `<language>:<version>`; optional; container-mode MAY ignore (base image authoritative); code-mode requires it or resolvable inference, else fail-closed.
- **T20 [P1]** — Supporting edits + **sitewide version bump draft.5 → draft.6 (one commit)** + T8 landing: §14.2 open decision #6 promotion candidates (`protocol`, `maxLifetime`); /docs/agentfile/ platform-scoped paragraph + JWT example; /profiles/core/ accepts unknown `substrate.<token>.*`; `Agentfile.code-reviewer` commented `substrate.aws.*` + `agent.auth.*` block; CHANGELOG draft.6 entry; implement T8 choice (default Option A).
- **T21 [P0* demo]** — Flag rename `--substrate` → `--backend` everywhere (CLI code, help, /cli/, quickstart step 5, tooling README): `arc run <ref> --backend local|bedrock|kubernetes` with per-backend flags; `--dry-run` prints translated config and exits. Record: GCP dropped, Docker Compose dropped. Verify `rg -- '--substrate' -l` → 0; `rg 'POLICY substrate\.' spec/` intact.
- **T22 [P1]** — Backend `local`: wire the existing microsandbox VMM MVP under `--backend local` (default). Plumbing only + §0.8 positioning line.
- **T23 [P1]** — Backend `bedrock`: map `org.agentrc.*` labels + image config → Bedrock `CreateAgentRuntime` fields (13/13 mapping). Fail-closed on missing `roleArn`, unenforceable `agent.auth.mode=jwt`, code-mode without resolvable language. `--dry-run` emits translated config.
- **T24 [P1]** — Backend `kubernetes`: emit (dry-run) or apply Deployment / Service / deny-by-default NetworkPolicy (from `POLICY network dns:*`) / ServiceAccount / MCP-server sidecars. ONE emission format (manifests OR Helm), not both.
- **T25 [P1]** — CLI docs table: `run` → `implemented (local, bedrock, kubernetes — reference translators)`; `sign`/`verify` stay `planned`; §0.8 positioning line above the table.
- **T26 [P2]** — One-agent/three-backends demo: `arc build`, then `--backend local --isolation microvm`, `--backend bedrock --dry-run`, `--backend kubernetes --dry-run`. Verbatim narrative on "Same artifact, same labels, three substrates."
- **T27 [P0]** — RELEASE + live verification: full §V incl. checks 8–9; deploy; live re-verify each Sprint-2 task; per-task pass/fail; owner sign-off; retrospective.

## Constraints

Global invariants from work-order §0 (never violate):

1. **Version discipline** — stays `0.1.0-draft.5` through Phase 3; bumps to
   `0.1.0-draft.6` **only in T20, sitewide, in ONE commit**. The syntax line
   `# syntax=agentrc.agentfile/v0.1` is frozen throughout.
2. **Exactly four keywords** — IDENTITY, CAPABILITY, SOP, POLICY. No new keywords ever.
3. **EXACTLY three new POLICY namespaces** — `substrate.<platform>.*` (T17),
   `agent.auth.*` (T18), `substrate.runtime.language` (T19). No others. No fourth.
4. **Terminology split (M-003 grep-first)** — `--backend` is the CLI flag
   (`local|bedrock|kubernetes`); the old `--substrate <driver>` flag + "driver" wording
   are renamed in T21. The `substrate.*` **POLICY namespace** and the concept
   **"substrate-neutral"** are **NEVER renamed** — the rename is CLI-flag-only.
5. **POLICY lines are requests** — Cedar is platform-side only; secrets are deferred;
   no inline Cedar, no secret keyword. `agent.auth.*` is generic fail-closed authZ
   config, NOT a secret.
6. **§0.8 positioning line, VERBATIM** on backend docs: "Reference translators — a
   proof of concept until platforms read `org.agentrc.*` labels natively. Not
   production runners." (never paraphrased).
7. **ONE canonical hello** — byte-identical wherever rendered inline and identical to
   `examples/Agentfile.minimal` (M-001: diff, don't grep-substring).
8. **Every example passes `arc lint`** via `go run ./cmd/agentrc lint <file>` from repo
   root (CLI is at `./cmd/agentrc`, NOT `tooling/`).
9. **Fail-closed backends** — a translator must NOT emit valid/exposed config when a
   required grant is missing (no `roleArn`; unenforceable `jwt`; code-mode without
   language).
10. **Master requires a PR**; grep-locate every edit target first (M-003); build-site +
    §V before commit (build-dependent checks delegated to CI + live, no local Jekyll —
    M-004/M-006).

## Failure scenarios

- **§0.2/0.4 violation** — the T21 sweep accidentally renames the `substrate.*` POLICY
  namespace or the "substrate-neutral" concept (e.g. blind find-replace of "substrate"),
  breaking the spec/terminology split. (`rg 'POLICY substrate\.' spec/` must stay intact.)
- **Fourth namespace** — a translator or spec edit introduces a POLICY namespace beyond
  the three sanctioned ones (e.g. inventing `substrate.kubernetes.serviceAccount` as a
  *new namespace* rather than a key under the existing `substrate.<platform>.*`).
- **Split-version drift** — the draft.6 bump misses a surface (homepage, spec, cli.md,
  examples, CHANGELOG, schema `0.1.0-draft.5` in `lock.go`), leaving mixed draft.5/draft.6;
  §V check 3 must find exactly one `draft.N` value.
- **Backend not fail-closed** — a translator emits an invocation endpoint / valid config
  despite missing `roleArn`, an unenforceable `agent.auth.mode=jwt`, or code-mode without
  resolvable `substrate.runtime.language`.
- **Double k8s emission** — T24 emits BOTH manifests and Helm (work order mandates one).
- **Invalid dry-run output** — `--backend bedrock --dry-run` fails `python3 -m json.tool`
  or `--backend kubernetes --dry-run` fails kubeconform/yaml-parse.
- **Example drift / lint break** — new examples fail `arc lint`, use non-spec keys, or
  drift the canonical hello (M-001); or examples-index links don't resolve.
- **§0.8 paraphrase** — the positioning line is reworded instead of quoted verbatim on
  any backend doc surface.
- **Stale status prose (M-002)** — leaving `run` as `planned`/`not implemented` in
  cli.md/run.go help/spec §10 after backends land, or leaving the pre-rename
  `--substrate` help text.

## Success scenarios

Acceptance = §V passes (locally where possible; CI + live otherwise), including the two
new checks:

- **Check 8 (terminology split)** — `rg -- '--substrate' -l | wc -l` → 0 across CLI/docs;
  `rg 'POLICY substrate\.' spec/` still returns the namespace intact.
- **Check 9 (backends)** — `arc run --help | grep -- '--backend'` present;
  `arc run <ref> --backend bedrock --dry-run | python3 -m json.tool` parses;
  `arc run <ref> --backend kubernetes --dry-run | kubeconform -` (or yaml-parse) valid.
- **Check 3 (version coherence)** — exactly one `draft.N` value = `draft.6` after T20
  (was `draft.5` before).
- Both new examples pass `go run ./cmd/agentrc lint` and appear in the examples index in
  the existing card style; canonical hello unchanged (byte-diff).
- Spec §8.5/§8.6/§8.7 define exactly the three namespaces with their documented keys and
  fail-closed semantics; no keyword count change; §0.8 line verbatim on backend docs.
- Three-backend demo (T26) runs end-to-end with the verbatim narrative.
- **T27 live re-verify**: version = draft.6 sitewide live; `--backend` in live CLI docs;
  new examples lint + render; addendum committed iff owner-provided; per-task pass/fail
  recorded; owner sign-off; retrospective written.

## Connections

- **New namespaces precede the backends that consume them**: T17–T19 (spec) must land
  before T23/T24, which read `substrate.aws.*` (T17), `agent.auth.jwt.*` (T18), and
  `substrate.runtime.language` (T19) as translation inputs.
- **Rename precedes translators**: T21 (`--substrate` → `--backend`) is the surface T22–T24
  hang off; it must land first.
- **Version bump is a late single commit**: T20 is the ONE sitewide draft.5 → draft.6
  commit; it also lands the T8 A/B choice and the supporting spec/doc edits. It gates the
  version-coherence check but must not be split across earlier tasks.
- **Docs after backends**: T25 (CLI docs table) reflects T22–T24's real status; T26 (demo)
  needs all three backends; T27 (release) after everything + full §V.
- **Toolchain split**: T21–T24 are Go work requiring the real toolchain and TDD
  (`go build/vet/test`, table-driven, `-race`); T15–T16 and T17–T20/T25 are
  markdown/Agentfile edits validated by `arc lint` + grep + CI/live.
- **Grounding reality**: `cmd/agentrc/run.go` currently returns "not implemented" and
  carries `--isolation`/`--substrate`; `build.go` shells to `docker build`; `lock.go`
  hardcodes `0.1.0-draft.5` — all are draft.6-bump / rename targets to grep-locate.
- **Memory ties**: M-001 (diff snippets vs canonical), M-002 (re-derive status),
  M-003 (grep-first), M-004/M-006 (build + live-check caveats) all apply directly.

## In scope

- T15–T27 exactly, as specified in `docs/agents/sprint2/work-order.md`.

## Out of scope

- Anything beyond T27.
- Renaming the `substrate.*` POLICY namespace or the "substrate-neutral" concept
  (CLI flag rename only).
- Adding a fourth new POLICY namespace or any new keyword.
- Production-grade runners — the three translators are proof-of-concept only.
- §H owner/infra items (CDN purge for /spec/, agentrc.io DNS 301, name-collision).
- Secrets / inline Cedar / a secret keyword.
- GCP and Docker Compose backends (explicitly dropped in T21).

## Open questions

These are genuine owner blockers surfaced for the planning gate; none is resolved here.

1. **Missing addendum (provenance blocker).** `agentrc-draft6-addendum-2026-07-02.md`
   (the Phase-4 "full text" source) is **ABSENT** from the repo. T17–T19 can be authored
   from the work order's inline detail, but "commit the addendum for provenance" cannot.
   Default: author T17–T19 from inline detail and **skip committing a separate addendum**.
   Owner: provide the file, or confirm the inline-only default?
2. **ENV / Go build feasibility.** Local disk is near-full and Go builds pulling the
   buildkit dep tree have hit ENOSPC. Is real-TDD local execution of T21–T24 feasible, or
   should those tasks be authored + validated **via CI only** (or deferred)?
3. **T8 lockfile A/B.** Proceeding with **Option A** (informative spec §9 subsection
   "Reproducible builds / `agentrc.lock`", marked "Status: informative; format TODO; not
   consumed by `build` today", keep homepage slogan) unless owner picks **Option B** (cut
   slogan; lockfile on /cli/ only). Confirm Option A?
4. **k8s emission format.** T24 requires ONE of manifests OR Helm. Default: **manifests**.
   Owner: confirm manifests, or prefer Helm?
5. **`local` backend wiring depth.** Does the owner want the microsandbox `local` backend
   (T22) **actually wired** even if its VMM MVP isn't runnable in CI/this env, or is
   **dry-run / plumbing-only** acceptable for this sprint?

> Note for the planner: the work order assigns T17→§8.5, T18→§8.6, T19→§8.7, but the
> current `spec/index.md` already uses §8.5 ("Auto-derived egress") and §8.6 ("Why POLICY
> is a keyword"). Section renumbering vs. insertion is a plan-level decision, flagged here
> for awareness only (not an intake blocker).
</content>
</invoke>
