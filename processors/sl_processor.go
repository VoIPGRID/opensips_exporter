package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// slProcessor provides metrics for its stateless UA
// doc: http://www.opensips.org/html/docs/modules/1.11.x/sl.html#idp158896
// src: https://github.com/OpenSIPS/opensips/blob/1.11/modules/sl/sl.c#L91
type slProcessor struct {
	statistics map[string]opensips.Statistic
}

var slLabelNames = []string{}
var slMetrics = map[string]metric{
	"xxx_replies":      newMetric("sl", "xxx_replies", "The number of replies that don't match any other reply status.", slLabelNames, prometheus.CounterValue),
	"1xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"2xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"200_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"202_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"3xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"300_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"301_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"302_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"4xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"400_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"401_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"403_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"404_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"407_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"408_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"483_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"5xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"500_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"6xx_replies":      newMetric("sl", "replies", "The number of replies.", []string{"type"}, prometheus.CounterValue),
	"sent_replies":     newMetric("sl", "sent_replies_total", "The total number of sent_replies.", slLabelNames, prometheus.CounterValue),
	"sent_err_replies": newMetric("sl", "sent_err_replies_total", "The total number of sent_err_replies.", slLabelNames, prometheus.CounterValue),
	"received_ACKs":    newMetric("sl", "received_ACKs", "The number of received_ACKs.", slLabelNames, prometheus.CounterValue),
	"failures":         newMetric("sl", "failures", "The number of failures.", slLabelNames, prometheus.CounterValue),
}

func init() {
	for metric := range slMetrics {
		OpensipsProcessors[metric] = slProcessorFunc
	}
	OpensipsProcessors["sl:"] = slProcessorFunc
}

// Describe implements prometheus.Collector.
func (p slProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range slMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p slProcessor) Collect(ch chan<- prometheus.Metric) {
	for _, s := range p.statistics {
		if s.Module == "sl" {
			switch s.Name {
			case "xxx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["xxx_replies"].Desc,
					slMetrics["xxx_replies"].ValueType,
					s.Value,
				)
			case "1xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["1xx_replies"].Desc,
					slMetrics["1xx_replies"].ValueType,
					s.Value,
					"1xx",
				)
			case "2xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["2xx_replies"].Desc,
					slMetrics["2xx_replies"].ValueType,
					s.Value,
					"2xx",
				)
			case "200_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["200_replies"].Desc,
					slMetrics["200_replies"].ValueType,
					s.Value,
					"200",
				)
			case "202_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["202_replies"].Desc,
					slMetrics["202_replies"].ValueType,
					s.Value,
					"202",
				)
			case "3xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["3xx_replies"].Desc,
					slMetrics["3xx_replies"].ValueType,
					s.Value,
					"3xx",
				)
			case "300_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["300_replies"].Desc,
					slMetrics["300_replies"].ValueType,
					s.Value,
					"300",
				)
			case "301_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["301_replies"].Desc,
					slMetrics["301_replies"].ValueType,
					s.Value,
					"301",
				)
			case "302_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["302_replies"].Desc,
					slMetrics["302_replies"].ValueType,
					s.Value,
					"302",
				)
			case "4xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["4xx_replies"].Desc,
					slMetrics["4xx_replies"].ValueType,
					s.Value,
					"4xx",
				)
			case "400_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["400_replies"].Desc,
					slMetrics["400_replies"].ValueType,
					s.Value,
					"400",
				)
			case "401_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["401_replies"].Desc,
					slMetrics["401_replies"].ValueType,
					s.Value,
					"401",
				)
			case "403_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["403_replies"].Desc,
					slMetrics["403_replies"].ValueType,
					s.Value,
					"403",
				)
			case "404_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["404_replies"].Desc,
					slMetrics["404_replies"].ValueType,
					s.Value,
					"404",
				)
			case "407_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["407_replies"].Desc,
					slMetrics["407_replies"].ValueType,
					s.Value,
					"407",
				)
			case "408_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["408_replies"].Desc,
					slMetrics["408_replies"].ValueType,
					s.Value,
					"408",
				)
			case "483_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["483_replies"].Desc,
					slMetrics["483_replies"].ValueType,
					s.Value,
					"483",
				)
			case "5xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["5xx_replies"].Desc,
					slMetrics["5xx_replies"].ValueType,
					s.Value,
					"5xx",
				)
			case "500_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["500_replies"].Desc,
					slMetrics["500_replies"].ValueType,
					s.Value,
					"500",
				)
			case "6xx_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["6xx_replies"].Desc,
					slMetrics["6xx_replies"].ValueType,
					s.Value,
					"6xx",
				)
			case "sent_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["sent_replies"].Desc,
					slMetrics["sent_replies"].ValueType,
					s.Value,
				)
			case "sent_err_replies":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["sent_err_replies"].Desc,
					slMetrics["sent_err_replies"].ValueType,
					s.Value,
				)
			case "received_ACKs":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["received_ACKs"].Desc,
					slMetrics["received_ACKs"].ValueType,
					s.Value,
				)
			case "failures":
				ch <- prometheus.MustNewConstMetric(
					slMetrics["failures"].Desc,
					slMetrics["failures"].ValueType,
					s.Value,
				)
			}
		}
	}
}

func slProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &slProcessor{
		statistics: s,
	}
}
