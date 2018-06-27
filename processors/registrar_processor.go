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
	for _, s := range p.statistics {
		if s.Module == "registrar" {
			switch s.Name {
			case "max_expires":
				ch <- prometheus.MustNewConstMetric(
					registrarMetrics["max_expires"].Desc,
					registrarMetrics["max_expires"].ValueType,
					s.Value,
				)
			case "max_contacts":
				ch <- prometheus.MustNewConstMetric(
					registrarMetrics["max_contacts"].Desc,
					registrarMetrics["max_contacts"].ValueType,
					s.Value,
				)
			case "default_expire":
				ch <- prometheus.MustNewConstMetric(
					registrarMetrics["default_expire"].Desc,
					registrarMetrics["default_expire"].ValueType,
					s.Value,
				)
			case "accepted_regs":
				ch <- prometheus.MustNewConstMetric(
					registrarMetrics["accepted_regs"].Desc,
					registrarMetrics["accepted_regs"].ValueType,
					s.Value,
					"accepted",
				)
			case "rejected_regs":
				ch <- prometheus.MustNewConstMetric(
					registrarMetrics["rejected_regs"].Desc,
					registrarMetrics["rejected_regs"].ValueType,
					s.Value,
					"rejected",
				)
			}
		}
	}
}

func registrarProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &registrarProcessor{
		statistics: s,
	}
}
