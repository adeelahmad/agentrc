package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

func newLintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint [path]",
		Short: "Check an Agentfile for keyword and request errors before building",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := agentfilePathArg(args)
			f, err := readAndExtract(path)
			if err != nil {
				return err
			}

			issues := agentfile.Validate(f)
			for _, issue := range issues {
				fmt.Fprintln(cmd.OutOrStdout(), issue.String())
			}
			if agentfile.HasErrors(issues) {
				return fmt.Errorf("%s failed validation", path)
			}
			if len(issues) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: ok\n", path)
			}
			return nil
		},
	}
}
