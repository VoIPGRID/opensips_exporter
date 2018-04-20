package processors

import (
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus"
)

type ShmemProcessor struct {
	statistics map[string]opensips.Statistic
}

var shmemLabelNames = []string{}
var shmemMetrics = map[string]metric{
	"total_size":     newMetric("shmem", "total_size", "Total size of shared memory available to OpenSIPS processes.", shmemLabelNames, prometheus.GaugeValue),
	"used_size":      newMetric("shmem", "used_size", "Amount of shared memory requested and used by OpenSIPS processes.", shmemLabelNames, prometheus.GaugeValue),
	"real_used_size": newMetric("shmem", "real_used_size", "Amount of shared memory requested by OpenSIPS processes + malloc overhead", shmemLabelNames, prometheus.GaugeValue),
	"max_used_size":  newMetric("shmem", "max_used_size", "Maximum amount of shared memory ever used by OpenSIPS processes.", shmemLabelNames, prometheus.GaugeValue),
	"free_size":      newMetric("shmem", "free_size", "Free memory available. Computed as total_size - real_used_size", shmemLabelNames, prometheus.GaugeValue),
	"fragments":      newMetric("shmem", "fragments", "Total number of fragments in the shared memory.", shmemLabelNames, prometheus.GaugeValue),
}

func init() {
	for metric := range shmemMetrics {
		Processors[metric] = shmemProcessorFunc
	}
	Processors["shmem:"] = shmemProcessorFunc
}

func (c ShmemProcessor) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range shmemMetrics {
		ch <- metric.Desc
	}
}

func (p ShmemProcessor) Collect(ch chan<- prometheus.Metric) {
	for key, metric := range shmemMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.ValueType,
			p.statistics[key].Value,
		)
	}
}

func shmemProcessorFunc(s map[string]opensips.Statistic) prometheus.Collector {
	return &ShmemProcessor{
		statistics: s,
	}
}
