package signature

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSealLogFile_Basic(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "test.log")
	if err := os.WriteFile(logFile, []byte("test log content"), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	privKeyFile := filepath.Join(tmpdir, "test_key")
	testKey := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUtbm9uZS1ub25lAAAAAAAAAEMAAAAAy8qR
-----END OPENSSH PRIVATE KEY-----`
	if err := os.WriteFile(privKeyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("failed to write test private key: %v", err)
	}

	err := SealLogFile(logFile, privKeyFile)
	if err != nil {
		t.Logf("Expected error with test key: %v", err)
	}
}

func TestSealLogFile_NonexistentLog(t *testing.T) {
	tmpdir := t.TempDir()
	nonexistentLog := filepath.Join(tmpdir, "nonexistent.log")
	privKey := filepath.Join(tmpdir, "nonexistent.key")

	err := SealLogFile(nonexistentLog, privKey)
	if err == nil {
		t.Fatalf("expected error for nonexistent log file, got nil")
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
}
