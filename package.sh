#!/bin/bash

# This utility script can generate the conjur-env,
# placed in the 'vendor' directory,
# then fully package the buildpack for usage.
#
# The buildpack-packager expects all buildpack relevant files
# and folders to be housed in the top-level directory.

cd "$(dirname $0)"

echo "Removing previous builds..."
rm -rf ./conjur-env/vendor
rm -f "conjur_buildpack-v$(cat VERSION)"

echo "Building the conjur-env..."
./conjur-env/build.sh

echo "Building the image for buildpack-packager..."
docker build -t packager -f Dockerfile.packager .

echo "Packaging the conjur-buildpack as a zip file..."
docker run --rm \
  -w /cyberark \
  -v $(pwd):/cyberark \
  packager \
  /bin/bash -c "buildpack-packager build -any-stack"
