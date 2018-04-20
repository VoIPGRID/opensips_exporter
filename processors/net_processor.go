package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type NetProcessor struct {
	statistics map[string]opensips.Statistic
}

var netLabelNames = []string{"protocol"}
var netMetrics = map[string]metric{
	"waiting_udp": newMetric("net", "waiting", "Number of bytes waiting to be consumed on an interface that OpenSIPS is listening on.", netLabelNames, prometheus.GaugeValue),
	"waiting_tcp": newMetric("net", "waiting", "Number of bytes waiting to be consumed on an interface that OpenSIPS is listening on.", netLabelNames, prometheus.GaugeValue),
	"waiting_tls": newMetric("net", "waiting", "Number of bytes waiting to be consumed on an interface that OpenSIPS is listening on.", netLabelNames, prometheus.GaugeValue),
}

func init() {
	for metric := range netMetrics {
		Processors[metric] = netProcessorFunc
	}
	Processors["net:"] = netProcessorFunc
}

func (c NetProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range netMetrics {
		ch <- metric.Desc
	}
}

func (p NetProcessor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		netMetrics["waiting_udp"].Desc,
		netMetrics["waiting_udp"].ValueType,
		p.statistics["waiting_udp"].Value,
		"udp",
	)
	ch <- prometheus.MustNewConstMetric(
		netMetrics["waiting_tcp"].Desc,
		netMetrics["waiting_tcp"].ValueType,
		p.statistics["waiting_tcp"].Value,
		"tcp",
	)
	ch <- prometheus.MustNewConstMetric(
		netMetrics["waiting_tls"].Desc,
		netMetrics["waiting_tls"].ValueType,
		p.statistics["waiting_tls"].Value,
		"tls",
	)

}

func netProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &NetProcessor{
		statistics: s,
	}
}
