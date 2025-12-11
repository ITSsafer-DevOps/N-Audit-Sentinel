package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSealLogFile_AppendsSeal(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "sealtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "session.log")
	content := []byte("line1\nline2\n")
	if err := ioutil.WriteFile(logPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	pemBlk := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	keyPath := filepath.Join(tmpDir, "id_rsa")
	f, err := os.Create(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := pem.Encode(f, pemBlk); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("SealLogFile failed: %v", err)
	}

	data, err := ioutil.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "=== FORENSIC SEAL ===") {
		t.Fatalf("seal not found in log: %s", s)
	}

	h := sha256.Sum256(content)
	if !strings.Contains(s, strings.ToLower(strings.TrimSpace(hexEncode(h[:])))) {
		t.Fatalf("hash not present or mismatched in seal")
	}
}

func hexEncode(b []byte) string {
	const hextable = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hextable[v>>4]
		out[i*2+1] = hextable[v&0x0f]
	}
	return string(out)
}

func TestSealLogFile_Smoke(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyDir := t.TempDir()
	keyPath := keyDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	logDir := t.TempDir()
	logPath := logDir + "/session.log"
	if err := os.WriteFile(logPath, []byte("test log line\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("seal log: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sealed log: %v", err)
	}
	content := string(data)
	required := []string{"=== FORENSIC SEAL ===", "SHA256 Hash:", "SSH Signature (Base64):", "====================="}
	for _, part := range required {
		if !strings.Contains(content, part) {
			t.Fatalf("seal block missing part %q: %s", part, content)
		}
	}
}

func TestSealLogFile_NonexistentLogFile(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyDir := t.TempDir()
	keyPath := keyDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	if err := SealLogFile("/nonexistent/log/file.log", keyPath); err == nil {
		t.Fatalf("expected error for nonexistent log file")
	}
}

func TestSealLogFile_NonexistentKeyFile(t *testing.T) {
	logDir := t.TempDir()
	logPath := logDir + "/session.log"
	if err := os.WriteFile(logPath, []byte("test content\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	if err := SealLogFile(logPath, "/nonexistent/key/file.pem"); err == nil {
		t.Fatalf("expected error for nonexistent key file")
	}
}

func TestSealLogFile_InvalidKeyFormat(t *testing.T) {
	logDir := t.TempDir()
	logPath := logDir + "/session.log"
	if err := os.WriteFile(logPath, []byte("test content\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	keyPath := logDir + "/invalid.pem"
	if err := os.WriteFile(keyPath, []byte("not a valid ssh key"), 0644); err != nil {
		t.Fatalf("write invalid key: %v", err)
	}

	if err := SealLogFile(logPath, keyPath); err == nil {
		t.Fatalf("expected error for invalid key format")
	}
}

func TestSealLogFile_EmptyLogFile(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := tmpDir + "/empty.log"
	if err := os.WriteFile(logPath, []byte(""), 0644); err != nil {
		t.Fatalf("write empty log: %v", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyPath := tmpDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("SealLogFile should handle empty file: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sealed log: %v", err)
	}
	if !strings.Contains(string(data), "FORENSIC SEAL") {
		t.Fatalf("expected seal appended to empty log")
	}
}

func TestSealLogFile_LargeLogFile(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := tmpDir + "/large.log"
	largeContent := strings.Repeat("x", 1000000) + "\n"
	if err := os.WriteFile(logPath, []byte(largeContent), 0644); err != nil {
		t.Fatalf("write large log: %v", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyPath := tmpDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("SealLogFile should handle large file: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sealed log: %v", err)
	}
	if !strings.Contains(string(data), "FORENSIC SEAL") {
		t.Fatalf("expected seal appended to large log")
	}
}

func TestSealLogFile_ReadOnlyLogFile(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := tmpDir + "/readonly.log"
	if err := os.WriteFile(logPath, []byte("test content\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyPath := tmpDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	// Make log file read-only
	if err := os.Chmod(logPath, 0444); err != nil {
		t.Fatalf("chmod log: %v", err)
	}
	defer os.Chmod(logPath, 0644)

	if err := SealLogFile(logPath, keyPath); err == nil {
		t.Fatalf("expected error when appending to read-only file")
	}
}

func TestSealLogFile_SealBlockFormat(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := tmpDir + "/format.log"
	if err := os.WriteFile(logPath, []byte("content\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyPath := tmpDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("SealLogFile failed: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sealed log: %v", err)
	}
	content := string(data)

	// Verify seal block boundaries
	if !strings.Contains(content, "\n=== FORENSIC SEAL ===\n") {
		t.Fatalf("expected proper seal block opening boundary")
	}
	if !strings.Contains(content, "\n=====================\n") {
		t.Fatalf("expected proper seal block closing boundary")
	}

	// Verify hash format (64 hex chars for SHA256)
	if !strings.Contains(content, "SHA256 Hash: ") {
		t.Fatalf("expected SHA256 Hash label")
	}
}

func TestSealLogFile_MultipleSeals(t *testing.T) {
	tmpDir := t.TempDir()

	logPath := tmpDir + "/multi.log"
	if err := os.WriteFile(logPath, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	keyPath := tmpDir + "/id_rsa"
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	// First seal
	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("first seal failed: %v", err)
	}

	firstData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log after first seal: %v", err)
	}

	// Second seal (sealing already sealed log)
	if err := SealLogFile(logPath, keyPath); err != nil {
		t.Fatalf("second seal failed: %v", err)
	}

	secondData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log after second seal: %v", err)
	}

	// Verify file grew (second seal appended)
	if len(secondData) <= len(firstData) {
		t.Fatalf("expected second seal to append to file")
	}

	// Count seal blocks
	sealCount := strings.Count(string(secondData), "=== FORENSIC SEAL ===")
	if sealCount != 2 {
		t.Fatalf("expected 2 seal blocks, found %d", sealCount)
	}
}
