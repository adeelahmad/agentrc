---
layout: default
title: agentrc
description: "Agent Run Config: an open specification for portable, governed AI agents."
permalink: /
---
<section class="hero">
  <div>
    <div class="eyebrow">Agent Run Config</div>
    <h1>Portable, governed AI agents.</h1>
    <p class="lead">agentrc is an open specification for declaring, packaging, securing, and sharing AI agents. It defines the contract an agent declares; compatible runners decide how to execute it.</p>
    <div class="cta-row">
      <a class="button primary" href="{{ '/spec/' | relative_url }}">Read the specification</a>
      <a class="button" href="{{ '/docs/quickstart/' | relative_url }}">Start with an Agentfile</a>
      <a class="button" href="https://github.com/adeelahmad/agentrc" target="_blank" rel="noopener">View on GitHub</a>
      <a class="button" href="https://discord.gg/jWx6Qak5D" target="_blank" rel="noopener">Join the Discord</a>
    </div>
  </div>
  <div class="hero-card" markdown="1">
```dockerfile
# syntax=agentrc.agentfile/v0.1

AGENT code-reviewer
CMD claude --print

TOOL utcp:file_read
TOOL utcp:shell
MOUNT /workspace rw
CRED github_token env:GITHUB_TOKEN host:api.github.com
AUDIT all

POLICY
  permit(
    principal == AgentRC::Agent::"code-reviewer",
    action == AgentRC::Action::"tool.invoke",
    resource == AgentRC::Tool::"file_read"
  );
END
```
  </div>
</section>

<section class="grid">
  <div class="card">
    <h3>Not a runtime</h3>
    <p>agentrc does not implement containers, microVMs, cloud sandboxes, or model loops. It declares what must be true before a runner executes an agent.</p>
  </div>
  <div class="card">
    <h3>Security by declaration</h3>
    <p>Tools, mounts, network egress, credentials, rate limits, and audit requirements are declared in a reviewable file and pinned into a package.</p>
  </div>
  <div class="card">
    <h3>Registry-native</h3>
    <p>agentrc packages are designed for OCI-compatible registries, so agents, bases, tools, policies, and skills can be shared like container images.</p>
  </div>
</section>

## The separation agentrc creates

| Concern | Defined by | Implemented by |
|---|---|---|
| Agent identity and entrypoint | `Agentfile` | agent author |
| Tools, skills, functions, MCP servers | `Agentfile` + lockfile | package builder / runner |
| Security boundaries | policy and declarations | compatible runner |
| Packaging and sharing | agentrc package profile | OCI registry |
| Execution substrate | runner profile | Docker, gVisor, Firecracker, cloud runners, local runners |
| Multi-agent workflow | future workflow profile | workflow engines |

<div class="callout">
<strong>Core claim:</strong> The Agentfile declares one agent. The lockfile pins dependencies. The package makes it portable. The policy makes boundaries reviewable. The registry makes it shareable. Compatible runners execute it.
</div>

## Current draft

<div class="pill-row">
  <span class="pill">Working Draft 0.1</span>
  <span class="pill">Agentfile</span>
  <span class="pill">OCI package</span>
  <span class="pill">Cedar profile</span>
  <span class="pill">Runner conformance</span>
  <span class="pill">Workflow draft</span>
</div>

The project is ready to publish as a standards-style repository: specification first, reference tooling second.
