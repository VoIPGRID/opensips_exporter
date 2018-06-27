# opensips_exporter 
This exporter exposes OpenSIPS metrics for consumption by Prometheus using the Unix socket 
provided by OpenSIPS. It uses the
OpenSIPS [Management Interface](http://www.opensips.org/Documentation/Interface-MI-1-11) to gather
these statistics.

Tested and developed for OpenSIPS 1.11. Though the Management Interface to gather metrics is available
in other OpenSIPS versions.

## Table of contents

- [opensips_exporter](#opensipsexporter)
    - [Table of contents](#table-of-contents)
    - [Status](#status)
    - [Usage](#usage)
        - [Exported Metrics](#exported-metrics)
        - [Processors](#processors)
            - [Filtering enabled processors](#filtering-enabled-processors)
        - [Development](#development)
    - [Contributing](#contributing)
    - [Contributors](#contributors)
    - [License](#license)

## Status

[![Go Report Card](https://goreportcard.com/badge/github.com/VoIPGRID/opensips_exporter)](https://goreportcard.com/report/github.com/VoIPGRID/opensips_exporter) 
[![Docker Pulls](https://img.shields.io/docker/pulls/voipgrid/opensips_exporter.svg?maxAge=604800)](https://hub.docker.com/r/voipgrid/opensips_exporter/)

Active / maintained

This project is considered stable for use in production environments.

## Usage

Make sure `$GOPATH/bin` is in your `$PATH`.

```text
Usage of opensips_exporter:
  -path string
        The path where metrics will be served (default "/metrics")
  -addr string
        Address on which the OpenSIPS exporter listens. (e.g. 127.0.0.1:9434) (default ":9434")
  -socket string
        Path to the socket file for OpenSIPS. (default "/var/run/ser-fg/ser.sock")
```

### Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| opensips_up | Whether the opensips exporter could read metrics from the Management Interface socket. (i.e. is OpenSIPS up) | |
| opensips_core_bad_URIs_rcvd | Number of URIs that OpenSIPS failed to parse. | |
| opensips_core_bad_msg_hdr | Number of SIP headers that OpenSIPS failed to parse. | |
| opensips_core_replies | Number of received replies by OpenSIPS. | kind|
| opensips_core_replies_total | Total number of received replies by OpenSIPS. | |
| opensips_core_request | Number of requests by OpenSIPS. | kind |
| opensips_core_requests_total | Total number of received requests by OpenSIPS. |  |
| opensips_core_unsupported_methods | Number of non-standard methods encountered by OpenSIPS while parsing SIP methods. | |
| opensips_core_uptime_seconds | Number of seconds elapsed from OpenSIPS starting. | |
| opensips_dialog_dialogs | Number of dialogs. | status |
| opensips_dialog_received | The number of dialog events received from other OpenSIPS instances. | event |
| opensips_dialog_sent | Number of replicated dialog requests send to other OpenSIPS instances. | event |
| opensips_load_load | Percentage of UDP children that are awake and processing SIP messages on the specific UDP interface. |ip, port, protocol|
| opensips_load_tcp_load | Percentage of TCP children that are awake and processing SIP messages. | | 
| opensips_net_waiting | Number of bytes waiting to be consumed on an interface that OpenSIPS is listening on. | protocol |
| opensips_pkmem_fragments | Currently available number of free fragments in the private memory for OpenSIPS process. | pid |
| opensips_pkmem_free_size | Free private memory available for the OpenSIPS process. Computed as total_size - real_used_size. | pid |
| opensips_pkmem_max_used_size | The maximum amount of private memory ever used by the OpenSIPS process. | pid |
| opensips_pkmem_real_used_size | Amount of private memory requested by the OpenSIPS process, including allocator-specific metadata. | pid |
| opensips_pkmem_total_size | Total size of private memory available to the OpenSIPS process. | pid |
| opensips_pkmem_used_size | Amount of private memory requested and used by the OpenSIPS process. | pid |
| opensips_registrar_default_expire | Value of default_expire parameter. | |
| opensips_registrar_max_contacts | Value of max_contacts parameter. | |
| opensips_registrar_max_expires | Value of max_expires parameter. | |
| opensips_registrar_registrations | Number of registrations. | type |
| opensips_shmem_fragments | Total number of fragments in the shared memory. | |
| opensips_shmem_free_size | Free memory available. Computed as total_size - real_used_size | |
| opensips_shmem_max_used_size | Maximum amount of shared memory ever used by OpenSIPS processes. | |
| opensips_shmem_real_used_size | Amount of shared memory requested by OpenSIPS processes + malloc overhead | |
| opensips_shmem_total_size | Total size of shared memory available to OpenSIPS processes. | |
| opensips_shmem_used_size | Amount of shared memory requested and used by OpenSIPS processes. | |
| opensips_sl_received_ACKs | The number of received_ACKs. | |
| opensips_sl_replies | The number of replies. | type |
| opensips_sl_xxx_replies | The number of replies that don't match any other reply status. | |
| opensips_sl_failures | The number of failures. | |
| opensips_sl_sent_err_replies_total | The total number of sent_err_replies. | |
| opensips_sl_sent_replies_total | The total number of sent_replies. | |
| opensips_tm_inuse_transactions | Number of transactions existing in memory at current time. | |
| opensips_tm_local_replies_total | Total number of replies local generated by TM module. | |
| opensips_tm_received_replies_total | Total number of total replies received by TM module. | |
| opensips_tm_relayed_replies_total | Total number of replies received and relayed by TM module. | |
| opensips_tm_transactions_total | Total number of transactions. (TM module)| type |
| opensips_tmx_transactions_total | Total number of transactions. (TMX module) | type |
| opensips_tmx_UAS_transactions | Total number of transactions created by received requests. | type |
| opensips_tmx_UAC_transactions | Total number of transactions created by local generated requests. | |
| opensips_tmx_UAC_transactions | Total number of transactions created by local generated requests. | |
| opensips_tmx_inuse_transactions | Number of transactions existing in memory at current time. | |
| opensips_tmx_active_transactions | Number of ongoing transactions at current time. | |
| opensips_tmx_replies | Total number of replies. | type |
| opensips_uri_negative_checks | Amount of negative URI checks. | |
| opensips_uri_positive_checks | Amount of positive URI checks. | |
| opensips_userloc_registered_users_total | Total number of AOR existing in the USRLOC memory cache for all domains. | |
| opensips_usrloc_contacts | Number of contacts existing in the USRLOC memory cache for that domain. | domain |
| opensips_usrloc_expires | Total number of expired contacts for that domain. | domain |
| opensips_usrloc_users | Number of AOR existing in the USRLOC memory cache for that domain. | domain |

### Processors

There are processors available per 'module' of OpenSIPS. The processors take the
statistics from the OpenSIPS socket and turn it into a Prometheus metric. You can 
recognise a processor by the naming convention of the metrics. 
For example the`opensips_core_replies` metric comes from the core module and its 
processor can be found in `./processors/core_processor`.

You can find out more about the available modules in the OpenSIPS documentation.
At this time not all modules have a processor written for them. (help wanted!)

#### Filtering enabled processors

It is possible to select what processors you want metrics from. You can do this by
appending `collect[]` parameters to your request. If for example you only want to
get metrics about the [core](http://www.opensips.org/Documentation/Interface-CoreStatistics-1-11) 
and [usrloc](http://www.opensips.org/html/docs/modules/1.11.x/usrloc.html#idp5699792) 
module you can do this as follows:

```bash
curl localhost:9434/metrics?collect[]=core:&collect[]=usrloc:
```

**_Note: You have to append `:` to the module name for this to work._**

### Development

To work on opensips_exporter, get a recent [Go], get a recent [dep], and
run:

    go get -u github.com/VoIPGRID/opensips_exporter

While developing, make sure to run `dep ensure` often enough to keep
dependencies up-to-date.

The `github.com/VoIPGRID/opensips_exporter/opensips` package contains the
implementation of the interactions with OpenSIPS needed to get statistics from
the mi_datagram Unix socket of a running OpenSIPS. For tests, there is a mock
in the `./internal/mock` package.

Metrics from different OpenSIPS modules are extracted by processors defined in
the `./processors` package. To extend this exporter with metrics from other modules
create your own processor and implement the `Collector` interface. See the other
processors for inspiration. 

## Contributing

See the [CONTRIBUTING.md](.github/CONTRIBUTING.md) file on how to contribute to this project.

## Contributors

See the [CONTRIBUTORS.md](CONTRIBUTORS.md) file for a list of contributors to the project.

[Go]: https://golang.org/doc/install (Getting Started - The Go Programming Language)
[dep]: https://golang.github.io/dep/docs/installation.html (Installation · dep)

## License

opensips_exporter is made available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
