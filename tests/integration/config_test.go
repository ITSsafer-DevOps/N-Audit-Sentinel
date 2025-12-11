package integration

import (
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/config"
	"os"
	"testing"
)

func TestGetEnvFallback(t *testing.T) {
	key := "NA_TEST_FOO"
	os.Unsetenv(key)
	v := config.GetEnv(key, "fallback")
	if v != "fallback" {
		t.Fatalf("expected fallback, got %q", v)
	}
}
