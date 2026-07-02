package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <ref>",
		Short: "Run an artifact on a chosen substrate — not yet implemented",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("agentrc run is not implemented yet — agentrc declares agents, it does not ship a runtime (docs/non-goals.md); see cli.md")
		},
	}
	cmd.Flags().String("isolation", "", "local|container|microvm (not yet implemented)")
	cmd.Flags().String("substrate", "", "substrate driver (not yet implemented)")
	return cmd
}
