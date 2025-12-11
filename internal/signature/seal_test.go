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

    // create a sample log file
    logPath := filepath.Join(tmpDir, "session.log")
    content := []byte("line1\nline2\n")
    if err := ioutil.WriteFile(logPath, content, 0644); err != nil {
        t.Fatal(err)
    }

    // generate an RSA private key and write PEM file
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

    // verify appended seal
    data, err := ioutil.ReadFile(logPath)
    if err != nil {
        t.Fatal(err)
    }
    s := string(data)
    if !strings.Contains(s, "=== FORENSIC SEAL ===") {
        t.Fatalf("seal not found in log: %s", s)
    }

    // verify SHA256 matches original content
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
