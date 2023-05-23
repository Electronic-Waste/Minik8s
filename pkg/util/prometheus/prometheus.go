package prometheus

import (
	"fmt"
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/namespaces"
	//"github.com/containerd/containerd"
	"net/http"
	//"minik8s.io/pkg/constant"
	"minik8s.io/pkg/cli/remote_cli"
	"minik8s.io/pkg/apis/core"
)

func NewPrometheusClient() {
	// register a new handler for the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	// start an http server
	http.ListenAndServe(":9090", nil)
}

func RunPrometheus(){


	//client, _ := containerd.New(constant.Cli_uri)
	//defer client.Close()
	//image, _ := client.Pull(context.Background(), "docker.io/library/prom/prometheus:latest", containerd.WithPullUnpack)
	//container, _ := client.NewContainer(
    //    context.Background(),
    //    "prometheus",
    //    containerd.WithNewSnapshot("prometheus-snapshot", image),
    //    containerd.WithNewSpec(oci.WithImageConfig(image)),
    //)

	//image_getted, err := client.GetImage(ctx, image)
}