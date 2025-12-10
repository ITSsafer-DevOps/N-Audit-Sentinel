#!/usr/bin/env bash
set -euo pipefail

echo "=========================================="
echo "N-Audit Sentinel Project Backup Tool"
echo "=========================================="
echo ""

# Configuration
BACKUP_DIR="$HOME/n-audit-backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
PROJECT_DIR="$(pwd)"
PROJECT_NAME="n-audit-sentinel"
BACKUP_NAME="${PROJECT_NAME}_${TIMESTAMP}"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}"

# Create backup directory
mkdir -p "$BACKUP_DIR"

echo "[1/5] Preparing backup: $BACKUP_NAME"
echo "Project: $PROJECT_DIR"
echo ""

# Create temporary staging directory
STAGE_DIR=$(mktemp -d)
trap "rm -rf $STAGE_DIR" EXIT

echo "[2/5] Copying project files..."
# Copy entire project
cp -r "$PROJECT_DIR" "$STAGE_DIR/$PROJECT_NAME"

# Clean up build artifacts and cache
echo "[3/5] Cleaning temporary files..."
cd "$STAGE_DIR/$PROJECT_NAME"
rm -rf bin/ *.tar.gz *.sha256 .terraform/ .terraform.lock.hcl terraform.tfstate* deploy/terraform/.terraform/ deploy/terraform/.terraform.lock.hcl deploy/terraform/terraform.tfstate* || true

echo "[4/5] Creating compressed archive..."
cd "$STAGE_DIR"
tar -czf "${BACKUP_PATH}.tar.gz" "$PROJECT_NAME"

# Calculate checksum
cd "$BACKUP_DIR"
sha256sum "${BACKUP_NAME}.tar.gz" > "${BACKUP_NAME}.tar.gz.sha256"

echo "[5/5] Backup complete!"
echo ""
echo "=========================================="
echo "Backup Details"
echo "=========================================="
echo "Location: ${BACKUP_PATH}.tar.gz"
echo "Size: $(du -h ${BACKUP_NAME}.tar.gz | cut -f1)"
echo "Checksum: ${BACKUP_NAME}.tar.gz.sha256"
echo ""
echo "To restore:"
echo "  cd $HOME"
echo "  tar -xzf ${BACKUP_PATH}.tar.gz"
echo ""
echo "All backups:"
ls -lh "$BACKUP_DIR" | grep "^-" || echo "  (none)"
echo ""
