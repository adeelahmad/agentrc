---
layout: doc
title: Implementation mapping
description: "How the current AIO work maps into agentrc."
permalink: /docs/implementation-mapping/
---
# Implementation mapping

The uploaded current work contains useful implementation material: Agentfile parsing, runner experimentation, Cedar policy gates, secret broker work, OCI-related pieces, and projection ideas.

In the agentrc framing, these should be renamed and positioned carefully:

| Current implementation idea | agentrc framing |
|---|---|
| AIO runtime | reference tooling or runner experiment |
| SubstrateDriver | runner adapter interface |
| isolation flag | runner-selected execution substrate |
| tools patched into VM | optional tool projection profile |
| Cedar gate | security profile enforcement example |
| Vault broker | secret resolution profile example |
| OCI image/package work | package profile implementation |
| microVM support | one compatible runner backend, not product identity |

The implementation can still exist. It should not define the public identity of the project.
