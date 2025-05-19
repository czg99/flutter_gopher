#!/bin/bash

cd $(dirname $0)

if ! command -v go &> /dev/null; then
    echo "Error: Go compiler not found. Please install Go."
    exit 1
fi

if ! command -v zig &> /dev/null; then
    echo "Error: Zig compiler not found. Please install Zig."
    exit 1
fi

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.so"
OUTPUT_DIR="$(pwd)"
GO_SRC="../src"
TIMESTAMP_FILE=".last_build_time"

# Check if source code has been updated
check_source_changes() {
    if [ ! -f "${TIMESTAMP_FILE}" ]; then
        return 0
    fi
    
    LAST_BUILD_TIME=$(cat "${TIMESTAMP_FILE}")
    NEWEST_FILE=$(find ${GO_SRC} -type f -name "*.go" -printf "%T@ %p\n" 2>/dev/null | sort -nr | head -1)
    if [ -z "${NEWEST_FILE}" ]; then
        return 0
    fi
    
    NEWEST_TIMESTAMP=$(echo ${NEWEST_FILE} | cut -d' ' -f1 | cut -d'.' -f1)
    if [ "${NEWEST_TIMESTAMP}" -gt "${LAST_BUILD_TIME}" ]; then
        return 0
    else
        if [ ! -f "${OUTPUT_DIR}/${OUTPUT_FILE}" ]; then
            return 0
        fi
        return 1
    fi
}

# Save current build timestamp
save_build_time() {
    date +%s > "${TIMESTAMP_FILE}"
}

# Check if source code has been updated
if ! check_source_changes; then
    echo "Source code unchanged, skipping compilation"
    exit 0
fi

echo "Detecting system architecture..."
GOARCH="amd64"
ZIG_TARGET="x86_64-linux-musl"

# Detecting ARM architecture
if [[ $(uname -m) == *"arm"* ]] || [[ $(uname -m) == *"aarch64"* ]]; then
    GOARCH="arm64"
    ZIG_TARGET="aarch64-linux-musl"
    echo "ARM architecture detected"
else
    echo "x86_64 architecture detected"
fi

export CGO_ENABLED=1
export GOOS="linux"
export GOARCH="${GOARCH}"

export CC="zig cc -target ${ZIG_TARGET}"
export CXX="zig c++ -target ${ZIG_TARGET}"

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

# Save current build timestamp
save_build_time

echo "Build process completed"