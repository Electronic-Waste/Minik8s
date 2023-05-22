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
	image_manager,_ := remote_cli.NewRemoteImageService(remote_cli.IdenticalErrorDelay)
	runtime_manager := remote_cli.NewRemoteImageServiceByImageService(image_manager)
	ctx := namespaces.WithNamespace(context.Background(), "default")

	res, err := image_manager.PullImage(ctx, "prom/prometheus:latest")
	if err != nil {
		fmt.Println("pull image err")
		panic(err)
	}
	fmt.Printf("pull image %s success\n", res)

	// finish port map first and port random assign after that
	// format : nervctl run image name  port-map path:path netnamespace command arg...
	//            arg0  arg1 arg2 arg3     arg4    arg5      arg6         arg7  arg...
	// construct the Container Object
	Container := core.Container{}
	Container.Image = "docker.io/library/prom/prometheus:latest"
	Container.Name = "prometheus"
	Container.Ports, err = core.ConstructPorts(os.Args[4])
	if err != nil {
		panic(err)
	}
	Container.Mounts, err = core.ConstructMounts(os.Args[5])
	if err != nil {
		panic(err)
	}
	NetNameSpace := os.Args[6]
	if len(os.Args) < 8 {

	} else {
		Container.Command = append(Container.Command, os.Args[7])
		for i, arg := range os.Args {
			if i < 8 {
				continue
			}
			Container.Args = append(Container.Args, arg)
		}
	}
	fmt.Printf("get the cmd is \n %s \n", Container.String())
	err := runtime_manager.StartContainer(ctx, Container, NetNameSpace)
	if err != nil {
		fmt.Println("start a container failed")
		panic(err)
	}
	

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