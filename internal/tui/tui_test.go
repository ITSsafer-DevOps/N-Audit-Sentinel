package tui

import (
	"bytes"
	"strings"
	"testing"
)

func TestShowBanner(t *testing.T) {
	var buf bytes.Buffer
	ShowBanner(&buf)
	out := buf.String()
	if !strings.Contains(out, "Developer:") || !strings.Contains(out, "Company:") {
		t.Fatalf("banner missing metadata: %s", out)
	}
}

func TestGetPentesterInfo(t *testing.T) {
	in := "Alice\nAcme Corp\n"
	var buf bytes.Buffer
	p, c, err := GetPentesterInfo(strings.NewReader(in), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != "Alice" || c != "Acme Corp" {
		t.Fatalf("unexpected result: %q %q", p, c)
	}
}

func TestGetScopeDoubleEnter(t *testing.T) {
	// IP lines: 1.2.3.4, double enter to move to domains, domain example.com, double enter to finish
	in := "1.2.3.4\n\n\nexample.com\n\n\n"
	var buf bytes.Buffer
	ips, domains, err := GetScope(strings.NewReader(in), &buf)
	if err != nil {
		t.Fatalf("GetScope error: %v", err)
	}
	if len(ips) != 1 || ips[0] != "1.2.3.4" {
		t.Fatalf("unexpected ips: %v", ips)
	}
	if len(domains) != 1 || domains[0] != "example.com" {
		t.Fatalf("unexpected domains: %v", domains)
	}
}
