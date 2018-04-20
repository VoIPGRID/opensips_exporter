package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/VoIPGRID/opensips_exporter/processors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var o *opensips.OpenSIPS
var collectAll = []string{"core:", "shmem:", "net:", "uri:", "tm:", "sl:", "usrloc:", "dialog:", "registrar:", "pkmem:", "load:"}

func handler(w http.ResponseWriter, r *http.Request) {
	collect := r.URL.Query()["collect[]"]
	collectors := make(map[prometheus.Collector]bool)

	if len(collect) == 0 {
		// Collect everything if nothing is specified
		collect = collectAll
	}

	statistics, err := o.GetStatistics(collect...)
	if err != nil {
		http.Error(w, "Could not collect statistics", http.StatusInternalServerError)
		return
	}

	for _, processor := range collect {
		if p, ok := processors.Processors[processor]; ok {
			collectors[p(statistics)] = true
		}
	}

	registry := prometheus.NewRegistry()
	for collector := range collectors {
		err := registry.Register(collector)
		if err != nil {
			log.Fatalf("Couldn't register collector: %v", err)
			http.Error(w, fmt.Sprintf("Couldn't register collector: %s", err), http.StatusInternalServerError)
			return
		}
	}

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		})
	h.ServeHTTP(w, r)
}

func main() {
	listenAddress := ":9737"                 // TODO: make this a flag
	metricsPath := "/metrics"                // TODO: maybe make this a flag
	socketPath := "/var/run/ser-fg/ser.sock" // TODO: make this a flag or even a mandatory argument

	// This part is to mock up setting up and using the Management
	// Interface. Replace/remove this eventually.
	var err error
	o, err = opensips.New(socketPath)
	if err != nil {
		log.Fatalf("failed to open socket: %v\n", err)
	}

	http.HandleFunc(metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>OpenSIPS Exporter</title></head>
			<body>
			<h1>OpenSIPS Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	http.ListenAndServe(listenAddress, nil)
}
