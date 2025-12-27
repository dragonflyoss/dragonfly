FROM golang:1.24.11-alpine3.21@sha256:59a9237590705e00dd04bb43f9f75562caced1393b8767b4666a658daeb93e61

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

WORKDIR /go/src/d7y.io/dragonfly/v2

RUN apk --no-cache add bash make gcc libc-dev git

COPY . /go/src/d7y.io/dragonfly/v2

ENV CGO_ENABLED="1"
ENV GO111MODULE="on"

COPY build/plugin-builder/build.sh /
