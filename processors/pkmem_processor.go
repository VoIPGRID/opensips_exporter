package processors

import (
	"strings"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// pkmemProcessor provides metrics about private memory usage and fragments
// doc: http://www.opensips.org/Documentation/Interface-CoreStatistics-1-11#toc28
// src: https://github.com/OpenSIPS/opensips/blob/b917c70ba8d5797dc6364aaf702c3415539be52a/core_stats.c#L165
type pkmemProcessor struct {
	statistics map[string]opensips.Statistic
}

type pkmemMetric struct {
	metric metric
	pid    string
}

func init() {
	OpensipsProcessors["pkmem:"] = pkmemProcessorFunc
	OpensipsProcessors["total_size"] = pkmemProcessorFunc
	OpensipsProcessors["used_size"] = pkmemProcessorFunc
	OpensipsProcessors["real_used_size"] = pkmemProcessorFunc
	OpensipsProcessors["max_used_size"] = pkmemProcessorFunc
	OpensipsProcessors["free_size"] = pkmemProcessorFunc
	OpensipsProcessors["fragments"] = pkmemProcessorFunc
}

// Describe implements prometheus.Collector.
func (p pkmemProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.pkmemMetrics() {
		ch <- m.metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p pkmemProcessor) Collect(ch chan<- prometheus.Metric) {
	for key, u := range p.pkmemMetrics() {
		ch <- prometheus.MustNewConstMetric(
			u.metric.Desc,
			u.metric.ValueType,
			p.statistics[key].Value,
			u.pid,
		)
	}
}

func pkmemProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &pkmemProcessor{
		statistics: s,
	}
}

func (p pkmemProcessor) pkmemMetrics() map[string]pkmemMetric {
	var metrics = map[string]pkmemMetric{}

	// Get all pkmem statistics
	var stats []opensips.Statistic
	for _, s := range p.statistics {
		if s.Module == "pkmem" {
			stats = append(stats, s)
		}
	}

	for _, s := range stats {
		split := strings.Index(s.Name, "-")

		if split == -1 {
			continue
		}

		pid := s.Name[:split]
		metricType := s.Name[split+1:]

		switch metricType {
		case "total_size":
			metric := newMetric("pkmem", metricType, "Total size of private memory available to the OpenSIPS process.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		case "used_size":
			metric := newMetric("pkmem", metricType, "Amount of private memory requested and used by the OpenSIPS process.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		case "real_used_size":
			metric := newMetric("pkmem", metricType, "Amount of private memory requested by the OpenSIPS process, including allocator-specific metadata.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		case "max_used_size":
			metric := newMetric("pkmem", metricType, "The maximum amount of private memory ever used by the OpenSIPS process.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		case "free_size":
			metric := newMetric("pkmem", metricType, "Free private memory available for the OpenSIPS process. Computed as total_size - real_used_size.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		case "fragments":
			metric := newMetric("pkmem", metricType, "Currently available number of free fragments in the private memory for OpenSIPS process.", []string{"pid"}, prometheus.GaugeValue)
			metrics[s.Name] = pkmemMetric{
				metric: metric,
				pid:    pid,
			}
		}
	}

	return metrics
}
