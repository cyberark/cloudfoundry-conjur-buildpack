FROM cloudfoundry/cflinuxfs2

ENV RUBY_PACKAGES bundler

# Update and install all of the required packages.
RUN gem install $RUBY_PACKAGES

ENV WD_PATH /conjurinc/cloudfoundry-conjur-buildpack
ENV CI_PATH $WD_PATH/ci
RUN mkdir -p $CI_PATH
WORKDIR $WD_PATH

COPY Gemfile $CI_PATH
RUN cd $CI_PATH && bundle install
