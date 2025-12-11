# Tools and Utilities

This document lists the small command-line utilities included in the repository and their purpose. These tools are implemented under `cmd/` and are intended to be used by maintainers, CI, and release engineers.

Discovered tools (from `cmd/`):

- `n-audit-sentinel` (cmd/n-audit-sentinel)
  - Main PID 1 runtime for the project. Runs inside the audit container, presents the interactive TUI, discovers Kubernetes environment, applies Cilium policies, records PTY sessions, and performs log sealing on exit.

- `n-audit` (cmd/n-audit-cli)
  - Lightweight CLI that signals the running PID 1 process (SIGUSR1) to trigger a controlled forensic seal and exit sequence.

- `n-audit-release` (cmd/n-audit-release)
  - Release helper that builds product binaries (linux/amd64) and packages them into a versioned tarball with SHA256 checksum. Intended for simple release automation.

- `release-manager` (cmd/release-manager)
  - Go-based release manager that builds targets and produces deterministic tarballs; wraps internal `releasemgr` helpers for CI usage.

- `backup-manager` (cmd/backup-manager)
  - Creates a Gold Master source archive (via `git archive`) and writes a checksum file for reproducible backups.

Usage examples

Build a single tool (Go example):

```go
// Build a single tool programmatically
package main

import (
  "log"
  "os/exec"
)

func main() {
  if err := exec.Command("go", "build", "-o", "out/n-audit-sentinel", "./cmd/n-audit-sentinel").Run(); err != nil {
    log.Fatal(err)
  }
}
```

Create a release tarball (Go example):

```go
// Run release-manager from Go
package main

import (
  "log"
  "os/exec"
)

func main() {
  if err := exec.Command("go", "run", "./cmd/release-manager", "--version", "v1.0.0-Beta", "--out", "out").Run(); err != nil {
    log.Fatal(err)
  }
}
```

Notes

- These utilities are small and focused; they are intended for maintainers and CI. Binaries should not be committed to the repository (they belong in Releases or packaging artifacts).
- See `cmd/*/README.md` for additional usage details where present.
