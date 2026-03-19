#!/bin/bash

set -e

cd $(dirname $0)/../../

if ! command -v go &>/dev/null; then
	echo "Error: Go compiler not found. Please install Go."
	exit 1
fi

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.a"
OUTPUT_DIR="${PWD}/darwin"
GO_SRC="gosrc"
MIN_VERSION=11

mkdir -p "$OUTPUT_DIR"

export CGO_ENABLED=1
export GOOS=ios

SIMULATOR_LIBS=""
DEVICE_LIBS=""

for ARCH in "amd64:x86_64:iphonesimulator:simulator" "arm64:arm64:iphonesimulator:simulator" "arm64:arm64:iphoneos:"; do
	IFS=: read GOARCH CARCH SDK SIMULATOR <<<"$ARCH"

	export GOARCH=$GOARCH
	TARGET="$CARCH-apple-ios$MIN_VERSION"
	if [ -n "$SIMULATOR" ]; then
		TARGET="$TARGET-$SIMULATOR"
	fi

	SDK_PATH=$(xcrun --sdk "$SDK" --show-sdk-path)
	CLANG_PATH=$(xcrun --sdk "$SDK" --find clang)

	export CC="$CLANG_PATH -target $TARGET -isysroot $SDK_PATH $@"
	export CXX="$CLANG_PATH++ -target $TARGET -isysroot $SDK_PATH $@"

	if [ "$SDK" = "iphonesimulator" ]; then
		OUTPUT_FILE_TMP="${OUTPUT_DIR}/ios-simulator/${GOARCH}_${OUTPUT_FILE}"
		SIMULATOR_LIBS="$SIMULATOR_LIBS \"$OUTPUT_FILE_TMP\""
	else
		OUTPUT_FILE_TMP="${OUTPUT_DIR}/ios-arm64/${GOARCH}_${OUTPUT_FILE}"
		DEVICE_LIBS="$DEVICE_LIBS \"$OUTPUT_FILE_TMP\""
	fi

	go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-archive -o "${OUTPUT_FILE_TMP}"

	if [ $? -ne 0 ]; then
		echo "Error: Go compilation failed, error code: $?"
		echo "Please check Go source code or compilation environment settings"
		exit $?
	else
		echo "Success: Go compilation completed, output file: ${OUTPUT_FILE_TMP}"
	fi
done

echo "Creating fat libraries for simulator and device..."
if [ ! -z "$SIMULATOR_LIBS" ]; then
	lipo -create $SIMULATOR_LIBS -output "${OUTPUT_DIR}/ios-simulator/${OUTPUT_FILE}"
fi

if [ ! -z "$DEVICE_LIBS" ]; then
	lipo -create $DEVICE_LIBS -output "${OUTPUT_DIR}/ios-arm64/${OUTPUT_FILE}"
fi

echo "Creating XCFramework..."

rm -rf "${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"

xcodebuild -create-xcframework \
	-library "${OUTPUT_DIR}/ios-simulator/${OUTPUT_FILE}" \
	-library "${OUTPUT_DIR}/ios-arm64/${OUTPUT_FILE}" \
	-output "${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"

rm -rf "${OUTPUT_DIR}/ios-arm64"
rm -rf "${OUTPUT_DIR}/ios-simulator"

echo "Created ${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework"
