# K3s E2E tests

These tests run N-Audit Sentinel end-to-end on a K3s cluster. The harness expects a running K3s cluster accessible via `kubectl`.

Prerequisites:
- `kubectl` configured to target the K3s cluster
- `docker` or `nerdctl` available for building images

See `tests/README.md` for global instructions.