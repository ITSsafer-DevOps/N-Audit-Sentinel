#!/usr/bin/env bash
################################################################################
# N-Audit Sentinel - Release & Push Script
#
# Non-interactive script that authenticates GH, cleans artifacts, creates a
# Gold Master backup, and force-pushes the repository to GitHub.
#
# Usage: ./release-and-push.sh
################################################################################

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GITHUB_ORG="ITSsafer-DevOps"
REPO_NAME="N-Audit-Sentinel"
GITHUB_REMOTE_URL="https://github.com/${GITHUB_ORG}/${REPO_NAME}.git"
BACKUP_SCRIPT="${PROJECT_DIR}/scripts/backup-project.sh"

# Color helpers
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}==> STEP 1: Authenticate GitHub (non-interactive)${NC}"
# Authenticate gh using provided token (non-interactive)
echo "Authenticating with GitHub CLI..."
echo "$TOKEN" | gh auth login --with-token >/dev/null 2>&1 || {
    echo -e "${RED}✗ ERROR: gh authentication failed. Aborting.${NC}"
    exit 1
}
echo -e "${GREEN}✓ gh authenticated${NC}"

echo -e "\n${YELLOW}==> STEP 2: Deep Clean - remove old artifacts from project root${NC}"
cd "$PROJECT_DIR"

# Remove binary and directories (idempotent)
rm -f n-audit-release || true
rm -rf bin/ || true
rm -rf logs/ || true

# Remove any old/dirty archives or checksums in the project root only
shopt -s nullglob
deleted_count=0
for f in *.tar.gz *.sha256; do
    if [[ -e "$f" ]]; then
        rm -f -- "$f" && ((deleted_count++))
    fi
done
shopt -u nullglob

echo -e "${GREEN}✓ Old artifacts/backups removed from project root.${NC}"

echo -e "\n${YELLOW}==> STEP 3: Maintenance - go mod tidy${NC}"
if [[ ! -f go.mod ]]; then
    echo -e "${RED}✗ ERROR: go.mod not found in project root. Aborting.${NC}"
    exit 1
fi
go mod tidy
echo -e "${GREEN}✓ go mod tidy completed${NC}"

echo -e "\n${YELLOW}==> STEP 4: Create Gold Master archive${NC}"
if [[ ! -f "$BACKUP_SCRIPT" ]]; then
    echo -e "${RED}✗ ERROR: Backup script not found at: $BACKUP_SCRIPT${NC}"
    exit 1
fi
chmod +x "$BACKUP_SCRIPT"

# Record current latest backup timestamp before running
BACKUP_DIR="$HOME/n-audit-backups"
mkdir -p "$BACKUP_DIR"
prev_latest=""
if compgen -G "$BACKUP_DIR/*.tar.gz" > /dev/null; then
    prev_latest=$(ls -1t "$BACKUP_DIR"/*.tar.gz | head -n1 || true)
fi

echo "Running backup script..."
bash "$BACKUP_SCRIPT"

# Find newest archive after backup
new_latest=""
if compgen -G "$BACKUP_DIR/*.tar.gz" > /dev/null; then
    new_latest=$(ls -1t "$BACKUP_DIR"/*.tar.gz | head -n1 || true)
fi

if [[ -z "$new_latest" ]]; then
    echo -e "${RED}✗ ERROR: No backup archive was created in $BACKUP_DIR${NC}"
    exit 1
fi

if [[ "$new_latest" == "$prev_latest" ]]; then
    echo -e "${YELLOW}⚠ Warning: Backup did not create a new archive (latest unchanged).${NC}"
else
    echo -e "${GREEN}✓ New Gold Master archive created at: ${new_latest}${NC}"
fi

echo -e "\n${YELLOW}==> STEP 5: Git initialization, commit, and force-push to ${GITHUB_REMOTE_URL}${NC}"
cd "$PROJECT_DIR"

# Initialize git repo if missing
if [[ ! -d .git ]]; then
    git init -b main
    git config user.email "noreply@itssafer-devops.com"
    git config user.name "N-Audit Release Bot"
    echo "Initialized new git repository (main)"
else
    echo "Git repository already initialized"
fi

# Stage all files
git add .

# Commit if there are changes
if git diff --cached --quiet; then
    echo "No changes to commit"
else
    git commit -m "feat: initial beta release"
    echo "Created commit: feat: initial beta release"
fi

# Ensure remote is set (ignore error if exists)
if git remote get-url origin >/dev/null 2>&1; then
    echo "Remote 'origin' already set to: $(git remote get-url origin)"
else
    git remote add origin "$GITHUB_REMOTE_URL" || true
    echo "Added remote origin -> $GITHUB_REMOTE_URL"
fi

# Force push to ensure remote matches our clean local state
echo "Pushing to origin main (force)..."
git push -u origin main --force

echo -e "${GREEN}✓ Force push completed. Release published to ${GITHUB_REMOTE_URL}${NC}"

echo -e "\n${GREEN}All steps completed successfully.${NC}"
