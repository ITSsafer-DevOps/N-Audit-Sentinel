package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/releasemgr"
)

func main() {
	version := flag.String("version", "", "Release version (e.g. v1.0.0)")
	outdir := flag.String("out", "out", "Output directory for build artifacts")
	flag.Parse()
	if *version == "" {
		fmt.Fprintln(os.Stderr, "version is required")
		os.Exit(2)
	}
	if err := os.MkdirAll(*outdir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "failed to create out dir:", err)
		os.Exit(2)
	}
	// Build binaries (limited to newly scaffolded CLIs to avoid external dep resolution in this environment)
	bin1 := filepath.Join(*outdir, "release-manager")
	if err := releasemgr.BuildTarget("./cmd/release-manager", bin1, "", ""); err != nil {
		fmt.Fprintln(os.Stderr, "build failed:", err)
		os.Exit(2)
	}
	bin2 := filepath.Join(*outdir, "backup-manager")
	if err := releasemgr.BuildTarget("./cmd/backup-manager", bin2, "", ""); err != nil {
		fmt.Fprintln(os.Stderr, "build failed:", err)
		os.Exit(2)
	}

	// Package
	tarPath := fmt.Sprintf("n-audit-sentinel-%s-linux-amd64.tar.gz", *version)
	files := []string{bin1, bin2}
	if err := releasemgr.CreateTarGz(tarPath, files); err != nil {
		fmt.Fprintln(os.Stderr, "package failed:", err)
		os.Exit(2)
	}
	if _, err := releasemgr.WriteChecksumFile(tarPath); err != nil {
		fmt.Fprintln(os.Stderr, "checksum failed:", err)
		os.Exit(2)
	}
	fmt.Println("Created:", tarPath)
}
