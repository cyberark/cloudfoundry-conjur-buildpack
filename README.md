The CyberArk Conjur Buildpack is a [supply buildpack](https://docs.cloudfoundry.org/buildpacks/custom.html#contract) that provides convenient and secure access to secrets stored in Conjur.

The buildpack supplies scripts to your application that do the following:

+ Examine your app to determine the secrets to fetch using a [`secrets.yml`](https://cyberark.github.io/summon/#secrets.yml) file in the app root folder.
+ Retrieve credentials stored in your app's [`VCAP_SERVICES`](https://docs.run.pivotal.io/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES) environment variable to communicate with the bound `cyberark-conjur` service.
+ Authenticate using the Conjur credentials, fetch the relevant secrets from Conjur, and inject them into the session environment variables at the start of the app. The secrets are only available to the app process.

## Requirements

+ Your app must be bound to a Conjur service instance. For more information on binding your application to a Conjur service instance, see the [Conjur Service Broker documentation](https://github.com/cyberark/conjur-service-broker#bind-your-application-to-the-conjur-service)

+ Your app must have a `secrets.yml` file in its root directory when deployed

## How Does the Buildpack Work ?

The buildpack uses a [supply script](https://docs.cloudfoundry.org/buildpacks/understand-buildpacks.html#supply-script) to copy files into the application's dependency directory under a subdirectory corresponding to the buildpack's index. The `lib/0001_retrieve-secrets.sh` script is copied into a `.profile.d` subdirectory so that it will run automatically when the app starts and the `conjur-env` binary is copied to a `vendor` subdirectory. In other words, your application will end up with the following two files:

```
- $DEPS_DIR/$BUILDPACK_INDEX/.profile.d/0001-retrieve-secrets.sh
- $DEPS_DIR/$BUILDPACK_INDEX/vendor/conjur-env
```

The `.profile.d` script is run automatically when the application starts and is responsible for retrieving secrets and injecting them into the app's session environment variables.

The `conjur-env` binary leverages the [Conjur Go API](https://github.com/cyberark/conjur-api-go) and [Summon](https://github.com/cyberark/summon)
to authenticate with Conjur and retrieve secrets.

The buildpack has a cucumber test suite. This validates the functionality and also offers great insight into the intended functionality of the buildpack. Please see `./ci/features`.

## Getting Started

### Installing the Conjur Buildpack

**Before you begin, ensure you are logged into your CF deployment via the CF CLI.**

To install the Conjur Buildpack, download a ZIP of [the latest release](https://github.com/cyberark/cloudfoundry-conjur-buildpack/releases),
unzip the release into its own directory, and run the `upload.sh` script:
```
mkdir conjur-buildpack
cd conjur-buildpack/
curl -L $(curl -s https://api.github.com/repos/cyberark/cloudfoundry-conjur-buildpack/releases/latest | \
          jq .assets[0].browser_download_url | \
          sed 's/"//g') \
          > conjur-buildpack.zip
unzip conjur-buildpack.zip
./upload.sh
```

Earlier versions of the Conjur Buildpack (v0.x) may be installed by cloning the repository and running `./upload.sh`.

### Using the Conjur Buildpack

#### Create a `secrets.yml` File

For each application that will be using the Conjur Buildpack you must create a `secrets.yml` file. The `secrets.yml` file gives a mapping of **environment variable name** to a **location where a secret is stored in Conjur**. For more information about creating this file, [see the Summon documentation](https://cyberark.github.io/summon/#secrets.yml). There are no sensitive values in the file itself, so it can safely be checked into source control.

The following is an example of a `secrets.yml` file

```
AWS_ACCESS_KEY_ID: !var aws/$environment/iam/user/robot/access_key_id
AWS_SECRET_ACCESS_KEY: !var aws/$environment/iam/user/robot/secret_access_key
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

#### Invoke the Buildpack at Deploy Time

When you deploy your application, ensure it is bound to a Conjur service instance and add the Conjur Buildpack to your `cf push` command:

```
cf push my-app -b conjur-buildpack -b final-buildpack
```

When your application starts, the Conjur Buildpack will inject the secrets specified in the `secrets.yml` file into the application process as environment variables.

## Development

Before getting started, you should install some developer tools. These are not required to deploy the Conjur Service Broker but they will let you develop using a standardized, expertly configured environment.

1. [git][get-git] to manage source code
2. [Docker][get-docker] to manage dependencies and runtime environments
3. [Docker Compose][get-docker-compose] to orchestrate Docker environments

[get-docker]: https://docs.docker.com/engine/installation
[get-git]: https://git-scm.com/downloads
[get-docker-compose]: https://docs.docker.com/compose/install

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Make sure your Pull Request includes an update to the [CHANGELOG](https://github.com/cyberark/cloudfoundry-conjur-buildpack/blob/master/CHANGELOG.md) describing your changes.

### Updating the `summon` and `conjur-api-go` binaries

To update the versions of `summon` / `conjur-api-go` that are included in the buildpack, update to the appropriate version number in `conjur-env/Gopkg.toml`. To update `conjur-env/vendor/` and `conjur-env/Gopkg.lock`, start up a `conjur-env` container and run `dep ensure`:
```
$ docker-compose -f conjur-env/docker-compose.yml run --entrypoint bash conjur-env-builder
# dep ensure
```
You can verify that the correct dependencies are being used by running `dep status` from within the container before exiting and committing your changes.

Once you have done this, the next time `./package.sh` is run the `vendor/conjur-env` directory will be created with updated dependencies.

### Testing

To test the usage of the Conjur Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

To run the test suite on your local machine:
```
$ ./package.sh   # Create the conjur-env binary in the vendor dir and a ZIP of the project contents
$ ./test.sh      # Run the test suite
```

### Releasing

1. Based on the unreleased content, determine the new version number and update the [VERSION](VERSION) file.
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

## License

Copyright 2018 CyberArk

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
