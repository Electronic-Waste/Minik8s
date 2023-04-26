package remote

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	runtime "minik8s.io/cri-api/pkg/apis/runtime/v1"
)

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
