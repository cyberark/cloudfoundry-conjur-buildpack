#!/bin/bash -e

. ./buildpack_vars.sh

echo "Buildpack name: $NAME"
echo "Source directory: $SRC_DIR"
echo "Target file: $ZIP_FILE"

pushd ${SRC_DIR}
  ls deps/conjur-env || {
    echo "ERROR: conjur-env isn't present in ${SRC_DIR}/deps."
    echo "ERROR: conjur-env should be built and placed in ${SRC_DIR}/deps before running this script";
    exit 1;
  }
  rm -f "$ZIP_FILE"
  zip -r "$ZIP_FILE" bin lib deps
  rm -rf deps
popd
