package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type UriProcessor struct {
	statistics map[string]opensips.Statistic
}

var uriLabelNames = []string{}
var uriMetrics = map[string]metric{
	"positive checks": newMetric("uri", "positive_checks", "Amount of positive URI checks.", uriLabelNames, prometheus.CounterValue),
	"positive_checks": newMetric("uri", "positive_checks", "Amount of positive URI checks.", uriLabelNames, prometheus.CounterValue),
	"negative_checks": newMetric("uri", "negative_checks", "Amount of negative URI checks.", uriLabelNames, prometheus.CounterValue),
}

func init() {
	for metric := range uriMetrics {
		Processors[metric] = uriProcessorFunc
	}
	Processors["uri:"] = uriProcessorFunc
}

func (c UriProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range uriMetrics {
		ch <- metric.Desc
	}
}

func (p UriProcessor) Collect(ch chan<- prometheus.Metric) {
	for key, metric := range uriMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.ValueType,
			p.statistics[key].Value,
		)
	}
}

func uriProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &UriProcessor{
		statistics: s,
	}
}
