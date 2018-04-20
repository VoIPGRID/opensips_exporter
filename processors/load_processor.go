package processors

import (
	"strings"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type LoadProcessor struct {
	statistics map[string]opensips.Statistic
}

type loadMetric struct {
	metric   metric
	ip       string
	port     string
	protocol string
}

func init() {
	Processors["load:"] = loadProcessorFunc
	Processors["tcp-load"] = loadProcessorFunc
	Processors["load"] = loadProcessorFunc
}

func (p LoadProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.loadMetrics() {
		ch <- m.metric.Desc
	}
}

func (p LoadProcessor) Collect(ch chan<- prometheus.Metric) {
	for key, u := range p.loadMetrics() {
		if u.ip != "" {
			ch <- prometheus.MustNewConstMetric(
				u.metric.Desc,
				u.metric.ValueType,
				p.statistics[key].Value,
				u.ip, u.port, u.protocol,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				u.metric.Desc,
				u.metric.ValueType,
				p.statistics[key].Value,
			)
		}
	}
}

func loadProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &LoadProcessor{
		statistics: s,
	}
}

func (p LoadProcessor) loadMetrics() map[string]loadMetric {
	var metrics = map[string]loadMetric{}

	var stats []opensips.Statistic
	for _, s := range p.statistics {
		if s.Module == "load" {
			stats = append(stats, s)
		}
	}

	for _, s := range stats {

		if !strings.Contains(s.Name, ":") {
			continue
		}

		split := strings.Split(s.Name, ":")
		protocol := split[0]
		ip := split[1]
		port := split[2]
		port = strings.Trim(port, "-load")

		metric := newMetric("load", "load", "Percentage of UDP children that are awake and processing SIP messages on the specific UDP interface.", []string{"ip", "port", "protocol"}, prometheus.GaugeValue)
		metrics[s.Name] = loadMetric{
			metric:   metric,
			ip:       ip,
			port:     port,
			protocol: protocol,
		}
	}

	metrics["tcp-load"] = loadMetric{
		metric:   newMetric("load", "tcp_load", "Percentage of TCP children that are awake and processing SIP messages.", []string{}, prometheus.GaugeValue),
		ip:       "",
		protocol: "",
		port:     "",
	}
	return metrics
}
