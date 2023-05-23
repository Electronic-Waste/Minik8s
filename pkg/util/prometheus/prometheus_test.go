package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"context"
	"testing"
	"fmt"
)

func TestPrometheus(t *testing.T) {
	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", "localhost:9090/api/v1/query?query=node_cpu_seconds_total")
	if err != nil {
		log.Fatal(err)
	}
}
