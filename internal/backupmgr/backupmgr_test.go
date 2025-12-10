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
