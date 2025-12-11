package validation

import (
	"testing"
)

func TestValidateScope_ValidIPs(t *testing.T) {
	ips := []string{"192.168.1.1", "10.0.0.0/8"}
	domains := []string{}
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(validIPs) != 2 {
		t.Fatalf("expected 2 valid IPs, got %d", len(validIPs))
	}
	if len(warnings) != 0 {
		t.Fatalf("expected 0 warnings, got %d: %v", len(warnings), warnings)
	}
	if len(validDomains) != 0 {
		t.Fatalf("expected 0 valid domains, got %d", len(validDomains))
	}
}

func TestValidateScope_ValidDomains(t *testing.T) {
	ips := []string{}
	domains := []string{"example.com", "test.org"}
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(validDomains) != 2 {
		t.Fatalf("expected 2 valid domains, got %d", len(validDomains))
	}
	if len(warnings) != 0 {
		t.Fatalf("expected 0 warnings, got %d: %v", len(warnings), warnings)
	}
	if len(validIPs) != 0 {
		t.Fatalf("expected 0 valid IPs, got %d", len(validIPs))
	}
}

func TestValidateScope_InvalidInputs(t *testing.T) {
	ips := []string{"invalid-ip", "256.256.256.256"}
	domains := []string{".example", "example.", ""}
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(validIPs) != 0 {
		t.Fatalf("expected 0 valid IPs for invalid inputs, got %d", len(validIPs))
	}
	if len(validDomains) != 0 {
		t.Fatalf("expected 0 valid domains for invalid inputs, got %d", len(validDomains))
	}
	if len(warnings) < 3 {
		t.Fatalf("expected at least 3 warnings, got %d", len(warnings))
	}
}

func TestValidateScope_MixedInputs(t *testing.T) {
	ips := []string{"192.168.1.1", "invalid-ip", "10.0.0.0/8"}
	domains := []string{"example.com", ".test.org", "test.org"}
	validIPs, validDomains, warnings := ValidateScope(ips, domains)
	if len(validIPs) != 2 {
		t.Fatalf("expected 2 valid IPs, got %d", len(validIPs))
	}
	if len(validDomains) != 2 {
		t.Fatalf("expected 2 valid domains, got %d", len(validDomains))
	}
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestValidateScope_SingleIPConverted(t *testing.T) {
	ips := []string{"8.8.8.8"}
	domains := []string{}
	validIPs, _, _ := ValidateScope(ips, domains)
	if len(validIPs) != 1 {
		t.Fatalf("expected 1 valid IP, got %d", len(validIPs))
	}
	if validIPs[0] != "8.8.8.8/32" {
		t.Fatalf("expected single IP converted to /32, got %s", validIPs[0])
	}
}
