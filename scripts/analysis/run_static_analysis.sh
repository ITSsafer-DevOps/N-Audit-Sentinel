#!/usr/bin/env bash
set -euo pipefail

# scripts/analysis/run_static_analysis.sh
# Run a stepped static analysis and collect artifacts under analysis/

outdir="analysis"
mkdir -p "$outdir"

echo "[analysis] workspace: $(pwd)" > "$outdir/summary.txt"
echo "[analysis] go version: $(go version)" >> "$outdir/summary.txt"

echo "[analysis] Running go vet..."
{ go vet ./... 2>&1 || true; } | tee "$outdir/go-vet.txt"

echo "[analysis] Running go test (unit + cover)..."
{ go test ./... -coverprofile="$outdir/cover.out" -covermode=count 2>&1 || true; } | tee "$outdir/go-test.txt"

echo "[analysis] Listing module dependencies..."
go list -m all > "$outdir/deps.txt" 2>/dev/null || true

if command -v golangci-lint >/dev/null 2>&1; then
  echo "[analysis] Running golangci-lint... (this may take a while)"
  golangci-lint run --out-format json ./... > "$outdir/golangci.json" 2> "$outdir/golangci.err" || true
else
  echo "[analysis] golangci-lint not found; skipping. To run, install golangci-lint: https://golangci-lint.run/usage/install/" | tee -a "$outdir/summary.txt"
fi

if command -v staticcheck >/dev/null 2>&1; then
  echo "[analysis] Running staticcheck..."
  staticcheck ./... > "$outdir/staticcheck.txt" 2>&1 || true
else
  echo "[analysis] staticcheck not found; skipping." >> "$outdir/summary.txt"
fi

if command -v govulncheck >/dev/null 2>&1; then
  echo "[analysis] Running govulncheck..."
  govulncheck ./... > "$outdir/govulncheck.txt" 2>&1 || true
else
  echo "[analysis] govulncheck not found; skipping." >> "$outdir/summary.txt"
fi

echo "[analysis] Generated artifacts in $outdir"
