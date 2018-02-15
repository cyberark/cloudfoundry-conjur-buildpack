#!/bin/bash -e

# get context
. ./docker_vars.sh

docker-compose up -d conjur

# wait for conjur
docker-compose exec -T conjur conjurctl wait -r 45 -p 80
