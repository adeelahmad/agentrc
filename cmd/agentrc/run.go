package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var validBackends = []string{"local", "bedrock", "kubernetes"}

func isValidBackend(b string) bool {
	for _, v := range validBackends {
		if v == b {
			return true
		}
	}
	return false
}

// translate is the pure per-backend seam: it maps ai.agentrc.* labels to the
// backend's native config (bedrock JSON, kubernetes manifests, local exec plan).
// It is fail-closed — an unenforceable request yields an error and no config.
func translate(backend string, labels map[string]string) (string, error) {
	switch backend {
	case "local":
		return translateLocal(labels)
	case "bedrock":
		return translateBedrock(labels)
	case "kubernetes":
		return translateKubernetes(labels)
	default:
		return "", fmt.Errorf("unknown backend %q; valid backends are %s", backend, strings.Join(validBackends, ", "))
	}
}

// representativeLabels builds a fully-populated demonstration label set for
// --dry-run. agentrc ships no runtime and dry-run resolves no registry, so the
// reference translators render against a representative agent to show the shape
// of the config each backend would produce for the given ref.
func representativeLabels(ref string) map[string]string {
	return map[string]string{
		"ai.agentrc.identity.name":                       "code-reviewer",
		"ai.agentrc.identity.description":                "Reviews pull requests",
		"image.ref":                                       ref,
		"ai.agentrc.substrate.aws.roleArn":               "arn:aws:iam::123456789012:role/agent-exec",
		"ai.agentrc.substrate.aws.networkMode":           "PUBLIC",
		"ai.agentrc.substrate.aws.securityGroup":         "sg-0abc123,sg-0def456",
		"ai.agentrc.substrate.aws.subnet":                "subnet-0abc123,subnet-0def456",
		"ai.agentrc.substrate.aws.protocol":              "HTTP",
		"ai.agentrc.substrate.aws.maxLifetime":           "1h",
		"ai.agentrc.substrate.aws.deployment.mode":       "code",
		"ai.agentrc.substrate.aws.code.s3.uri":           "s3://acme-agents/code-reviewer.zip",
		"ai.agentrc.substrate.runtime.language":          "python:3.11",
		"ai.agentrc.substrate.runtime.memory":            "8gb",
		"ai.agentrc.substrate.runtime.cpu":               "2",
		"env.LOG_LEVEL":                                   "info",
		"ai.agentrc.agent.idle_timeout":                  "5m",
		"ai.agentrc.agent.auth.mode":                     "jwt",
		"ai.agentrc.agent.auth.jwt.discovery_url":        "https://auth.acme/.well-known/openid-configuration",
		"ai.agentrc.agent.auth.jwt.allowed_audience":     "agentrc://code-reviewer",
		"ai.agentrc.agent.auth.jwt.allowed_client":       "acme-ci,acme-bot",
		"ai.agentrc.network.dns.api.github.com":          "443",
		"ai.agentrc.substrate.kubernetes.serviceAccount": "agent-sa",
		"ai.agentrc.mcp.github":                          "runtime:https://registry.agentrc.io/mcp/github:latest",
	}
}

func newRunCmd() *cobra.Command {
	var (
		backend    string
		isolation  string
		region     string
		profile    string
		kubeconfig string
		namespace  string
		dryRun     bool
	)
	cmd := &cobra.Command{
		Use:   "run <ref>",
		Short: "Run an artifact on a chosen backend (reference translators; --dry-run prints the translated config)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isValidBackend(backend) {
				return fmt.Errorf("unknown --backend %q; valid backends are %s", backend, strings.Join(validBackends, ", "))
			}
			if isolation != "" && backend != "local" {
				return fmt.Errorf("--isolation is only valid for --backend local (got backend %q)", backend)
			}
			if dryRun {
				labels := representativeLabels(args[0])
				if real, env, ok := localImageConfig(args[0]); ok {
					labels = agentLabelSet(args[0], real, env)
				} else {
					fmt.Fprintf(cmd.ErrOrStderr(), "note: %s is not a local image; showing a representative example. Build it (arc build -t %s .) or push it, then re-run.\n", args[0], args[0])
				}
				out, err := translate(backend, labels)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), out)
				return nil
			}
			return fmt.Errorf("agentrc run has no bundled runtime; use --dry-run to print the translated config (agentrc declares agents, it does not ship a runtime — docs/non-goals.md)")
		},
	}
	cmd.Flags().StringVar(&backend, "backend", "local", "execution backend: local|bedrock|kubernetes")
	cmd.Flags().StringVar(&isolation, "isolation", "", "local|container|microvm (only valid with --backend local)")
	cmd.Flags().StringVar(&region, "region", "", "AWS region (--backend bedrock)")
	cmd.Flags().StringVar(&profile, "profile", "", "AWS profile (--backend bedrock)")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig (--backend kubernetes)")
	cmd.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace (--backend kubernetes)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print the translated config and exit (all backends)")
	return cmd
}
