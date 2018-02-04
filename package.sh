#!/bin/bash -e

cd $(dirname $0)

. ./docker_vars.sh

# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
docker run --rm \
 -v "$PWD/deps":/deps \
 -v "$PWD/conjur-env":/go/src/conjur-env \
 -w /go/src/conjur-env \
 -e CGO_ENABLED=0 \
 -e GOOS=linux \
 golang:1.8-alpine3.5 \
 go build -o /deps/conjur-env -a -ldflags '-extldflags "-static"' .

docker-compose build
docker-compose run --rm tester ./package_buildpack.sh
