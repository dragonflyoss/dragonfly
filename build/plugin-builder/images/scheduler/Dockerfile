ARG BASE_IMAGE=alpine:3.23.7@sha256:ff7f0c727dd796aa36787bfcd9a64bf35a6ce72a73451b3359fde849b6703d63

FROM ${BASE_IMAGE}

ENV PATH=/opt/dragonfly/bin:$PATH
RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /usr/local/dragonfly/plugins/

COPY ./artifacts/binaries/scheduler /opt/dragonfly/bin/scheduler
COPY ./artifacts/plugins/d7y-scheduler-plugin-* /usr/local/dragonfly/plugins/

EXPOSE 8002

ENTRYPOINT ["/opt/dragonfly/bin/scheduler"]
