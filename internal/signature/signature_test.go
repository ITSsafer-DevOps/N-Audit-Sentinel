package signature

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSealLogFile_NonexistentLog(t *testing.T) {
	tmpdir := t.TempDir()
	nonexistentLog := filepath.Join(tmpdir, "nonexistent.log")
	privKey := filepath.Join(tmpdir, "nonexistent.key")

	err := SealLogFile(nonexistentLog, privKey)
	if err == nil {
		t.Fatalf("expected error for nonexistent log file, got nil")
	}
	if !contains(err.Error(), "read log file") {
		t.Fatalf("expected 'read log file' error, got: %v", err)
	}
}

func TestSealLogFile_NonexistentKey(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "test.log")
	if err := os.WriteFile(logFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	nonexistentKey := filepath.Join(tmpdir, "nonexistent.key")
	err := SealLogFile(logFile, nonexistentKey)
	if err == nil {
		t.Fatalf("expected error for nonexistent private key, got nil")
	}
	if !contains(err.Error(), "read private key") {
		t.Fatalf("expected 'read private key' error, got: %v", err)
	}
}

func TestSealLogFile_InvalidKey(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "test.log")
	if err := os.WriteFile(logFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	keyFile := filepath.Join(tmpdir, "invalid.key")
	if err := os.WriteFile(keyFile, []byte("not a valid ssh key"), 0600); err != nil {
		t.Fatalf("failed to write invalid key: %v", err)
	}

	err := SealLogFile(logFile, keyFile)
	if err == nil {
		t.Fatalf("expected error for invalid SSH key, got nil")
	}
	if !contains(err.Error(), "parse private key") {
		t.Fatalf("expected 'parse private key' error, got: %v", err)
	}
}

func TestSealLogFile_LogNotWritable(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "readonly.log")

	if err := os.WriteFile(logFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	// Make the log file read-only
	if err := os.Chmod(logFile, 0444); err != nil {
		t.Fatalf("failed to chmod: %v", err)
	}
	defer os.Chmod(logFile, 0644) // restore for cleanup

	keyFile := filepath.Join(tmpdir, "key.pem")
	if err := os.WriteFile(keyFile, []byte("dummy"), 0600); err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	err := SealLogFile(logFile, keyFile)
	// Will fail at SSH parsing or file append, but should error
	if err == nil {
		t.Fatalf("expected error for read-only log with invalid key, got nil")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
