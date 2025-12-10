# Tools: release-manager and backup-manager

This document describes the new Go utilities included in the repository.

## release-manager

- Location: `cmd/release-manager`
- Purpose: build local binaries (deterministic flags), package artifacts into `.tar.gz`, and generate `.sha256` checksums.

Usage example:

```bash
# Build the local release artifacts into current directory
go run ./cmd/release-manager --version v1.0.0-Beta --out out

# The tool will create: n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz and the .sha256 file
```

## backup-manager

- Location: `cmd/backup-manager`
- Purpose: create a Gold Master source archive using `git archive` and write a checksum.

Usage example:

```bash
go run ./cmd/backup-manager --out gold-master-20251210T235959Z.tar.gz --ref HEAD
```

These utilities are intentionally minimal and testable. They are intended to be run locally and from CI.
