FROM alpine:3.2

ENV BUILD_PACKAGES bash curl-dev ruby-dev build-base
ENV RUBY_PACKAGES ruby ruby-io-console ruby-bundler jq
ENV WD_PATH /conjurinc/cloudfoundry-conjur-buildpack
# Update and install all of the required packages.
# At the end, remove the apk cache
RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    apk add $RUBY_PACKAGES && \
    rm -rf /var/cache/apk/*

RUN mkdir -p $WD_PATH/ci
WORKDIR $WD_PATH

COPY Gemfile $WD_PATH/ci
RUN cd $WD_PATH/ci && bundle install
