package config

import (
	"os"
	"testing"
)

func TestGetEnvWithFallback(t *testing.T) {
	key := "N_AUDIT_TEST_FOO"
	os.Unsetenv(key)
	got := GetEnv(key, "bar")
	if got != "bar" {
		t.Fatalf("expected fallback 'bar', got %q", got)
	}
	os.Setenv(key, "baz")
	got2 := GetEnv(key, "bar")
	if got2 != "baz" {
		t.Fatalf("expected env 'baz', got %q", got2)
	}
	os.Unsetenv(key)
}
