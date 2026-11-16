FROM golang:1.26.10-alpine
LABEL maintainer="CyberArk Software, Inc."

ARG SKIP_GOMOD_DOWNLOAD=false
ARG GOPROXY=https://proxy.golang.org,direct

WORKDIR /conjur-env

COPY go.mod go.sum /conjur-env/

RUN --mount=type=secret,id=gonetrc if [ "$SKIP_GOMOD_DOWNLOAD" != "true" ]; then \
        if [ -s /run/secrets/gonetrc ]; then export NETRC=/run/secrets/gonetrc; fi; \
        go mod download; \
    fi
