package unit

import (
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/seal"
	"testing"
)

func TestHashSHA256(t *testing.T) {
	data := []byte("hello world")
	h := seal.HashSHA256(data)
	if len(h) != 64 {
		t.Fatalf("unexpected hash length: %d", len(h))
	}
}
