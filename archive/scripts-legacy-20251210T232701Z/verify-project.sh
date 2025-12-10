#!/usr/bin/env bash
set -uo pipefail

echo "=========================================="
echo "N-Audit Sentinel Dry-Run Verification"
echo "=========================================="
echo "This script validates the project without requiring Docker/Kubernetes"
echo ""

PASS=0
FAIL=0

check_pass() {
    echo "✓ $1"
    ((PASS++))
}

check_fail() {
    echo "✗ $1"
    ((FAIL++))
}

# Phase 1: Project Structure
echo "[1/7] Validating project structure..."
if [ -f "go.mod" ]; then check_pass "go.mod exists"; else check_fail "go.mod missing"; fi
if [ -f "Dockerfile" ]; then check_pass "Dockerfile exists"; else check_fail "Dockerfile missing"; fi
if [ -f "Makefile" ]; then check_pass "Makefile exists"; else check_fail "Makefile missing"; fi
if [ -f "README.md" ]; then check_pass "README.md exists"; else check_fail "README.md missing"; fi
if [ -f "DEPLOYMENT.md" ]; then check_pass "DEPLOYMENT.md exists"; else check_fail "DEPLOYMENT.md missing"; fi
if [ -f "LICENSE" ]; then check_pass "LICENSE exists"; else check_fail "LICENSE missing"; fi
if [ -d "cmd/n-audit-sentinel" ]; then check_pass "cmd/n-audit-sentinel exists"; else check_fail "cmd/n-audit-sentinel missing"; fi
if [ -d "cmd/n-audit-cli" ]; then check_pass "cmd/n-audit-cli exists"; else check_fail "cmd/n-audit-cli missing"; fi
if [ -d "internal" ]; then check_pass "internal package exists"; else check_fail "internal missing"; fi
echo ""

# Phase 2: Go Build
echo "[2/7] Building Go binaries..."
if go build -o bin/n-audit-sentinel ./cmd/n-audit-sentinel 2>&1; then
    check_pass "n-audit-sentinel builds successfully"
else
    check_fail "n-audit-sentinel build failed"
fi

if go build -o bin/n-audit ./cmd/n-audit-cli 2>&1; then
    check_pass "n-audit CLI builds successfully"
else
    check_fail "n-audit CLI build failed"
fi

if go build -o bin/n-audit-release ./cmd/n-audit-release 2>&1; then
    check_pass "n-audit-release builds successfully"
else
    check_fail "n-audit-release build failed"
fi
echo ""

# Phase 3: Run Tests
echo "[3/7] Running test suite..."
if go test -count=1 ./... > /tmp/test-output.log 2>&1; then
    check_pass "All tests pass"
    TEST_COUNT=$(grep -c "^ok" /tmp/test-output.log || echo "0")
    echo "  → $TEST_COUNT package(s) tested"
else
    check_fail "Tests failed"
    echo "  → See /tmp/test-output.log for details"
fi
echo ""

# Phase 4: Race Detector
echo "[4/7] Running race detector..."
if go test -race -count=1 ./... > /tmp/race-output.log 2>&1; then
    check_pass "No race conditions detected"
else
    check_fail "Race conditions found"
    echo "  → See /tmp/race-output.log for details"
fi
echo ""

# Phase 5: Code Formatting
echo "[5/7] Checking code formatting..."
UNFORMATTED=$(gofmt -l . | grep -v vendor || true)
if [ -z "$UNFORMATTED" ]; then
    check_pass "All Go files are formatted"
else
    check_fail "Some files need formatting"
    echo "$UNFORMATTED"
fi
echo ""

# Phase 6: Terraform Validation
echo "[6/7] Validating Terraform configuration..."
cd deploy/terraform
if [ -f "main.tf" ] && [ -f "variables.tf" ] && [ -f "outputs.tf" ]; then
    check_pass "Terraform files exist"
    
    if terraform fmt -check > /dev/null 2>&1; then
        check_pass "Terraform files are formatted"
    else
        check_fail "Terraform files need formatting"
    fi
    
    if terraform init > /tmp/tf-init.log 2>&1 && terraform validate > /tmp/tf-validate.log 2>&1; then
        check_pass "Terraform configuration is valid"
    else
        check_fail "Terraform validation failed"
        echo "  → See /tmp/tf-init.log and /tmp/tf-validate.log"
    fi
else
    check_fail "Terraform files incomplete"
fi
cd ../..
echo ""

# Phase 7: Documentation Validation
echo "[7/7] Validating documentation..."
if grep -q "Security Considerations" README.md; then
    check_pass "README includes Security Considerations"
else
    check_fail "README missing Security Considerations"
fi

if grep -q "internal/validation" README.md; then
    check_pass "README documents validation package"
else
    check_fail "README missing validation package docs"
fi

if grep -q "terraform apply" DEPLOYMENT.md; then
    check_pass "DEPLOYMENT.md includes Terraform instructions"
else
    check_fail "DEPLOYMENT.md incomplete"
fi

if [ -f "deploy/terraform/terraform.tfvars.example" ]; then
    check_pass "Example tfvars file exists"
else
    check_fail "Example tfvars missing"
fi
echo ""

# Summary
echo "=========================================="
echo "Verification Summary"
echo "=========================================="
echo "✓ Passed: $PASS"
echo "✗ Failed: $FAIL"
echo ""

if [ $FAIL -eq 0 ]; then
    echo "🎉 All checks passed! Project is ready for deployment."
    echo ""
    echo "Next steps:"
    echo "  1. Set up a Kubernetes cluster with Cilium CNI"
    echo "  2. Install docker, kubectl, and terraform"
    echo "  3. Run: ./local-deploy-test.sh"
    exit 0
else
    echo "⚠️  Some checks failed. Review the output above."
    exit 1
fi
