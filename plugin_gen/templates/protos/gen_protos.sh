#!/bin/bash
# This script generates Protocol Buffers code for various platforms
# It checks for required tools and installs them if missing
#
# Before using this script, please ensure the following directories are added to your environment variables:
# 1. The Go bin directory (usually $GOPATH/bin or $HOME/go/bin)
# 2. The Dart .pub-cache/bin directory (usually $HOME/.pub-cache/bin)
#
# These directories must be included in your PATH environment variable for the script to properly run the required tools.
# Usage: ./gen_protos.sh

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
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-dart is installed
case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*|Windows*)
        # Windows check
        if ! where protoc-gen-dart.bat &> /dev/null; then
            echo "protoc-gen-dart.bat not found, installing..."
            dart pub global activate protoc_plugin
        fi
        ;;
    *)
        # Unix-like check
        if ! command -v protoc-gen-dart &> /dev/null; then
            echo "protoc-gen-dart not found, installing..."
            dart pub global activate protoc_plugin
        fi
        ;;
esac

if [ -d "src" ]; then
    protoc --go_out=src protos/*.proto
    go mod -C src tidy
fi

if [ -d "lib" ]; then
    if [ ! -d "lib/models" ]; then
        mkdir -p lib/models
    fi
    protoc --dart_out=lib/models --proto_path=protos protos/*.proto
fi

if [ -d "android" ]; then
    if [ ! -d "android/src/main/java" ]; then
        mkdir -p android/src/main/java
    fi
    protoc --java_out=android/src/main/java protos/*.proto
fi

if [ -d "ios" ]; then
    if [ ! -d "ios/Classes/models" ]; then
        mkdir -p ios/Classes/models
    fi
    protoc --objc_out=ios/Classes/models --proto_path=protos protos/*.proto
fi

if [ -d "macos" ]; then
    if [ ! -d "macos/Classes/models" ]; then
        mkdir -p macos/Classes/models
    fi
    protoc --objc_out=macos/Classes/models --proto_path=protos protos/*.proto
fi

if [ -d "linux" ]; then
    protoc --go_out=linux/src protos/*.proto
    go mod -C linux/src tidy
fi

if [ -d "windows" ]; then
    protoc --go_out=windows/src protos/*.proto
    go mod -C windows/src tidy
fi
