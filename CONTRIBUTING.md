# Contributing to N-Audit Sentinel

This repository accepts contributions under the following guidelines:

- All code, documentation, and comments must be in English.
- Follow Go formatting: `make fmt`.
- Run linters and tests before submitting PRs: `make lint && make test`.

## Development Workflow

### 1. Setup
Clone the repository and ensure Go 1.25+ is installed. Example (Go):

```go
// Clone repository and run go mod tidy programmatically (illustrative)
package main

import (
	"log"
	"os/exec"
)

func main() {
	if err := exec.Command("git", "clone", "https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git").Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c", "cd N-Audit-Sentinel && go mod tidy").Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 2. Running Tests Locally
Before submitting a pull request, run the full test suite. Example (Go):

```go
// Run common development targets from Go
package main

import (
	"log"
	"os/exec"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		log.Fatalf("command failed: %v %v: %v", name, args, err)
	}
}

func main() {
	run("make", "fmt")
	run("make", "lint")
	run("make", "test")
	run("make", "test-e2e", "ENV=k3s")
	run("make", "security-scan")
	run("make", "verify-deps")
}
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
