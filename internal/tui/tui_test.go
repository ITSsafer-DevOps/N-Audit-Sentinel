package tui

import (
	"bytes"
	"strings"
	"testing"
)

// TestAsciiBannerContainsMagenta ensures the banner includes the magenta ANSI escape.
func TestAsciiBannerContainsMagenta(t *testing.T) {
	want := "\x1b[35m"
	if !strings.Contains(asciiBanner, want) {
		t.Fatalf("asciiBanner does not contain magenta ANSI escape %q", want)
	}
}

func TestGetPentesterInfo(t *testing.T) {
	input := "Alice\nACME Corp\n"
	in := bytes.NewBufferString(input)
	out := &bytes.Buffer{}
	pentester, client, err := GetPentesterInfo(in, out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pentester != "Alice" {
		t.Fatalf("expected pentester Alice got %q", pentester)
	}
	if client != "ACME Corp" {
		t.Fatalf("expected client ACME Corp got %q", client)
	}
	outStr := out.String()
	if !strings.Contains(outStr, "Pentester Name:") || !strings.Contains(outStr, "Client Name:") {
		t.Errorf("missing prompts in output: %q", outStr)
	}
}

func TestGetScopeDoubleEnter(t *testing.T) {
	inputs := []string{
		"192.168.0.0/24",
		"10.0.0.1",
		"", // first empty
		"", // second empty -> transition to domains
		"example.com",
		"sub.example.com",
		"", // first empty
		"", // second empty -> finish
	}
	in := bytes.NewBufferString(strings.Join(inputs, "\n") + "\n")
	out := &bytes.Buffer{}
	ips, domains, err := GetScope(in, out)
	if err != nil {
		t.Fatalf("GetScope error: %v", err)
	}
	if len(ips) != 2 || ips[0] != "192.168.0.0/24" || ips[1] != "10.0.0.1" {
		t.Fatalf("unexpected ips: %#v", ips)
	}
	if len(domains) != 2 || domains[0] != "example.com" || domains[1] != "sub.example.com" {
		t.Fatalf("unexpected domains: %#v", domains)
	}
	outStr := out.String()
	if strings.Count(outStr, "IP/CIDR>") < 2 {
		t.Errorf("expected multiple IP/CIDR prompts; output: %q", outStr)
	}
	if strings.Count(outStr, "Domain>") < 2 {
		t.Errorf("expected multiple Domain prompts; output: %q", outStr)
	}
}
