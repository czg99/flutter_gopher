#!/bin/bash
set -e

# Receive parameters: OHOS NDK home path
OHOS_NDK_HOME=$1

if [ -z "$OHOS_NDK_HOME" ]; then
    echo "Error: Please provide OHOS NDK home path as the first parameter"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Error: Go compiler not found. Please install Go."
    exit 1
fi

echo "OHOS_NDK_HOME: ${OHOS_NDK_HOME}"

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
        for ARCH in "arm64-v8a" "armeabi-v7a" "x86_64"; do
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
export GOOS=android

CC="${OHOS_NDK_HOME}/native/llvm/bin/clang"
CXX="${OHOS_NDK_HOME}/native/llvm/bin/clang++"

ARCHS=("arm64-v8a" "armeabi-v7a" "x86_64")
for ARCH in "${ARCHS[@]}"; do
    echo "Compiling ${ARCH} architecture..."
    
    if [ "$ARCH" == "arm64-v8a" ]; then
        export GOARCH=arm64
        CC_TARGET="aarch64-linux-ohos"
        elif [ "$ARCH" == "armeabi-v7a" ]; then
        export GOARCH=arm
        export GOARM=7
        CC_TARGET="armv7-linux-ohos"
        elif [ "$ARCH" == "x86_64" ]; then
        export GOARCH=amd64
        CC_TARGET="x86_64-linux-ohos"
    fi
    
    mkdir -p "${OUTPUT_DIR}/${ARCH}"
    
    export CC="${CC} --target=${CC_TARGET} --sysroot=${OHOS_NDK_HOME}/native/sysroot"
    export CXX="${CXX} --target=${CC_TARGET} --sysroot=${OHOS_NDK_HOME}/native/sysroot"
    
    export CGO_CFLAGS="-I${PWD}/log-adaptor/include"
    export CGO_LDFLAGS="-L${PWD}/log-adaptor/dist/${ARCH}"

    cp -f ${PWD}/log-adaptor/dist/${ARCH}/* ${OUTPUT_DIR}/${ARCH}
    
    go build -C ${GO_SRC} -ldflags "-s -w -extldflags '-Wl,-soname,${OUTPUT_FILE}'" -trimpath -buildmode=c-shared -o "${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"
    
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