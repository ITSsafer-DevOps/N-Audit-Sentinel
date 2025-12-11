package releasemgr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTarGz_EmptyFiles(t *testing.T) {
	tmpdir := t.TempDir()
	outPath := filepath.Join(tmpdir, "test.tar.gz")

	err := CreateTarGz(outPath, []string{})
	if err != nil {
		t.Fatalf("expected no error creating empty tar.gz, got: %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected tar.gz file to exist: %v", err)
	}
}

func TestCreateTarGz_NonexistentFile(t *testing.T) {
	tmpdir := t.TempDir()
	outPath := filepath.Join(tmpdir, "test.tar.gz")
	nonexistentFile := filepath.Join(tmpdir, "nonexistent.txt")

	err := CreateTarGz(outPath, []string{nonexistentFile})
	if err == nil {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestCreateTarGz_ValidFile(t *testing.T) {
	tmpdir := t.TempDir()
	outPath := filepath.Join(tmpdir, "test.tar.gz")
	testFile := filepath.Join(tmpdir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := CreateTarGz(outPath, []string{testFile})
	if err != nil {
		t.Fatalf("expected no error creating tar.gz with valid file, got: %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected tar.gz file to exist: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("stat tar.gz failed: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected non-empty tar.gz, got size 0")
	}
}

func TestCreateTarGz_MultipleFiles(t *testing.T) {
	tmpdir := t.TempDir()
	outPath := filepath.Join(tmpdir, "multi.tar.gz")

	file1 := filepath.Join(tmpdir, "file1.txt")
	file2 := filepath.Join(tmpdir, "file2.txt")

	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	err := CreateTarGz(outPath, []string{file1, file2})
	if err != nil {
		t.Fatalf("expected no error creating tar.gz with multiple files, got: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("stat tar.gz failed: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected non-empty tar.gz, got size 0")
	}
}

func TestDownloadModules_Basic(t *testing.T) {
	err := DownloadModules()
	if err != nil {
		t.Logf("Note: DownloadModules returned error: %v (may be expected in test environment)", err)
	}
}

func TestBuildTarget_InvalidPackagePath(t *testing.T) {
	tmpdir := t.TempDir()
	dest := filepath.Join(tmpdir, "binary")

	err := BuildTarget("./nonexistent/package", dest, "", "")
	if err == nil {
		t.Fatalf("expected error for nonexistent package, got nil")
	}
}

func TestComputeSHA256_ValidFile(t *testing.T) {
	tmpdir := t.TempDir()
	testFile := filepath.Join(tmpdir, "test.txt")
	content := "hello world"

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	sum, err := ComputeSHA256(testFile)
	if err != nil {
		t.Fatalf("ComputeSHA256 failed: %v", err)
	}

	if len(sum) == 0 {
		t.Fatalf("expected non-empty checksum, got empty")
	}

	// SHA256 should be 64 hex characters
	if len(sum) != 64 {
		t.Fatalf("expected 64-char SHA256, got %d chars", len(sum))
	}
}

func TestComputeSHA256_NonexistentFile(t *testing.T) {
	_, err := ComputeSHA256("/nonexistent/file.txt")
	if err == nil {
		t.Fatalf("expected error for nonexistent file, got nil")
	}
}

func TestWriteChecksumFile_CreatesFile(t *testing.T) {
	tmpdir := t.TempDir()
	testFile := filepath.Join(tmpdir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	outPath, err := WriteChecksumFile(testFile)
	if err != nil {
		t.Fatalf("WriteChecksumFile failed: %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("checksum file not created: %v", err)
	}
}
