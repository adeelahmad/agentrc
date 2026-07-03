// Package llb translates an extracted agentrc Agentfile into a BuildKit
// llb.State plus OCI image config. It compiles the (already agentrc-keyword
// stripped) CleanedSource with BuildKit's own real dockerfile2llb compiler —
// giving full, correct standard Dockerfile semantics for free — and then
// layers on the four agentrc keywords' effects: ai.agentrc.* labels, the
// embedded /mnt/SOP file, and ADD --remote --cached fetch/embed.
package llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

const defaultHTTPTimeout = 30 * time.Second

// Options carries everything Translate needs beyond the parsed Agentfile.
type Options struct {
	// Convert is passed straight through to dockerfile2llb.Dockerfile2LLB.
	// The caller (cmd/agentrc-frontend) wires Convert.Client,
	// Convert.MainContext, Convert.TargetPlatform, etc. from its dockerui.Client.
	Convert dockerfile2llb.ConvertOpt

	// GatewayClient and MainContext are used to read a file-backed SOP's
	// source content (for its digest label) directly from the build
	// context, the same way dockerui reads the Agentfile itself.
	GatewayClient client.Client
	MainContext   llb.State

	// PolicyMode is "inline" (default) or "digest" — see
	// spec/index.md §9.4 / profiles/oci-package.md §6.
	PolicyMode string

	SourceURL, Revision, Created string
	HTTPTimeout                  time.Duration

	// FetchRemote fetches an ADD --remote --cached resource's content.
	// Defaults to a real http(s) GET; overridable so tests don't need
	// network access.
	FetchRemote func(ctx context.Context, source string, timeout time.Duration) ([]byte, error)
}

// Result is the translated build: the LLB state and the image config to
// export it with, plus any non-fatal warnings (from --warn-if-unavailable
// resources that couldn't be resolved).
type Result struct {
	State    llb.State
	Image    *dockerspec.DockerOCIImage
	Warnings []string
}

// Translate compiles f into LLB. f.CleanedSource is compiled with
// BuildKit's own dockerfile2llb for real Dockerfile semantics (FROM,
// multi-stage, ARG/ENV, COPY, HEALTHCHECK, ...); the agentrc surface
// (IDENTITY/CAPABILITY/SOP/POLICY/ADD --remote) is layered on afterward.
func Translate(ctx context.Context, f *agentfile.File, opts Options) (*Result, error) {
	if err := agentfile.PopulateLocalResources(f); err != nil {
		return nil, err
	}

	convRes, err := dockerfile2llb.Dockerfile2LLB(ctx, f.CleanedSource, opts.Convert)
	if err != nil {
		return nil, fmt.Errorf("compiling Dockerfile instructions: %w", err)
	}
	state := convRes.State
	img := convRes.Image

	var warnings []string
	httpTimeout := opts.HTTPTimeout
	if httpTimeout <= 0 {
		httpTimeout = defaultHTTPTimeout
	}
	fetch := opts.FetchRemote
	if fetch == nil {
		fetch = fetchRemote
	}

	sopSHA256, err := resolveFileBackedSOPDigest(ctx, f, opts)
	if err != nil {
		return nil, err
	}
	if f.SOP != nil && !f.SOP.FileBacked {
		state = writeFile(state, "/mnt/SOP", []byte(f.SOP.Content))
	}

	for i := range f.RemoteAdds {
		ra := &f.RemoteAdds[i]
		if ra.Runtime {
			continue // reference-only; internal/agentfile.BuildLabels handles the label
		}

		// Defense in depth: the caller is expected to have already run
		// agentfile.Validate (which rejects a Dest outside /mnt), but this
		// writes fetched remote content to an absolute path in the image,
		// so it never trusts the raw, uncleaned Dest even if that
		// invariant is somehow violated upstream.
		cleanedDest, ok := agentfile.CleanMntPath(ra.Dest)
		if !ok {
			return nil, fmt.Errorf("line %d: ADD --remote destination must be under /mnt, got %q", ra.Line, ra.Dest)
		}
		ra.Dest = cleanedDest

		data, ferr := fetch(ctx, ra.Source, httpTimeout)
		if ferr != nil {
			msg := fmt.Sprintf("line %d: ADD --remote %s: %v", ra.Line, ra.Source, ferr)
			if ra.FailIfUnavailable {
				return nil, fmt.Errorf("%s", msg)
			}
			warnings = append(warnings, msg)
			continue // ResolvedDigest stays empty; BuildLabels degrades to a runtime: reference
		}

		digest := "sha256:" + hashHex(data)
		ra.ResolvedDigest = digest
		if ra.Dest == "/mnt/SOP" {
			sopSHA256 = digest
		}
		state, err = embedFile(state, ra.Dest, data, ra.Chmod, ra.Chown)
		if err != nil {
			return nil, fmt.Errorf("line %d: ADD --remote: %w", ra.Line, err)
		}
	}

	labels, err := agentfile.BuildLabels(f, agentfile.LabelOptions{
		SOPSHA256: sopSHA256,
		SourceURL: opts.SourceURL,
		Revision:  opts.Revision,
		Created:   opts.Created,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range agentfile.BuildOCIAnnotations(f, agentfile.LabelOptions{SourceURL: opts.SourceURL, Revision: opts.Revision, Created: opts.Created}) {
		labels[k] = v
	}

	policyMode := opts.PolicyMode
	if policyMode == "" {
		policyMode = "inline"
	}
	switch policyMode {
	case "inline":
	case "digest":
		state, labels, err = applyDigestPolicyMode(state, labels)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown --policy-mode %q (want inline or digest)", policyMode)
	}

	if img.Config.Labels == nil {
		img.Config.Labels = map[string]string{}
	}
	for k, v := range labels {
		img.Config.Labels[k] = v
	}

	return &Result{State: state, Image: img, Warnings: warnings}, nil
}

// resolveFileBackedSOPDigest handles the file-backed SOP form (`COPY
// ./sop.md /mnt/SOP`): the content lives in a build-context file, so its
// digest is computed by reading that file back through the gateway client,
// the same way dockerui reads the Agentfile itself. ADD --remote-backed SOP
// is handled by the RemoteAdd loop instead (its digest comes from the fetch).
func resolveFileBackedSOPDigest(ctx context.Context, f *agentfile.File, opts Options) (string, error) {
	if f.SOP == nil || !f.SOP.FileBacked {
		return "", nil
	}
	var sopSource string
	for _, r := range f.LocalResources {
		if r.Dest == "/mnt/SOP" {
			sopSource = r.Source
			break
		}
	}
	if sopSource == "" {
		return "", nil // file-backed via ADD --remote; resolved in the RemoteAdd loop
	}
	if opts.GatewayClient == nil {
		return "", fmt.Errorf("file-backed SOP requires a gateway client to read %s", sopSource)
	}
	data, err := readContextFile(ctx, opts.GatewayClient, opts.MainContext, sopSource)
	if err != nil {
		return "", fmt.Errorf("reading file-backed SOP %s: %w", sopSource, err)
	}
	return "sha256:" + hashHex(data), nil
}

func readContextFile(ctx context.Context, gw client.Client, mainCtx llb.State, filename string) ([]byte, error) {
	def, err := mainCtx.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	res, err := gw.Solve(ctx, client.SolveRequest{Definition: def.ToPB()})
	if err != nil {
		return nil, err
	}
	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}
	return ref.ReadFile(ctx, client.ReadRequest{Filename: filename})
}

// maxFetchSize bounds how much a single ADD --remote --cached fetch will
// embed, so a malicious or misbehaving endpoint can't exhaust build memory.
const maxFetchSize = 200 * 1024 * 1024 // 200MiB

// fetchClient rejects redirects to disallowed targets (see
// agentfile.ValidateRemoteSourceURL) instead of blindly following them —
// the initial URL can pass validation and still redirect somewhere internal.
var fetchClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		return agentfile.ValidateRemoteSourceURL(req.URL.String())
	},
}

func fetchRemote(ctx context.Context, source string, timeout time.Duration) ([]byte, error) {
	if err := agentfile.ValidateRemoteSourceURL(source); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
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

func writeFile(state llb.State, dest string, data []byte) llb.State {
	dir := path.Dir(dest)
	return state.File(
		llb.Mkdir(dir, 0o755, llb.WithParents(true)).
			Mkfile(dest, 0o644, data),
	)
}

func embedFile(state llb.State, dest string, data []byte, chmod, chown string) (llb.State, error) {
	mode := os.FileMode(0o644)
	if chmod != "" {
		m, err := strconv.ParseUint(chmod, 8, 32)
		if err != nil {
			return state, fmt.Errorf("invalid --chmod %q: %w", chmod, err)
		}
		mode = os.FileMode(m)
	}
	var opts []llb.MkfileOption
	if chown != "" {
		opts = append(opts, llb.WithUser(chown))
	}
	dir := path.Dir(dest)
	return state.File(
		llb.Mkdir(dir, 0o755, llb.WithParents(true)).
			Mkfile(dest, mode, data, opts...),
	), nil
}

// applyDigestPolicyMode implements --policy-mode digest (spec/index.md
// §9.4): the POLICY-derived labels (agent.*/substrate.*/model.*/network.*)
// are pulled out of the inline label set, written as a JSON manifest layer,
// and replaced with a single digest pointer label.
func applyDigestPolicyMode(state llb.State, labels map[string]string) (llb.State, map[string]string, error) {
	policyLabels := map[string]string{}
	kept := map[string]string{}
	for k, v := range labels {
		if strings.HasPrefix(k, "ai.agentrc.agent.") ||
			strings.HasPrefix(k, "ai.agentrc.substrate.") ||
			strings.HasPrefix(k, "ai.agentrc.model.") ||
			strings.HasPrefix(k, "ai.agentrc.network.") {
			policyLabels[k] = v
			continue
		}
		kept[k] = v
	}
	if len(policyLabels) == 0 {
		return state, kept, nil
	}

	// encoding/json marshals map[string]string keys in sorted order, so this
	// manifest (and therefore its digest) is deterministic for the same input.
	manifest, err := json.Marshal(policyLabels)
	if err != nil {
		return state, nil, fmt.Errorf("marshaling policy manifest: %w", err)
	}
	digest := "sha256:" + hashHex(manifest)
	state = writeFile(state, "/.agentrc/policy-manifest.json", manifest)
	kept["ai.agentrc.policy.manifest.sha256"] = digest
	return state, kept, nil
}

func hashHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
