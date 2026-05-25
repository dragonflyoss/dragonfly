ARG BASE_IMAGE=alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

FROM ${BASE_IMAGE}

ENV PATH=/opt/dragonfly/bin:$PATH
RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/scheduler /opt/dragonfly/bin/scheduler
COPY ./artifacts/plugins/d7y-scheduler-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8002

ENTRYPOINT ["/opt/dragonfly/bin/scheduler"]
