#!/bin/sh

cd $(dirname $0)

if ! command -v go &> /dev/null; then
    echo "Error: Go compiler not found. Please install Go."
    exit 1
fi

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.a"
OUTPUT_DIR="$(pwd)"
GO_SRC="../src"
TIMESTAMP_FILE=".last_build_time"

MIN_VERSION=11

# Check if source code has been updated
check_source_changes() {
    if [ ! -f "${TIMESTAMP_FILE}" ]; then
        return 0
    fi
    
    LAST_BUILD_TIME=$(cat "${TIMESTAMP_FILE}")
    NEWEST_FILE=$(find ${GO_SRC} -type f -name "*.go" -exec sh -c 'for file; do echo "$(date -r "$file" +%s) $file"; done' _ {} + | sort -nr | head -1)
    if [ -z "${NEWEST_FILE}" ]; then
        return 0
    fi
    
    NEWEST_TIMESTAMP=$(echo ${NEWEST_FILE} | cut -d' ' -f1 | cut -d'.' -f1)
    if [ "${NEWEST_TIMESTAMP}" -gt "${LAST_BUILD_TIME}" ]; then
        return 0
    else
        if [ ! -d "${OUTPUT_DIR}/${OUTPUT_NAME}.xcframework" ]; then
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

export CGO_ENABLED=1
export GOOS=ios

SIMULATOR_LIBS=""
DEVICE_LIBS=""

for ARCH in "amd64:x86_64:iphonesimulator:simulator" "arm64:arm64:iphonesimulator:simulator" "arm64:arm64:iphoneos:"
do
    IFS=: read GOARCH CARCH SDK SIMULATOR <<< "$ARCH"
    
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
        OUTPUT_FILE_TMP="ios-simulator/${GOARCH}_${OUTPUT_FILE}"
        SIMULATOR_LIBS="$SIMULATOR_LIBS $OUTPUT_FILE_TMP"
    else
        OUTPUT_FILE_TMP="ios-arm64/${GOARCH}_${OUTPUT_FILE}"
        DEVICE_LIBS="$DEVICE_LIBS $OUTPUT_FILE_TMP"
    fi
    
    go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-archive -o "$OUTPUT_DIR/$OUTPUT_FILE_TMP"

    if [ $? -ne 0 ]; then
        echo "Error: Go compilation failed, error code: $?"
        echo "Please check Go source code or compilation environment settings"
        exit $?
    else
        echo "Success: Go compilation completed, output file: $OUTPUT_DIR/$OUTPUT_FILE_TMP"
    fi
done

echo "Creating fat libraries for simulator and device..."
if [ ! -z "$SIMULATOR_LIBS" ]; then
    lipo -create $SIMULATOR_LIBS -output ios-simulator/${OUTPUT_FILE}
fi

if [ ! -z "$DEVICE_LIBS" ]; then
    lipo -create $DEVICE_LIBS -output ios-arm64/${OUTPUT_FILE}
fi

echo "Creating XCFramework..."

rm -rf ${OUTPUT_NAME}.xcframework

xcodebuild -create-xcframework \
    -library ios-simulator/${OUTPUT_FILE} \
    -library ios-arm64/${OUTPUT_FILE} \
    -output ${OUTPUT_NAME}.xcframework

rm -rf ios-arm64
rm -rf ios-simulator

# Save current build timestamp
save_build_time

echo "Created ${OUTPUT_NAME}.xcframework"
