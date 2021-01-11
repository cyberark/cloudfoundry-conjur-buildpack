#!/bin/bash -e

# This utility script can create a buildpack
# in a remote PCF stack. This requires proper
# configuration of the `cf` utility and access
# to a remote stack.

cd "$(dirname $0)"

# BUILDPACK_NAME is provided by our bin/test_e2e script
# It not set, it will default to `conjur_buildpack`
NAME="${BUILDPACK_NAME:-conjur_buildpack}"
ZIP_FILE="conjur_buildpack-v$(cat VERSION).zip"

echo "Deleting previous instances of the conjur-buildpack..."
cf delete-buildpack -f "$NAME"

echo "Creating a new buildpack named $NAME..."
# The `1` specifies where to place the buildpack
# in the detection priority list
cf create-buildpack "$NAME" "$ZIP_FILE" 1
