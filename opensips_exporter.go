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

var o *opensips.OpenSIPS
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
	switch *protocol {
	case "mi_datagram":
		o, err = opensips.New(*socketPath)
		if err != nil {
			log.Fatalf("Could not create datagram socket: %v", err)
		}
		statistics, err = o.GetStatistics(collect...)
	case "mi_http":
		j = jsonrpc.New(*httpEndpoint)
		statistics, err = j.GetStatistics(collect...)
	}

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
	socketPath   *string
	metricsPath  *string
	addr         *string
	protocol     *string
	httpEndpoint *string
)

func main() {
	addr = strflag("addr", ":9434", "Address on which the OpenSIPS exporter listens. (e.g. 127.0.0.1:9434)")
	metricsPath = strflag("path", "/metrics", "The path where metrics will be served.")
	socketPath = strflag("socket", "/var/run/ser-fg/ser.sock", "Path to the socket file for OpenSIPS.)")
	httpEndpoint = strflag("http_address", "http://127.0.0.1:8888/mi/", "Address to query the Management Interface through HTTP with (e.g. http://127.0.0.1:8888/mi/)")
	protocol = strflag("protocol", "", "Which protocol to use to get data from the Management Interface (mi_datagram & mi_http currently supported)")
	flag.Parse()

	switch *protocol {
	case "mi_datagram":
		if *socketPath == "" {
			log.Fatalf("The -protocol flag is set to mi_datagram but the -socket param is not set. Exiting.")
		}
	case "mi_http":
		if *httpEndpoint == "" {
			log.Fatalf("The -protocol is set to mi_http but the -http_address flag is not set. Exiting.")
		}
	default:
		log.Fatalf("Please set the -protocol flag to define which protocol the exporter should use to query for metrics. (mi_datagram or mi_http)")
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
