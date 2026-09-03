#!/bin/sh
# Canvas CLI installer for macOS and Linux.
#
#   curl -fsSL https://raw.githubusercontent.com/chiptoe-svg/canvas-cli/release/audited/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/chiptoe-svg/canvas-cli/release/audited/install.sh | sh -s -- --open-auth
#
# Downloads the audited release archive matching the OS and architecture,
# verifies its SHA-256 checksum, and installs the binary. Configure via:
#   CANVAS_VERSION=v1.13.0+audited.1   pin a version (the audited release is the default)
#   INSTALL_DIR=/usr/local/bin         install location (falls back to ~/.local/bin)
#   --open-auth                        macOS only: visibly open Terminal for optional token setup
set -eu

REPO="chiptoe-svg/canvas-cli"
BINARY="canvas"
VERSION="${CANVAS_VERSION:-v1.13.0+audited.2}"
OPEN_AUTH=0

die() { printf 'error: %s\n' "$1" >&2; exit 1; }

usage() {
  cat <<'EOF'
Usage: install.sh [--open-auth]

Installs the audited Canvas CLI release. Authentication is never started unless
--open-auth is supplied. On macOS that option opens a visible Terminal window
for the user to enter a Canvas access token locally; the token is not placed in
this installer, its arguments, a shell-history entry, or a persistent file.
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --open-auth) OPEN_AUTH=1 ;;
    -h|--help) usage; exit 0 ;;
    *) die "unknown option: $1 (use --help)" ;;
  esac
  shift
done

command -v curl >/dev/null 2>&1 || die "curl is required"
command -v tar >/dev/null 2>&1 || die "tar is required"
if command -v sha256sum >/dev/null 2>&1; then shacmd="sha256sum"
elif command -v shasum >/dev/null 2>&1; then shacmd="shasum -a 256"
else die "need sha256sum or shasum to verify the download"; fi

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux | darwin) ;;
  *) die "unsupported OS: $os (this installer covers Linux and macOS; use 'go install' otherwise)" ;;
esac
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64"; archre="(amd64|x86_64)" ;;
  arm64 | aarch64) arch="arm64"; archre="(arm64|aarch64)" ;;
  *) die "unsupported architecture: $arch" ;;
esac

expected_version="${VERSION#v}"
active_canvas="$(command -v "$BINARY" 2>/dev/null || true)"
if [ -n "$active_canvas" ] && "$active_canvas" version 2>/dev/null | grep -Fq "canvas-cli ${expected_version}"; then
  installed_canvas="$active_canvas"
  dir="$(dirname "$installed_canvas")"
  printf '%s %s is already active at %s; skipping installation\n' "$BINARY" "$VERSION" "$installed_canvas"
else
  base="https://github.com/${REPO}/releases/download/${VERSION}"
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT INT TERM

  curl -fsSL "${base}/checksums.txt" -o "${tmp}/checksums.txt" \
    || die "could not download checksums.txt for ${VERSION}"
  archive="$(awk '{print $2}' "${tmp}/checksums.txt" | grep -iE "_${os}_${archre}\.tar\.gz$" | head -1)"
  [ -n "$archive" ] || die "no ${os}/${arch} archive in ${VERSION}"

  curl -fsSL "${base}/${archive}" -o "${tmp}/${archive}" || die "download failed: ${archive}"
  want="$(awk -v f="$archive" '$2 == f {print $1}' "${tmp}/checksums.txt")"
  got="$(cd "$tmp" && $shacmd "$archive" | awk '{print $1}')"
  [ -n "$want" ] || die "no checksum recorded for ${archive}"
  [ "$want" = "$got" ] || die "checksum mismatch for ${archive}"

  tar -xzf "${tmp}/${archive}" -C "$tmp"
  [ -f "${tmp}/${BINARY}" ] || die "binary '${BINARY}' not found in the archive"
  chmod +x "${tmp}/${BINARY}"

  dir="${INSTALL_DIR:-/usr/local/bin}"
  if mkdir -p "$dir" 2>/dev/null && [ -w "$dir" ]; then
    mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
  elif command -v sudo >/dev/null 2>&1 && [ -t 0 ]; then
    printf 'installing to %s (needs sudo)\n' "$dir" >&2
    sudo mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
  else
    dir="${HOME}/.local/bin"
    mkdir -p "$dir"
    mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
  fi
  installed_canvas="${dir}/${BINARY}"
  printf '%s %s installed to %s\n' "$BINARY" "$VERSION" "$installed_canvas"
fi

case ":${PATH}:" in
  *":${dir}:"*) ;;
  *) printf 'note: %s is not on your PATH — add:\n  export PATH="%s:$PATH"\n' "$dir" "$dir" >&2 ;;
esac

if [ "$OPEN_AUTH" -eq 1 ]; then
  if [ "$os" != "darwin" ]; then
    printf '\nTo connect Canvas, run this command in your own terminal:\n  %s auth token set clemson --url https://clemson.instructure.com\n' "$installed_canvas"
    printf 'After it succeeds, confirm the connection with:\n  %s courses list --no-cache --columns id,name,course_code\n' "$installed_canvas"
    exit 0
  fi

  command -v open >/dev/null 2>&1 || die "macOS Launch Services command 'open' is unavailable"
  # macOS mktemp requires the Xs at the end of its template. Rename the
  # securely created file so Launch Services recognizes it as a .command file.
  auth_launcher="$(mktemp "${TMPDIR:-/tmp}/canvas-auth.XXXXXX")"
  mv "$auth_launcher" "${auth_launcher}.command"
  auth_launcher="${auth_launcher}.command"
  cat >"$auth_launcher" <<EOF
#!/bin/sh
trap 'rm -f "\$0"' EXIT HUP INT TERM
"$installed_canvas" auth token set clemson --url https://clemson.instructure.com
auth_status=\$?
if [ "\$auth_status" -eq 0 ]; then
  printf '\\nCanvas is connected. Checking your current classes...\\n\\n'
  "$installed_canvas" courses list --no-cache --columns id,name,course_code
  exit \$?
fi
printf '\\nCanvas connection was not completed. No class list was requested.\\n' >&2
exit "\$auth_status"
EOF
  chmod 700 "$auth_launcher"
  open "$auth_launcher"
  printf '\nA Terminal window has opened for the one-time Canvas connection. Enter the token there; it will not appear on screen. After a successful connection, that window will list your current classes and remove this temporary launcher.\n'
fi
