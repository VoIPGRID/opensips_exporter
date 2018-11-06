package processors

import (
	"strings"

	"log"

	"fmt"

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
	process  string
}

func init() {
	OpensipsProcessors["load:"] = loadProcessorFunc
	OpensipsProcessors["tcp-load"] = loadProcessorFunc
	OpensipsProcessors["load"] = loadProcessorFunc
	OpensipsProcessors["load1m"] = loadProcessorFunc
	OpensipsProcessors["load10m"] = loadProcessorFunc
	OpensipsProcessors["load-all"] = loadProcessorFunc
	OpensipsProcessors["load1m-all"] = loadProcessorFunc
	OpensipsProcessors["load10m-all"] = loadProcessorFunc
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
		} else if u.process != "" {
			ch <- prometheus.MustNewConstMetric(
				u.metric.Desc,
				u.metric.ValueType,
				p.statistics[key].Value,
				u.process,
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

		if s.Name == "tcp-load" {
			// tcp-load is a OpenSIPS < 2 metric.
			metrics["tcp-load"] = loadMetric{
				metric:   newMetric("load", "tcp_load", "Percentage of TCP children that are awake and processing SIP messages.", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if strings.Contains(s.Name, "udp:") || strings.Contains(s.Name, "tcp:") {
			// The load metric is an OpenSIPS < 2 metric. The format for it is:
			// load:udp:127.0.0.1:5060-load = 0
			metrics[s.Name] = parseLegacyLoadFormat(s)
			continue
		} else if strings.Contains(s.Name, "proc") {
			// This metric is an OpenSIPS > 2 metric. The format for it is:
			// load:load-proc-1:: 0
			// load:load1m-proc-1:: 0
			// load:load10m-proc-1:: 0
			metrics[s.Name] = parseNewLoadFormat(s)
			continue
		} else if s.Name == "load" {
			metrics["load"] = loadMetric{
				metric:   newMetric("load", "core", "The realtime load of entire OpenSIPS - this counts all the core processes of OpenSIPS; the additional processes requested by modules are not counted in this load.", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if s.Name == "load1m" {
			metrics["load1m"] = loadMetric{
				metric:   newMetric("load", "core_1m", "The last minute average load of core OpenSIPS (covering only core/SIP processes)", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if s.Name == "load10m" {
			metrics["load10m"] = loadMetric{
				metric:   newMetric("load", "core_10m", "The last 10 minute average load of core OpenSIPS (covering only core/SIP processes)", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if s.Name == "load-all" {
			metrics["load-all"] = loadMetric{
				metric:   newMetric("load", "all", "The realtime load of entire OpenSIPS, counting both core and module processes.", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if s.Name == "load1m-all" {
			metrics["load1m-all"] = loadMetric{
				metric:   newMetric("load", "all_1m", "The last minute average load of entire OpenSIPS (covering all processes).", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		} else if s.Name == "load10m-all" {
			metrics["load10m-all"] = loadMetric{
				metric:   newMetric("load", "all_10m", "The last 10 minute average load of entire OpenSIPS (covering all processes).", []string{}, prometheus.GaugeValue),
				ip:       "",
				protocol: "",
				port:     "",
				process:  "",
			}
			continue
		}
	}
	return metrics
}

func parseLegacyLoadFormat(statistic opensips.Statistic) loadMetric {
	var ret loadMetric
	split := strings.Split(statistic.Name, ":")
	if len(split) >= 2 {
		protocol := split[0]
		ip := split[1]
		port := split[2]
		port = strings.Trim(port, "-load")

		metric := newMetric("load", "load", "The realtime load of entire OpenSIPS - this counts all the core processes of OpenSIPS; the additional processes requested by modules are not counted in this load.", []string{"ip", "port", "protocol"}, prometheus.GaugeValue)
		ret = loadMetric{
			metric:   metric,
			ip:       ip,
			port:     port,
			protocol: protocol,
			process:  "",
		}
	} else {
		log.Printf("Unable to parse metric '%v'in loadProcessor. Reason: Not enough fields received (protocol, ip, port)", statistic.Name)
	}
	return ret
}

func parseNewLoadFormat(statistic opensips.Statistic) loadMetric {
	split := strings.Split(statistic.Name, "-")
	process := split[len(split)-1]
	load := split[0]

	switch load {
	case "load":
		return loadMetric{
			metric:   newMetric("load", "process", "The realtime load of the process ID.", []string{"process"}, prometheus.GaugeValue),
			ip:       "",
			port:     "",
			protocol: "",
			process:  process,
		}
	case "load1m":
		return loadMetric{
			metric:   newMetric("load", "1m", "The last minute average load of the process ID.", []string{"process"}, prometheus.GaugeValue),
			ip:       "",
			port:     "",
			protocol: "",
			process:  process,
		}
	case "load10m":
		return loadMetric{
			metric:   newMetric("load", "10m", "The last 10 minutes average load of the process ID.", []string{"process"}, prometheus.GaugeValue),
			ip:       "",
			port:     "",
			protocol: "",
			process:  process,
		}
	}
	fmt.Errorf("could not parse load metric for %v", statistic.Name)
	return loadMetric{}
}
