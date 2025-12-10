package logger

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestStripANSI(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  string
	}{
		{"color_codes", "\x1b[31mred text\x1b[0m", "red text"},
		{"no_codes", "plain", "plain"},
		{"empty", "", ""},
		{"mixed", "start \x1b[1;32mGREEN\x1b[0m end", "start GREEN end"},
	}
	for _, tc := range cases {
		got := string(StripANSI([]byte(tc.in)))
		if got != tc.out {
			// show escaped for clarity
			if len(got) != len(tc.out) || got != tc.out {
				// Provide diff style output
				t.Errorf("%s: expected %q got %q", tc.name, tc.out, got)
			}
		}
	}
}

func TestTimestampedWriter(t *testing.T) {
	under := &bytes.Buffer{}
	w := NewTimestampedWriter(under)

	// Write partial then complete lines to test buffering.
	w.Write([]byte("partial"))
	if under.Len() != 0 {
		// No newline yet, should not flush.
		// Accept but note if output exists unexpectedly.
		if strings.Count(under.String(), "\n") > 0 {
			// This would indicate premature flush.
			t.Errorf("unexpected flushed line before newline")
		}
	}
	w.Write([]byte(" line one\nline two"))
	w.Write([]byte(" continued\nfinal line\n"))

	output := under.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 3 { // Expect 3 complete lines
		// Lines: "partial line one", "line two continued", "final line"
		// partial + line one merged due to no newline after first write.
		// second line constructed across two writes.
		// third line final.
		// If mismatch, fail.
		t.Fatalf("expected 3 lines got %d: %q", len(lines), output)
	}

	// Simple RFC3339Nano-like regex: 2025-01-01T12:30:00.123456789Z <rest>
	tsRe := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[^ ]+Z `)
	for i, line := range lines {
		if !tsRe.MatchString(line) {
			t.Errorf("line %d missing timestamp prefix: %q", i, line)
		}
	}

	// Verify content parts (without strict timestamp checks)
	contents := make([]string, 0, len(lines))
	for _, line := range lines {
		c := tsRe.ReplaceAllString(line, "")
		contents = append(contents, c)
	}
	expected := []string{
		"partial line one",   // first merged line
		"line two continued", // second merged line from split writes
		"final line",         // third line
	}
	for i := range expected {
		if contents[i] != expected[i]+"\n" { // newline should be retained per original lines
			// Because we kept newline in writing, content ends with \n.
			if contents[i] != expected[i] { // fallback if newline trimmed somehow
				t.Errorf("line %d expected %q got %q", i, expected[i]+"\\n", contents[i])
			}
		}
	}
}
