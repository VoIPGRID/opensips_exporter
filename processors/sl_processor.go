package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type SlProcessor struct {
	statistics map[string]opensips.Statistic
}

var slLabelNames = []string{}
var slMetrics = map[string]metric{
	"1xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"2xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"3xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"4xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"5xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"6xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"sent_replies":     newMetric("sl", "sent_replies_total", "The total number of sent_replies.", slLabelNames, prometheus.CounterValue),
	"sent_err_replies": newMetric("sl", "sent_err_replies_total", "The total number of sent_err_replies.", slLabelNames, prometheus.CounterValue),
	"received_ACKs":    newMetric("sl", "received_ACKs", "The number of received_ACKs.", slLabelNames, prometheus.CounterValue),
}

func init() {
	for metric := range slMetrics {
		Processors[metric] = slProcessorFunc
	}
	Processors["sl:"] = slProcessorFunc
}

func (c SlProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range slMetrics {
		ch <- metric.Desc
	}
}

func (p SlProcessor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		slMetrics["1xx_replies"].Desc,
		slMetrics["1xx_replies"].ValueType,
		p.statistics["1xx_replies"].Value,
		"1xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["2xx_replies"].Desc,
		slMetrics["2xx_replies"].ValueType,
		p.statistics["2xx_replies"].Value,
		"2xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["3xx_replies"].Desc,
		slMetrics["3xx_replies"].ValueType,
		p.statistics["3xx_replies"].Value,
		"3xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["4xx_replies"].Desc,
		slMetrics["4xx_replies"].ValueType,
		p.statistics["4xx_replies"].Value,
		"4xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["5xx_replies"].Desc,
		slMetrics["5xx_replies"].ValueType,
		p.statistics["5xx_replies"].Value,
		"5xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["6xx_replies"].Desc,
		slMetrics["6xx_replies"].ValueType,
		p.statistics["6xx_replies"].Value,
		"6xx",
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["sent_replies"].Desc,
		slMetrics["sent_replies"].ValueType,
		p.statistics["sent_replies"].Value,
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["sent_err_replies"].Desc,
		slMetrics["sent_err_replies"].ValueType,
		p.statistics["sent_err_replies"].Value,
	)
	ch <- prometheus.MustNewConstMetric(
		slMetrics["received_ACKs"].Desc,
		slMetrics["received_ACKs"].ValueType,
		p.statistics["received_ACKs"].Value,
	)
}

func slProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &SlProcessor{
		statistics: s,
	}
}
