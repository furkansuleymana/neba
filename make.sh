#!/bin/bash

set -e

BUILD_DIR="./build"
BINARY_PATH="./build/neba.exe"
MAIN_PACKAGE_PATH="."

env=(
    "CGO_ENABLED=0"
    "GOOS=windows"
    "GOARCH=amd64"
)

function clean() {
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
    fi
}

function tidy() {
    go mod tidy -v && \
    go vet ./... && \
    go fmt ./...
}

function run() {
    clean
    go build -o $BINARY_PATH $MAIN_PACKAGE_PATH
    ./"$BINARY_PATH"
}

function prod() {
    clean
    tidy
    # TODO: cross-platform build with ldflags
    # TODO: auto versioning
}

task_name="${1:-run}"
tasks=(
    "clean"
    "tidy"
    "run"
    "production"
)

if [[ ! " ${tasks[*]} " =~ " $task_name " ]]; then
    echo "Invalid task name: $task_name" >&2
    exit 1
fi

${task_name}
