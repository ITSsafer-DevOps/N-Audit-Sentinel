#!/usr/bin/env bash
################################################################################
# N-Audit Sentinel - Release Preparation & GitHub Publication Script
# 
# Purpose: Automates the entire beta release workflow including security audit,
#          artifact cleanup, Go module synchronization, Gold Master archival,
#          and automatic GitHub repository creation & push.
#
# Usage: ./prepare-and-publish.sh
#
# Requirements: bash 4+, git, gh (GitHub CLI), sudo privileges for apt-get
################################################################################

set -euo pipefail

# Color codes for output formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

################################################################################
# CONFIGURATION VARIABLES
################################################################################

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GITHUB_ORG="ITSsafer-DevOps"
REPO_NAME="N-Audit-Sentinel"
GITHUB_REPO_URL="https://github.com/${GITHUB_ORG}/${REPO_NAME}"
BACKUP_SCRIPT="${PROJECT_DIR}/scripts/backup-project.sh"

# Detect package manager based on OS
detect_package_manager() {
    if command -v dnf &> /dev/null; then
        echo "dnf"
    elif command -v yum &> /dev/null; then
        echo "yum"
    elif command -v apt-get &> /dev/null; then
        echo "apt-get"
    else
        echo "unknown"
    fi
}

PKG_MANAGER=$(detect_package_manager)

################################################################################
# UTILITY FUNCTIONS
################################################################################

# Print colored output with consistent formatting
print_header() {
    echo -e "\n${BLUE}===============================================================================${NC}"
    echo -e "${BLUE}→ $1${NC}"
    echo -e "${BLUE}===============================================================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Exit with error message
exit_error() {
    print_error "$1"
    exit 1
}

################################################################################
# STEP 1: PREREQUISITES CHECK & INSTALLATION
################################################################################

check_prerequisites() {
    print_header "STEP 1: Checking Prerequisites & Installing Missing Tools"

    # Detect OS and package manager
    print_info "Detected package manager: $PKG_MANAGER"
    
    if [[ "$PKG_MANAGER" == "unknown" ]]; then
        exit_error "Unable to detect package manager. Supported: dnf, yum, apt-get"
    fi

    # Helper function to install packages based on package manager
    install_package() {
        local package=$1
        case "$PKG_MANAGER" in
            dnf|yum)
                sudo "$PKG_MANAGER" install -y "$package" > /dev/null 2>&1
                ;;
            apt-get)
                sudo apt-get update > /dev/null 2>&1
                sudo apt-get install -y "$package" > /dev/null 2>&1
                ;;
        esac
    }

    # Check and install git
    if ! command -v git &> /dev/null; then
        print_warning "git is not installed. Installing via $PKG_MANAGER..."
        install_package "git"
        print_success "git installed successfully"
    else
        print_success "git is installed"
    fi

    # Check and install gh (GitHub CLI)
    if ! command -v gh &> /dev/null; then
        print_warning "GitHub CLI (gh) is not installed. Installing via $PKG_MANAGER..."
        install_package "gh"
        print_success "GitHub CLI (gh) installed successfully"
    else
        print_success "GitHub CLI (gh) is installed"
    fi

    # Check for Go (required for go mod tidy)
    if ! command -v go &> /dev/null; then
        exit_error "Go is not installed. Please install Go 1.16+ before running this script"
    fi
    print_success "Go is installed"

    # Handle GitHub CLI authentication
    print_info "Verifying GitHub CLI authentication..."
    if ! gh auth status > /dev/null 2>&1; then
        print_warning "Not authenticated with GitHub CLI"
        echo ""
        echo -e "${YELLOW}Starting interactive GitHub authentication...${NC}"
        echo "You will be guided through the authentication process."
        echo ""
        
        # Run gh auth login interactively
        gh auth login
        
        # Verify authentication after login
        if ! gh auth status > /dev/null 2>&1; then
            exit_error "GitHub authentication failed. Please try again."
        fi
        print_success "GitHub CLI authentication successful"
    else
        print_success "GitHub CLI authentication verified"
    fi
}

################################################################################
# STEP 2: SECURITY AUDIT (.gitignore)
################################################################################

audit_gitignore() {
    print_header "STEP 2: Security Audit (.gitignore)"

    local gitignore_path="${PROJECT_DIR}/.gitignore"

    # Check if .gitignore exists
    if [[ ! -f "$gitignore_path" ]]; then
        print_warning ".gitignore not found. Creating minimal .gitignore..."
        touch "$gitignore_path"
    else
        print_success ".gitignore exists"
    fi

    # Critical security check: terraform.tfvars
    if ! grep -q "terraform.tfvars" "$gitignore_path"; then
        print_warning "terraform.tfvars is NOT in .gitignore. Adding for security..."
        echo "terraform.tfvars" >> "$gitignore_path"
        print_success "terraform.tfvars added to .gitignore"
    else
        print_success "terraform.tfvars is properly ignored"
    fi

    # Ensure build artifacts are ignored
    local required_ignores=(
        "bin/"
        "n-audit-release"
        "*.log"
        "*.tar.gz"
        "*.sha256"
    )

    for pattern in "${required_ignores[@]}"; do
        if ! grep -q "^${pattern}$" "$gitignore_path"; then
            print_warning "Adding missing ignore pattern: $pattern"
            echo "$pattern" >> "$gitignore_path"
        fi
    done

    print_success "Security audit of .gitignore completed"
}

################################################################################
# STEP 3: DEEP CLEANING (ARTIFACT REMOVAL)
################################################################################

deep_clean() {
    print_header "STEP 3: Deep Cleaning - Removing Artifacts & Dev Artifacts"

    # Remove compiled binary
    if [[ -f "${PROJECT_DIR}/n-audit-release" ]]; then
        print_warning "Removing compiled binary: n-audit-release"
        rm -f "${PROJECT_DIR}/n-audit-release"
        print_success "Binary removed"
    else
        print_info "No binary found to remove"
    fi

    # Remove release artifacts
    local artifacts_removed=0
    while IFS= read -r artifact; do
        if [[ -f "$artifact" ]]; then
            print_warning "Removing artifact: $(basename "$artifact")"
            rm -f "$artifact"
            ((artifacts_removed++))
        fi
    done < <(find "${PROJECT_DIR}" -maxdepth 1 -type f \( -name "*.tar.gz" -o -name "*.sha256" \))

    if [[ $artifacts_removed -gt 0 ]]; then
        print_success "Removed $artifacts_removed artifact(s)"
    else
        print_info "No additional artifacts found to remove"
    fi

    # Recursively delete logs directory
    if [[ -d "${PROJECT_DIR}/logs" ]]; then
        print_warning "Removing logs directory (dev/test artifacts)..."
        rm -rf "${PROJECT_DIR}/logs"
        print_success "logs/ directory removed"
    else
        print_info "No logs directory found"
    fi

    # Clean bin directory
    if [[ -d "${PROJECT_DIR}/bin" ]]; then
        print_warning "Cleaning bin/ directory..."
        rm -rf "${PROJECT_DIR}/bin"
        print_success "bin/ directory cleaned"
    else
        print_info "No bin/ directory found"
    fi

    print_success "Deep cleaning completed"
}

################################################################################
# STEP 4: CODEBASE MAINTENANCE
################################################################################

maintain_codebase() {
    print_header "STEP 4: Codebase Maintenance - Synchronizing Go Modules"

    # Check if go.mod exists
    if [[ ! -f "${PROJECT_DIR}/go.mod" ]]; then
        print_warning "go.mod not found. Initializing Go module..."
        cd "${PROJECT_DIR}" && go mod init github.com/ITSsafer-DevOps/N-Audit-Sentinel
        print_success "Go module initialized"
    fi

    # Execute go mod tidy
    print_info "Running 'go mod tidy' to synchronize dependencies..."
    cd "${PROJECT_DIR}" && go mod tidy
    print_success "Go modules synchronized and cleaned"
}

################################################################################
# STEP 5: GOLD MASTER ARCHIVAL
################################################################################

create_gold_master() {
    print_header "STEP 5: Gold Master Archival - Creating Pristine Archive"

    # Verify backup script exists
    if [[ ! -f "$BACKUP_SCRIPT" ]]; then
        exit_error "Backup script not found at: $BACKUP_SCRIPT"
    fi

    # Check if backup script is executable
    if [[ ! -x "$BACKUP_SCRIPT" ]]; then
        print_warning "Making backup script executable..."
        chmod +x "$BACKUP_SCRIPT"
    fi

    # Execute the backup script
    print_info "Executing backup script: $BACKUP_SCRIPT"
    cd "${PROJECT_DIR}" && bash "$BACKUP_SCRIPT"

    # Find the backup location
    local backup_dir="$HOME/n-audit-backups"
    if [[ ! -d "$backup_dir" ]]; then
        exit_error "Backup directory not created at: $backup_dir"
    fi

    # Get the latest backup archive
    local latest_backup=$(find "$backup_dir" -maxdepth 1 -name "*.tar.gz" -type f -printf '%T@ %p\n' | sort -rn | head -1 | cut -d' ' -f2-)
    
    if [[ -z "$latest_backup" ]] || [[ ! -f "$latest_backup" ]]; then
        exit_error "No backup archive found in: $backup_dir"
    fi

    print_success "Gold Master archive created successfully"
    export GOLD_MASTER_PATH="$latest_backup"
}

################################################################################
# STEP 6: GIT INITIALIZATION & REPOSITORY SETUP
################################################################################

setup_git_repository() {
    print_header "STEP 6: Git Initialization & Repository Setup"

    cd "${PROJECT_DIR}"

    # Initialize git if not already initialized
    if [[ ! -d .git ]]; then
        print_warning "Git repository not initialized. Initializing..."
        git init
        git config user.email "${GIT_EMAIL:-noreply@itssafer-devops.com}"
        git config user.name "${GIT_USER:-N-Audit Bot}"
        print_success "Git repository initialized"
    else
        print_success "Git repository already initialized"
    fi

    # Configure git to handle file permissions
    git config core.fileMode true

    print_success "Git configuration completed"
}

################################################################################
# STEP 7: REPOSITORY CREATION & PUSH
################################################################################

create_and_push_repository() {
    print_header "STEP 7: Repository Creation & GitHub Push"

    cd "${PROJECT_DIR}"

    # Check if repository already exists on GitHub
    if gh repo view "${GITHUB_ORG}/${REPO_NAME}" 2>/dev/null; then
        print_warning "Repository ${GITHUB_ORG}/${REPO_NAME} already exists on GitHub"
        print_info "Ensuring remote is properly configured..."
    else
        print_info "Repository does not exist. Creating: ${GITHUB_ORG}/${REPO_NAME}"
        
        # Create repository with gh cli
        gh repo create "${GITHUB_ORG}/${REPO_NAME}" \
            --public \
            --source=. \
            --remote=origin \
            --description "N-Audit Sentinel - Network Audit & Security Monitoring Tool" \
            || print_warning "Repository creation returned a status (may already exist)"
    fi

    # Verify and configure remote
    if git remote get-url origin &> /dev/null; then
        print_info "Remote 'origin' already configured"
        print_info "Remote URL: $(git remote get-url origin)"
    else
        print_warning "Adding remote origin..."
        git remote add origin "https://github.com/${GITHUB_ORG}/${REPO_NAME}.git"
        print_success "Remote origin configured"
    fi

    # Stage all files
    print_info "Staging all project files..."
    git add .
    print_success "Files staged for commit"

    # Check if there are changes to commit
    if git diff --cached --quiet; then
        print_warning "No changes to commit. Repository already up-to-date"
    else
        # Commit with descriptive message
        print_info "Creating initial commit..."
        git commit -m "feat: initial beta release of N-Audit Sentinel" \
                   -m "- Complete source code for N-Audit Sentinel v1.0.0-beta
- Network security auditing and monitoring tool
- Includes CLI, release manager, and sentinel deployment
- Ready for open-source beta testing"
        print_success "Initial commit created"
    fi

    # Push to main branch
    print_info "Pushing to GitHub (main branch)..."
    
    # Check if main branch exists locally
    if ! git rev-parse --verify main &> /dev/null; then
        print_warning "main branch doesn't exist locally. Creating from current branch..."
        git branch -M main
    fi

    # Push with force-with-lease for safety
    git push -u origin main || {
        # Fallback to force push if needed
        print_warning "Standard push failed. Attempting force push..."
        git push -u origin main --force-with-lease
    }

    print_success "Repository pushed to GitHub successfully"
}

################################################################################
# STEP 8: FINAL SUMMARY
################################################################################

final_summary() {
    print_header "STEP 8: Release Preparation Complete - Summary"

    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✓ RELEASE PREPARATION SUCCESSFUL${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    echo -e "${BLUE}Repository Details:${NC}"
    echo -e "  Organization:  ${GITHUB_ORG}"
    echo -e "  Repository:     ${REPO_NAME}"
    echo -e "  GitHub URL:     ${GITHUB_REPO_URL}"
    echo ""

    if [[ -n "${GOLD_MASTER_PATH:-}" ]]; then
        echo -e "${BLUE}Gold Master Archive:${NC}"
        echo -e "  Location:       ${GOLD_MASTER_PATH}"
        echo -e "  Size:           $(du -h "${GOLD_MASTER_PATH}" | cut -f1)"
        echo ""
    fi

    echo -e "${BLUE}Completed Actions:${NC}"
    echo -e "  ✓ Prerequisites verified (git, gh, Go)"
    echo -e "  ✓ GitHub CLI authentication confirmed"
    echo -e "  ✓ Security audit (.gitignore) completed"
    echo -e "  ✓ Project artifacts cleaned"
    echo -e "  ✓ Go modules synchronized (go mod tidy)"
    echo -e "  ✓ Gold Master archive created"
    echo -e "  ✓ Git repository initialized/configured"
    echo -e "  ✓ Repository created on GitHub"
    echo -e "  ✓ Code pushed to main branch"
    echo ""

    echo -e "${BLUE}Next Steps:${NC}"
    echo -e "  1. Verify repository: ${GITHUB_REPO_URL}"
    echo -e "  2. Create GitHub release from the latest commit"
    echo -e "  3. Add release notes and attach Gold Master archive"
    echo -e "  4. Configure branch protection rules for main"
    echo -e "  5. Add CODEOWNERS and contribution guidelines"
    echo ""

    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

################################################################################
# ERROR HANDLER
################################################################################

error_handler() {
    local line_number=$1
    print_error "Script failed at line $line_number"
    echo ""
    echo "Troubleshooting tips:"
    echo "  • Ensure you have GitHub CLI installed and authenticated (gh auth login)"
    echo "  • Verify internet connectivity"
    echo "  • Check write permissions in the project directory"
    echo "  • Review error message above for specific issues"
    exit 1
}

trap 'error_handler ${LINENO}' ERR

################################################################################
# MAIN EXECUTION FLOW
################################################################################

main() {
    clear
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════════════╗"
    echo "║                   N-AUDIT SENTINEL RELEASE PREP                        ║"
    echo "║              Automated Beta Release & GitHub Publication               ║"
    echo "╚════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    check_prerequisites
    audit_gitignore
    deep_clean
    maintain_codebase
    create_gold_master
    setup_git_repository
    create_and_push_repository
    final_summary
}

# Execute main function
main "$@"
