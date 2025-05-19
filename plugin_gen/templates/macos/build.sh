#!/bin/sh

cd $(dirname $0)

if ! command -v go &> /dev/null; then
    echo "Error: Go compiler not found. Please install Go."
    exit 1
fi

OUTPUT_NAME="{{.LibName}}"
OUTPUT_FILE="lib${OUTPUT_NAME}.dylib"
OUTPUT_DIR="$(pwd)"
GO_SRC="../src"
TIMESTAMP_FILE=".last_build_time"

MIN_VERSION=10.11

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

export CGO_ENABLED=1
export GOOS=darwin

LIB_FILES=""

for ARCH in "amd64:x86_64" "arm64:arm64"
do
    IFS=: read GOARCH CARCH <<< "$ARCH"
    
    export GOARCH=$GOARCH
    TARGET="$CARCH-apple-macos$MIN_VERSION"
    
    SDK_PATH=$(xcrun --sdk macosx --show-sdk-path)
    CLANG_PATH=$(xcrun --sdk macosx --find clang)
    
    export CC="$CLANG_PATH -target $TARGET -isysroot $SDK_PATH $@"
    export CXX="$CLANG_PATH++ -target $TARGET -isysroot $SDK_PATH $@"
    
    echo "Compiling $GOARCH architecture..."

    if [ "$GOARCH" = "amd64" ]; then
        OUTPUT_FILE_TMP="macos-x86_64/${OUTPUT_FILE}"
    else
        OUTPUT_FILE_TMP="macos-arm64/${OUTPUT_FILE}"
    fi

    LIB_FILES="$LIB_FILES $OUTPUT_FILE_TMP"

    go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-shared -o "$OUTPUT_DIR/$OUTPUT_FILE_TMP"

    if [ $? -ne 0 ]; then
        echo "Error: Go compilation failed, error code: $?"
        echo "Please check Go source code or compilation environment settings"
        exit $?
    else
        echo "Success: Go compilation completed, output file: $OUTPUT_DIR/$OUTPUT_FILE_TMP"
    fi
done

rm -rf ${OUTPUT_FILE}


echo "Merging all architecture library files..."
lipo -create $LIB_FILES -output ${OUTPUT_FILE}

install_name_tool -id @rpath/${OUTPUT_FILE} ${OUTPUT_FILE}

rm -rf macos-arm64
rm -rf macos-x86_64

# Save current build timestamp
save_build_time

echo "Created ${OUTPUT_FILE}"