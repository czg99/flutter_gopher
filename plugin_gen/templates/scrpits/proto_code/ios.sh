#!/bin/bash

cd $(dirname $0)/../../

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler first."
    echo "Visit https://github.com/protocolbuffers/protobuf/releases for installation instructions."
    exit 1
fi

ProtoDir="protos"

outPath="ios/Classes/protos"
if [ ! -d $outPath ]; then
	mkdir -p $outPath
fi

protoc --objc_out=$outPath --proto_path=$ProtoDir $ProtoDir/*.proto
