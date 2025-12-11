# Contributing to N-Audit Sentinel

This repository accepts contributions under the following guidelines:

- All code, documentation, and comments must be in English.
- Follow Go formatting: `make fmt`.
- Run linters and tests before submitting PRs: `make lint && make test`.

## Development Workflow

### 1. Setup
Clone the repository and ensure Go 1.25+ is installed:
```bash
git clone https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel
go mod tidy
```

### 2. Running Tests Locally
Before submitting a pull request, run the full test suite:

```bash
# Format code
make fmt

# Lint code (runs `go vet` and basic static analysis)
make lint

# Run all unit and integration tests
make test

# Run end-to-end tests (requires K3s environment or set RUN_E2E=true to skip external deps)
make test-e2e ENV=k3s

# Security vulnerability scanning (if govulncheck is installed)
make security-scan

# Verify dependencies (go, docker, kubectl availability)
make verify-deps
```

### 3. Test Coverage
# Contributing to N-Audit Sentinel

This repository accepts contributions under the following guidelines:

- All code, documentation, and comments must be in English.
- Follow Go formatting: `make fmt`.
- Run linters and tests before submitting PRs: `make lint && make test`.

## Development Workflow

### 1. Setup
Clone the repository and ensure Go 1.25+ is installed:

```bash
git clone https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel
go mod tidy
```

### 2. Running Tests Locally
Before submitting a pull request, run the full test suite:

```bash
# Format code
make fmt

# Lint code (runs `go vet` and basic static analysis)
make lint

# Run all unit and integration tests
make test

# Run end-to-end tests (requires K3s environment or set RUN_E2E=true to skip external deps)
make test-e2e ENV=k3s

# Security vulnerability scanning (if govulncheck is installed)
make security-scan

# Verify dependencies (go, docker, kubectl availability)
make verify-deps
```

### 3. Test Coverage
Unit tests are located in:

- `tests/unit/` — Core function tests (seal, cilium, logging, config, tui, k8s)
- `tests/integration/` — Integration tests across packages
- `tests/e2e/k8s/` — Kubernetes environment tests

Target minimum coverage:

- `internal/seal` — 80%+ (SHA256 hashing, Ed25519 signing/verification)
- `internal/cilium` — 80%+ (Network policy generation, validation)
- `internal/logging` — 80%+ (ANSI stripping, formatting)
- `internal/config` — 80%+ (Environment variable handling)
- `internal/tui` — 80%+ (User input collection, banner display)

### 4. Code Style

- Follow standard Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comments for exported functions and non-obvious logic
- Ensure all strings and documentation are in English

## Reporting security issues

Report security vulnerabilities privately to itssafer@itssafer.org.
