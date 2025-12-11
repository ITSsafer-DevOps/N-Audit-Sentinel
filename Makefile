GO := go
BINARY := bin/n-audit-sentinel
VERSION ?= v1.0.0-Beta
BUILD_TOOL := cmd/build-tool/main.go

BIN_DIR := bin
RELEASE_DIR := releases

.PHONY: all build clean fmt lint test test-e2e release verify-deps security-scan help backup-final build-tool

all: build

build:
	@echo "Building binaries..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/n-audit-sentinel ./cmd/n-audit-sentinel
	$(GO) build -o $(BIN_DIR)/n-audit ./cmd/n-audit-cli
	$(GO) build -o $(BIN_DIR)/n-audit-release ./cmd/n-audit-release

clean:
	@echo "Cleaning artifacts..."
	@$(GO) run $(BUILD_TOOL) -cmd verify-deps
	rm -rf $(BIN_DIR)
	find $(RELEASE_DIR) -maxdepth 1 \( -name "*.tar.gz" -o -name "*.sha256" \) -delete 2>/dev/null || true

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

lint:
	@echo "Linting..."
	@$(GO) run $(BUILD_TOOL) -cmd check-tool golangci-lint 2>/dev/null | grep -q FOUND && \
		golangci-lint run ./... || \
		$(GO) vet ./...

test:
	@echo "Running unit and integration tests..."
	$(GO) test ./... -count=1 -v

test-e2e:
	@echo "Running E2E tests for ENV=$(ENV)"
	$(GO) run $(BUILD_TOOL) -cmd verify-e2e-env -env $(ENV)
	$(GO) test ./tests/e2e/... -run Test -v

security-scan:
	@echo "Running security scans (govulncheck if installed)..."
	@$(GO) run $(BUILD_TOOL) -cmd check-tool govulncheck 2>/dev/null | grep -q FOUND && \
		govulncheck ./... || true || \
		echo "govulncheck not installed - skipping"

verify-deps:
	@echo "Verifying dependencies: go, docker, kubectl"
	$(GO) run $(BUILD_TOOL) -cmd verify-deps

release: clean build
	@echo "Creating release artifacts..."
	$(GO) run $(BUILD_TOOL) -cmd release-bin -version $(VERSION) -bin-dir $(BIN_DIR) -release-dir $(RELEASE_DIR)

backup-final:
	@echo "Creating final deterministic backup (gold master)..."
	$(GO) run $(BUILD_TOOL) -cmd backup-source -version $(VERSION) -release-dir $(RELEASE_DIR)

help:
	@echo "N-Audit Sentinel - Makefile Targets"
	@echo ""
	@echo "Build & Test:"
	@echo "  make build            Build all binaries"
	@echo "  make test             Run all unit and integration tests"
	@echo "  make test-e2e ENV=k3s Run end-to-end tests"
	@echo "  make fmt              Format Go source code"
	@echo "  make lint             Run linters"
	@echo ""
	@echo "Security & Verification:"
	@echo "  make security-scan    Run vulnerability scans"
	@echo "  make verify-deps      Check for required tools"
	@echo ""
	@echo "Release:"
	@echo "  make release          Create tarballs with checksums (in releases/)"
	@echo "  make backup-final     Create deterministic goldmaster backup (in releases/)"
	@echo "  make clean            Remove build artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make fmt && make lint && make test"
	@echo "  make test-e2e ENV=k3s"
	@echo "  make release"
