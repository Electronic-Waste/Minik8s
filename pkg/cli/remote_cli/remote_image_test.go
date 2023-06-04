package remote_cli

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"testing"
	"time"
)

func TestGetImage(t *testing.T) {
	cli, err := NewRemoteImageService(1 * time.Second)
	if err != nil {
		t.Error(err)
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	image, err := cli.GetImage(ctx, "docker.io/library/jobserver:latest")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("get image is %s\n", image.Name())
}
