ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/minecraft_exporter /bin/minecraft_exporter

EXPOSE      9940
USER        nobody
ENTRYPOINT  [ "/bin/minecraft_exporter" ]
