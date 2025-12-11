SHELL := /bin/bash
GO := go
BINARY := bin/n-audit-sentinel
VERSION ?= v1.0.0-Beta

.PHONY: all build clean fmt lint test test-e2e release release-dry-run verify-deps security-scan help

all: build

build:
	@echo "Building binaries..."
	@mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/n-audit-sentinel || exit 1

clean:
	@echo "Cleaning artifacts..."
	rm -rf bin *.tar.gz *.sha256

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

lint:
	@echo "Linting..."
	if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		$(GO) vet ./...; \
		$(GO) list ./... | xargs -n1 $(GO) test -run Test -bench . -benchmem >/dev/null 2>&1 || true; \
	fi

test:
	@echo "Running unit and integration tests..."
	$(GO) test ./... -count=1

test-e2e:
	@echo "Running E2E tests for ENV=$(ENV)"
	if [ "$(ENV)" = "k3s" ]; then \
		echo "K3s environment selected - ensure kubeconfig points to k3s"; \
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
	@echo "Verifying dependencies: go, docker, kubectl";
	command -v go >/dev/null || (echo "go missing" && exit 1);
	command -v docker >/dev/null || echo "docker not found - optional";
	command -v kubectl >/dev/null || echo "kubectl not found - optional";

release: clean build
	@echo "Creating release artifacts..."
	tar -czf n-audit-sentinel-beta-bin.tar.gz -C bin n-audit-sentinel
	sha256sum n-audit-sentinel-beta-bin.tar.gz > n-audit-sentinel-beta-bin.tar.gz.sha256
	tar --exclude='.git' --exclude='bin' -czf n-audit-sentinel-beta-source-code.tar.gz .
	sha256sum n-audit-sentinel-beta-source-code.tar.gz > n-audit-sentinel-beta-source-code.tar.gz.sha256

release-dry-run:
	@echo "Dry-run release (no upload)"
	@echo VERSION=$(VERSION)

help:
	@echo "Makefile targets:"
	@echo "  build           Build primary binary"
	@echo "  test            Run unit+integration tests"
	@echo "  test-e2e ENV=   Run e2e tests for environment (k8s,k3s)"
	@echo "  fmt             Format sources"
	@echo "  lint            Run linters (golangci-lint if available)"
	@echo "  security-scan   Run govulncheck if installed"
	@echo "  release         Create tar.gz release artifacts"
	@echo "  verify-deps     Verify required tools are present"
## N-Audit Sentinel Makefile
## Usage examples:
##   make build                # Build local host binaries
##   make test                 # Run all tests
##   make fmt                  # Format Go code
##   make release VERSION=v1.2.3  # Produce versioned linux/amd64 tarball + .sha256
##   make clean                # Remove build artifacts

BIN_DIR := bin
VERSION ?=

.PHONY: build test fmt release clean help all tf-validate deploy-local attach exec-exit destroy-local

all: build

build:
	@echo "[make] Building host binaries..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/n-audit-sentinel ./cmd/n-audit-sentinel
	go build -o $(BIN_DIR)/n-audit ./cmd/n-audit-cli
	go build -o $(BIN_DIR)/n-audit-release ./cmd/n-audit-release

test:
	@echo "[make] Running tests..."
	go test -count=1 -v ./...

fmt:
	@echo "[make] Formatting source..."
	go fmt ./...

release:
	@if [ -z "$(VERSION)" ]; then echo "ERROR: VERSION required (e.g. make release VERSION=v1.2.3)"; exit 1; fi
	@echo "[make] Building release artifacts for $(VERSION)..."
	go run ./cmd/n-audit-release $(VERSION)

clean:
	@echo "[make] Cleaning artifacts..."
	rm -rf $(BIN_DIR) *.tar.gz *.sha256

help:
	@echo "N-Audit Sentinel Makefile Targets:"; \
	echo "  build    - Build local binaries into $(BIN_DIR)"; \
	echo "  test     - Run all Go tests"; \
	echo "  fmt      - Format Go code"; \
	echo "  release  - Build versioned linux/amd64 tarball (requires VERSION)"; \
	echo "  clean    - Remove build artifacts"; \
	echo "  all      - Alias for build"; \
	echo "  help     - Show this help"; \
	echo "  tf-validate - Terraform fmt (check) and validate in deploy/terraform"; \
	echo "  deploy-local - Build Docker image, apply Terraform for local cluster test"; \
	echo "  attach   - Attach to the sentinel pod interactively"; \
	echo "  exec-exit - Trigger graceful exit via n-audit CLI"; \
	echo "  destroy-local - Destroy local Terraform deployment"

tf-validate:
	@echo "[make] Validating Terraform configuration..."
	@cd deploy/terraform && terraform fmt -check && terraform validate

deploy-local:
	@echo "[make] Building and deploying for local test..."
	@command -v docker >/dev/null 2>&1 || { echo "ERROR: docker not found in PATH"; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo "ERROR: kubectl not found in PATH"; exit 1; }
	@command -v terraform >/dev/null 2>&1 || { echo "ERROR: terraform not found in PATH"; exit 1; }
	docker build -t n-audit-sentinel:local-test .
	@echo "[make] Load image into your local cluster (uncomment one line in Makefile):"
	@echo "  # minikube image load n-audit-sentinel:local-test"
	@echo "  # k3d image import n-audit-sentinel:local-test"
	@echo "  # kind load docker-image n-audit-sentinel:local-test"
	@echo "[make] Applying Terraform with local image settings..."
	@cd deploy/terraform && \
		echo 'namespace       = "default"' > terraform.tfvars && \
		echo 'image_name      = "n-audit-sentinel"' >> terraform.tfvars && \
		echo 'image_tag       = "local-test"' >> terraform.tfvars && \
		echo 'pvc_storage_size = "1Gi"' >> terraform.tfvars && \
		terraform init && terraform apply -auto-approve

attach:
	@echo "[make] Attaching to sentinel pod..."
	@command -v kubectl >/dev/null 2>&1 || { echo "ERROR: kubectl not found in PATH"; exit 1; }
	kubectl attach -it n-audit-sentinel

exec-exit:
	@echo "[make] Triggering graceful exit..."
	@command -v kubectl >/dev/null 2>&1 || { echo "ERROR: kubectl not found in PATH"; exit 1; }
	kubectl exec -it n-audit-sentinel -- /usr/local/bin/n-audit

destroy-local:
	@echo "[make] Destroying local deployment..."
	@command -v terraform >/dev/null 2>&1 || { echo "ERROR: terraform not found in PATH"; exit 1; }
	@if [ -d deploy/terraform/.terraform ]; then \
		cd deploy/terraform && terraform destroy -auto-approve; \
	else \
		echo "[make] No Terraform state found, skipping destroy."; \
	fi
