package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Companion files the scaffolded Agentfile references, so `arc init && arc build`
// succeeds out of the box (the Agentfile COPYs ./tools/file_read and its CMD
// runs ./agent.py).
const toolStub = `#!/bin/sh
# Minimal example tool projected at /mnt/tools/file_read.
# Replace with your real tool; --agentrc-schema is used by the HEALTHCHECK.
if [ "$1" = "--agentrc-schema" ]; then echo '{"tool":"file_read","ok":true}'; exit 0; fi
cat "$1" 2>/dev/null
`

const agentStub = `# Minimal example agent entrypoint (referenced by CMD).
# Replace with your real agent loop.
print("hello from agentrc")
`

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

			// Scaffold the companion files the Agentfile references so a fresh
			// `arc build` works immediately. Never clobber existing files.
			dir := filepath.Dir(path)
			toolPath := filepath.Join(dir, "tools", "file_read")
			if err := writeCompanion(cmd, toolPath, toolStub, 0o755); err != nil {
				return err
			}
			if err := writeCompanion(cmd, filepath.Join(dir, "agent.py"), agentStub, 0o644); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "agent name (default: hello)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}

// writeCompanion writes a scaffolded companion file, creating its parent dir.
// It never clobbers an existing file (so re-running init won't wipe your tool).
func writeCompanion(cmd *cobra.Command, path, content string, mode os.FileMode) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", path)
	return nil
}
