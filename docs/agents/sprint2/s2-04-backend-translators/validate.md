---
type: validate
story: S2-04
---
# S2-04 Validation — Backend translators (T22–T24)

## Pre-flight

- [ ] §8.7/§8.8/§8.9 (S2-02 T17–T19) landed — the translators read those namespaces.
- [ ] `--backend` surface (S2-03 T21) landed — the dispatch seam exists.
- [ ] RED fail-closed tests written first and FAIL before implementation (real TDD).
- [ ] Chosen k8s format = manifests ONLY (no Helm; open question #4 UNRESOLVED).
- [ ] Version unchanged: single `draft.5` (§0.1); `go mod tidy` clean if any dep added.

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | bedrock dry-run valid JSON | `go run ./cmd/agentrc run <ref> --backend bedrock --dry-run \| python3 -m json.tool` | parses, exit 0 (§V.9) |
| 2 | bedrock fail-closed roleArn | `go run ./cmd/agentrc run <ref-without-roleArn> --backend bedrock --dry-run` | non-zero, no config emitted |
| 3 | bedrock fail-closed jwt | `… <ref agent.auth.mode=jwt, no discovery_url> --backend bedrock --dry-run` | non-zero, no endpoint |
| 4 | bedrock fail-closed code-mode | `… <ref deployment.mode=code, no runtime.language> --backend bedrock --dry-run` | non-zero |
| 5 | k8s dry-run yaml valid | `go run ./cmd/agentrc run <ref> --backend kubernetes --dry-run \| python3 -c 'import sys,yaml;list(yaml.safe_load_all(sys.stdin))'` | exit 0 (kubeconform in CI) |
| 6 | k8s deny-by-default | `… --backend kubernetes --dry-run \| grep -c 'NetworkPolicy'` | `>= 1` (deny-by-default) |
| 7 | k8s single format | `… --backend kubernetes --dry-run \| grep -ci 'helm\|Chart.yaml'` | `0` (manifests only) |
| 8 | local §0.8 verbatim | `grep -c 'Reference translators — a proof of concept until platforms read' cmd/agentrc/*.go cli.md` | `>= 1` |
| 9 | scoped build | `go build ./cmd/agentrc/...` | exit 0 (full in CI) |
| 10 | scoped tests | `go test ./cmd/agentrc/...` | PASS (`-race` full in CI) |
| 11 | mod tidy | `go mod tidy && git diff --exit-code go.mod go.sum` | no diff |
| 12 | no fourth namespace | `rg -oE 'org\.agentrc\.[a-z]+' cmd/agentrc/*.go \| sort -u` | only known namespaces (identity/capability/sop/agent/substrate/model/network) |
</content>
