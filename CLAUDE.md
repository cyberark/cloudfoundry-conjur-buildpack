# CLAUDE.md

## Overview

This repo is the **CyberArk Secrets Manager (Conjur) Buildpack for Cloud
Foundry** — a Cloud Foundry [supply buildpack](https://docs.cloudfoundry.org/buildpacks/understand-buildpacks.html#supply-script).
When bound to an app alongside a `cyberark-conjur` service instance (from the
[Conjur Service Broker](https://github.com/cyberark/conjur-service-broker)),
it installs a `profile.d` startup script and a `conjur-env` binary that
authenticate to Conjur/Secrets Manager (API key or mTLS `authn-cert`), fetch
secrets described in the app's `secrets.yml` (via [Summon](https://github.com/cyberark/summon)
semantics and the [conjur-api-go](https://github.com/cyberark/conjur-api-go)
client), and inject them as environment variables before the app process
starts. It must be used as a *non-final* buildpack, stacked with a language
buildpack.

## Tech stack

- **Shell (bash)** — the buildpack scripts themselves (`bin/`, `lib/`, `ci/`).
- **Go** (see `default_versions`/`dependencies` in `manifest.yml`, currently
  Go 1.22.x/1.23.x) — the `conjur-env` binary lives in `conjur-env/` (Go
  module, `go.mod`/`go.sum`).
- **Cucumber/Ruby** — integration/behavior tests in `tests/`.
- Docker & Docker Compose — used to build `conjur-env`, run unit tests, and
  spin up a local Conjur + Cloud Foundry dev environment.

## Build / packaging

- `./package.sh` — builds `conjur-env` (`./conjur-env/build.sh`), builds the
  `buildpack-packager` Docker image from `Dockerfile.packager`, and produces
  `conjur_buildpack-v<VERSION>.zip` (version from `utils.sh`'s
  `project_semantic_version`).
- `./conjur-env/build.sh` — builds the Go binary standalone (also invoked by
  `bin/supply` at buildpack-supply time if `vendor/conjur-env` is missing).
- `./upload.sh` — `cf delete-buildpack` + `cf create-buildpack` to push the
  packaged zip to a target CF foundation (requires `cf` CLI login).
- `manifest.yml` drives `buildpack-packager` (lists `include_files`, Go
  dependency versions/SHAs per `cf_stack`). Update it (and
  `lib/install_go.sh`) together when bumping the Go version — see
  CONTRIBUTING.md.

## Tests

- `./ci/test_unit` — runs everything below (Go vet/lint + Go unit tests +
  `retrieve-secrets` shell tests). This is the main entry point.
- `./ci/test_conjur-env` — Go-only: `go vet`, `go test -coverprofile`, plus
  golint-style checks, all inside a `golang` Docker container.
- `./tests/retrieve-secrets/start` — Dockerized unit tests for
  `lib/0001_retrieve-secrets.sh` (see `tests/retrieve-secrets/README.md`).
- `./ci/test_integration` — Cucumber features tagged `not @integration`
  against a local CF + Conjur dev environment (see
  `tests/integration/`, `tests/docker-compose.yml`).
- `./ci/test_e2e` — end-to-end tests requiring a real PCF/CF instance.
- `./ci/start_dev_environment` — brings up local Conjur + CF stack containers
  for manual testing/dev.

## Repo structure

- `bin/supply` — the actual buildpack "supply" script CF invokes; validates
  `secrets.yml` presence and `VCAP_SERVICES` binding, ensures `conjur-env` is
  built, copies `lib/0001_retrieve-secrets.sh` into `profile.d` and
  `conjur-env` into `vendor` under `$DEPS_DIR/$INDEX`.
- `bin/compile` — deprecated no-op, kept only for legacy Heroku/CF compat.
- `bin/supply.bat` / `bin/supply.ps1`, `lib/0001_retrieve-secrets.bat` —
  Windows Server stack equivalents.
- `lib/0001_retrieve-secrets.sh` — the `profile.d` script that runs at app
  start, invokes `conjur-env`, and `export`s returned secrets into the app
  session (careful error handling to avoid leaking secret values; disables
  `xtrace` while running).
- `lib/install_go.sh` — installs Go at supply time if `conjur-env` needs to
  be compiled on the fly (e.g. online buildpack usage).
- `conjur-env/` — the Go module/binary: `main.go`, `authn.go` /
  `authn_api_key.go` / `authn_cert.go` (auth strategies), `build.sh`,
  `Dockerfile`, tests (`main_test.go`, `authn_cert_test.go`).
- `ci/` — all CI/dev shell entry points (`test_unit`, `test_conjur-env`,
  `test_integration`, `test_e2e`, `start_dev_environment`, `utils`).
- `tests/` — Cucumber integration suite (`tests/integration/`) and shell unit
  tests for the profile.d script (`tests/retrieve-secrets/`).
- `manifest.yml` — buildpack-packager metadata: Go dependency versions/SHAs
  per CF stack (`cflinuxfs3`/`cflinuxfs4`) and packaged file list.
- `package.sh`, `upload.sh`, `unpack.sh`, `utils.sh` — packaging/release
  helper scripts.
- `Dockerfile.packager` — image used to run `buildpack-packager`.
- `Jenkinsfile` — CI pipeline definition (primary CI, mirrored partly by
  `.github/workflows/unit-tests.yml`).
- `kics.config` — IaC security-scan exclusions (test/dev compose files only).

## Conventions / gotchas

- This is a **supply-only** buildpack — it never runs standalone; always
  paired with a language buildpack, and detection priority (`cf
  create-buildpack ... 1`) should keep it first.
- `lib/0001_retrieve-secrets.sh` contains an `__BUILDPACK_INDEX__` placeholder
  that `bin/supply` replaces via `sed` — don't hardcode index paths.
- Secret values must never be echoed/logged; the script explicitly disables
  `xtrace` and sanitizes error messages (see `export_err`/`conjur_env_err` in
  `lib/0001_retrieve-secrets.sh`) — preserve this pattern when editing it.
- Environment variable names from `secrets.yml` must be valid shell
  identifiers (exported via bash `export`).
- When bumping the Go version for `conjur-env`, update it consistently in
  `conjur-env/Dockerfile`, `conjur-env/go.mod`, `manifest.yml`, and
  `lib/install_go.sh` (see CONTRIBUTING.md "Updating the conjur-env Binary").
- `CONJUR_BUILDPACK_BYPASS_SERVICE_CHECK=true` skips the `VCAP_SERVICES`
  binding check in `bin/supply` — used for CI testing without a real broker.
- Two authentication modes are supported end-to-end (API key vs. `authn-cert`
  mTLS/SPIFFE); auth method is selected via the `authn_type` credential field
  in `VCAP_SERVICES` — see README.md "Authentication Methods" for the exact
  credential fields and URL construction rules.
