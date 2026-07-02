---
type: plan
story: S2-02
scope: "tests only"
---
# S2-02 Test Plan — Spec draft.6 (RED)

Tests only. Each check is a function in `scripts/verify-sprint2.sh` (created by the RED
worker). T17–T19 checks run at `draft.5`; T20 checks assert the post-bump `draft.6` state.

## T17 — §8.7 `substrate.<platform>.*`

- [ ] `scripts/verify-sprint2.sh::t17_substrate_platform_section` — assert `### 8.7`
      heading naming `substrate.<platform>` exists in `spec/index.md`. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint2.sh::t17_aws_key_registry` — assert `roleArn`, `networkMode`,
      `securityGroup`, `subnet`, `protocol`, `maxLifetime`, `deployment.mode`, `code.s3.uri`
      all documented in the section. RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t17_unknown_tokens_parse` — assert prose states unknown
      platform tokens MUST parse / foreign keys ignored (never error). RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t17_no_collision_85_86` — assert existing `### 8.5` and
      `### 8.6` headings are unchanged (count == 2, still egress + why-POLICY). RED now: green;
      regression guard against reusing 8.5/8.6.

## T18 — §8.8 `agent.auth.*`

- [ ] `scripts/verify-sprint2.sh::t18_agent_auth_section` — assert `### 8.8` heading naming
      `agent.auth` exists. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint2.sh::t18_auth_modes_and_jwt_keys` — assert `mode`
      (`platform|jwt|none`), `jwt.discovery_url`, `allowed_audience`, `allowed_client`
      documented. RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t18_auth_failclosed_prose` — assert "MUST NOT expose the
      invocation endpoint" fail-closed sentence present; NOT described as a secret. RED now: FAILS.

## T19 — §8.9 `substrate.runtime.language`

- [ ] `scripts/verify-sprint2.sh::t19_runtime_language_section` — assert `### 8.9` heading
      naming `substrate.runtime.language` exists. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint2.sh::t19_codemode_failclosed` — assert prose states code-mode
      requires the language or resolvable inference else fail-closed; container-mode MAY ignore.
      RED now: FAILS.

## T20 — Supporting edits + version bump + T8 (post-bump state)

- [ ] `scripts/verify-sprint2.sh::t20_version_draft6_sitewide` — assert exactly one `draft.N`
      value tree-wide AND it equals `draft.6`. RED now: FAILS (tree is `draft.5` pre-T20).
- [ ] `scripts/verify-sprint2.sh::t20_lockgo_bumped` — assert `internal/**/lock.go` no longer
      hardcodes `0.1.0-draft.5`. RED now: FAILS — still `draft.5`.
- [ ] `scripts/verify-sprint2.sh::t8a_spec_lock_subsection` — assert the informative
      "Reproducible builds / `agentrc.lock`" §9 subsection exists, marked "format TODO", and the
      homepage slogan is kept. RED now: FAILS — subsection absent.
- [ ] `scripts/verify-sprint2.sh::t20_docs_agentfile_platform_para` — assert `docs/agentfile.md`
      gains a platform-scoped paragraph + a `agent.auth.jwt.*` example. RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t20_code_reviewer_commented_block` — assert
      `examples/Agentfile.code-reviewer` has a COMMENTED `substrate.aws.*` + `agent.auth.*` block
      and still lints. RED now: FAILS — block absent.
- [ ] `scripts/verify-sprint2.sh::t20_changelog_draft6` — assert CHANGELOG has a `0.1.0-draft.6`
      entry. RED now: FAILS.

## Suite guards

- [ ] `scripts/verify-sprint2.sh::v3_version_draft6` — assert single `draft.N` value ==
      `draft.6` AFTER T20 (§V.3). RED now: FAILS pre-T20.
- [ ] `scripts/verify-sprint2.sh::v8_terminology_split` — assert `rg 'POLICY substrate\.'
      spec/index.md` intact (the spec sections must NOT rename the namespace). RED now: green;
      regression guard.
- [ ] `scripts/verify-sprint2.sh::v_keyword_count` — assert exactly four keywords remain
      (IDENTITY, CAPABILITY, SOP, POLICY); no fifth. RED now: green; regression guard.
</content>
