#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
#  gbx installer — install Git Branch Explorer from GitHub
#
#  Usage:
#    curl -fsSL https://raw.githubusercontent.com/Pixie2468/git-branch-explorer/main/install.sh | bash
#    curl -fsSL https://raw.githubusercontent.com/Pixie2468/git-branch-explorer/main/install.sh | bash -s -- --dir ~/.local/bin
# ──────────────────────────────────────────────────────────────
set -euo pipefail

REPO="Pixie2468/git-branch-explorer"
APP="gbx"
INSTALL_DIR="/usr/local/bin"

# ── Helpers ───────────────────────────────────────────────────

info()  { printf "\033[1;34m→\033[0m %s\n" "$*"; }
ok()    { printf "\033[1;32m✓\033[0m %s\n" "$*"; }
err()   { printf "\033[1;31m✗\033[0m %s\n" "$*" >&2; exit 1; }

need() {
  command -v "$1" >/dev/null 2>&1 || err "Required tool '$1' not found. Please install it."
}

# ── Parse flags ───────────────────────────────────────────────

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dir)  INSTALL_DIR="$2"; shift 2 ;;
    --help|-h)
      echo "Usage: install.sh [--dir <path>]"
      echo "  --dir    Installation directory (default: /usr/local/bin)"
      exit 0 ;;
    *) err "Unknown flag: $1" ;;
  esac
done

# ── Detect platform ──────────────────────────────────────────

need curl
need tar

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *)      err "Unsupported OS: $OS" ;;
esac

case "$ARCH" in
  x86_64|amd64)   ARCH="amd64" ;;
  aarch64|arm64)   ARCH="arm64" ;;
  *)               err "Unsupported architecture: $ARCH" ;;
esac

info "Detected platform: ${OS}/${ARCH}"

# ── Resolve latest release tag ────────────────────────────────

info "Fetching latest release..."
RELEASE_URL="https://api.github.com/repos/${REPO}/releases/latest"
TAG=$(curl -fsSL "$RELEASE_URL" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/')

if [[ -z "$TAG" ]]; then
  err "Could not determine latest release. Check https://github.com/${REPO}/releases"
fi

ok "Latest release: ${TAG}"

# ── Download binary ───────────────────────────────────────────

BINARY_NAME="${APP}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${BINARY_NAME}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

info "Downloading ${BINARY_NAME}..."
HTTP_CODE=$(curl -fsSL -w "%{http_code}" -o "${TMPDIR}/${APP}" "$DOWNLOAD_URL" 2>/dev/null || true)

if [[ "$HTTP_CODE" != "200" && ! -s "${TMPDIR}/${APP}" ]]; then
  # Fallback: try .tar.gz archive
  ARCHIVE_NAME="${BINARY_NAME}.tar.gz"
  ARCHIVE_URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE_NAME}"
  info "Trying archive format: ${ARCHIVE_NAME}..."
  curl -fsSL "$ARCHIVE_URL" -o "${TMPDIR}/${ARCHIVE_NAME}" || err "Download failed. Asset not found for ${OS}/${ARCH} at ${TAG}."
  tar -xzf "${TMPDIR}/${ARCHIVE_NAME}" -C "${TMPDIR}"
  # Find the binary inside the archive
  EXTRACTED=$(find "${TMPDIR}" -name "${APP}" -type f | head -1)
  if [[ -z "$EXTRACTED" ]]; then
    EXTRACTED=$(find "${TMPDIR}" -name "${APP}-*" -type f | head -1)
  fi
  if [[ -z "$EXTRACTED" ]]; then
    err "Could not find '${APP}' binary inside the archive."
  fi
  mv "$EXTRACTED" "${TMPDIR}/${APP}"
fi

chmod +x "${TMPDIR}/${APP}"

# ── Verify it runs ────────────────────────────────────────────

"${TMPDIR}/${APP}" --help >/dev/null 2>&1 || err "Downloaded binary failed to execute. Possible platform mismatch."

# ── Install ───────────────────────────────────────────────────

mkdir -p "$INSTALL_DIR"

if [[ -w "$INSTALL_DIR" ]]; then
  mv "${TMPDIR}/${APP}" "${INSTALL_DIR}/${APP}"
else
  info "Elevated permissions required for ${INSTALL_DIR}"
  sudo mv "${TMPDIR}/${APP}" "${INSTALL_DIR}/${APP}"
fi

ok "Installed ${APP} ${TAG} to ${INSTALL_DIR}/${APP}"

# ── Post-install check ────────────────────────────────────────

if command -v "$APP" >/dev/null 2>&1; then
  ok "Run '${APP} run' inside a git repo to start exploring!"
else
  echo ""
  echo "  '${APP}' is not in your PATH."
  echo "  Add this to your shell profile:"
  echo ""
  echo "    export PATH=\"${INSTALL_DIR}:\$PATH\""
  echo ""
fi
