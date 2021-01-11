# CyberArk Conjur Buildpack

The CyberArk Conjur Buildpack is a [supply buildpack](https://docs.cloudfoundry.org/buildpacks/understand-buildpacks.html#supply-script)
that installs scripts to provide convenient and secure access to secrets stored
in Conjur.

The buildpack supplies scripts to your application that do the following:

+ Examine your app to determine the secrets to fetch using a [`secrets.yml`](https://cyberark.github.io/summon/#secrets.yml)
  file in the app root folder or [configured location](#secrets_yaml).

+ Retrieve credentials stored in your app's [`VCAP_SERVICES`](https://docs.run.pivotal.io/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES)
  environment variable to communicate with the bound `cyberark-conjur` service
  instance of the [Conjur Service Broker](https://github.com/cyberark/conjur-service-broker).

+ Authenticate using the Conjur credentials, fetch the relevant secrets from
  Conjur, and inject them into the session environment variables at the start of
  the app. The secrets are only available to the app process.
  
## Requirements

+ Your app must be bound to a Conjur service instance. For more information on
  binding your application to a Conjur service instance, see the [Conjur Service Broker documentation](https://github.com/cyberark/conjur-service-broker#bind-your-application-to-the-conjur-service).

+ Your app must have a `secrets.yml` file in its root directory or in a [configured location](#secrets_yaml).

+ Using this buildpack **requires** using multiple buildpacks, because this is not a
  final buildpack and your app will still need to also invoke a language buildpack. To use Cloud Foundry
  with multiple buildpacks, you **must** ensure your [Cloud Foundry CLI](https://github.com/cloudfoundry/cli)
  version is greater than `6.38`. See
  the [CloudFoundry documentation on multiple buildpacks](https://docs.cloudfoundry.org/buildpacks/use-multiple-buildpacks.html)
  for more information.

## How Does the Buildpack Work?

The buildpack uses a [supply script](https://docs.cloudfoundry.org/buildpacks/understand-buildpacks.html#supply-script)
to copy files into the application's dependency directory under a subdirectory
corresponding to the buildpack's index. The `lib/0001_retrieve-secrets.sh`
script is copied into a `profile.d` subdirectory so that it will run automatically
when the app starts and the `conjur-env` binary is copied to a `vendor`
subdirectory. In other words, your application will end up with the following
two files:

```
- $DEPS_DIR/$BUILDPACK_INDEX/profile.d/0001-retrieve-secrets.sh
- $DEPS_DIR/$BUILDPACK_INDEX/vendor/conjur-env
```

The `profile.d` script is run automatically when the application starts and is
responsible for retrieving secrets and injecting them into the app's session
environment variables.

The `conjur-env` binary leverages the [Conjur Go API](https://github.com/cyberark/conjur-api-go)
and [Summon](https://github.com/cyberark/summon) to authenticate with Conjur and
retrieve secrets using the application identity provided by the Conjur Service
Broker.

## Getting Started

The Conjur Buildpack can be included in a CloudFoundry application as an online
buildpack, using the GitHub repository address, or installed into a
CloudFoundry foundation.

For documentation on how to use the online buildpack, see [using](#online)
below for details.

### Installing the Conjur Buildpack

**Before you begin, ensure you are logged into your CF deployment via the CF CLI.**

To install the Conjur Buildpack, download a ZIP of [the latest release](https://github.com/cyberark/cloudfoundry-conjur-buildpack/releases),
unzip the release into its own directory, and run the `upload.sh` script:

```shell
# Download latest version of the Conjur Buildpack
wget -q --show-progress \
  -O "${PWD}/conjur_buildpack.zip" \
  $(curl -s \
  "https://api.github.com/repos/cyberark/cloudfoundry-conjur-buildpack/releases/latest" \
  | jq '.assets[0].browser_download_url' \
  | sed 's/"//g')

# Create the buildpack in your remote stack
cf create-buildpack conjur_buildpack conjur_buildpack.zip 1
```

The '1' will place it at the top of the detection priority.
This is recommended to ensure proper installation of the
secret retrieval script.

Alternatively, you can clone the entire repository, and run the following commands:

```shell
./package.sh
./upload.sh
```

The `./package.sh` script will run [buildpack-packager](https://github.com/cloudfoundry/buildpack-packager)
within the `buildpack` directory and create a `.ZIP` file. `upload.sh` will run
`cf create-buildpack` similar to the command above, as well as removing prior
instances of the Conjur Buildpack with `cf delete-buildpack`.

Earlier versions of the Conjur Buildpack (v0.x) may be installed by cloning the
repository and running `./upload.sh`.

### Using the Conjur Buildpack

#### Create a `secrets.yml` File

For each application that will be using the Conjur Buildpack you must create a
`secrets.yml` file. The `secrets.yml` file gives a mapping of **environment
variable name** to a **location where a secret is stored in Conjur**. For more
information about creating this file, [see the Summon documentation](https://cyberark.github.io/summon/#secrets.yml).
There are no sensitive values in the file itself, so it can safely be checked into source control.

The following is an example of a `secrets.yml` file

```
AWS_ACCESS_KEY_ID: !var aws/prod/iam/user/robot/access_key_id
AWS_SECRET_ACCESS_KEY: !var aws/prod/iam/user/robot/secret_access_key
AWS_REGION: us-east-1
SSL_CERT: !var:file ssl/certs/private
```

The above example could resolve to the following environment variables:

```
AWS_ACCESS_KEY_ID: AKIAI44QH8DHBEXAMPLE
AWS_SECRET_ACCESS_KEY: je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
AWS_REGION: us-east-1
SSL_CERT: /tmp/ssl-cert.pem
```

**Note:** Since the buildpack injects secrets into the application runtime
environment using the [bash export method](https://www.gnu.org/savannah-checkouts/gnu/bash/manual/bash.html#index-export),
environment variable names included in the secret.yml file **must** be valid
shell variable names. In particular, they may contain upper or lowercase letters,
numbers, and underscores **only**.

##### <a name="secrets_yaml"></a> Configuring the `secrets.yml` Location

Some final buildpacks do not allow deploying the `secrets.yml` file to the application
root directory at runtime. In this case, the runtime location of the `secrets.yml`
file may be configured by setting the `SECRETS_YAML_PATH` environment variable to
its relative path.

This can be configured in the application's `manifest.yml`:
```yaml
---
applications:
- name: my-app
  services:
  - conjur
  buildpacks:
  - conjur_buildpack
  - php_buildpack
  env:
    SECRETS_YAML_PATH: lib/secrets.yml
```

Alternatively, this may be set using the Cloud Foundry CLI:
```
$ cf set-env {Application Name} SECRETS_YAML_PATH {Relative Path to secrets.yml}
$ cf restage {Application Name}
```

#### Invoke the Installed Buildpack at Deploy Time

When you deploy your application, ensure it is bound to a Conjur service instance
and add the Conjur Buildpack to your `cf push` command:

```sh
cf push my-app -b conjur_buildpack ... -b final_buildpack
```

Alternatively, the buildpacks may be specified in the application manifest, for
example:

```yaml
---
applications:
- name: my-app
  services:
  - conjur
  buildpacks:
  - conjur_buildpack
  - ruby_buildpack
```

When your application starts, the Conjur Buildpack will inject the secrets
specified in the `secrets.yml` file into the application process as environment
variables.

**Note:** If you add the Conjur buildpack to your manifest or `cf push` command
but don't also explicitly include the language buildpack, Cloud Foundry will see
that the Conjur buildpack is not a "final buildpack" and will fail to invoke it.
To use this buildpack, you **must** follow the instructions for
[using multiple buildpacks](https://docs.cloudfoundry.org/buildpacks/use-multiple-buildpacks.html)
and specify _both_ the Conjur buildpack **and** the final, language buildpack for
your app (in that order).

##### <a name="online"></a> Invoking the Online Buildpack at Deploy Time

To use the CyberArk Conjur Buildpack as an online buildpack, use the GitHub
repository address instead of specifying the installed buildpack name. This may
be done with the `cf push` command or using the manifest file.

```sh
cf push my-app -b https://github.com/cyberark/cloudfoundry-conjur-buildpack#latest ... -b final_buildpack
```

```yaml
---
applications:
- name: my-app
  services:
  - conjur
  buildpacks:
  - https://github.com/cyberark/cloudfoundry-conjur-buildpack#latest
  - ruby_buildpack
```

We recommend users specifically reference the `latest` branch when using the
online version of the buildpack. The `latest` branch is up-to-date with the
latest tagged buildpack release. You may also opt to reference a specific
release version `TAG` by referring to the online buildpack as:
```
https://github.com/cyberark/cloudfoundry-conjur-buildpack#TAG
```

## Contributing

We welcome contributions of all kinds to the Conjur Buildpack. For instructions on
how to get started and descriptions of our development workflows, please see our
[contributing guide](CONTRIBUTING.md). 

## License

This repository is licensed under Apache License 2.0 - see [`LICENSE`](LICENSE) for more details.
