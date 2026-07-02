package agentfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func readExample(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "examples", name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return b
}

// assertCleanedSourceIsPlainDockerfile confirms CleanedSource has no
// agentrc-specific lines left — it must parse as a standard Dockerfile with
// BuildKit's own instruction parser without error.
func assertCleanedSourceIsPlainDockerfile(t *testing.T, cleaned []byte) []instructions.Stage {
	t.Helper()
	res, err := parser.Parse(strings.NewReader(string(cleaned)))
	if err != nil {
		t.Fatalf("CleanedSource is not valid Dockerfile syntax: %v\n%s", err, cleaned)
	}
	stages, _, err := instructions.Parse(res.AST, nil)
	if err != nil {
		t.Fatalf("CleanedSource has an instruction BuildKit doesn't recognize (agentrc keyword not stripped?): %v\n%s", err, cleaned)
	}
	return stages
}

func TestExtractMinimal(t *testing.T) {
	f, err := Extract(readExample(t, "Agentfile.minimal"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if f.Identity["name"] != "hello" || f.Identity["version"] != "0.1" || f.Identity["author"] != "acme" {
		t.Errorf("Identity = %+v", f.Identity)
	}
	if f.Identity["description"] != "Minimal AgentRC agent" {
		t.Errorf("Identity[description] = %q", f.Identity["description"])
	}
	if len(f.Capabilities) != 1 || f.Capabilities[0] != "text" {
		t.Errorf("Capabilities = %v", f.Capabilities)
	}
	if f.SOP == nil || f.SOP.FileBacked {
		t.Fatalf("SOP = %+v, want inline SOP", f.SOP)
	}
	if !strings.Contains(f.SOP.Content, "minimal example agent") {
		t.Errorf("SOP.Content = %q", f.SOP.Content)
	}
	if len(f.Policies) != 3 {
		t.Fatalf("len(Policies) = %d, want 3: %+v", len(f.Policies), f.Policies)
	}
	if f.Policies[0].Key != "model.name" || f.Policies[0].Value != "claude-sonnet-4" {
		t.Errorf("Policies[0] = %+v", f.Policies[0])
	}

	stages := assertCleanedSourceIsPlainDockerfile(t, f.CleanedSource)
	if len(stages) != 1 {
		t.Fatalf("len(stages) = %d, want 1", len(stages))
	}
	foundCopy, foundHealthcheck := false, false
	for _, cmd := range stages[0].Commands {
		switch cmd.(type) {
		case *instructions.CopyCommand:
			foundCopy = true
		case *instructions.HealthCheckCommand:
			foundHealthcheck = true
		}
	}
	if !foundCopy || !foundHealthcheck {
		t.Errorf("expected COPY and HEALTHCHECK to survive cleaning: copy=%v healthcheck=%v", foundCopy, foundHealthcheck)
	}
}

func TestExtractCodeReviewer(t *testing.T) {
	f, err := Extract(readExample(t, "Agentfile.code-reviewer"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if f.Identity["name"] != "code-reviewer" {
		t.Errorf("Identity[name] = %q", f.Identity["name"])
	}
	if len(f.Capabilities) != 3 {
		t.Errorf("Capabilities = %v", f.Capabilities)
	}
	if f.SOP == nil {
		t.Fatal("SOP is nil, want heredoc SOP")
	}
	if !strings.Contains(f.SOP.Content, "senior code reviewer") {
		t.Errorf("SOP.Content = %q", f.SOP.Content)
	}

	if len(f.RemoteAdds) != 2 {
		t.Fatalf("len(RemoteAdds) = %d, want 2: %+v", len(f.RemoteAdds), f.RemoteAdds)
	}
	skill := f.RemoteAdds[0]
	if skill.Dest != "/mnt/skills/code-review" || !skill.Cached || !skill.FailIfUnavailable {
		t.Errorf("RemoteAdds[0] (skill) = %+v", skill)
	}
	mcp := f.RemoteAdds[1]
	if mcp.Dest != "/mnt/mcp/github" || !mcp.Runtime || !mcp.FailIfUnavailable {
		t.Errorf("RemoteAdds[1] (mcp) = %+v", mcp)
	}
	if mcp.Source != "mcp://registry.internal.acme/servers/github:latest" {
		t.Errorf("RemoteAdds[1].Source = %q", mcp.Source)
	}

	wantPolicyKeys := map[string]bool{
		"model.name": false, "model.min_context": false, "model.fallback": false,
		"agent.idle_timeout": false, "agent.tool_timeout": false, "agent.max_retries": false,
		"network": false,
	}
	for _, p := range f.Policies {
		if _, ok := wantPolicyKeys[p.Key]; ok {
			wantPolicyKeys[p.Key] = true
		}
	}
	for k, found := range wantPolicyKeys {
		if !found {
			t.Errorf("expected POLICY %s, not found in %+v", k, f.Policies)
		}
	}

	stages := assertCleanedSourceIsPlainDockerfile(t, f.CleanedSource)
	if stages[0].BaseName != "ghcr.io/acme/pii-redacted-base:1.4" {
		t.Errorf("FROM = %q", stages[0].BaseName)
	}
}

func TestExtractSecureWorkspaceAndVaultAgent(t *testing.T) {
	for _, name := range []string{"Agentfile.secure-workspace", "Agentfile.vault-agent"} {
		t.Run(name, func(t *testing.T) {
			f, err := Extract(readExample(t, name))
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}
			if f.Identity["name"] == "" {
				t.Error("Identity[name] is empty")
			}
			if f.SOP == nil {
				t.Error("SOP is nil")
			}
			assertCleanedSourceIsPlainDockerfile(t, f.CleanedSource)

			issues := Validate(f)
			if HasErrors(issues) {
				t.Errorf("Validate() reported errors on a spec example: %v", issues)
			}
		})
	}
}

func TestExtractRejectsDuplicateSOP(t *testing.T) {
	cases := []string{
		"IDENTITY name=a\nSOP one\nSOP two\nCMD run\n",
		"IDENTITY name=a\nSOP <<EOF\nfirst\nEOF\nSOP <<EOF\nsecond\nEOF\nCMD run\n",
		"IDENTITY name=a\nSOP <<EOF\nfirst\nEOF\nSOP inline second\nCMD run\n",
	}
	for _, src := range cases {
		if _, err := Extract([]byte(src)); err == nil {
			t.Errorf("expected error for duplicate SOP in %q", src)
		}
	}
}

func TestExtractRejectsUnterminatedHeredoc(t *testing.T) {
	_, err := Extract([]byte("IDENTITY name=a\nSOP <<EOF\nno terminator\nCMD run\n"))
	if err == nil {
		t.Fatal("expected error for unterminated SOP heredoc")
	}
}

func TestExtractMultilineIdentityMerges(t *testing.T) {
	f, err := Extract([]byte("IDENTITY name=a version=1.0\nIDENTITY description=\"a thing\"\nSOP x\nCMD run\n"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	if f.Identity["name"] != "a" || f.Identity["version"] != "1.0" || f.Identity["description"] != "a thing" {
		t.Errorf("Identity = %+v", f.Identity)
	}
}

func TestExtractAddRemoteRequiresMntDest(t *testing.T) {
	f, err := Extract([]byte("IDENTITY name=a\nSOP x\nCMD run\nADD --remote https://example.com/x /opt/x\n"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	issues := Validate(f)
	if !HasErrors(issues) {
		t.Fatal("expected validation error for ADD --remote destination outside /mnt")
	}
}

func TestExtractAddRemoteRejectsPathTraversalOutOfMnt(t *testing.T) {
	// "/mnt/tools/../../etc/passwd" satisfies a naive strings.HasPrefix(dest,
	// "/mnt/") check but resolves outside /mnt entirely once cleaned.
	f, err := Extract([]byte("IDENTITY name=a\nSOP x\nCMD run\nADD --remote https://example.com/x /mnt/tools/../../etc/passwd\n"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	issues := Validate(f)
	if !HasErrors(issues) {
		t.Fatal("expected validation error for an ADD --remote destination that escapes /mnt via path traversal")
	}
}

func TestExtractPreservesLineNumbersForStandardInstructions(t *testing.T) {
	src := "IDENTITY name=a\nCAPABILITY text\nSOP x\n\nCMD run\n"
	f, err := Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	lines := strings.Split(string(f.CleanedSource), "\n")
	if len(lines) != 6 { // 5 content lines + trailing empty from final \n
		t.Fatalf("len(lines) = %d, want 6: %q", len(lines), lines)
	}
	if strings.TrimSpace(lines[4]) != "CMD run" {
		t.Errorf("CMD should remain on its original line (5): %q", lines)
	}
}
