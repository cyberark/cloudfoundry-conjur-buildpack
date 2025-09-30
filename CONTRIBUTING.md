# Contributing to the Secrets Manager Buildpack

Thanks for your interest in contributing to the Secrets Manager Buildpack! Here
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

The following prerequisites and all sections below pertain to building and running the Secrets Manager Buildpack locally,
unless otherwise specified.

Before getting started, you should install some developer tools. These are not required to deploy the Secrets Manager Buildpack but they will let you develop using a standardized, expertly configured environment.

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

To test the usage of the Secrets Manager Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

### Running the Dev Environment

To test your changes within a running instance of [Cloud Foundry Stack](https://docs.cloudfoundry.org/devguide/deploy-apps/stacks.html)
and Secrets Manager, run:

```shell script
./ci/start_dev_environment
```

This starts Secrets Manager and Cloud Foundry Stack containers, and provides terminal
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

This starts Secrets Manager and Cloud Foundry Stack containers, and 
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

## Releases

### Verify and update dependencies

1.  Review the changes to `conjur-env/go.mod` since the last release and make any needed
    updates to [NOTICES.txt](./NOTICES.txt):
    *   Verify that dependencies fit into supported licenses types:
        ```shell
         cd conjur-env && \
         go-licenses check ./conjur-env/... --allowed_licenses="MIT,ISC,Apache-2.0,BSD-3-Clause" \
            --ignore github.com/cyberark/cloudfoundry-conjur-buildpack/conjur-env \
            --ignore $(go list std | awk 'NR > 1 { printf(",") } { printf("%s",$0) } END { print "" }')
        ```
        If there is new dependency having unsupported license, such license should be included to [notices.tpl](./notices.tpl)
        file in order to get generated in NOTICES.txt.  

        NOTE: The second ignore flag tells the command to ignore standard library packages, which
        may or may not be necessary depending on your local Go installation and toolchain.

    *   If no errors occur, proceed to generate updated NOTICES.txt:
        ```shell
         go-licenses report ./... --template ../notices.tpl > ../NOTICES.txt \
            --ignore github.com/cyberark/cloudfoundry-conjur-buildpack/conjur-env \
            --ignore $(go list std | awk 'NR > 1 { printf(",") } { printf("%s",$0) } END { print "" }')

### Update the version and changelog

1.  Create a new branch for the version bump.

2.  Ensure the [changelog](CHANGELOG.md) is up to date with the changes included in the release, and that
    the unreleased changes are captured under an appropriate version number. This project uses [semantic versioning](https://semver.org/).

3.  Ensure the [open source acknowledgements](NOTICES.txt) are up to date with the dependencies in the
    [conjur-env binary](conjur-env/go.mod) and update the file if there have been any new or changed 
    dependencies since the last release.

4.  Commit these changes - `Bump version to x.y.z` is an acceptable commit message - and open a PR
    for review. Your PR should include updates to
    `CHANGELOG.md`, and if there are any license updates, to `NOTICES.txt`.

### Release and Promote

1.  Jenkins build parameters can be utilized to release and promote successful builds.

2.  Merging into main/master branches will automatically trigger a release.

3.  Reference the [internal automated release doc](https://github.com/conjurinc/docs/blob/master/reference/infrastructure/automated_releases.md#release-and-promotion-process)
    for releasing and promoting.

   **IMPORTANT** Do not upload any artifacts besides the ZIP to the GitHub
   release. At this time, the tile build assumes the project ZIP is the only
   artifact.
