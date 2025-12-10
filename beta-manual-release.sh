#!/usr/bin/env bash
set -euo pipefail

# beta-manual-release.sh
# Build release artifacts and create a GitHub Release (beta) with assets.

VERSION="${1:-v1.0.0-Beta}"
echo "[release] Version: $VERSION"

REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"

command -v gh >/dev/null 2>&1 || { echo "gh CLI not found in PATH"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "go not found in PATH"; exit 1; }

echo "[release] Checking GitHub auth..."
gh auth status || { echo "Please run 'gh auth login' first"; exit 1; }

# Build artifacts in a temporary directory to avoid leaving large files in the repo working tree
TMPDIR="$(mktemp -d)"
echo "[release] Building artifacts into $TMPDIR"
(cd "$REPO_ROOT" && go run "./cmd/n-audit-release" "$VERSION" && mv "n-audit-sentinel-${VERSION}-linux-amd64.tar.gz"* "$TMPDIR/")

TARBALL="$TMPDIR/n-audit-sentinel-${VERSION}-linux-amd64.tar.gz"
SIGFILE="$TARBALL.sha256"

[ -f "$TARBALL" ] || { echo "Artifact $TARBALL not found"; exit 1; }
[ -f "$SIGFILE" ] || { echo "Signature $SIGFILE not found"; exit 1; }

echo "[release] Creating/updating tag $VERSION"
git tag -f "$VERSION"
git push origin refs/tags/$VERSION --force

echo "[release] Creating GitHub release (prerelease) and uploading assets"
if ! gh release create "$VERSION" "$TARBALL" "$SIGFILE" --title "$VERSION" --notes "Beta release $VERSION" --prerelease; then
  echo "gh release create failed; attempting to update existing release"
  gh release upload "$VERSION" "$TARBALL" "$SIGFILE" --clobber
fi

echo "[release] Done. Release $VERSION should be available on GitHub (if account/actions allowed)."
