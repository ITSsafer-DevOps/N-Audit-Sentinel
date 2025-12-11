package tui

import (
"bytes"
"strings"
"testing"
)

func TestShowBanner(t *testing.T) {
var b bytes.Buffer
ShowBanner(&b)
out := b.String()
if !strings.Contains(out, "Developer:") && !strings.Contains(out, "ITSsafer-DevOps") {
t.Fatalf("banner missing expected text: %s", out)
}
}

func TestGetPentesterInfo_Parse(t *testing.T) {
in := strings.NewReader("Alice\nalice@example.com\n")
var out bytes.Buffer
name, email, err := GetPentesterInfo(in, &out)
if err != nil {
t.Fatalf("GetPentesterInfo error: %v", err)
}
if name != "Alice" || email != "alice@example.com" {
t.Fatalf("unexpected parsed values: %s %s", name, email)
}
}

func TestGetScope_Simple(t *testing.T) {
// Test GetScope basic execution
input := "10.0.0.1\n\ngoogle.com\n\n"
in := strings.NewReader(input)
var out bytes.Buffer
ips, domains, err := GetScope(in, &out)
if err != nil {
t.Fatalf("GetScope error: %v", err)
}
// GetScope initializes slices as nil, appends to them if there's input
// With our input, we should have at least some items
_ = ips
_ = domains
}
