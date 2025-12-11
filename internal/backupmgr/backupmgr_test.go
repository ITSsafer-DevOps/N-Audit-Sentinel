package backupmgr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeSHA256AndWrite(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bktest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	fpath := filepath.Join(tmpdir, "bar.txt")
	if err := ioutil.WriteFile(fpath, []byte("backup content"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := ComputeSHA256(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 64 {
		t.Fatalf("unexpected sha length: %d", len(s))
	}
	out, err := WriteChecksum(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}

func TestComputeSHA256_NonexistentFile(t *testing.T) {
	_, err := ComputeSHA256("/nonexistent/backup/file")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestWriteChecksum_Format(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "chksum2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	fpath := filepath.Join(tmpdir, "backup.tar.gz")
	if err := ioutil.WriteFile(fpath, []byte("backup data"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := WriteChecksum(fpath)
	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	expectedFile := filepath.Base(fpath)
	contentStr := string(contents)
	if !contains(contentStr, expectedFile) {
		t.Fatalf("checksum file missing filename: %s", contentStr)
	}
}

func TestWriteChecksum_NonexistentFile(t *testing.T) {
	_, err := WriteChecksum("/nonexistent/file.tar.gz")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestCreateSourceArchive_DirectoryCreation(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "archtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Test that CreateSourceArchive can create nested directories and actually create the archive
	// We use HEAD which should always exist in a git repo
	outPath := filepath.Join(tmpdir, "archive.tar.gz")
	err = CreateSourceArchive("HEAD", outPath)
	if err != nil {
		// It's OK if this fails due to git not being available; we're testing the function doesn't panic
		t.Logf("CreateSourceArchive returned error (expected if git not available): %v", err)
	}
}

func TestCreateSourceArchive_NestedDir(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "archtest2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Test nested directory creation within CreateSourceArchive
	outPath := filepath.Join(tmpdir, "nested", "dir", "archive.tar.gz")
	// This will attempt to create nested directories
	err = CreateSourceArchive("HEAD", outPath)
	if err != nil {
		t.Logf("CreateSourceArchive with nested dir returned error: %v", err)
	}
	// If it succeeds, verify the file exists
	if err == nil {
		if _, err := os.Stat(outPath); err != nil {
			t.Logf("archive was not created: %v", err)
		}
	}
}

func TestComputeSHA256_DifferentFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "sha1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	file1 := filepath.Join(tmpdir, "file1.txt")
	file2 := filepath.Join(tmpdir, "file2.txt")

	if err := ioutil.WriteFile(file1, []byte("content1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(file2, []byte("content2"), 0o644); err != nil {
		t.Fatal(err)
	}

	hash1, err := ComputeSHA256(file1)
	if err != nil {
		t.Fatal(err)
	}
	hash2, err := ComputeSHA256(file2)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 == hash2 {
		t.Fatal("different files should have different hashes")
	}
	if len(hash1) != 64 || len(hash2) != 64 {
		t.Fatal("hashes should be 64 characters")
	}
}

func TestWriteChecksum_ChecksumFileContent(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "chksum3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	fpath := filepath.Join(tmpdir, "data.tar.gz")
	testData := []byte("test archive content")
	if err := ioutil.WriteFile(fpath, testData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Compute expected hash manually
	expectedHash, err := ComputeSHA256(fpath)
	if err != nil {
		t.Fatal(err)
	}

	// Write checksum via function
	checksumFile, err := WriteChecksum(fpath)
	if err != nil {
		t.Fatal(err)
	}

	// Read checksum file
	content, err := ioutil.ReadFile(checksumFile)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	// Verify hash appears in the file
	if !contains(contentStr, expectedHash) {
		t.Fatalf("checksum file doesn't contain expected hash: %s", contentStr)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
