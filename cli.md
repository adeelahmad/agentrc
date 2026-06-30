---
layout: doc
title: CLI
description: "The agentrc command-line interface — coming soon."
permalink: /cli/
---
# agentrc CLI

<div class="pill-row">
  <span class="pill">Coming soon</span>
</div>

A first-class command-line interface for agentrc is in progress. It will let you create, validate, lock, build, sign, inspect, and publish agents straight from the terminal.

<div class="callout">
<strong>Coming soon.</strong> The <code>agentrc</code> CLI is not released yet. The specification, schemas, and examples on this site are the source of truth today — the CLI will implement them.
</div>

## What it will do

| Command | Purpose |
|---|---|
| `agentrc init` | Scaffold a new `Agentfile`. |
| `agentrc lint` | Validate an Agentfile against the Core profile and warn on risky declarations. |
| `agentrc lock` | Resolve dependencies and write `agentrc.lock`. |
| `agentrc build` | Produce an OCI-compatible agent package. |
| `agentrc inspect` | Show declared tools, mounts, hosts, secrets, SOPs, and policy before running. |
| `agentrc sign` / `verify` | Sign and verify packages and provenance. |
| `agentrc push` / `pull` | Publish to and fetch from an OCI registry. |

## In the meantime

- Read [What is agentrc?](/docs/what-is-agentrc/) for the why.
- Start with the [Quickstart](/docs/quickstart/) and the [Specification](/spec/).
- Track or contribute on [GitHub](https://github.com/adeelahmad/agentrc), or say hi on [Discord](https://discord.gg/jWx6Qak5D).
