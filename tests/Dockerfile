FROM cloudfoundry/cflinuxfs4

ENV WD_PATH="/cyberark/cloudfoundry-conjur-buildpack" \
    TEST_PATH="/cyberark/cloudfoundry-conjur-buildpack/tests"

WORKDIR "$WD_PATH"

# Install certificate tools
RUN apt-get update && \
    apt-get install -y apt-transport-https \
                       ca-certificates \
                       openssl \
                       ruby-dev

# Install the Cloud Foundry CLI
RUN wget -q https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key && \
    apt-key add cli.cloudfoundry.org.key && \
    echo "deb https://packages.cloudfoundry.org/debian stable main" > /etc/apt/sources.list.d/cloudfoundry-cli.list && \
    apt-get update && \
    apt-get install -y cf-cli

# Install ruby dependencies
RUN mkdir -p "$TEST_PATH" && \
    gem install bundler -v 2.3.26 && \
    # This will prevent warnings about installing gems as root
    bundle config --global silence_root_warning 1

COPY Gemfile "$TEST_PATH"

RUN cd "$TEST_PATH" && \
    bundle install
