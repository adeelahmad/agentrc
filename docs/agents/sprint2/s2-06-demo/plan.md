---
type: plan
story: S2-06
scope: "tests only"
---
# S2-06 Test Plan — Demo (RED)

Tests only. Each check is a function in `scripts/verify-sprint2.sh`. Every bullet FAILS now
because the demo section does not yet exist.

## T26 — one agent, three backends

- [ ] `scripts/verify-sprint2.sh::t26_demo_narrative_verbatim` — assert the verbatim string
      "Same artifact, same labels, three substrates. The translators are the proof of concept;
      the labels are the standard." is present exactly once in the demo doc. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint2.sh::t26_three_backend_commands` — assert all four commands present
      (`arc build -t ghcr.io/agentrc/code-reviewer:1.0`, `--backend local --isolation microvm`,
      `--backend bedrock --dry-run`, `--backend kubernetes --dry-run`). RED now: FAILS.
- [ ] `scripts/verify-sprint2.sh::t26_no_dropped_backends` — assert the demo uses no
      `--substrate` and no `--backend gcp`/`--backend compose`. RED now: FAILS (doc absent).

## Suite guard

- [ ] `scripts/verify-sprint2.sh::v9_backend_dryruns` — assert the demo's bedrock JSON and
      kubernetes YAML dry-runs parse (§V.9). RED now: FAILS — translators/demo not ready; shared
      guard with S2-04.
</content>
