#!/bin/bash

# Borrowing this script from the Node.js Buildpack.
# https://raw.githubusercontent.com/cloudfoundry/nodejs-buildpack/master/scripts/install_go.sh

set -euo pipefail

if [ "${CF_STACK}" == "cflinuxfs3" ]; then
    GO_VERSION="1.25.2"
    GO_SHA256="385184a62bdcb565860663d365e2b28cdfbb6919d4439dae7e5cc87694a3dca6"
elif [ "${CF_STACK}" == "cflinuxfs4" ]; then
    GO_VERSION="1.25.12"
    GO_SHA256="72adea67aa47ec63a07a301b958a60ea3ac7ca7ec80a39fd4533a4c26a4cbb90"
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
