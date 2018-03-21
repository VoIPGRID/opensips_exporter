# opensips_exporter

Work In Progress

Currently, this opensips_exporter doesn't work. More documentation to come.

## Development

To work on opensips_exporter, get a recent [Go], get a recent [dep], and
run:

    go get -u github.com/VoIPGRID/opensips_exporter

While developing, make sure to run `dep ensure` often enough to keep
dependencies up-to-date.

The `github.com/VoIPGRID/opensips_exporter/opensips` package contains the
implementation of the interactions with OpenSIPS needed to get statistics from
the mi_datagram Unix socket of a running OpenSIPS. For tests, there is a mock
in the `./internal/mock` package.

[Go]: https://golang.org/doc/install (Getting Started - The Go Programming Language)
[dep]: https://golang.github.io/dep/docs/installation.html (Installation · dep)
