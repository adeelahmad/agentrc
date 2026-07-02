---
type: tasks
story: S2-06
---
# S2-06 Tasks — Demo (T26)

One agent, three backends. Narrative VERBATIM. Runs in Wave 5 (after T20 + S2-04 + S2-05).
Uses the real `Agentfile.code-reviewer` (with the T20 commented `substrate.aws.*` +
`agent.auth.*` block). Version = `draft.6` by now.

## T26 — One agent, three backends [P2]

Author the demo section/doc (grep-locate the natural home — e.g. `docs/` or `examples/` demo
prose) with these commands:
- `arc build -t ghcr.io/agentrc/code-reviewer:1.0 .`
- `arc run ghcr.io/agentrc/code-reviewer:1.0 --backend local --isolation microvm`
- `arc run ghcr.io/agentrc/code-reviewer:1.0 --backend bedrock --dry-run`
- `arc run ghcr.io/agentrc/code-reviewer:1.0 --backend kubernetes --dry-run`

Narrative VERBATIM (no paraphrase): "Same artifact, same labels, three substrates. The
translators are the proof of concept; the labels are the standard."

Constraints: no removed `--substrate` flag; no dropped backend (gcp/compose) in the commands;
both dry-run outputs must satisfy §V.9 (bedrock JSON, kubernetes YAML). §0.8 line already lives
on the backend surfaces (S2-04/S2-05).

Verify: the four commands appear exactly as above; the narrative string is verbatim; the two
dry-runs parse (§V.9).
</content>
