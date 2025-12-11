package signature

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/seal"
)

// SealLogFile computes SHA-256 of the log file and appends a seal block containing
// the hex-encoded hash and an SSH signature generated from the provided private key.
//
// Parameters:
// - logPath: absolute path to the session log file to seal.
// - privateKeyPath: absolute path to an SSH private key (e.g., Ed25519).
//
// Returns:
// - error: non-nil if reading, hashing, key parsing, signing, or appending fails.
func SealLogFile(logPath, privateKeyPath string) error {
	logData, err := os.ReadFile(logPath)
	if err != nil {
		return fmt.Errorf("read log file: %w", err)
	}

	// Use the canonical internal/seal package for hashing
	hashHex := seal.HashSHA256(logData)

	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	// Decode the hex hash back to bytes for signing
	logHash := []byte(hashHex)
	sig, err := signer.Sign(nil, logHash)
	if err != nil {
		return fmt.Errorf("sign hash: %w", err)
	}
	b64Sig := base64.StdEncoding.EncodeToString(sig.Blob)

	sealBlock := fmt.Sprintf("\n\n=== FORENSIC SEAL ===\nSHA256 Hash: %s\nSSH Signature (Base64): %s\n=====================\n", hashHex, b64Sig)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log for append: %w", err)
	}
	defer f.Close()

	if _, err := io.WriteString(f, sealBlock); err != nil {
		return fmt.Errorf("append seal: %w", err)
	}
	return nil
}

// signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
// go func() {
//     sig := <-sigCh
//     fmt.Fprintf(os.Stdout, "\n[N-Audit] Received signal %s. Initiating shutdown...\n", sig)
//     cancel()
// }()
