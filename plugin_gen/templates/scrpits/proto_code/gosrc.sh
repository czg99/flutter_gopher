#!/bin/bash

cd $(dirname $0)/../../

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler first."
    echo "Visit https://github.com/protocolbuffers/protobuf/releases for installation instructions."
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "protoc-gen-go not found, installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
fi

ProtoDir="protos"

outPath="gosrc/protos"
if [ ! -d $outPath ]; then
	mkdir -p $outPath
fi

protoc --go_out=$outPath --proto_path=$ProtoDir $ProtoDir/*.proto

go mod -C gosrc tidy
