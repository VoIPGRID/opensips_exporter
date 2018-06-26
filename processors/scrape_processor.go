package processors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type scrapeProcessor struct {
	upMetric metric
	upStatus float64
}

// Describe implements prometheus.Collector.
func (p scrapeProcessor) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.upMetric.Desc
}

// Collect implements prometheus.Collector.
func (p scrapeProcessor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		p.upMetric.Desc,
		p.upMetric.ValueType,
		p.upStatus,
	)
}

func NewScrapeProcessor(upStatus float64) prometheus.Collector {
	return &scrapeProcessor{
		upMetric: newMetric("", "up", "Whether the opensips exporter could read metrics from the Management Interface socket. (i.e. is OpenSIPS up)", []string{}, prometheus.GaugeValue),
		upStatus: upStatus,
	}
}
