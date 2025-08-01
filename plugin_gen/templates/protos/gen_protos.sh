#!/bin/bash
# This script generates Protocol Buffers code for various platforms
# It checks for required tools and installs them if missing
#
# Before using this script, please ensure the following directories are added to your environment variables:
# 1. The Go bin directory (usually $GOPATH/bin or $HOME/go/bin)
# 2. The Dart .pub-cache/bin directory (usually $HOME/.pub-cache/bin)
#
# These directories must be included in your PATH environment variable for the script to properly run the required tools.
# Usage: sh gen_protos.sh

cd $(dirname $0)/../

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler first."
    echo "Visit https://github.com/protocolbuffers/protobuf/releases for installation instructions."
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "protoc-gen-go not found, installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
fi

# Check if protoc-gen-dart is installed
case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*|Windows*)
        # Windows check
        if ! where protoc-gen-dart.bat &> /dev/null; then
            echo "protoc-gen-dart.bat not found, installing..."
            dart pub global activate protoc_plugin 21.1.2
        fi
        ;;
    *)
        # Unix-like check
        if ! command -v protoc-gen-dart &> /dev/null; then
            echo "protoc-gen-dart not found, installing..."
            dart pub global activate protoc_plugin 21.1.2
        fi
        ;;
esac

protoc --go_out=. --proto_path=protos protos/proto/*.proto
go mod -C protos tidy

if [ -d "src" ]; then
    go mod -C src tidy
fi

if [ -d "linux/src" ]; then
    go mod -C linux/src tidy
fi

if [ -d "windows/src" ]; then
    go mod -C windows/src tidy
fi

if [ -d "lib" ]; then
    if [ ! -d "lib/protos" ]; then
        mkdir -p lib/protos
    fi
    protoc --dart_out=lib/protos --proto_path=protos protos/proto/*.proto
fi

if [ -d "android" ]; then
    if [ ! -d "android/src/main/java" ]; then
        mkdir -p android/src/main/java
    fi
    protoc --java_out=android/src/main/java --proto_path=protos protos/proto/*.proto
fi

if [ -d "ios" ]; then
    if [ ! -d "ios/Classes/protos" ]; then
        mkdir -p ios/Classes/protos
    fi
    protoc --objc_out=ios/Classes/protos --proto_path=protos protos/proto/*.proto
fi

if [ -d "macos" ]; then
    if [ ! -d "macos/Classes/protos" ]; then
        mkdir -p macos/Classes/protos
    fi
    protoc --objc_out=macos/Classes/protos --proto_path=protos protos/proto/*.proto
fi
