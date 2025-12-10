#!/usr/bin/env bash
################################################################################
# N-Audit Sentinel - Final Release Script
# 
# Purpose: Complete release workflow with robust cleanup, Go module sync,
#          Gold Master archival, and GitHub push.
#
# Usage: ./finish-release.sh
#
# Requirements: RHEL, bash 4+, git, gh (authenticated), go
################################################################################

set -e

################################################################################
# COLOR CODES
################################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

################################################################################
# CONFIGURATION
################################################################################

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GITHUB_ORG="ITSsafer-DevOps"
REPO_NAME="N-Audit-Sentinel"
GITHUB_REPO_URL="https://github.com/${GITHUB_ORG}/${REPO_NAME}"

################################################################################
# UTILITY FUNCTIONS
################################################################################

print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}→ $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

exit_error() {
    print_error "$1"
    exit 1
}

################################################################################
# STEP 1: SYSTEM CHECK
################################################################################

step_system_check() {
    print_header "STEP 1: System Check"

    # Check if running on RHEL
    if [[ ! -f /etc/redhat-release ]]; then
        exit_error "This system does not appear to be RHEL. Aborting."
    fi
    print_success "Running on RHEL"

    # Check git
    if ! command -v git &> /dev/null; then
        print_info "Installing git..."
        sudo dnf install -y git > /dev/null 2>&1
    fi
    print_success "git is installed"

    # Check gh
    if ! command -v gh &> /dev/null; then
        print_info "Installing GitHub CLI..."
        sudo dnf install -y gh > /dev/null 2>&1
    fi
    print_success "GitHub CLI (gh) is installed"

    # Check go
    if ! command -v go &> /dev/null; then
        print_info "Installing Go..."
        sudo dnf install -y golang > /dev/null 2>&1
    fi
    print_success "Go is installed"

    # Check gh authentication
    if ! gh auth status > /dev/null 2>&1; then
        exit_error "GitHub CLI is not authenticated. Run 'gh auth login' first."
    fi
    print_success "GitHub CLI is authenticated"
}

################################################################################
# STEP 2: SECURITY HARDENING (.gitignore)
################################################################################

step_security_hardening() {
    print_header "STEP 2: Security Hardening (.gitignore)"

    local gitignore_path="${PROJECT_DIR}/.gitignore"

    # Create .gitignore if it doesn't exist
    if [[ ! -f "$gitignore_path" ]]; then
        print_info "Creating .gitignore..."
        touch "$gitignore_path"
    fi

    # List of patterns that must be in .gitignore
    local patterns=(
        "terraform.tfvars"
        "bin/"
        "logs/"
        "*.log"
        "n-audit-release"
        "*.sha256"
    )

    local added_count=0
    for pattern in "${patterns[@]}"; do
        if ! grep -q "^${pattern}$" "$gitignore_path" 2>/dev/null; then
            echo "$pattern" >> "$gitignore_path"
            print_info "Added to .gitignore: $pattern"
            ((added_count++))
        fi
    done

    if [[ $added_count -eq 0 ]]; then
        print_success ".gitignore already contains all required patterns"
    else
        print_success "Added $added_count pattern(s) to .gitignore"
    fi
}

################################################################################
# STEP 3: ROBUST CLEANUP
################################################################################

step_robust_cleanup() {
    print_header "STEP 3: Robust Cleanup - Removing Artifacts"

    cd "${PROJECT_DIR}"

    # List of items to remove
    local items_to_remove=(
        "n-audit-release"
        "bin"
        "logs"
    )

    local removed_count=0
    for item in "${items_to_remove[@]}"; do
        if [[ -e "$item" ]]; then
            print_info "Removing: $item"
            rm -rf "$item"
            ((removed_count++))
        fi
    done

    # Remove *.tar.gz and *.sha256 files in root
    if find . -maxdepth 1 -name "*.tar.gz" -o -name "*.sha256" 2>/dev/null | grep -q .; then
        print_info "Removing release artifacts (*.tar.gz, *.sha256)..."
        find . -maxdepth 1 \( -name "*.tar.gz" -o -name "*.sha256" \) -delete
        ((removed_count++))
    fi

    if [[ $removed_count -eq 0 ]]; then
        print_success "No cleanup needed - directory is already clean"
    else
        print_success "Cleaned up $removed_count item(s)"
    fi
}

################################################################################
# STEP 4: GO MODULE MAINTENANCE
################################################################################

step_go_maintenance() {
    print_header "STEP 4: Go Module Maintenance"

    cd "${PROJECT_DIR}"

    # Check if go.mod exists
    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found in project root"
        exit_error "Please initialize Go module first: go mod init github.com/ITSsafer-DevOps/N-Audit-Sentinel"
    fi

    print_info "Running 'go mod tidy'..."
    go mod tidy
    print_success "Go modules synchronized"
}

################################################################################
# STEP 5: GOLD MASTER ARCHIVAL
################################################################################

step_gold_master() {
    print_header "STEP 5: Gold Master Archival"

    local backup_script="${PROJECT_DIR}/scripts/backup-project.sh"

    if [[ ! -f "$backup_script" ]]; then
        exit_error "Backup script not found: $backup_script"
    fi

    if [[ ! -x "$backup_script" ]]; then
        chmod +x "$backup_script"
    fi

    print_info "Executing backup script..."
    cd "${PROJECT_DIR}" && bash "$backup_script"

    # Find the latest backup
    local backup_dir="$HOME/n-audit-backups"
    if [[ ! -d "$backup_dir" ]]; then
        exit_error "Backup directory not found: $backup_dir"
    fi

    local latest_backup=$(find "$backup_dir" -maxdepth 1 -name "*.tar.gz" -type f -printf '%T@ %p\n' | sort -rn | head -1 | cut -d' ' -f2-)

    if [[ -z "$latest_backup" ]] || [[ ! -f "$latest_backup" ]]; then
        exit_error "No backup archive found in: $backup_dir"
    fi

    export GOLD_MASTER_PATH="$latest_backup"
    print_success "Gold Master archive created: $GOLD_MASTER_PATH"
}

################################################################################
# STEP 6: GIT INITIALIZATION
################################################################################

step_git_init() {
    print_header "STEP 6: Git Initialization"

    cd "${PROJECT_DIR}"

    if [[ ! -d .git ]]; then
        print_info "Initializing git repository..."
        git init -b main
        git config user.email "noreply@itssafer-devops.com"
        git config user.name "N-Audit Release Bot"
        print_success "Git repository initialized"
    else
        print_success "Git repository already initialized"
    fi
}

################################################################################
# STEP 7: GIT COMMIT
################################################################################

step_git_commit() {
    print_header "STEP 7: Git Commit"

    cd "${PROJECT_DIR}"

    print_info "Staging all files..."
    git add .

    # Check if there are changes to commit
    if git diff --cached --quiet; then
        print_info "No changes to commit - repository already up-to-date"
        return 0
    fi

    print_info "Creating commit..."
    git commit -m "feat: initial beta release of N-Audit Sentinel" \
               -m "- Complete source code for N-Audit Sentinel v1.0.0-beta
- Network security auditing and monitoring tool
- Includes CLI, release manager, and sentinel deployment
- Ready for open-source beta testing"
    print_success "Commit created"
}

################################################################################
# STEP 8: GITHUB PUSH (SMART CREATION)
################################################################################

step_github_push() {
    print_header "STEP 8: GitHub Push - Smart Repository Creation"

    cd "${PROJECT_DIR}"

    # Attempt to create repository
    print_info "Attempting to create repository: ${GITHUB_ORG}/${REPO_NAME}"
    
    if gh repo create "${GITHUB_ORG}/${REPO_NAME}" \
        --public \
        --source=. \
        --remote=origin \
        --push 2>/dev/null; then
        print_success "Repository created and code pushed successfully"
        return 0
    fi

    # Fallback: Repository may already exist
    print_info "Repository may already exist or creation failed. Checking remote..."

    # Verify remote is configured
    if ! git remote get-url origin &> /dev/null; then
        print_info "Configuring remote origin..."
        git remote add origin "https://github.com/${GITHUB_ORG}/${REPO_NAME}.git"
    else
        print_info "Remote origin already configured: $(git remote get-url origin)"
    fi

    # Ensure main branch exists locally
    if ! git rev-parse --verify main &> /dev/null; then
        print_info "Creating main branch..."
        git branch -M main
    fi

    # Push to main
    print_info "Pushing to main branch..."
    git push -u origin main || {
        print_info "Standard push failed, attempting force-with-lease..."
        git push -u origin main --force-with-lease
    }

    print_success "Code pushed to GitHub successfully"
}

################################################################################
# FINAL SUMMARY
################################################################################

step_final_summary() {
    print_header "STEP 9: Release Complete - Summary"

    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✓ RELEASE WORKFLOW COMPLETED SUCCESSFULLY${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    echo -e "${BLUE}Repository Information:${NC}"
    echo -e "  Organization:      ${GITHUB_ORG}"
    echo -e "  Repository:         ${REPO_NAME}"
    echo -e "  GitHub URL:         ${GITHUB_REPO_URL}"
    echo ""

    if [[ -n "${GOLD_MASTER_PATH:-}" ]]; then
        echo -e "${BLUE}Gold Master Archive:${NC}"
        echo -e "  Location:           ${GOLD_MASTER_PATH}"
        echo -e "  Size:               $(du -h "${GOLD_MASTER_PATH}" | cut -f1)"
        echo ""
    fi

    echo -e "${BLUE}Completed Actions:${NC}"
    echo -e "  ✓ System check (RHEL, git, gh, go)"
    echo -e "  ✓ GitHub CLI authentication verified"
    echo -e "  ✓ Security hardening (.gitignore)"
    echo -e "  ✓ Robust cleanup (artifacts removed)"
    echo -e "  ✓ Go modules synchronized"
    echo -e "  ✓ Gold Master archive created"
    echo -e "  ✓ Git repository initialized"
    echo -e "  ✓ Initial commit created"
    echo -e "  ✓ Code pushed to GitHub"
    echo ""

    echo -e "${BLUE}Next Steps:${NC}"
    echo -e "  1. Verify repository: ${GITHUB_REPO_URL}"
    echo -e "  2. Create GitHub release from main branch"
    echo -e "  3. Attach Gold Master archive to release"
    echo -e "  4. Configure branch protection rules"
    echo -e "  5. Add CODEOWNERS file"
    echo ""

    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

################################################################################
# MAIN EXECUTION
################################################################################

main() {
    clear
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════════╗"
    echo "║           N-AUDIT SENTINEL - FINAL RELEASE WORKFLOW             ║"
    echo "║                RHEL-Optimized Release Preparation               ║"
    echo "╚══════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"

    step_system_check
    step_security_hardening
    step_robust_cleanup
    step_go_maintenance
    step_gold_master
    step_git_init
    step_git_commit
    step_github_push
    step_final_summary
}

# Execute
main "$@"
