package logger

import (
	"bytes"
	"regexp"
	"testing"
	"time"
)

func TestStripANSI(t *testing.T) {
	in := []byte("\x1b[31mred\x1b[0m\n")
	out := StripANSI(in)
	if bytes.Contains(out, []byte("\x1b")) {
		t.Fatalf("ansi not stripped: %q", out)
	}
}

func TestTimestampedWriter(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTimestampedWriter(&buf)
	line := "hello world\n"
	if _, err := tw.Write([]byte(line)); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	// Expect timestamp prefix followed by space and the line
	// timestamp is RFC3339Nano; check prefix is a timestamp
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)
	if !re.MatchString(out) {
		t.Fatalf("missing timestamp prefix: %q", out)
	}
	if !bytes.Contains([]byte(out), []byte("hello world")) {
		t.Fatalf("missing line content: %q", out)
	}
	// ensure timestamp parseable
	parts := bytes.SplitN([]byte(out), []byte(" "), 2)
	if len(parts) < 2 {
		t.Fatalf("unexpected output format")
	}
	if _, err := time.Parse(time.RFC3339Nano, string(parts[0])); err != nil {
		t.Fatalf("timestamp parse error: %v", err)
	}
}
