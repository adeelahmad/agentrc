package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <ref>",
		Short: "Verify an artifact's signature and provenance — not yet implemented",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("agentrc verify is not implemented yet; see cli.md")
		},
	}
}
