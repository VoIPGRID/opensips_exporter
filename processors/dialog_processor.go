package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// dialogProcessor exposes metrics about SIP dialogs.
// doc: http://www.opensips.org/html/docs/modules/1.11.x/dialog.html#idp5859728
// src: https://github.com/OpenSIPS/opensips/blob/1.11/modules/dialog/dialog.c#L283
type dialogProcessor struct {
	statistics map[string]opensips.Statistic
}

var dialogLabelNames = []string{}
var dialogMetrics = map[string]metric{
	"active_dialogs":    newMetric("dialog", "dialogs", "Number of dialogs.", []string{"status"}, prometheus.GaugeValue),
	"early_dialogs":     newMetric("dialog", "dialogs", "Number of dialogs.", []string{"status"}, prometheus.GaugeValue),
	"processed_dialogs": newMetric("dialog", "dialogs", "Number of dialogs.", []string{"status"}, prometheus.GaugeValue),
	"expired_dialogs":   newMetric("dialog", "dialogs", "Number of dialogs.", []string{"status"}, prometheus.GaugeValue),
	"failed_dialogs":    newMetric("dialog", "dialogs", "Number of dialogs.", []string{"status"}, prometheus.GaugeValue),
	"create_sent":       newMetric("dialog", "sent", "Number of replicated dialog requests send to other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
	"update_sent":       newMetric("dialog", "sent", "Number of replicated dialog requests send to other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
	"delete_sent":       newMetric("dialog", "sent", "Number of replicated dialog requests send to other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
	"create_rcv":        newMetric("dialog", "received", "The number of dialog events received from other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
	"update_rcv":        newMetric("dialog", "received", "The number of dialog events received from other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
	"delete_rcv":        newMetric("dialog", "received", "The number of dialog events received from other OpenSIPS instances.", []string{"event"}, prometheus.CounterValue),
}

func init() {
	for metric := range dialogMetrics {
		OpensipsProcessors[metric] = dialogProcessorFunc
	}
	OpensipsProcessors["dialog:"] = dialogProcessorFunc
}

// Describe implements prometheus.Collector.
func (p dialogProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range dialogMetrics {
		ch <- metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p dialogProcessor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["active_dialogs"].Desc,
		dialogMetrics["active_dialogs"].ValueType,
		p.statistics["active_dialogs"].Value,
		"active",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["early_dialogs"].Desc,
		dialogMetrics["early_dialogs"].ValueType,
		p.statistics["early_dialogs"].Value,
		"early",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["processed_dialogs"].Desc,
		dialogMetrics["processed_dialogs"].ValueType,
		p.statistics["processed_dialogs"].Value,
		"processed",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["expired_dialogs"].Desc,
		dialogMetrics["expired_dialogs"].ValueType,
		p.statistics["expired_dialogs"].Value,
		"expired",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["failed_dialogs"].Desc,
		dialogMetrics["failed_dialogs"].ValueType,
		p.statistics["failed_dialogs"].Value,
		"failed",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["create_sent"].Desc,
		dialogMetrics["create_sent"].ValueType,
		p.statistics["create_sent"].Value,
		"create",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["update_sent"].Desc,
		dialogMetrics["update_sent"].ValueType,
		p.statistics["update_sent"].Value,
		"update",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["delete_sent"].Desc,
		dialogMetrics["delete_sent"].ValueType,
		p.statistics["delete_sent"].Value,
		"delete",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["create_rcv"].Desc,
		dialogMetrics["create_rcv"].ValueType,
		p.statistics["create_rcv"].Value,
		"create",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["update_rcv"].Desc,
		dialogMetrics["update_rcv"].ValueType,
		p.statistics["update_rcv"].Value,
		"update",
	)
	ch <- prometheus.MustNewConstMetric(
		dialogMetrics["delete_rcv"].Desc,
		dialogMetrics["delete_rcv"].ValueType,
		p.statistics["delete_rcv"].Value,
		"delete",
	)
}

func dialogProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &dialogProcessor{
		statistics: s,
	}
}
