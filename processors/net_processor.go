package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// netProcessor provides metrics about network packets.
// doc: http://www.opensips.org/Documentation/Interface-CoreStatistics-1-11#toc17
// src: https://github.com/OpenSIPS/opensips/blob/b917c70ba8d5797dc6364aaf702c3415539be52a/core_stats.c#L95
type netProcessor struct {
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
		OpensipsProcessors[metric] = netProcessorFunc
	}
	OpensipsProcessors["net:"] = netProcessorFunc
}

// Describe implements prometheus.Collector.
func (p netProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range netMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p netProcessor) Collect(ch chan<- prometheus.Metric) {
	for _, s := range p.statistics {
		if s.Module == "net" {
			switch s.Name {
			case "waiting_udp":
				ch <- prometheus.MustNewConstMetric(
					netMetrics["waiting_udp"].Desc,
					netMetrics["waiting_udp"].ValueType,
					s.Value,
					"udp",
				)
			case "waiting_tcp":
				ch <- prometheus.MustNewConstMetric(
					netMetrics["waiting_tcp"].Desc,
					netMetrics["waiting_tcp"].ValueType,
					s.Value,
					"tcp",
				)
			case "waiting_tls":
				ch <- prometheus.MustNewConstMetric(
					netMetrics["waiting_tls"].Desc,
					netMetrics["waiting_tls"].ValueType,
					s.Value,
					"tls",
				)
			}
		}
	}
}

func netProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &netProcessor{
		statistics: s,
	}
}
