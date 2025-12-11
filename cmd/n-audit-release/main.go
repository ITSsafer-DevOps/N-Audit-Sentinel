// N-Audit Sentinel Release Tool
// Developer: Kristian Kasnik
// Company: ITSsafer-DevOps
// License: MIT License
// Automates building and packaging binaries into a versioned tarball with SHA256 signature.
package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// cmdRunner abstracts exec.Command for testability.
type cmdRunner interface {
	Run(name string, args ...string) error
}

// realCmdRunner implements cmdRunner using exec.Command.
type realCmdRunner struct{}

func (r *realCmdRunner) Run(name string, args ...string) error {
	return runCmd(name, args...)
}

// BuildBinariesWithRunner builds binaries using the provided cmdRunner.
// This allows tests to inject a fake command runner without actual go build calls.
func BuildBinariesWithRunner(runner cmdRunner, buildDir string) (sentinelBin, cliBin string, err error) {
	sentinelBin = filepath.Join(buildDir, "n-audit-sentinel")
	cliBin = filepath.Join(buildDir, "n-audit")

	fmt.Println("[release] Building binaries for linux/amd64...")
	if err := runner.Run("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", sentinelBin, "./cmd/n-audit-sentinel"); err != nil {
		return "", "", fmt.Errorf("failed to build sentinel: %w", err)
	}
	if err := runner.Run("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", cliBin, "./cmd/n-audit-cli"); err != nil {
		return "", "", fmt.Errorf("failed to build cli: %w", err)
	}

	return sentinelBin, cliBin, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: n-audit-release <version>\n")
		os.Exit(1)
	}
	version := os.Args[1]
	if version == "" {
		fmt.Fprintln(os.Stderr, "empty version string")
		os.Exit(1)
	}

	// Build output paths

	buildDir, err := os.MkdirTemp("", "n-audit-release-build-")
	must(err)
	defer os.RemoveAll(buildDir)

	// Cross-compile binaries for linux/amd64
	realRunner := &realCmdRunner{}
	sentinelBin, cliBin, err := BuildBinariesWithRunner(realRunner, buildDir)
	if err != nil {
		must(err)
	}

	// Tarball name
	tarName := fmt.Sprintf("n-audit-sentinel-%s-linux-amd64.tar.gz", version)
	tarFile, err := os.Create(tarName)
	must(err)
	defer tarFile.Close()

	gz := gzip.NewWriter(tarFile)
	tw := tar.NewWriter(gz)

	// Add binaries to tar
	addFile := func(srcPath, name string) {
		info, err := os.Stat(srcPath)
		must(err)
		hdr, err := tar.FileInfoHeader(info, "")
		must(err)
		hdr.Name = name
		must(tw.WriteHeader(hdr))
		f, err := os.Open(srcPath)
		must(err)
		defer f.Close()
		_, err = io.Copy(tw, f)
		must(err)
	}

	fmt.Println("[release] Packaging binaries...")
	addFile(sentinelBin, "n-audit-sentinel")
	addFile(cliBin, "n-audit")
	must(tw.Close())
	must(gz.Close())

	// Compute SHA-256 of tarball
	fmt.Println("[release] Computing SHA-256 signature...")
	tarData, err := os.ReadFile(tarName)
	must(err)
	sum := sha256.Sum256(tarData)
	sigHex := hex.EncodeToString(sum[:])
	sigFile := tarName + ".sha256"
	must(os.WriteFile(sigFile, []byte(sigHex+"\n"), 0644))
	fmt.Printf("[release] Done. Artifacts: %s, %s\n", tarName, sigFile)
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func must(err error) {
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			fmt.Fprintf(os.Stderr, "[release] command failed: %v\n", ee)
		} else {
			fmt.Fprintf(os.Stderr, "[release] error: %v\n", err)
		}
		os.Exit(1)
	}
}
