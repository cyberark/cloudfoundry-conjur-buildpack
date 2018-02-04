The CyberArk Conjur Buildpack is a [decorator buildpack](https://github.com/cf-platform-eng/meta-buildpack#what-is-a-decorator) that provides convenient and secure access to secrets stored in Conjur.

The buildpack carries out the following:
 
+ Examines your app to determine the secrets to fetch using a [`secrets.yml`](https://cyberark.github.io/summon/#secrets.yml) file in the app root folder.
+ Utilizes credentials stored in your app's [`VCAP_SERVICES`](https://docs.run.pivotal.io/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES) environment variable to communicate with the bound `cyberark-conjur` service.
+ The buildpack will authenticate using the aforementioned credentials, fetch the relevant secrets and inject them into the session environment variables at the start of the app. This means the secrets are only available to the app process.

Internally, the Conjur Buildpack uses [Summon](https://cyberark.github.io/summon/) to load secrets into the environment of CF-deployed applications based on the app's `secrets.yml` file.

## Requirements

+ The Conjur Buildpack requires the meta-buildpack. If you are unable to use the meta-buildpack (because you are using a custom buildpack, for example), see the [instructions below](#custom-buildpack-usage) for manually loading the buildpack scripts.

+ Ensure that your app is bound to a Conjur service instance. For more information on binding your application to a Conjur service instance, see the [Conjur Service Broker documentation](https://github.com/conjurinc/conjur-service-broker#binding-your-application-to-the-conjur-service)

+ A `secrets.yml` file exists within the root folder of your app

## How does the buildpack work ?

### `meta-buildpack`

`meta-buildpack` makes it possible to "decorate" the application with the ability to retrieve secrets on startup. Please read [`meta-buildpack` documentation](https://github.com/cf-platform-eng/meta-buildpack#how-it-works) for a quick run through of how `meta-buildpack` works. 

`meta-buildpack` relies on automatic detection of the language buildpack. The first language buildpack in buildpack index order to detect and claim the build will be used to build the droplet and run the application. With this in mind, **it is not currently possible to specify custom buildpacks in the manifest or as an argument to `cf push`**. In these instances `meta-buildpack` will not be invoked; consequently the Conjur buildpack will also not be invoked.

### Lifecycle scripts

The buildpack is comprised of the [3 lifecycle scripts](https://github.com/cf-platform-eng/meta-buildpack#how-to-write-a-decorator) that are required for decorator buildpacks.

+ Detect script always returns a non-zero exit status.
+ Compile script installs summon, summon-conjur in the app's root directory and a `.profile.d` script at `./.profile.d/0001_retrieve-secrets.sh` relative to the app's root directory.
+ Decorate script returns a 0 when `secrets.yml` exists in the root folder of your app and the "cyberark-conjur" key is present in `VCAP_SERVICES`, otherwise non-zero exit status

The `.profile.d` script is responsible for retrieving secrets and injecting them into the session environment variables at the start of the app.

The buildpack has a cucumber test suite. This validates the functionality and also offers great insight into the intended functionality of the buildpack. Please see `./ci/features`.

## Getting Started

### Installing the Conjur Buildpack

**Before you begin, ensure you are logged into your CF deployment via the CF CLI.**

Install the [`meta-buildpack`](https://github.com/cf-platform-eng/meta-buildpack):
```
git clone git@github.com:cf-platform-eng/meta-buildpack
cd meta-buildpack
./build
./upload
```

Install the [Conjur buildpack](https://github.com/conjurinc/cloudfoundry-conjur-buildpack):
```
git clone git@github.com:conjurinc/cloudfoundry-conjur-buildpack
cd cloudfoundry-conjur-buildpack
./upload.sh
```

### Buildpack Usage

#### Create a `secrets.yml` file

To use the Conjur Buildpack with a CF-deployed application, a `secrets.yml` file is required. The `secrets.yml` file gives a mapping of **environment variable name** to a **location where a secret is stored in Conjur**. For more information about creating this file, [see the Summon documentation](https://cyberark.github.io/summon/#secrets.yml). There are no sensitive values in the file itself, so it can safely be checked into source control.

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

#### Push with `meta-buildpack`

As described in the [`meta-buildpack` documentation](https://github.com/cf-platform-eng/meta-buildpack#how-it-works) the Conjur decorate buildpack will **ONLY** be invoked in the following are true:

+ `cf push` with `meta-buildpack` installed and at the top of list of buildpacks. 
+ no buildpack is specified in app's manifest.
+ no `-b` flag value is provided to the CLI. 

Assuming the app is bound to a Conjur service instance, pushing the app will result in the Conjur buildpack being invoked, as well as the language buildpack. The secrets specified in the `secrets.yml` file will now be available in the session environment variables at the start of the app.

### <a name="custom-buildpack-usage"></a>Usage for users with custom buildpacks

To retrieve secrets without using `meta-buildpack` you can simply:

+ copy the contents of the `./lib` directory in a release to the `./.profile.d` directory relative to your app's root directory.
+ copy the contents of the `./vendor` directory in a release to the `./vendor` directory relative to your app's root directory.
 
When your application starts the secrets specified in the `secrets.yml` file will now be available in the session environment variables at the start of the app.

## Development

Before getting started, you should install some developer tools. These are not required to deploy the Conjur Service Broker but they will let you develop using a standardized, expertly configured environment.

1. [git][get-git] to manage source code
2. [Docker][get-docker] to manage dependencies and runtime environments
3. [Docker Compose][get-docker-compose] to orchestrate Docker environments

[get-docker]: https://docs.docker.com/engine/installation
[get-git]: https://git-scm.com/downloads
[get-docker-compose]: https://docs.docker.com/compose/install

To test the usage of the Conjur Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

To run the test suite, call `./test.sh` from your local machine - the script will stand up the needed containers and run the full suite of rspec and cucumber tests.

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

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
