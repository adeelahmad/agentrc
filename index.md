---
layout: default
title: agentrc
description: "Agent Run Config: an open specification for packaging one AI agent as a portable, governed artifact."
permalink: /
---
<section class="hero">
  <div>
    <div class="eyebrow">Open specification</div>
    <h1 style="margin:.3rem 0 .8rem;min-height:2.3em"><span id="rotator">portable, governed ai agents</span><span class="term-cursor" aria-hidden="true">_</span></h1>
    <p class="lead"><strong style="color:var(--text-strong)">Like <code>bashrc</code> or <code>zshrc</code>, but for an agent.</strong> An Agentfile declares one AI agent's identity, capabilities, system prompt, and tools, plus its requests for models, resources, and network — as typed policy a security team can review. Package it as an OCI artifact; compatible runners execute and enforce it. Not a runtime, cloud, model provider, or agent framework.</p>
    <div class="cta-row">
      <a class="button primary" href="{{ '/spec/' | relative_url }}">Read the specification →</a>
      <a class="button" href="{{ '/docs/quickstart/' | relative_url }}">Start with an Agentfile</a>
    </div>
    <div class="cta-row" style="margin-top:.7rem">
      <a class="button" href="https://github.com/adeelahmad/agentrc" target="_blank" rel="noopener">View on GitHub</a>
      <a class="button" href="https://discord.gg/jWx6Qak5D" target="_blank" rel="noopener">Join the Discord</a>
    </div>
  </div>
  <div class="term-col">
    <pre class="ascii-rain" aria-hidden="true">01 &gt; { } _ /mnt
=&gt; POLICY 10
arc build ..
0x1f 4a2b &gt;_
IDENTITY ::
{ agent } //
CAPABILITY 1
SOP -&gt; run
label 0.1.0
$ arc lint _
network:443
grant|narrow
oci://ghcr..
04a2 &gt; 1101
deny-default
&gt;_ cedar ok</pre>
    <div class="term-window">
      <div class="term-bar">
        <span class="term-dots"><i></i><i></i><i></i></span>
        <span class="term-title">agentrc:~$ <b id="tcmd">cat Agentfile</b></span>
      </div>
      <div class="term-body" id="tbody"><pre><code><span class="c"># syntax=agentrc.agentfile/v0.1</span>
<span class="k">FROM</span> python:3.11-slim
<span class="k">IDENTITY</span> name=hello version=0.1 author=acme
<span class="k">IDENTITY</span> description=<span class="s">"Minimal agentrc agent"</span>
<span class="k">CAPABILITY</span> text
<span class="k">SOP</span> You are a minimal example agent. Read a
    file when asked; do nothing else.
<span class="k">CMD</span> python ./agent.py

<span class="c"># Tool (local, embedded) → /mnt/tools/</span>
<span class="k">COPY</span> --chmod=755 ./tools/file_read /mnt/tools/file_read

<span class="c"># Requests: the platform grants, narrows, or rejects</span>
<span class="k">POLICY</span> model.name         claude-sonnet-4
<span class="k">POLICY</span> agent.tool_timeout 30s
<span class="k">POLICY</span> network dns:api.example.com:443

<span class="k">HEALTHCHECK</span> --interval=60s CMD /mnt/tools/file_read --agentrc-schema</code></pre></div>
      <div class="term-status" id="tstatus"><span>arc lint: ok</span><span>compiles to OCI + <code>ai.agentrc.*</code> labels</span><span>policy reviewable</span></div>
    </div>
  </div>
</section>

<section class="grid">
  <div class="card" style="display:flex;flex-direction:column;gap:.7rem;min-height:150px">
    <span style="display:grid;place-items:center;width:38px;height:38px;border-radius:9px;border:1px solid var(--line);background:var(--panel-2);color:var(--accent)"><svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M4 5h16v14H4zM8 10l2.5 2L8 14m5 0h3"/></svg></span>
    <h3 style="font-size:1rem;margin:0">Declarative &amp; reproducible</h3>
    <p class="muted" style="margin:0;font-size:.84rem;flex:1">One Agentfile captures identity, capability, policy, tools, and resources — reusing standard Dockerfile keywords plus four agent-native ones: <code>IDENTITY</code>, <code>CAPABILITY</code>, <code>SOP</code>, <code>POLICY</code>.</p>
  </div>
  <div class="card" style="display:flex;flex-direction:column;gap:.7rem;min-height:150px">
    <span style="display:grid;place-items:center;width:38px;height:38px;border-radius:9px;border:1px solid var(--line);background:var(--panel-2);color:var(--accent)"><svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M12 3l7 3v5c0 4.3-2.9 7.6-7 9-4.1-1.4-7-4.7-7-9V6z"/></svg></span>
    <h3 style="font-size:1rem;margin:0">Policy, not hope</h3>
    <p class="muted" style="margin:0;font-size:.84rem;flex:1">A <code>POLICY</code> line requests a model, resource, or constraint. The platform <strong>grants, narrows, or rejects</strong> it and enforces the decision with Cedar, deny-by-default.</p>
  </div>
  <div class="card" style="display:flex;flex-direction:column;gap:.7rem;min-height:150px">
    <span style="display:grid;place-items:center;width:38px;height:38px;border-radius:9px;border:1px solid var(--line);background:var(--panel-2);color:var(--accent)"><svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M12 3l8 4.5v9L12 21l-8-4.5v-9zM12 12l8-4.5M12 12v9M12 12L4 7.5"/></svg></span>
    <h3 style="font-size:1rem;margin:0">Portable everywhere</h3>
    <p class="muted" style="margin:0;font-size:.84rem;flex:1">The build translates intent into namespaced <code>ai.agentrc.*</code> OCI labels. Platforms read the labels — never the Agentfile — so agents ship, sign, and mirror like any container image.</p>
  </div>
</section>

<section class="get-started">
  <h2 id="install" class="flush">Install the CLI</h2>
  <p class="lead">One binary — <span style="color:var(--text-strong)">agentrc</span> (alias <span style="color:var(--text-strong)">arc</span>). It scaffolds, validates, and builds Agentfiles, inspects what an agent requests, and translates an artifact into a backend's deploy config.</p>
  <div class="install-oneliner">
    <span class="install-label">macOS &amp; Linux</span>
    <pre><code class="cmd">curl -fsSL https://agentrc.ai/install.sh | sh</code><button class="icon-button" type="button" data-copy="curl -fsSL https://agentrc.ai/install.sh | sh" aria-label="Copy install command"><svg viewBox="0 0 24 24" width="15" height="15" aria-hidden="true"><path fill="currentColor" d="M16 1H4a2 2 0 0 0-2 2v12h2V3h12V1Zm3 4H8a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h11a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2Zm0 16H8V7h11v14Z"/></svg></button></pre>
  </div>
  <div class="grid">
    <div class="card">
      <span class="install-label">Homebrew</span>
      <pre><code class="cmd">brew install
adeelahmad/tap/agentrc</code></pre>
    </div>
    <div class="card">
      <span class="install-label">Go 1.25+</span>
      <pre><code class="cmd">go install
github.com/adeelahmad/agentrc/cmd/agentrc@latest</code></pre>
    </div>
    <div class="card">
      <span class="install-label">From source</span>
      <pre><code class="cmd">git clone https://github.com/adeelahmad/agentrc
cd agentrc &amp;&amp; go build -o arc ./cmd/agentrc</code></pre>
    </div>
  </div>
  <p class="muted" style="font-size:.84rem">Prebuilt, checksum-verified binaries for macOS &amp; Linux (amd64 / arm64). Confirm with <code>arc version</code>. Prefer to read first? <code>curl -fsSL https://agentrc.ai/install.sh</code> and inspect it.</p>
</section>

<section class="get-started">
  <h2 id="build">Build and run — locally</h2>
  <p class="lead">Scaffold, validate, and compile an agent into a portable OCI artifact, then preview exactly what a local runner would execute.</p>
  <ol class="steps">
    <li><div><span class="step-k">Scaffold</span><pre><code>arc init            <span style="color:var(--muted)">› writes ./Agentfile</span></code></pre></div></li>
    <li><div><span class="step-k">Validate</span><pre><code>arc lint Agentfile  <span style="color:var(--muted)">› identity, policy &amp; schema</span></code></pre></div></li>
    <li><div><span class="step-k">Build</span><pre><code>arc build -t ghcr.io/you/hello:0.1 .  <span style="color:var(--muted)">› OCI artifact</span></code></pre></div></li>
    <li><div><span class="step-k">Preview the run</span><pre><code>arc run ghcr.io/you/hello:0.1 --backend local --dry-run</code></pre></div></li>
  </ol>
  <p class="muted" style="font-size:.84rem"><code>arc build</code> produces a real OCI image (via <code>docker build</code> and the agentrc BuildKit frontend). <code>--dry-run</code> prints the config a runner would apply — agentrc declares and translates; it ships no runtime of its own.</p>
</section>

<section class="get-started">
  <h2 id="cloud">Ship the same artifact to the cloud</h2>
  <p class="lead">The build writes <code>ai.agentrc.*</code> labels once. Point <code>arc run</code> at any backend to translate those labels into that platform's deploy form.</p>
  <div class="grid">
    <div class="card">
      <h3 style="font-size:1rem">Push once</h3>
      <pre><code class="cmd">arc push
ghcr.io/you/hello:0.1</code></pre>
      <p class="muted" style="margin:.55rem 0 0;font-size:.8rem">→ any OCI registry</p>
    </div>
    <div class="card">
      <h3 style="font-size:1rem">AWS Bedrock</h3>
      <pre><code class="cmd">arc run …hello:0.1
--backend bedrock --dry-run</code></pre>
      <p class="muted" style="margin:.55rem 0 0;font-size:.8rem">→ CreateAgentRuntime JSON</p>
    </div>
    <div class="card">
      <h3 style="font-size:1rem">Kubernetes</h3>
      <pre><code class="cmd">arc run …hello:0.1
--backend kubernetes --dry-run</code></pre>
      <p class="muted" style="margin:.55rem 0 0;font-size:.8rem">→ deploy manifests</p>
    </div>
  </div>
  <div class="callout"><strong>Same artifact, same labels, three substrates.</strong> Reference translators — a proof of concept until platforms read <code>ai.agentrc.*</code> labels natively. Not production runners.</div>
</section>

<h2>The separation agentrc creates</h2>
<table>
  <thead>
    <tr><th>Concern</th><th>Defined by</th><th>Read / enforced by</th><th>Why it matters</th></tr>
  </thead>
  <tbody>
    <tr><td>Agent identity, capabilities, objective</td><td><code>IDENTITY / CAPABILITY / SOP</code></td><td>agent author</td><td>Clear purpose and scope</td></tr>
    <tr><td>Tools, skills, MCP servers</td><td><code>COPY / ADD --remote</code> into <code>/mnt</code></td><td>compiler → layers + labels</td><td>Portable across stacks</td></tr>
    <tr><td>Resource, model, network, lifecycle requests</td><td><code>POLICY</code> (typed namespaces)</td><td>platform (grant / narrow / reject)</td><td>Governed and reviewable</td></tr>
    <tr><td>Enforcement</td><td>typed requests compiled to Cedar</td><td>platform (deny-by-default, <code>forbid &gt; permit</code>)</td><td>Least privilege by design</td></tr>
    <tr><td>Packaging and sharing</td><td>OCI artifact + <code>ai.agentrc.*</code> labels</td><td>any OCI registry</td><td>Interoperable distribution</td></tr>
    <tr><td>Execution substrate</td><td>run-time choice (<code>arc run --backend</code>)</td><td>local, container, microVM, cloud runners</td><td>Freedom with guardrails</td></tr>
  </tbody>
</table>

<div class="callout"><strong>Core slogan:</strong> The Agentfile declares one agent. The lockfile pins dependencies. The package makes it portable. The policy makes boundaries reviewable. The registry makes it shareable. Compatible runners execute it.</div>

<h2 id="draft">Current draft</h2>
<div class="pill-row">
  <span class="pill">Working Draft 0.1.0-draft.6</span>
  <span class="pill">Four keywords</span>
  <span class="pill">/mnt projection</span>
  <span class="pill">OCI labels</span>
  <span class="pill">Cedar enforcement</span>
  <span class="pill">Secrets deferred</span>
</div>
<p class="muted" style="font-size:.84rem;margin-top:1.1rem">The project is published as a standards-style repository: specification first, reference tooling second.</p>

<script>
/* Hero headline typewriter — animates on load (runs even with reduced-motion, by request). */
(function () {
  var el = document.getElementById('rotator');
  if (!el) return;
  var phrases = [
    'portable, governed ai agents',
    'a federation of trust, developer to agent',
    'secure boundaries for agentic apps'
  ];
  var pi = 0, ci = 0, deleting = false;
  el.textContent = '';
  function tick() {
    var p = phrases[pi];
    if (!deleting) {
      ci++; el.textContent = p.slice(0, ci);
      if (ci === p.length) { deleting = true; return setTimeout(tick, 1900); }
      return setTimeout(tick, 55);
    }
    ci--; el.textContent = p.slice(0, ci);
    if (ci === 0) { deleting = false; pi = (pi + 1) % phrases.length; return setTimeout(tick, 420); }
    return setTimeout(tick, 28);
  }
  setTimeout(tick, 350);
})();

/* Hero terminal — cycles: cat Agentfile → compiled OCI labels → Bedrock deploy config. */
(function () {
  var cmd = document.getElementById('tcmd'),
      body = document.getElementById('tbody'),
      status = document.getElementById('tstatus');
  if (!cmd || !body || !status) return;
  var frames = [
    {
      cmd: 'cat Agentfile',
      body: body.innerHTML, // the Agentfile as authored (frame 0, also the no-JS view)
      status: '<span>arc lint: ok</span><span>Dockerfile-shaped</span><span>policy reviewable</span>'
    },
    {
      cmd: 'arc build -t ghcr.io/you/hello:0.1 .',
      body: '<pre><code><span class="c"># the build translates intent into OCI labels</span>\n' +
            '<span class="k">ai.agentrc.identity.name</span>=hello\n' +
            '<span class="k">ai.agentrc.capability</span>=text\n' +
            '<span class="k">ai.agentrc.model.name</span>=claude-sonnet-4\n' +
            '<span class="k">ai.agentrc.agent.tool_timeout</span>=30s\n' +
            '<span class="k">ai.agentrc.network.dns.api.example.com</span>=443\n' +
            '<span class="k">ai.agentrc.tool.file_read</span>=/mnt/tools/file_read\n' +
            '<span class="k">ai.agentrc.sop.sha256</span>=<span class="s">a1b2c3…</span></code></pre>',
      status: '<span>compiled → OCI</span><span>7 <code>ai.agentrc.*</code> labels</span><span>signed &amp; portable</span>'
    },
    {
      cmd: 'arc run …hello:0.1 --backend bedrock --dry-run',
      body: '<pre><code><span class="c"># same labels → a backend deploy config</span>\n' +
            '<span class="p">{</span>\n' +
            '  <span class="k">"agentRuntimeName"</span>: <span class="s">"hello"</span>,\n' +
            '  <span class="k">"containerUri"</span>: <span class="s">"ghcr.io/you/hello:0.1"</span>,\n' +
            '  <span class="k">"roleArn"</span>: <span class="s">"arn:aws:iam::…:role/agent"</span>,\n' +
            '  <span class="k">"networkMode"</span>: <span class="s">"PUBLIC"</span>,\n' +
            '  <span class="k">"serverProtocol"</span>: <span class="s">"HTTP"</span>,\n' +
            '  <span class="k">"environmentVariables"</span>: <span class="p">{</span> <span class="k">"LOG_LEVEL"</span>: <span class="s">"info"</span> <span class="p">}</span>\n' +
            '<span class="p">}</span></code></pre>',
      status: '<span>translated → Bedrock</span><span>CreateAgentRuntime</span><span>same artifact, any substrate</span>'
    }
  ];
  var i = 0;
  body.style.transition = 'opacity .4s ease';
  function show(next) {
    body.style.opacity = '0';
    setTimeout(function () {
      cmd.textContent = frames[next].cmd;
      body.innerHTML = frames[next].body;
      status.innerHTML = frames[next].status;
      body.style.opacity = '1';
    }, 400);
  }
  setInterval(function () { i = (i + 1) % frames.length; show(i); }, 4200);
})();
</script>
