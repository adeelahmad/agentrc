---
type: tasks
story: S2-02
---
# S2-02 Tasks — Spec draft.6 (T17, T18, T19, T20)

**§8.x renumbering decision:** the existing spec already uses §8.5 ("Auto-derived
egress", `spec/index.md:400`) and §8.6 ("Why POLICY is a keyword", `spec/index.md:415`).
The three new subsections are INSERTED at the next FREE numbers after §8.6 —
**§8.7 / §8.8 / §8.9** — leaving §8.5/§8.6 untouched. Do NOT reuse 8.5/8.6.

T17–T19 land at `0.1.0-draft.5`. The version bump is T20 ONLY (single sitewide commit).
Grep-locate every target first (§0.9, M-003).

## T17 — §8.7 `substrate.<platform>.*` [P1]

Insert a new `### 8.7 \`substrate.<platform>.*\`` subsection after §8.6 in `spec/index.md`:
- Platform tokens `aws | gcp | azure | kubernetes | local`; UNKNOWN tokens MUST parse;
  foreign-platform keys are ignored, NEVER an error (linter MAY warn).
- Labels: `ai.agentrc.substrate.<platform>.<key>=<value>`.
- Platform-scoped beats generic (`substrate.*`) ON THAT PLATFORM ONLY.
- Tightening-only across `FROM` per namespace.
- AWS key registry: `roleArn`, `networkMode`, `securityGroup` (repeatable), `subnet`
  (repeatable), `protocol`, `maxLifetime`, `deployment.mode` (`container` default | `code`),
  `code.s3.uri`.
- This is a KEY namespace under the existing `substrate.*` — NOT a rename of it (§0.2).
  `substrate.kubernetes.serviceAccount` is a KEY here, not a separate namespace (§0.3).

## T18 — §8.8 `agent.auth.*` (generic, fail-closed) [P1]

Insert `### 8.8 \`agent.auth.*\`` after §8.7:
- `agent.auth.mode` (`platform` default | `jwt` | `none`); `agent.auth.jwt.discovery_url`;
  `agent.auth.jwt.allowed_audience` (repeatable); `agent.auth.jwt.allowed_client` (repeatable).
- Fail-closed: a platform that CANNOT enforce a requested `jwt` authorizer MUST NOT expose
  the invocation endpoint.
- Explicitly NOT a secret — generic authZ config (§0.4); no inline Cedar, no secret keyword.

## T19 — §8.9 `substrate.runtime.language` [P1]

Insert `### 8.9 \`substrate.runtime.language\`` after §8.8:
- Value `<language>:<version>`. Optional.
- Container-mode MAY ignore it (base image authoritative).
- Code-mode requires it or a resolvable inference, else FAIL CLOSED.

## T20 — Supporting edits + version bump + T8 landing [P1] (LATE — single commit, Wave 4)

Do this as ONE sitewide commit, only after S2-04 translators land:
- §14.2 open decision #6 area: add `protocol`, `maxLifetime` as promotion candidates (they
  live under §8.7 today; note the promotion question is open — do not resolve, §0.5).
- `/docs/agentfile/` (`docs/agentfile.md`): add a platform-scoped `substrate.<platform>.*`
  paragraph + a `agent.auth.jwt.*` example.
- `/profiles/core/` (`profiles/core.md` or equivalent — grep-locate): state it accepts
  unknown `substrate.<token>.*` (parse, ignore foreign).
- `examples/Agentfile.code-reviewer`: add a COMMENTED `substrate.aws.*` + `agent.auth.*`
  block (commented so it stays lint-clean and demonstrates the new keys). Sequence after S2-01.
- CHANGELOG: add a `0.1.0-draft.6` entry summarizing §8.7/§8.8/§8.9 + the flag/backends.
- **Sitewide bump `0.1.0-draft.5` → `0.1.0-draft.6`** — grep-locate EVERY occurrence first
  (M-003): `spec/index.md:4,9`, page descriptions in `_config.yml`/front-matter, `cli.md`,
  `examples/*` descriptions, CHANGELOG, and the schema string `0.1.0-draft.5` in
  `internal/**/lock.go`. The syntax line `# syntax=agentrc.agentfile/v0.1` stays frozen.
- **T8 Option A:** add an informative `### 9.x Reproducible builds / \`agentrc.lock\``
  subsection under §9 documenting ONLY what tooling does (`arc lock` writes `agentrc.lock`;
  `arc build` does NOT consume it today — see `docs/agents/sprint1/lockfile-report.md`),
  marked "Status: informative in this draft; format TODO". Keep the homepage slogan.
  (Owner may override to Option B — open question #3, UNRESOLVED.)

Verify: §8.7/§8.8/§8.9 present; `rg 'POLICY substrate\.' spec/index.md` intact; keyword
count unchanged; AFTER T20 `grep -rhoE 'draft\.[0-9]+' .` single value == `draft.6`; T8
subsection present + marked informative; slogan kept.
</content>
