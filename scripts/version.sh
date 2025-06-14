#!/bin/bash

set -e

# Get current version from git
CURRENT_VERSION=$(git describe --tags --always --dirty)
if [[ $CURRENT_VERSION == v* ]]; then
    CURRENT_VERSION=${CURRENT_VERSION:1}
fi

# Function to update version in files
update_version() {
    local version=$1
    local file=$2
    local pattern=$3

    if [[ -f $file ]]; then
        if [[ $OSTYPE == "darwin"* ]]; then
            # macOS
            sed -i '' -E "s/$pattern/$version/g" "$file"
        else
            # Linux
            sed -i -E "s/$pattern/$version/g" "$file"
        fi
    fi
}

# Update version in various files
update_version "$CURRENT_VERSION" "cmd/templater/version.go" "Version = \"[0-9]+\.[0-9]+\.[0-9]+\""
update_version "$CURRENT_VERSION" "Dockerfile" "LABEL version=\"[0-9]+\.[0-9]+\.[0-9]+\""

# Create new version
if [ "$1" != "" ]; then
    NEW_VERSION=$1
    if [[ $NEW_VERSION != v* ]]; then
        NEW_VERSION="v$NEW_VERSION"
    fi

    # Update version in files
    update_version "${NEW_VERSION:1}" "cmd/templater/version.go" "Version = \"[0-9]+\.[0-9]+\.[0-9]+\""
    update_version "${NEW_VERSION:1}" "Dockerfile" "LABEL version=\"[0-9]+\.[0-9]+\.[0-9]+\""

    # Create git tag
    git add cmd/templater/version.go Dockerfile
    git commit -m "Bump version to $NEW_VERSION"
    git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

    echo "Version bumped to $NEW_VERSION"
    echo "Don't forget to push the tag: git push origin $NEW_VERSION"
else
    echo "Current version: $CURRENT_VERSION"
    echo "To create a new version: ./scripts/version.sh x.y.z"
fi 