---
layout: default
title: agentrc
description: "Agent Run Config: an open specification for packaging one AI agent as a portable, governed artifact."
permalink: /
---
<section class="hero">
  <div>
    <div class="eyebrow">Agent Run Config</div>
    <h1>Portable, governed AI agents.</h1>
    <p class="lead">agentrc is an open specification for packaging one AI agent as a portable, governed artifact. An Agentfile declares the agent's identity, capabilities, system prompt, tools and skills, and its requests for models, resources, and network access — as typed policy a security team can review. The package is distributable through OCI-compatible registries; compatible runners decide how to execute it and enforce the declared boundaries. agentrc is not a runtime, cloud, model provider, or agent framework.</p>
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
IDENTITY name=hello version=0.1 author=acme
IDENTITY description="Minimal AgentRC agent"
CAPABILITY text
SOP You are a minimal example agent. Read a file when asked; do nothing else.
CMD python ./agent.py

# Tool (local, embedded) — projected under /mnt/tools/
COPY --chmod=755 ./tools/file_read /mnt/tools/file_read

# Model + operational requests (platform grants, narrows, or rejects)
POLICY model.name         claude-sonnet-4
POLICY agent.tool_timeout 30s

# Network egress request
POLICY network dns:api.example.com:443

HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/ping
```
  </div>
</section>

<section class="grid">
  <div class="card">
    <h3>Dockerfile-shaped</h3>
    <p>Reuse standard Dockerfile keywords for everything they can express. agentrc adds just four agent-native ones: <code>IDENTITY</code>, <code>CAPABILITY</code>, <code>SOP</code>, and <code>POLICY</code>. Build with <code>docker build</code> or <code>arc build</code>.</p>
  </div>
  <div class="card">
    <h3>Requests, not enforcement</h3>
    <p>A <code>POLICY</code> line asks the platform for a model, a resource, or a constraint. The platform grants, narrows, or rejects it — and enforces the decision with Cedar, deny-by-default.</p>
  </div>
  <div class="card">
    <h3>Labels are the manifest</h3>
    <p>The build translates intent into namespaced <code>org.agentrc.*</code> OCI labels. The platform reads the labels — never the Agentfile — so agents ship, sign, and mirror like any container image.</p>
  </div>
</section>

## The separation agentrc creates

| Concern | Defined by | Read / enforced by |
|---|---|---|
| Agent identity, capabilities, objective | `IDENTITY` / `CAPABILITY` / `SOP` | agent author |
| Tools, skills, MCP servers | `COPY` / `ADD --remote` into `/mnt` | compiler → layers + labels |
| Resource, model, network, lifecycle requests | `POLICY` (typed namespaces) | platform (grant / narrow / reject) |
| Enforcement | typed requests compiled to Cedar | platform (deny-by-default, `forbid` > `permit`) |
| Packaging and sharing | OCI artifact + `org.agentrc.*` labels | any OCI registry |
| Execution substrate | run-time choice (`arc run --substrate`) | local, container, microVM, cloud runners |

<div class="callout">
<strong>Core slogan:</strong> The Agentfile declares one agent. The lockfile pins dependencies. The package makes it portable. The policy makes boundaries reviewable. The registry makes it shareable. Compatible runners execute it.
</div>

## Current draft

<div class="pill-row">
  <span class="pill">Working Draft 0.1.0-draft.5</span>
  <span class="pill">Four keywords</span>
  <span class="pill">/mnt projection</span>
  <span class="pill">OCI labels</span>
  <span class="pill">Cedar enforcement</span>
  <span class="pill">Secrets deferred</span>
</div>

The project is published as a standards-style repository: specification first, reference tooling second.
