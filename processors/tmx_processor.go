package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// tmxProcessor exposes metrics for stateful processing of SIP transactions.
// doc: https://kamailio.org/docs/modules/4.4.x/modules/tmx.html#idp23886596
type tmxProcessor struct {
	statistics map[string]opensips.Statistic
}

var tmxLabelNames = []string{}
var tmxMetrics = map[string]metric{
	"UAS_transactions":    newMetric("tmx", "UAS_transactions", "Total number of transactions created by received requests.", tmxLabelNames, prometheus.CounterValue),
	"UAC_transactions":    newMetric("tmx", "UAC_transactions", "Total number of transactions created by local generated requests.", tmxLabelNames, prometheus.CounterValue),
	"2xx_transactions":    newMetric("tmx", "transactions_total", "Total number of transactions.", []string{"type"}, prometheus.CounterValue),
	"3xx_transactions":    newMetric("tmx", "transactions_total", "Total number of transactions.", []string{"type"}, prometheus.CounterValue),
	"4xx_transactions":    newMetric("tmx", "transactions_total", "Total number of transactions.", []string{"type"}, prometheus.CounterValue),
	"5xx_transactions":    newMetric("tmx", "transactions_total", "Total number of transactions.", []string{"type"}, prometheus.CounterValue),
	"6xx_transactions":    newMetric("tmx", "transactions_total", "Total number of transactions.", []string{"type"}, prometheus.CounterValue),
	"inuse_transactions":  newMetric("tmx", "inuse_transactions", "Number of transactions existing in memory at current time.", tmxLabelNames, prometheus.GaugeValue),
	"active_transactions": newMetric("tmx", "active_transactions", "Number of ongoing transactions at current time.", tmxLabelNames, prometheus.GaugeValue),
	"rpl_received":        newMetric("tmx", "replies", "Total number of replies.", []string{"type"}, prometheus.CounterValue),
	"rpl_absorbed":        newMetric("tmx", "replies", "Total number of replies.", []string{"type"}, prometheus.CounterValue),
	"rpl_relayed":         newMetric("tmx", "replies", "Total number of replies.", []string{"type"}, prometheus.CounterValue),
	"rpl_generated":       newMetric("tmx", "replies", "Total number of replies.", []string{"type"}, prometheus.CounterValue),
	"rpl_sent":            newMetric("tmx", "replies", "Total number of replies.", []string{"type"}, prometheus.CounterValue),
}

func init() {
	for metric := range tmxMetrics {
		OpensipsProcessors[metric] = tmxProcessorFunc
	}
	OpensipsProcessors["tmx:"] = tmxProcessorFunc
}

// Describe implements prometheus.Collector.
func (p tmxProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range tmxMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p tmxProcessor) Collect(ch chan<- prometheus.Metric) {
	for _, s := range p.statistics {
		if s.Module == "tmx" {
			switch s.Name {
			case "UAS_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["UAS_transactions"].Desc,
					tmxMetrics["UAS_transactions"].ValueType,
					s.Value,
				)
			case "UAC_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["UAC_transactions"].Desc,
					tmxMetrics["UAC_transactions"].ValueType,
					s.Value,
				)
			case "2xx_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["2xx_transactions"].Desc,
					tmxMetrics["2xx_transactions"].ValueType,
					s.Value,
					"2xx",
				)
			case "3xx_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["3xx_transactions"].Desc,
					tmxMetrics["3xx_transactions"].ValueType,
					s.Value,
					"3xx",
				)
			case "4xx_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["4xx_transactions"].Desc,
					tmxMetrics["4xx_transactions"].ValueType,
					s.Value,
					"4xx",
				)
			case "5xx_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["5xx_transactions"].Desc,
					tmxMetrics["5xx_transactions"].ValueType,
					s.Value,
					"5xx",
				)
			case "6xx_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["6xx_transactions"].Desc,
					tmxMetrics["6xx_transactions"].ValueType,
					s.Value,
					"6xx",
				)
			case "inuse_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["inuse_transactions"].Desc,
					tmxMetrics["inuse_transactions"].ValueType,
					s.Value,
				)
			case "active_transactions":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["active_transactions"].Desc,
					tmxMetrics["active_transactions"].ValueType,
					s.Value,
				)
			case "rpl_received":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["rpl_received"].Desc,
					tmxMetrics["rpl_received"].ValueType,
					s.Value,
					"received",
				)
			case "rpl_absorbed":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["rpl_absorbed"].Desc,
					tmxMetrics["rpl_absorbed"].ValueType,
					s.Value,
					"absorbed",
				)
			case "rpl_relayed":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["rpl_relayed"].Desc,
					tmxMetrics["rpl_relayed"].ValueType,
					s.Value,
					"relayed",
				)
			case "rpl_generated":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["rpl_generated"].Desc,
					tmxMetrics["rpl_generated"].ValueType,
					s.Value,
					"generated",
				)
			case "rpl_sent":
				ch <- prometheus.MustNewConstMetric(
					tmxMetrics["rpl_sent"].Desc,
					tmxMetrics["rpl_sent"].ValueType,
					s.Value,
					"sent",
				)
			}
		}
	}
}

func tmxProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &tmxProcessor{
		statistics: s,
	}
}
