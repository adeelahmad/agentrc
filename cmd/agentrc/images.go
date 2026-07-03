package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

// localImageConfig reads a local Docker image's labels and env through the
// user's own `docker` (which knows the active context, e.g. colima). ok is
// false when the image is not present locally or docker is unavailable.
func localImageConfig(ref string) (labels map[string]string, env []string, ok bool) {
	out, err := exec.Command("docker", "image", "inspect", ref, "--format", "{{json .Config}}").Output()
	if err != nil {
		return nil, nil, false
	}
	var cfg struct {
		Labels map[string]string `json:"Labels"`
		Env    []string          `json:"Env"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &cfg); err != nil {
		return nil, nil, false
	}
	return cfg.Labels, cfg.Env, true
}

// agentLabelSet assembles the label map the backend translators consume: the
// image's org.agentrc.* labels, the image reference (image.ref), and env.<NAME>
// entries derived from the image's Env.
func agentLabelSet(ref string, labels map[string]string, env []string) map[string]string {
	out := make(map[string]string, len(labels)+len(env)+1)
	for k, v := range labels {
		out[k] = v
	}
	out["image.ref"] = ref
	for _, e := range env {
		if i := strings.IndexByte(e, '='); i > 0 {
			out["env."+e[:i]] = e[i+1:]
		}
	}
	return out
}
