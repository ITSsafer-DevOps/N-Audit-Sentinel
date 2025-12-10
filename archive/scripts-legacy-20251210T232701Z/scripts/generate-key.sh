#!/usr/bin/env bash
set -euo pipefail
KEY_DIR="/var/lib/n-audit/signing"
KEY_FILE="$KEY_DIR/id_ed25519"

if [ -f "$KEY_FILE" ]; then
  echo "[keygen] Existing key found at $KEY_FILE. Skipping generation." >&2
  exit 0
fi

echo "[keygen] Generating new ed25519 key at $KEY_FILE" >&2
mkdir -p "$KEY_DIR"
chmod 700 "$KEY_DIR"

ssh-keygen -t ed25519 -N '' -f "$KEY_FILE" >/dev/null 2>&1 || {
  echo "[keygen] ssh-keygen failed" >&2
  exit 1
}
chmod 600 "$KEY_FILE"
echo "[keygen] Key generation complete." >&2
