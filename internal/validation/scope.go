// Package validation provides helpers for validating scope inputs (IPs/CIDRs and domains).
package validation

import (
	"net"
	"net/url"
	"strings"
)

// ValidateScope validates raw IP/CIDR and domain inputs.
// Returns normalized valid IPs/CIDRs, valid domains, and human-readable warnings for any skipped entries.
func ValidateScope(ips, domains []string) (validIPs, validDomains []string, warnings []string) {
	// Validate IPs/CIDRs
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		if strings.Contains(ip, "/") {
			if _, _, err := net.ParseCIDR(ip); err == nil {
				validIPs = append(validIPs, ip)
			} else {
				warnings = append(warnings, "Invalid CIDR skipped: "+ip)
			}
		} else {
			if parsed := net.ParseIP(ip); parsed != nil {
				validIPs = append(validIPs, ip+"/32")
			} else {
				warnings = append(warnings, "Invalid IP skipped: "+ip)
			}
		}
	}

	// Validate domains
	for _, d := range domains {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		if strings.Contains(d, "/") || strings.HasPrefix(d, ".") || strings.HasSuffix(d, ".") {
			warnings = append(warnings, "Invalid domain skipped: "+d)
			continue
		}
		if !strings.Contains(d, ".") {
			warnings = append(warnings, "Domain without TLD skipped: "+d)
			continue
		}
		if _, err := url.Parse("https://" + d); err != nil {
			warnings = append(warnings, "Unparsable domain skipped: "+d)
			continue
		}
		validDomains = append(validDomains, d)
	}
	return validIPs, validDomains, warnings
}
