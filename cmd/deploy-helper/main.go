package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func PrepareStorageAndKeys(base string) error {
	storage := filepath.Join(base, "n-audit-data")
	signing := filepath.Join(storage, "signing")
	if err := os.MkdirAll(signing, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	keyPath := filepath.Join(signing, "id_ed25519")
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-N", "", "-f", keyPath, "-C", "n-audit-sentinel@localhost")
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ssh-keygen failed: %w", err)
	}
	if err := os.Chmod(keyPath, 0600); err != nil {
		return fmt.Errorf("chmod private key: %w", err)
	}
	if err := os.Chmod(keyPath+".pub", 0644); err != nil {
		return fmt.Errorf("chmod public key: %w", err)
	}
	fmt.Printf("storage prepared: %s\nkeys: %s\n", storage, signing)
	return nil
}

func main() {
	base := flag.String("prepare-storage", "", "Base path to prepare storage (e.g. /mnt)")
	flag.Parse()
	if *base == "" {
		fmt.Println("Usage: deploy-helper --prepare-storage /mnt")
		os.Exit(1)
	}
	if err := PrepareStorageAndKeys(*base); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}
