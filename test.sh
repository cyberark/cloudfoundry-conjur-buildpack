#!/bin/bash -e

SRCDIR=$(cd "$(dirname $0)/."  && pwd)
NAME=$(basename "$SRCDIR")
ESCAPED_NAME=$(echo $NAME | sed s/-/_/g)

docker-compose -f ci/docker-compose.yml build

rm -rf unzipped-buildpack-build buildpack-build
mkdir -p unzipped-buildpack-build buildpack-build

docker-compose -f ci/docker-compose.yml run --rm -v $(pwd):/conjurinc/cloudfoundry-conjur-buildpack tester unzip $ESCAPED_NAME -d unzipped-buildpack-build

BUILT_PACKAGE=$(find unzipped-buildpack-build -type d -name $NAME)
[ ! -z $BUILT_PACKAGE ] && mv $BUILT_PACKAGE/* buildpack-build

rm -rf unzipped-buildpack-build

docker-compose -f ci/docker-compose.yml run --rm -w /ci -v $(pwd)/ci:/ci -v $(pwd)/buildpack-build:/buildpack-build tester cucumber --format pretty --format junit --out /ci/features/reports
