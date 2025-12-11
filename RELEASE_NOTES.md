# Release Notes — N-Audit Sentinel v1.0.0-Beta

**Release Date:** 2025-12-11
**Build Status:** ✅ All tests passing (49.5% coverage)
**Artifacts Location:** `releases/`

## Build Artifacts

- **Binary Bundle:** `n-audit-sentinel-v1.0.0-Beta-bin.tar.gz` (30M)
  - SHA256: `n-audit-sentinel-v1.0.0-Beta-bin.tar.gz.sha256`
  - Contents: Pre-compiled `n-audit-sentinel` binary (linux/amd64)

- **Gold Master Source:** `n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz` (79M)
  - SHA256: `n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz.sha256`
  - Contents: Complete source tree at commit HEAD, deterministic (git archive)

## Build & Test Results

```
✅ Build: All binaries compiled successfully (go build)
✅ Tests: 100% pass rate (go test ./...)
✅ Coverage: 49.5% overall (improved from 46.8% baseline)
✅ Format: All code formatted (gofmt)
✅ Artifacts: Release tarballs + SHA256 checksums verified
```

## Key Features

- **Forensic Session Sealing:** Log file signing with Ed25519 or RSA keys
- **Network Policy Enforcement:** 3-zone Cilium CNI model (Infrastructure, Maintenance, Target)
- **Interactive TUI:** Pentester identification + scope definition
- **Protected Session Recording:** Bash session with output capture and sealing
- **Kubernetes-Native:** In-cluster deployment via StatefulSet/Pod

## Documentation

- **README.md:** Architecture overview + deployment guide (4 Mermaid diagrams)
- **DEPLOYMENT.md:** Step-by-step Kubernetes/Terraform setup
- **SECURITY.md:** Threat model + operational security (Mermaid diagram)
- **TESTING_AND_VERIFICATION.md:** Test suite guide + CI/CD patterns
- **docs/ARCHITECTURE_DIAGRAMS.md:** 4 advanced architecture diagrams
- **COMPREHENSIVE_ENHANCEMENT_REPORT.md:** Full project status + coverage metrics

## Go Package Coverage

| Package | Coverage |
|---------|----------|
| internal/cilium | 100.0% |
| internal/config | 100.0% |
| internal/k8s | 100.0% |
| internal/logging | 100.0% |
| internal/seal | 87.5% |
| internal/signature | 87.0% |
| internal/discovery | 86.4% |
| internal/recorder | 85.4% |
| internal/logger | 82.6% |
| internal/backupmgr | 80.0% |
| internal/releasemgr | 78.3% |
| internal/tui | 77.4% |
| internal/validation | 76.9% |
| internal/policy | 72.6% |
| cmd/n-audit-cli | 42.9% |
| cmd/deploy-helper | 27.6% |
| cmd/n-audit-release | 10.8% |
| cmd/n-audit-sentinel | 4.3% |
| **Total** | **49.5%** |

## Installation

### From Binary Bundle

```bash
tar -xzf n-audit-sentinel-v1.0.0-Beta-bin.tar.gz
cp n-audit-sentinel /usr/local/bin/
chmod +x /usr/local/bin/n-audit-sentinel
```

### From Source

```bash
tar -xzf n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz
cd n-audit-sentinel-v1.0.0-Beta-source
make build
./bin/n-audit-sentinel
```

## Verification

### Verify Artifact Integrity

```bash
sha256sum -c n-audit-sentinel-v1.0.0-Beta-bin.tar.gz.sha256
sha256sum -c n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz.sha256
```

### Verify Build

```bash
cd n-audit-sentinel-v1.0.0-Beta-source
make clean
make test      # Run full test suite
make release   # Create release artifacts
```

## Known Limitations

- E2E tests are skipped by default (set `RUN_E2E=true` to run)
- Command package coverage remains low (refactored for DI testability)
- Requires Go 1.25+ and Linux/amd64 platform

## Future Work

- Increase command package coverage through integration tests
- Add Windows/Darwin cross-compile targets
- Expand policy zones for multi-tenant isolation
- Machine-readable audit log export (JSON/JSONL)

---

**Git Commit:** `d970df4` (main branch)
**Release Channel:** Production — v1.0.0-Beta
**Support:** See CONTRIBUTING.md for development guidelines
