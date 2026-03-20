#!/bin/bash

set -e

cd $(dirname $0)/../../

# Receive parameters: Android NDK home path and minimum API version
ANDROID_NDK_HOME=$1
MIN_API=$2

if [ -z "$ANDROID_NDK_HOME" ]; then
	echo "Error: Please provide NDK path as the first parameter"
	exit 1
fi

if [ -z "$MIN_API" ]; then
	MIN_API=21
fi

if ! command -v go &>/dev/null; then
	echo "Error: Go compiler not found. Please install Go."
	exit 1
fi

echo "ANDROID_NDK_HOME: ${ANDROID_NDK_HOME}"

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.so"
OUTPUT_HEADER="lib${OUTPUT_NAME}.h"
OUTPUT_DIR="${PWD}/android/libs"
GO_SRC="gosrc"

mkdir -p "$OUTPUT_DIR"

export CGO_ENABLED=1
export GOOS=android
export CGO_LDFLAGS="$CGO_LDFLAGS -Wl,-z,max-page-size=16384"

HOST_OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
	HOST_OS="darwin"
elif [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "cygwin"* ]]; then
	HOST_OS="windows"
fi

CC="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/${HOST_OS}-x86_64/bin/clang"
CXX="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/${HOST_OS}-x86_64/bin/clang++"

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

	export CC="$CC --target=$CC_TARGET"
	export CXX="$CXX --target=$CC_TARGET"

	go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-shared -o "${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"

	rm -rf "${OUTPUT_DIR}/${ARCH}/${OUTPUT_HEADER}"

	if [ $? -ne 0 ]; then
		echo "Error: Go compilation failed for architecture ${ARCH}"
		exit 1
	fi

	echo "${ARCH} compiled successfully: ${OUTPUT_DIR}/${ARCH}/${OUTPUT_FILE}"
done

echo "Build process completed"
