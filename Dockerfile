FROM        quay.io/prometheus/busybox:latest
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY opensips_exporter /bin/opensips_exporter

ENTRYPOINT  ["/bin/opensips_exporter"]
USER        nobody
EXPOSE      9434
