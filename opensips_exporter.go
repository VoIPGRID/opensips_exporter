package main

import (
	"log"
	"net/http"

	"flag"
	"os"
	"strings"

	"fmt"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/VoIPGRID/opensips_exporter/processors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var o *opensips.OpenSIPS
var collectAll = []string{"core:", "shmem:", "net:", "uri:", "tm:", "sl:", "usrloc:", "dialog:", "registrar:", "pkmem:", "load:"}

const envPrefix = "OPENSIPS_EXPORTER"

func handler(w http.ResponseWriter, r *http.Request) {
	collect := r.URL.Query()["collect[]"]
	collectors := make(map[prometheus.Collector]bool)

	if len(collect) == 0 {
		// Collect everything if nothing is specified
		collect = collectAll
	}

	var scrapeProcessor prometheus.Collector
	statistics, err := o.GetStatistics(collect...)
	if err != nil {
		scrapeProcessor = processors.NewScrapeProcessor(0)
		log.Printf("Error encountered while reading statistics from opensips socket: %v", err)
	} else {
		scrapeProcessor = processors.NewScrapeProcessor(1)
	}
	collectors[scrapeProcessor] = true

	var selectedProcessors = map[string]bool{}
	for _, processor := range collect {
		if p, ok := processors.OpensipsProcessors[processor]; ok {
			processorFunc := fmt.Sprintf("%p", p)
			if _, ok := selectedProcessors[processorFunc]; !ok {
				selectedProcessors[processorFunc] = true
				collectors[p(statistics)] = true
			}
		}
	}

	registry := prometheus.NewRegistry()
	for collector := range collectors {
		err := registry.Register(collector)
		if err != nil {
			log.Printf("Couldn't register collector for %v: %v\n", collector, err)
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

// strflag is like flag.String, with value overridden by an environment
// variable (when present). e.g. with socket path, the env var used as default
// is OPENSIPS_EXPORTER_SOCKET_PATH, if present in env.
func strflag(name string, value string, usage string) *string {
	if v, ok := os.LookupEnv(envPrefix + strings.ToUpper(name)); ok {
		return flag.String(name, v, usage)
	}
	return flag.String(name, value, usage)
}

var (
	socketPath  *string
	metricsPath *string
	addr        *string
)

func main() {
	addr = strflag("addr", ":9434", "Address on which the OpenSIPS exporter listens. (e.g. 127.0.0.1:9434)")
	metricsPath = strflag("path", "/metrics", "The path where metrics will be served.")
	socketPath = strflag("socket", "/var/run/ser-fg/ser.sock", "Path to the socket file for OpenSIPS.)")
	flag.Parse()

	// This part is to mock up setting up and using the Management
	// Interface. Replace/remove this eventually.
	var err error
	o, err = opensips.New(*socketPath)
	if err != nil {
		log.Fatalf("failed to open socket: %v\n", err)
	}

	http.HandleFunc(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>OpenSIPS Exporter</title></head>
			<body>
			<h1>OpenSIPS Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Printf("Started OpenSIPS exporter, listening on %v", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
