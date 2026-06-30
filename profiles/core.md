---
layout: doc
title: Core
description: "Core profile."
permalink: /profiles/core/
---
# Core

**Status:** Working Draft  
**Version:** 0.1.0-draft.2

## Purpose

The Core Agentfile Profile defines the minimum behavior required to parse, validate, and inspect a single-agent Agentfile.

It does not define execution.

## Required behavior

An implementation claiming this profile MUST:

1. read a text file named `Agentfile` or an explicitly supplied path;
2. ignore blank lines and comments;
3. parse known directives case-sensitively;
4. preserve directive order for lockfile generation and review;
5. capture `POLICY ... END` blocks without interpreting inner lines as Agentfile directives;
6. reject unknown directives unless an extension profile is enabled;
7. produce a structured parse tree suitable for linting and package generation;
8. report line numbers for parse failures.

## Core directive set

The v0.1 core directive set is:

```text
AGENT FROM SHELL CMD TOOL TOOLSET FUNCTION SKILL SERVER MCP URL CRED
BIND MOUNT PLUGIN POLICY ALLOW DENY RATELIMIT AUDIT TIMEOUT LIMIT
SLICE IMAGE ISOLATION BROKER BACKEND TRACE MEMORY OPTIMIZER HEALTHCHECK
```

## Compatibility notes

The current AgentRC implementation already recognizes this directive family, though not every directive has full semantic enforcement. This profile defines target semantics; implementations should publish a support matrix rather than implying all recognized directives are enforced.

## Validation recommendations

A linter SHOULD warn when:

1. `AGENT` is missing in a package intended for publication;
2. `FROM` uses a mutable tag such as `latest`;
3. a `CRED` contains plaintext-like material;
4. `BIND` lacks a mode;
5. a `TOOL` is declared but no policy permits it under deny-by-default mode;
6. `AUDIT` is absent;
7. `ISOLATION`, `IMAGE`, `SLICE`, `PLUGIN`, or `BACKEND` couples the package to one runner unnecessarily.

