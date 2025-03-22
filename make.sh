#!/bin/bash

set -e

# Define variables
ROOT_DIR="$(pwd)"
BUILD_DIR_GO="${ROOT_DIR}/build"
BINARY_NAME="neba"

# Set environment variables
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
LDFLAGS="-s"

# Build project for production
build() {
  clean
  tidy
  go build -ldflags="${LDFLAGS}" -o "${BUILD_DIR_GO}/${BINARY_NAME}" "${ROOT_DIR}"
}

# Run project in development mode
dev() {
  go build -o "${BUILD_DIR_GO}/${BINARY_NAME}" "${ROOT_DIR}" && "${BUILD_DIR_GO}/${BINARY_NAME}"
}

# Tidy project
tidy() {
  go mod tidy -v
  go fmt ./...
}

# Clean build directories
clean() {
  rm -rf "$BUILD_DIR_GO"
}

# Get task name (default to "dev") and
# execute task if valid, else print error
TASK_NAME="${1:-dev}"
if ! declare -F "$TASK_NAME" >/dev/null; then
  echo "Invalid task name: $TASK_NAME" >&2
  exit 1
fi
$TASK_NAME
