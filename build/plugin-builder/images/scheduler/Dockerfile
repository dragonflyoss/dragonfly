ARG BASE_IMAGE=alpine:3.17@sha256:8fc3dacfb6d69da8d44e42390de777e48577085db99aa4e4af35f483eb08b989

FROM ${BASE_IMAGE}

ENV PATH=/opt/dragonfly/bin:$PATH
RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/scheduler /opt/dragonfly/bin/scheduler
COPY ./artifacts/plugins/d7y-scheduler-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8002

ENTRYPOINT ["/opt/dragonfly/bin/scheduler"]
