package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const initTemplate = `# syntax=agentrc.agentfile/v0.1
FROM python:3.11-slim

IDENTITY name=%s version=0.1 author=you
CAPABILITY text
SOP You are a concise local assistant. Answer in one paragraph.
CMD python ./agent.py

# A local tool, projected at /mnt/tools/file_read
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

POLICY model.name claude-opus-4
POLICY network dns:api.github.com:443
`

func newInitCmd() *cobra.Command {
	var name string
	var force bool

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Scaffold a starter Agentfile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "Agentfile"
			if len(args) == 1 {
				path = args[0]
			}
			if !force {
				if _, err := os.Stat(path); err == nil {
					return fmt.Errorf("%s already exists (use --force to overwrite)", path)
				}
			}
			if name == "" {
				name = "hello"
			}
			content := fmt.Sprintf(initTemplate, name)
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", path)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "agent name (default: hello)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}
