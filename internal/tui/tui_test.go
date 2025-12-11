package tui

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestShowBanner(t *testing.T) {
	var b bytes.Buffer
	ShowBanner(&b)
	out := b.String()
	if !strings.Contains(out, "Developer:") {
		t.Error("banner missing Developer")
	}
}

func TestGetPentesterInfo_Parse(t *testing.T) {
	in := strings.NewReader("Alice\nalice@example.com\n")
	var out bytes.Buffer
	name, email, err := GetPentesterInfo(in, &out)
	if err != nil || name != "Alice" || email != "alice@example.com" {
		t.Errorf("GetPentesterInfo failed: %v, %q, %q", err, name, email)
	}
}

func TestGetPentesterInfo_EOFFirstLine(t *testing.T) {
	in := strings.NewReader("")
	var out bytes.Buffer
	_, _, err := GetPentesterInfo(in, &out)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestGetPentesterInfo_EOFSecondLine(t *testing.T) {
	in := strings.NewReader("Bob\n")
	var out bytes.Buffer
	name, _, err := GetPentesterInfo(in, &out)
	if err != io.EOF || name != "Bob" {
		t.Errorf("expected EOF with Bob, got %v/%q", err, name)
	}
}

func TestGetPentesterInfo_Whitespace(t *testing.T) {
	in := strings.NewReader("  Charlie  \n  client@test  \n")
	var out bytes.Buffer
	name, email, err := GetPentesterInfo(in, &out)
	if err != nil || name != "Charlie" || email != "client@test" {
		t.Errorf("whitespace trimming failed: %q %q", name, email)
	}
}

func TestGetScope_Basic(t *testing.T) {
	in := strings.NewReader("10.0.0.1\n\nexample.com\n\n")
	var out bytes.Buffer
	ips, domains, err := GetScope(in, &out)
	if err != nil || (len(ips) == 0 && len(domains) == 0) {
		t.Errorf("GetScope failed")
	}
}

func TestGetScope_EOFHandle(t *testing.T) {
	in := strings.NewReader("192.168.1.1\n")
	var out bytes.Buffer
	ips, _, err := GetScope(in, &out)
	if err != nil || len(ips) == 0 {
		t.Errorf("GetScope EOF failed")
	}
}

func TestGetScope_Multiple(t *testing.T) {
	in := strings.NewReader("10.0.0.1\n10.0.0.2\n\ntest.com\n\n")
	var out bytes.Buffer
	ips, _, err := GetScope(in, &out)
	if err != nil || len(ips) < 1 {
		t.Errorf("GetScope multiple failed")
	}
}

func TestGetScope_Empty(t *testing.T) {
	in := strings.NewReader("\n\n\n")
	var out bytes.Buffer
	_, _, err := GetScope(in, &out)
	if err != nil {
		t.Errorf("GetScope empty failed")
	}
}

func TestGetScope_Whitespace(t *testing.T) {
	in := strings.NewReader("10.0.0.1\n   \n\nexample.com\n\n")
	var out bytes.Buffer
	ips, _, err := GetScope(in, &out)
	if err != nil || len(ips) == 0 {
		t.Errorf("GetScope whitespace failed")
	}
}

func TestGetScope_CIDR(t *testing.T) {
	in := strings.NewReader("172.16.0.0/16\n10.0.0.0/8\n\n")
	var out bytes.Buffer
	ips, _, err := GetScope(in, &out)
	if err != nil || len(ips) < 1 {
		t.Errorf("GetScope CIDR failed")
	}
}

func TestGetScope_LongDomains(t *testing.T) {
	in := strings.NewReader("10.0.0.1\n\nvery.long.subdomain.example.internal.domain.com\nshort.io\n\n")
	var out bytes.Buffer
	_, _, err := GetScope(in, &out)
	if err != nil {
		t.Errorf("GetScope long domains failed")
	}
}

func TestGetScope_PromptCheck(t *testing.T) {
	in := strings.NewReader("10.0.0.1\n\ntest.com\n\n")
	var out bytes.Buffer
	_, _, err := GetScope(in, &out)
	if err != nil {
		t.Errorf("GetScope prompt check failed")
	}
	if !strings.Contains(out.String(), "Layer 3") {
		t.Error("Layer 3 prompt missing")
	}
}
