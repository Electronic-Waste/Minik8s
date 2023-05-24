package prometheus

import (
	//"fmt"
	//"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//"github.com/containerd/containerd/oci"
	//"github.com/containerd/containerd/namespaces"
	//"github.com/containerd/containerd"
	"net/http"
	//"minik8s.io/pkg/constant"
	//"minik8s.io/pkg/cli/remote_cli"
	//"minik8s.io/pkg/apis/core"
)

//暴露端口给prometheus
func NewPrometheusClient() {
	// register a new handler for the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	// start an http server
	http.ListenAndServe(":9100", nil)
}

func RunPrometheus(){
	go NewPrometheusClient()

}