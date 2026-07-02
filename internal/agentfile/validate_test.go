package agentfile

import "testing"

func TestValidateRequiresIdentityName(t *testing.T) {
	f := &File{Identity: Identity{}, SOP: &SOP{Content: "x"}, Capabilities: []string{"text"}}
	issues := Validate(f)
	if !HasErrors(issues) {
		t.Fatal("expected error for missing IDENTITY name")
	}
}

func TestValidateWarnsOnMissingCapabilityAndSOP(t *testing.T) {
	f := &File{Identity: Identity{"name": "a"}}
	issues := Validate(f)
	if HasErrors(issues) {
		t.Fatalf("missing CAPABILITY/SOP should warn, not error: %v", issues)
	}
	if len(issues) != 2 {
		t.Fatalf("expected 2 warnings (capability, sop), got %v", issues)
	}
}

func TestValidateNetworkPolicyShape(t *testing.T) {
	cases := []struct {
		value   string
		wantErr bool
	}{
		{"dns:api.github.com:443", false},
		{"dns:internal.acme:*", false},
		{"api.github.com:443", true}, // missing dns: prefix
		{"dns:api.github.com", true}, // missing port
		{"dns::443", true},           // empty host
	}
	for _, c := range cases {
		f := &File{
			Identity:     Identity{"name": "a"},
			Capabilities: []string{"text"},
			SOP:          &SOP{Content: "x"},
			Policies:     []Policy{{Key: "network", Value: c.value, Line: 1}},
		}
		issues := Validate(f)
		if got := HasErrors(issues); got != c.wantErr {
			t.Errorf("network policy %q: HasErrors() = %v, want %v (%v)", c.value, got, c.wantErr, issues)
		}
	}
}

func TestValidateWarnsOnUnknownPolicyNamespace(t *testing.T) {
	f := &File{
		Identity:     Identity{"name": "a"},
		Capabilities: []string{"text"},
		SOP:          &SOP{Content: "x"},
		Policies:     []Policy{{Key: "totallyunknown.thing", Value: "1", Line: 1}},
	}
	issues := Validate(f)
	if HasErrors(issues) {
		t.Fatalf("unknown namespace should warn (extensible), not error: %v", issues)
	}
	found := false
	for _, i := range issues {
		if i.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected a warning for the unrecognized policy namespace")
	}
}

func TestValidateKnownPolicyNamespacesDoNotWarn(t *testing.T) {
	f := &File{
		Identity:     Identity{"name": "a"},
		Capabilities: []string{"text"},
		SOP:          &SOP{Content: "x"},
		Policies: []Policy{
			{Key: "agent.idle_timeout", Value: "5m", Line: 1},
			{Key: "substrate.runtime.memory", Value: "8gb", Line: 2},
			{Key: "model.name", Value: "claude-opus-4", Line: 3},
		},
	}
	issues := Validate(f)
	if len(issues) != 0 {
		t.Errorf("expected no issues for well-formed known-namespace policies, got %v", issues)
	}
}
