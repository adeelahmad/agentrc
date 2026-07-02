package agentfile

import "testing"

func TestCleanMntPath(t *testing.T) {
	cases := []struct {
		dest    string
		wantOK  bool
		cleaned string
	}{
		{"/mnt", true, "/mnt"},
		{"/mnt/tools/x", true, "/mnt/tools/x"},
		{"/mnt/tools/../skills/x", true, "/mnt/skills/x"},
		{"/mnt/tools/../../etc/passwd", false, "/etc/passwd"},
		{"/mnt/../etc/passwd", false, "/etc/passwd"},
		{"/opt/x", false, "/opt/x"},
		{"/mntx/y", false, "/mntx/y"},
	}
	for _, c := range cases {
		cleaned, ok := CleanMntPath(c.dest)
		if ok != c.wantOK || cleaned != c.cleaned {
			t.Errorf("CleanMntPath(%q) = (%q, %v), want (%q, %v)", c.dest, cleaned, ok, c.cleaned, c.wantOK)
		}
	}
}

func TestResourceKindForDestRejectsTraversal(t *testing.T) {
	if kind := ResourceKindForDest("/mnt/tools/../../etc/passwd"); kind != KindOther {
		t.Errorf("ResourceKindForDest(traversal) = %q, want KindOther", kind)
	}
	if kind := ResourceKindForDest("/mnt/tools/file_read"); kind != KindTool {
		t.Errorf("ResourceKindForDest(/mnt/tools/file_read) = %q, want KindTool", kind)
	}
}
