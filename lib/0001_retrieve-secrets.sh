#!/usr/bin/env bash
set +x

echo "cyberark-conjur-buildpack: retrieving & injecting secrets"

err_report() {
  local previous_exit=$?
  trap - ERR
  echo "${BASH_SOURCE}: Error on line $1" 1>&2
  exit ${previous_exit}
}
trap 'err_report $LINENO' ERR

#export VCAP_SERVICES='
#{
#  "cyberark-conjur": [{
#    "credentials": {
#      "appliance_url": "https://conjur.myorg.com/",
#      "authn_api_key": "2389fh3hf9283niiejwfhjsb83ydbn23u",
#      "authn_login": "3F20D12E-A470-4B7B-8778-C8885769887F",
#      "account": "brokered-services"
#    }
#  }]
#}
#'

# inject secrets into environment
pushd $1
  eval "$(./vendor/conjur-env)"
popd

trap - ERR
