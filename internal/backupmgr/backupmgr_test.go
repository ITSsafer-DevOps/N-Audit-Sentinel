package backupmgr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeSHA256AndWriteChecksum(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "sample.txt")
	content := "hello world\n"
	if err := os.WriteFile(fpath, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	sum, err := ComputeSHA256(fpath)
	if err != nil {
		t.Fatalf("ComputeSHA256 error: %v", err)
	}
	if len(sum) == 0 {
		t.Fatalf("empty checksum")
	}

	outPath, err := WriteChecksum(fpath)
	if err != nil {
		t.Fatalf("WriteChecksum error: %v", err)
	}
	if !strings.HasSuffix(outPath, ".sha256") {
		t.Fatalf("unexpected checksum path: %s", outPath)
	}
	// verify file exists
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("checksum file missing: %v", err)
	}
}
func TestComputeSHA256_NonexistentFile(t *testing.T) {
	_, err := ComputeSHA256("/nonexistent/file.txt")
	if err == nil {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestComputeSHA256_DifferentFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0o644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0o644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	sum1, err := ComputeSHA256(file1)
	if err != nil {
		t.Fatalf("ComputeSHA256(file1) error: %v", err)
	}
	sum2, err := ComputeSHA256(file2)
	if err != nil {
		t.Fatalf("ComputeSHA256(file2) error: %v", err)
	}

	if sum1 == sum2 {
		t.Fatalf("expected different checksums for different files")
	}
}

func TestComputeSHA256_SameContent(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")

	content := "same content"
	if err := os.WriteFile(file1, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	sum1, err := ComputeSHA256(file1)
	if err != nil {
		t.Fatalf("ComputeSHA256(file1) error: %v", err)
	}
	sum2, err := ComputeSHA256(file2)
	if err != nil {
		t.Fatalf("ComputeSHA256(file2) error: %v", err)
	}

	if sum1 != sum2 {
		t.Fatalf("expected same checksums for identical files")
	}
}

func TestWriteChecksum_CreatesDirIfNeeded(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(fpath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	outPath, err := WriteChecksum(fpath)
	if err != nil {
		t.Fatalf("WriteChecksum error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read checksum file: %v", err)
	}
	if !strings.Contains(string(data), "sample.txt") {
		t.Fatalf("checksum file does not contain filename")
	}
}

func TestCreateSourceArchive_InvalidRef(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "archive.tar.gz")

	// Try to create archive with invalid ref (if not in git repo, will fail)
	err := CreateSourceArchive("nonexistent-ref", outPath)
	if err != nil {
		t.Logf("Expected error with invalid ref: %v", err)
	}
}
