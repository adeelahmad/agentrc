---
layout: doc
title: Changelog
description: "AgentRC site and spec changelog."
permalink: /changelog/
---
# Changelog

## 0.1.0-draft.3 — 2026-06-30

Rewritten after reviewing the current AgentRC source archive.

### Changed

- Reframed AgentRC as the **Agentfile and agent-package specification**, not an agent isolation orchestrator.
- Moved runtime/driver concerns into a **Runner Conformance Profile** instead of the core specification.
- Kept current implementation directives, but classified them into identity, capability, boundary, governance, runner-hint, and experimental groups.
- Replaced the unsupported `SPEC` directive with parser-compatible `# syntax=agentrc.agentfile/v0.1`.
- Standardized package distribution around OCI-compatible artifacts/images.
- Added explicit lockfile and package schemas.
- Added current implementation mapping and gap list.
- Treated tools-as-files as an optional **Tool Projection Profile**, not the whole product.
- Added a non-normative workflow companion draft for ASL-like multi-agent orchestration.

### Corrected

- Avoided claiming AgentRC provides isolation directly.
- Avoided claiming AgentRC is a runtime or runner.
- Avoided hard-binding AgentRC to microsandbox, Firecracker, Docker, Strands, or any cloud.
- Made `ISOLATION`, `IMAGE`, `SLICE`, `PLUGIN`, and `BACKEND` runner hints/profile directives rather than core product identity.

### Known open issues

- Real Cedar evaluator is required before Cedar profile conformance can be claimed.
- Current `CRED` parsing has inconsistent named-reference behavior.
- Current `BIND` parsing is inconsistent across implementation paths.
- Multi-agent workflow policy composition is intentionally left open.

## 0.1.0-draft.1 — 2026-06-30

Initial GitHub-ready working draft.


## 0.1.0-draft.3 — AgentRC rebrand

- Renamed public identity from AIO to AgentRC — Agent Run Config.
- Reframed the project as a publishable Jekyll documentation site.
- Added GitHub Pages configuration, CNAME, documentation pages, layouts, CSS, and publishing workflow.
- Preserved implementation work as reference tooling, not runtime identity.

## v0.1-site.2 — Hybrid brand and theme support

- Reworked the visual system into the selected hybrid direction: standards-style documentation with a subtle holographic AgentRC mark.
- Added light and dark theme support with a persistent user toggle.
- Added system-theme detection and early theme application to reduce flash during page load.
- Replaced the placeholder `rc` square with an AgentRC `R` glyph and matching SVG favicon.
- Updated CSS tokens so future brand assets can be adjusted centrally.
