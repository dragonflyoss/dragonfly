ARG BASE_IMAGE=alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

FROM golang:1.25.5-alpine3.23@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

WORKDIR /go/src/d7y.io/dragonfly/v2

RUN apk --no-cache add bash make gcc libc-dev git

COPY . /go/src/d7y.io/dragonfly/v2

RUN make build-scheduler && make install-scheduler

FROM ${BASE_IMAGE} AS health

ENV GRPC_HEALTH_PROBE_VERSION=v0.4.48

RUN if [ "$(uname -m)" = "ppc64le" ]; then \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-ppc64le; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-arm64; \
    else \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64; \
    fi && \
    chmod +x /bin/grpc_health_probe

FROM ${BASE_IMAGE}

ENV PATH=/opt/dragonfly/bin:$PATH

RUN echo "hosts: files dns" > /etc/nsswitch.conf

RUN apk --no-cache add curl

COPY --from=builder /opt/dragonfly/bin/scheduler /opt/dragonfly/bin/scheduler
COPY --from=health /bin/grpc_health_probe /bin/grpc_health_probe

EXPOSE 8002

ENTRYPOINT ["/opt/dragonfly/bin/scheduler"]
