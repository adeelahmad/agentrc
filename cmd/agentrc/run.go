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

// translate is the pure per-backend seam. For S2-03 this is a STUB returning a
// minimal non-empty placeholder; S2-04 fills the real per-backend bodies
// (bedrock JSON, kubernetes YAML, local exec plan).
func translate(backend string, labels map[string]string) (string, error) {
	if !isValidBackend(backend) {
		return "", fmt.Errorf("unknown backend %q; valid backends are %s", backend, strings.Join(validBackends, ", "))
	}
	_ = labels
	return fmt.Sprintf("{\"backend\": %q, \"note\": \"translate stub — S2-04 fills the per-backend body\"}\n", backend), nil
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
				out, err := translate(backend, nil)
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
