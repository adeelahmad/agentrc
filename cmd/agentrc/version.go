package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// version is the CLI build version, injected at release time via
// -ldflags "-X main.version=<tag>". It is "dev" for local/source builds.
var version = "dev"

// specDraft is the Agentfile spec draft this CLI targets.
const specDraft = "0.1.0-draft.6"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the agentrc CLI version and the Agentfile spec draft it targets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "agentrc %s (spec %s, %s/%s)\n",
				version, specDraft, runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}
