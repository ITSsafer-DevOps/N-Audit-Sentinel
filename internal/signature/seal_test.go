package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"
	"testing"
)

func TestSealLogFile(t *testing.T) {
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
