#!/bin/bash
# xօ թαթí installer — download the matching release binary and run `xopapi install`
set -euo pipefail

REPO="abini121621/vpn4net-slipgate"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[1;36m'
NC='\033[0m'

info()    { echo -e "${GREEN}[+]${NC} $1"; }
error()   { echo -e "${RED}[-]${NC} $1"; exit 1; }

# Check root
[[ $EUID -ne 0 ]] && error "This script must be run as root (sudo)"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *)       error "Unsupported architecture: $ARCH" ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
[[ "$OS" != "linux" ]] && error "SlipGate only supports Linux"

BINARY="xopapi-${OS}-${ARCH}"

# Override with: XOPAPI_RELEASE_TAG=v1.0.0 bash install.sh
RELEASE_TAG="${XOPAPI_RELEASE_TAG:-${SLIPGATE_RELEASE_TAG:-}}"
CHANNEL=""  # ← set to "dev" on dev branch, empty on main

if [[ -n "$RELEASE_TAG" ]]; then
    URL="https://github.com/${REPO}/releases/download/${RELEASE_TAG}/${BINARY}"
elif [[ "$CHANNEL" == "dev" ]]; then
    # Find the latest dev pre-release tag via GitHub API
    DEV_TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases" \
        | grep -o '"tag_name": *"[^"]*-dev"' | head -1 | grep -o '"[^"]*-dev"' | tr -d '"')
    if [[ -n "$DEV_TAG" ]]; then
        URL="https://github.com/${REPO}/releases/download/${DEV_TAG}/${BINARY}"
        info "Dev channel: using release ${DEV_TAG}"
    else
        URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"
        info "No dev release found, falling back to latest stable"
    fi
else
    URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"
fi

echo -e "${CYAN}"
echo " __  __       _        _                 "
echo " \\ \\/ /  ___ | |_ __ _| |__   ___        "
echo "  \\  /  / _ \\| __/ _  | '_ \\ / _ \\     "
echo "  /  \\ | (_) | || (_| | |_) |  __/       "
echo " /_/\\_\\ \\___/ \\__\\__,_|_.__/ \\___|       "
echo "             xօ թαթí                    "
echo -e "${NC}"

# Stop only the previous xopapi process; existing SlipGate services remain
# untouched until the operator explicitly migrates them.
killall xopapi 2>/dev/null || true
killall dnstt-server 2>/dev/null || true
killall slipstream-server 2>/dev/null || true

info "Downloading xopapi ($OS/$ARCH)..."
if command -v curl &>/dev/null; then
    curl -fsSL "$URL" -o "${INSTALL_DIR}/xopapi"
elif command -v wget &>/dev/null; then
    wget -qO "${INSTALL_DIR}/xopapi" "$URL"
else
    error "Neither curl nor wget found"
fi

chmod +x "${INSTALL_DIR}/xopapi"
ln -sfn xopapi "${INSTALL_DIR}/slipgate"

info "Running xopapi install..."
if ! "${INSTALL_DIR}/xopapi" install </dev/tty; then
    error "xopapi install failed — run 'sudo xopapi install' to retry"
fi

info "Done! Run 'sudo xopapi' to see the menu."
