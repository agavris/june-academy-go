#!/bin/bash

PLATFORMS="darwin/amd64 linux/amd64 windows/amd64 linux/arm64 darwin/arm64"
OUTPUT_DIR="bin"

mkdir -p ${OUTPUT_DIR}

for PLATFORM in ${PLATFORMS}; do
    IFS="/" read -r GOOS GOARCH <<< "${PLATFORM}"
    OUTPUT_NAME='scheduler-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi

    echo "Building $OUTPUT_NAME..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ${OUTPUT_DIR}/${OUTPUT_NAME} .
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "Compilation finished."
