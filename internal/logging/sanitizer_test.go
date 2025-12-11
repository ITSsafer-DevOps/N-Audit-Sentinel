package logging

import "testing"

func TestStripANSI(t *testing.T) {
	in := "\x1b[31mERROR\x1b[0m: something happened"
	out := StripANSI(in)
	if out == in {
		t.Fatalf("expected ANSI to be stripped, got same string")
	}
	if out != "ERROR: something happened" {
		t.Fatalf("unexpected sanitized output: %q", out)
	}
}
