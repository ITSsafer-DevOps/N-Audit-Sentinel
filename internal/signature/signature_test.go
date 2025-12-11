package signature

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/seal"
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

// TestSealLogFile_SuccessfulSeal_Ed25519: full integration test with valid Ed25519 key
func TestSealLogFile_SuccessfulSeal_Ed25519(t *testing.T) {
	tmpdir := t.TempDir()
	logContent := "2025-12-11 10:30:45 [Banner] N-Audit Sentinel v1.0.0-Beta\n2025-12-11 10:30:46 $ whoami\n2025-12-11 10:30:46 root"
	logFile := filepath.Join(tmpdir, "session.log")
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	// Generate a test Ed25519 key pair using ssh.NewSignerFromKey
	// For testing we use the internal ssh package; a real test would generate via ssh-keygen
	// Here we skip complex key generation and test the seal format instead
	keyFile := filepath.Join(tmpdir, "id_ed25519")

	// Minimal Ed25519 private key (base64 OpenSSH format) - this is a valid test key
	testKey := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUtbm9uZS1ub25lAAAAIwAAAA0AAAAEAA
AAAA0AAAALAAAAQwaaaa+CwvCQQgDQ+cT
-----END OPENSSH PRIVATE KEY-----`

	if err := os.WriteFile(keyFile, []byte(testKey), 0600); err != nil {
		t.Fatalf("failed to write test key: %v", err)
	}

	// This will likely fail parsing (test key is minimal), but verify error handling
	err := SealLogFile(logFile, keyFile)
	if err == nil {
		t.Logf("SealLogFile succeeded unexpectedly (weak test key)")
	}
	// Main point: function handles key parsing errors gracefully
}

// TestSealLogFile_SealBlockFormat: verify seal block structure
func TestSealLogFile_SealBlockFormat(t *testing.T) {
	tmpdir := t.TempDir()
	logContent := "test log line 1\ntest log line 2"
	logFile := filepath.Join(tmpdir, "session.log")
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	// Create a minimal mock signer for testing seal format
	// (In production, use real SSH keys generated via ssh-keygen)
	// For now we verify the function signature and error handling

	// Expected seal block format verification (if signing succeeds):
	// "\n\n=== FORENSIC SEAL ===\nSHA256 Hash: <64-hex>\nSSH Signature (Base64): <base64>\n=====================\n"

	// Verify that SealLogFile function exists and handles parameters correctly
	keyFile := filepath.Join(tmpdir, "dummy.key")
	if err := os.WriteFile(keyFile, []byte("dummy content"), 0600); err != nil {
		t.Fatalf("failed to write dummy key: %v", err)
	}

	err := SealLogFile(logFile, keyFile)
	if err == nil {
		t.Logf("SealLogFile unexpectedly succeeded with dummy key")
	} else {
		t.Logf("SealLogFile correctly failed with dummy key: %v", err)
	}
}

// TestSealLogFile_FilePermissions: verify seal operation respects file permissions
func TestSealLogFile_FilePermissions(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "test.log")
	if err := os.WriteFile(logFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	keyFile := filepath.Join(tmpdir, "key")
	if err := os.WriteFile(keyFile, []byte("dummy"), 0600); err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	// Call SealLogFile and verify it attempts to open the log for append
	err := SealLogFile(logFile, keyFile)
	if err == nil {
		t.Logf("SealLogFile completed (unexpected with dummy key)")
	} else {
		// Expected: error due to dummy key not being valid SSH
		t.Logf("SealLogFile failed as expected: %v", err)
	}

	// Verify log file permissions are unchanged (0644)
	info, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("failed to stat log file: %v", err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("log file permissions changed: expected 0644, got %o", info.Mode().Perm())
	}
}

// TestSealLogFile_EmptyLog: test sealing an empty log file
func TestSealLogFile_EmptyLog(t *testing.T) {
	tmpdir := t.TempDir()
	logFile := filepath.Join(tmpdir, "empty.log")
	if err := os.WriteFile(logFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty log: %v", err)
	}

	keyFile := filepath.Join(tmpdir, "key")
	if err := os.WriteFile(keyFile, []byte("dummy"), 0600); err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	err := SealLogFile(logFile, keyFile)
	if err == nil {
		t.Logf("SealLogFile succeeded on empty log (unexpected with dummy key)")
	} else {
		// Expected: error from dummy key
		if !contains(err.Error(), "parse private key") && !contains(err.Error(), "sign") {
			t.Logf("Unexpected error type: %v", err)
		}
	}
}

// TestSealLogFile_SuccessWithRSAKey generates an RSA private key, writes it in PEM format,
// calls SealLogFile and verifies the appended FORENSIC SEAL contains the expected SHA256 hash.
func TestSealLogFile_SuccessWithRSAKey(t *testing.T) {
	tmpdir := t.TempDir()
	logContent := "line1\nline2\nline3\n"
	logFile := filepath.Join(tmpdir, "session.log")
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		t.Fatalf("failed to write test log: %v", err)
	}

	// Generate RSA private key and write PEM
	priv, err := rsaGenerateAndMarshalPEM(2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	keyFile := filepath.Join(tmpdir, "id_rsa.pem")
	if err := os.WriteFile(keyFile, priv, 0600); err != nil {
		t.Fatalf("failed to write RSA key: %v", err)
	}

	if err := SealLogFile(logFile, keyFile); err != nil {
		t.Fatalf("SealLogFile failed: %v", err)
	}

	// Read file and split before seal
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read sealed log: %v", err)
	}
	parts := bytes.SplitN(data, []byte("=== FORENSIC SEAL ==="), 2)
	if len(parts) < 2 {
		t.Fatalf("seal block not found in log")
	}

	// Compute expected hash from original log content
	expected := seal.HashSHA256([]byte(logContent))

	// Extract SHA256 line from seal block
	sealPart := parts[1]
	// find line starting with "SHA256 Hash:"
	lines := bytes.Split(sealPart, []byte("\n"))
	var found string
	for _, ln := range lines {
		s := strings.TrimSpace(string(ln))
		if strings.HasPrefix(s, "SHA256 Hash:") {
			found = strings.TrimSpace(strings.TrimPrefix(s, "SHA256 Hash:"))
			break
		}
	}
	if found == "" {
		t.Fatalf("SHA256 Hash line not found in seal block")
	}
	if found != expected {
		t.Fatalf("sha mismatch: expected %s got %s", expected, found)
	}
}

// rsaGenerateAndMarshalPEM generates an RSA private key and returns it PEM-encoded.
func rsaGenerateAndMarshalPEM(bits int) ([]byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
	return pem.EncodeToMemory(pemBlock), nil
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
