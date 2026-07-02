// Package agentfile extracts the agentrc-specific surface (IDENTITY,
// CAPABILITY, SOP, POLICY, and the ADD --remote extension) from a
// Dockerfile-shaped Agentfile, per grammar/Agentfile.ebnf (0.1.0-draft.6).
// Everything else is a standard Dockerfile instruction, left untouched for
// BuildKit's own dockerfile2llb compiler (see internal/llb).
package agentfile

import (
	"path"
	"strings"
)

// Identity is the IDENTITY keyword's accumulated key=value pairs (repeatable
// across lines; later lines add or overwrite keys).
type Identity map[string]string

// SOP is the agent's system prompt, from whichever of the three forms
// (inline, heredoc, or file-backed via COPY/ADD to /mnt/SOP) was used.
type SOP struct {
	// Content is set for the inline/heredoc keyword forms; empty for the
	// file-backed form (where the content lives in a build-context file
	// instead and is hashed separately — see internal/llb).
	Content string
	// FileBacked is true when SOP was authored as `COPY .../ADD --remote ...
	// /mnt/SOP` rather than the SOP keyword.
	FileBacked bool
	Line       int
}

// Policy is one `POLICY <namespaced.key> <value>` request.
type Policy struct {
	Key   string
	Value string
	Line  int
}

// RemoteAdd is one `ADD --remote ...` extension line.
type RemoteAdd struct {
	Source string
	Dest   string
	// Cached is true for --cached (or the default when neither is given):
	// fetch at build time and embed as a layer.
	Cached bool
	// Runtime is true for --runtime: record a reference only, fetched at
	// agent bootstrap. Mutually exclusive with Cached.
	Runtime bool
	// FailIfUnavailable is true (the default) if an unresolved resource
	// must abort the build; false (--warn-if-unavailable) degrades to a
	// runtime-style reference with a warning instead.
	FailIfUnavailable bool
	Chmod             string
	Chown             string
	Line              int

	// ResolvedDigest is populated by internal/llb after an actual build-time
	// fetch of a --cached resource (sha256 of the fetched content). Empty
	// when unresolved — e.g. an unsupported source scheme degraded to a
	// runtime-style reference under --warn-if-unavailable.
	ResolvedDigest string
}

// ResourceKind classifies a /mnt destination path.
type ResourceKind string

const (
	KindTool  ResourceKind = "tool"
	KindSkill ResourceKind = "skill"
	KindMCP   ResourceKind = "mcp"
	KindSOP   ResourceKind = "sop" // /mnt/SOP itself, not a tool/skill/mcp
	KindOther ResourceKind = ""    // under /mnt but not a recognized subtree
)

// LocalResource is a COPY or plain ADD (no --remote) instruction whose
// destination lands under /mnt — a locally embedded tool/skill/mcp/SOP file.
type LocalResource struct {
	Source string
	Dest   string
	Kind   ResourceKind
	Name   string
	Line   int
}

// File is the parsed agentrc surface of one Agentfile.
type File struct {
	Source string // original, unmodified Agentfile text (for agentfile_sha256)

	Identity     Identity
	Capabilities []string
	SOP          *SOP
	Policies     []Policy
	RemoteAdds   []RemoteAdd

	// LocalResources is populated by internal/llb after parsing
	// CleanedSource with the real Dockerfile instruction parser — it
	// requires quote-aware COPY/ADD argument parsing that only the real
	// parser does correctly, so it isn't filled in by Extract itself.
	LocalResources []LocalResource

	// CleanedSource is Source with every agentrc-specific line (IDENTITY,
	// CAPABILITY, POLICY, SOP in any form, ADD --remote) blanked out but
	// still present as an empty line, so remaining standard Dockerfile
	// instructions keep their original line numbers. It is valid input to
	// BuildKit's own dockerfile parser/compiler.
	CleanedSource []byte
}

// CleanMntPath cleans dest (resolving "." and ".." segments) and reports
// whether the result is actually contained under /mnt. A raw
// strings.HasPrefix(dest, "/mnt/") check is not sufficient: a destination
// like "/mnt/tools/../../etc/passwd" satisfies that prefix check textually
// but resolves outside /mnt entirely (path.Clean("/mnt/tools/../../etc/passwd")
// == "/etc/passwd"). Callers must validate and use the cleaned path, never
// the raw one, anywhere a destination is written to or classified from.
func CleanMntPath(dest string) (cleaned string, ok bool) {
	cleaned = path.Clean(dest)
	if cleaned == "/mnt" {
		return cleaned, true
	}
	return cleaned, strings.HasPrefix(cleaned, "/mnt/")
}

// ResourceKindForDest classifies a /mnt destination path per
// docs/agentfile.md's "/mnt projection layout". Destinations that escape
// /mnt after cleaning (see CleanMntPath) are never classified as a resource.
func ResourceKindForDest(dest string) ResourceKind {
	cleaned, ok := CleanMntPath(dest)
	if !ok {
		return KindOther
	}
	switch {
	case cleaned == "/mnt/SOP":
		return KindSOP
	case hasDirPrefix(cleaned, "/mnt/tools/"):
		return KindTool
	case hasDirPrefix(cleaned, "/mnt/skills/"):
		return KindSkill
	case hasDirPrefix(cleaned, "/mnt/mcp/"):
		return KindMCP
	default:
		return KindOther
	}
}

func hasDirPrefix(s, prefix string) bool {
	return len(s) > len(prefix) && s[:len(prefix)] == prefix
}
