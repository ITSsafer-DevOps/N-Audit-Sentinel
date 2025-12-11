package seal

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHashSHA256(t *testing.T) {
	data := []byte("abc")
	expected := sha256.Sum256(data)
	got := HashSHA256(data)
	if got != hex.EncodeToString(expected[:]) {
		t.Fatalf("HashSHA256 mismatch: got=%s want=%s", got, hex.EncodeToString(expected[:]))
	}
}

func TestEd25519SignVerify(t *testing.T) {
	pub, priv, err := GenerateEd25519KeyPair()
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}
	msg := []byte("hello world")
	sig := SignEd25519(priv, msg)
	if !VerifyEd25519(pub, msg, sig) {
		t.Fatalf("valid signature failed to verify")
	}
	// tamper
	if VerifyEd25519(pub, []byte("bad"), sig) {
		t.Fatalf("tampered message verified")
	}
}
