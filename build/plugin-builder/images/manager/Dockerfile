ARG BASE_IMAGE=alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

FROM ${BASE_IMAGE}

WORKDIR /opt/dragonfly

ENV PATH=/opt/dragonfly/bin:$PATH

RUN echo "hosts: files dns" > /etc/nsswitch.conf  && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/manager /opt/dragonfly/bin/server
COPY ./artifacts/plugins/d7y-manager-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8080 65003

ENTRYPOINT ["/opt/dragonfly/bin/server"]
