FROM cloudfoundry/cflinuxfs3

ENV WD_PATH="/conjurinc/cloudfoundry-conjur-buildpack" \
    CI_PATH="/conjurinc/cloudfoundry-conjur-buildpack/ci"

WORKDIR "$WD_PATH"

# Install the Cloud Foundry CLI
RUN apt-get update && \
    apt-get install -y apt-transport-https \
                       ca-certificates \
                       openssl \
                       ruby2.5-dev

# Install a specific older version, since our remote foundation currently uses
# TAS 2.4
RUN curl -L "https://packages.cloudfoundry.org/stable?release=linux64-binary&version=6.52.0&source=github-rel" | tar -zx && \
    mv cf /usr/local/bin && \

    mkdir -p "$CI_PATH" && \

    gem install bundler && \
    # This will prevent warnings about installing gems as root
    bundle config --global silence_root_warning 1

COPY Gemfile \
     Gemfile.lock "$CI_PATH/"

RUN cd "$CI_PATH" && \
    bundle install
