package recorder

import (
    "context"
    "bytes"
    "testing"
    "time"
)

func TestStartSession_ImmediateCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    var buf bytes.Buffer
    // Should return quickly with nil on cancelled context
    if err := StartSession(ctx, &buf, "/bin/sh", "-c", "echo hi"); err != nil {
        t.Fatalf("expected nil on immediate cancel, got %v", err)
    }
}

func TestStartSession_RunAndCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    var buf bytes.Buffer

    done := make(chan error, 1)
    go func() {
        done <- StartSession(ctx, &buf, "/bin/sh", "-c", "sleep 5")
    }()

    // let the shell start
    time.Sleep(200 * time.Millisecond)
    cancel()

    select {
    case err := <-done:
        if err != nil && err != context.Canceled {
            t.Fatalf("expected nil or context.Canceled, got %v", err)
        }
    case <-time.After(3 * time.Second):
        t.Fatalf("StartSession did not return after cancel")
    }
}
