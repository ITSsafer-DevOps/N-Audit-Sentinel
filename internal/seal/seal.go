package seal

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 returns the hex-encoded SHA-256 of the input bytes
func HashSHA256(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// GenerateEd25519KeyPair generates a public/private Ed25519 keypair.
func GenerateEd25519KeyPair() (pub, priv []byte, err error) {
	pubk, privk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pubk, privk, nil
}

// SignEd25519 signs the provided message using the Ed25519 private key.
func SignEd25519(priv []byte, msg []byte) []byte {
	return ed25519.Sign(priv, msg)
}

// VerifyEd25519 verifies the signature for msg using the public key.
func VerifyEd25519(pub []byte, msg []byte, sig []byte) bool {
	return ed25519.Verify(pub, msg, sig)
}
