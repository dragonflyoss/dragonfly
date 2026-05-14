ARG BASE_IMAGE=alpine:3.21.7@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d

FROM ${BASE_IMAGE}

ENV PATH=/opt/dragonfly/bin:$PATH
RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/scheduler /opt/dragonfly/bin/scheduler
COPY ./artifacts/plugins/d7y-scheduler-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8002

ENTRYPOINT ["/opt/dragonfly/bin/scheduler"]
