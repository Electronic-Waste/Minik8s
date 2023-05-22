package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewPrometheusClient() {
	// register a new handler for the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	// start an http server
	http.ListenAndServe(":9001", nil)
}
