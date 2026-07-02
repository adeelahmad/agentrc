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
