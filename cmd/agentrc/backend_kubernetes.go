package main

import (
	"fmt"
	"sort"
	"strings"
)

// mcpSidecar is one MCP server delivered at runtime as a sidecar container
// mounting under /mnt/mcp/<name> (§9.3).
type mcpSidecar struct {
	name  string
	image string
}

// translateKubernetes emits plain multi-doc Kubernetes manifests (no Helm) for
// `--backend kubernetes`: a ServiceAccount, a Deployment (resources + env + MCP
// sidecars), a Service, and a deny-by-default egress NetworkPolicy. Manifests
// only — one format (open question #4, UNRESOLVED).
func translateKubernetes(labels map[string]string) (string, error) {
	name := labels["ai.agentrc.identity.name"]
	if name == "" {
		name = "agent"
	}
	image := labels["image.ref"]
	sa := labels["ai.agentrc.substrate.kubernetes.serviceAccount"]
	sidecars := mcpSidecars(labels)

	var b strings.Builder

	// ServiceAccount (§8.7: serviceAccount is a KEY, never a new namespace).
	if sa != "" {
		fmt.Fprintf(&b, "---\napiVersion: v1\nkind: ServiceAccount\nmetadata:\n  name: %s\n", sa)
	}

	// Deployment.
	fmt.Fprintf(&b, "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: %s\n  labels:\n    app: %s\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: %s\n  template:\n    metadata:\n      labels:\n        app: %s\n    spec:\n", name, name, name, name)
	if sa != "" {
		fmt.Fprintf(&b, "      serviceAccountName: %s\n", sa)
	}
	b.WriteString("      containers:\n")
	writeContainer(&b, name, image, labels, sidecars)
	for _, s := range sidecars {
		fmt.Fprintf(&b, "        - name: mcp-%s\n          image: %q\n          volumeMounts:\n            - name: mcp-%s\n              mountPath: /mnt/mcp/%s\n", s.name, s.image, s.name, s.name)
	}
	if len(sidecars) > 0 {
		b.WriteString("      volumes:\n")
		for _, s := range sidecars {
			fmt.Fprintf(&b, "        - name: mcp-%s\n          emptyDir: {}\n", s.name)
		}
	}

	// Service.
	fmt.Fprintf(&b, "---\napiVersion: v1\nkind: Service\nmetadata:\n  name: %s\nspec:\n  selector:\n    app: %s\n  ports:\n    - port: 80\n      targetPort: 8080\n", name, name)

	// deny-by-default egress NetworkPolicy: policyTypes carries Egress and only
	// the requested hosts/ports are allowed; everything else is denied.
	fmt.Fprintf(&b, "---\napiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: %s-egress\nspec:\n  podSelector:\n    matchLabels:\n      app: %s\n  policyTypes:\n    - Egress\n  egress:\n", name, name)
	hosts := dnsAllows(labels)
	if len(hosts) == 0 {
		b.WriteString("    []\n")
	}
	for _, h := range hosts {
		fmt.Fprintf(&b, "    # allow egress to %s\n    - ports:\n        - protocol: TCP\n          port: %s\n", h.host, h.port)
	}

	return b.String(), nil
}

// writeContainer emits the primary agent container with resources, env, and MCP
// volume mounts.
func writeContainer(b *strings.Builder, name, image string, labels map[string]string, sidecars []mcpSidecar) {
	fmt.Fprintf(b, "        - name: %s\n          image: %q\n", name, image)

	mem := k8sMemory(labels["ai.agentrc.substrate.runtime.memory"])
	cpu := labels["ai.agentrc.substrate.runtime.cpu"]
	if mem != "" || cpu != "" {
		b.WriteString("          resources:\n            limits:\n")
		if mem != "" {
			fmt.Fprintf(b, "              memory: %q\n", mem)
		}
		if cpu != "" {
			fmt.Fprintf(b, "              cpu: %q\n", cpu)
		}
	}

	if env := envVars(labels); len(env) > 0 {
		b.WriteString("          env:\n")
		names := make([]string, 0, len(env))
		for n := range env {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Fprintf(b, "            - name: %s\n              value: %q\n", n, env[n])
		}
	}

	if len(sidecars) > 0 {
		b.WriteString("          volumeMounts:\n")
		for _, s := range sidecars {
			fmt.Fprintf(b, "            - name: mcp-%s\n              mountPath: /mnt/mcp/%s\n", s.name, s.name)
		}
	}
}

type dnsAllow struct {
	host string
	port string
}

// dnsAllows collects ai.agentrc.network.dns.<host>=<port> labels, sorted for
// deterministic output.
func dnsAllows(labels map[string]string) []dnsAllow {
	const prefix = "ai.agentrc.network.dns."
	var out []dnsAllow
	for k, v := range labels {
		if host := strings.TrimPrefix(k, prefix); host != k {
			out = append(out, dnsAllow{host: host, port: v})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].host < out[j].host })
	return out
}

// mcpSidecars collects ai.agentrc.mcp.<name>=runtime:<url> labels into sidecar
// definitions, sorted for deterministic output. Cached digests
// (ai.agentrc.mcp.<name>.origin, or non-runtime values) are skipped.
func mcpSidecars(labels map[string]string) []mcpSidecar {
	const prefix = "ai.agentrc.mcp."
	var out []mcpSidecar
	for k, v := range labels {
		rest := strings.TrimPrefix(k, prefix)
		if rest == k || strings.Contains(rest, ".") {
			continue // not an mcp.<name> key (e.g. .origin)
		}
		if !strings.HasPrefix(v, "runtime:") {
			continue
		}
		out = append(out, mcpSidecar{name: rest, image: stripScheme(strings.TrimPrefix(v, "runtime:"))})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func stripScheme(u string) string {
	for _, p := range []string{"https://", "http://", "oci://"} {
		if strings.HasPrefix(u, p) {
			return strings.TrimPrefix(u, p)
		}
	}
	return u
}

// k8sMemory converts an Agentfile memory request (e.g. "8gb") to a Kubernetes
// quantity (e.g. "8Gi"); unknown shapes pass through unchanged.
func k8sMemory(v string) string {
	l := strings.ToLower(strings.TrimSpace(v))
	switch {
	case l == "":
		return ""
	case strings.HasSuffix(l, "gb"):
		return strings.TrimSuffix(l, "gb") + "Gi"
	case strings.HasSuffix(l, "mb"):
		return strings.TrimSuffix(l, "mb") + "Mi"
	default:
		return v
	}
}
