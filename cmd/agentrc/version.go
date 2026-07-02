package main

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is the CLI build version, injected at release time via
// -ldflags "-X main.version=<tag>". It is "dev" for local/source builds.
var version = "dev"

// specDraft is the Agentfile spec draft this CLI targets.
const specDraft = "0.1.0-draft.6"

// resolveVersion prefers the ldflags-injected release version, then falls back
// to the module version stamped by the Go toolchain (so `go install ...@v0.1.1`
// reports v0.1.1 rather than "dev"), then to "dev".
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return version
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the agentrc CLI version and the Agentfile spec draft it targets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "agentrc %s (spec %s, %s/%s)\n",
				resolveVersion(), specDraft, runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}
