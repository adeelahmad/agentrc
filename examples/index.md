---
layout: doc
title: Examples
description: "agentrc example Agentfiles and workflow drafts."
permalink: /examples/
---
# Examples

These examples demonstrate the intended style of agentrc declarations.

## Files

- [Minimal Agentfile](/examples/Agentfile.minimal)
- [Secure workspace Agentfile](/examples/Agentfile.secure-workspace)
- [Code reviewer Agentfile](/examples/Agentfile.code-reviewer)
- [Vault agent Agentfile](/examples/Agentfile.vault-agent)
- [Workflow draft YAML](/examples/agent-workflow.yaml)

## Minimal

```dockerfile
# syntax=agentrc.agentfile/v0.1

AGENT hello-local
CMD python ./agent.py
TOOL utcp:file_read
AUDIT basic
POLICY
  permit(
    principal == AgentRC::Agent::"hello-local",
    action == AgentRC::Action::"tool.invoke",
    resource == AgentRC::Tool::"file_read"
  );
END
```

More example files are included in this directory as raw examples.
