package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// bedrockJWTAuthorizer maps agent.auth.jwt.* labels onto the Bedrock
// AgentCore CustomJWTAuthorizer shape.
type bedrockJWTAuthorizer struct {
	DiscoveryUrl    string   `json:"discoveryUrl"`
	AllowedAudience []string `json:"allowedAudience,omitempty"`
	AllowedClients  []string `json:"allowedClients,omitempty"`
}

type bedrockS3Loc struct {
	Uri string `json:"uri"`
}

type bedrockCodeS3 struct {
	S3 bedrockS3Loc `json:"s3"`
}

// bedrockCodeConfig is emitted only in deployment.mode=code.
type bedrockCodeConfig struct {
	Runtime string        `json:"runtime"`
	Code    bedrockCodeS3 `json:"code"`
}

// bedrockRuntime is the CreateAgentRuntime request. Fields use the exact
// camelCase keys the Bedrock AgentCore API expects.
type bedrockRuntime struct {
	AgentRuntimeName          string                `json:"agentRuntimeName"`
	Description               string                `json:"description,omitempty"`
	ContainerUri              string                `json:"containerUri,omitempty"`
	RoleArn                   string                `json:"roleArn"`
	NetworkMode               string                `json:"networkMode,omitempty"`
	SecurityGroups            []string              `json:"securityGroups,omitempty"`
	Subnets                   []string              `json:"subnets,omitempty"`
	ServerProtocol            string                `json:"serverProtocol,omitempty"`
	EnvironmentVariables      map[string]string     `json:"environmentVariables,omitempty"`
	CustomJWTAuthorizer       *bedrockJWTAuthorizer `json:"customJWTAuthorizer,omitempty"`
	IdleRuntimeSessionTimeout string                `json:"idleRuntimeSessionTimeout,omitempty"`
	MaxLifetime               string                `json:"maxLifetime,omitempty"`
	CodeConfiguration         *bedrockCodeConfig    `json:"codeConfiguration,omitempty"`
}

// translateBedrock maps ai.agentrc.* labels + image config to a Bedrock
// AgentCore CreateAgentRuntime request (§8.7/§8.8/§8.9). It is fail-closed: on
// any unenforceable request it returns an error and EMPTY output rather than a
// partial config or an exposed invocation endpoint.
func translateBedrock(labels map[string]string) (string, error) {
	roleArn := labels["ai.agentrc.substrate.aws.roleArn"]
	if roleArn == "" {
		return "", fmt.Errorf("bedrock: substrate.aws.roleArn is required; refusing to emit CreateAgentRuntime config (fail-closed)")
	}

	rt := bedrockRuntime{
		AgentRuntimeName:          labels["ai.agentrc.identity.name"],
		Description:               labels["ai.agentrc.identity.description"],
		ContainerUri:              labels["image.ref"],
		RoleArn:                   roleArn,
		NetworkMode:               labels["ai.agentrc.substrate.aws.networkMode"],
		SecurityGroups:            splitCSV(labels["ai.agentrc.substrate.aws.securityGroup"]),
		Subnets:                   splitCSV(labels["ai.agentrc.substrate.aws.subnet"]),
		ServerProtocol:            labels["ai.agentrc.substrate.aws.protocol"],
		IdleRuntimeSessionTimeout: labels["ai.agentrc.agent.idle_timeout"],
		MaxLifetime:               labels["ai.agentrc.substrate.aws.maxLifetime"],
	}

	if env := envVars(labels); len(env) > 0 {
		rt.EnvironmentVariables = env
	}

	// §8.8 fail-closed: a jwt authorizer that cannot be enforced (no discovery
	// URL) MUST NOT expose the invocation endpoint.
	if labels["ai.agentrc.agent.auth.mode"] == "jwt" {
		disc := labels["ai.agentrc.agent.auth.jwt.discovery_url"]
		if disc == "" {
			return "", fmt.Errorf("bedrock: agent.auth.mode=jwt requires agent.auth.jwt.discovery_url; refusing to expose an invocation endpoint (fail-closed)")
		}
		rt.CustomJWTAuthorizer = &bedrockJWTAuthorizer{
			DiscoveryUrl:    disc,
			AllowedAudience: splitCSV(labels["ai.agentrc.agent.auth.jwt.allowed_audience"]),
			AllowedClients:  splitCSV(labels["ai.agentrc.agent.auth.jwt.allowed_client"]),
		}
	}

	// §8.9 fail-closed: code mode requires a resolvable runtime language;
	// refuse to guess rather than emit a config that cannot run.
	if labels["ai.agentrc.substrate.aws.deployment.mode"] == "code" {
		s3uri := labels["ai.agentrc.substrate.aws.code.s3.uri"]
		lang := labels["ai.agentrc.substrate.runtime.language"]
		if s3uri != "" && lang == "" {
			return "", fmt.Errorf("bedrock: deployment.mode=code with code.s3.uri requires a resolvable substrate.runtime.language; refusing to guess a runtime (fail-closed)")
		}
		rt.CodeConfiguration = &bedrockCodeConfig{
			Runtime: lang,
			Code:    bedrockCodeS3{S3: bedrockS3Loc{Uri: s3uri}},
		}
	}

	b, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("bedrock: marshaling CreateAgentRuntime: %w", err)
	}
	return string(b) + "\n", nil
}

// splitCSV splits a repeatable request value comma-joined into one label value
// (e.g. securityGroup, subnet, allowed_audience, allowed_client) back into a
// trimmed, non-empty slice.
func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// envVars collects CLI-injected env.<NAME> labels into a name→value map.
func envVars(labels map[string]string) map[string]string {
	env := map[string]string{}
	for k, v := range labels {
		if name := strings.TrimPrefix(k, "env."); name != k {
			env[name] = v
		}
	}
	return env
}
