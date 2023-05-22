package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
)

func TestPrometheus(t *testing.T) {
	// register a new handler for the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	// start an http server
	http.ListenAndServe(":9090", nil)
}
