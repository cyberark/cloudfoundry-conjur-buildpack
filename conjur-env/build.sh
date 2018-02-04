#!/bin/bash

cd $(dirname $0)

# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
CGO_ENABLED=0 GOOS=linux go build -o conjur-env -a -ldflags '-extldflags "-static"' .

mkdir -p ../deps
mv conjur-env ../deps
