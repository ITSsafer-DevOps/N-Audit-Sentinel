SHELL := /bin/bash
SHELL := /bin/bash
GO := go
BINARY := bin/n-audit-sentinel
VERSION ?= v1.0.0-Beta

BIN_DIR := bin
RELEASE_DIR := releases

.PHONY: all build clean fmt lint test test-e2e release verify-deps security-scan help backup-final

all: build

build:
	@echo "Building binaries..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/n-audit-sentinel ./cmd/n-audit-sentinel
	$(GO) build -o $(BIN_DIR)/n-audit ./cmd/n-audit-cli
	$(GO) build -o $(BIN_DIR)/n-audit-release ./cmd/n-audit-release

clean:
	@echo "Cleaning artifacts..."
	rm -rf $(BIN_DIR) $(RELEASE_DIR)/*.tar.gz $(RELEASE_DIR)/*.sha256

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

lint:
	@echo "Linting..."
	if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		$(GO) vet ./...; \
	fi

test:
	@echo "Running unit and integration tests..."
	$(GO) test ./... -count=1 -v

test-e2e:
	@echo "Running E2E tests for ENV=$(ENV)"
	if [ "$(ENV)" = "k3s" ]; then \
		echo "K3s environment selected"; \
	fi
	$(GO) test ./tests/e2e/... -run Test -v

security-scan:
	@echo "Running security scans (govulncheck if installed)..."
	if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./... || true; \
	else \
		echo "govulncheck not installed - skipping"; \
	fi

verify-deps:
	@echo "Verifying dependencies: go, docker, kubectl"
	command -v go >/dev/null || (echo "ERROR: go missing" && exit 1)
	command -v docker >/dev/null || echo "WARN: docker not found - optional"
	command -v kubectl >/dev/null || echo "WARN: kubectl not found - optional"

release: clean build
	@echo "Creating release artifacts..."
	@mkdir -p $(RELEASE_DIR)
	tar -czf $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-bin.tar.gz -C $(BIN_DIR) n-audit-sentinel
	sha256sum $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-bin.tar.gz > $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-bin.tar.gz.sha256
	@echo "Release artifacts created:"; ls -lh $(RELEASE_DIR) | grep -E '\.tar\.gz|\.sha256' | awk '{print "  " $$9, "(" $$5 ")"}'

backup-final:
	@echo "Creating final deterministic backup (gold master)"
	@mkdir -p $(RELEASE_DIR)
	git archive --format=tar --prefix=n-audit-sentinel-$(VERSION)-source/ HEAD | gzip -9 > $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-goldmaster.tar.gz
	sha256sum $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-goldmaster.tar.gz > $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-goldmaster.tar.gz.sha256
	@echo "Backup created:"; ls -lh $(RELEASE_DIR)/n-audit-sentinel-$(VERSION)-goldmaster* | awk '{print "  " $$9, "(" $$5 ")"}'

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
