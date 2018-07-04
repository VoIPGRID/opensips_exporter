package processors

import (
	"strings"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

// usrlocProcessor provides metrics about User locations via registring.
// doc: http://www.opensips.org/html/docs/modules/1.11.x/usrloc.html#idp5699792
// src: Not clear
type usrlocProcessor struct {
	statistics map[string]opensips.Statistic
}

type usrlocMetric struct {
	metric metric
	domain string
}

func init() {
	OpensipsProcessors["usrloc:"] = usrlocProcessorFunc
	OpensipsProcessors["contacts"] = usrlocProcessorFunc
	OpensipsProcessors["users"] = usrlocProcessorFunc
	OpensipsProcessors["expires"] = usrlocProcessorFunc
}

// Describe implements prometheus.Collector.
func (p usrlocProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.usrlocMetrics() {
		ch <- m.metric.Desc
	}
}

// Collect implements prometheus.Collector.
func (p usrlocProcessor) Collect(ch chan<- prometheus.Metric) {
	for key, u := range p.usrlocMetrics() {
		if u.domain != "" {
			ch <- prometheus.MustNewConstMetric(
				u.metric.Desc,
				u.metric.ValueType,
				p.statistics[key].Value,
				u.domain,
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

func usrlocProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &usrlocProcessor{
		statistics: s,
	}
}

func (p usrlocProcessor) usrlocMetrics() map[string]usrlocMetric {
	var metrics = map[string]usrlocMetric{}

	// Get all usrloc statistics
	var stats []opensips.Statistic
	for _, s := range p.statistics {
		if s.Module == "usrloc" {
			stats = append(stats, s)
		}
	}

	for _, s := range stats {
		split := strings.LastIndex(s.Name, "-")

		if split == -1 {
			continue
		}

		metricType := s.Name[split+1:]
		domain := s.Name[:split]

		switch metricType {
		case "users":
			metric := newMetric("usrloc", metricType, "Number of AOR existing in the USRLOC memory cache for that domain.", []string{"domain"}, prometheus.GaugeValue)
			metrics[s.Name] = usrlocMetric{
				metric: metric,
				domain: domain,
			}
		case "contacts":
			metric := newMetric("usrloc", metricType, "Number of contacts existing in the USRLOC memory cache for that domain.", []string{"domain"}, prometheus.GaugeValue)
			metrics[s.Name] = usrlocMetric{
				metric: metric,
				domain: domain,
			}
		case "expires":
			metric := newMetric("usrloc", metricType, "Total number of expired contacts for that domain.", []string{"domain"}, prometheus.GaugeValue)
			metrics[s.Name] = usrlocMetric{
				metric: metric,
				domain: domain,
			}
		}
	}

	metrics["registered_users"] = usrlocMetric{
		metric: newMetric("usrloc", "registered_users_total", " Total number of AOR existing in the USRLOC memory cache for all domains.", []string{}, prometheus.CounterValue),
		domain: "",
	}
	return metrics
}
