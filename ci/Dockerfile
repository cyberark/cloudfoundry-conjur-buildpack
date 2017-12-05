FROM alpine:3.2

ENV BUILD_PACKAGES bash curl-dev ruby-dev build-base
ENV RUBY_PACKAGES ruby ruby-io-console ruby-bundler jq zip

# Update and install all of the required packages.
# At the end, remove the apk cache
RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    apk add $RUBY_PACKAGES && \
    rm -rf /var/cache/apk/*


ENV WD_PATH /conjurinc/cloudfoundry-conjur-buildpack
ENV CI_PATH $WD_PATH/ci
RUN mkdir -p $CI_PATH
WORKDIR $WD_PATH

COPY Gemfile $CI_PATH
RUN cd $CI_PATH && bundle install

RUN mkdir /lib64/; ln -s /lib/ld-linux-x86-64.so.2 /lib64
