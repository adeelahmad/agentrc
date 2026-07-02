#!/usr/bin/env bash
# Sprint 1 §V LIVE verification (T14). Runs curl checks against the deployed site.
# Usage: scripts/verify-sprint1-live.sh [BASE_URL]   (default https://agentrc.ai)
# These are meant to run only AFTER deploy (T14). Running pre-deploy will fail/timeout.
# Return codes per check: 0 = PASS, 1 = FAIL, 2 = SKIP.
set -uo pipefail

cd "$(git rev-parse --show-toplevel 2>/dev/null || echo .)" || exit 1

BASE_URL="${1:-https://agentrc.ai}"
BASE_URL="${BASE_URL%/}"

PASS()  { echo "PASS $1"; }
FAILM() { echo "FAIL $1: $2"; }
SKIPM() { echo "SKIP $1 ($2)"; }

# Concatenated needles so a local tree-grep of this file does not self-match.
GHOST="tools/""ping"
RUNNER_A="Runner ""Conformance"
RUNNER_B="Runner ""conformance"

CURL="curl -sS --max-time 20"

_fetch() { $CURL "$1" 2>/dev/null; }

t14_live_ghost() {
  local body n
  body=$( { _fetch "$BASE_URL/"; _fetch "$BASE_URL/examples/"; } )
  n=$(printf '%s' "$body" | grep -c "$GHOST")
  if [ "${n:-0}" -eq 0 ]; then PASS t14_live_ghost; return 0; fi
  FAILM t14_live_ghost "$n ghost-tool occurrence(s) on live homepage/examples"; return 1
}

t14_live_sitemap_notes() {
  local n
  n=$(_fetch "$BASE_URL/sitemap.xml" | grep -c 'CURRENT_IMPLEMENTATION_MAPPING')
  if [ "${n:-0}" -eq 0 ]; then PASS t14_live_sitemap_notes; return 0; fi
  FAILM t14_live_sitemap_notes "$n notes URL(s) in live sitemap"; return 1
}

t14_live_sitemap_workflow() {
  local n
  n=$(_fetch "$BASE_URL/sitemap.xml" | grep -c 'workflow')
  if [ "${n:-0}" -eq 0 ]; then PASS t14_live_sitemap_workflow; return 0; fi
  FAILM t14_live_sitemap_workflow "$n workflow URL(s) in live sitemap"; return 1
}

t14_live_workflow_404() {
  local code
  code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 20 "$BASE_URL/docs/workflows/" 2>/dev/null)
  if [ "$code" = "404" ] || [ "$code" = "410" ]; then PASS t14_live_workflow_404; return 0; fi
  # A published-but-noindex mechanism is also acceptable per T12; check meta robots.
  if _fetch "$BASE_URL/docs/workflows/" | grep -qi 'name="robots"[^>]*noindex'; then
    PASS t14_live_workflow_404; return 0
  fi
  FAILM t14_live_workflow_404 "/docs/workflows/ returned $code and is not noindex"; return 1
}

t14_live_runner_label() {
  local n
  n=$(_fetch "$BASE_URL/profiles/" | grep -cE "$RUNNER_A|$RUNNER_B")
  if [ "${n:-0}" -eq 0 ]; then PASS t14_live_runner_label; return 0; fi
  FAILM t14_live_runner_label "$n 'Runner Conformance' label(s) on live /profiles/"; return 1
}

t14_all_tasks_recorded() {
  local f=docs/agents/sprint1/live-verification.md
  if [ ! -f "$f" ]; then SKIPM t14_all_tasks_recorded "$f absent — record T1-T12 live pass/fail there"; return 2; fi
  local missing="" t
  for t in T1 T2 T3 T4 T5 T6 T7 T8 T9 T10 T11 T12; do
    grep -qE "\b${t}\b" "$f" || missing="$missing $t"
  done
  if [ -z "$missing" ]; then PASS t14_all_tasks_recorded; return 0; fi
  FAILM t14_all_tasks_recorded "tasks not recorded in $f:$missing"; return 1
}

ALL_LIVE=(
  t14_live_ghost t14_live_sitemap_notes t14_live_sitemap_workflow
  t14_live_workflow_404 t14_live_runner_label t14_all_tasks_recorded
)

main() {
  echo "== live verification against $BASE_URL =="
  local p=0 fdone=0 s=0 fn rc
  for fn in "${ALL_LIVE[@]}"; do
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
  if [ "$#" -le 1 ]; then
    main
  else
    shift
    for fn in "$@"; do "$fn"; done
  fi
fi
