#!/bin/bash -e

SRC_DIR=$(cd "$(dirname $0)/."  && pwd)
NAME=$(basename "$SRC_DIR" | sed s/-/_/g)
TGT_DIR=${SRC_DIR}
ZIP_FILE="$TGT_DIR/$NAME.zip"