package unit

import (
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/logging"
	"testing"
)

func TestStripANSI(t *testing.T) {
	s := "Hello \x1b[31mRed\x1b[0m World"
	out := logging.StripANSI(s)
	if out != "Hello Red World" {
		t.Fatalf("unexpected sanitized output: %q", out)
	}
}
