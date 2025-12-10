package backupmgr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateSourceArchive uses `git archive` to create a tar.gz archive of the repository at ref
func CreateSourceArchive(ref, outPath string) error {
	// ensure output dir exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("git", "archive", "--format=tar.gz", "-o", outPath, ref)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ComputeSHA256 computes sha256 for a file
func ComputeSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// WriteChecksum writes checksum to path.sha256
func WriteChecksum(path string) (string, error) {
	s, err := ComputeSHA256(path)
	if err != nil {
		return "", err
	}
	out := path + ".sha256"
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "%s  %s\n", s, filepath.Base(path)); err != nil {
		return "", err
	}
	return out, nil
}
