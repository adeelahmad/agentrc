package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

// defaultFrontendImage is the BuildKit frontend image this build routes
// through. It must be built (Dockerfile.frontend) and available to the
// local Docker daemon before `agentrc build` can succeed; see
// tooling/README.md. Override with --frontend-image once published.
const defaultFrontendImage = "local/agentrc-frontend:dev"

func newBuildCmd() *cobra.Command {
	var file, tag, policyMode, frontendImage string

	cmd := &cobra.Command{
		Use:   "build [context]",
		Short: "Compile an Agentfile to an OCI artifact, emitting org.agentrc.* labels",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildContext := "."
			if len(args) == 1 {
				buildContext = args[0]
			}
			if policyMode != "inline" && policyMode != "digest" {
				return fmt.Errorf("--policy-mode must be inline or digest, got %q", policyMode)
			}

			// Validate locally first for fast feedback before invoking docker.
			f, err := readAndExtract(file)
			if err != nil {
				return err
			}
			if issues := agentfile.Validate(f); agentfile.HasErrors(issues) {
				return fmt.Errorf("%s failed validation; run `agentrc lint` for details", file)
			}

			dockerArgs := []string{
				"build",
				"-f", file,
				"--build-arg", "BUILDKIT_SYNTAX=" + frontendImage,
				"--build-arg", "AGENTRC_POLICY_MODE=" + policyMode,
			}
			if tag != "" {
				dockerArgs = append(dockerArgs, "-t", tag)
			}
			dockerArgs = append(dockerArgs, "--", buildContext)

			dockerCmd := exec.Command("docker", dockerArgs...)
			dockerCmd.Stdout = cmd.OutOrStdout()
			dockerCmd.Stderr = cmd.ErrOrStderr()
			dockerCmd.Stdin = os.Stdin
			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("docker %v: %w", dockerArgs, err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "Agentfile", "path to the Agentfile")
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "image reference, e.g. ghcr.io/org/agent:1.0")
	cmd.Flags().StringVar(&policyMode, "policy-mode", "inline", "how POLICY requests are encoded: inline or digest")
	cmd.Flags().StringVar(&frontendImage, "frontend-image", defaultFrontendImage, "BuildKit frontend image to build through")
	return cmd
}
