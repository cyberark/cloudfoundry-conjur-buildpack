FROM golang:1.22
MAINTAINER CyberArk Software, Inc.

WORKDIR /conjur-env

COPY go.mod go.sum /conjur-env/

RUN go mod download
