---
layout: doc
title: Documentation
description: "agentrc documentation index: the Dockerfile-shaped Agentfile, its OCI labels, and the platform that reads them."
permalink: /docs/
---
# Documentation

agentrc is organized as a small set of stable documents rather than a monolithic application manual.

<div class="grid two">
  <div class="card"><h3><a href="{{ '/docs/what-is-agentrc/' | relative_url }}">What is agentrc?</a></h3><p>Start here: a Dockerfile-shaped recipe for portable, governed agents, the problem it solves, and why it is needed.</p></div>
  <div class="card"><h3><a href="{{ '/spec/' | relative_url }}">Specification</a></h3><p>The full v1 working draft: the four new keywords, the <code>/mnt</code> layout, build &rarr; <code>org.agentrc.*</code> labels, and Cedar enforcement.</p></div>
  <div class="card"><h3><a href="{{ '/docs/quickstart/' | relative_url }}">Quickstart</a></h3><p>Write your first Agentfile, build it (BuildKit frontend or <code>arc build</code>), read the emitted labels, push, and run.</p></div>
  <div class="card"><h3><a href="{{ '/docs/agentfile/' | relative_url }}">Agentfile</a></h3><p>The single-agent declaration: <code>IDENTITY</code>, <code>CAPABILITY</code>, <code>SOP</code>, <code>POLICY</code> over standard Dockerfile keywords.</p></div>
  <div class="card"><h3><a href="{{ '/docs/security/' | relative_url }}">Security</a></h3><p><code>POLICY</code> requests and secret references become labels; the platform grants/narrows/rejects and enforces via Cedar, deny-by-default.</p></div>
  <div class="card"><h3><a href="{{ '/docs/package/' | relative_url }}">Package model</a></h3><p>A standard OCI artifact: layers carry <code>/mnt</code> resources, the config carries <code>org.agentrc.*</code> labels. Sign, mirror, review before run.</p></div>
  <div class="card"><h3><a href="{{ '/docs/runners/' | relative_url }}">Runners</a></h3><p>How a platform reads <code>org.agentrc.*</code> labels (never the Agentfile) and a substrate executes <code>CMD</code> &mdash; without agentrc becoming a runtime.</p></div>
</div>
