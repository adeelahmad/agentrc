package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// fullBedrockLabels returns a label set that populates every CreateAgentRuntime
// field, so the mapping test can assert 13/13. Callers mutate/delete keys to
// exercise the fail-closed paths. Keys follow the emitted-label contract:
//   - org.agentrc.identity.*       → agentRuntimeName / description
//   - image.ref (CLI-injected)     → containerUri
//   - env.<NAME> (CLI-injected)    → environmentVariables
//   - org.agentrc.substrate.aws.*  → roleArn/networkMode/securityGroups/subnets/
//     serverProtocol/maxLifetime/codeConfiguration
//   - org.agentrc.substrate.runtime.language → codeConfiguration.runtime
//   - org.agentrc.agent.auth.jwt.* → customJWTAuthorizer
//   - org.agentrc.agent.idle_timeout → idleRuntimeSessionTimeout
//
// Repeatable requests (securityGroup, subnet, allowed_audience, allowed_client)
// are comma-joined in a single label value.
func fullBedrockLabels() map[string]string {
	return map[string]string{
		"org.agentrc.identity.name":                   "code-reviewer",
		"org.agentrc.identity.description":            "Reviews pull requests",
		"image.ref":                                   "123456789012.dkr.ecr.us-east-1.amazonaws.com/code-reviewer:1.0",
		"org.agentrc.substrate.aws.roleArn":           "arn:aws:iam::123456789012:role/agent-exec",
		"org.agentrc.substrate.aws.networkMode":       "PUBLIC",
		"org.agentrc.substrate.aws.securityGroup":     "sg-0abc123,sg-0def456",
		"org.agentrc.substrate.aws.subnet":            "subnet-0abc123,subnet-0def456",
		"org.agentrc.substrate.aws.protocol":          "HTTP",
		"org.agentrc.substrate.aws.maxLifetime":       "1h",
		"org.agentrc.substrate.aws.deployment.mode":   "code",
		"org.agentrc.substrate.aws.code.s3.uri":       "s3://acme-agents/code-reviewer.zip",
		"org.agentrc.substrate.runtime.language":      "python:3.11",
		"env.LOG_LEVEL":                               "info",
		"org.agentrc.agent.idle_timeout":              "5m",
		"org.agentrc.agent.auth.mode":                 "jwt",
		"org.agentrc.agent.auth.jwt.discovery_url":    "https://auth.acme/.well-known/openid-configuration",
		"org.agentrc.agent.auth.jwt.allowed_audience": "agentrc://code-reviewer",
		"org.agentrc.agent.auth.jwt.allowed_client":   "acme-ci,acme-bot",
	}
}

// TestBedrockMapsAllThirteenFields asserts the full label set maps to all 13
// CreateAgentRuntime fields in the emitted JSON. FAILS now: translate is a stub
// that emits none of these fields.
func TestBedrockMapsAllThirteenFields(t *testing.T) {
	out, err := translate("bedrock", fullBedrockLabels())
	if err != nil {
		t.Fatalf("translate(bedrock) with a full label set should not error: %v", err)
	}
	var got map[string]json.RawMessage
	if uerr := json.Unmarshal([]byte(out), &got); uerr != nil {
		t.Fatalf("bedrock output is not valid JSON: %v\n%s", uerr, out)
	}
	fields := []string{
		"agentRuntimeName",
		"description",
		"containerUri",
		"roleArn",
		"networkMode",
		"securityGroups",
		"subnets",
		"serverProtocol",
		"environmentVariables",
		"customJWTAuthorizer",
		"idleRuntimeSessionTimeout",
		"maxLifetime",
		"codeConfiguration",
	}
	for _, f := range fields {
		raw, ok := got[f]
		if !ok {
			t.Errorf("missing CreateAgentRuntime field %q in bedrock JSON", f)
			continue
		}
		if len(raw) == 0 || string(raw) == "null" || string(raw) == `""` || string(raw) == "[]" || string(raw) == "{}" {
			t.Errorf("CreateAgentRuntime field %q is present but empty: %s", f, raw)
		}
	}
}

// TestBedrockDryRunEmitsValidJSON asserts the dry-run output unmarshals as JSON
// and carries the top-level agentRuntimeName field (proof it is real translated
// Bedrock config, not the placeholder stub). FAILS now: stub JSON has no
// agentRuntimeName.
func TestBedrockDryRunEmitsValidJSON(t *testing.T) {
	out, err := translate("bedrock", fullBedrockLabels())
	if err != nil {
		t.Fatalf("translate(bedrock) should not error: %v", err)
	}
	var got map[string]json.RawMessage
	if uerr := json.Unmarshal([]byte(out), &got); uerr != nil {
		t.Fatalf("bedrock --dry-run output must be valid JSON: %v\n%s", uerr, out)
	}
	if _, ok := got["agentRuntimeName"]; !ok {
		t.Errorf("bedrock JSON should carry agentRuntimeName (real CreateAgentRuntime config), got keys: %v", keysOf(got))
	}
}

// TestBedrockDryRunFailsClosedWithoutRoleArn asserts that a label set missing
// substrate.aws.roleArn fails closed: a non-nil error AND no config emitted.
// FAILS now: stub returns placeholder JSON and a nil error.
func TestBedrockDryRunFailsClosedWithoutRoleArn(t *testing.T) {
	labels := fullBedrockLabels()
	delete(labels, "org.agentrc.substrate.aws.roleArn")

	out, err := translate("bedrock", labels)
	if err == nil {
		t.Fatal("bedrock must fail closed without substrate.aws.roleArn, got nil error")
	}
	if out != "" {
		t.Errorf("bedrock must emit NO config when roleArn is missing, got:\n%s", out)
	}
}

// TestBedrockFailsClosedOnUnenforceableJWT asserts that agent.auth.mode=jwt
// without a resolvable discovery_url fails closed: an error and no invocation
// endpoint (in fact no config at all). FAILS now: stub never errors.
func TestBedrockFailsClosedOnUnenforceableJWT(t *testing.T) {
	labels := fullBedrockLabels()
	delete(labels, "org.agentrc.agent.auth.jwt.discovery_url")

	out, err := translate("bedrock", labels)
	if err == nil {
		t.Fatal("bedrock must fail closed on jwt mode without discovery_url, got nil error")
	}
	if out != "" {
		t.Errorf("bedrock must emit NO config/endpoint on unenforceable jwt, got:\n%s", out)
	}
	if strings.Contains(strings.ToLower(out), "endpoint") {
		t.Errorf("bedrock must NOT emit an invocation endpoint on unenforceable jwt, got:\n%s", out)
	}
}

// TestBedrockFailsClosedCodeModeWithoutLanguage asserts that deployment.mode=code
// with a code.s3.uri but no resolvable substrate.runtime.language fails closed
// (§8.9). FAILS now: stub never errors.
func TestBedrockFailsClosedCodeModeWithoutLanguage(t *testing.T) {
	labels := fullBedrockLabels()
	delete(labels, "org.agentrc.substrate.runtime.language")

	out, err := translate("bedrock", labels)
	if err == nil {
		t.Fatal("bedrock must fail closed in code mode without substrate.runtime.language, got nil error")
	}
	if out != "" {
		t.Errorf("bedrock must emit NO config in code mode without a runtime language, got:\n%s", out)
	}
}

// TestBedrockJWTAuthorizerFromAuthLabels asserts valid agent.auth.jwt.* labels
// populate customJWTAuthorizer with the discovery URL, allowed audiences, and
// allowed clients. FAILS now: stub emits no customJWTAuthorizer.
func TestBedrockJWTAuthorizerFromAuthLabels(t *testing.T) {
	out, err := translate("bedrock", fullBedrockLabels())
	if err != nil {
		t.Fatalf("translate(bedrock) should not error with valid jwt labels: %v", err)
	}
	var got struct {
		CustomJWTAuthorizer struct {
			DiscoveryUrl    string   `json:"discoveryUrl"`
			AllowedAudience []string `json:"allowedAudience"`
			AllowedClients  []string `json:"allowedClients"`
		} `json:"customJWTAuthorizer"`
	}
	if uerr := json.Unmarshal([]byte(out), &got); uerr != nil {
		t.Fatalf("bedrock output is not valid JSON: %v\n%s", uerr, out)
	}
	a := got.CustomJWTAuthorizer
	if a.DiscoveryUrl != "https://auth.acme/.well-known/openid-configuration" {
		t.Errorf("customJWTAuthorizer.discoveryUrl = %q, want the OIDC discovery URL", a.DiscoveryUrl)
	}
	if len(a.AllowedAudience) == 0 {
		t.Error("customJWTAuthorizer.allowedAudience should be populated from agent.auth.jwt.allowed_audience")
	}
	if len(a.AllowedClients) < 2 {
		t.Errorf("customJWTAuthorizer.allowedClients should carry both comma-joined clients, got %v", a.AllowedClients)
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
