# Contributing to the Conjur Buildpack

Thanks for your interest in contributing to the Conjur Buildpack! Here
are some guidelines on how to get started.

For general contribution and community guidelines, including our
pull request workflow, please see the [community repo](https://github.com/cyberark/community).

## Table of Contents

* [Table of Contents](#table-of-contents)
* [Prerequisites](#prerequisites)
* [Updating the `conjur-env` Binary](#updating-the-conjur-env-binary)
* [Testing](#testing)
  + [Running the Dev Environment](#running-the-dev-environment)
  + [Unit Testing](#unit-testing)
  + [Local Integration Testing](#local-integration-testing)
  + [End to End Testing](#end-to-end-testing)
  + [Cleanup](#cleanup)
* [Releasing](#releasing)

<!--
Table of contents generated with markdown-toc
http://ecotrust-canada.github.io/markdown-toc/
-->

## Prerequisites

The following prerequisites and all sections below pertain to building and running the Conjur Buildpack locally,
unless otherwise specified.

Before getting started, you should install some developer tools. These are not required to deploy the Conjur Buildpack but they will let you develop using a standardized, expertly configured environment.

1. [git][get-git] to manage source code
2. [Docker][get-docker] to manage dependencies and runtime environments
3. [Docker Compose][get-docker-compose] to orchestrate Docker environments

[get-docker]: https://docs.docker.com/engine/installation
[get-git]: https://git-scm.com/downloads
[get-docker-compose]: https://docs.docker.com/compose/install

In addition, if you will be making changes to the `conjur-env` binary, you should
ensure you have [Go installed](https://golang.org/doc/install#install) locally.
Our project uses Go modules, so you will want to install version 1.12+.

## Updating the `conjur-env` Binary

The `conjur-env` binary uses Go modules to manage dependencies.
To update the versions of `summon` / `conjur-api-go`
that are included in the `conjur-env` binary in the buildpack,
make sure you have Go installed locally (at least version 1.12) and run:

```
$ cd conjur-env/
$ go get github.com/cyberark/[repo]@v[version]
```

This will automatically update go.mod and go.sum.

Commit your changes, and the next time `./conjur-env/build.sh` is run the
`vendor/conjur-env`directory will be created with updated dependencies.

When upgrading the version of Go for `conjur-env`, the value needs to be updated
in a few places:

* Update the base image in `./conjur-env/Dockerfile`
* Update the Go version in `./conjur-env/go.mod`
* Update the version and file hashes in `manifest.yml` - available versions and
  hashes can be found [here][buildpacks], or see the manifest for the
  [official Go Buildpack][go-buildpack]. (This is for the offline version of
  the buildpack, which is built with buildpack-packager.)
* Update the version and SHA hash in `lib/install_go.sh` -- you can
  find the available versions and hashes on the [CF dependencies][deps] page.

[buildpacks]: https://buildpacks.cloudfoundry.org/#/buildpacks/
[go-buildpack]: https://github.com/cloudfoundry/go-buildpack/blob/master/manifest.yml
[deps]: https://buildpacks.cloudfoundry.org/#/dependencies

## Testing

The buildpack has a cucumber test suite. This validates the functionality and
also offers great insight into the intended functionality of the buildpack.
Please see `./tests/features`.

To test the usage of the Conjur Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

### Running the Dev Environment

To test your changes within a running instance of [Cloud Foundry Stack](https://docs.cloudfoundry.org/devguide/deploy-apps/stacks.html)
and Conjur, run:

```shell script
./ci/start_dev_environment
```

This starts Conjur and Cloud Foundry Stack containers, and provides terminal
access to the Cloud Foundry container. You do not need to restart the container
after you make changes to the project.

To run the local `cucumber` tests within the development environment, run the following 
command from the `tests/integration` directory, within the container:

```shell script
cucumber \
    --format pretty \
    --format junit \
    --out ./features/reports \
    --tags 'not @integration'
```

### Unit Testing

Unit tests are comprised of two categories:

- Unit tests, linting, and code coverage for `conjur-env` Golang module
- Unit tests for `lib/0001_retrieve-secrets.sh`

To run all tests for the `conjur-env` Golang module *and* for
`lib/0001_retrieve-secrets.sh`, you can run:

```shell script
./ci/test_unit
```

To run all tests for _only_ the `conjur-env` Golang module, run:

```shell script
./ci/test_conjur-env
```

To run all tests for _only_ `0001_retrieve-secrets.sh`, run:

```shell script
./tests/retrieve-secrets/start
```

See the [README.md](tests/retrieve-secrets/README.md) for more information.

### Local Integration Testing

To run the set of features marked with `not @integration`,
which are the subset of `cucumber` integration tests not dependent
on a remote PCF instance or privileged credentials. Run:

```shell script
./ci/test_integration
```

This starts Conjur and Cloud Foundry Stack containers, and 
runs the `cucumber` tests within. 

### End to End Testing

To run the Buildpack end-to-end tests, the test script needs to be given the API
endpoint and admin credentials for a CloudFoundry installation.
These are provided as environment variables to the script:

```shell script
export CF_API_ENDPOINT=https://api.sys.cloudfoundry.net
CF_ADMIN_PASSWORD=... ./ci/test_e2e
```

These variables may also be provided using [Summon](https://cyberark.github.io/summon/)
by updating the `ci/secrets.yml` file as needed and running:

```shell script
summon -f ./ci/secrets.yml ./ci/test_e2e
```

This requires access to privileged credentials.

### Cleanup

If integration tests fail, it's possible that some artifacts may not
be cleaned up properly. To clean up leftover components from running
integration tests on a remote PCF environment, run:

```shell script
./ci/clear_ci_artifacts
```

## Releasing

1. Based on the unreleased content, determine the new version number and update the [VERSION](VERSION) file. This project uses [semantic versioning](https://semver.org/).
1. Ensure the [changelog](CHANGELOG.md) is up to date with the changes included in the release.
1. Ensure the [open source acknowledgements](NOTICES.txt) are up to date with
   the dependencies in the [conjur-env binary](conjur-env/go.mod), and
   update the file if there have been any new or changed dependencies
   since the last release.
1. Commit these changes - `Bump version to x.y.z` is an acceptable commit message.
1. Once your changes have been reviewed and merged into master, tag the version
   using `git tag -s v0.1.1`. Note this requires you to be  able to sign releases.
   Consult the [github documentation on signing commits](https://help.github.com/articles/signing-commits-with-gpg/)
   on how to set this up. `vx.y.z` is an acceptable tag message.
1. Push the tag: `git push vx.y.z` (or `git push origin vx.y.z` if you are working
   from your local machine).
1. Retrieve the archived release ZIP from the Jenkins build. Upload this ZIP
   file to the GitHub release.

   **IMPORTANT** Do not upload any artifacts besides the ZIP to the GitHub
   release. At this time, the tile build assumes the project ZIP is the only
   artifact.
