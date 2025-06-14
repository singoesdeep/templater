#!/bin/bash

set -e

# Version from git tag
VERSION=$(git describe --tags --always --dirty)
if [[ $VERSION == v* ]]; then
    VERSION=${VERSION:1}
fi

# Build directory
BUILD_DIR="build"
mkdir -p $BUILD_DIR

# Build for multiple platforms
echo "Building for multiple platforms..."

# Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/templater-linux-amd64 ./cmd/templater
GOOS=linux GOARCH=arm64 go build -o $BUILD_DIR/templater-linux-arm64 ./cmd/templater

# Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/templater-windows-amd64.exe ./cmd/templater
GOOS=windows GOARCH=arm64 go build -o $BUILD_DIR/templater-windows-arm64.exe ./cmd/templater

# macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/templater-darwin-amd64 ./cmd/templater
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/templater-darwin-arm64 ./cmd/templater

# Create checksums
echo "Creating checksums..."
cd $BUILD_DIR
sha256sum templater-* > checksums.txt
cd ..

# Create release archive
echo "Creating release archive..."
RELEASE_NAME="templater-$VERSION"
mkdir -p $RELEASE_NAME
cp $BUILD_DIR/templater-* $RELEASE_NAME/
cp $BUILD_DIR/checksums.txt $RELEASE_NAME/
cp LICENSE README.md $RELEASE_NAME/
tar -czf $RELEASE_NAME.tar.gz $RELEASE_NAME
rm -rf $RELEASE_NAME

echo "Release $VERSION built successfully!"
echo "Files are in $BUILD_DIR/"
echo "Release archive: $RELEASE_NAME.tar.gz" 