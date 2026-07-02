---
type: validate
story: S2-06
---
# S2-06 Validation — Demo (T26)

## Pre-flight

- [ ] T20 (draft.6 bump + code-reviewer commented block) + S2-04 (all three backends) +
      S2-05 (CLI docs) landed.
- [ ] Confirmed the narrative string will be pasted VERBATIM (no paraphrase; §0.8).
- [ ] Grep-located the demo doc home before writing (M-003).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | narrative verbatim | `grep -c 'Same artifact, same labels, three substrates. The translators are the proof of concept; the labels are the standard.' <demo-file>` | `1` |
| 2 | build command | `grep -c 'arc build -t ghcr.io/agentrc/code-reviewer:1.0' <demo-file>` | `>= 1` |
| 3 | local command | `grep -c 'backend local --isolation microvm' <demo-file>` | `>= 1` |
| 4 | bedrock command | `grep -c 'backend bedrock --dry-run' <demo-file>` | `>= 1` |
| 5 | kubernetes command | `grep -c 'backend kubernetes --dry-run' <demo-file>` | `>= 1` |
| 6 | no --substrate | `rg -c -- '--substrate' <demo-file> 2>/dev/null \|\| echo 0` | `0` |
| 7 | bedrock dry-run parses | `go run ./cmd/agentrc run ghcr.io/agentrc/code-reviewer:1.0 --backend bedrock --dry-run \| python3 -m json.tool` | exit 0 (§V.9) |
| 8 | k8s dry-run parses | `go run ./cmd/agentrc run … --backend kubernetes --dry-run \| python3 -c 'import sys,yaml;list(yaml.safe_load_all(sys.stdin))'` | exit 0 (§V.9) |
| 9 | site build (CI) | `bundle exec jekyll build --trace` | exit 0 |
</content>
