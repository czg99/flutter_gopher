#!/bin/bash

cd $(dirname $0)/../../

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler first."
    echo "Visit https://github.com/protocolbuffers/protobuf/releases for installation instructions."
    exit 1
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

ProtoDir="protos"

outPath="lib/src/protos"
if [ ! -d $outPath ]; then
	mkdir -p $outPath
fi

protoc --dart_out=$outPath --proto_path=$ProtoDir $ProtoDir/*.proto
