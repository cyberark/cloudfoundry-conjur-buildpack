# cloudfoundry-conjur-buildpack

## How it works ?

Make sure PCF is running and has the meta-buildpack installed `https://github.com/cf-platform-eng/meta-buildpack`:

For easy testing without a Conjur instance running. Uncomment the `VCAP_SERVICES` in `./lib/0000_retrieve-secrets.sh` and `./bin/decorate`

Upload `cloudfoundry-conjur-buildpack`

```
./upload
```

Run the test application

```
cd test_app
cf push test_app
```

Get the test app route at

```
cf app test_app | grep routes | awk '{ print $2 }'
```

Visit the test app URL and observe that environment variables have been filled in.

---

SSH into the test app

```
cf ssh test_app
```

Then run `env` and observe that there are no environment variables leaked from the application.

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
