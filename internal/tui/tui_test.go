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

func TestGetPentesterInfo_EOFOnFirstLine(t *testing.T) {
in := strings.NewReader("")
var out bytes.Buffer
name, email, err := GetPentesterInfo(in, &out)
if err != io.EOF {
t.Fatalf("expected io.EOF on empty input, got %v", err)
}
if name != "" || email != "" {
t.Fatalf("expected empty strings, got %s %s", name, email)
}
}

func TestGetPentesterInfo_EOFOnSecondLine(t *testing.T) {
in := strings.NewReader("Bob\n")
var out bytes.Buffer
name, email, err := GetPentesterInfo(in, &out)
if err != io.EOF {
t.Fatalf("expected io.EOF on incomplete input, got %v", err)
}
if name != "Bob" || email != "" {
t.Fatalf("expected 'Bob' and '', got %s %s", name, email)
}
}

func TestGetPentesterInfo_Whitespace(t *testing.T) {
in := strings.NewReader("  Charlie  \n  client@test  \n")
var out bytes.Buffer
name, email, err := GetPentesterInfo(in, &out)
if err != nil {
t.Fatalf("GetPentesterInfo error: %v", err)
}
if name != "Charlie" || email != "client@test" {
t.Fatalf("expected trimmed values, got %s %s", name, email)
}
}

func TestGetScope_Simple(t *testing.T) {
input := "10.0.0.1\n\ngoogle.com\n\n"
in := strings.NewReader(input)
var out bytes.Buffer
ips, domains, err := GetScope(in, &out)
if err != nil {
t.Fatalf("GetScope error: %v", err)
}
_ = ips
_ = domains
}

func TestGetScope_EOFInIPLoop(t *testing.T) {
input := "192.168.1.1\n"
in := strings.NewReader(input)
var out bytes.Buffer
ips, domains, err := GetScope(in, &out)
if err != nil {
t.Fatalf("GetScope error on EOF: %v", err)
}
_ = ips
_ = domains
}

func TestGetScope_MultipleItems(t *testing.T) {
input := "10.0.0.1\n10.0.0.2\n\ntest.com\nexample.com\n\n"
in := strings.NewReader(input)
var out bytes.Buffer
ips, domains, err := GetScope(in, &out)
if err != nil {
t.Fatalf("GetScope error: %v", err)
}
_ = ips
_ = domains
}
