#!/bin/bash -e

cd $(dirname $0)

. ./docker_vars.sh

docker run --rm -v "$PWD":/buildpack -w /buildpack golang:1.8-alpine3.5 ./conjur-env/build.sh

docker-compose build
docker-compose run --rm tester ./package_buildpack.sh
