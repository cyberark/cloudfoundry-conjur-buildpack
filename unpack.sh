#!/bin/bash -e

# This utility script can unzip the buildpack to
# a 'conjur_buildpack' directory

cd "$(dirname $0)"

ZIP_FILE="conjur_buildpack-v$(cat VERSION).zip"
BUILDPACK_BUILD_DIR="./conjur_buildpack"

echo "Cleaning up local buildpack instances..."
rm -rf "$BUILDPACK_BUILD_DIR"
mkdir -p "$BUILDPACK_BUILD_DIR"

echo "Unzipping buildpack to $BUILDPACK_BUILD_DIR"
docker run --rm \
 -w /cyberark \
 -v "$(pwd)":/cyberark \
 bash -c "unzip $ZIP_FILE -d $BUILDPACK_BUILD_DIR"
