package main

import (
	"strings"
	"testing"
)

// fullK8sLabels returns a label set exercising every Kubernetes manifest path:
// runtime resources, env, a deny-by-default NetworkPolicy from a network POLICY,
// a ServiceAccount key, and an MCP sidecar mount. Keys follow the emitted-label
// contract:
//   - org.agentrc.substrate.runtime.*        → Deployment resources
//   - env.<NAME> (CLI-injected)              → Deployment env
//   - org.agentrc.network.dns.<host>=<port>  → NetworkPolicy egress allow
//   - org.agentrc.substrate.kubernetes.serviceAccount → ServiceAccount (a KEY)
//   - org.agentrc.mcp.<name>=runtime:<url>   → /mnt/mcp/<name> sidecar
func fullK8sLabels() map[string]string {
	return map[string]string{
		"org.agentrc.identity.name":                       "code-reviewer",
		"image.ref":                                       "ghcr.io/acme/code-reviewer:1.0",
		"org.agentrc.substrate.runtime.memory":            "8gb",
		"org.agentrc.substrate.runtime.cpu":               "2",
		"env.LOG_LEVEL":                                   "info",
		"org.agentrc.network.dns.api.github.com":          "443",
		"org.agentrc.substrate.kubernetes.serviceAccount": "agent-sa",
		"org.agentrc.mcp.github":                          "runtime:https://registry.agentrc.io/mcp/github:latest",
	}
}

// TestK8sEmitsCoreManifests asserts the dry-run output contains the four core
// manifest kinds. FAILS now: translate is a stub emitting placeholder JSON.
func TestK8sEmitsCoreManifests(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	for _, kind := range []string{
		"kind: Deployment",
		"kind: Service",
		"kind: NetworkPolicy",
		"kind: ServiceAccount",
	} {
		if !strings.Contains(out, kind) {
			t.Errorf("kubernetes manifests missing %q\n%s", kind, out)
		}
	}
}

// TestK8sDryRunYAMLParses asserts the dry-run output is structured as multi-doc
// YAML manifests (apiVersion/kind on each doc, separated by ---). Stdlib-only:
// no YAML dependency is added in the RED test. FAILS now: stub emits JSON with
// none of these structural markers.
func TestK8sDryRunYAMLParses(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	if !strings.Contains(out, "apiVersion:") {
		t.Errorf("kubernetes output should contain YAML apiVersion keys, got:\n%s", out)
	}
	if !strings.Contains(out, "kind:") {
		t.Errorf("kubernetes output should contain YAML kind keys, got:\n%s", out)
	}
	if !strings.Contains(out, "---") {
		t.Errorf("kubernetes output should be multi-doc YAML separated by ---, got:\n%s", out)
	}
}

// TestK8sDenyByDefaultNetworkPolicyFromPolicyNetwork asserts a NetworkPolicy is
// emitted from the network DNS labels and that it is deny-by-default: it declares
// an Egress policy type and only allows the requested host. FAILS now: no
// NetworkPolicy is emitted.
func TestK8sDenyByDefaultNetworkPolicyFromPolicyNetwork(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	if !strings.Contains(out, "kind: NetworkPolicy") {
		t.Fatalf("expected a NetworkPolicy derived from network dns labels, got:\n%s", out)
	}
	if !strings.Contains(out, "Egress") {
		t.Errorf("deny-by-default NetworkPolicy must declare an Egress policyType, got:\n%s", out)
	}
	// Deny-by-default: the requested host is explicitly allowed. An empty egress
	// allow-list (with Egress in policyTypes) denies all; a populated one allows
	// only the named hosts. Either shape must reference the requested port.
	if !strings.Contains(out, "443") {
		t.Errorf("NetworkPolicy should allow the requested egress port 443 for api.github.com, got:\n%s", out)
	}
}

// TestK8sServiceAccountFromSubstrateKey asserts substrate.kubernetes.serviceAccount
// is treated as a KEY (§8.7): it yields a ServiceAccount named agent-sa that the
// Deployment references — not a new namespace. FAILS now: absent.
func TestK8sServiceAccountFromSubstrateKey(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	if !strings.Contains(out, "kind: ServiceAccount") {
		t.Errorf("expected a ServiceAccount manifest from substrate.kubernetes.serviceAccount, got:\n%s", out)
	}
	if !strings.Contains(out, "agent-sa") {
		t.Errorf("ServiceAccount / Deployment should reference the serviceAccount key %q, got:\n%s", "agent-sa", out)
	}
	if !strings.Contains(out, "serviceAccountName") {
		t.Errorf("Deployment pod spec should set serviceAccountName from the substrate key, got:\n%s", out)
	}
}

// TestK8sMCPSidecarsFromMntMcp asserts org.agentrc.mcp.* entries become sidecar
// containers in the Deployment. FAILS now: absent.
func TestK8sMCPSidecarsFromMntMcp(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	if !strings.Contains(out, "github") {
		t.Errorf("expected an MCP sidecar container for the github MCP server, got:\n%s", out)
	}
	if !strings.Contains(out, "/mnt/mcp") {
		t.Errorf("MCP sidecar should mount under /mnt/mcp, got:\n%s", out)
	}
}

// TestK8sSingleFormatNoHelm asserts exactly ONE format is emitted — plain
// manifests, never Helm. Output must contain manifest markers (kind:) and MUST
// NOT contain Helm artifacts. FAILS now: stub emits JSON with no "kind:" marker,
// so the manifest-present assertion fails.
func TestK8sSingleFormatNoHelm(t *testing.T) {
	out, err := translate("kubernetes", fullK8sLabels())
	if err != nil {
		t.Fatalf("translate(kubernetes) should not error: %v", err)
	}
	if !strings.Contains(out, "kind:") {
		t.Fatalf("kubernetes output should be plain manifests (kind: …), got:\n%s", out)
	}
	lower := strings.ToLower(out)
	for _, helm := range []string{"chart.yaml", "{{ .values", "{{.values", "helm.sh/"} {
		if strings.Contains(lower, helm) {
			t.Errorf("kubernetes output must be manifests only, found Helm artifact %q:\n%s", helm, out)
		}
	}
}
