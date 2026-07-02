package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// positioningLine is the §0.8 line that MUST appear verbatim on the backend
// surface (docs and/or the local backend command), stating the translators are a
// proof of concept, not production runners.
const positioningLine = "Reference translators — a proof of concept until platforms read `org.agentrc.*` labels natively. Not production runners."

// TestLocalBackendDispatches asserts that translate("local", …) dispatches to the
// real local translator seam (microsandbox exec plan) rather than the S2-03
// placeholder stub. FAILS now: translate is a stub whose output contains "stub".
func TestLocalBackendDispatches(t *testing.T) {
	labels := map[string]string{
		"org.agentrc.identity.name":    "code-reviewer",
		"org.agentrc.identity.version": "1.0",
		"image.ref":                    "ghcr.io/acme/code-reviewer:1.0",
	}
	out, err := translate("local", labels)
	if err != nil {
		t.Fatalf("translate(local) should dispatch to the local translator, not error: %v", err)
	}
	if out == "" {
		t.Fatal("translate(local) should return a non-empty local exec plan")
	}
	if strings.Contains(out, "stub") {
		t.Errorf("translate(local) still returns the placeholder stub, not the local translator:\n%s", out)
	}
}

// TestLocalPositioningLineVerbatim asserts the §0.8 positioning line appears
// verbatim in the backend surface — one of the cmd/agentrc/*.go files or the
// top-level cli.md (the surface named by validate.md item 8). FAILS now: the line
// is absent from every candidate surface.
func TestLocalPositioningLineVerbatim(t *testing.T) {
	var surfaces []string

	goFiles, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob cmd/agentrc/*.go: %v", err)
	}
	for _, f := range goFiles {
		// The surface is production code + docs, not the test files (which
		// carry the constant verbatim and would trivially self-satisfy).
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		surfaces = append(surfaces, f)
	}
	surfaces = append(surfaces, filepath.Join("..", "..", "cli.md"))

	for _, path := range surfaces {
		b, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(b), positioningLine) {
			return // found verbatim
		}
	}
	t.Errorf("§0.8 positioning line not found verbatim in any backend surface (cmd/agentrc/*.go or cli.md); expected:\n%q", positioningLine)
}
