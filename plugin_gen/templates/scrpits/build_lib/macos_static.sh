#!/bin/bash

set -e

cd $(dirname $0)/../../

if ! command -v go &>/dev/null; then
	echo "Error: Go compiler not found. Please install Go."
	exit 1
fi

OUTPUT_NAME="{{.LibName}}MacOS"
OUTPUT_FILE="lib${OUTPUT_NAME}.a"
OUTPUT_DIR="${PWD}/darwin"
GO_SRC="gosrc"
MIN_VERSION=10.11

mkdir -p "$OUTPUT_DIR"

export CGO_ENABLED=1
export GOOS=darwin

LIB_FILES=""

for ARCH in "amd64:x86_64" "arm64:arm64"; do
	IFS=: read GOARCH CARCH <<<"$ARCH"

	export GOARCH=$GOARCH
	TARGET="$CARCH-apple-macos$MIN_VERSION"

	SDK_PATH=$(xcrun --sdk macosx --show-sdk-path)
	CLANG_PATH=$(xcrun --sdk macosx --find clang)

	export CC="$CLANG_PATH -target $TARGET -isysroot $SDK_PATH $@"
	export CXX="$CLANG_PATH++ -target $TARGET -isysroot $SDK_PATH $@"

	echo "Compiling $GOARCH architecture..."

	OUTPUT_FILE_TMP="${OUTPUT_DIR}/macos-${CARCH}/${OUTPUT_FILE}"

	LIB_FILES="$LIB_FILES \"$OUTPUT_FILE_TMP\""

	go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-archive -o "${OUTPUT_FILE_TMP}"

	if [ $? -ne 0 ]; then
		echo "Error: Go compilation failed, error code: $?"
		echo "Please check Go source code or compilation environment settings"
		exit $?
	else
		echo "Success: Go compilation completed, output file: ${OUTPUT_FILE_TMP}"
	fi
done

echo "Merging all architecture library files..."
lipo -create $LIB_FILES -output "${OUTPUT_DIR}/${OUTPUT_FILE}"

rm -rf "${OUTPUT_DIR}/macos-arm64"
rm -rf "${OUTPUT_DIR}/macos-x86_64"

echo "Creating XCFramework..."

rm -rf "${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"

xcodebuild -create-xcframework \
	-library "${OUTPUT_DIR}/${OUTPUT_FILE}" \
	-output "${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"

rm -rf "${OUTPUT_DIR}/${OUTPUT_FILE}"

echo "Created ${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"
