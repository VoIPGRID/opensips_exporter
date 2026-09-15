package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"fmt"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/VoIPGRID/opensips_exporter/opensips/jsonrpc"
	"github.com/VoIPGRID/opensips_exporter/processors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var j *jsonrpc.JSONRPC
var collectAll = []string{"core:", "shmem:", "net:", "uri:", "tm:", "sl:", "usrloc:", "dialog:", "registrar:", "pkmem:", "load:", "tmx:"}

const envPrefix = "OPENSIPS_EXPORTER"

func handler(w http.ResponseWriter, r *http.Request) {
	collect := r.URL.Query()["collect[]"]
	collectors := make(map[prometheus.Collector]bool)

	if len(collect) == 0 {
		// Collect everything if nothing is specified
		collect = collectAll
	}
	var scrapeProcessor prometheus.Collector

	var statistics map[string]opensips.Statistic
	var err error
	j = jsonrpc.New(*httpEndpoint)
	statistics, err = j.GetStatistics(collect...)

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
			log.Printf("Problems registering the %T processor (could be due to no metrics found for this processor). Error: %v\n", collector, err)
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
	metricsPath  *string
	addr         *string
	httpEndpoint *string
)

func main() {
	addr = strflag("addr", ":9434", "Address on which the OpenSIPS exporter listens. (e.g. 127.0.0.1:9434)")
	metricsPath = strflag("path", "/metrics", "The path where metrics will be served.")
	httpEndpoint = strflag("http_address", "http://127.0.0.1:8888/mi/", "Address to query the Management Interface through HTTP with (e.g. http://127.0.0.1:8888/mi/)")
	flag.Parse()

	if *httpEndpoint == "" {
		log.Fatalf("The -http_address flag is not set. Exiting.")
	}

	http.HandleFunc(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>OpenSIPS Exporter</title></head>
			<body>
			<h1>OpenSIPS Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Printf("Error writing response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	log.Printf("Started OpenSIPS exporter, listening on %v", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
