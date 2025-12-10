#!/usr/bin/env bash
################################################################################
# final-release-and-push.sh
#
# Complete release workflow for N-Audit Sentinel Beta:
#  - gh authentication
#  - README documentation updates (origin note, contact, mermaid)
#  - .gitignore security audit
#  - deep clean artifacts
#  - go mod tidy
#  - Gold Master backup
#  - git init/commit/force-push to GitHub
#
# Usage: ./final-release-and-push.sh
################################################################################

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Token must be passed via environment variable to avoid committing secrets
TOKEN="${GITHUB_TOKEN:-}"
if [[ -z "$TOKEN" ]]; then
  echo "ERROR: GITHUB_TOKEN environment variable not set."
  echo "Usage: export GITHUB_TOKEN=\"your_token_here\"; ./final-release-and-push.sh"
  exit 1
fi
GITHUB_ORG="ITSsafer-DevOps"
REPO_NAME="N-Audit-Sentinel"
REMOTE_URL="https://github.com/${GITHUB_ORG}/${REPO_NAME}.git"
BACKUP_SCRIPT="${PROJECT_DIR}/scripts/backup-project.sh"

################################################################################
# STEP 1: AUTHENTICATION
################################################################################
echo ""
echo "==> STEP 1: Authenticate gh (non-interactive)"
echo "$TOKEN" | gh auth login --with-token
echo "✓ gh authenticated"

################################################################################
# STEP 2: DOCUMENTATION POLISH
################################################################################
echo ""
echo "==> STEP 2: Documentation Polish"
cd "$PROJECT_DIR"

# 2a. Prepend Origin Note to README.md
if [[ -f README.md ]]; then
  echo "==> Updating README Origin Note..."
  ORIGIN_NOTE="**Origin Note:** This project was architected and developed as a **proactive R&D initiative** (\"going the extra mile\") within the recruitment process for **Nethemba s.r.o.**"
  
  if ! grep -Fq "$ORIGIN_NOTE" README.md; then
    printf "%s\n\n%s\n" "$ORIGIN_NOTE" "$(cat README.md)" > README.md
    echo "✓ Origin Note prepended to README.md"
  else
    echo "✓ Origin Note already present in README.md"
  fi
fi

# 2b. Replace Website with Contact
if [[ -f README.md ]]; then
  echo "==> Updating Contact Information..."
  if grep -q "Website:[[:space:]]*https://www.nethemba.com" README.md; then
    sed -i 's|Website:[[:space:]]*https://www.nethemba.com|Contact: itssafer@itssafer.org|g' README.md
    echo "✓ Website replaced with Contact in README.md"
  else
    echo "✓ No Website line found (already updated or not present)"
  fi
fi

# 2c. Fix Mermaid diagram syntax
echo "==> Fixing Mermaid Diagram Syntax..."
for md_file in README.md DEPLOYMENT.md VERIFICATION_GUIDE.md; do
  if [[ -f "$md_file" ]]; then
    # Replace any opening fence with mermaid keyword to standard ```mermaid
    awk '
      BEGIN{IGNORECASE=1}
      {
        if ($0 ~ /^```[[:space:]]*.*mermaid.*$/) {
          print "```mermaid"
          next
        }
        print $0
      }
    ' "$md_file" > "${md_file}.tmp" && mv "${md_file}.tmp" "$md_file"
    echo "✓ Normalized mermaid fences in $md_file"
  fi
done

################################################################################
# STEP 3: SECURITY AUDIT (.gitignore)
################################################################################
echo ""
echo "==> STEP 3: Security Audit (.gitignore)"
GITIGNORE="${PROJECT_DIR}/.gitignore"
touch "$GITIGNORE"

PATTERNS_TO_ADD=(
  "terraform.tfvars"
  "*.log"
  "bin/"
  "n-audit-release"
  "*.sha256"
)

added_count=0
for pat in "${PATTERNS_TO_ADD[@]}"; do
  if ! grep -qxF "$pat" "$GITIGNORE"; then
    echo "$pat" >> "$GITIGNORE"
    echo "✓ Added to .gitignore: $pat"
    ((added_count++))
  fi
done

if [[ $added_count -eq 0 ]]; then
  echo "✓ .gitignore already contains all required patterns"
else
  echo "✓ Added $added_count pattern(s) to .gitignore"
fi

################################################################################
# STEP 4: DEEP CLEAN
################################################################################
echo ""
echo "==> STEP 4: Deep Clean - Removing Artifacts"
cd "$PROJECT_DIR"

# Remove binary and directories
echo "Removing n-audit-release..."
rm -f n-audit-release || true
echo "✓ n-audit-release removed"

echo "Removing bin/ directory..."
rm -rf bin/ || true
echo "✓ bin/ removed"

echo "Removing logs/ directory..."
rm -rf logs/ || true
echo "✓ logs/ removed"

# Remove archives and checksums from project root
echo "Removing *.tar.gz and *.sha256 files..."
shopt -s nullglob
for archive in *.tar.gz *.sha256; do
  rm -f -- "$archive"
  echo "✓ Removed $archive"
done
shopt -u nullglob

echo "✓ Deep clean complete"

################################################################################
# STEP 5: MAINTENANCE
################################################################################
echo ""
echo "==> STEP 5: Maintenance - go mod tidy"
cd "$PROJECT_DIR"

if [[ ! -f go.mod ]]; then
  echo "✗ ERROR: go.mod not found. Exiting."
  exit 1
fi

go mod tidy
echo "✓ go mod tidy completed"

################################################################################
# STEP 6: GOLD MASTER ARCHIVE
################################################################################
echo ""
echo "==> STEP 6: Gold Master Archive"
if [[ ! -f "$BACKUP_SCRIPT" ]]; then
  echo "✗ ERROR: Backup script not found at $BACKUP_SCRIPT. Exiting."
  exit 1
fi

chmod +x "$BACKUP_SCRIPT"
echo "Running backup script..."
cd "$PROJECT_DIR" && bash "$BACKUP_SCRIPT"
echo "✓ Gold Master archive created"

################################################################################
# STEP 7: GIT PUBLISH
################################################################################
echo ""
echo "==> STEP 7: Git Publish (init/commit/force-push)"
cd "$PROJECT_DIR"

# Re-initialize git on main
echo "Initializing git repository..."
git init -b main
git config user.email "noreply@itssafer-devops.com"
git config user.name "N-Audit Release Bot"
echo "✓ Git initialized on main"

# Stage all files
echo "Staging all files..."
git add .
echo "✓ Files staged"

# Commit
echo "Creating commit..."
git commit -m "feat: initial beta release (Nethemba initiative)" || echo "✓ Commit (no changes or already committed)"
echo "✓ Commit created/updated"

# Set remote
echo "Configuring remote..."
if git remote get-url origin &>/dev/null; then
  git remote set-url origin "$REMOTE_URL"
  echo "✓ Remote origin updated to $REMOTE_URL"
else
  git remote add origin "$REMOTE_URL" 2>/dev/null || true
  echo "✓ Remote origin added: $REMOTE_URL"
fi

# Force push
echo "Force-pushing to GitHub..."
git push -u origin main --force
echo "✓ Force push completed"

################################################################################
# FINAL MESSAGE
################################################################################
echo ""
echo "==============================================================================="
echo "✓ FINAL RELEASE & PUSH WORKFLOW COMPLETE"
echo "==============================================================================="
echo "Repository: $REMOTE_URL"
echo "Branch: main"
echo ""
echo "Next steps:"
echo "  1. Verify repository on GitHub: https://github.com/${GITHUB_ORG}/${REPO_NAME}"
echo "  2. Create a GitHub Release from main branch"
echo "  3. Attach Gold Master archive from ~/n-audit-backups/"
echo "  4. Enable branch protection rules"
echo "  5. Add CODEOWNERS and governance files"
echo ""
echo "SECURITY NOTE: Rotate the token used in this script at:"
echo "  https://github.com/settings/tokens"
echo "==============================================================================="
echo ""
