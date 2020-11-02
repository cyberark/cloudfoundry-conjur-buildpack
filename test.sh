#!/bin/bash -ex

# Set up the containers to run in their own namespace
COMPOSE_PROJECT_NAME="$(basename "$PWD")_$(openssl rand -hex 3)"
export COMPOSE_PROJECT_NAME

export BRANCH_NAME=${BRANCH_NAME:-$(git symbolic-ref --short HEAD)}

# Sets up conjur and retrieves credentials
. ./setup-conjur.sh

# Run the install script to verify compatibility
docker-compose run --rm tester bash ./lib/install_go.sh

# Check for summon variables
if [ -n "$CF_API_ENDPOINT" ]; then
    # Make sure all of the environment are present for the integration tests
    : ${CF_ADMIN_PASSWORD?"Need to set CF_ADMIN_PASSWORD"}

    # Build the Java CI application
    pushd ci/apps/java
      ./bin/build
    popd
fi

# Build latest test images
docker-compose build

# Unpack latest build of buildpack
docker-compose run --rm tester bash ./unpack.sh

if [ -n "$CF_API_ENDPOINT" ]; then
     # Run tests against latest build of buildpack (including integration tests against remote foundation)
     docker-compose run --rm \
         -w "$BUILDPACK_ROOT_DIR/ci" \
         tester cucumber \
         --format pretty \
         --format junit \
         --out ./features/reports
else
    # Skip the integration tests if the Summon variables are not present
    docker-compose run --rm \
        -w "$BUILDPACK_ROOT_DIR/ci" \
        tester cucumber \
        --format pretty \
        --format junit \
        --out ./features/reports \
        --tags 'not @integration'
fi