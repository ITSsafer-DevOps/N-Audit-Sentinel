TITLE:
[New Tool Request] N-Audit Sentinel v1.0.0-Beta – Enterprise Kubernetes Forensic Auditing Tool

BODY:
## Project Information

**Name:** N-Audit Sentinel  
**Version:** v1.0.0-Beta  
**Homepage:** https://github.com/ITSsafer-DevOps/N-Audit-Sentinel  
**Repository:** https://github.com/ITSsafer-DevOps/N-Audit-Sentinel  
**Download:** https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases/tag/v1.0.0-Beta  
**License:** MIT  
**Author:** Kristián Kašník (ITSsafer-DevOps)  
**Contact:** itssafer@itssafer.org

## Description

N-Audit Sentinel is an enterprise-grade forensic auditing tool for Kubernetes clusters.

It secures audit log integrity through cryptographic SSH-based signatures (Ed25519) and enforces cloud-native microsegmentation via dynamic Cilium Network Policies.

### Primary Use Cases
- Penetration testing with scope enforcement  
- Compliance audit trails with cryptographic sealing  
- Red team operations with legally defensible logs  
- Kubernetes security assessments  
- Forensic investigation of cluster activity

## Key Features
- Kubernetes Environment Discovery & Enumeration  
- Cryptographic Log Sealing: SHA256 hashing + Ed25519 SSH signatures  
- Dynamic Cilium Network Policies (Layer 3 & Layer 7 enforcement)  
- Interactive TUI for session configuration  
- Persistent ANSI-stripped, timestamped session logs (hostPath or PVC)  
- Graceful shutdown with forensic seal block appended to logs  
- Multi-environment support: K3s, K8s, Talos Linux, MicroShift, OpenShift  
- Enterprise deployment automation via Terraform

## Technical Details
- Language: Go (module-based)  
- Minimum Go: 1.20+  
- Platforms: Linux amd64 (primary), Linux arm64 supported  
- Dependencies: kubectl, Docker/Podman for local builds; Cilium CNI v1.10+ for network policy features

## Testing & Quality
- Unit tests: 9+ (core packages) — all passing  
- Integration tests: 1+ (mock K8s clients) — passing  
- E2E framework scaffolded for K3s/K8s/Talos/MicroShift/OpenShift (ready)  
- Code quality: gofmt + go vet clean  
- CI: GitHub Actions (lint, unit tests, security scan)

## Packaging & Files
- Debian packaging scaffold: `contrib/debian/` (suitable for Kali integration)  
- Dockerfile and multi-stage builds included  
- Release tooling: `release-manager`, `backup-manager` utilities included in `cmd/`

## Artifacts & Checksums
- Release URL: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases/tag/v1.0.0-Beta

Checksums (SHA256):
- n-audit-sentinel-beta-bin.tar.gz: e48fc6794f5a46bae0a5340169dfd442cd60c2d14d03f3b1f5afb9de8ccb5d34  
- n-audit-sentinel-beta-source-code.tar.gz: f6444531bbff8411d07e506a3988966d8f5828546f9bb3c81e162d3554ebdeb9  
- Gold master backup: dcd206944e8caa46c1c25510f72fbb487f7ccfb9120684b5c3d8cad5803f1781

## Licensing & Compliance
- License: MIT (permissive)  
- All source code and documentation in English  
- Security policy present in `SECURITY.md`

## Notes for Kali Maintainers
- Build/test via `Makefile` (targets: `make build`, `make test`, `make release`)  
- Debian scaffold present: `contrib/debian/` — adjust control fields as needed for Kali packaging  
- Contact author for help: itssafer@itssafer.org

## Suggested Keywords / Tags
new-package, security, forensics, kubernetes, kali
