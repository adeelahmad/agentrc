package agentfile

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ValidateRemoteSourceURL rejects ADD --remote sources that would fetch from
// a loopback, link-local, private, or otherwise internal-looking address —
// best-effort SSRF mitigation for build-time fetches of a URL taken
// directly from Agentfile content (which, in a shared build service, may be
// attacker-controlled). This is not a substitute for network-level egress
// controls: the resolved IP is re-checked at request time, not pinned, so a
// DNS-rebinding attacker who controls both the DNS answer and the timing
// could still slip through between this check and the actual connection.
// Callers that follow redirects MUST re-run this check against every hop.
func ValidateRemoteSourceURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("only http/https URLs can be fetched at build time, got scheme %q", u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("URL has no host")
	}
	if strings.EqualFold(host, "localhost") {
		return fmt.Errorf("refusing to fetch from localhost")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("resolving host %q: %w", host, err)
	}
	for _, ip := range ips {
		if isDisallowedFetchTarget(ip) {
			return fmt.Errorf("refusing to fetch from %s: resolves to disallowed address %s", host, ip)
		}
	}
	return nil
}

func isDisallowedFetchTarget(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsPrivate() ||
		ip.IsUnspecified() ||
		ip.IsMulticast()
}
