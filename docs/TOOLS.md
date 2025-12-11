# Tools: release-manager and backup-manager

This document describes the new Go utilities included in the repository.

## release-manager

- Location: `cmd/release-manager`
- Purpose: build local binaries (deterministic flags), package artifacts into `.tar.gz`, and generate `.sha256` checksums.

Usage example (Go):

```go
// Run release-manager from Go
package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "./cmd/release-manager", "--version", "v1.0.0-Beta", "--out", "out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## backup-manager

- Location: `cmd/backup-manager`
- Purpose: create a Gold Master source archive using `git archive` and write a checksum.

Usage example (Go):

```go
// Run backup-manager from Go
package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "./cmd/backup-manager", "--out", "gold-master-20251210T235959Z.tar.gz", "--ref", "HEAD")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
```

These utilities are intentionally minimal and testable. They are intended to be run locally and from CI.
