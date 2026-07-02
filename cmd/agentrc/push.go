package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <ref>",
		Short: "Push the artifact to any OCI registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDocker(cmd, "push", "--", args[0])
		},
	}
}

func runDocker(cmd *cobra.Command, args ...string) error {
	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Stdout = cmd.OutOrStdout()
	dockerCmd.Stderr = cmd.ErrOrStderr()
	dockerCmd.Stdin = os.Stdin
	if err := dockerCmd.Run(); err != nil {
		return fmt.Errorf("docker %v: %w", args, err)
	}
	return nil
}
