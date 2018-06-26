package processors

import (
	"log"
	"strings"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// loadProcessor describes busy children
// doc: http://www.opensips.org/Documentation/Interface-CoreStatistics-1-11#toc14
type loadProcessor struct {
	statistics map[string]opensips.Statistic
}

type loadMetric struct {
	metric   metric
	ip       string
	port     string
	protocol string
}

func init() {
	OpensipsProcessors["load:"] = loadProcessorFunc
	OpensipsProcessors["tcp-load"] = loadProcessorFunc
	OpensipsProcessors["load"] = loadProcessorFunc
}

// Describe implements prometheus.Collector.
func (p loadProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.loadMetrics() {
		ch <- m.metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p loadProcessor) Collect(ch chan<- prometheus.Metric) {
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
	return &loadProcessor{
		statistics: s,
	}
}

func (p loadProcessor) loadMetrics() map[string]loadMetric {
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
		if len(split) >= 2 {
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
		} else {
			log.Printf("Unable to parse metric '%v'in loadProcessor. Reason: Not enough fields received (protocol, ip, port)", s.Name)
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
