FROM golang:1.22
MAINTAINER CyberArk Software, Inc.

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /mock-conjur-env

COPY go.mod /mock-conjur-env/

RUN go mod download
