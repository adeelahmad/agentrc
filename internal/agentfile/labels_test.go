package agentfile

import "testing"

func TestBuildLabelsMinimalExample(t *testing.T) {
	f, err := Extract(readExample(t, "Agentfile.minimal"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	// Simulate what internal/llb does after parsing CleanedSource for real
	// COPY/ADD instructions.
	f.LocalResources = []LocalResource{
		{Source: "./tools/file_read", Dest: "/mnt/tools/file_read", Kind: KindTool, Name: "file_read"},
	}

	labels, err := BuildLabels(f, LabelOptions{})
	if err != nil {
		t.Fatalf("BuildLabels() error = %v", err)
	}

	want := map[string]string{
		"org.agentrc.identity.name":               "hello",
		"org.agentrc.identity.version":            "0.1",
		"org.agentrc.capability.text":             "true",
		"org.agentrc.sop":                         "/mnt/SOP",
		"org.agentrc.tool.file_read":              "local",
		"org.agentrc.model.name":                  "claude-sonnet-4",
		"org.agentrc.agent.tool_timeout":          "30s",
		"org.agentrc.network.dns.api.example.com": "443",
	}
	for k, v := range want {
		if labels[k] != v {
			t.Errorf("labels[%q] = %q, want %q", k, labels[k], v)
		}
	}
	if labels["org.agentrc.sop.sha256"] == "" {
		t.Error("org.agentrc.sop.sha256 is empty")
	}
}

func TestBuildLabelsRemoteResources(t *testing.T) {
	f, err := Extract(readExample(t, "Agentfile.code-reviewer"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	f.LocalResources = []LocalResource{
		{Dest: "/mnt/tools/file_read", Kind: KindTool, Name: "file_read"},
	}
	// Simulate internal/llb resolving the --cached skill's digest but
	// leaving the --runtime mcp server unresolved (as spec requires).
	for i := range f.RemoteAdds {
		if f.RemoteAdds[i].Cached {
			f.RemoteAdds[i].ResolvedDigest = "sha256:deadbeef"
		}
	}

	labels, err := BuildLabels(f, LabelOptions{})
	if err != nil {
		t.Fatalf("BuildLabels() error = %v", err)
	}

	if labels["org.agentrc.skill.code-review"] != "sha256:deadbeef" {
		t.Errorf("skill label = %q", labels["org.agentrc.skill.code-review"])
	}
	if labels["org.agentrc.skill.code-review.origin"] != "https://registry.agentrc.io/skills/code-review:1.2.3" {
		t.Errorf("skill origin label = %q", labels["org.agentrc.skill.code-review.origin"])
	}
	if labels["org.agentrc.mcp.github"] != "runtime:mcp://registry.internal.acme/servers/github:latest" {
		t.Errorf("mcp label = %q", labels["org.agentrc.mcp.github"])
	}
	if _, ok := labels["org.agentrc.mcp.github.origin"]; ok {
		t.Error("runtime resources must not get a resolved digest+origin pair")
	}
}

func TestBuildLabelsAutoDerivedEgress(t *testing.T) {
	f := &File{
		Identity: Identity{"name": "a"},
		Policies: []Policy{
			{Key: "agent.hooks.pre", Value: "https://hooks.internal.example/pre-step"},
		},
	}
	labels, err := BuildLabels(f, LabelOptions{})
	if err != nil {
		t.Fatalf("BuildLabels() error = %v", err)
	}
	if labels["org.agentrc.agent.hooks.pre"] != "https://hooks.internal.example/pre-step" {
		t.Errorf("hooks.pre label = %q", labels["org.agentrc.agent.hooks.pre"])
	}
	if labels["org.agentrc.network.dns.hooks.internal.example"] != "443" {
		t.Errorf("auto-derived egress label = %q, labels=%v", labels["org.agentrc.network.dns.hooks.internal.example"], labels)
	}
	if labels["org.agentrc.network.dns.hooks.internal.example.source"] != "auto:agent.hooks.pre" {
		t.Errorf("auto-derived egress source label = %q", labels["org.agentrc.network.dns.hooks.internal.example.source"])
	}
}

func TestBuildLabelsSOPFileBackedRequiresDigest(t *testing.T) {
	f := &File{
		Identity: Identity{"name": "a"},
		SOP:      &SOP{FileBacked: true},
	}
	if _, err := BuildLabels(f, LabelOptions{}); err == nil {
		t.Fatal("expected error when a file-backed SOP has no resolved digest")
	}
	labels, err := BuildLabels(f, LabelOptions{SOPSHA256: "sha256:abc"})
	if err != nil {
		t.Fatalf("BuildLabels() error = %v", err)
	}
	if labels["org.agentrc.sop.sha256"] != "sha256:abc" {
		t.Errorf("sop.sha256 = %q", labels["org.agentrc.sop.sha256"])
	}
}

func TestBuildOCIAnnotations(t *testing.T) {
	f := &File{Identity: Identity{"name": "hello", "version": "1.0", "author": "acme"}}
	anno := BuildOCIAnnotations(f, LabelOptions{Created: "2026-07-01T00:00:00Z"})
	if anno["org.opencontainers.image.title"] != "hello" {
		t.Errorf("title = %q", anno["org.opencontainers.image.title"])
	}
	if anno["org.opencontainers.image.version"] != "1.0" {
		t.Errorf("version = %q", anno["org.opencontainers.image.version"])
	}
	if anno["org.opencontainers.image.created"] != "2026-07-01T00:00:00Z" {
		t.Errorf("created = %q", anno["org.opencontainers.image.created"])
	}
}
