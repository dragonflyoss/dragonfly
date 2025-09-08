FROM quay.io/iovisor/bpftrace:v0.13.0

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update \
    && apt-get install -y binutils bpfcc-tools curl dnsutils git iproute2 iputils-ping jq netcat socat tree vim wget \
    && apt-get install -y openssh-client openssh-server \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /run/sshd

ENV PATH=/go/bin:/usr/local/go/bin:$PATH

ENV GOPATH=/go
ENV GOLANG_VERSION 1.25.0
ENV GOLANG_DOWNLOAD_URL https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 3c9a05d0c5a4c4a66b2c5a21f9f0b8e5e5b5c5b5c5b5c5b5c5b5c5b5c5b5c5b5c5b

ARG GOPROXY

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz \
    && go install github.com/go-delve/delve/cmd/dlv@v1.7.0
