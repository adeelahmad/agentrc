package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

// resolvedManifest mirrors schemas/agentrc-lock.schema.json ("agentrc
// Resolved Manifest").
type resolvedManifest struct {
	Version         string             `json:"version"`
	AgentfileSHA256 string             `json:"agentfile_sha256"`
	Timestamp       string             `json:"timestamp"`
	PolicyMode      string             `json:"policy_mode,omitempty"`
	LabelsDigest    string             `json:"labels_digest,omitempty"`
	Base            *resolvedBase      `json:"base,omitempty"`
	Resources       []resolvedResource `json:"resources,omitempty"`
	SOP             *resolvedSOP       `json:"sop,omitempty"`
}

type resolvedBase struct {
	Ref    string `json:"ref"`
	Digest string `json:"digest,omitempty"`
}

type resolvedResource struct {
	Name     string `json:"name"`
	Kind     string `json:"kind,omitempty"`
	Dest     string `json:"dest,omitempty"`
	Delivery string `json:"delivery,omitempty"`
	Digest   string `json:"digest,omitempty"`
	Origin   string `json:"origin,omitempty"`
	FailMode string `json:"fail_mode,omitempty"`
}

type resolvedSOP struct {
	SHA256 string `json:"sha256,omitempty"`
}

func newLockCmd() *cobra.Command {
	var out, policyMode string

	cmd := &cobra.Command{
		Use:   "lock [path]",
		Short: "Pin ADD --remote resources to digests for reproducible builds",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := agentfilePathArg(args)
			f, err := readAndExtract(path)
			if err != nil {
				return err
			}
			if err := agentfile.PopulateLocalResources(f); err != nil {
				return err
			}

			ctx, cancel := withRegistryTimeout(cmd)
			defer cancel()

			manifest := &resolvedManifest{
				Version:         "0.1.0-draft.5",
				AgentfileSHA256: hashHex([]byte(f.Source)),
				Timestamp:       time.Now().UTC().Format(time.RFC3339),
				PolicyMode:      policyMode,
			}

			if ref, err := agentfile.BaseImageRef(f); err == nil {
				base := &resolvedBase{Ref: ref}
				if digest, derr := resolveImageDigest(ctx, ref); derr == nil {
					base.Digest = digest
				} else {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not resolve base digest for %s: %v\n", ref, derr)
				}
				manifest.Base = base
			}

			sopDir := filepath.Dir(path)
			var sopSHA256 string

			for _, r := range f.LocalResources {
				entry := resolvedResource{Name: r.Name, Kind: string(r.Kind), Dest: r.Dest, Delivery: "local"}
				if r.Dest == "/mnt/SOP" {
					if sopPath, perr := safeJoin(sopDir, r.Source); perr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: refusing to read file-backed SOP %s: %v\n", r.Source, perr)
					} else if data, rerr := os.ReadFile(sopPath); rerr == nil {
						sopSHA256 = "sha256:" + hashHex(data)
					} else {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not read file-backed SOP %s: %v\n", r.Source, rerr)
					}
					continue // SOP is reported via manifest.SOP, not as a generic resource
				}
				manifest.Resources = append(manifest.Resources, entry)
			}

			for _, ra := range f.RemoteAdds {
				failMode := "fail"
				if !ra.FailIfUnavailable {
					failMode = "warn"
				}
				entry := resolvedResource{
					Name:     filepath.Base(ra.Dest),
					Dest:     ra.Dest,
					Origin:   ra.Source,
					FailMode: failMode,
				}
				if ra.Runtime {
					entry.Delivery = "runtime"
				} else {
					entry.Delivery = "cached"
					if data, ferr := fetchHTTP(ctx, ra.Source); ferr == nil {
						entry.Digest = "sha256:" + hashHex(data)
						if ra.Dest == "/mnt/SOP" {
							sopSHA256 = entry.Digest
						}
					} else {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not resolve %s: %v\n", ra.Source, ferr)
					}
				}
				if ra.Dest != "/mnt/SOP" {
					manifest.Resources = append(manifest.Resources, entry)
				}
			}

			if f.SOP != nil {
				if !f.SOP.FileBacked {
					sopSHA256 = "sha256:" + hashHex([]byte(f.SOP.Content))
				}
				if sopSHA256 != "" {
					manifest.SOP = &resolvedSOP{SHA256: sopSHA256}
				}
			}

			if labels, lerr := agentfile.BuildLabels(f, agentfile.LabelOptions{SOPSHA256: sopSHA256}); lerr == nil {
				if b, merr := json.Marshal(labels); merr == nil {
					manifest.LabelsDigest = "sha256:" + hashHex(b)
				}
			}

			data, err := json.MarshalIndent(manifest, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling lock manifest: %w", err)
			}
			if out == "" {
				out = "agentrc.lock"
			}
			if err := writeFile(out, data); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", out)
			return nil
		},
	}

	cmd.Flags().StringVar(&out, "out", "", "output path (default: agentrc.lock)")
	cmd.Flags().StringVar(&policyMode, "policy-mode", "inline", "recorded policy encoding: inline or digest")
	return cmd
}

func resolveImageDigest(ctx context.Context, ref string) (string, error) {
	r, err := name.ParseReference(ref)
	if err != nil {
		return "", err
	}
	img, err := remote.Image(r, remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", err
	}
	h, err := img.Digest()
	if err != nil {
		return "", err
	}
	return h.String(), nil
}

// maxFetchSize bounds how much a single resource fetch will read while
// resolving a digest, so a malicious or misbehaving endpoint can't exhaust
// memory during `agentrc lock`.
const maxFetchSize = 200 * 1024 * 1024 // 200MiB

// fetchClient rejects redirects to disallowed targets (see
// agentfile.ValidateRemoteSourceURL) instead of blindly following them.
var fetchClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		return agentfile.ValidateRemoteSourceURL(req.URL.String())
	},
}

// safeJoin resolves rel against base and refuses to return a path outside
// base — a COPY source is normally a relative path inside the build
// context, but a crafted Agentfile (e.g. one from an untrusted fork in a CI
// pipeline) could set it to something like "../../../../etc/shadow".
func safeJoin(base, rel string) (string, error) {
	joined := filepath.Join(base, rel)
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	if absJoined != absBase && !strings.HasPrefix(absJoined, absBase+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes %q", rel, base)
	}
	return joined, nil
}

func fetchHTTP(ctx context.Context, source string) ([]byte, error) {
	if err := agentfile.ValidateRemoteSourceURL(source); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}
	resp, err := fetchClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxFetchSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxFetchSize {
		return nil, fmt.Errorf("response exceeds %d byte limit", maxFetchSize)
	}
	return data, nil
}

func hashHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
