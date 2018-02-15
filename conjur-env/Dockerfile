FROM golang:1.8
MAINTAINER Conjur Inc.

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/cyberark/conjur-env
WORKDIR /go/src/github.com/cyberark/conjur-env

COPY . .
RUN dep ensure
