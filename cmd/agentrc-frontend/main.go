// Command agentrc-frontend is a BuildKit gateway frontend for
// Dockerfile-shaped Agentfiles (0.1.0-draft.6): it compiles the standard
// Dockerfile surface with BuildKit's own dockerfile2llb compiler and layers
// on the four agentrc keywords (IDENTITY, CAPABILITY, SOP, POLICY) plus the
// ADD --remote extension, emitting org.agentrc.* labels. It is invoked
// either directly via `buildctl build --frontend gateway.v0 --opt
// source=<this-image-ref>`, or automatically by `docker build -f Agentfile
// .` once the Agentfile's first line is `# syntax=<this-image-ref>`.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/dockerui"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	digest "github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/adeelahmad/agentrc/internal/agentfile"
	agentllb "github.com/adeelahmad/agentrc/internal/llb"
)

func main() {
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), build); err != nil {
		panic(err)
	}
}

func build(ctx context.Context, c client.Client) (*client.Result, error) {
	bc, err := dockerui.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("initializing dockerui client: %w", err)
	}

	src, err := bc.ReadEntrypoint(ctx, "dockerfile")
	if err != nil {
		return nil, fmt.Errorf("reading Agentfile (pass -f/--opt filename=<path> if it isn't named Agentfile): %w", err)
	}

	f, err := agentfile.Extract(src.Data)
	if err != nil {
		return nil, fmt.Errorf("parsing Agentfile: %w", err)
	}
	if issues := agentfile.Validate(f); agentfile.HasErrors(issues) {
		var msgs []string
		for _, issue := range issues {
			if issue.Severity == "error" {
				msgs = append(msgs, issue.String())
			}
		}
		return nil, fmt.Errorf("Agentfile is not conformant:\n%s", strings.Join(msgs, "\n"))
	}

	policyMode := c.BuildOpts().Opts["build-arg:AGENTRC_POLICY_MODE"]

	rb, err := bc.Build(ctx, func(ctx context.Context, platform *ocispecs.Platform, _ int) (*dockerui.BuildResult, error) {
		p := ocispecs.Platform{OS: "linux", Architecture: "amd64"}
		if platform != nil {
			p = *platform
		}

		mainCtx, err := bc.MainContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving build context: %w", err)
		}
		llbCaps := c.BuildOpts().LLBCaps

		translated, err := agentllb.Translate(ctx, f, agentllb.Options{
			Convert: dockerfile2llb.ConvertOpt{
				Config: bc.Config,
				// Client alone: dockerfile2llb derives the main context from it.
				// Passing Client AND MainContext is rejected by BuildKit
				// ("Client and MainContext cannot both be provided", convert.go).
				Client:         bc,
				TargetPlatform: &p,
				MetaResolver:   c,
				LLBCaps:        &llbCaps,
			},
			GatewayClient: c,
			MainContext:   *mainCtx,
			PolicyMode:    policyMode,
		})
		if err != nil {
			return nil, fmt.Errorf("translating Agentfile: %w", err)
		}
		for _, w := range translated.Warnings {
			c.Warn(ctx, digest.Digest(""), w, client.WarnOpts{})
		}

		def, err := translated.State.Marshal(ctx)
		if err != nil {
			return nil, fmt.Errorf("marshaling LLB state: %w", err)
		}

		solveRes, err := c.Solve(ctx, client.SolveRequest{Definition: def.ToPB()})
		if err != nil {
			return nil, fmt.Errorf("solving: %w", err)
		}
		ref, err := solveRes.SingleRef()
		if err != nil {
			return nil, err
		}

		return &dockerui.BuildResult{Reference: ref, Image: translated.Image}, nil
	})
	if err != nil {
		return nil, err
	}

	return rb.Finalize()
}
