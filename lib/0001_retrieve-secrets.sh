#!/usr/bin/env bash

# Check whether the conjur-env output has already been sourced
# into the environment. If so, return early.
if [ ! -z "$CONJUR_ENV_SOURCED" ]; then
  return
fi

echo "[cyberark-conjur-buildpack]: retrieving & injecting secrets"

# Report an error. This is the default error handler.
err_report() {
    local previous_exit=$?
    trap - ERR
    printf "%s: Error on line %s" "${BASH_SOURCE[@]}" "$1" 1>&2
    exit ${previous_exit}
}
trap 'err_report $LINENO' ERR

# Create a temp file for capturing stderr
temp_err_file=$(mktemp)
trap 'rm -f "$temp_err_file"' EXIT

# Report an error in execution of the conjur-env binary. The error is
# available in a temp file.
conjur_env_err() {
    local previous_exit=$?
    trap - ERR
    printf "%s: Error on line %s: $(<"$temp_err_file")" "${BASH_SOURCE[@]}" "$1"
    exit ${previous_exit}
}

# Report an error in exporting an environmental setting. Sanitize the output
# by extracting only the variable name (exclude the value, which is probably
# sensitive).
export_err() {
    local previous_exit=$?
    trap - ERR
    printf "%s: Error on line %s: Unable to export \`%s\`; value may not be a valid identifier" "${BASH_SOURCE[@]}" "$1" "$2"
    exit ${previous_exit}
}

# __BUILDPACK_INDEX__ is replaced by sed in the 'supply' script
if [[ -z "$CONJUR_ENV_DIR" ]]; then
  CONJUR_ENV_DIR="${DEPS_DIR}/__BUILDPACK_INDEX__/vendor/conjur-env"
fi

# Prevent tracing to ensures secrets won't be leaked.
declare xtrace=""
case $- in
  (*x*) xtrace="xtrace";;
esac
set +x

# $HOME points to the app directory, which should contains a secrets.yml file.
pushd "$HOME"
  # Retrieve environmental settings
  trap 'conjur_env_err $LINENO' ERR
  env="$(${CONJUR_ENV_DIR} 2>"$temp_err_file")"

  # Iterate through and export each statement silently. If there is an
  # error, generate a sanitized error report that includes only the
  # environment variable name.
  IFS="
" # Use newline as separator
  for line in $env; do
    # Values should be base64 encoded when passed to this script,
    # otherwise special characters may be misinterpreted by the shell.
    # Expect lines to be of the form "<variable>: <base64-encoded-value>"
    if [[ $line =~ (.*):[[:space:]]*(.*) ]]; then
        var="${BASH_REMATCH[1]}"
        trap 'export_err $LINENO $var' ERR
        value="$(base64 -d <<< "${BASH_REMATCH[2]}")"
        export "$var=$value" 2> /dev/null
    else
        # Invalid format. Raise an error but don't display sensitive info. 
        echo "Error: ${CONJUR_ENV_DIR} output is not of the form \"<variable>: <base64-encoded-value>\""
        exit 1
    fi
  done
popd

[ ! -z "$xtrace" ] && set -x

export CONJUR_ENV_SOURCED=true
