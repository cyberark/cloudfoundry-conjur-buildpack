#!/bin/bash -e

cd $(dirname $0)

. ./docker_vars.sh

rm -rf vendor/conjur-env
# http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
docker-compose -f conjur-env/docker-compose.yml build
docker-compose -f conjur-env/docker-compose.yml run --rm conjur-env-builder

docker-compose build
docker-compose run --rm tester ./package_buildpack.sh
