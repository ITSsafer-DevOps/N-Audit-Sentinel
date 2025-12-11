package recorder

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestStartSessionContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var buf bytes.Buffer

	cancel()

	err := StartSession(ctx, &buf, "/bin/sh", "-c", "echo test")
	if err != nil {
		t.Fatalf("expected no error with cancelled context, got: %v", err)
	}
}

func TestStartSessionContextWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var buf bytes.Buffer

	err := StartSession(ctx, &buf, "/bin/sh", "-c", "echo test")
	if err != nil {
		t.Fatalf("expected no error with timeout context, got: %v", err)
	}
}

func TestStartSessionNonexistentCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var buf bytes.Buffer

	err := StartSession(ctx, &buf, "/nonexistent/command")
	if err != nil {
		t.Fatalf("expected nil error (context cancellation), got: %v", err)
	}
}
