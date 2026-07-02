#!/bin/sh
# agentrc installer — downloads the latest released `agentrc` CLI binary and
# installs it (with the `arc` alias) into a bin directory on your PATH.
#
#   curl -fsSL https://agentrc.ai/install.sh | sh
#
# Environment overrides:
#   AGENTRC_VERSION      pin a release tag (e.g. v0.1.1); default: latest
#   AGENTRC_INSTALL_DIR  target bin dir; default: /usr/local/bin (falls back to
#                        $HOME/.local/bin when /usr/local/bin is not writable)
#
# agentrc packages one AI agent as a portable, governed OCI artifact. The CLI
# builds and inspects Agentfiles and translates them to backend deploy configs
# (`arc run --backend <b> --dry-run`); it does not ship an agent runtime.
set -eu

REPO="adeelahmad/agentrc"
BINARY="agentrc"
ALIAS="arc"

info() { printf '\033[1;36m==>\033[0m %s\n' "$1"; }
warn() { printf '\033[1;33mwarning:\033[0m %s\n' "$1" >&2; }
err()  { printf '\033[1;31merror:\033[0m %s\n' "$1" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || err "curl is required"
command -v tar  >/dev/null 2>&1 || err "tar is required"

# --- detect platform -------------------------------------------------------
os=$(uname -s)
case "$os" in
  Darwin) os="darwin" ;;
  Linux)  os="linux" ;;
  *) err "unsupported OS: $os (agentrc ships darwin and linux binaries; use 'go install $REPO/cmd/agentrc@latest')" ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) err "unsupported architecture: $arch (use 'go install $REPO/cmd/agentrc@latest')" ;;
esac

# --- resolve version -------------------------------------------------------
version="${AGENTRC_VERSION:-}"
if [ -z "$version" ]; then
  info "Resolving latest release..."
  version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
  [ -n "$version" ] || err "could not resolve latest release tag (set AGENTRC_VERSION to pin one)"
fi
info "Installing agentrc $version ($os/$arch)"

asset="agentrc_${version}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$version"

# --- download + optional checksum verification -----------------------------
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
curl -fsSL "$base/$asset" -o "$tmp/$asset" || err "download failed: $base/$asset"

if curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt" 2>/dev/null; then
  expected=$(awk -v a="$asset" '$2 == a {print $1}' "$tmp/checksums.txt")
  if [ -z "$expected" ]; then
    warn "checksums.txt has no entry for $asset; skipping verification"
  else
    if command -v shasum >/dev/null 2>&1; then actual=$(shasum -a 256 "$tmp/$asset" | awk '{print $1}')
    elif command -v sha256sum >/dev/null 2>&1; then actual=$(sha256sum "$tmp/$asset" | awk '{print $1}')
    else actual=""; warn "no sha256 tool found; skipping checksum verification"; fi
    if [ -n "$actual" ] && [ "$actual" != "$expected" ]; then
      err "checksum mismatch for $asset (expected $expected, got $actual)"
    fi
    [ -n "$actual" ] && info "Checksum verified"
  fi
else
  warn "no checksums.txt published for $version; skipping verification"
fi

tar -xzf "$tmp/$asset" -C "$tmp" || err "extract failed"
[ -f "$tmp/$BINARY" ] || err "archive did not contain expected binary '$BINARY'"
chmod +x "$tmp/$BINARY"

# --- choose install dir ----------------------------------------------------
dir="${AGENTRC_INSTALL_DIR:-/usr/local/bin}"
if [ -n "${AGENTRC_INSTALL_DIR:-}" ]; then
  mkdir -p "$dir" 2>/dev/null || err "cannot create install dir: $dir"
  [ -w "$dir" ] || err "install dir not writable: $dir"
elif [ ! -d "$dir" ] || [ ! -w "$dir" ]; then
  home="${HOME:-}"
  [ -n "$home" ] || err "/usr/local/bin not writable and HOME is unset; set AGENTRC_INSTALL_DIR"
  dir="$home/.local/bin"
  mkdir -p "$dir" || err "cannot create install dir: $dir"
  warn "/usr/local/bin not writable; installing to $dir"
fi

# --- install binary + alias ------------------------------------------------
mv "$tmp/$BINARY" "$dir/$BINARY" || err "failed to install binary to $dir"
ln -sf "$dir/$BINARY" "$dir/$ALIAS" 2>/dev/null || cp "$dir/$BINARY" "$dir/$ALIAS"
info "Installed $dir/$BINARY (alias: $ALIAS)"

# --- PATH hint + verify ----------------------------------------------------
case ":$PATH:" in
  *":$dir:"*) ;;
  *) warn "$dir is not on your PATH — add it, e.g.:"
     # shellcheck disable=SC2016  # literal $PATH is intended in the printed advice
     printf '       export PATH="%s:$PATH"\n' "$dir" >&2 ;;
esac

if "$dir/$BINARY" version >/dev/null 2>&1; then
  info "$("$dir/$BINARY" version 2>/dev/null | head -1)"
fi
info "Done. Try: $ALIAS --help"
