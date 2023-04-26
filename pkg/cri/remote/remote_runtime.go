package remote

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	internalapi "minik8s.io/cri-api/pkg/apis"
	runtime "minik8s.io/cri-api/pkg/apis/runtime/v1"
)

const uri = "unix:///var/run/containerd/containerd.sock"

// remoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient runtime.RuntimeServiceClient
}

const (
	// How frequently to report identical errors
	IdenticalErrorDelay = 1 * time.Minute
)

func NewRemoteRuntimeService(connectionTimeout time.Duration) (internalapi.RuntimeService, error) {
	// build a new cri client
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cli := runtime.NewRuntimeServiceClient(conn)
	return &remoteRuntimeService{
		timeout:       connectionTimeout,
		runtimeClient: cli,
	}, nil
}

func (cli *remoteRuntimeService) StartContainer(ctx context.Context, containerID string) error {
	_, err := cli.runtimeClient.StartContainer(ctx, &runtime.StartContainerRequest{ContainerId: containerID})
	if err != nil {
		log.Fatalf("could not start a container: %v", err)
		return err
	}
	log.Printf("Success to start a new Container")
	return nil
}

// set filter to nil and list all containers
func (cli *remoteRuntimeService) ListContainers(ctx context.Context, filter *runtime.ContainerFilter) ([]*runtime.Container, error) {
	res, err := cli.runtimeClient.ListContainers(ctx, &runtime.ListContainersRequest{Filter: filter})
	if err != nil {
		log.Fatalf("could not start a container: %v", err)
		return nil, err
	}

	// for debug use
	for _, con := range res.Containers {
		fmt.Printf("the container id is %s\n", con.GetId())
	}

	return res.Containers, nil
}

func (cli *remoteRuntimeService) CreateContainer(ctx context.Context, podSandboxID string, config *runtime.ContainerConfig,
	sandboxConfig *runtime.PodSandboxConfig) (string, error) {
	res, err := cli.runtimeClient.CreateContainer(ctx, &runtime.CreateContainerRequest{
		PodSandboxId:  podSandboxID,
		Config:        config,
		SandboxConfig: sandboxConfig,
	})
	if err != nil {
		log.Fatalf("could not finish this rpc: %v", err)
		return "", err
	}
	log.Printf("Success to Create a Container\n")
	return res.ContainerId, nil
}

func (cli *remoteRuntimeService) StopContainer(ctx context.Context, containerID string, timeout int64) error {
	_, err := cli.runtimeClient.StopContainer(ctx, &runtime.StopContainerRequest{
		ContainerId: containerID,
		Timeout:     timeout,
	})
	if err != nil {
		log.Fatalf("could not stop a container: %v", err)
		return err
	}
	log.Printf("Success to stop a Container which id is %s\n", containerID)
	return nil
}

func (cli *remoteRuntimeService) RemoveContainer(ctx context.Context, containerID string) error {
	_, err := cli.runtimeClient.RemoveContainer(ctx, &runtime.RemoveContainerRequest{
		ContainerId: containerID,
	})
	if err != nil {
		log.Fatalf("could not remove a container: %v", err)
		return err
	}
	log.Printf("Success to remove a Container which id is %s\n", containerID)
	return nil
}

func Convert2State(state uint32) string {
	if state == 1 {
		return "NotReady"
	} else {
		return "Ready"
	}
}

// set filter to nil and list all sandbox
func (cli *remoteRuntimeService) ListPodSandbox(ctx context.Context, filter *runtime.PodSandboxFilter) ([]*runtime.PodSandbox, error) {
	res, err := cli.runtimeClient.ListPodSandbox(ctx, &runtime.ListPodSandboxRequest{Filter: filter})
	if err != nil {
		log.Fatalf("could not finish this rpc: %v", err)
		return nil, err
	}

	// for debug use
	log.Printf("Success to List all PodSandBox\n")
	for _, con := range res.Items {
		fmt.Printf("the SandBox id is %s and status is %s\n", con.GetId(), Convert2State((uint32(con.State))))
	}

	return res.Items, nil
}

func (cli *remoteRuntimeService) RunPodSandbox(ctx context.Context,
	config *runtime.PodSandboxConfig, runtimeHandler string) (string, error) {
	res, err := cli.runtimeClient.RunPodSandbox(ctx, &runtime.RunPodSandboxRequest{
		Config:         config,
		RuntimeHandler: runtimeHandler,
	})
	if err != nil {
		log.Fatalf("could not start a pod: %v", err)
		return "", err
	}
	log.Printf("Success to start a new Pod %s", res.PodSandboxId)
	return res.PodSandboxId, nil
}

// UpdateRuntimeConfig updates runtime configuration if specified
func (cli *remoteRuntimeService) UpdateRuntimeConfig(ctx context.Context, runtimeConfig *runtime.RuntimeConfig) error {
	_, err := cli.runtimeClient.UpdateRuntimeConfig(ctx, &runtime.UpdateRuntimeConfigRequest{
		RuntimeConfig: runtimeConfig,
	})

	if err != nil {
		fmt.Println("err in the updateRuntimeConfig")
		return err
	}
	return nil
}

// Status returns the status of the runtime.
func (cli *remoteRuntimeService) Status(ctx context.Context, verbose bool) (*runtime.StatusResponse, error) {
	res, err := cli.runtimeClient.Status(ctx, &runtime.StatusRequest{
		Verbose: verbose,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
