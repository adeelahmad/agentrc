package llb

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/sourceresolver"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	digest "github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/adeelahmad/agentrc/internal/agentfile"
)

var testPlatform = ocispecs.Platform{OS: "linux", Architecture: "amd64"}

// fakeMetaResolver avoids real registry lookups in tests — the example
// Agentfiles reference fictional base images (ghcr.io/acme/...).
type fakeMetaResolver struct{}

func (fakeMetaResolver) ResolveImageConfig(_ context.Context, ref string, _ sourceresolver.Opt) (string, digest.Digest, []byte, error) {
	cfg := []byte(`{"architecture":"amd64","os":"linux","config":{}}`)
	return ref, digest.FromBytes(cfg), cfg, nil
}

func parseExample(t *testing.T, name string) *agentfile.File {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "examples", name))
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	f, err := agentfile.Extract(b)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	return f
}

func testConvertOpt() dockerfile2llb.ConvertOpt {
	mainCtx := llb.Local("context")
	return dockerfile2llb.ConvertOpt{
		MainContext:    &mainCtx,
		TargetPlatform: &testPlatform,
		MetaResolver:   fakeMetaResolver{},
	}
}

var errUnavailable = errors.New("resource unavailable")

// fakeFetch avoids real network access in tests: it returns deterministic
// content for any source, or an error for sources containing "unavailable".
func fakeFetch(_ context.Context, source string, _ time.Duration) ([]byte, error) {
	if strings.Contains(source, "unavailable") {
		return nil, errUnavailable
	}
	return []byte("fake-content-for:" + source), nil
}

func TestTranslateAllExamples(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{
		"Agentfile.minimal", "Agentfile.code-reviewer", "Agentfile.secure-workspace", "Agentfile.vault-agent",
	} {
		t.Run(name, func(t *testing.T) {
			f := parseExample(t, name)
			res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
			if err != nil {
				t.Fatalf("Translate() error = %v", err)
			}
			if _, err := res.State.Marshal(ctx); err != nil {
				t.Fatalf("State.Marshal() error = %v", err)
			}
			if res.Image.Config.Labels["org.agentrc.identity.name"] == "" {
				t.Error("missing org.agentrc.identity.name label")
			}
		})
	}
}

func TestTranslateMinimalLabels(t *testing.T) {
	ctx := context.Background()
	f := parseExample(t, "Agentfile.minimal")
	res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	labels := res.Image.Config.Labels

	want := map[string]string{
		"org.agentrc.identity.name":   "hello",
		"org.agentrc.capability.text": "true",
		"org.agentrc.sop":             "/mnt/SOP",
		"org.agentrc.tool.file_read":  "local",
		"org.agentrc.model.name":      "claude-sonnet-4",
	}
	for k, v := range want {
		if labels[k] != v {
			t.Errorf("labels[%q] = %q, want %q", k, labels[k], v)
		}
	}
	if labels["org.agentrc.sop.sha256"] == "" {
		t.Error("org.agentrc.sop.sha256 missing")
	}
}

func TestTranslateCodeReviewerRemoteResources(t *testing.T) {
	ctx := context.Background()
	f := parseExample(t, "Agentfile.code-reviewer")
	res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	labels := res.Image.Config.Labels

	// The --cached skill must be embedded: digest + origin.
	if labels["org.agentrc.skill.code-review"] == "" || labels["org.agentrc.skill.code-review"] == "local" {
		t.Errorf("skill label = %q, want a resolved digest", labels["org.agentrc.skill.code-review"])
	}
	if labels["org.agentrc.skill.code-review.origin"] != "https://registry.agentrc.io/skills/code-review:1.2.3" {
		t.Errorf("skill origin = %q", labels["org.agentrc.skill.code-review.origin"])
	}
	// The --runtime mcp server must stay reference-only.
	if labels["org.agentrc.mcp.github"] != "runtime:mcp://registry.internal.acme/servers/github:latest" {
		t.Errorf("mcp label = %q", labels["org.agentrc.mcp.github"])
	}
}

func TestTranslateFailIfUnavailableAbortsBuild(t *testing.T) {
	ctx := context.Background()
	src := "IDENTITY name=a\nFROM alpine\nCAPABILITY text\nSOP x\nCMD run\n" +
		"ADD --remote --cached --fail-if-unavailable https://example.com/unavailable-thing /mnt/tools/x\n"
	f, err := agentfile.Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	_, err = Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
	if err == nil {
		t.Fatal("expected Translate() to fail for an unavailable --fail-if-unavailable resource")
	}
}

func TestTranslateWarnIfUnavailableDegradesToRuntimeLabel(t *testing.T) {
	ctx := context.Background()
	src := "IDENTITY name=a\nFROM alpine\nCAPABILITY text\nSOP x\nCMD run\n" +
		"ADD --remote --cached --warn-if-unavailable https://example.com/unavailable-thing /mnt/tools/x\n"
	f, err := agentfile.Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if len(res.Warnings) != 1 {
		t.Fatalf("Warnings = %v, want 1", res.Warnings)
	}
	if res.Image.Config.Labels["org.agentrc.tool.x"] != "runtime:https://example.com/unavailable-thing" {
		t.Errorf("tool.x label = %q", res.Image.Config.Labels["org.agentrc.tool.x"])
	}
}

func TestTranslatePolicyModeDigest(t *testing.T) {
	ctx := context.Background()
	f := parseExample(t, "Agentfile.secure-workspace")
	res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch, PolicyMode: "digest"})
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	labels := res.Image.Config.Labels
	if labels["org.agentrc.policy.manifest.sha256"] == "" {
		t.Error("expected a policy manifest digest label in digest mode")
	}
	if _, ok := labels["org.agentrc.substrate.runtime.memory"]; ok {
		t.Error("inline policy labels should be replaced by the digest pointer in digest mode")
	}
	// Non-policy labels must still be present.
	if labels["org.agentrc.identity.name"] == "" {
		t.Error("identity labels must survive digest mode")
	}
	if _, err := res.State.Marshal(ctx); err != nil {
		t.Fatalf("State.Marshal() error = %v", err)
	}
}

func TestTranslateFileBackedSOPViaRemoteAdd(t *testing.T) {
	ctx := context.Background()
	src := "IDENTITY name=a\nFROM alpine\nCAPABILITY text\nCMD run\n" +
		"ADD --remote --cached --fail-if-unavailable https://example.com/sop.md /mnt/SOP\n"
	f, err := agentfile.Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	res, err := Translate(ctx, f, Options{Convert: testConvertOpt(), FetchRemote: fakeFetch})
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if res.Image.Config.Labels["org.agentrc.sop"] != "/mnt/SOP" {
		t.Errorf("org.agentrc.sop = %q", res.Image.Config.Labels["org.agentrc.sop"])
	}
	if res.Image.Config.Labels["org.agentrc.sop.sha256"] == "" {
		t.Error("org.agentrc.sop.sha256 is empty for a remote-ADD-backed SOP")
	}
	// A file-backed SOP must never appear as a regular mcp/tool/skill
	// resource label alongside the sop pointer.
	if _, ok := res.Image.Config.Labels["org.agentrc.tool.SOP"]; ok {
		t.Error("file-backed SOP incorrectly also labeled as a generic resource")
	}
}
