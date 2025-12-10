package releasemgr

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildTarget builds a go package at pkgPath into dest binary path
func BuildTarget(pkgPath, dest string, goos, goarch string) error {
	env := os.Environ()
	if goos != "" {
		env = append(env, "GOOS="+goos)
	}
	if goarch != "" {
		env = append(env, "GOARCH="+goarch)
	}
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", dest, pkgPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreateTarGz creates a tar.gz archive at outPath containing the files in files (paths are relative or absolute)
func CreateTarGz(outPath string, files []string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, fn := range files {
		fi, err := os.Stat(fn)
		if err != nil {
			return err
		}
		var file *os.File
		if fi.IsDir() {
			// skip directories for now
			continue
		}
		file, err = os.Open(fn)
		if err != nil {
			return err
		}
		hdr := &tar.Header{
			Name: filepath.Base(fn),
			Mode: int64(fi.Mode().Perm()),
			Size: fi.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			file.Close()
			return err
		}
		if _, err := io.Copy(tw, file); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

// ComputeSHA256 computes the SHA256 checksum of a file at path
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

// WriteChecksumFile writes the sha256 sum to a .sha256 file alongside archive
func WriteChecksumFile(archivePath string) (string, error) {
	sum, err := ComputeSHA256(archivePath)
	if err != nil {
		return "", err
	}
	out := archivePath + ".sha256"
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "%s  %s\n", sum, filepath.Base(archivePath)); err != nil {
		return "", err
	}
	return out, nil
}
