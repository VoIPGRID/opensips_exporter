package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// registrarProcessor provides metrcs about SIP registrations
// doc: http://www.opensips.org/html/docs/modules/1.11.x/registrar.html#idp5702944
// src: https://github.com/OpenSIPS/opensips/blob/1.11/modules/registrar/reg_mod.c#L202
type registrarProcessor struct {
	statistics map[string]opensips.Statistic
}

var registrarLabelNames = []string{}
var registrarMetrics = map[string]metric{
	"max_expires":    newMetric("registrar", "max_expires", "Value of max_expires parameter.", registrarLabelNames, prometheus.GaugeValue),
	"max_contacts":   newMetric("registrar", "max_contacts", "Value of max_contacts parameter.", registrarLabelNames, prometheus.GaugeValue),
	"default_expire": newMetric("registrar", "default_expire", "Value of default_expire parameter.", registrarLabelNames, prometheus.GaugeValue),
	"accepted_regs":  newMetric("registrar", "registrations", "Number of registrations.", []string{"type"}, prometheus.CounterValue),
	"rejected_regs":  newMetric("registrar", "registrations", "Number of registrations.", []string{"type"}, prometheus.CounterValue),
}

func init() {
	for metric := range registrarMetrics {
		OpensipsProcessors[metric] = registrarProcessorFunc
	}
	OpensipsProcessors["registrar:"] = registrarProcessorFunc
}

// Describe implements prometheus.Collector.
func (p registrarProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range registrarMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p registrarProcessor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		registrarMetrics["max_expires"].Desc,
		registrarMetrics["max_expires"].ValueType,
		p.statistics["max_expires"].Value,
	)
	ch <- prometheus.MustNewConstMetric(
		registrarMetrics["max_contacts"].Desc,
		registrarMetrics["max_contacts"].ValueType,
		p.statistics["max_contacts"].Value,
	)
	ch <- prometheus.MustNewConstMetric(
		registrarMetrics["default_expire"].Desc,
		registrarMetrics["default_expire"].ValueType,
		p.statistics["default_expire"].Value,
	)
	ch <- prometheus.MustNewConstMetric(
		registrarMetrics["accepted_regs"].Desc,
		registrarMetrics["accepted_regs"].ValueType,
		p.statistics["accepted_regs"].Value,
		"accepted",
	)
	ch <- prometheus.MustNewConstMetric(
		registrarMetrics["rejected_regs"].Desc,
		registrarMetrics["rejected_regs"].ValueType,
		p.statistics["rejected_regs"].Value,
		"rejected",
	)
}

func registrarProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &registrarProcessor{
		statistics: s,
	}
}
