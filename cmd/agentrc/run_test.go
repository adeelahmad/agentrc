package main

import (
	"bytes"
	"strings"
	"testing"
)

// executeRun builds a fresh run command, wires stdout/stderr to a single buffer,
// sets args, and executes. It returns the combined output and the execute error.
// All tests reference only symbols that exist today (newRunCmd) so the file
// COMPILES against the current run.go; they FAIL on assertions until T21 lands.
func executeRun(args ...string) (string, error) {
	cmd := newRunCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

// TestBackendFlagReplacesSubstrate asserts the post-T21 help surface: `--backend`
// is documented and `--substrate` is gone. FAILS now: run.go:19 declares
// `--substrate` and no `--backend` flag exists, so --help shows the opposite.
func TestBackendFlagReplacesSubstrate(t *testing.T) {
	out, err := executeRun("--help")
	if err != nil {
		t.Fatalf("run --help returned error: %v", err)
	}
	if !strings.Contains(out, "--backend") {
		t.Errorf("help output should document --backend, got:\n%s", out)
	}
	if strings.Contains(out, "--substrate") {
		t.Errorf("help output must NOT mention --substrate (CLI-flag rename), got:\n%s", out)
	}
}

// TestBackendDefaultsToLocal asserts the `--backend` flag exists and defaults to
// "local". FAILS now: the flag is absent (Lookup returns nil).
func TestBackendDefaultsToLocal(t *testing.T) {
	f := newRunCmd().Flags().Lookup("backend")
	if f == nil {
		t.Fatal("expected --backend flag to be registered")
	}
	if f.DefValue != "local" {
		t.Errorf("expected --backend default %q, got %q", "local", f.DefValue)
	}
}

// TestBackendRejectsUnknownValue asserts an unknown backend value is rejected with
// an error naming the valid values. FAILS now: `--backend` is not a known flag, so
// cobra returns "unknown flag: --backend" which does NOT name the valid values.
func TestBackendRejectsUnknownValue(t *testing.T) {
	_, err := executeRun("someref", "--backend", "bogus")
	if err == nil {
		t.Fatal("expected error for unknown --backend value")
	}
	msg := err.Error()
	for _, valid := range []string{"local", "bedrock", "kubernetes"} {
		if !strings.Contains(msg, valid) {
			t.Errorf("error should name valid backend %q, got: %v", valid, msg)
		}
	}
}

// TestIsolationScopedToLocalBackend asserts `--isolation` is only meaningful for
// `--backend local`: using it with a non-local backend must error and the message
// must reference isolation/local scoping. FAILS now: `--backend` is unknown so the
// error is "unknown flag: --backend", not the isolation-scoping validation error.
func TestIsolationScopedToLocalBackend(t *testing.T) {
	_, err := executeRun("someref", "--backend", "bedrock", "--isolation", "microvm")
	if err == nil {
		t.Fatal("expected error: --isolation is only valid for --backend local")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "isolation") || !strings.Contains(msg, "local") {
		t.Errorf("error should explain --isolation is scoped to backend local, got: %v", err)
	}
}

// TestPerBackendFlagsParse asserts the per-backend flags land on the run command:
// bedrock (--region, --profile), kubernetes (--kubeconfig, --namespace), and the
// shared --dry-run. FAILS now: none of these flags are registered.
func TestPerBackendFlagsParse(t *testing.T) {
	flags := newRunCmd().Flags()
	for _, name := range []string{"region", "profile", "kubeconfig", "namespace", "dry-run"} {
		if flags.Lookup(name) == nil {
			t.Errorf("expected --%s flag to be registered", name)
		}
	}
}

// TestDryRunExitsAfterPrint asserts `--dry-run` prints the translated config and
// exits 0 WITHOUT the hard "not implemented" runtime error. FAILS now: `--backend`
// is unknown (flag parse error) and RunE still returns the "not implemented" error.
//
// GREEN seam: run.go should dispatch dry-run through a pure
//
//	translate(backend string, labels map[string]string) (string, error)
//
// stub (S2-04 fills the per-backend body). This test asserts the observable
// invariant that dry-run succeeds (err == nil) rather than reaching a runtime.
func TestDryRunExitsAfterPrint(t *testing.T) {
	_, err := executeRun("someref", "--backend", "bedrock", "--dry-run")
	if err != nil {
		t.Fatalf("--dry-run should exit 0 (print translated config, no runtime), got: %v", err)
	}
}
