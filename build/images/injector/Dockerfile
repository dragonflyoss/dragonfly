ARG BASE_IMAGE=alpine:3.21

FROM golang:1.23.8-alpine3.21 AS builder

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

WORKDIR /go/src/d7y.io/dragonfly/v2

RUN apk --no-cache add bash make gcc libc-dev git

COPY . /go/src/d7y.io/dragonfly/v2

RUN make build-injector && make install-injector

FROM ${BASE_IMAGE}

WORKDIR /opt/dragonfly

ENV PATH=/opt/dragonfly/bin:$PATH

RUN mkdir -p /opt/dragonfly/bin \
    && echo "hosts: files dns" > /etc/nsswitch.conf

RUN apk --no-cache add ca-certificates

COPY --from=builder /opt/dragonfly/bin/injector /opt/dragonfly/bin/injector

USER 65532:65532

ENTRYPOINT ["/opt/dragonfly/bin/injector"]
