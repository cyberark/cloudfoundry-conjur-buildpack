#!/bin/bash

# Borrowing this script from the Node.js Buildpack.
# https://raw.githubusercontent.com/cloudfoundry/nodejs-buildpack/master/scripts/install_go.sh

set -euo pipefail

GO_VERSION="1.17.9"

if [ "${CF_STACK}" == "cflinuxfs3" ]; then
    GO_SHA256="a620909b3dc54e17a7769dbcb95bc9fbac60dd376c730b09afaee984d4d7da91"
else
  echo "       **ERROR** Unsupported stack"
  echo "                 See https://docs.cloudfoundry.org/devguide/deploy-apps/stacks.html for more info"
  exit 1
fi

export GoInstallDir="/tmp/go${GO_VERSION}"
mkdir -p "${GoInstallDir}"

if [ ! -f "${GoInstallDir}/go/bin/go" ]; then
  URL="https://buildpacks.cloudfoundry.org/dependencies/go/go_${GO_VERSION}_linux_x64_${CF_STACK}_${GO_SHA256:0:8}.tgz"

  echo "-----> Download go ${GO_VERSION}"
  curl -s -L --retry 15 --retry-delay 2 "${URL}" -o /tmp/go.tar.gz

  DOWNLOAD_SHA256=$(shasum -a 256 /tmp/go.tar.gz | cut -d ' ' -f 1)

  if [[ "${DOWNLOAD_SHA256}" != "${GO_SHA256}" ]]; then
    echo "       **ERROR** SHA256 mismatch: got ${DOWNLOAD_SHA256} expected ${GO_SHA256}"
    exit 1
  fi

  tar xzf /tmp/go.tar.gz -C "${GoInstallDir}"
  rm /tmp/go.tar.gz
fi
if [ ! -f "${GoInstallDir}/bin/go" ]; then
  echo "       **ERROR** Could not download go"
  exit 1
fi
