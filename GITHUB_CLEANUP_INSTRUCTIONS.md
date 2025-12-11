# GitHub Repository Cleanup Instructions

## Objective
Remove large release artifacts from Git history and use GitHub Releases API for distribution.

## Why?
- GitHub recommends files < 50MB in repositories
- Our binary artifacts (79M) exceed this limit
- GitHub Releases API is the proper way to distribute large files

## Steps to Clean Up GitHub Repository

### 1. Remove Large Files from Git History (BFG Recommended)

**Using BFG Repo-Cleaner (recommended):**

```bash
# Install bfg (https://rtyley.github.io/bfg-repo-cleaner/)
brew install bfg  # macOS
# or download from https://rtyley.github.io/bfg-repo-cleaner/

# Clone a fresh mirror
git clone --mirror https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel.git

# Remove all .tar.gz and .sha256 files > 50MB
bfg --delete-files 'n-audit-*.tar.gz'
bfg --delete-files '*.sha256'

# Clean and push
git reflog expire --expire=now --all
git gc --prune=now --aggressive
git push --mirror

# Verify
cd ..
git clone https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel
git log --pretty=oneline -20
```

### 2. Create GitHub Release (Web UI or CLI)

**Using GitHub CLI:**

```bash
gh release create v1.0.0-Beta \
  --title "N-Audit Sentinel v1.0.0-Beta" \
  --notes-file RELEASE_NOTES.md \
  ./releases/n-audit-sentinel-v1.0.0-Beta-bin.tar.gz \
  ./releases/n-audit-sentinel-v1.0.0-Beta-bin.tar.gz.sha256 \
  ./releases/n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz \
  ./releases/n-audit-sentinel-v1.0.0-Beta-goldmaster.tar.gz.sha256
```

**Or via GitHub Web UI:**
1. Go to `https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases`
2. Click "Create a new release"
3. Tag version: `v1.0.0-Beta`
4. Title: "N-Audit Sentinel v1.0.0-Beta"
5. Description: Copy from `RELEASE_NOTES.md`
6. Upload binary files (4 artifacts)
7. Publish release

### 3. Verify Repository is Clean

```bash
git log --all --pretty=format:"%h %s" | grep -i "release\|artifact\|backup"
git ls-files | grep -E '\.(tar\.gz|sha256)$'  # Should be empty
```

### 4. Document in README

Add to README.md Installation section:

```markdown
## Installation

### Download Pre-built Binaries

Visit [Releases](https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases) page:

- Download `n-audit-sentinel-v1.0.0-Beta-bin.tar.gz`
- Verify SHA256: `n-audit-sentinel-v1.0.0-Beta-bin.tar.gz.sha256`

### Or Build from Source

\`\`\`bash
git clone https://github.com/ITSsafer-DevOps/N-Audit-Sentinel.git
cd N-Audit-Sentinel
make build
make test
\`\`\`
```

## Alternative: Use Git LFS (Lighter Weight)

If you prefer to keep large files in Git:

```bash
# Install git-lfs
brew install git-lfs  # macOS

# Initialize LFS in repository
git lfs install

# Track large files
git lfs track "releases/*.tar.gz"
git lfs track "releases/*.sha256"

# Add to git
git add .gitattributes
git add releases/
git commit -m "track: large release artifacts via git-lfs"
git push
```

## Status Check

After cleanup:

```bash
# Should show only RELEASE_NOTES.md in releases/
git ls-files releases/

# Should be 0B or have only .gitkeep
ls -lh releases/
```

---

**Current Status:** Ready for GitHub cleanup
**Branch:** main
**Recommendation:** Use GitHub Releases API (preferred) over LFS
