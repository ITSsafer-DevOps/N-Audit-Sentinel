# Contributing to N-Audit Sentinel

This repository accepts contributions under the following guidelines:

- All code, documentation, and comments must be in English.
- Follow Go formatting: `make fmt`.
- Run linters and tests before submitting PRs: `make lint && make test`.

## Development

1. Clone repository and set up Go environment.
2. Run unit tests: `make test`.
3. Run e2e tests on local KinD cluster: `make test-e2e ENV=k8s`.

## Reporting security issues

Report security vulnerabilities privately to itssafer@itssafer.org.

