ARG BASE_IMAGE=alpine:3.23.7@sha256:ff7f0c727dd796aa36787bfcd9a64bf35a6ce72a73451b3359fde849b6703d63

FROM golang:1.25.5-alpine3.23@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder

ARG GOPROXY
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/d7y.io/dragonfly/v2
COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} GOPROXY=${GOPROXY} \
    go build -trimpath -o /out/injector ./cmd/injector

FROM ${BASE_IMAGE}

RUN echo "hosts: files dns" > /etc/nsswitch.conf

COPY --from=builder /out/injector /opt/dragonfly/bin/injector

USER 65532:65532
EXPOSE 9443 8081

ENTRYPOINT ["/opt/dragonfly/bin/injector"]
