# Tests

This directory contains testing scaffolding for the N-Audit Sentinel project.

Structure:

- `tests/unit/` - unit tests (fast, no cluster required)
- `tests/integration/` - integration tests (mocked or lightweight cluster)
- `tests/e2e/` - end-to-end tests organized by environment (k3s, k8s, talos, microshift, openshift)

See the top-level `Makefile` for commands to run tests.
