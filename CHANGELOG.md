# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.6] - 2020-01-11
### Added
- A [`manifest.yml`](./manifest.yml) has been added, allowing for
  the usage of [buildpack-packager](https://github.com/cloudfoundry/buildpack-packager)
  and other native CloudFoundry features. Please refer
  to [`manifest.yml`](./manifest.yml) for information on
  dependencies and deprecation notices thereof, as well as a list
  of files included in the Buildpack.
  [cyberark/cloudfoundry-conjur-buildpack#79](https://github.com/cyberark/cloudfoundry-conjur-buildpack/issues/79)

### Changed
- The project has been reorganized to respect Cloudfoundry Buildpack
  best practices and improve maintainability. This also should
  reduce overall build times, and slightly reduces the size
  of the Conjur Buildpack `.ZIP`.
  [PR cyberark/cloudfoundry-conjur-buildpack#99](https://github.com/cyberark/cloudfoundry-conjur-buildpack/pull/99)
- The default go version has been bumped to `1.15.x` in the manifest,
  with other supported version listed as well.
  [cyberark/cloudfoundry-conjur-buildpack#41](https://github.com/cyberark/cloudfoundry-conjur-buildpack/issues/41)
- Release tags will now be auto-archived in the `latest` branch. Users
  consuming this buildpack via the online buildpack functionality should now
  point their manifests to the `latest` branch and only consume release versions
  of this buildpack.
  [cyberark/cloudfoundry-conjur-buildpack#101](https://github.com/cyberark/cloudfoundry-conjur-buildpack/issues/101)

### Deprecated
- Support for using the Conjur Buildpack with Conjur Enterprise v4 is now deprecated.
  Support will be removed in the next release.
  [cyberark/cloudfoundry-conjur-buildpack#73](https://github.com/cyberark/cloudfoundry-conjur-buildpack/issues/73)

## [2.1.5] - 2020-11-06
### Added
- Added sanity tests and additional validation of environment variable keys specified
  in the `secrets.yml` to ensure compatibility with Bash commands

### Security
- Added custom handling to sanitize error messages related to retrieved secrets
  being exported to the environment

## [2.1.4] - 2020-07-16

### Changed
- Go dependencies updated for the `conjur-env` binary.
  [cyberark/cloudfoundry-conjur-buildpack#66](https://github.com/cyberark/cloudfoundry-conjur-buildpack/issues/66)

## [2.1.3] - 2020-01-29

### Added
- Added a NOTICES.txt file for open source acknowledgements that is included
  in the release ZIP file

### Changed
- The buildpack now properly reads Conjur credentials from `VCAP_SERVICES` when
  `VCAP_SERVICES` contains credentials for other services with the same field
  names (e.g. `version`).
- Go version for the `conjur-env` binary bumped to 1.13.6
- Go dependencies updated for the `conjur-env` binary

## [2.1.2] - 2019-10-28

### Added
- Buildpack supply step now scans build directory for candidate secrets.yml files
  and reports them to the output.
- The runtime location for `secrets.yml` can now be configured by setting the
  `SECRETS_YAML_PATH` environment variable for the Cloud Foundry application. See
  the [README](README.md) for more information.

## [2.1.1] - 2019-05-15

### Changed
- Go version for online buildpack bumped to 1.12

## [2.1.0] - 2019-05-13

### Added
- Buildpack now searches for `secrets.yml` in `BOOT-INF/classes/` to better
  support Java applications by default.
- Added support to use the Conjur buildpack as an online buildpack by referencing
  the github repository directly. See the [README](README.md#online) for more
  information.

### Changed
- Buildpack now copies the secrets retrieval profile script into the application
  directory. This works around a missing feature in the Java buildpack, where it
  does not correctly source from the buildpacks profile directories.
- Go version of conjur-env binary bumped to 1.12
- Go binary updated to use native os homedir method instead of mitchellh lib

## [2.0.1] - 2019-03-19

### Fixed
- bin/compile script is made executable

## [2.0.0] - 2019-02-15

### Changed
- Buildpack is converted to a supply buildpack to support multi-buildpack usage
- Conjur-env binary dependencies are updated
- Conjur-env binary converted to use Go modules

## [1.0.0] - 2018-03-01

### Changed
- Buildpack uses `conjur-env` binary built from the guts of `summon` and `conjur-api-go` instead of installing Summon and Summon-Conjur each time it is invoked.

## [0.3.0] - 2018-02-13

### Added
- Added support for v4 Conjur

## [0.2.0] - 2018-01-29

### Added
- Added supporting files and documentation for the custom buildpack use case

## 0.1.0 - 2018-01-24
### Added
- The first tagged version.

[Unreleased]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.6...HEAD
[2.1.6]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.5...v2.1.6
[2.1.5]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.4...v2.1.5
[2.1.4]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.3...v2.1.4
[2.1.3]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.2...v2.1.3
[2.1.2]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.0.1...v2.1.0
[2.0.1]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.3.0...v1.0.0
[0.3.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.1.0...v0.2.0
