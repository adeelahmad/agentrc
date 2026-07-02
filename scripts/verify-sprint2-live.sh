#!/usr/bin/env bash
# Sprint 2 live re-verification (POST-DEPLOY). Curl checks against the deployed
# site. RED phase: written before the sprint lands, so every live assertion
# FAILS until https://agentrc.ai is deployed at draft.6.
# Return codes: 0 = PASS, 1 = FAIL, 2 = SKIP.
#
# BASE_URL defaults to https://agentrc.ai; override to test a preview deploy.
set -uo pipefail

BASE_URL="${BASE_URL:-https://agentrc.ai}"

PASS()  { echo "PASS $1"; }
FAILM() { echo "FAIL $1: $2"; }
SKIPM() { echo "SKIP $1 ($2)"; }

# Fetch a URL; echo body on success, non-zero on transport failure.
_get() {
  if ! command -v curl >/dev/null 2>&1; then return 3; fi
  curl -fsSL --max-time 20 "$1" 2>/dev/null
}

# =============================================================================
# T27 — live aggregate re-verification
# =============================================================================
t27_live() {
  if ! command -v curl >/dev/null 2>&1; then
    SKIPM t27_live "curl unavailable"; return 2
  fi
  local home cli fails=""

  home=$(_get "$BASE_URL/") || { FAILM t27_live "$BASE_URL unreachable"; return 1; }

  # version = draft.6 sitewide (homepage + spec)
  echo "$home" | grep -q 'draft.6' || fails="$fails home-not-draft6"
  local spec; spec=$(_get "$BASE_URL/spec/") || spec=""
  echo "$spec" | grep -q 'draft.6' || fails="$fails spec-not-draft6"

  # §8.7 / §8.8 / §8.9 rendered in the spec
  echo "$spec" | grep -q '8.7' || fails="$fails no-8.7"
  echo "$spec" | grep -q '8.8' || fails="$fails no-8.8"
  echo "$spec" | grep -q '8.9' || fails="$fails no-8.9"

  # /cli/ shows --backend and no --substrate (token-split aware, M-006)
  cli=$(_get "$BASE_URL/cli/") || cli=""
  echo "$cli" | grep -q -- '--backend'  || fails="$fails cli-no-backend"
  echo "$cli" | grep -q -- '--substrate' && fails="$fails cli-has-substrate"

  # the two new examples reachable + rendered
  _get "$BASE_URL/examples/Agentfile.hooked"    >/dev/null || fails="$fails hooked-unreachable"
  _get "$BASE_URL/examples/Agentfile.delegator" >/dev/null || fails="$fails delegator-unreachable"

  # arc lint of the live examples needs the local binary + fetched file: delegated.
  echo "NOTE t27_live: live 'arc lint' of served examples delegated to a local file record (fetch + go run ./cmd/agentrc lint)."

  if [ -z "$fails" ]; then PASS t27_live; return 0; fi
  FAILM t27_live "live checks failed:$fails"; return 1
}

t27_no_artifact_leak_live() {
  if ! command -v curl >/dev/null 2>&1; then
    SKIPM t27_no_artifact_leak_live "curl unavailable"; return 2
  fi
  local sm n
  sm=$(_get "$BASE_URL/sitemap.xml") || { FAILM t27_no_artifact_leak_live "$BASE_URL/sitemap.xml unreachable"; return 1; }
  n=$(echo "$sm" | grep -c 'docs/agents' | tr -d ' ')
  if [ "${n:-0}" -eq 0 ]; then PASS t27_no_artifact_leak_live; return 0; fi
  FAILM t27_no_artifact_leak_live "$n docs/agents URL(s) in live sitemap (expected 0)"; return 1
}

# =============================================================================
# Runner
# =============================================================================
ALL_CHECKS=(t27_live t27_no_artifact_leak_live)

main() {
  local p=0 fdone=0 s=0 fn rc
  for fn in "${ALL_CHECKS[@]}"; do
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

if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
  if [ "$#" -eq 0 ]; then main; else for fn in "$@"; do "$fn"; done; fi
fi
