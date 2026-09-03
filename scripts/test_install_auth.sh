#!/bin/sh
# Network-free regression checks for the audited installer authentication flow.
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
installer="$repo_dir/install.sh"

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

sh -n "$installer"

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT HUP INT TERM
mock_bin="$tmp_dir/bin"
mkdir -p "$mock_bin"

cat >"$mock_bin/canvas" <<'EOF'
#!/bin/sh
if [ "${1:-}" = "version" ]; then
  printf 'canvas-cli 1.13.0+audited.1\n'
  exit 0
fi
printf 'unexpected canvas invocation: %s\n' "$*" >&2
exit 99
EOF

cat >"$mock_bin/curl" <<'EOF'
#!/bin/sh
printf 'curl must not run when the pinned binary is already active\n' >&2
exit 99
EOF

cat >"$mock_bin/uname" <<'EOF'
#!/bin/sh
case "${1:-}" in
  -s) printf 'Darwin\n' ;;
  -m) printf 'arm64\n' ;;
  *) exit 2 ;;
esac
EOF

cat >"$mock_bin/open" <<'EOF'
#!/bin/sh
[ "$#" -eq 1 ] || exit 2
cp "$1" "$OPEN_CAPTURE"
printf '%s\n' "$1" >"$OPEN_PATH"
EOF

chmod 700 "$mock_bin/canvas" "$mock_bin/curl" "$mock_bin/uname" "$mock_bin/open"

OPEN_CAPTURE="$tmp_dir/captured.command" \
OPEN_PATH="$tmp_dir/opened-path" \
TMPDIR="$tmp_dir" \
PATH="$mock_bin:$PATH" \
CANVAS_VERSION=v1.13.0+audited.1 \
sh "$installer" --open-auth >"$tmp_dir/output" 2>&1

grep -Fq 'already active' "$tmp_dir/output" || fail 'already-installed path was not used'
[ -s "$tmp_dir/opened-path" ] || fail 'Launch Services was not invoked'
[ -s "$tmp_dir/captured.command" ] || fail 'launcher was not created'
grep -Fq 'auth token set clemson --url https://clemson.instructure.com' "$tmp_dir/captured.command" \
  || fail 'launcher did not contain the token-free auth command'
grep -Fq 'courses list --no-cache --columns id,name,course_code' "$tmp_dir/captured.command" \
  || fail 'launcher did not contain the post-auth class-list check'
grep -Fq 'rm -f "$0"' "$tmp_dir/captured.command" \
  || fail 'launcher did not self-delete'
if grep -Fq -- '--token' "$tmp_dir/captured.command"; then
  fail 'launcher exposed a token argument'
fi

printf 'PASS: installer syntax, already-installed path, and macOS Terminal-launch path\n'
