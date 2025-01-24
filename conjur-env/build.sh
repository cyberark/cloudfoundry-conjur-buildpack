#!/bin/bash -e

cd "$(dirname "$0")"

# Delete existing binaries
rm -rf ../vendor/conjur-env
rm -rf ../vendor/conjur-win-env.exe

# Build the conjur-env binary for Linux and Windows
docker compose build --build-arg SKIP_GOMOD_DOWNLOAD="$SKIP_GOMOD_DOWNLOAD"
docker compose run --rm conjur-env-builder
docker compose run --rm conjur-win-env-builder
