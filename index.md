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
    <p class="lead"><strong>Like <code>bashrc</code> or <code>zshrc</code>, but for an agent.</strong> agentrc is an open specification for packaging one AI agent as a portable, governed artifact. An Agentfile declares the agent's identity, capabilities, system prompt, tools and skills, and its requests for models, resources, and network access — as typed policy a security team can review. The package is distributable through OCI-compatible registries; compatible runners decide how to execute it and enforce the declared boundaries. agentrc is not a runtime, cloud, model provider, or agent framework.</p>
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
FROM python:3.11-slim
IDENTITY name=hello version=0.1 author=acme
IDENTITY description="Minimal agentrc agent"
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

HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema
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

<section class="get-started">
  <div class="eyebrow">Get started</div>
  <h2 id="install" class="flush">Install the CLI</h2>
  <p class="lead">One binary — <code>agentrc</code> (alias <code>arc</code>). It scaffolds and builds Agentfiles, inspects what an agent requests, and translates an artifact into a backend's deploy config.</p>

  <div class="install-oneliner">
    <span class="install-label">macOS &amp; Linux</span>
    <pre><code>curl -fsSL https://agentrc.ai/install.sh | sh</code></pre>
  </div>

  <div class="grid">
    <div class="card">
      <h3>Homebrew</h3>
<pre><code>brew install \
  adeelahmad/tap/agentrc</code></pre>
    </div>
    <div class="card">
      <h3>Go 1.25+</h3>
<pre><code>go install \
  github.com/adeelahmad/agentrc/cmd/agentrc@latest</code></pre>
    </div>
    <div class="card">
      <h3>From source</h3>
<pre><code>git clone https://github.com/adeelahmad/agentrc
cd agentrc &amp;&amp; go build -o arc ./cmd/agentrc</code></pre>
    </div>
  </div>
  <p class="muted">Prebuilt binaries for macOS and Linux (amd64 / arm64), checksum-verified. Confirm with <code>arc version</code>. Prefer to read first? <code>curl -fsSL https://agentrc.ai/install.sh</code> and inspect it.</p>
</section>

<section class="get-started">
  <h2 id="build">Build and run — locally</h2>
  <p class="lead">Scaffold, validate, and compile an agent into a portable OCI artifact, then preview exactly what a local runner would execute.</p>
  <ol class="steps">
    <li><div><span class="step-k">Scaffold</span><pre><code>arc init            # writes ./Agentfile</code></pre></div></li>
    <li><div><span class="step-k">Validate</span><pre><code>arc lint Agentfile</code></pre></div></li>
    <li><div><span class="step-k">Build</span><pre><code>arc build -t ghcr.io/you/hello:0.1 .</code></pre></div></li>
    <li><div><span class="step-k">Preview the run</span><pre><code>arc run ghcr.io/you/hello:0.1 --backend local --dry-run</code></pre></div></li>
  </ol>
  <p class="muted"><code>arc build</code> produces a real OCI image (via <code>docker build</code> and the agentrc BuildKit frontend). <code>--dry-run</code> prints the config a runner would apply — agentrc declares and translates; it ships no runtime of its own.</p>
</section>

<section class="get-started">
  <h2 id="cloud">Ship the same artifact to the cloud</h2>
  <p class="lead">The build writes <code>org.agentrc.*</code> labels once. Point <code>arc run</code> at any backend to translate those labels into that platform's deploy form.</p>
  <div class="grid">
    <div class="card">
      <h3>Push once</h3>
<pre><code>arc push \
  ghcr.io/you/hello:0.1</code></pre>
    </div>
    <div class="card">
      <h3>AWS Bedrock</h3>
<pre><code>arc run ghcr.io/you/hello:0.1 \
  --backend bedrock --dry-run</code></pre>
      <p class="muted">→ CreateAgentRuntime JSON</p>
    </div>
    <div class="card">
      <h3>Kubernetes</h3>
<pre><code>arc run ghcr.io/you/hello:0.1 \
  --backend kubernetes --dry-run</code></pre>
      <p class="muted">→ deploy manifests</p>
    </div>
  </div>
  <div class="callout">
    <strong>Same artifact, same labels, three substrates.</strong> Reference translators — a proof of concept until platforms read <code>org.agentrc.*</code> labels natively. Not production runners.
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
| Execution substrate | run-time choice (`arc run --backend`) | local, container, microVM, cloud runners |

<div class="callout">
<strong>Core slogan:</strong> The Agentfile declares one agent. The lockfile pins dependencies. The package makes it portable. The policy makes boundaries reviewable. The registry makes it shareable. Compatible runners execute it.
</div>

## Current draft

<div class="pill-row">
  <span class="pill">Working Draft 0.1.0-draft.6</span>
  <span class="pill">Four keywords</span>
  <span class="pill">/mnt projection</span>
  <span class="pill">OCI labels</span>
  <span class="pill">Cedar enforcement</span>
  <span class="pill">Secrets deferred</span>
</div>

The project is published as a standards-style repository: specification first, reference tooling second.
