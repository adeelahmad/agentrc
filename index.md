---
layout: default
title: agentrc
description: "Agent Run Config: an open, Dockerfile-shaped specification for portable, governed AI agents."
permalink: /
---
<section class="hero">
  <div>
    <div class="eyebrow">Agent Run Config</div>
    <h1>Portable, governed AI agents.</h1>
    <p class="lead">agentrc is an open, Dockerfile-shaped specification for declaring, packaging, securing, and sharing AI agents. Four new keywords over the Dockerfile you already know. The build emits OCI labels; the platform reads the labels and decides what to honour.</p>
    <div class="cta-row">
      <a class="button primary" href="{{ '/spec/' | relative_url }}">Read the specification</a>
      <a class="button" href="{{ '/docs/quickstart/' | relative_url }}">Start with an Agentfile</a>
      <a class="button" href="https://github.com/adeelahmad/agentrc" target="_blank" rel="noopener">View on GitHub</a>
      <a class="button" href="https://discord.gg/jWx6Qak5D" target="_blank" rel="noopener">Join the Discord</a>
    </div>
  </div>
  <div class="hero-card" markdown="1">
```dockerfile
# syntax=agentrc.io/agentfile:v1
FROM ghcr.io/acme/pii-redacted-base:1.4

IDENTITY name=code-reviewer version=1.0 author=acme
CAPABILITY text
CAPABILITY streaming

SOP <<EOF
Review pull requests for correctness and security.
Escalate anything ambiguous to a human.
EOF

CMD claude --print

COPY --chmod=755 ./tools/file_read /mnt/tools/file_read
ADD --remote --cached --fail-if-unavailable \
    https://registry.agentrc.io/skills/code-review:1.2.3 /mnt/skills/code-review

LABEL org.agentrc.secret.github_token=host:api.github.com

POLICY model.name      claude-opus-4
POLICY model.fallback  claude-sonnet-4
POLICY agent.tool_timeout 30s
POLICY network dns:api.github.com:443
```
  </div>
</section>

<section class="grid">
  <div class="card">
    <h3>Dockerfile-shaped</h3>
    <p>Reuse standard Dockerfile keywords for everything they can express. agentrc adds just four: <code>IDENTITY</code>, <code>CAPABILITY</code>, <code>SOP</code>, and <code>POLICY</code>. Build with <code>docker build</code> or <code>arc build</code>.</p>
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
| Secrets | `LABEL org.agentrc.secret.*` | platform secret broker |
| Enforcement | typed requests compiled to Cedar | platform (deny-by-default, `forbid` > `permit`) |
| Packaging and sharing | OCI artifact + `org.agentrc.*` labels | any OCI registry |
| Execution substrate | run-time choice (`arc run --substrate`) | local, container, microVM, cloud runners |

<div class="callout">
<strong>Core claim:</strong> The Agentfile is the human-authored recipe. The build emits OCI labels. The platform reads the labels and decides what to honour — granting, narrowing, or rejecting each request and enforcing the result with Cedar. The agent never carries a secret, and never writes a policy language.
</div>

## Current draft

<div class="pill-row">
  <span class="pill">Working Draft v1</span>
  <span class="pill">Four keywords</span>
  <span class="pill">/mnt projection</span>
  <span class="pill">OCI labels</span>
  <span class="pill">Cedar enforcement</span>
  <span class="pill">BuildKit frontend + arc CLI</span>
</div>

The project is published as a standards-style repository: specification first, reference tooling second.
