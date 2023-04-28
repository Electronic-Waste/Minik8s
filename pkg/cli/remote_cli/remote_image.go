package remote_cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/containerd/containerd"
	constant "minik8s.io/pkg/const"
)

// remoteImageService is a gRPC implementation of internalapi.ImageManagerService.
// !!!: there is no need to separate two service of image and runtime container,
// !!!: however, we plan to replace to use cri implement in the future......
type remoteImageService struct {
	timeout     time.Duration
	imageClient *containerd.Client
}

const (
	// How frequently to report identical errors
	IdenticalErrorDelay = 1 * time.Minute
)

func NewRemoteImageService(connectionTimeout time.Duration) (*remoteImageService, error) {
	// build a new cri client
	client, err := containerd.New(constant.Cli_uri)
	// need to call client.Close() to gc this object
	if err != nil {
		return nil, err
	}
	return &remoteImageService{
		timeout:     connectionTimeout,
		imageClient: client,
	}, nil
}

func NewRemoteImageServiceByRunTime(cli *remoteRuntimeService) *remoteImageService {
	return &remoteImageService{
		timeout:     cli.timeout,
		imageClient: cli.runtimeClient,
	}
}

func (cli *remoteImageService) PullImage(ctx context.Context, image string) (string, error) {
	fmt.Printf("want to pull image %s\n", image)
	image_getted, err := cli.imageClient.Pull(ctx, image, containerd.WithPullUnpack)
	if err != nil {
		return "", err
	}
	return image_getted.Name(), nil
}

func (cli *remoteImageService) GetImage(ctx context.Context, image string) (containerd.Image, error) {
	log.Printf("try to find image %s\n", image)
	// get the image if the image exist and create a new image if not
	image_getted, err := cli.imageClient.GetImage(ctx, image)
	if err != nil {
		return nil, err
	}
	return image_getted, nil
}
