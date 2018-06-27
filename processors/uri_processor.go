package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// uriProcessor metrics related to SIP URI processing.
// doc: http://www.opensips.org/html/docs/modules/1.11.x/uri.html
// src: https://github.com/OpenSIPS/opensips/blob/1.11/modules/uri/uri_mod.c#L191
type uriProcessor struct {
	statistics map[string]opensips.Statistic
}

var uriLabelNames = []string{}
var uriMetrics = map[string]metric{
	"positive":        newMetric("uri", "positive_checks", "Amount of positive URI checks.", uriLabelNames, prometheus.CounterValue),
	"negative_checks": newMetric("uri", "negative_checks", "Amount of negative URI checks.", uriLabelNames, prometheus.CounterValue),
}

func init() {
	for metric := range uriMetrics {
		OpensipsProcessors[metric] = uriProcessorFunc
	}
	OpensipsProcessors["uri:"] = uriProcessorFunc
}

// Describe implements prometheus.Collector.
func (p uriProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range uriMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p uriProcessor) Collect(ch chan<- prometheus.Metric) {
	for _, s := range p.statistics {
		if s.Module == "uri" {
			switch s.Name {
			case "positive":
				ch <- prometheus.MustNewConstMetric(
					uriMetrics["positive"].Desc,
					uriMetrics["positive"].ValueType,
					s.Value,
				)
			case "negative_checks":
				ch <- prometheus.MustNewConstMetric(
					uriMetrics["negative_checks"].Desc,
					uriMetrics["negative_checks"].ValueType,
					s.Value,
				)
			}
		}
	}
}

func uriProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &uriProcessor{
		statistics: s,
	}
}
