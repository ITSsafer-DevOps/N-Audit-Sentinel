# Upstream Kubernetes E2E tests

These tests run N-Audit Sentinel end-to-end on an upstream Kubernetes cluster (KinD/GKE/EKS/AKS). The harness expects cluster access via `kubectl`.

Prerequisites:
- `kubectl` configured to target the cluster
- `kind` (optional) for local testing

See `tests/README.md` for global instructions.