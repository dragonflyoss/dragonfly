ARG BASE_IMAGE=alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

FROM golang:1.25.5-alpine3.23@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

WORKDIR /go/src/d7y.io/dragonfly/v2

RUN apk --no-cache add bash make gcc libc-dev git

COPY . /go/src/d7y.io/dragonfly/v2

RUN make build-datacontroller && make install-datacontroller

FROM ${BASE_IMAGE}

WORKDIR /opt/dragonfly

ENV PATH=/opt/dragonfly/bin:$PATH

RUN mkdir -p /opt/dragonfly/bin /var/log/dragonfly \
    && echo "hosts: files dns" > /etc/nsswitch.conf \
    && chown -R 65532:65532 /var/log/dragonfly

COPY --from=builder /opt/dragonfly/bin/datacontroller /opt/dragonfly/bin/datacontroller

USER 65532:65532

EXPOSE 8000 8001

ENTRYPOINT ["/opt/dragonfly/bin/datacontroller"]
