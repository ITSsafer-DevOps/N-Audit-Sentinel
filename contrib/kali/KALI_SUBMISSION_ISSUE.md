TITLE: [New Tool Request] N-Audit Sentinel v1.0.0-Beta – Enterprise Kubernetes Forensic Auditing Tool

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

It secures audit log integrity through **cryptographic SSH-based signatures** (Ed25519) and enforces **cloud-native microsegmentation** via dynamic **Cilium Network Policies**.

### Primary Use Cases
- Penetration testing with scope enforcement
- Compliance audit trails with cryptographic sealing
- Red team operations with legally defensible logs
- Kubernetes security assessments
- Forensic investigation of cluster activity

## Key Features

✅ **Kubernetes Environment Discovery** – Automatically discovers and enumerates K8s API server and DNS
✅ **Cryptographic Log Sealing** – SHA256 hashing + Ed25519 SSH signatures for tamper-evident logs
✅ **Dynamic Cilium Network Policies** – Layer 3 (IP/CIDR) and Layer 7 (domain) scope enforcement
✅ **Interactive TUI** – User-friendly terminal UI for session configuration
✅ **Persistent Logging** – ANSI-stripped, timestamped session logs on hostPath or PVC
✅ **Graceful Shutdown** – Clean exit with forensic seal block appended to logs
✅ **Multi-Environment Support** – K3s, K8s, Talos Linux, MicroShift, OpenShift
✅ **Enterprise Deployment** – Terraform-based automation for repeatable deployments

## Technical Details

**Language:** Go 1.20+  
**Dependencies:** kubectl, Docker/Podman, Cilium CNI v1.10+  
**Platforms:** Linux amd64, arm64  
**Architecture:** Cloud-native, Kubernetes-first, container-based

## Project Maturity & Quality Assurance

### Testing
- ✅ 9+ unit tests (80%+ coverage of core packages)
- ✅ 1+ integration tests (mock Kubernetes clients)
- ✅ E2E framework scaffolded for all platforms
- ✅ 100% test pass rate

### Documentation
- ✅ README.md (comprehensive, enterprise-grade)
- ✅ DEPLOYMENT.md (multi-environment setup guides)
- ✅ VERIFICATION_GUIDE.md (testing & validation procedures)
- ✅ SECURITY.md (vulnerability disclosure policy, CVSS, SLA)
- ✅ CONTRIBUTING.md (development guidelines)
- ✅ TOOLS.md (internal utilities: release-manager, backup-manager)

### Build System & CI/CD
- ✅ Makefile with comprehensive targets (build, test, test-e2e, lint, security-scan, release)
- ✅ GitHub Actions CI workflow (lint, unit tests, security scan, release automation)
- ✅ Dockerfile (multi-stage, production-optimized)

### Security & Compliance
- ✅ Full security policy (CVSS v3.1, 72h acknowledgment SLA, 30d fix target)
- ✅ Debian packaging scaffold for Kali Linux integration
- ✅ MIT License (permissive, compatible with Kali)
- ✅ All source code in English

## Files Available

Repository structure:

(see repository tree in GitHub)

## Artifacts & Checksums
- Binary backup: e48fc6794f5a46bae0a5340169dfd442cd60c2d14d03f3b1f5afb9de8ccb5d34
- Source backup: f6444531bbff8411d07e506a3988966d8f5828546f9bb3c81e162d3554ebdeb9
- Gold master: dcd206944e8caa46c1c25510f72fbb487f7ccfb9120684b5c3d8cad5803f1781

