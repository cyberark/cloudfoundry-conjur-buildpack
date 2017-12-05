#!/bin/bash -e

docker-compose -f ci/docker-compose.yml run --rm -v $(pwd):/conjurinc/cloudfoundry-conjur-buildpack -v $(pwd)/buildpack-build:/buildpack-build tester bash
