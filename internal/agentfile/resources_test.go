package agentfile

import "testing"

func TestPopulateLocalResourcesFromExample(t *testing.T) {
	f, err := Extract(readExample(t, "Agentfile.minimal"))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	if err := PopulateLocalResources(f); err != nil {
		t.Fatalf("PopulateLocalResources() error = %v", err)
	}
	if len(f.LocalResources) != 1 {
		t.Fatalf("len(LocalResources) = %d, want 1: %+v", len(f.LocalResources), f.LocalResources)
	}
	r := f.LocalResources[0]
	if r.Kind != KindTool || r.Name != "file_read" || r.Dest != "/mnt/tools/file_read" || r.Source != "./tools/file_read" {
		t.Errorf("LocalResources[0] = %+v", r)
	}
}

func TestPopulateLocalResourcesClassifiesAllKinds(t *testing.T) {
	src := "IDENTITY name=a\nFROM python:3.11-slim\nCOPY ./tools/x /mnt/tools/x\nCOPY ./skills/y /mnt/skills/y\nCOPY ./mcp/z /mnt/mcp/z\nCOPY ./sop.md /mnt/SOP\nCOPY ./app.py /app/app.py\nCMD run\n"
	f, err := Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	if err := PopulateLocalResources(f); err != nil {
		t.Fatalf("PopulateLocalResources() error = %v", err)
	}

	want := map[string]ResourceKind{
		"/mnt/tools/x":  KindTool,
		"/mnt/skills/y": KindSkill,
		"/mnt/mcp/z":    KindMCP,
		"/mnt/SOP":      KindSOP,
	}
	got := map[string]ResourceKind{}
	for _, r := range f.LocalResources {
		got[r.Dest] = r.Kind
	}
	for dest, kind := range want {
		if got[dest] != kind {
			t.Errorf("resource at %s: kind = %q, want %q", dest, got[dest], kind)
		}
	}
	// /app/app.py is a totally ordinary COPY unrelated to agentrc resources
	// and must not be picked up.
	if _, ok := got["/app/app.py"]; ok {
		t.Error("non-/mnt COPY was incorrectly classified as a resource")
	}
	if len(f.LocalResources) != 4 {
		t.Errorf("len(LocalResources) = %d, want 4: %+v", len(f.LocalResources), f.LocalResources)
	}
}

func TestPopulateLocalResourcesSynthesizesFileBackedSOP(t *testing.T) {
	// No SOP keyword at all — SOP is declared purely via COPY .../mnt/SOP,
	// the spec's third documented form (spec/index.md §7).
	src := "IDENTITY name=a\nFROM python:3.11-slim\nCOPY ./sop.md /mnt/SOP\nCMD run\n"
	f, err := Extract([]byte(src))
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}
	if f.SOP != nil {
		t.Fatalf("Extract() should not itself populate a file-backed SOP: %+v", f.SOP)
	}
	if err := PopulateLocalResources(f); err != nil {
		t.Fatalf("PopulateLocalResources() error = %v", err)
	}
	if f.SOP == nil || !f.SOP.FileBacked {
		t.Fatalf("f.SOP = %+v, want a synthesized file-backed SOP", f.SOP)
	}
}

func TestPopulateLocalResourcesRejectsDuplicateSOPDeclarations(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{
			"keyword + local COPY",
			"IDENTITY name=a\nFROM python:3.11-slim\nSOP inline text\nCOPY ./sop.md /mnt/SOP\nCMD run\n",
		},
		{
			"keyword + remote ADD",
			"IDENTITY name=a\nFROM python:3.11-slim\nSOP inline text\nADD --remote https://example.com/sop.md /mnt/SOP\nCMD run\n",
		},
		{
			"local COPY + remote ADD",
			"IDENTITY name=a\nFROM python:3.11-slim\nCOPY ./sop.md /mnt/SOP\nADD --remote https://example.com/other-sop.md /mnt/SOP\nCMD run\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, err := Extract([]byte(c.src))
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}
			if err := PopulateLocalResources(f); err == nil {
				t.Fatal("expected an error for a duplicate SOP declaration")
			}
		})
	}
}
