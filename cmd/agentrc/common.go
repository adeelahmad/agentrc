package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

func readAndExtract(path string) (*agentfile.File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	f, err := agentfile.Extract(data)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return f, nil
}

func agentfilePathArg(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return "Agentfile"
}

func writeFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// registryTimeout bounds registry network operations so a hung or slow
// registry connection can't block the CLI indefinitely.
const registryTimeout = 5 * time.Minute

func withRegistryTimeout(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	return context.WithTimeout(cmd.Context(), registryTimeout)
}
