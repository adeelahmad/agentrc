package agentfile

import (
	"net"
	"testing"
)

func TestValidateRemoteSourceURLRejectsInternalTargets(t *testing.T) {
	cases := []string{
		"http://169.254.169.254/latest/meta-data/",
		"http://localhost:8080/admin",
		"http://127.0.0.1/",
		"https://127.0.0.1/",
		"http://[::1]/",
		"http://10.0.0.5/",
		"http://192.168.1.1/",
		"http://0.0.0.0/",
		"ftp://example.com/file",
	}
	for _, u := range cases {
		if err := ValidateRemoteSourceURL(u); err == nil {
			t.Errorf("ValidateRemoteSourceURL(%q) = nil, want an error", u)
		}
	}
}

func TestValidateRemoteSourceURLAllowsPublicHTTPS(t *testing.T) {
	// A real public hostname that resolves to a public IP; this check
	// necessarily does a live DNS lookup, so skip rather than fail if this
	// environment has no network/DNS access at all.
	const publicURL = "https://github.com/"
	if _, err := net.LookupIP("github.com"); err != nil {
		t.Skipf("skipping (no DNS access in this environment): %v", err)
	}
	if err := ValidateRemoteSourceURL(publicURL); err != nil {
		t.Errorf("ValidateRemoteSourceURL(%q) = %v, want nil", publicURL, err)
	}
}
