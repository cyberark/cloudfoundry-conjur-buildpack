#!/bin/bash -e

SRCDIR=$(cd "$(dirname $0)/."  && pwd)
TGTDIR=$SRCDIR
NAME=$(basename "$SRCDIR" | sed s/-/_/g)
ZIPFILE="$TGTDIR/$NAME.zip"

echo "Buildpack name: $NAME"
echo "Source directory: $SRCDIR"
echo "Target file: $ZIPFILE"

rm -f "$ZIPFILE"
zip -r "$ZIPFILE" "$SRCDIR"/bin "$SRCDIR"/lib
