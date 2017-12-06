#!/bin/bash -e

function finish {
  echo 'Removing environment'
  echo '-----'
  docker-compose down -v
}
trap finish EXIT

. ./setup-conjur.sh

docker-compose run --rm -v $(pwd):/conjurinc/cloudfoundry-conjur-buildpack tester bash <<'EOF'
UNZIP_DIR="/tmp/unzipped-buildpack-build"
BUILDPACK_DIR="$(pwd)/buildpack-build"
SRC_DIR=$(cd "$(dirname $0)/." && pwd)
NAME=$(basename "$SRC_DIR")
ESCAPED_NAME=$(echo $NAME | sed s/-/_/g)

rm -rf $BUILDPACK_DIR
mkdir -p $UNZIP_DIR $BUILDPACK_DIR

unzip $ESCAPED_NAME -d $UNZIP_DIR

BUILT_PACKAGE=$(find $UNZIP_DIR -type d -name $NAME)
[ ! -z $BUILT_PACKAGE ] && mv $BUILT_PACKAGE/* $BUILDPACK_DIR
EOF

docker-compose run --rm \
 -e CONJUR_CREDENTIALS_JSON="$CONJUR_CREDENTIALS_JSON" \
 -w /ci \
 -v $(pwd)/ci:/ci \
 -v $(pwd)/buildpack-build:/buildpack-build \
 tester cucumber --format pretty --format junit --out /ci/features/reports
