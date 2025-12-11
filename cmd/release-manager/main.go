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
	// Ensure module dependencies are downloaded so we can build the full product
	fmt.Println("Downloading module dependencies...")
	if err := releasemgr.DownloadModules(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to download modules:", err)
		// continue — sometimes network issues; attempt builds anyway
	}

	// Build the product binaries
	bin1 := filepath.Join(*outdir, "n-audit-sentinel")
	if err := releasemgr.BuildTarget("./cmd/n-audit-sentinel", bin1, "linux", "amd64"); err != nil {
		fmt.Fprintln(os.Stderr, "build n-audit-sentinel failed:", err)
		os.Exit(2)
	}
	bin2 := filepath.Join(*outdir, "n-audit")
	if err := releasemgr.BuildTarget("./cmd/n-audit-cli", bin2, "linux", "amd64"); err != nil {
		fmt.Fprintln(os.Stderr, "build n-audit-cli failed:", err)
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
