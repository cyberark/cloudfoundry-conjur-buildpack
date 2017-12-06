#!/bin/bash -e

SRC_DIR=$(cd "$(dirname $0)/."  && pwd)
TGT_DIR=$SRC_DIR
NAME=$(basename "$SRC_DIR" | sed s/-/_/g)
ZIP_FILE="$TGT_DIR/$NAME.zip"

echo "Buildpack name: $NAME"
echo "Source directory: $SRC_DIR"
echo "Target file: $ZIP_FILE"

rm -f "$ZIP_FILE"
zip -r "$ZIP_FILE" "$SRC_DIR"/bin "$SRC_DIR"/lib
