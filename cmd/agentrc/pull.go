package main

import "github.com/spf13/cobra"

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <ref>",
		Short: "Pull an artifact from any OCI registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDocker(cmd, "pull", args[0])
		},
	}
}
