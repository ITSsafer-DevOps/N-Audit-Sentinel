package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	command := flag.String("cmd", "", "Command: verify-deps, release-bin, backup-source, list-artifacts")
	version := flag.String("version", "v1.0.0-Beta", "Version string")
	binDir := flag.String("bin-dir", "bin", "Binary directory")
	releaseDir := flag.String("release-dir", "releases", "Release directory")
	env := flag.String("env", "", "Environment for e2e tests (e.g., 'k3s')")
	flag.Parse()

	switch *command {
	case "verify-deps":
		verifyDeps()
	case "release-bin":
		releaseBinaries(*version, *binDir, *releaseDir)
	case "backup-source":
		backupSource(*version, *releaseDir)
	case "list-artifacts":
		listArtifacts(*releaseDir)
	case "check-tool":
		checkTool(flag.Arg(0))
	case "verify-e2e-env":
		verifyE2EEnv(*env)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func verifyDeps() {
	tools := map[string]bool{
		"go":      true,  // required
		"docker":  false, // optional
		"kubectl": false, // optional
	}

	for tool, required := range tools {
		path, err := exec.LookPath(tool)
		if err != nil {
			if required {
				log.Fatalf("ERROR: %s not found (required)", tool)
			}
			fmt.Printf("WARN: %s not found (optional)\n", tool)
		} else {
			fmt.Printf("OK: %s found at %s\n", tool, path)
		}
	}
}

func checkTool(tool string) {
	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("NOTFOUND\n")
		return
	}
	fmt.Printf("FOUND %s\n", path)
}

func releaseBinaries(version, binDir, releaseDir string) {
	// Create release directory
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		log.Fatalf("Failed to create release dir: %v", err)
	}

	// Tar and compress binary
	binFile := filepath.Join(binDir, "n-audit-sentinel")
	if _, err := os.Stat(binFile); err != nil {
		log.Fatalf("Binary not found: %s", binFile)
	}

	archiveName := fmt.Sprintf("%s/n-audit-sentinel-%s-bin.tar.gz", releaseDir, version)
	cmd := exec.Command("tar", "-czf", archiveName, "-C", binDir, "n-audit-sentinel")
	if err := cmd.Run(); err != nil {
		log.Fatalf("tar failed: %v", err)
	}

	// Generate SHA256 checksum
	file, err := os.Open(archiveName)
	if err != nil {
		log.Fatalf("Failed to open archive: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalf("Failed to hash: %v", err)
	}

	checksumFile := archiveName + ".sha256"
	checksumData := fmt.Sprintf("%x  %s\n", hash.Sum(nil), filepath.Base(archiveName))
	if err := os.WriteFile(checksumFile, []byte(checksumData), 0644); err != nil {
		log.Fatalf("Failed to write checksum: %v", err)
	}

	fmt.Printf("Released: %s\n", archiveName)
	fmt.Printf("Checksum: %s\n", checksumFile)

	// List artifacts
	listArtifacts(releaseDir)
}

func backupSource(version, releaseDir string) {
	// Create release directory
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		log.Fatalf("Failed to create release dir: %v", err)
	}

	archiveName := fmt.Sprintf("%s/n-audit-sentinel-%s-goldmaster.tar.gz", releaseDir, version)

	// git archive --format=tar --prefix=... HEAD | gzip > file
	gitCmd := exec.Command("git", "archive", "--format=tar",
		fmt.Sprintf("--prefix=n-audit-sentinel-%s-source/", version), "HEAD")
	gzipCmd := exec.Command("gzip", "-9")

	// Pipe git output to gzip
	var err error
	gzipCmd.Stdin, err = gitCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create pipe: %v", err)
	}

	outFile, err := os.Create(archiveName)
	if err != nil {
		log.Fatalf("Failed to create archive file: %v", err)
	}
	defer outFile.Close()

	gzipCmd.Stdout = outFile

	if err := gitCmd.Start(); err != nil {
		log.Fatalf("git archive failed: %v", err)
	}
	if err := gzipCmd.Run(); err != nil {
		log.Fatalf("gzip failed: %v", err)
	}
	if err := gitCmd.Wait(); err != nil {
		log.Fatalf("git archive wait failed: %v", err)
	}

	// Generate SHA256 checksum
	file, err := os.Open(archiveName)
	if err != nil {
		log.Fatalf("Failed to open archive: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalf("Failed to hash: %v", err)
	}

	checksumFile := archiveName + ".sha256"
	checksumData := fmt.Sprintf("%x  %s\n", hash.Sum(nil), filepath.Base(archiveName))
	if err := os.WriteFile(checksumFile, []byte(checksumData), 0644); err != nil {
		log.Fatalf("Failed to write checksum: %v", err)
	}

	fmt.Printf("Backed up: %s\n", archiveName)
	fmt.Printf("Checksum: %s\n", checksumFile)
}

func listArtifacts(releaseDir string) {
	entries, err := os.ReadDir(releaseDir)
	if err != nil {
		log.Fatalf("Failed to list release dir: %v", err)
	}

	fmt.Println("\nArtifacts in " + releaseDir + ":")
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".tar.gz") || strings.HasSuffix(entry.Name(), ".sha256")) {
			info, _ := entry.Info()
			size := formatSize(info.Size())
			fmt.Printf("  %s (%s)\n", entry.Name(), size)
		}
	}
}

func verifyE2EEnv(env string) {
	if env == "" {
		fmt.Println("No environment specified - skipping E2E test env checks")
		return
	}

	if env == "k3s" {
		fmt.Println("K3s environment selected")
		path, err := exec.LookPath("k3s")
		if err != nil {
			fmt.Printf("WARN: k3s not found - E2E tests may fail\n")
		} else {
			fmt.Printf("OK: k3s found at %s\n", path)
		}
	}
}

func formatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIdx := 0

	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}

	if unitIdx == 0 {
		return fmt.Sprintf("%d%s", int64(size), units[unitIdx])
	}
	return fmt.Sprintf("%.1f%s", size, units[unitIdx])
}
