package config

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

func TestGetEnvFallbackAlternate(t *testing.T) {
    key := "NA_TEST_FOO"
    os.Unsetenv(key)
    v := GetEnv(key, "bar")
    if v != "bar" {
        t.Fatalf("expected fallback 'bar', got '%s'", v)
    }
    os.Setenv(key, "baz")
    v2 := GetEnv(key, "bar")
    if v2 != "baz" {
        t.Fatalf("expected env 'baz', got '%s'", v2)
    }
    os.Unsetenv(key)
}
        t.Fatalf("expected env 'baz', got '%s'", v2)
