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

// remoteImageService is a gRPC implementation of internalapi.ImageManagerService.
type remoteImageService struct {
	timeout     time.Duration
	imageClient runtime.ImageServiceClient
}

func NewRemoteImageService(connectionTimeout time.Duration) (internalapi.ImageManagerService, error) {
	// build a new cri client
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cli := runtime.NewImageServiceClient(conn)
	return &remoteImageService{
		timeout:     connectionTimeout,
		imageClient: cli,
	}, nil
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

func (cli *remoteImageService) PullImage(ctx context.Context, image *runtime.ImageSpec, auth *runtime.AuthConfig,
	podSandboxConfig *runtime.PodSandboxConfig) (string, error) {
	res, err := cli.imageClient.PullImage(ctx, &runtime.PullImageRequest{
		Image:         image,
		Auth:          auth,
		SandboxConfig: podSandboxConfig,
	})
	if err != nil {
		log.Fatalf("could not pull a image: %v", err)
		return "", err
	}
	log.Printf("Success to pull a Image")
	return res.ImageRef, nil
}
