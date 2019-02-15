# Contributing to the Conjur Buildpack

Thanks for your interest in contributing to the Conjur Buildpack! Here
are some guidelines on how to get started.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Pull Request Workflow](#pull-request-workflow)
- [Updating the `conjur-env` Binary](#updating-the-conjur-env-binary)
- [Testing](#testing)
- [Releasing](#releasing)

## Prerequisites

Before getting started, you should install some developer tools. These are not required to deploy the Conjur Buildpack but they will let you develop using a standardized, expertly configured environment.

1. [git][get-git] to manage source code
2. [Docker][get-docker] to manage dependencies and runtime environments
3. [Docker Compose][get-docker-compose] to orchestrate Docker environments

[get-docker]: https://docs.docker.com/engine/installation
[get-git]: https://git-scm.com/downloads
[get-docker-compose]: https://docs.docker.com/compose/install

In addition, if you will be making changes to the `conjur-env` binary, you should
ensure you have [Go installed](https://golang.org/doc/install#install) locally.
Our project uses Go modules, so you will want to install version 1.11+.

### Pull Request Workflow

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Make sure your Pull Request includes an update to the [CHANGELOG](https://github.com/cyberark/cloudfoundry-conjur-buildpack/blob/master/CHANGELOG.md) describing your changes.

### Updating the `conjur-env` Binary

To update the versions of `summon` / `conjur-api-go` that are included in the `conjur-env` binary in the buildpack:

- Bump the appropriate version number in `conjur-env/Gopkg.toml`
- Start up a `conjur-env` container and run `dep ensure` to update `conjur-env/vendor/` and `conjur-env/Gopkg.lock`:
  ```
  $ docker-compose -f conjur-env/docker-compose.yml run --entrypoint bash conjur-env-builder
  # dep ensure
  ```
  You can verify that the correct dependencies are being used by running `dep status` from within the container before exiting and committing your changes.
- Commit your changes
- The next time `./package.sh` is run the `vendor/conjur-env` directory will be created with updated dependencies.

### Testing

To test the usage of the Conjur Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

To run the test suite on your local machine:
```
$ ./package.sh   # Create the conjur-env binary in the vendor dir and a ZIP of the project contents
$ ./test.sh      # Run the test suite
```

### Releasing

1. Based on the unreleased content, determine the new version number and update the [VERSION](VERSION) file. This project uses [semantic versioning](https://semver.org/).
1. Ensure the [changelog](CHANGELOG.md) is up to date with the changes included in the release.
1. Commit these changes - `Bump version to x.y.z` is an acceptable commit message.
1. Once your changes have been reviewed and merged into master, tag the version
   using `git tag -s v0.1.1`. Note this requires you to be  able to sign releases.
   Consult the [github documentation on signing commits](https://help.github.com/articles/signing-commits-with-gpg/)
   on how to set this up. `vx.y.z` is an acceptable tag message.
1. Push the tag: `git push vx.y.z` (or `git push origin vx.y.z` if you are working
   from your local machine).
1. From a **clean checkout of master** run `./package.sh` to generate the release ZIP. Upload this to the GitHub
   release.
   
   **IMPORTANT** Do not upload any artifacts besides the ZIP to the GitHub release. At this time, the tile build
   assumes the project ZIP is the only artifact.
