#!/bin/bash

set -e

cd $(dirname $0)/../../

# Receive parameters: OHOS NDK home path
OHOS_NDK_HOME=$1

if [ -z "$OHOS_NDK_HOME" ]; then
	echo "Error: Please provide OHOS NDK home path as the first parameter"
	exit 1
fi

if ! command -v go &>/dev/null; then
	echo "Error: Go compiler not found. Please install Go."
	exit 1
fi

echo "OHOS_NDK_HOME: ${OHOS_NDK_HOME}"

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.so"
OUTPUT_HEADER="lib${OUTPUT_NAME}.h"
OUTPUT_DIR="${PWD}/ohos/libs"
GO_SRC="gosrc"

mkdir -p "${OUTPUT_DIR}"

export CGO_ENABLED=1
export GOOS=openharmony

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

	go build -C $GO_SRC -ldflags "-s -w -extldflags '-Wl,-soname,${OUTPUT_FILE}'" -trimpath -buildmode=c-shared -o "${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"

	rm -rf "${OUTPUT_DIR}/${ARCH}/${OUTPUT_HEADER}"

	if [ $? -ne 0 ]; then
		echo "Error: Go compilation failed for architecture ${ARCH}"
		exit 1
	fi

	echo "${ARCH} compiled successfully: ${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"
done

echo "Build process completed"
