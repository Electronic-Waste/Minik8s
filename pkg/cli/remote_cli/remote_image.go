package remote_cli

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
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

func (cli *remoteImageService) PullImage(image string) (string, error) {
	// use a trick namespace to test for namespace isolation
	// this ctx will in the specific namespace
	// which can be seen by ctr ns ls
	// ps. we will get our image in the minik8s.io namespace in release
	// and put the test case in test namespace
	ctx := namespaces.WithNamespace(context.Background(), "test")
	fmt.Printf("want to pull image %s\n", image)
	image_getted, err := cli.imageClient.Pull(ctx, image, containerd.WithPullUnpack)
	if err != nil {
		return "", err
	}
	return image_getted.Name(), nil
}
