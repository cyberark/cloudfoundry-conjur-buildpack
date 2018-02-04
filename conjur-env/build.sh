#!/bin/bash

cd $(dirname $0)

# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
docker run -e CGO_ENABLED=0 -e GOOS=linux -v "$PWD":/go/src/conjur-env -w /go/src/conjur-env golang:1.8 go build -o conjur-env -a -ldflags '-extldflags "-static"' .

mkdir -p ../deps
mv conjur-env ../deps
