#!/bin/bash -e
function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose down -v
}
trap finish EXIT

. ./setup-conjur.sh

docker-compose run --rm -e CONJUR_CREDENTIALS_JSON="$CONJUR_CREDENTIALS_JSON" -v $(pwd):/conjurinc/cloudfoundry-conjur-buildpack -v $(pwd)/buildpack-build:/buildpack-build tester bash
