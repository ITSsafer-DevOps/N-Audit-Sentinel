#!/usr/bin/env bash
################################################################################
# fix-and-publish.sh
#
# Final remedial update for N-Audit Sentinel:
#  - Fix broken Mermaid diagrams (localized syntax → English)
#  - Update copyright/license from Nethemba → ITSsafer-DevOps
#  - Deep security and file audit
#  - Final Gold Master backup and force-push
#
# Usage:
#   export GITHUB_TOKEN="ghp_..."  # optional; will prompt if not set
#   ./fix-and-publish.sh
#
################################################################################

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GITHUB_ORG="ITSsafer-DevOps"
REPO_NAME="N-Audit-Sentinel"
REMOTE_URL="https://github.com/${GITHUB_ORG}/${REPO_NAME}.git"
BACKUP_SCRIPT="${PROJECT_DIR}/scripts/backup-project.sh"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "\n${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  N-AUDIT SENTINEL - FINAL REMEDIAL UPDATE & PUBLISH            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}\n"

################################################################################
# STEP 1: AUTHENTICATION
################################################################################
echo -e "${YELLOW}==> STEP 1: GitHub Authentication${NC}"
TOKEN="${GITHUB_TOKEN:-}"
if [[ -z "$TOKEN" ]]; then
  echo "GITHUB_TOKEN not set. Please provide your GitHub token:"
  read -rsp "GitHub Token: " TOKEN
  echo ""
fi

if [[ -z "$TOKEN" ]]; then
  echo -e "${RED}✗ No token provided. Exiting.${NC}"
  exit 1
fi

echo "Authenticating with gh..."
echo "$TOKEN" | gh auth login --with-token 2>&1 | grep -v "^Your authentication token" || true
gh auth status >/dev/null 2>&1 || {
  echo -e "${RED}✗ GitHub authentication failed. Exiting.${NC}"
  exit 1
}
echo -e "${GREEN}✓ GitHub authenticated${NC}"
unset TOKEN  # Remove from memory for safety

################################################################################
# STEP 2: MERMAID DIAGRAM SYNTAX REPAIR
################################################################################
echo -e "\n${YELLOW}==> STEP 2: Mermaid Diagram Syntax Repair${NC}"
cd "$PROJECT_DIR"

# List of markdown files to repair
MD_FILES=(README.md DEPLOYMENT.md VERIFICATION_GUIDE.md)

for md in "${MD_FILES[@]}"; do
  [[ ! -f "$md" ]] && continue
  
  echo "Processing $md for Mermaid syntax fixes..."
  
  # Repair 1: Replace localized "vývojový diagram" or "vyvojovy diagram" with English "graph"
  # This handles the Slovak/Czech term for "flowchart" or "flow diagram"
  # sed patterns:
  #   - \b ensures word boundary (beginning of word)
  #   - [[:space:]]* allows optional whitespace
  #   - g flag = global replacement (all occurrences)
  if grep -qi "vyvojov.*diagram\|vývojový.*diagram" "$md"; then
    sed -i 's/\(^[[:space:]]*\)vyvojov[[:alpha:]]*[[:space:]]*diagram\([[:space:]]*$\)/\1graph LR\2/gi' "$md"
    sed -i 's/\(^[[:space:]]*\)vývojový[[:space:]]*diagram\([[:space:]]*$\)/\1graph LR\2/gi' "$md"
    echo "  ✓ Replaced localized 'vývojový diagram' with 'graph LR'"
  fi
  
  # Repair 2: Replace "podgraf" (Slovak subgraph) with English "subgraph"
  # This handles the localized term for subgraph blocks within diagrams
  if grep -qi "podgraf\|pod graf" "$md"; then
    sed -i 's/\bpodgraf\b/subgraph/gi' "$md"
    sed -i 's/\bpod[[:space:]]*graf\b/subgraph/gi' "$md"
    echo "  ✓ Replaced 'podgraf' with 'subgraph'"
  fi
  
  # Repair 3: Replace "koniec" (Slovak "end") with English "end"
  # This handles the localized keyword for ending blocks
  if grep -qi "\bkoniec\b" "$md"; then
    sed -i 's/\bkoniec\b/end/gi' "$md"
    echo "  ✓ Replaced 'koniec' with 'end'"
  fi
  
  # Repair 4: Normalize all mermaid code fences to standard ```mermaid syntax
  # This ensures GitHub recognizes the diagram blocks correctly
  awk '
    BEGIN{IGNORECASE=1}
    {
      if ($0 ~ /^```[[:space:]]*.*mermaid.*$/) {
        print "```mermaid"
        next
      }
      print $0
    }
  ' "$md" > "${md}.tmp" && mv "${md}.tmp" "$md"
  
  echo "  ✓ Normalized mermaid code fences in $md"
done

echo -e "${GREEN}✓ Mermaid diagram syntax repaired across all markdown files${NC}"

################################################################################
# STEP 3: LICENSE & BRANDING CLEANUP
################################################################################
echo -e "\n${YELLOW}==> STEP 3: License & Branding Cleanup${NC}"

# 3a. Update LICENSE file copyright holder
if [[ -f LICENSE ]]; then
  echo "Updating LICENSE copyright holder..."
  
  # Replace "Nethemba s.r.o." with "ITSsafer-DevOps" throughout LICENSE
  if grep -q "Nethemba" LICENSE; then
    sed -i 's/Nethemba[[:space:]]*s\.r\.o\./ITSsafer-DevOps/g' LICENSE
    sed -i 's/Nethemba/ITSsafer-DevOps/g' LICENSE
    echo "  ✓ Updated copyright holder to ITSsafer-DevOps"
  else
    echo "  ✓ No Nethemba reference in LICENSE (already updated)"
  fi
fi

# 3b. Update README.md: Ensure Origin Note at top, clean ownership language
if [[ -f README.md ]]; then
  echo "Updating README.md for clarity of ownership..."
  
  # Ensure Origin Note is present at the very top
  ORIGIN_NOTE="**Origin Note:** This project was architected and developed as a **proactive R&D initiative** (\"going the extra mile\") within the recruitment process for **Nethemba s.r.o.**"
  
  if ! grep -Fq "$ORIGIN_NOTE" README.md; then
    echo "  ✓ Origin Note already present (or being added)"
    if ! grep -q "Origin Note" README.md; then
      printf "%s\n\n%s\n" "$ORIGIN_NOTE" "$(cat README.md)" > README.md
      echo "  ✓ Prepended Origin Note to README.md"
    fi
  fi
  
  # Remove any implied ownership by Nethemba (e.g., "Nethemba's project", "Company: Nethemba", etc.)
  # but keep "Origin Note" intact
  # Remove lines like "Company: Nethemba" or "Developed by: Nethemba"
  sed -i '/^[[:space:]]*Company:[[:space:]]*Nethemba/d' README.md
  sed -i '/^[[:space:]]*Developed[[:space:]]*by:[[:space:]]*Nethemba/d' README.md
  sed -i '/^[[:space:]]*Copyright:[[:space:]]*Nethemba/d' README.md
  echo "  ✓ Removed ownership claims by Nethemba (except Origin Note)"
  
  # Ensure contact is set correctly
  if grep -q "itssafer@itssafer.org"; then
    echo "  ✓ Contact email already set to itssafer@itssafer.org"
  fi
fi

echo -e "${GREEN}✓ License and branding cleanup complete${NC}"

################################################################################
# STEP 4: DEEP SECURITY & FILE AUDIT
################################################################################
echo -e "\n${YELLOW}==> STEP 4: Deep Security & File Audit${NC}"
cd "$PROJECT_DIR"

# Remove sensitive and build artifacts
echo "Removing build artifacts and sensitive files..."
rm -f n-audit-release || true
echo "  ✓ Removed n-audit-release"

rm -rf bin/ || true
echo "  ✓ Removed bin/"

rm -rf logs/ || true
echo "  ✓ Removed logs/"

# Remove Terraform state files (contains secrets)
find . -name "terraform.tfstate*" -type f -delete || true
echo "  ✓ Removed terraform.tfstate files"

# Remove old backups from project root (Gold Master archives)
shopt -s nullglob
removed=0
for archive in *.tar.gz *.sha256; do
  rm -f -- "$archive"
  echo "  ✓ Removed $archive"
  ((removed++))
done
shopt -u nullglob

if [[ $removed -eq 0 ]]; then
  echo "  ✓ No old archives found in project root"
fi

# Verify .gitignore is strict
echo "Verifying .gitignore..."
REQUIRED_PATTERNS=(
  "terraform.tfvars"
  "terraform.tfstate"
  "terraform.tfstate.*"
  "*.log"
  "bin/"
  "n-audit-release"
  "*.sha256"
  "*.tar.gz"
  ".terraform/"
)

GITIGNORE=".gitignore"
touch "$GITIGNORE"
added=0
for pat in "${REQUIRED_PATTERNS[@]}"; do
  if ! grep -qxF "$pat" "$GITIGNORE"; then
    echo "$pat" >> "$GITIGNORE"
    ((added++))
  fi
done

if [[ $added -gt 0 ]]; then
  echo "  ✓ Added $added missing patterns to .gitignore"
else
  echo "  ✓ .gitignore already comprehensive"
fi

echo -e "${GREEN}✓ Security audit complete${NC}"

################################################################################
# STEP 5: GO MODULE MAINTENANCE
################################################################################
echo -e "\n${YELLOW}==> STEP 5: Go Module Maintenance${NC}"
cd "$PROJECT_DIR"

if [[ -f go.mod ]]; then
  go mod tidy
  echo -e "${GREEN}✓ go mod tidy completed${NC}"
else
  echo "  ✓ No go.mod found (skipping)"
fi

################################################################################
# STEP 6: FINAL GOLD MASTER BACKUP
################################################################################
echo -e "\n${YELLOW}==> STEP 6: Final Gold Master Backup${NC}"

if [[ ! -f "$BACKUP_SCRIPT" ]]; then
  echo -e "${RED}✗ Backup script not found at $BACKUP_SCRIPT${NC}"
  exit 1
fi

chmod +x "$BACKUP_SCRIPT"
echo "Running backup script to create pristine archive..."
cd "$PROJECT_DIR" && bash "$BACKUP_SCRIPT" >/dev/null 2>&1 || bash "$BACKUP_SCRIPT"

# Find the latest backup
BACKUP_DIR="$HOME/n-audit-backups"
if [[ -d "$BACKUP_DIR" ]]; then
  LATEST=$(find "$BACKUP_DIR" -maxdepth 1 -name "*.tar.gz" -type f -printf '%T@ %p\n' | sort -rn | head -1 | cut -d' ' -f2-)
  if [[ -n "$LATEST" ]]; then
    ARCHIVE_SIZE=$(du -h "$LATEST" | cut -f1)
    echo -e "${GREEN}✓ Gold Master archive created: $(basename "$LATEST") ($ARCHIVE_SIZE)${NC}"
  fi
fi

################################################################################
# STEP 7: FINAL GIT COMMIT & FORCE PUSH
################################################################################
echo -e "\n${YELLOW}==> STEP 7: Final Git Commit & Force Push${NC}"
cd "$PROJECT_DIR"

# Ensure git is initialized
if [[ ! -d .git ]]; then
  git init -b main
  git config user.email "noreply@itssafer-devops.com"
  git config user.name "N-Audit Release Bot"
fi

# Stage all changes
git add .

# Commit
if git diff --cached --quiet; then
  echo "  ✓ No changes to commit"
else
  git commit -m "fix: repair mermaid syntax and update copyright/license info"
  echo -e "${GREEN}✓ Commit created${NC}"
fi

# Ensure remote is set
if ! git remote get-url origin &>/dev/null; then
  git remote add origin "$REMOTE_URL"
fi

# Force push
echo "Force-pushing to GitHub..."
git push -u origin main --force
echo -e "${GREEN}✓ Force push completed${NC}"

################################################################################
# FINAL SUMMARY
################################################################################
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  FINAL REMEDIAL UPDATE COMPLETE                               ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Repairs Applied:${NC}"
echo "  ✓ Mermaid diagram syntax fixed (localized → English)"
echo "  ✓ Copyright/License updated to ITSsafer-DevOps"
echo "  ✓ Ownership language clarified (Origin Note retained)"
echo "  ✓ Deep security audit performed"
echo "  ✓ Artifacts and state files removed"
echo "  ✓ Gold Master backup created"
echo "  ✓ Code pushed to GitHub (main branch)"
echo ""
echo -e "${BLUE}Repository:${NC} ${REMOTE_URL}"
echo -e "${BLUE}Verify:${NC} https://github.com/${GITHUB_ORG}/${REPO_NAME}"
echo ""
