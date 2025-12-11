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
