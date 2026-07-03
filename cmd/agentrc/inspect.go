package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

func newInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <path-or-ref>",
		Short: "Read an artifact's org.agentrc.* labels to review what an agent requests before it runs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			if _, err := os.Stat(target); err == nil {
				return inspectLocalAgentfile(cmd, target)
			}
			return inspectRemoteArtifact(cmd, target)
		},
	}
}

func inspectLocalAgentfile(cmd *cobra.Command, path string) error {
	f, err := readAndExtract(path)
	if err != nil {
		return err
	}
	if err := agentfile.PopulateLocalResources(f); err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "identity:\n")
	for _, k := range sortedKeys(f.Identity) {
		fmt.Fprintf(w, "  %s: %s\n", k, f.Identity[k])
	}
	fmt.Fprintf(w, "capabilities: %s\n", joinOrNone(f.Capabilities))
	if f.SOP != nil {
		fmt.Fprintf(w, "sop: declared (file-backed=%v)\n", f.SOP.FileBacked)
	} else {
		fmt.Fprintf(w, "sop: (none)\n")
	}
	fmt.Fprintf(w, "resources:\n")
	for _, r := range f.LocalResources {
		fmt.Fprintf(w, "  %s %s -> %s (local)\n", r.Kind, r.Name, r.Dest)
	}
	for _, ra := range f.RemoteAdds {
		delivery := "cached"
		if ra.Runtime {
			delivery = "runtime"
		}
		fmt.Fprintf(w, "  %s -> %s (%s, %s)\n", ra.Source, ra.Dest, delivery, failModeString(ra.FailIfUnavailable))
	}
	fmt.Fprintf(w, "policy:\n")
	for _, p := range f.Policies {
		fmt.Fprintf(w, "  %s = %s\n", p.Key, p.Value)
	}
	return nil
}

func failModeString(failIfUnavailable bool) string {
	if failIfUnavailable {
		return "fail-if-unavailable"
	}
	return "warn-if-unavailable"
}

func inspectRemoteArtifact(cmd *cobra.Command, ref string) error {
	ctx, cancel := withRegistryTimeout(cmd)
	defer cancel()

	// Prefer a locally-built image (e.g. `arc build -t hello .` before pushing),
	// read through the user's own docker (which knows the active context);
	// fall back to the registry for a pushed reference.
	if labels, _, ok := localImageConfig(ref); ok {
		printAgentrcLabels(cmd, labels)
		return nil
	}

	r, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %w", ref, err)
	}
	img, err := remote.Image(r, remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("inspecting %s: not found as a local image and not pullable from a registry: %w", ref, err)
	}
	cfg, err := img.ConfigFile()
	if err != nil {
		return fmt.Errorf("reading config for %s: %w", ref, err)
	}
	printAgentrcLabels(cmd, cfg.Config.Labels)
	return nil
}

func printAgentrcLabels(cmd *cobra.Command, labels map[string]string) {
	w := cmd.OutOrStdout()
	for _, k := range sortedKeys(labels) {
		if strings.HasPrefix(k, "org.agentrc.") {
			fmt.Fprintf(w, "%s=%s\n", k, labels[k])
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func joinOrNone(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	return strings.Join(items, ", ")
}
