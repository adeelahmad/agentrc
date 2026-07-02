package agentfile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"strings"
)

// LabelOptions carries values BuildLabels needs but can't derive from a
// *File alone.
type LabelOptions struct {
	// SOPSHA256 is the digest of the /mnt/SOP file content. Required when
	// f.SOP.FileBacked is true (the content lives in a build-context file
	// internal/agentfile never reads); ignored for inline/heredoc SOP, whose
	// digest is computed directly from f.SOP.Content.
	SOPSHA256 string

	SourceURL string // -> org.opencontainers.image.source
	Revision  string // -> org.opencontainers.image.revision
	Created   string // -> org.opencontainers.image.created (RFC 3339)
}

// BuildLabels translates f into the org.agentrc.* label set per
// spec/index.md §9 and profiles/oci-package.md §3. Callers must populate
// f.LocalResources (from the real Dockerfile COPY/ADD instructions —
// internal/llb does this after parsing f.CleanedSource) and each RemoteAdd's
// ResolvedDigest (for successfully embedded --cached resources) before
// calling this.
func BuildLabels(f *File, opts LabelOptions) (map[string]string, error) {
	labels := map[string]string{}

	for k, v := range f.Identity {
		labels["org.agentrc.identity."+k] = v
	}
	for _, c := range f.Capabilities {
		labels["org.agentrc.capability."+c] = "true"
	}

	if f.SOP != nil {
		labels["org.agentrc.sop"] = "/mnt/SOP"
		digest := opts.SOPSHA256
		if !f.SOP.FileBacked {
			digest = hashString(f.SOP.Content)
		}
		if digest == "" {
			return nil, fmt.Errorf("SOP is file-backed but no digest was resolved")
		}
		labels["org.agentrc.sop.sha256"] = digest
	}

	for _, r := range f.LocalResources {
		if r.Kind == KindSOP || r.Kind == KindOther || r.Name == "" {
			continue
		}
		labels["org.agentrc."+string(r.Kind)+"."+r.Name] = "local"
	}

	for _, ra := range f.RemoteAdds {
		kind := ResourceKindForDest(ra.Dest)
		if kind == KindSOP || kind == KindOther {
			continue // SOP file-backed digest is handled above; KindOther isn't a labeled resource
		}
		name := path.Base(ra.Dest)
		prefix := "org.agentrc." + string(kind) + "." + name

		if ra.Runtime || ra.ResolvedDigest == "" {
			// --runtime, or a --cached resource that couldn't be resolved
			// (degraded per --warn-if-unavailable; --fail-if-unavailable
			// would already have aborted the build before this is called).
			labels[prefix] = "runtime:" + ra.Source
			continue
		}
		labels[prefix] = ra.ResolvedDigest
		labels[prefix+".origin"] = ra.Source
	}

	for _, p := range f.Policies {
		if p.Key == "network" {
			host, port, err := parseNetworkPolicyValue(p.Value)
			if err != nil {
				return nil, fmt.Errorf("line %d: POLICY network: %w", p.Line, err)
			}
			labels["org.agentrc.network.dns."+host] = port
			continue
		}

		labels["org.agentrc."+p.Key] = p.Value

		if isAutoEgressKey(p.Key) {
			if host, port, ok := hostPortFromURL(p.Value); ok {
				key := "org.agentrc.network.dns." + host
				if _, exists := labels[key]; !exists {
					labels[key] = port
					labels[key+".source"] = "auto:" + p.Key
				}
			}
		}
	}

	return labels, nil
}

// BuildOCIAnnotations builds the standard org.opencontainers.image.*
// annotations profiles/oci-package.md §5 says SHOULD accompany the agentrc
// labels. These are conventional metadata, not the authoritative manifest.
func BuildOCIAnnotations(f *File, opts LabelOptions) map[string]string {
	anno := map[string]string{}
	if v := f.Identity["name"]; v != "" {
		anno["org.opencontainers.image.title"] = v
	}
	if v := f.Identity["version"]; v != "" {
		anno["org.opencontainers.image.version"] = v
	}
	if v := f.Identity["author"]; v != "" {
		anno["org.opencontainers.image.authors"] = v
	}
	if opts.SourceURL != "" {
		anno["org.opencontainers.image.source"] = opts.SourceURL
	}
	if opts.Revision != "" {
		anno["org.opencontainers.image.revision"] = opts.Revision
	}
	if opts.Created != "" {
		anno["org.opencontainers.image.created"] = opts.Created
	}
	return anno
}

// isAutoEgressKey reports whether a POLICY key's URL value must auto-derive
// a network egress label, per spec/index.md §8.5.
func isAutoEgressKey(key string) bool {
	return key == "agent.interrupt_endpoint" || strings.HasPrefix(key, "agent.hooks.")
}

func hostPortFromURL(value string) (host, port string, ok bool) {
	u, err := url.Parse(value)
	if err != nil || u.Host == "" {
		return "", "", false
	}
	host = u.Hostname()
	port = u.Port()
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	return host, port, true
}

// parseNetworkPolicyValue parses a `dns:<host>:<port>` POLICY network value.
func parseNetworkPolicyValue(value string) (host, port string, err error) {
	rest, ok := strings.CutPrefix(value, "dns:")
	if !ok {
		return "", "", fmt.Errorf("expected dns:<host>:<port>, got %q", value)
	}
	idx := strings.LastIndex(rest, ":")
	if idx < 0 {
		return "", "", fmt.Errorf("expected dns:<host>:<port>, got %q", value)
	}
	host, port = rest[:idx], rest[idx+1:]
	if host == "" || port == "" {
		return "", "", fmt.Errorf("expected dns:<host>:<port>, got %q", value)
	}
	return host, port, nil
}

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
