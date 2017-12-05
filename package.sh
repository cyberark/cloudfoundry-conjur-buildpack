#!/bin/bash -e

docker-compose -f ci/docker-compose.yml build
docker-compose -f ci/docker-compose.yml run --rm -v $(pwd):/conjurinc/cloudfoundry-conjur-buildpack tester ./build.sh
