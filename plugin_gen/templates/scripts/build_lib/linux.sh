#!/bin/bash

set -e

cd $(dirname $0)/../../

if ! command -v go &>/dev/null; then
	echo "Error: Go compiler not found. Please install Go."
	exit 1
fi

if ! command -v zig &>/dev/null; then
	echo "Error: Zig compiler not found. Please install Zig."
	exit 1
fi

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.so"
OUTPUT_DIR="${PWD}/linux"
GO_SRC="gosrc"

mkdir -p "$OUTPUT_DIR"

echo "Detecting system architecture..."
GOARCH="amd64"
ZIG_TARGET="x86_64-linux-gnu"

# Detecting ARM architecture
if [[ $(uname -m) == *"arm"* ]] || [[ $(uname -m) == *"aarch64"* ]]; then
	GOARCH="arm64"
	ZIG_TARGET="aarch64-linux-gnu"
	echo "ARM architecture detected"
else
	echo "x86_64 architecture detected"
fi

export CGO_ENABLED=1
export GOOS="linux"
export GOARCH="$GOARCH"

export CC="zig cc -target $ZIG_TARGET"
export CXX="zig c++ -target $ZIG_TARGET"

echo "Compiling Go code to shared library..."

go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-shared -o "${OUTPUT_DIR}/${OUTPUT_FILE}"

if [ $? -ne 0 ]; then
	echo "Error: Go compilation failed, error code: $?"
	echo "Please check Go source code or compilation environment settings"
	exit $?
else
	echo "Success: Go compilation completed, output file: ${OUTPUT_DIR}/${OUTPUT_FILE}"
fi

echo "Cleaning header files..."
if [ -f "${OUTPUT_DIR}/${OUTPUT_NAME}.h" ]; then
	rm "${OUTPUT_DIR}/${OUTPUT_NAME}.h"
fi

echo "Build process completed"
