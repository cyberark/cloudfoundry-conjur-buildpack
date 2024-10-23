#!/bin/bash

# This utility script can generate the conjur-env,
# placed in the 'vendor' directory,
# then fully package the buildpack for usage.
#
# The buildpack-packager expects all buildpack relevant files
# and folders to be housed in the top-level directory.

# When running in Jenkins, we need to skip the go mod download command since we
# already fetch the latest dependencies with updatePrivateGoDependencies(). We
# use the --skip-gomod-download flag for this purpose.
SKIP_GOMOD_DOWNLOAD=false
while true ; do
  case "$1" in
    --skip-gomod-download ) SKIP_GOMOD_DOWNLOAD=true ; shift ;;
     * ) if [ -z "$1" ]; then break; else echo "$1 is not a valid option"; exit 1; fi;;
  esac
done

cd "$(dirname $0)"

echo "Removing previous builds..."
rm -f "conjur_buildpack-v$(cat VERSION)"

echo "Building the conjur-env..."
export SKIP_GOMOD_DOWNLOAD
./conjur-env/build.sh

echo "Building the image for buildpack-packager..."
docker build -t packager -f Dockerfile.packager .

echo "Packaging the conjur-buildpack as a zip file..."
docker run --rm \
  -w /cyberark \
  -v $(pwd):/cyberark \
  packager \
  /bin/bash -c "buildpack-packager build -any-stack"
