#!/bin/bash -e

. build.sh

cf delete-buildpack -f "$NAME"
cf create-buildpack "$NAME" "$ZIPFILE" $(($(cf buildpacks | grep meta_buildpack | awk '{ print $2 }') + 1)) --enable
