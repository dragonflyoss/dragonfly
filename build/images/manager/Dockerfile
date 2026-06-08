ARG BASE_IMAGE=alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

FROM node:20-alpine@sha256:658d0f63e501824d6c23e06d4bb95c71e7d704537c9d9272f488ac03a370d448 AS console-builder

WORKDIR /build

COPY ./manager/console/package.json /build

RUN yarn install --frozen-lockfile --network-timeout 1000000

COPY ./manager/console /build

RUN yarn build

FROM golang:1.25.5-alpine3.23@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS server-builder

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

WORKDIR /go/src/d7y.io/dragonfly/v2

RUN apk --no-cache add bash make gcc libc-dev git

COPY . /go/src/d7y.io/dragonfly/v2

COPY --from=console-builder /build/dist /go/src/d7y.io/dragonfly/v2/manager/dist

RUN make build-manager-server && make install-manager

FROM ${BASE_IMAGE} AS health

ENV GRPC_HEALTH_PROBE_VERSION=v0.4.48

RUN if [ "$(uname -m)" = "ppc64le" ]; then \
    probe_arch=ppc64le; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
    probe_arch=arm64; \
    else \
    probe_arch=amd64; \
    fi && \
    probe_url=https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${probe_arch} && \
    for i in 1 2 3 4 5; do \
    wget -q -T 30 -O /bin/grpc_health_probe "${probe_url}" && break; \
    [ "${i}" -eq 5 ] && exit 1; \
    sleep 2; \
    done && \
    chmod +x /bin/grpc_health_probe

FROM ${BASE_IMAGE}

WORKDIR /opt/dragonfly

ENV PATH=/opt/dragonfly/bin:$PATH

RUN mkdir -p /opt/dragonfly/bin/manager/console \
    && echo "hosts: files dns" > /etc/nsswitch.conf

RUN apk --no-cache add curl

COPY --from=server-builder /opt/dragonfly/bin/manager /opt/dragonfly/bin/server
COPY --from=health /bin/grpc_health_probe /bin/grpc_health_probe

EXPOSE 8080 65003

ENTRYPOINT ["/opt/dragonfly/bin/server"]
