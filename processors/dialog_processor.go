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
	for _, s := range p.statistics {
		if s.Module == "dialog" {
			switch s.Name {
			case "active_dialogs":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["active_dialogs"].Desc,
					dialogMetrics["active_dialogs"].ValueType,
					s.Value,
					"active",
				)
			case "early_dialogs":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["early_dialogs"].Desc,
					dialogMetrics["early_dialogs"].ValueType,
					s.Value,
					"early",
				)
			case "processed_dialogs":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["processed_dialogs"].Desc,
					dialogMetrics["processed_dialogs"].ValueType,
					s.Value,
					"processed",
				)
			case "expired_dialogs":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["expired_dialogs"].Desc,
					dialogMetrics["expired_dialogs"].ValueType,
					s.Value,
					"expired",
				)
			case "failed_dialogs":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["failed_dialogs"].Desc,
					dialogMetrics["failed_dialogs"].ValueType,
					s.Value,
					"failed",
				)
			case "create_sent":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["create_sent"].Desc,
					dialogMetrics["create_sent"].ValueType,
					s.Value,
					"create",
				)
			case "update_sent":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["update_sent"].Desc,
					dialogMetrics["update_sent"].ValueType,
					s.Value,
					"update",
				)
			case "delete_sent":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["delete_sent"].Desc,
					dialogMetrics["delete_sent"].ValueType,
					s.Value,
					"delete",
				)
			case "create_rcv":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["create_rcv"].Desc,
					dialogMetrics["create_rcv"].ValueType,
					s.Value,
					"create",
				)
			case "update_rcv":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["update_rcv"].Desc,
					dialogMetrics["update_rcv"].ValueType,
					s.Value,
					"update",
				)
			case "delete_rcv":
				ch <- prometheus.MustNewConstMetric(
					dialogMetrics["delete_rcv"].Desc,
					dialogMetrics["delete_rcv"].ValueType,
					s.Value,
					"delete",
				)
			}
		}
	}
}

func dialogProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &dialogProcessor{
		statistics: s,
	}
}
