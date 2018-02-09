#!/usr/bin/env bash

echo "cyberark-conjur-buildpack: installing summon & summon-conjur..."

err_report() {
  local previous_exit=$?
  trap - ERR
  echo "${BASH_SOURCE}: Error on line $1" 1>&2
  exit ${previous_exit}
}
trap 'err_report $LINENO' ERR

_conjur_BUILD_DIR=$1
_conjur_BIN_DIR="$_conjur_BUILD_DIR/bin"
_conjur_ESCAPED_BIN_DIR=$(echo "$_conjur_BIN_DIR" | sed -e 's/[\/&]/\\&/g')

mkdir -p ${_conjur_BIN_DIR}
curl -sSL https://raw.githubusercontent.com/cyberark/summon/v0.6.6/install.sh | sed -e 's/sudo //' -e "s/\/usr\/local\/bin/$_conjur_ESCAPED_BIN_DIR/" -e "s/-d \"\/etc\/bash_completion.d\"/\! 1/" | bash
curl -sSL https://raw.githubusercontent.com/cyberark/summon-conjur/v0.5.0/install.sh | sed -e 's/sudo //' -e "s/\/usr\/local\/lib\/summon/$_conjur_ESCAPED_BIN_DIR/" | bash

unset -f _conjur_ESCAPED_BIN_DIR _conjur_BIN_DIR _conjur_ESCAPED_BIN_DIR
trap - ERR
