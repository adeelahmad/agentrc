package main

import "encoding/json"

// localPositioning is the §0.8 positioning line, kept verbatim on a production
// backend surface so the reference-translator framing travels with the code:
//
// Reference translators — a proof of concept until platforms read `ai.agentrc.*`
// labels natively. Not production runners.
const localPositioning = "Reference translators — a proof of concept until platforms read `ai.agentrc.*` labels natively. Not production runners."

// translateLocal renders the microsandbox exec plan for `--backend local` (the
// default). This is the plumbing seam (open question #5, UNRESOLVED): dry-run /
// plumbing depth is acceptable — it wires the VMM MVP under the translate seam
// rather than shipping a production runner.
func translateLocal(labels map[string]string) (string, error) {
	name := labels["ai.agentrc.identity.name"]
	if name == "" {
		name = "agent"
	}
	plan := map[string]any{
		"backend":   "local",
		"runtime":   "microsandbox",
		"isolation": "microvm",
		"name":      name,
		"image":     labels["image.ref"],
	}
	b, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
