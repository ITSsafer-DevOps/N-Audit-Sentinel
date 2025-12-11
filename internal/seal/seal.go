package seal

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 returns the hex-encoded SHA-256 of the input bytes
func HashSHA256(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
