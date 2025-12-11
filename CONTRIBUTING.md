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

---

## Community & Support

**Connect with the maintainers and community:**

- 💼 **LinkedIn:** [ITSsafer DevOps Team](https://www.linkedin.com/in/itsafer-devops/) — Follow for updates and industry insights
- 🐙 **GitHub Discussions:** [Project Discussions](https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/discussions) — Ask questions, share ideas
- 📚 **Documentation:** [DEPLOYMENT.md](./DEPLOYMENT.md), [SECURITY.md](./SECURITY.md), [TESTING_AND_VERIFICATION.md](./TESTING_AND_VERIFICATION.md)
- 🔗 **Related Resources:** [Architecture Diagrams](./docs/ENTERPRISE_ARCHITECTURE.md), [Enterprise-Grade Auditing Guide](./VERIFICATION_GUIDE.md)

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Make changes and test: `make test`
4. Format code: `make fmt`
5. Submit PR with description linking to related issues
6. Await review from maintainers
7. Address feedback and push updates
8. Merge upon approval

## Code Review Standards

- **Functionality:** Code must pass all tests and not introduce regressions
- **Documentation:** Updated markdown and inline comments required
- **Security:** No hardcoded credentials or unsafe patterns; vulnerability scans must pass
- **Performance:** No significant performance degradation; benchmarks welcome
- **Maintainability:** Clear naming, proper error handling, DI patterns for testing

