package remote

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"

	jsonutil "k8s.io/pkg/tools/json"
)

const uri = "unix:///var/run/containerd/containerd.sock"

// remoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient runtime.RuntimeServiceClient
}

func StartContainer(id string) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	_, err = cli.StartContainer(context.TODO(), &runtime.StartContainerRequest{ContainerId: id})
	if err != nil {
		log.Fatalf("could not start a container: %v", err)
	}
	log.Printf("Success to start a new Container")
}

// set filter to nil and list all containers
func ListContainers() {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	res, err := cli.ListContainers(context.TODO(), &runtime.ListContainersRequest{Filter: nil})
	if err != nil {
		log.Fatalf("could not start a container: %v", err)
	}
	log.Printf("Success to List all Container\n")
	for _, con := range res.Containers {
		fmt.Printf("the container id is %s\n", con.GetId())
	}
}

func Convert2State(state uint32) string {
	if state == 1 {
		return "NotReady"
	} else {
		return "Ready"
	}
}

// set filter to nil and list all sandbox
func ListPodSandbox() {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	res, err := cli.ListPodSandbox(context.TODO(), &runtime.ListPodSandboxRequest{Filter: nil})
	if err != nil {
		log.Fatalf("could not finish this rpc: %v", err)
	}
	log.Printf("Success to List all PodSandBox\n")
	for _, con := range res.Items {
		fmt.Printf("the SandBox id is %s and status is %s\n", con.GetId(), Convert2State((uint32(con.State))))
	}
}

// only set the image flag
func PullImage(image string) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewImageServiceClient(conn)

	_, err = cli.PullImage(context.TODO(), &runtime.PullImageRequest{Image: &runtime.ImageSpec{
		Image: image,
	}})
	if err != nil {
		log.Fatalf("could not finish this rpc: %v", err)
	}
	log.Printf("Success to Pull Image %s\n", image)
}

func CreateContainer(pod_id string, con_config string, pod_config string) {
	con_fig := new(runtime.ContainerConfig)
	pod_fig := new(runtime.PodSandboxConfig)
	jsonutil.ParseContainerConfig(con_config, con_fig, os.Stdout)
	jsonutil.ParseSandBoxConfig(pod_config, pod_fig, os.Stdout)

	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	_, err = cli.CreateContainer(context.TODO(), &runtime.CreateContainerRequest{
		PodSandboxId:  pod_id,
		Config:        con_fig,
		SandboxConfig: pod_fig,
	})
	if err != nil {
		log.Fatalf("could not finish this rpc: %v", err)
	}
	log.Printf("Success to Create a Container\n")
}

func StopContainer(con_id string) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	_, err = cli.StopContainer(context.TODO(), &runtime.StopContainerRequest{
		ContainerId: con_id,
	})
	if err != nil {
		log.Fatalf("could not stop a container: %v", err)
	}
	log.Printf("Success to stop a Container which id is %s\n", con_id)
}

func RemoveContainer(con_id string) {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	_, err = cli.RemoveContainer(context.TODO(), &runtime.RemoveContainerRequest{
		ContainerId: con_id,
	})
	if err != nil {
		log.Fatalf("could not remove a container: %v", err)
	}
	log.Printf("Success to remove a Container which id is %s\n", con_id)
}

func RunPodSandbox(fig_path string) {
	sand_fig := new(runtime.PodSandboxConfig)
	jsonutil.ParseSandBoxConfig(fig_path, sand_fig, os.Stdout)

	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connect to containerd")
	defer conn.Close()
	cli := runtime.NewRuntimeServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	res, err := cli.RunPodSandbox(context.TODO(), &runtime.RunPodSandboxRequest{
		Config: sand_fig,
	})
	if err != nil {
		log.Fatalf("could not start a pod: %v", err)
	}
	log.Printf("Success to start a new Pod %s", res.PodSandboxId)
}
