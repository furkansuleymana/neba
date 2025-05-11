#!/bin/bash

set -e

# Define variables
ROOT_DIR="$(pwd)"
BUILD_DIR="${ROOT_DIR}/build"
BINARY_PATH="${BUILD_DIR}/neba"

# Build for production
prod() {
  clean
  tidy

  local ldflags="-w -s"
  local targets=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
  )

  echo "Building for production..."

  local failures=()

  for target in "${targets[@]}"; do
    IFS="/" read -r GOOS GOARCH <<<"${target}"

    local ext=$([ "$GOOS" == "windows" ] && echo ".exe" || echo "")
    local binary="${BINARY_PATH}_${GOOS}_${GOARCH}${ext}"

    if (
      export GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0
      go build -a -ldflags="$ldflags" -o "$binary" "$ROOT_DIR" 2>/dev/null
    ); then
      echo "✓ $(basename "$binary")"
    else
      echo "✗ Failed"
      failures+=("$target")
    fi
  done

  echo
  if ((${#failures[@]})); then
    echo "Build failed for: ${failures[*]}" >&2
    return 1
  else
    echo "All builds completed successfully!"
  fi
}

# Build and run project
dev() {
  go build -o "${BINARY_PATH}" "${ROOT_DIR}"
  cd "${BUILD_DIR}" && "${BINARY_PATH}"
}

# Tidy project
tidy() {
  go mod tidy -v
  go fmt ./...
}

# Clean build directories
clean() {
  rm -rf "$BUILD_DIR"
}

# Get task name (default to "dev") and
# execute task if valid, else print error
TASK_NAME="${1:-dev}"
if ! declare -F "$TASK_NAME" >/dev/null; then
  echo "Invalid task name: $TASK_NAME" >&2
  exit 1
fi
$TASK_NAME
