FROM golang:1.23
LABEL maintainer="CyberArk Software, Inc."

ARG SKIP_GOMOD_DOWNLOAD=false

WORKDIR /conjur-env

COPY go.mod go.sum /conjur-env/

RUN if [ "$SKIP_GOMOD_DOWNLOAD" != "true" ]; then \
        go mod download; \
    fi
