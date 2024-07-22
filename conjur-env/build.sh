#!/bin/bash -e

cd "$(dirname "$0")"

rm -rf ../vendor/conjur-env

docker compose build
docker compose run --rm conjur-env-builder
docker compose run --rm conjur-win-env-builder
