ARG BASE_IMAGE=alpine:3.21@sha256:5405e8f36ce1878720f71217d664aa3dea32e5e5df11acbf07fc78ef5661465b

FROM ${BASE_IMAGE}

WORKDIR /opt/dragonfly

ENV PATH=/opt/dragonfly/bin:$PATH

RUN echo "hosts: files dns" > /etc/nsswitch.conf  && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/manager /opt/dragonfly/bin/server
COPY ./artifacts/plugins/d7y-manager-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8080 65003

ENTRYPOINT ["/opt/dragonfly/bin/server"]
