package processors

import (
	"strings"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type PkmemProcessor struct {
	statistics map[string]opensips.Statistic
}

type pkmemMetric struct {
	metric metric
	pid    string
}

func init() {
	Processors["pkmem:"] = pkmemProcessorFunc
	Processors["total_size"] = pkmemProcessorFunc
	Processors["used_size"] = pkmemProcessorFunc
	Processors["real_used_size"] = pkmemProcessorFunc
	Processors["max_used_size"] = pkmemProcessorFunc
	Processors["free_size"] = pkmemProcessorFunc
	Processors["fragments"] = pkmemProcessorFunc
}

func (p PkmemProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.pkmemMetrics() {
		ch <- m.metric.Desc
	}
}

func (p PkmemProcessor) Collect(ch chan<- prometheus.Metric) {
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
	return &PkmemProcessor{
		statistics: s,
	}
}

func (p PkmemProcessor) pkmemMetrics() map[string]pkmemMetric {
	var metrics = map[string]pkmemMetric{}

	// Get all usrloc statistics
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
