package validation

import (
	"reflect"
	"testing"
)

func TestValidateScope_ValidIPsAndCIDRs(t *testing.T) {
	ips := []string{"192.168.1.10", "10.0.0.0/24", "2001:db8::1"}
	var domains []string
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	wantIPs := []string{"192.168.1.10/32", "10.0.0.0/24", "2001:db8::1/32"}
	if !reflect.DeepEqual(validIPs, wantIPs) {
		t.Fatalf("validIPs mismatch: got %v want %v", validIPs, wantIPs)
	}
	if len(validDomains) != 0 {
		t.Fatalf("expected no domains, got %v", validDomains)
	}
}

func TestValidateScope_InvalidIPsAndCIDRs(t *testing.T) {
	ips := []string{"999.999.999.999", "10.0.0.0/"}
	var domains []string
	validIPs, _, warnings := ValidateScope(ips, domains)
	if len(validIPs) != 0 {
		t.Fatalf("expected no valid IPs, got %v", validIPs)
	}
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings, got %d (%v)", len(warnings), warnings)
	}
}

func TestValidateScope_ValidDomains(t *testing.T) {
	var ips []string
	domains := []string{"example.com", "sub.domain.org"}
	_, validDomains, warnings := ValidateScope(ips, domains)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	wantDomains := []string{"example.com", "sub.domain.org"}
	if !reflect.DeepEqual(validDomains, wantDomains) {
		t.Fatalf("validDomains mismatch: got %v want %v", validDomains, wantDomains)
	}
}

func TestValidateScope_InvalidDomains(t *testing.T) {
	var ips []string
	domains := []string{"bad/domain", "localhost", ".leadingdot.com", "trailingdot.com."}
	_, validDomains, warnings := ValidateScope(ips, domains)
	if len(validDomains) != 0 {
		t.Fatalf("expected no valid domains, got %v", validDomains)
	}
	if len(warnings) != 4 {
		t.Fatalf("expected 4 warnings, got %d (%v)", len(warnings), warnings)
	}
}

func TestValidateScope_MixedInputs(t *testing.T) {
	ips := []string{"10.0.0.1", "badip", "10.0.0.0/24", "10.0.0.0/"}
	domains := []string{"example.com", "bad/domain", "no-tld", "sub.valid.org"}
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	wantIPs := []string{"10.0.0.1/32", "10.0.0.0/24"}
	wantDomains := []string{"example.com", "sub.valid.org"}
	if !reflect.DeepEqual(validIPs, wantIPs) {
		t.Fatalf("validIPs mismatch: got %v want %v", validIPs, wantIPs)
	}
	if !reflect.DeepEqual(validDomains, wantDomains) {
		t.Fatalf("validDomains mismatch: got %v want %v", validDomains, wantDomains)
	}
	if len(warnings) != 4 {
		t.Fatalf("expected 4 warnings, got %d (%v)", len(warnings), warnings)
	}
}

func TestValidateScope_EmptyInputs(t *testing.T) {
	var ips []string
	var domains []string
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(validIPs) != 0 || len(validDomains) != 0 || len(warnings) != 0 {
		t.Fatalf("expected all empty, got ips=%v domains=%v warnings=%v", validIPs, validDomains, warnings)
	}
}
