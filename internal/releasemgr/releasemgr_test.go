package releasemgr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeSHA256AndWrite(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rltest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	fpath := filepath.Join(tmpdir, "foo.txt")
	if err := ioutil.WriteFile(fpath, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	sum, err := ComputeSHA256(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if len(sum) != 64 {
		t.Fatalf("unexpected sha length: %d", len(sum))
	}
	out, err := WriteChecksumFile(fpath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}
