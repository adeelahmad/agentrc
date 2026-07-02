#!/usr/bin/env bash
# Sprint 2 §V verification suite (LOCAL). RED phase: written before the fixes.
# Each function asserts the DESIRED post-fix end-state and prints PASS/FAIL/SKIP.
# Return codes per check: 0 = PASS, 1 = FAIL, 2 = SKIP.
#
# NOTE ON SELF-MATCH: a few checks grep the whole tree for a verbatim literal
# (e.g. the T26 demo narrative). To stop this script matching itself, those
# needles are built from concatenated fragments so the exact literal never
# appears verbatim in this file — including in messages.
#
# BUILD-DEPENDENT checks (anything needing a rendered _site) SKIP when _site is
# absent and are delegated to CI + the T27 live suite, exactly as in Sprint 1.
set -uo pipefail

cd "$(git rev-parse --show-toplevel 2>/dev/null || echo .)" || exit 1

# ---- output helpers ---------------------------------------------------------
PASS()  { echo "PASS $1"; }
FAILM() { echo "FAIL $1: $2"; }
SKIPM() { echo "SKIP $1 ($2)"; }

SITE="_site"
SPEC="spec/index.md"
CLI="cli.md"

# Robust tree-grep exclusion (works whether grep prefixes paths with ./ or not),
# also drops this scripts/ dir so verbatim-needle checks never self-match.
_notree() { grep -Ev '(^|/)\.git/' | grep -Ev '(^|/)docs/agents/' | grep -Ev '(^|/)scripts/'; }

# The Sprint-1 healthcheck line every example carries (verbatim).
HEALTHCHECK="CMD /mnt/tools/file_read --agentrc-schema"

# T26 demo narrative — fragmented to dodge self-match on the tree-wide grep.
NARR="Same artifact, same labels, three ""substrates. The translators are the proof of concept; the labels are the standard."

# Allowed org.agentrc.* label namespaces (no invented fourth — §0.3).
KNOWN_NS="agent capability identity mcp model network policy skill sop substrate tool"

# List (non-docs/agents, non-scripts) files that contain the demo narrative.
_demo_files() {
  grep -rlF "$NARR" --include='*.md' --include='*.html' . 2>/dev/null | _notree
}

# =============================================================================
# T15 — examples/Agentfile.hooked
# =============================================================================
T15=examples/Agentfile.hooked

t15_hooked_exists() {
  if [ -f "$T15" ]; then PASS t15_hooked_exists; return 0; fi
  FAILM t15_hooked_exists "$T15 absent"; return 1
}

t15_hooked_lints() {
  if [ ! -f "$T15" ]; then FAILM t15_hooked_lints "$T15 absent"; return 1; fi
  if go run ./cmd/agentrc lint "$T15" >/dev/null 2>&1; then
    PASS t15_hooked_lints; return 0
  fi
  FAILM t15_hooked_lints "go run ./cmd/agentrc lint $T15 exited non-zero"; return 1
}

t15_hooked_hook_keys() {
  if [ ! -f "$T15" ]; then FAILM t15_hooked_hook_keys "$T15 absent"; return 1; fi
  local n
  n=$(grep -cE 'POLICY agent\.hooks\.(on_tool_call|pre|post)' "$T15" 2>/dev/null | tr -d ' ')
  if [ "${n:-0}" -ge 3 ]; then PASS t15_hooked_hook_keys; return 0; fi
  FAILM t15_hooked_hook_keys "found ${n:-0} POLICY agent.hooks.* line(s) (expected >= 3)"; return 1
}

t15_hooked_explicit_egress() {
  if [ ! -f "$T15" ]; then FAILM t15_hooked_explicit_egress "$T15 absent"; return 1; fi
  if grep -qE 'POLICY network dns:' "$T15" 2>/dev/null; then
    PASS t15_hooked_explicit_egress; return 0
  fi
  FAILM t15_hooked_explicit_egress "no explicit 'POLICY network dns:' line for contrast"; return 1
}

t15_hooked_source_attribution() {
  if [ ! -f "$T15" ]; then FAILM t15_hooked_source_attribution "$T15 absent"; return 1; fi
  if grep -qE '\.source' "$T15" 2>/dev/null; then
    PASS t15_hooked_source_attribution; return 0
  fi
  FAILM t15_hooked_source_attribution "no .source auto-derivation attribution comment (§8.5)"; return 1
}

t15_hooked_healthcheck() {
  if [ ! -f "$T15" ]; then FAILM t15_hooked_healthcheck "$T15 absent"; return 1; fi
  if grep -qF "$HEALTHCHECK" "$T15" 2>/dev/null; then
    PASS t15_hooked_healthcheck; return 0
  fi
  FAILM t15_hooked_healthcheck "verbatim schema healthcheck line absent"; return 1
}

# =============================================================================
# T16 — examples/Agentfile.delegator
# =============================================================================
T16=examples/Agentfile.delegator

t16_delegator_exists() {
  if [ -f "$T16" ]; then PASS t16_delegator_exists; return 0; fi
  FAILM t16_delegator_exists "$T16 absent"; return 1
}

t16_delegator_lints() {
  if [ ! -f "$T16" ]; then FAILM t16_delegator_lints "$T16 absent"; return 1; fi
  if go run ./cmd/agentrc lint "$T16" >/dev/null 2>&1; then
    PASS t16_delegator_lints; return 0
  fi
  FAILM t16_delegator_lints "go run ./cmd/agentrc lint $T16 exited non-zero"; return 1
}

t16_delegator_subagent_keys() {
  if [ ! -f "$T16" ]; then FAILM t16_delegator_subagent_keys "$T16 absent"; return 1; fi
  local miss=""
  grep -qE 'agent\.sub_agents([^.]|$)' "$T16" 2>/dev/null || miss="$miss agent.sub_agents"
  grep -qE 'agent\.sub_agents\.max' "$T16" 2>/dev/null   || miss="$miss agent.sub_agents.max"
  grep -qE 'agent\.sub_agent_timeout' "$T16" 2>/dev/null || miss="$miss agent.sub_agent_timeout"
  if [ -z "$miss" ]; then PASS t16_delegator_subagent_keys; return 0; fi
  FAILM t16_delegator_subagent_keys "missing sub-agent key(s):$miss"; return 1
}

t16_delegator_healthcheck() {
  if [ ! -f "$T16" ]; then FAILM t16_delegator_healthcheck "$T16 absent"; return 1; fi
  if grep -qF "$HEALTHCHECK" "$T16" 2>/dev/null; then
    PASS t16_delegator_healthcheck; return 0
  fi
  FAILM t16_delegator_healthcheck "verbatim schema healthcheck line absent"; return 1
}

# =============================================================================
# Examples index + lint / version guards
# =============================================================================
t15_16_index_cards() {
  local f=examples/index.md miss=""
  grep -q 'Agentfile.hooked' "$f" 2>/dev/null    || miss="$miss hooked"
  grep -q 'Agentfile.delegator' "$f" 2>/dev/null || miss="$miss delegator"
  if [ -z "$miss" ]; then PASS t15_16_index_cards; return 0; fi
  FAILM t15_16_index_cards "examples/index.md missing card link(s):$miss"; return 1
}

v5_example_lint() {
  local fail="" f
  for f in examples/Agentfile.*; do
    [ -e "$f" ] || continue
    go run ./cmd/agentrc lint "$f" >/dev/null 2>&1 || fail="$fail $f"
  done
  if [ -z "$fail" ]; then PASS v5_example_lint; return 0; fi
  FAILM v5_example_lint "lint failed for:$fail"; return 1
}

# unique draft.N values across *.md, excluding .git, docs/agents, CHANGELOG.
_draft_vals() {
  grep -rlE 'draft\.[0-9]+' --include='*.md' . 2>/dev/null \
    | grep -Ev '(^|/)docs/agents/' | grep -v 'CHANGELOG.md' \
    | tr '\n' '\0' | xargs -0 grep -hoE 'draft\.[0-9]+' 2>/dev/null | sort -u
}

v3_version_draft5_guard() {
  local vals count; vals=$(_draft_vals); count=$(echo "$vals" | grep -c .)
  if [ "$count" -eq 1 ] && echo "$vals" | grep -q '^draft\.5$'; then
    PASS v3_version_draft5_guard; return 0
  fi
  FAILM v3_version_draft5_guard "expected exactly draft.5; found: $(echo "$vals" | tr '\n' ' ')"; return 1
}

# =============================================================================
# T17 — §8.7 substrate.<platform>.*
# =============================================================================
t17_substrate_platform_section() {
  if grep -qE '^### 8\.7' "$SPEC" 2>/dev/null \
     && grep -qE 'substrate\.<platform>' "$SPEC" 2>/dev/null; then
    PASS t17_substrate_platform_section; return 0
  fi
  FAILM t17_substrate_platform_section "### 8.7 substrate.<platform> heading absent from $SPEC"; return 1
}

t17_aws_key_registry() {
  local miss="" k
  for k in roleArn networkMode securityGroup subnet protocol maxLifetime 'deployment\.mode' 'code\.s3\.uri'; do
    grep -qE "$k" "$SPEC" 2>/dev/null || miss="$miss ${k//\\/}"
  done
  if [ -z "$miss" ]; then PASS t17_aws_key_registry; return 0; fi
  FAILM t17_aws_key_registry "AWS key(s) not documented:$miss"; return 1
}

t17_unknown_tokens_parse() {
  if grep -qiE 'unknown[^.]*token' "$SPEC" 2>/dev/null \
     && grep -qiE 'MUST parse|ignored|never error' "$SPEC" 2>/dev/null; then
    PASS t17_unknown_tokens_parse; return 0
  fi
  FAILM t17_unknown_tokens_parse "prose on unknown-token parse / foreign-key ignore absent"; return 1
}

t17_no_collision_85_86() {
  local c85 c86
  c85=$(grep -cE '^### 8\.5' "$SPEC" 2>/dev/null | tr -d ' ')
  c86=$(grep -cE '^### 8\.6' "$SPEC" 2>/dev/null | tr -d ' ')
  if [ "${c85:-0}" -eq 1 ] && [ "${c86:-0}" -eq 1 ] \
     && grep -qE '^### 8\.5.*egress' "$SPEC" 2>/dev/null \
     && grep -qE '^### 8\.6.*POLICY' "$SPEC" 2>/dev/null; then
    PASS t17_no_collision_85_86; return 0
  fi
  FAILM t17_no_collision_85_86 "existing ### 8.5 (egress) / ### 8.6 (why POLICY) not intact (got 8.5=$c85 8.6=$c86)"; return 1
}

# =============================================================================
# T18 — §8.8 agent.auth.*
# =============================================================================
t18_agent_auth_section() {
  if grep -qE '^### 8\.8' "$SPEC" 2>/dev/null \
     && grep -qE 'agent\.auth' "$SPEC" 2>/dev/null; then
    PASS t18_agent_auth_section; return 0
  fi
  FAILM t18_agent_auth_section "### 8.8 agent.auth heading absent from $SPEC"; return 1
}

t18_auth_modes_and_jwt_keys() {
  local miss=""
  grep -qE 'platform\|jwt\|none|`platform`.*`jwt`.*`none`' "$SPEC" 2>/dev/null || miss="$miss modes"
  grep -qE 'jwt\.discovery_url' "$SPEC" 2>/dev/null || miss="$miss jwt.discovery_url"
  grep -qE 'allowed_audience' "$SPEC" 2>/dev/null   || miss="$miss allowed_audience"
  grep -qE 'allowed_client' "$SPEC" 2>/dev/null     || miss="$miss allowed_client"
  if [ -z "$miss" ]; then PASS t18_auth_modes_and_jwt_keys; return 0; fi
  FAILM t18_auth_modes_and_jwt_keys "auth mode/jwt key(s) not documented:$miss"; return 1
}

t18_auth_failclosed_prose() {
  if grep -qF 'MUST NOT expose the invocation endpoint' "$SPEC" 2>/dev/null; then
    PASS t18_auth_failclosed_prose; return 0
  fi
  FAILM t18_auth_failclosed_prose "fail-closed 'MUST NOT expose the invocation endpoint' sentence absent"; return 1
}

# =============================================================================
# T19 — §8.9 substrate.runtime.language
# =============================================================================
t19_runtime_language_section() {
  if grep -qE '^### 8\.9' "$SPEC" 2>/dev/null \
     && grep -qE 'substrate\.runtime\.language' "$SPEC" 2>/dev/null; then
    PASS t19_runtime_language_section; return 0
  fi
  FAILM t19_runtime_language_section "### 8.9 substrate.runtime.language heading absent from $SPEC"; return 1
}

t19_codemode_failclosed() {
  if grep -qiE 'code[- ]?mode' "$SPEC" 2>/dev/null \
     && grep -qiE 'resolvable inference|fail closed|fail-closed' "$SPEC" 2>/dev/null \
     && grep -qiE 'container[- ]?mode MAY ignore|container[- ]?mode.*ignore' "$SPEC" 2>/dev/null; then
    PASS t19_codemode_failclosed; return 0
  fi
  FAILM t19_codemode_failclosed "code-mode fail-closed / container-mode-MAY-ignore prose absent"; return 1
}

# =============================================================================
# T20 — supporting edits + version bump + T8 lockfile subsection
# =============================================================================
t20_version_draft6_sitewide() {
  local vals count; vals=$(_draft_vals); count=$(echo "$vals" | grep -c .)
  if [ "$count" -eq 1 ] && echo "$vals" | grep -q '^draft\.6$'; then
    PASS t20_version_draft6_sitewide; return 0
  fi
  FAILM t20_version_draft6_sitewide "expected exactly draft.6 tree-wide; found: $(echo "$vals" | tr '\n' ' ')"; return 1
}

t20_lockgo_bumped() {
  local hits
  hits=$(grep -rn '0\.1\.0-draft\.5' --include='lock.go' . 2>/dev/null | grep -Ev '(^|/)\.git/')
  if [ -z "$hits" ]; then PASS t20_lockgo_bumped; return 0; fi
  FAILM t20_lockgo_bumped "lock.go still hardcodes 0.1.0-draft.5: $(echo "$hits" | tr '\n' ' ')"; return 1
}

t8a_spec_lock_subsection() {
  local ok=1
  grep -qiE 'Reproducible builds' "$SPEC" 2>/dev/null || ok=0
  grep -qF 'agentrc.lock' "$SPEC" 2>/dev/null || ok=0
  grep -qiE 'format TODO' "$SPEC" 2>/dev/null || ok=0
  # slogan kept on the homepage
  grep -qF 'The lockfile pins dependencies' index.md 2>/dev/null || ok=0
  if [ "$ok" -eq 1 ]; then PASS t8a_spec_lock_subsection; return 0; fi
  FAILM t8a_spec_lock_subsection "informative agentrc.lock §9 subsection (Reproducible builds / format TODO) or homepage slogan missing"; return 1
}

t20_docs_agentfile_platform_para() {
  local f=docs/agentfile.md
  if grep -qE 'substrate\.<platform>|platform-scoped' "$f" 2>/dev/null \
     && grep -qE 'agent\.auth\.jwt' "$f" 2>/dev/null; then
    PASS t20_docs_agentfile_platform_para; return 0
  fi
  FAILM t20_docs_agentfile_platform_para "$f lacks platform-scoped paragraph and/or agent.auth.jwt.* example"; return 1
}

t20_code_reviewer_commented_block() {
  local f=examples/Agentfile.code-reviewer
  if [ ! -f "$f" ]; then FAILM t20_code_reviewer_commented_block "$f absent"; return 1; fi
  if grep -qE '^#.*substrate\.aws' "$f" 2>/dev/null \
     && grep -qE '^#.*agent\.auth' "$f" 2>/dev/null \
     && go run ./cmd/agentrc lint "$f" >/dev/null 2>&1; then
    PASS t20_code_reviewer_commented_block; return 0
  fi
  FAILM t20_code_reviewer_commented_block "commented substrate.aws.* + agent.auth.* block absent or file no longer lints"; return 1
}

t20_changelog_draft6() {
  if grep -qE '0\.1\.0-draft\.6' CHANGELOG.md 2>/dev/null; then
    PASS t20_changelog_draft6; return 0
  fi
  FAILM t20_changelog_draft6 "no 0.1.0-draft.6 entry in CHANGELOG.md"; return 1
}

# =============================================================================
# §V suite guards (version / terminology / keywords / namespaces / backends)
# =============================================================================
v3_version_draft6() {
  local vals count; vals=$(_draft_vals); count=$(echo "$vals" | grep -c .)
  if [ "$count" -eq 1 ] && echo "$vals" | grep -q '^draft\.6$'; then
    PASS v3_version_draft6; return 0
  fi
  FAILM v3_version_draft6 "expected single draft.6; found: $(echo "$vals" | tr '\n' ' ')"; return 1
}

# §V.8: --substrate CLI flag purged from cmd/ docs/ cli.md AND spec POLICY
# namespace intact. Shared guard across S2-03 / S2-05 / S2-07.
v8_terminology_split() {
  local n intact=1
  # Exclude *_test.go: tests legitimately reference the removed flag string to
  # assert its ABSENCE from --help; that is not a surviving flag usage.
  n=$( { grep -rl -- '--substrate' cmd/ docs/ "$CLI" 2>/dev/null || true; } \
        | grep -Ev '(^|/)docs/agents/' | grep -Ev '_test\.go$' | grep -c . )
  grep -qE 'POLICY substrate\.' "$SPEC" 2>/dev/null || intact=0
  if [ "${n:-0}" -eq 0 ] && [ "$intact" -eq 1 ]; then
    PASS v8_terminology_split; return 0
  fi
  FAILM v8_terminology_split "--substrate present in $n file(s) (expected 0); spec POLICY-substrate intact=$intact"; return 1
}

v_keyword_count() {
  local miss=""
  grep -qF 'four new keywords' "$SPEC" 2>/dev/null || miss="$miss four-new-keywords-phrase"
  local kw
  for kw in IDENTITY CAPABILITY SOP POLICY; do
    grep -qF "$kw" "$SPEC" 2>/dev/null || miss="$miss $kw"
  done
  # no fifth keyword ever introduced (SECRET/CRED/MOUNT/MEMORY/RATELIMIT as a keyword)
  if grep -qiE 'add(s|ing)? (a )?`?(SECRET|CRED)`? keyword' "$SPEC" 2>/dev/null; then
    miss="$miss fifth-keyword-introduced"
  fi
  if [ -z "$miss" ]; then PASS v_keyword_count; return 0; fi
  FAILM v_keyword_count "keyword invariant broken:$miss"; return 1
}

# §0.3: translator code emits only known org.agentrc.* namespaces.
v_no_fourth_namespace() {
  local bad="" ns
  for ns in $(grep -rhoE 'org\.agentrc\.[a-z_]+' cmd/ internal/ 2>/dev/null \
              | sed -E 's/^org\.agentrc\.//' | sort -u); do
    echo " $KNOWN_NS " | grep -q " $ns " || bad="$bad $ns"
  done
  if [ -z "$bad" ]; then PASS v_no_fourth_namespace; return 0; fi
  FAILM v_no_fourth_namespace "unknown org.agentrc namespace(s) emitted:$bad"; return 1
}

# §V.9: bedrock JSON + kubernetes YAML dry-runs parse.
v9_backend_dryruns() {
  local help
  help=$(go run ./cmd/agentrc run --help 2>&1)
  if ! echo "$help" | grep -q -- '--backend'; then
    FAILM v9_backend_dryruns "arc run --help does not advertise --backend"; return 1
  fi
  local ref="ghcr.io/agentrc/code-reviewer:1.0" jok=0 yok=0
  if go run ./cmd/agentrc run "$ref" --backend bedrock --dry-run 2>/dev/null \
       | python3 -m json.tool >/dev/null 2>&1; then jok=1; fi
  if go run ./cmd/agentrc run "$ref" --backend kubernetes --dry-run 2>/dev/null \
       | python3 -c 'import sys,yaml;list(yaml.safe_load_all(sys.stdin))' >/dev/null 2>&1; then yok=1; fi
  if [ "$jok" -eq 1 ] && [ "$yok" -eq 1 ]; then PASS v9_backend_dryruns; return 0; fi
  FAILM v9_backend_dryruns "bedrock-json parse=$jok kubernetes-yaml parse=$yok (expected 1/1)"; return 1
}

# =============================================================================
# T25 — CLI docs table
# =============================================================================
POS_LINE='Reference translators — a proof of concept until platforms read `org.agentrc.*` labels natively. Not production runners.'

t25_run_implemented() {
  # the `run` row status reads "implemented" with a "reference translators" qualifier
  if grep -E '`?(agentrc|arc) run`?' "$CLI" 2>/dev/null \
       | grep -qiE 'implemented.*reference translators|reference translators.*implemented'; then
    PASS t25_run_implemented; return 0
  fi
  FAILM t25_run_implemented "run row not 'implemented (reference translators)' in $CLI"; return 1
}

t25_positioning_line_verbatim() {
  if grep -qF "$POS_LINE" "$CLI" 2>/dev/null; then
    PASS t25_positioning_line_verbatim; return 0
  fi
  FAILM t25_positioning_line_verbatim "§0.8 positioning line absent (verbatim) from $CLI"; return 1
}

t25_sign_verify_stay_planned() {
  local ok=1
  grep -E '`?(agentrc|arc) sign`?' "$CLI" 2>/dev/null   | grep -q 'planned' || ok=0
  grep -E '`?(agentrc|arc) verify`?' "$CLI" 2>/dev/null | grep -q 'planned' || ok=0
  if [ "$ok" -eq 1 ]; then PASS t25_sign_verify_stay_planned; return 0; fi
  FAILM t25_sign_verify_stay_planned "sign/verify rows no longer read 'planned'"; return 1
}

t25_no_substrate_in_cli() {
  local n
  n=$(grep -c -- '--substrate' "$CLI" 2>/dev/null | tr -d ' ')
  if [ "${n:-0}" -eq 0 ]; then PASS t25_no_substrate_in_cli; return 0; fi
  FAILM t25_no_substrate_in_cli "$n '--substrate' occurrence(s) remain in $CLI (expected 0)"; return 1
}

# =============================================================================
# T26 — one agent, three backends (demo)
# =============================================================================
t26_demo_narrative_verbatim() {
  local n
  n=$(grep -rlF "$NARR" --include='*.md' --include='*.html' . 2>/dev/null | _notree | wc -l | tr -d ' ')
  if [ "${n:-0}" -eq 1 ]; then PASS t26_demo_narrative_verbatim; return 0; fi
  FAILM t26_demo_narrative_verbatim "demo narrative present in ${n:-0} served file(s) (expected exactly 1)"; return 1
}

t26_three_backend_commands() {
  local f miss=""
  f=$(_demo_files | head -1)
  if [ -z "$f" ]; then FAILM t26_three_backend_commands "demo doc absent"; return 1; fi
  grep -qF 'arc build -t ghcr.io/agentrc/code-reviewer:1.0' "$f" 2>/dev/null || miss="$miss build"
  grep -qF -- '--backend local --isolation microvm' "$f" 2>/dev/null || miss="$miss local"
  grep -qF -- '--backend bedrock --dry-run' "$f" 2>/dev/null || miss="$miss bedrock"
  grep -qF -- '--backend kubernetes --dry-run' "$f" 2>/dev/null || miss="$miss kubernetes"
  if [ -z "$miss" ]; then PASS t26_three_backend_commands; return 0; fi
  FAILM t26_three_backend_commands "demo doc missing command(s):$miss"; return 1
}

t26_no_dropped_backends() {
  local f
  f=$(_demo_files | head -1)
  if [ -z "$f" ]; then FAILM t26_no_dropped_backends "demo doc absent"; return 1; fi
  local bad=""
  grep -q -- '--substrate' "$f" 2>/dev/null && bad="$bad --substrate"
  grep -qE -- '--backend (gcp|compose)' "$f" 2>/dev/null && bad="$bad dropped-backend"
  if [ -z "$bad" ]; then PASS t26_no_dropped_backends; return 0; fi
  FAILM t26_no_dropped_backends "demo doc references dropped surface(s):$bad"; return 1
}

# =============================================================================
# §V — artifact leak guard (M-004)
# =============================================================================
v_no_artifact_leak() {
  if [ -d "$SITE" ]; then
    local n
    n=$(grep -rl 'docs/agents' "$SITE" 2>/dev/null | wc -l | tr -d ' ')
    if [ "${n:-0}" -eq 0 ]; then PASS v_no_artifact_leak; return 0; fi
    FAILM v_no_artifact_leak "$n built page(s) leak docs/agents (expected 0)"; return 1
  fi
  # no _site: assert the source config excludes docs/agents from the build.
  if grep -qE '^\s*-\s*docs/agents\s*$' _config.yml 2>/dev/null; then
    PASS v_no_artifact_leak; return 0
  fi
  FAILM v_no_artifact_leak "_config.yml does not exclude docs/agents (and no _site to verify)"; return 1
}

# =============================================================================
# Runner
# =============================================================================
ALL_CHECKS=(
  t15_hooked_exists t15_hooked_lints t15_hooked_hook_keys
  t15_hooked_explicit_egress t15_hooked_source_attribution t15_hooked_healthcheck
  t16_delegator_exists t16_delegator_lints t16_delegator_subagent_keys t16_delegator_healthcheck
  t15_16_index_cards v5_example_lint v3_version_draft5_guard
  t17_substrate_platform_section t17_aws_key_registry t17_unknown_tokens_parse t17_no_collision_85_86
  t18_agent_auth_section t18_auth_modes_and_jwt_keys t18_auth_failclosed_prose
  t19_runtime_language_section t19_codemode_failclosed
  t20_version_draft6_sitewide t20_lockgo_bumped t8a_spec_lock_subsection
  t20_docs_agentfile_platform_para t20_code_reviewer_commented_block t20_changelog_draft6
  v3_version_draft6 v8_terminology_split v_keyword_count
  v_no_fourth_namespace v9_backend_dryruns
  t25_run_implemented t25_positioning_line_verbatim t25_sign_verify_stay_planned t25_no_substrate_in_cli
  t26_demo_narrative_verbatim t26_three_backend_commands t26_no_dropped_backends
  v_no_artifact_leak
)

# Local-applicable subset used by the v_all_local aggregator (S2-07 pre-deploy).
LOCAL_CHECKS=("${ALL_CHECKS[@]}")

_run_set() {
  local -n checks=$1
  local p=0 fdone=0 s=0 fn rc
  for fn in "${checks[@]}"; do
    "$fn"; rc=$?
    case "$rc" in
      0) p=$((p+1));;
      2) s=$((s+1));;
      *) fdone=$((fdone+1));;
    esac
  done
  echo "----------------------------------------"
  echo "$p passed, $fdone failed, $s skipped"
  [ "$fdone" -eq 0 ]
}

v_all_local() {
  echo "== v_all_local: full Sprint-2 local §V suite =="
  _run_set LOCAL_CHECKS
}

main() {
  _run_set ALL_CHECKS
}

# If sourced, expose functions; if executed, run requested target(s).
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
  if [ "$#" -eq 0 ]; then
    main
  else
    for fn in "$@"; do "$fn"; done
  fi
fi
