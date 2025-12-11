# N-Audit Sentinel v1.0.0-Beta – Kali Linux Integration Request

## Project Overview

**Name:** N-Audit Sentinel
**Version:** v1.0.0-Beta
**Repository:** https://github.com/ITSsafer-DevOps/N-Audit-Sentinel
**License:** MIT
**Author:** Kristián Kašník / ITSsafer-DevOps

## Description

Enterprise-grade forensic auditing tool for Kubernetes clusters with cryptographic log sealing and cloud-native microsegmentation.

### Use Cases

- Penetration testing with scope enforcement
- Compliance audit trails with cryptographic sealing
- Red team operations with legally defensible logs
- Kubernetes security assessments
- Forensic investigation of cluster activity

## Compliance Status

✅ KALI LINUX READY (Audit passed)

## Requirements Met

- ✅ Source code: 100% English
- ✅ License: MIT (permissive)
- ✅ Documentation: Comprehensive (README, DEPLOYMENT, SECURITY, CONTRIBUTING)
- ✅ Tests: Unit, Integration, E2E (all passing)
- ✅ Build System: Professional (Makefile, CI/CD)
- ✅ Multi-platform: K3s, K8s, Talos, MicroShift, OpenShift
- ✅ Security Policy: Present (vulnerability disclosure, CVSS, SLA)
- ✅ No secrets: Verified
- ✅ Documentation: Diagrams, badges, professional styling

## Key Features

- Kubernetes environment discovery & enumeration
- Cryptographic log sealing (SHA256 + Ed25519 signatures)
- Dynamic Cilium Network Policy enforcement (Layer 3/7)
- Interactive TUI for audit session configuration
- Multi-environment support (K3s, K8s, Talos, MicroShift, OpenShift)
- Enterprise-ready Terraform deployment automation
- Full test coverage (80%+ core packages)

## Links

- Repository: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel
- Latest Release: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases/tag/v1.0.0-Beta
- Audit Report: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/blob/main/GITHUB_KALI_AUDIT_REPORT.md
- Deployment Guide: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/blob/main/DEPLOYMENT.md
- Security Policy: https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/blob/main/SECURITY.md

## Verification

Clone and test:

```bash
git clone https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel
make build
make test
```

Expected: All tests pass ✅

## Request

Please consider N-Audit Sentinel for inclusion in the Kali Linux toolkit.
