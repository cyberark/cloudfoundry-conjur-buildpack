#!/bin/bash

set -euo pipefail

function repo_root() {
	git rev-parse --show-toplevel
}

function project_version() {
	# VERSION derived from CHANGELOG and automated release library
	echo "$(<"$(repo_root)/VERSION")"
}

function project_semantic_version() {
	local version
	version="$(project_version)"

	# Remove Jenkins build number from VERSION
	echo "${version/-*/}"
}

# Ensure VERSION file exists for local builds
if [ ! -f VERSION ]; then
  echo "0.0.0-dev" > VERSION
  echo "Generated dev VERSION file..."
fi
