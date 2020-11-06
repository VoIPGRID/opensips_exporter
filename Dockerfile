FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Ruben Homs <ruben.homs@wearespindle.com>

COPY opensips_exporter /bin/opensips_exporter

ENTRYPOINT  ["/bin/opensips_exporter"]
USER        nobody
EXPOSE      9434
