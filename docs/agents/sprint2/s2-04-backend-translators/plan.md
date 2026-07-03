---
type: plan
story: S2-04
scope: "tests only"
---
# S2-04 Test Plan — Backend translators (RED)

Tests only. REAL Go test handles under `cmd/agentrc/`. Table-driven; fail-closed cases are
explicit. Every bullet FAILS now (translators not implemented). Scope `./cmd/agentrc/...`
(ENOSPC; `-race`/full build in CI).

## T22 — Backend `local`

- [ ] `cmd/agentrc/backend_local_test.go::TestLocalBackendDispatches` — given `--backend local`,
      assert dispatch reaches the local translator (microsandbox seam), not the old "not
      implemented" error. FAILS now: no local translator.
- [ ] `cmd/agentrc/backend_local_test.go::TestLocalPositioningLineVerbatim` — assert the §0.8
      line "Reference translators — a proof of concept until platforms read `ai.agentrc.*`
      labels natively. Not production runners." appears verbatim in the surface. FAILS now: absent.

## T23 — Backend `bedrock` (mapping + fail-closed)

- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockMapsAllThirteenFields` — given a full
      label set, assert all 13 CreateAgentRuntime fields (agentRuntimeName, containerUri,
      roleArn, networkMode, securityGroups, subnets, serverProtocol, env, customJWTAuthorizer,
      idleRuntimeSessionTimeout, maxLifetime, codeConfiguration, description) are populated.
      FAILS now: translator absent.
- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockDryRunEmitsValidJSON` — assert
      `--dry-run` output unmarshals as JSON. FAILS now: no output.
- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockDryRunFailsClosedWithoutRoleArn` — given
      labels lacking `substrate.aws.roleArn`, bedrock `--dry-run` returns an error and emits NO
      config. FAILS now: backend not implemented (fail-closed path absent).
- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockFailsClosedOnUnenforceableJWT` — given
      `agent.auth.mode=jwt` without a resolvable `agent.auth.jwt.discovery_url`, assert an error
      and NO invocation endpoint emitted (§8.8 fail-closed). FAILS now: absent.
- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockFailsClosedCodeModeWithoutLanguage` —
      given `deployment.mode=code` + `code.s3.uri` but no resolvable `substrate.runtime.language`,
      assert an error (§8.9 fail-closed). FAILS now: absent.
- [ ] `cmd/agentrc/backend_bedrock_test.go::TestBedrockJWTAuthorizerFromAuthLabels` — given
      valid `agent.auth.jwt.*`, assert `customJWTAuthorizer` carries discovery_url + audiences +
      clients. FAILS now: absent.

## T24 — Backend `kubernetes` (manifests + deny-by-default)

- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sEmitsCoreManifests` — assert dry-run
      output contains Deployment, Service, NetworkPolicy, ServiceAccount kinds. FAILS now: absent.
- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sDryRunYAMLParses` — assert `--dry-run`
      output parses as multi-doc YAML. FAILS now: no output.
- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sDenyByDefaultNetworkPolicyFromPolicyNetwork`
      — given `POLICY network dns:api.github.com:443`, assert a NetworkPolicy that denies all
      egress except the requested host (deny-by-default, fail-closed). FAILS now: absent.
- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sServiceAccountFromSubstrateKey` — given
      `substrate.kubernetes.serviceAccount=agent-sa`, assert a ServiceAccount + Deployment
      reference; assert it is treated as a KEY, not a new namespace. FAILS now: absent.
- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sMCPSidecarsFromMntMcp` — given `/mnt/mcp/*`
      entries, assert sidecar containers in the Deployment. FAILS now: absent.
- [ ] `cmd/agentrc/backend_kubernetes_test.go::TestK8sSingleFormatNoHelm` — assert output
      contains NO Helm/Chart.yaml artifacts (one format only). FAILS now: absent.

## Suite guards

- [ ] `scripts/verify-sprint2.sh::v9_backend_dryruns` — assert `--backend bedrock --dry-run |
      python3 -m json.tool` and `--backend kubernetes --dry-run` yaml-parse both pass (§V.9).
      FAILS now: translators absent.
- [ ] `scripts/verify-sprint2.sh::v_no_fourth_namespace` — assert translator code emits only
      known `ai.agentrc.*` namespaces (no invented fourth; §0.3). RED now: green; regression guard.
</content>
