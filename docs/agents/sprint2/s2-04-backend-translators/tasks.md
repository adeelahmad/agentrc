---
type: tasks
story: S2-04
---
# S2-04 Tasks — Backend translators (T22, T23, T24)

REAL TDD; each translator MUST have fail-closed tests written first. Pure labels→config
mapping, unit-testable without a live cloud; assert through `--dry-run` output. Prefer
plain-struct/YAML emission over heavy SDK/k8s clients (ENOSPC dep-tree risk); keep
`go mod tidy` clean. Consumes §8.7/§8.8/§8.9. Version stays `draft.5` (§0.1).

## T22 — Backend `local` [P1]

`cmd/agentrc/backend_local.go`:
- Wire the existing microsandbox VMM MVP under `--backend local` (default) via the
  `translate`/dispatch seam from T21. Plumbing only (dry-run/plumbing depth acceptable —
  open question #5, UNRESOLVED).
- Include the §0.8 positioning line VERBATIM in the command/docs surface: "Reference
  translators — a proof of concept until platforms read `ai.agentrc.*` labels natively.
  Not production runners."

## T23 — Backend `bedrock` (labels → CreateAgentRuntime) [P1]

`cmd/agentrc/backend_bedrock.go`: map `ai.agentrc.*` labels + image config → Bedrock
`CreateAgentRuntime` fields (13/13):
- agentRuntimeName / description ← IDENTITY; containerUri ← OCI ref; roleArn ←
  `substrate.aws.roleArn`; networkMode ← `substrate.aws.networkMode`; securityGroups/subnets
  ← `substrate.aws.securityGroup`/`subnet`; serverProtocol ← `substrate.aws.protocol`; env ←
  image Env; customJWTAuthorizer ← `agent.auth.jwt.*`; idleRuntimeSessionTimeout ←
  `agent.idle_timeout`; maxLifetime ← `substrate.aws.maxLifetime`; codeConfiguration ←
  `deployment.mode=code` + `code.s3.uri` + `substrate.runtime.language`.
- **Fail closed (MUST NOT emit config / endpoint):** missing `roleArn`; unenforceable
  `agent.auth.mode=jwt` (e.g. no `discovery_url`); code-mode without a resolvable
  `substrate.runtime.language`.
- `--dry-run` emits the translated JSON (must satisfy `python3 -m json.tool`).

## T24 — Backend `kubernetes` [P1]

`cmd/agentrc/backend_kubernetes.go`: emit (dry-run) or apply, in ONE format = **manifests**
(Helm NOT emitted — open question #4, UNRESOLVED):
- Deployment (resources from `substrate.runtime.*`, env from image config).
- Service.
- **deny-by-default NetworkPolicy** derived from `POLICY network dns:*` (a policy that denies
  all egress except the requested hosts — the fail-closed default).
- ServiceAccount from `substrate.kubernetes.serviceAccount` (a KEY under §8.7, not a new
  namespace).
- MCP servers from `/mnt/mcp/*` as sidecar containers.
- `--dry-run` emits the manifests (must yaml-parse / kubeconform-validate).

Verify: `--backend bedrock --dry-run | python3 -m json.tool` parses; `--backend kubernetes
--dry-run` yaml-parses and contains a deny-by-default NetworkPolicy; all fail-closed tests
green; §0.8 line verbatim; single k8s format; scoped `go build/vet/test ./cmd/agentrc/...`.
</content>
