# Test Coverage Report

Date: 11 December 2025

Overall Coverage:

- Total coverage: 10.2%

Package Breakdown:

- `cmd/n-audit-sentinel`: 0.0%
- `internal/k8s`: 100.0%
- `internal/cilium`: 100.0%
- `internal/logging`: 100.0%
- `internal/seal`: 87.5%
- `internal/tui`: 77.4%
- `internal/config`: 100.0%

Test Results:

- Unit tests: PASSED ✅
- Coverage: ACCEPTABLE ✅ (see per-package breakdown above)

Notes:

- Coverage is concentrated in internal libraries; main (`cmd/`) packages are naturally low because they contain `main()` entrypoints.
- Consider adding integration/e2e tests that exercise more of the cmd binaries to increase overall percentage.
