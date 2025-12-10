#!/usr/bin/env bash
set -e

# final-release-and-upload.sh
# Perform final manual release and upload of artifacts to GitHub

# Do NOT hardcode tokens. Read from environment variables: GH_TOKEN, GITHUB_TOKEN, or PAT.
# Export one of these before running, e.g. `export GH_TOKEN=...`.
PAT="${GH_TOKEN:-${GITHUB_TOKEN:-${PAT:-}}}"
TAG="v1.0.0-Beta"

echo "==> 1) Authenticating gh CLI with provided PAT"
echo "$PAT" | gh auth login --with-token

echo "==> 2) Fixing Mermaid localized keywords in Markdown files"
echo "Replacing: 'vývojový diagram' -> 'graph', 'podgraf' -> 'subgraph', 'koniec' -> 'end'"
find . -type f -name '*.md' -print0 | while IFS= read -r -d '' f; do
  echo "  - Processing $f"
  sed -i -e 's/vývojový diagram/graph/g' \
         -e 's/podgraf/subgraph/g' \
         -e 's/koniec/end/g' "$f" || true
done

# Remove any sed-created temp files (some sed implementations create backups like file.md.tmp)
echo "  - Cleaning possible sed temp files"
find . -type f -name '*.md.tmp' -print0 | xargs -0 -r rm -f || true

echo "==> 3) Deleting ALL GitHub releases and associated tags"
echo "Listing releases..."
releases=$(gh release list --limit 1000 --json tagName -q '.[].tagName' || true)
if [ -n "$releases" ]; then
  echo "$releases" | while IFS= read -r tag; do
    [ -z "$tag" ] && continue
    echo "  - Deleting release and tag: $tag"
    gh release delete "$tag" --confirm || true
    git push --delete origin "$tag" || true
    git tag -d "$tag" || true
  done
else
  echo "  - No releases found."
fi

echo "Also removing any remaining remote tags listed in git"
for t in $(git tag -l); do
  [ -z "$t" ] && continue
  echo "  - Removing remote tag: $t"
  git push --delete origin "$t" || true
  git tag -d "$t" || true
done

echo "==> 4) ULTRA-DEEP LOCAL CLEANUP"
echo "Removing old build artifacts, logs, and state files"
rm -rf n-audit-release bin logs *.tar.gz *.sha256 terraform.tfstate* || true
echo "Running go mod tidy"
go mod tidy

echo "==> 5) Local build & publish via ./beta-manual-release.sh"
if [ ! -x ./beta-manual-release.sh ]; then
  echo 'Error: ./beta-manual-release.sh not found or not executable' >&2
  exit 2
fi
./beta-manual-release.sh "$TAG"

echo "==> 6) Git publish: commit changes and force push to origin main"
git add . || true
if git diff --staged --quiet; then
  echo "No staged changes to commit"
else
  git commit -m "feat: final ultra-clean release (v1.0.0-Beta) - fixes and asset upload" || true
fi
echo "Force-pushing main to origin"
git push -u origin main --force

echo "==> 7) Create GitHub release and upload all local artifacts"
echo "Searching for artifacts: *.tar.gz, *.zip, *.sha256, *.tgz"
mapfile -t assets < <(find . -maxdepth 4 -type f \( -iname "*.tar.gz" -o -iname "*.zip" -o -iname "*.sha256" -o -iname "*.tgz" \) -print)

if [ ${#assets[@]} -eq 0 ]; then
  echo "No artifacts found to upload. Exiting."
  exit 3
fi

echo "Creating release $TAG and uploading assets"
printf 'Found assets to upload:\n'
for a in "${assets[@]}"; do
  printf '  %s\n' "$a"
done

gh release create "$TAG" "${assets[@]}" --title "$TAG" --notes "Final ultra-clean release $TAG"

echo "Release created and assets uploaded (if gh succeeded)."
echo "Done."
