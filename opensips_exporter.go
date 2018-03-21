package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VoIPGRID/opensips_exporter/opensips"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenAddress := ":9737"                 // TODO: make this a flag
	metricsPath := "/metrics"                // TODO: maybe make this a flag
	socketPath := "/var/run/ser-fg/ser.sock" // TODO: make this a flag or even a mandatory argument

	// This part is to mock up setting up and using the Management
	// Interface. Replace/remove this eventually.
	o, err := opensips.New(socketPath)
	if err != nil {
		log.Fatalf("failed to open socket: %v\n", err)
	}
	statistics, err := o.GetStatistics("all")
	fmt.Printf("%q\n", statistics)
	if err != nil {
		log.Fatalf("failed to get statistics: %v\n", err)
	}
	err = o.Close()
	if err != nil {
		log.Fatalf("failed to close: %v\n", err)
	}

	// TODO: set up a collector and register it (e.g.
	// prometheus.MustRegister())
	http.Handle(metricsPath, promhttp.Handler())
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
