package remote_cli

import (
	"context"
	"time"

	"github.com/containerd/containerd"
	constant "minik8s.io/pkg/const"
)

// remoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient *containerd.Client
}

func NewRemoteRuntimeService(connectionTimeout time.Duration) (*remoteRuntimeService, error) {
	// build a new cri client
	client, err := containerd.New(constant.Cri_uri)
	// need to call client.Close() to gc this object
	if err != nil {
		return nil, err
	}
	return &remoteRuntimeService{
		connectionTimeout,
		client,
	}, nil
}

func NewRemoteImageServiceByImageService(cli *remoteImageService) *remoteRuntimeService {
	return &remoteRuntimeService{
		timeout:       cli.timeout,
		runtimeClient: cli.imageClient,
	}
}

// set filter to nil and list all containers
func (cli *remoteRuntimeService) ListContainers(ctx context.Context, filters ...string) ([]containerd.Container, error) {
	res, err := cli.runtimeClient.Containers(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cli *remoteRuntimeService) StartContainer(ctx context.Context, containerID string) error {
	return nil
}
