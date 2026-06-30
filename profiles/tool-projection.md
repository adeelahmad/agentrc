---
layout: doc
title: Tool Projection
description: "Tool Projection profile."
permalink: /profiles/tool-projection/
---
# Tool Projection

**Status:** Working Draft  
**Version:** 0.1.0-draft.2

## Purpose

This optional profile defines a command-oriented way for agents and humans to discover and invoke declared capabilities.

This profile does not make AgentRC a runtime. It defines a portable surface a compatible runner may expose.

## Projection root

A runner implementing this profile SHOULD expose a projection root such as:

```text
/agentrc
```

A local implementation MAY expose the same model through a FUSE mount, generated directory, or CLI commands.

## Recommended layout

```text
/agentrc/
  tools/
    file_read
    file_read.schema.json
    file_read.help
    file_read.events
  functions/
    summarize_meeting
    summarize_meeting.schema.json
    summarize_meeting.help
    summarize_meeting.events
  skills/
    pr-review/
      SKILL.md
  mcp/
    github/
      get_issue
      get_issue.schema.json
  memory/
    writing.json
  proc/
    audit
    policy
    limits
    status
    tools
    functions
```

## Invocation semantics

A projected tool SHOULD be executable.

The invocation contract SHOULD support:

1. argv-style arguments;
2. stdin input for structured payloads;
3. stdout structured output;
4. stderr diagnostics;
5. exit code `0` for success;
6. non-zero exit for failure;
7. a distinct denied/policy exit code when possible.

## CLI equivalence

A runner MAY provide equivalent CLI commands instead of filesystem projection:

| Projection | CLI equivalent |
|---|---|
| `/agentrc/tools/file_read --path x` | `agentrc run tool utcp:file_read --path x` |
| `/agentrc/functions/foo --json args.json` | `agentrc run function foo --json args.json` |
| `cat /agentrc/proc/audit` | `agentrc dmesg` |
| `cat /agentrc/proc/policy` | `agentrc policy show` |

## Security

Every projected invocation remains subject to the active policy profile.

