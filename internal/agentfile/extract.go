package agentfile

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Extract parses raw Agentfile bytes and pulls out the agentrc-specific
// surface (IDENTITY, CAPABILITY, SOP, POLICY, ADD --remote), leaving a
// CleanedSource that is a plain, valid Dockerfile for BuildKit's own parser.
//
// This runs in two phases because BuildKit's generic Dockerfile tokenizer
// (frontend/dockerfile/parser) only understands heredocs (`<<DELIM`) for
// instructions in its own dispatch table — SOP isn't one, so a raw SOP
// heredoc's body would otherwise be mis-tokenized as bogus top-level
// instructions. Phase 1 strips SOP heredocs directly from the raw bytes
// (blanking their lines, preserving line numbers) using this package's own
// scanner; phase 2 parses the result with BuildKit's real tokenizer, which
// correctly handles comments, quoting, and line continuation for
// everything else (including a multi-line `ADD --remote \` continuation).
func Extract(src []byte) (*File, error) {
	phase1, heredocSOP, err := stripSOPHeredocs(src)
	if err != nil {
		return nil, err
	}

	result, err := parser.Parse(strings.NewReader(string(phase1)))
	if err != nil {
		return nil, fmt.Errorf("parsing Dockerfile syntax: %w", err)
	}

	f := &File{
		Source:   string(src),
		Identity: Identity{},
		SOP:      heredocSOP,
	}

	lines := strings.Split(string(phase1), "\n")
	blank := func(start, end int) {
		for i := start; i <= end && i-1 < len(lines); i++ {
			lines[i-1] = ""
		}
	}

	for _, child := range result.AST.Children {
		keyword := strings.ToUpper(child.Value)
		switch keyword {
		case "IDENTITY":
			pairs, err := parseKeyValuePairs(afterKeyword(child.Original, child.Value))
			if err != nil {
				return nil, fmt.Errorf("line %d: IDENTITY: %w", child.StartLine, err)
			}
			for k, v := range pairs {
				f.Identity[k] = v
			}
			blank(child.StartLine, child.EndLine)

		case "CAPABILITY":
			v := strings.TrimSpace(afterKeyword(child.Original, child.Value))
			if v == "" {
				return nil, fmt.Errorf("line %d: CAPABILITY requires a value", child.StartLine)
			}
			f.Capabilities = append(f.Capabilities, v)
			blank(child.StartLine, child.EndLine)

		case "POLICY":
			rest := strings.TrimSpace(afterKeyword(child.Original, child.Value))
			sp := strings.IndexAny(rest, " \t")
			if sp < 0 {
				return nil, fmt.Errorf("line %d: POLICY requires <namespaced.key> <value>", child.StartLine)
			}
			key := rest[:sp]
			value := strings.TrimSpace(rest[sp:])
			if key == "" || value == "" {
				return nil, fmt.Errorf("line %d: POLICY requires <namespaced.key> <value>", child.StartLine)
			}
			f.Policies = append(f.Policies, Policy{Key: key, Value: value, Line: child.StartLine})
			blank(child.StartLine, child.EndLine)

		case "SOP":
			if f.SOP != nil {
				return nil, fmt.Errorf("line %d: SOP declared more than once (already declared at line %d)", child.StartLine, f.SOP.Line)
			}
			content := afterKeyword(child.Original, child.Value)
			if content == "" {
				return nil, fmt.Errorf("line %d: SOP requires inline text or a <<DELIM heredoc", child.StartLine)
			}
			f.SOP = &SOP{Content: content, Line: child.StartLine}
			blank(child.StartLine, child.EndLine)

		case "ADD":
			if !hasFlag(child.Flags, "--remote") {
				continue // standard local ADD; leave for the real Dockerfile parser
			}
			ra, err := parseRemoteAdd(child)
			if err != nil {
				return nil, err
			}
			f.RemoteAdds = append(f.RemoteAdds, *ra)
			blank(child.StartLine, child.EndLine)
		}
	}

	f.CleanedSource = []byte(strings.Join(lines, "\n"))
	return f, nil
}

var sopHeredocStartRe = regexp.MustCompile(`(?i)^SOP\s+<<(\S+)\s*$`)

// stripSOPHeredocs finds `SOP <<DELIM ... DELIM` blocks in raw source,
// blanks their lines (preserving line count/numbers), and returns the
// extracted SOP (nil if none found). The heredoc body is never parsed as
// instructions.
func stripSOPHeredocs(src []byte) ([]byte, *SOP, error) {
	lines := strings.Split(string(src), "\n")
	var sop *SOP

	for i := 0; i < len(lines); i++ {
		m := sopHeredocStartRe.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if m == nil {
			continue
		}
		if sop != nil {
			return nil, nil, fmt.Errorf("line %d: SOP declared more than once (already declared at line %d)", i+1, sop.Line)
		}
		delim := m[1]
		startLine := i + 1

		j := i + 1
		found := false
		var body []string
		for ; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == delim {
				found = true
				break
			}
			body = append(body, lines[j])
		}
		if !found {
			return nil, nil, fmt.Errorf("line %d: SOP heredoc <<%s missing terminating %s", startLine, delim, delim)
		}

		sop = &SOP{Content: strings.Join(body, "\n"), Line: startLine}
		for k := i; k <= j; k++ {
			lines[k] = ""
		}
		i = j
	}

	return []byte(strings.Join(lines, "\n")), sop, nil
}

// afterKeyword strips the leading `<keyword>` token (matched by exact length,
// so casing/spelling as tokenized by the caller is preserved) and following
// whitespace from a raw instruction line, per docs/agentfile.md — Dockerfile
// instruction arguments run to end of line; there is no inline-comment
// stripping inside an instruction's own argument text.
func afterKeyword(original, keyword string) string {
	if len(original) < len(keyword) {
		return ""
	}
	rest := original[len(keyword):]
	return strings.TrimLeft(rest, " \t")
}

func hasFlag(flags []string, want string) bool {
	for _, f := range flags {
		if f == want {
			return true
		}
	}
	return false
}

// parseKeyValuePairs parses `key=value key2="quoted value"` pairs per the
// IDENTITY grammar (key-value = key, "=", (quoted-string | bare-token)).
func parseKeyValuePairs(s string) (map[string]string, error) {
	out := map[string]string{}
	rest := s
	for {
		rest = strings.TrimLeft(rest, " \t")
		if rest == "" {
			break
		}
		eq := strings.IndexByte(rest, '=')
		if eq < 0 {
			return nil, fmt.Errorf("expected key=value, got %q", rest)
		}
		key := rest[:eq]
		rest = rest[eq+1:]
		var value string
		if strings.HasPrefix(rest, `"`) {
			end := strings.IndexByte(rest[1:], '"')
			if end < 0 {
				return nil, fmt.Errorf("unterminated quoted value for key %q", key)
			}
			value = rest[1 : 1+end]
			rest = rest[1+end+1:]
		} else {
			sp := strings.IndexAny(rest, " \t")
			if sp < 0 {
				value = rest
				rest = ""
			} else {
				value = rest[:sp]
				rest = rest[sp:]
			}
		}
		if key == "" {
			return nil, fmt.Errorf("empty key in %q", s)
		}
		out[key] = value
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("requires at least one key=value pair")
	}
	return out, nil
}

// parseRemoteAdd builds a RemoteAdd from an AST node already confirmed to
// have the --remote flag.
func parseRemoteAdd(child *parser.Node) (*RemoteAdd, error) {
	rest := afterKeyword(child.Original, child.Value)
	// rest still contains the flag tokens (Original is the raw, unsplit
	// line); strip exactly the flags recorded in child.Flags off the front.
	for _, flag := range child.Flags {
		rest = strings.TrimLeft(rest, " \t")
		if !strings.HasPrefix(rest, flag) {
			return nil, fmt.Errorf("line %d: internal error: expected flag %q in %q", child.StartLine, flag, rest)
		}
		rest = rest[len(flag):]
	}
	fields := strings.Fields(rest)
	if len(fields) != 2 {
		return nil, fmt.Errorf("line %d: ADD --remote requires exactly <source> <destination>, got %q", child.StartLine, rest)
	}

	ra := &RemoteAdd{Source: fields[0], Dest: fields[1], FailIfUnavailable: true, Cached: true, Line: child.StartLine}
	var sawCached, sawRuntime, sawFail, sawWarn bool
	for _, flag := range child.Flags {
		switch {
		case flag == "--cached":
			ra.Cached, ra.Runtime = true, false
			sawCached = true
		case flag == "--runtime":
			ra.Cached, ra.Runtime = false, true
			sawRuntime = true
		case flag == "--fail-if-unavailable":
			ra.FailIfUnavailable = true
			sawFail = true
		case flag == "--warn-if-unavailable":
			ra.FailIfUnavailable = false
			sawWarn = true
		case strings.HasPrefix(flag, "--chmod="):
			ra.Chmod = strings.TrimPrefix(flag, "--chmod=")
		case strings.HasPrefix(flag, "--chown="):
			ra.Chown = strings.TrimPrefix(flag, "--chown=")
		case flag == "--remote":
			// already used to select this branch
		default:
			return nil, fmt.Errorf("line %d: ADD --remote: unrecognized flag %q", child.StartLine, flag)
		}
	}
	if sawCached && sawRuntime {
		return nil, fmt.Errorf("line %d: ADD --remote: --cached and --runtime are mutually exclusive", child.StartLine)
	}
	if sawFail && sawWarn {
		return nil, fmt.Errorf("line %d: ADD --remote: --fail-if-unavailable and --warn-if-unavailable are mutually exclusive", child.StartLine)
	}
	return ra, nil
}
