// Command agentrc (and its alias arc, per cli.md) is the native CLI for
// building, inspecting, and publishing agentrc agents — a second front door
// to the same compiler as the BuildKit frontend (cmd/agentrc-frontend);
// `build` shells out to `docker build` with that frontend so both paths are
// guaranteed to produce identical artifacts (spec/index.md §10).
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "agentrc",
		Short:         "Build, inspect, and publish agentrc agents (Dockerfile-shaped Agentfiles)",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newInitCmd(),
		newLintCmd(),
		newLockCmd(),
		newBuildCmd(),
		newInspectCmd(),
		newPushCmd(),
		newPullCmd(),
		newSignCmd(),
		newVerifyCmd(),
		newRunCmd(),
		newVersionCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
