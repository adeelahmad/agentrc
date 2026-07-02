package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Signing/verification (Sigstore) and running an artifact on a chosen
// substrate are deliberately not implemented yet — each needs its own
// design pass (crypto integration; a real execution runtime, which is
// explicitly out of scope for agentrc itself per docs/non-goals.md) rather
// than a partial implementation bolted onto this change. These stubs exist
// so the commands fail with a clear, honest message instead of not
// existing at all (cli.md lists all three as planned).
func newSignCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sign <ref>",
		Short: "Sign an artifact (Sigstore) — not yet implemented",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("agentrc sign is not implemented yet; see cli.md")
		},
	}
}
