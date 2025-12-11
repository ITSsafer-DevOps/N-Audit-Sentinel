package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareStorageAndKeysWithKeygen_Success(t *testing.T) {
	dir, err := os.MkdirTemp("", "deploy-helper-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	fakeKeygen := func(keyPath string) error {
		// create private and public files with expected permissions
		priv := []byte("PRIVATE")
		pub := []byte("PUBLIC")
		if err := os.WriteFile(keyPath, priv, 0600); err != nil {
			return fmt.Errorf("write priv: %w", err)
		}
		if err := os.WriteFile(keyPath+".pub", pub, 0644); err != nil {
			return fmt.Errorf("write pub: %w", err)
		}
		return nil
	}

	if err := PrepareStorageAndKeysWithKeygen(dir, fakeKeygen); err != nil {
		t.Fatalf("PrepareStorageAndKeysWithKeygen failed: %v", err)
	}

	storage := filepath.Join(dir, "n-audit-data")
	signing := filepath.Join(storage, "signing")
	privPath := filepath.Join(signing, "id_ed25519")
	pubPath := privPath + ".pub"

	if _, err := os.Stat(privPath); err != nil {
		t.Fatalf("expected private key file: %v", err)
	}
	if _, err := os.Stat(pubPath); err != nil {
		t.Fatalf("expected public key file: %v", err)
	}

	// check perms
	st, err := os.Stat(privPath)
	if err != nil {
		t.Fatalf("stat priv: %v", err)
	}
	if st.Mode().Perm() != 0600 {
		t.Fatalf("private key perms expected 0600 got %o", st.Mode().Perm())
	}
	st, err = os.Stat(pubPath)
	if err != nil {
		t.Fatalf("stat pub: %v", err)
	}
	if st.Mode().Perm() != 0644 {
		t.Fatalf("public key perms expected 0644 got %o", st.Mode().Perm())
	}
}

func TestPrepareStorageAndKeysWithKeygen_KeygenError(t *testing.T) {
	dir, err := os.MkdirTemp("", "deploy-helper-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	fakeKeygen := func(keyPath string) error {
		return fmt.Errorf("injected failure")
	}

	if err := PrepareStorageAndKeysWithKeygen(dir, fakeKeygen); err == nil {
		t.Fatalf("expected error from PrepareStorageAndKeysWithKeygen when keygen fails")
	}
}
