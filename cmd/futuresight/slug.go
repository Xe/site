package main

import "strings"

// dnsLabelMax is the maximum length of a single DNS label (RFC 1035).
const dnsLabelMax = 63

// slugifyBranch converts a git branch name into a DNS-safe label: lowercase
// ASCII alphanumerics, with every other run of characters collapsed to a single
// hyphen, trimmed of leading/trailing hyphens, and truncated to one DNS label.
//
//	feat/Foo_Bar -> feat-foo-bar
func slugifyBranch(branch string) string {
	var b strings.Builder
	lastHyphen := true // leading hyphens are trimmed

	for _, r := range strings.ToLower(branch) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastHyphen = false
		default:
			if !lastHyphen {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}

	slug := strings.TrimRight(b.String(), "-")

	if len(slug) > dnsLabelMax {
		slug = strings.TrimRight(slug[:dnsLabelMax], "-")
	}

	return slug
}
