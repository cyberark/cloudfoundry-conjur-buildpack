# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

The first tagged version.

[Unreleased]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v1.0.0...HEAD
[2.0.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.3.0...v1.0.0
[0.3.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cyberark/cloudfoundry-conjur-buildpack/compare/v0.1.0...v0.2.0
