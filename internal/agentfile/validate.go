package agentfile

import (
	"fmt"
	"strings"
)

// Issue is a single lint finding.
type Issue struct {
	Line     int
	Severity string // "error" | "warning"
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("line %d: %s: %s", i.Line, i.Severity, i.Message)
}

func errIssue(line int, format string, args ...any) Issue {
	return Issue{Line: line, Severity: "error", Message: fmt.Sprintf(format, args...)}
}

func warnIssue(line int, format string, args ...any) Issue {
	return Issue{Line: line, Severity: "warning", Message: fmt.Sprintf(format, args...)}
}

var knownPolicyNamespaces = map[string]bool{"agent": true, "substrate": true, "model": true}

// Validate lints a *File extracted by Extract. It only checks what Extract
// itself can see — SOP conflicts with a file-backed COPY/ADD target (which
// requires the real Dockerfile instruction parse) are caught later, at
// translate time, in internal/llb.
func Validate(f *File) []Issue {
	var issues []Issue

	if f.Identity["name"] == "" {
		issues = append(issues, errIssue(0, "IDENTITY must declare a name (IDENTITY name=<value>)"))
	}

	if len(f.Capabilities) == 0 {
		issues = append(issues, warnIssue(0, "no CAPABILITY declared"))
	}

	if f.SOP == nil {
		issues = append(issues, warnIssue(0, "no SOP declared"))
	}

	for _, ra := range f.RemoteAdds {
		// CleanMntPath, not a raw prefix check: a destination like
		// "/mnt/tools/../../etc/passwd" satisfies strings.HasPrefix(dest,
		// "/mnt/") textually but resolves outside /mnt entirely.
		if _, ok := CleanMntPath(ra.Dest); !ok {
			issues = append(issues, errIssue(ra.Line, "ADD --remote destination must be under /mnt, got %q", ra.Dest))
		}
	}

	for _, p := range f.Policies {
		issues = append(issues, validatePolicy(p)...)
	}

	issues = append(issues, validateInstructions(f.CleanedSource)...)

	return issues
}

// retiredKeywords are directives from pre-draft.5 agentrc (and never-valid
// keywords from other tools) that must not appear as leading tokens. SHELL is
// intentionally excluded — it is a valid Dockerfile instruction.
var retiredKeywords = map[string]bool{
	"AGENT": true, "TOOL": true, "TOOLSET": true, "FUNCTION": true, "SKILL": true,
	"SERVER": true, "MCP": true, "URL": true, "CRED": true, "BIND": true,
	"MOUNT": true, "PLUGIN": true, "ALLOW": true, "DENY": true, "RATELIMIT": true,
	"TIMEOUT": true, "LIMIT": true, "SLICE": true, "IMAGE": true, "ISOLATION": true,
	"BROKER": true, "BACKEND": true, "TRACE": true, "MEMORY": true, "OPTIMIZER": true,
	"AUDIT": true, "SECRET": true,
	// Never-were agentrc keywords that show up in YAML-era / hand-written mistakes.
	"SYSTEM": true, "RUNTIME": true, "MODEL": true, "NETWORK": true,
	"CAPABILITIES": true, "PROMPT": true, "PACKAGING": true,
}

// validateInstructions scans the cleaned Dockerfile source (agent keywords
// already blanked, line numbers preserved) for a required FROM and for retired
// or never-valid leading directives that would otherwise fail silently at lint
// and only blow up at build time.
func validateInstructions(cleaned []byte) []Issue {
	var issues []Issue
	hasFROM := false
	cont := false
	for i, raw := range strings.Split(string(cleaned), "\n") {
		endsBackslash := strings.HasSuffix(strings.TrimRight(raw, " \t"), "\\")
		if cont { // continuation of a previous instruction (e.g. a multi-line RUN)
			cont = endsBackslash
			continue
		}
		cont = endsBackslash
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tok := strings.ToUpper(strings.Fields(line)[0])
		if tok == "FROM" {
			hasFROM = true
			continue
		}
		if retiredKeywords[tok] {
			issues = append(issues, errIssue(i+1,
				"%q is a retired/legacy directive, not a valid instruction — the four agentrc keywords are IDENTITY, CAPABILITY, SOP, POLICY; tools, skills, and MCP servers are COPY/ADD --remote into /mnt", tok))
		}
	}
	if !hasFROM {
		issues = append(issues, errIssue(0, "Agentfile must declare a FROM base image (spec §2)"))
	}
	return issues
}

func validatePolicy(p Policy) []Issue {
	if p.Key == "network" {
		if _, _, err := parseNetworkPolicyValue(p.Value); err != nil {
			return []Issue{errIssue(p.Line, "POLICY network: %v", err)}
		}
		return nil
	}

	ns, _, ok := strings.Cut(p.Key, ".")
	if !ok || !knownPolicyNamespaces[ns] {
		return []Issue{warnIssue(p.Line, "POLICY %s: %q is not one of the documented namespaces (agent, substrate, model, network) — allowed as an extension, but check for a typo", p.Key, ns)}
	}
	return nil
}

// HasErrors reports whether any issue has severity "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}
