#!/bin/bash
set -e

# Receive parameters: NDK path and minimum API version
NDK_PATH=$1
MIN_API=$2

if [ -z "$NDK_PATH" ]; then
    echo "Error: Please provide NDK path as the first parameter"
    exit 1
fi

if [ -z "$MIN_API" ]; then
    echo "Error: Please provide minimum API version as the second parameter"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Error: Go compiler not found. Please install Go."
    exit 1
fi

echo "NDK_PATH: ${NDK_PATH}"

cd $(dirname $0)

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.so"
OUTPUT_HEADER="lib${OUTPUT_NAME}.h"
OUTPUT_DIR="$(pwd)/libs"
GO_SRC="../gosrc"
TIMESTAMP_FILE=".last_build_time"

mkdir -p "${OUTPUT_DIR}"

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
        for ARCH in "arm64-v8a" "armeabi-v7a" "x86" "x86_64"; do
            if [ ! -f "${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}" ]; then
                return 0
            fi
        done
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

export CGO_ENABLED=1
export CGO_LDFLAGS="$CGO_LDFLAGS -Wl,-z,max-page-size=16384"
export GOOS=android

HOST_OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
    HOST_OS="darwin"
elif [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "cygwin"* ]]; then
    HOST_OS="windows"
fi

CC="${NDK_PATH}/toolchains/llvm/prebuilt/${HOST_OS}-x86_64/bin/clang"
CXX="${NDK_PATH}/toolchains/llvm/prebuilt/${HOST_OS}-x86_64/bin/clang++"

ARCHS=("arm64-v8a" "armeabi-v7a" "x86" "x86_64")
for ARCH in "${ARCHS[@]}"; do
    echo "Compiling ${ARCH} architecture..."
    
    if [ "$ARCH" == "arm64-v8a" ]; then
        export GOARCH=arm64
        CC_TARGET="aarch64-linux-android${MIN_API}"
    elif [ "$ARCH" == "armeabi-v7a" ]; then
        export GOARCH=arm
        export GOARM=7
        CC_TARGET="armv7a-linux-androideabi${MIN_API}"
    elif [ "$ARCH" == "x86" ]; then
        export GOARCH=386
        CC_TARGET="i686-linux-android${MIN_API}"
    elif [ "$ARCH" == "x86_64" ]; then
        export GOARCH=amd64
        CC_TARGET="x86_64-linux-android${MIN_API}"
    fi
    
    mkdir -p "${OUTPUT_DIR}/${ARCH}"
    
    export CC="${CC} --target=${CC_TARGET}"
    export CXX="${CXX} --target=${CC_TARGET}"
    
    go build -C ${GO_SRC} -ldflags "-s -w" -trimpath -buildmode=c-shared -o "${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"

    rm -rf "${OUTPUT_DIR}/${ARCH}/${OUTPUT_HEADER}"
    
    if [ $? -ne 0 ]; then
        echo "Error: Go compilation failed for architecture ${ARCH}"
        exit 1
    fi
    
    echo "${ARCH} compiled successfully: ${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"
done

# Save current build timestamp
save_build_time

echo "Build process completed"