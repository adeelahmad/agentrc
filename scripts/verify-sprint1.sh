#!/usr/bin/env bash
# Sprint 1 §V verification suite (LOCAL). RED phase: written before the fixes.
# Each function asserts the DESIRED post-fix end-state and prints PASS/FAIL/SKIP.
# Return codes per check: 0 = PASS, 1 = FAIL, 2 = SKIP.
#
# NOTE ON SELF-MATCH: several checks grep the whole tree for a literal (e.g. the
# ghost tool path, the old profiles conformance label). To avoid this script
# matching itself, those needles are built from concatenated string fragments so
# the exact literal never appears verbatim in this file — including in messages.
set -uo pipefail

cd "$(git rev-parse --show-toplevel 2>/dev/null || echo .)" || exit 1

# ---- output helpers ---------------------------------------------------------
PASS()  { echo "PASS $1"; }
FAILM() { echo "FAIL $1: $2"; }
SKIPM() { echo "SKIP $1 ($2)"; }

# ---- shared needles (concatenated to dodge self-match on tree-wide greps) ----
GHOST="tools/""ping"                 # the ghost healthcheck tool path
BRAND="Minimal AgentRC"" agent"      # the wrong brand-casing string
HELLO="IDENTITY name=""hello"        # canonical hello identity line
RUNNER_A="Runner ""Conformance"
RUNNER_B="Runner ""conformance"

SITE="_site"

# The six T1 surfaces that carry a healthcheck (spec has two). secure-workspace is
# already correct and is NOT a target surface, so it is excluded here.
T1_SURFACES=(index.md spec/index.md examples/index.md examples/Agentfile.minimal examples/Agentfile.code-reviewer)

# Robust tree-grep exclusion (works whether grep prefixes paths with ./ or not).
_notree() { grep -Ev '(^|/)\.git/' | grep -Ev '(^|/)docs/agents/'; }

# ---- shared helpers ---------------------------------------------------------
# List (non-docs/agents) md/html files that render the hello agent.
hello_files() {
  grep -rln "$HELLO" --include='*.md' --include='*.html' . 2>/dev/null \
    | grep -v '/docs/agents/' | grep -v '/.git/'
}

# Normalise instruction lines: strip comments, blanks, trim, collapse spaces.
norm_instr() {
  sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' \
    | grep -v '^#' | grep -v '^$' \
    | sed -e 's/[[:space:]]\{1,\}/ /g'
}

# =============================================================================
# T1 — ghost tool purge
# =============================================================================
t1_no_ghost_ping() {
  local n
  n=$(grep -rn "$GHOST" . 2>/dev/null | _notree | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t1_no_ghost_ping; return 0; fi
  FAILM t1_no_ghost_ping "$n ghost-tool occurrence(s) remain (expected 0)"; return 1
}

t1_schema_healthcheck_present() {
  local n f present=()
  for f in "${T1_SURFACES[@]}"; do [ -e "$f" ] && present+=("$f"); done
  n=$(grep -h "file_read --agentrc-schema" "${present[@]}" 2>/dev/null | wc -l | tr -d ' ')
  if [ "$n" -eq 6 ]; then PASS t1_schema_healthcheck_present; return 0; fi
  FAILM t1_schema_healthcheck_present "found $n schema-healthcheck line(s) across the 6 T1 surfaces (expected 6)"; return 1
}

# =============================================================================
# T2 — hello + FROM
# =============================================================================
t2_hello_has_from() {
  local from="FROM python:3.11-slim"
  local bad="" f ln start
  for f in $(hello_files); do
    ln=$(grep -n "$HELLO" "$f" | head -1 | cut -d: -f1)
    [ -z "$ln" ] && continue
    start=$((ln - 3)); [ "$start" -lt 1 ] && start=1
    if ! sed -n "${start},${ln}p" "$f" | grep -q "$from"; then
      bad="$bad $f"
    fi
  done
  if [ -n "$bad" ]; then FAILM t2_hello_has_from "missing '$from' within 3 lines above hello in:$bad"; return 1; fi
  PASS t2_hello_has_from; return 0
}

t2_spec_from_sentence() {
  if grep -q 'every Agentfile MUST contain a `FROM`' spec/index.md 2>/dev/null; then
    PASS t2_spec_from_sentence; return 0
  fi
  FAILM t2_spec_from_sentence "FROM-required MUST sentence absent from spec/index.md"; return 1
}

t2_hello_diff_clean() {
  local canon tmp bad="" f blk got
  canon=$(norm_instr < examples/Agentfile.minimal)
  tmp=$(mktemp -d)
  for f in $(hello_files); do
    awk -v dir="$tmp" '
      /^[[:space:]]*```dockerfile/ { inb=1; n++; fn=dir"/blk_"n; next }
      /^[[:space:]]*```/           { if (inb) { inb=0; close(fn) } next }
      inb                          { print > fn }
    ' "$f"
    for blk in "$tmp"/blk_*; do
      [ -e "$blk" ] || continue
      if grep -q "$HELLO" "$blk"; then
        got=$(norm_instr < "$blk")
        [ "$got" != "$canon" ] && bad="$bad $f"
      fi
      rm -f "$blk"
    done
  done
  rm -rf "$tmp"
  if [ -n "$bad" ]; then FAILM t2_hello_diff_clean "inline hello block drift vs examples/Agentfile.minimal in:$bad"; return 1; fi
  PASS t2_hello_diff_clean; return 0
}

# =============================================================================
# T3 — quickstart honesty callouts
# =============================================================================
t3_step2_callout() {
  # open-q1 resolved: image is PUBLISHED. Accept EITHER the not-yet-published
  # callout OR the published-form `# syntax=ghcr.io/adeelahmad/agentrc-frontend`.
  local f=docs/quickstart.md
  if grep -q 'not yet published' "$f" 2>/dev/null \
     || grep -q 'syntax=ghcr.io/adeelahmad/agentrc-frontend' "$f" 2>/dev/null; then
    PASS t3_step2_callout; return 0
  fi
  FAILM t3_step2_callout "neither not-yet-published nor published-image callout present in $f"; return 1
}

t3_step5_callout() {
  local f=docs/quickstart.md
  if grep -q 'arc run' "$f" 2>/dev/null && grep -q 'is planned' "$f" 2>/dev/null; then
    PASS t3_step5_callout; return 0
  fi
  FAILM t3_step5_callout "'arc run is planned' callout absent in $f"; return 1
}

t3_no_present_tense_claims() {
  local f=docs/quickstart.md n
  n=$(grep -cE 'nothing to install|no extra tooling' "$f" 2>/dev/null | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t3_no_present_tense_claims; return 0; fi
  FAILM t3_no_present_tense_claims "$n present-tense install/tooling guarantee(s) remain in $f"; return 1
}

# =============================================================================
# T4 — polish
# =============================================================================
t4_brand_casing() {
  local n
  n=$(grep -rn "$BRAND" . 2>/dev/null | _notree | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t4_brand_casing; return 0; fi
  FAILM t4_brand_casing "$n wrong-casing brand string(s) remain (expected 0)"; return 1
}

t4_quickstart_hello_canonical() {
  local f=docs/quickstart.md
  if grep -q "$HELLO version=0.1 author=acme" "$f" 2>/dev/null \
     && grep -q 'claude-sonnet-4' "$f" 2>/dev/null; then
    PASS t4_quickstart_hello_canonical; return 0
  fi
  FAILM t4_quickstart_hello_canonical "quickstart hello not aligned to canonical (version=0.1 author=acme, model.name claude-sonnet-4)"; return 1
}

# =============================================================================
# T5 — sidebar + docs-index coverage
# =============================================================================
_sidebar_has() { grep -q "$1" _layouts/doc.html 2>/dev/null; }

t5_sidebar_quickstart() {
  if _sidebar_has '/docs/quickstart/'; then PASS t5_sidebar_quickstart; return 0; fi
  FAILM t5_sidebar_quickstart "_layouts/doc.html has no /docs/quickstart/ link"; return 1
}

t5_sidebar_conformance() {
  if _sidebar_has '/docs/conformance/'; then PASS t5_sidebar_conformance; return 0; fi
  FAILM t5_sidebar_conformance "_layouts/doc.html has no /docs/conformance/ link"; return 1
}

t5_sidebar_impl_mapping() {
  if _sidebar_has '/docs/implementation-mapping/'; then PASS t5_sidebar_impl_mapping; return 0; fi
  FAILM t5_sidebar_impl_mapping "_layouts/doc.html has no /docs/implementation-mapping/ link"; return 1
}

t5_sidebar_no_workflows() {
  local n
  n=$(grep -c '/docs/workflows/' _layouts/doc.html 2>/dev/null | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t5_sidebar_no_workflows; return 0; fi
  FAILM t5_sidebar_no_workflows "_layouts/doc.html still links /docs/workflows/ ($n)"; return 1
}

t5_docsindex_cards() {
  local f=docs/index.md
  if grep -q '/docs/conformance/' "$f" 2>/dev/null \
     && grep -q '/docs/implementation-mapping/' "$f" 2>/dev/null; then
    PASS t5_docsindex_cards; return 0
  fi
  FAILM t5_docsindex_cards "docs/index.md missing conformance and/or implementation-mapping card"; return 1
}

# =============================================================================
# T6 — de-orphan /tooling/
# =============================================================================
t6_tooling_inbound() {
  local bad=""
  grep -q '/tooling/' cli.md 2>/dev/null || bad="$bad cli.md"
  grep -q '/tooling/' docs/index.md 2>/dev/null || bad="$bad docs/index.md"
  if [ -z "$bad" ]; then PASS t6_tooling_inbound; return 0; fi
  FAILM t6_tooling_inbound "no /tooling/ anchor in:$bad"; return 1
}

# =============================================================================
# T7 — kill duplicate URL (BUILD-DEPENDENT)
# =============================================================================
t7_sitemap_no_notes() {
  if [ ! -d "$SITE" ]; then SKIPM t7_sitemap_no_notes "no _site; delegated to CI build + T14 live"; return 2; fi
  local n
  n=$(grep -c 'CURRENT_IMPLEMENTATION_MAPPING' "$SITE/sitemap.xml" 2>/dev/null | tr -d ' ')
  if [ "${n:-0}" -eq 0 ]; then PASS t7_sitemap_no_notes; return 0; fi
  FAILM t7_sitemap_no_notes "$n notes URL(s) in sitemap (expected 0)"; return 1
}

# =============================================================================
# T9 — Runner → Platform relabel
# =============================================================================
t9_no_runner_conformance_label() {
  local n
  n=$(grep -rnE "$RUNNER_A|$RUNNER_B" . 2>/dev/null | _notree | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t9_no_runner_conformance_label; return 0; fi
  FAILM t9_no_runner_conformance_label "$n visible '$RUNNER_A' label(s) remain (expected 0)"; return 1
}

t9_url_preserved() {
  if grep -q 'permalink: /profiles/runner-conformance/' profiles/runner-conformance.md 2>/dev/null; then
    PASS t9_url_preserved; return 0
  fi
  FAILM t9_url_preserved "permalink /profiles/runner-conformance/ was removed/changed"; return 1
}

# =============================================================================
# T10 — canonical dedupe (BUILD-DEPENDENT)
# =============================================================================
t10_one_canonical() {
  if [ ! -d "$SITE" ]; then SKIPM t10_one_canonical "no _site; delegated to CI build + T14 live"; return 2; fi
  local bad="" p c
  while IFS= read -r p; do
    c=$(grep -c 'rel="canonical"' "$p" 2>/dev/null | tr -d ' ')
    [ "${c:-0}" -eq 1 ] || bad="$bad $p($c)"
  done < <(find "$SITE" -name '*.html')
  if [ -z "$bad" ]; then PASS t10_one_canonical; return 0; fi
  FAILM t10_one_canonical "pages without exactly one canonical:$bad"; return 1
}

# =============================================================================
# T12 — park the workflow draft
# =============================================================================
t12_no_served_workflow_refs() {
  if [ ! -d "$SITE" ]; then SKIPM t12_no_served_workflow_refs "no _site; delegated to CI build + T14 live"; return 2; fi
  local n
  n=$(grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' "$SITE" 2>/dev/null | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t12_no_served_workflow_refs; return 0; fi
  FAILM t12_no_served_workflow_refs "$n served workflow reference(s) remain (expected 0)"; return 1
}

t12_pages_unpublished() {
  if [ ! -d "$SITE" ]; then SKIPM t12_pages_unpublished "no _site; delegated to CI build + T14 live"; return 2; fi
  local bad=""
  [ -e "$SITE/docs/workflows/index.html" ] && bad="$bad docs/workflows"
  [ -e "$SITE/profiles/workflow-draft/index.html" ] && bad="$bad profiles/workflow-draft"
  if [ -z "$bad" ]; then PASS t12_pages_unpublished; return 0; fi
  FAILM t12_pages_unpublished "still-rendered pages:$bad"; return 1
}

t12_yaml_dropped() {
  if [ ! -d "$SITE" ]; then SKIPM t12_yaml_dropped "no _site; delegated to CI build + T14 live"; return 2; fi
  local bad=""
  [ -e "$SITE/examples/agent-workflow.yaml" ] && bad="$bad served-yaml"
  grep -q '/examples/agent-workflow.yaml' examples/index.md 2>/dev/null && bad="$bad examples/index.md-link"
  if [ -z "$bad" ]; then PASS t12_yaml_dropped; return 0; fi
  FAILM t12_yaml_dropped "workflow yaml still present:$bad"; return 1
}

t12_sidebar_unlinked() {
  local n
  n=$(grep -c '/docs/workflows/' _layouts/doc.html 2>/dev/null | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS t12_sidebar_unlinked; return 0; fi
  FAILM t12_sidebar_unlinked "_layouts/doc.html still links /docs/workflows/ ($n)"; return 1
}

t12_llms_swept() {
  local n
  n=$(grep -cE 'docs/workflows/|profiles/workflow-draft/' llms.txt 2>/dev/null | tr -d ' ')
  if [ "${n:-0}" -eq 0 ]; then PASS t12_llms_swept; return 0; fi
  FAILM t12_llms_swept "$n workflow link entr(y/ies) remain in llms.txt (expected 0)"; return 1
}

t12_sitemap_no_workflow() {
  if [ ! -d "$SITE" ]; then SKIPM t12_sitemap_no_workflow "no _site; delegated to CI build + T14 live"; return 2; fi
  local n
  n=$(grep -cE 'workflow-draft|docs/workflows' "$SITE/sitemap.xml" 2>/dev/null | tr -d ' ')
  if [ "${n:-0}" -eq 0 ]; then PASS t12_sitemap_no_workflow; return 0; fi
  FAILM t12_sitemap_no_workflow "$n workflow URL(s) in sitemap (expected 0)"; return 1
}

t12_changelog_line() {
  if grep -q 'Workflow draft parked' CHANGELOG.md 2>/dev/null; then
    PASS t12_changelog_line; return 0
  fi
  FAILM t12_changelog_line "'Workflow draft parked' line absent from CHANGELOG.md"; return 1
}

t12_sources_kept() {
  # Sources may move under a parked/ area — search by basename among tracked files.
  local tracked bad=""
  tracked=$(git ls-files 2>/dev/null)
  echo "$tracked" | grep -q '/workflows\.md$\|^workflows\.md$\|/workflows/' || bad="$bad workflows.md"
  echo "$tracked" | grep -q 'workflow-draft\.md$' || bad="$bad workflow-draft.md"
  echo "$tracked" | grep -q 'agent-workflow\.yaml$' || bad="$bad agent-workflow.yaml"
  if [ -z "$bad" ]; then PASS t12_sources_kept; return 0; fi
  FAILM t12_sources_kept "workflow source(s) no longer tracked:$bad"; return 1
}

t12_no_dangling_links() {
  if [ ! -d "$SITE" ]; then SKIPM t12_no_dangling_links "no _site; delegated to CI build + T14 live"; return 2; fi
  if command -v htmlproofer >/dev/null 2>&1; then
    if htmlproofer "$SITE" --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https >/dev/null 2>&1; then
      PASS t12_no_dangling_links; return 0
    fi
    FAILM t12_no_dangling_links "htmlproofer reported dangling internal links"; return 1
  fi
  SKIPM t12_no_dangling_links "htmlproofer not installed"; return 2
}

# =============================================================================
# T8 — lockfile investigation report (report-only)
# =============================================================================
T8_REPORT=docs/agents/sprint1/lockfile-report.md

t8_report_exists() {
  if [ -f "$T8_REPORT" ]; then PASS t8_report_exists; return 0; fi
  FAILM t8_report_exists "$T8_REPORT not written"; return 1
}

t8_report_documents_output() {
  if [ ! -f "$T8_REPORT" ]; then FAILM t8_report_documents_output "report absent"; return 1; fi
  # filename + format + records facts (case-insensitive)
  if grep -qi 'filename\|\.lock' "$T8_REPORT" \
     && grep -qi 'format' "$T8_REPORT" \
     && grep -qi 'record' "$T8_REPORT"; then
    PASS t8_report_documents_output; return 0
  fi
  FAILM t8_report_documents_output "report does not document lock output (filename/format/records)"; return 1
}

t8_report_documents_build() {
  if [ ! -f "$T8_REPORT" ]; then FAILM t8_report_documents_build "report absent"; return 1; fi
  if grep -qi 'build' "$T8_REPORT" && grep -qi 'consume\|consumes\|reads\|uses' "$T8_REPORT"; then
    PASS t8_report_documents_build; return 0
  fi
  FAILM t8_report_documents_build "report does not document how build consumes the lock"; return 1
}

t8_report_ab_recommendation() {
  if [ ! -f "$T8_REPORT" ]; then FAILM t8_report_ab_recommendation "report absent"; return 1; fi
  if grep -qi 'Option A' "$T8_REPORT" && grep -qi 'Option B' "$T8_REPORT"; then
    PASS t8_report_ab_recommendation; return 0
  fi
  FAILM t8_report_ab_recommendation "report missing Option A / Option B framing"; return 1
}

t8_no_site_edits() {
  # Informational guard: report-only boundary. Not counted as a hard fail.
  local changed
  changed=$(git diff --name-only master...HEAD -- index.md spec/ cli.md 'docs/*.md' 2>/dev/null)
  if [ -z "$changed" ]; then
    PASS t8_no_site_edits; return 0
  fi
  echo "INFO t8_no_site_edits: site files differ from master (expected once S1-01/02/03 land): $(echo "$changed" | tr '\n' ' ')"
  PASS t8_no_site_edits; return 0
}

# =============================================================================
# §V suite gates
# =============================================================================
v3_version_coherence() {
  # exactly one draft.N value across *.md (excl .git, docs/agents, CHANGELOG); == draft.5
  local vals
  vals=$(grep -rlE 'draft\.[0-9]+' --include='*.md' . 2>/dev/null \
    | grep -Ev '(^|/)docs/agents/' | grep -v 'CHANGELOG.md' \
    | tr '\n' '\0' | xargs -0 grep -hoE 'draft\.[0-9]+' 2>/dev/null | sort -u)
  local count; count=$(echo "$vals" | grep -c . )
  if [ "$count" -eq 1 ] && echo "$vals" | grep -q '^draft\.5$'; then
    PASS v3_version_coherence; return 0
  fi
  FAILM v3_version_coherence "expected exactly draft.5; found: $(echo "$vals" | tr '\n' ' ')"; return 1
}

_lint_examples() {
  local fail="" f
  for f in examples/Agentfile.*; do
    [ -e "$f" ] || continue
    if ! go run ./cmd/agentrc lint "$f" >/dev/null 2>&1; then
      fail="$fail $f"
    fi
  done
  echo "$fail"
}

v5_example_lint() {
  local fail; fail=$(_lint_examples)
  if [ -z "$fail" ]; then PASS v5_example_lint; return 0; fi
  FAILM v5_example_lint "lint failed for:$fail"; return 1
}
v_example_lint() { v5_example_lint; }

v1_ghost_tool() {
  # §V check 1 alias — source-level; identical assertion to t1_no_ghost_ping.
  local n
  n=$(grep -rn "$GHOST" . 2>/dev/null | _notree | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS v1_ghost_tool; return 0; fi
  FAILM v1_ghost_tool "$n ghost-tool occurrence(s) remain (expected 0)"; return 1
}

v4_hello_from() {
  # §V check 4 alias — source-level; every hello render carries FROM.
  local from="FROM python:3.11-slim" bad="" f ln start
  for f in $(hello_files); do
    ln=$(grep -n "$HELLO" "$f" | head -1 | cut -d: -f1)
    [ -z "$ln" ] && continue
    start=$((ln - 3)); [ "$start" -lt 1 ] && start=1
    sed -n "${start},${ln}p" "$f" | grep -q "$from" || bad="$bad $f"
  done
  if [ -z "$bad" ]; then PASS v4_hello_from; return 0; fi
  FAILM v4_hello_from "missing FROM near hello in:$bad"; return 1
}

v6_one_canonical() {
  if [ ! -d "$SITE" ]; then SKIPM v6_one_canonical "no _site; delegated to CI build + T14 live"; return 2; fi
  local bad="" p c
  while IFS= read -r p; do
    c=$(grep -c 'rel="canonical"' "$p" 2>/dev/null | tr -d ' ')
    [ "${c:-0}" -eq 1 ] || bad="$bad $p($c)"
  done < <(find "$SITE" -name '*.html')
  if [ -z "$bad" ]; then PASS v6_one_canonical; return 0; fi
  FAILM v6_one_canonical "pages without exactly one canonical:$bad"; return 1
}

v7_internal_links() {
  if [ ! -d "$SITE" ]; then SKIPM v7_internal_links "no _site; delegated to CI build + T14 live"; return 2; fi
  if command -v htmlproofer >/dev/null 2>&1; then
    if htmlproofer "$SITE" --disable-external --allow-hash-href --ignore-empty-alt --no-enforce-https >/dev/null 2>&1; then
      PASS v7_internal_links; return 0
    fi
    FAILM v7_internal_links "htmlproofer reported dangling internal links"; return 1
  fi
  SKIPM v7_internal_links "htmlproofer not installed"; return 2
}

v11_workflow_parked() {
  if [ ! -d "$SITE" ]; then SKIPM v11_workflow_parked "no _site; delegated to CI build + T14 live"; return 2; fi
  local n
  n=$(grep -rn 'workflow-draft\|docs/workflows\|agent-workflow.yaml' "$SITE" 2>/dev/null | wc -l | tr -d ' ')
  if [ "$n" -eq 0 ]; then PASS v11_workflow_parked; return 0; fi
  FAILM v11_workflow_parked "$n served workflow reference(s) remain (expected 0)"; return 1
}

# =============================================================================
# T13 — release gate
# =============================================================================
t13_pr_merged() {
  if ! command -v gh >/dev/null 2>&1; then SKIPM t13_pr_merged "gh not available"; return 2; fi
  local hits
  hits=$(gh pr list --state merged --base master --json headRefName,title 2>/dev/null \
    | grep -iE 'sprint1|sprint 1|feat/sprint1')
  if [ -n "$hits" ]; then PASS t13_pr_merged; return 0; fi
  FAILM t13_pr_merged "no merged Sprint 1 PR against master yet"; return 1
}

# =============================================================================
# Runner
# =============================================================================
ALL_CHECKS=(
  t1_no_ghost_ping t1_schema_healthcheck_present
  t2_hello_has_from t2_spec_from_sentence t2_hello_diff_clean
  t3_step2_callout t3_step5_callout t3_no_present_tense_claims
  t4_brand_casing t4_quickstart_hello_canonical
  t5_sidebar_quickstart t5_sidebar_conformance t5_sidebar_impl_mapping
  t5_sidebar_no_workflows t5_docsindex_cards
  t6_tooling_inbound
  t7_sitemap_no_notes
  t9_no_runner_conformance_label t9_url_preserved
  t10_one_canonical
  t12_no_served_workflow_refs t12_pages_unpublished t12_yaml_dropped
  t12_sidebar_unlinked t12_llms_swept t12_sitemap_no_workflow
  t12_changelog_line t12_sources_kept t12_no_dangling_links
  t8_report_exists t8_report_documents_output t8_report_documents_build
  t8_report_ab_recommendation t8_no_site_edits
  v1_ghost_tool v3_version_coherence v4_hello_from v5_example_lint
  v6_one_canonical v7_internal_links v11_workflow_parked
  t13_pr_merged
)

# Local-applicable subset for v_all_local (checks 1-7, 10 local, 11).
LOCAL_CHECKS=(
  t1_no_ghost_ping t2_hello_has_from t2_spec_from_sentence t2_hello_diff_clean
  t3_step2_callout t3_step5_callout t3_no_present_tense_claims
  t4_brand_casing t4_quickstart_hello_canonical
  t5_sidebar_quickstart t5_sidebar_conformance t5_sidebar_impl_mapping
  t5_sidebar_no_workflows t5_docsindex_cards t6_tooling_inbound
  t9_no_runner_conformance_label t9_url_preserved
  t12_sidebar_unlinked t12_llms_swept t12_changelog_line t12_sources_kept
  v1_ghost_tool v3_version_coherence v4_hello_from v5_example_lint
)

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
  echo "== v_all_local: local-applicable §V subset =="
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
