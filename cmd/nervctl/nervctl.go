package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containerd/containerd/namespaces"
	"minik8s.io/pkg/cli/remote_cli"
)

// !!!:this version of ctl tools is just for test used

func main() {
	if len(os.Args) < 2 {
		panic("error args num")
	}

	image_manager, err := remote_cli.NewRemoteImageService(remote_cli.IdenticalErrorDelay)
	runtime_manager := remote_cli.NewRemoteImageServiceByImageService(image_manager)
	// use a trick namespace to test for namespace isolation
	// this ctx will in the specific namespace
	// which can be seen by ctr ns ls
	// ps. we will get our image in the minik8s.io namespace in release
	// and put the test case in test namespace
	ctx := namespaces.WithNamespace(context.Background(), "test")

	if err != nil {
		panic(err)
	}
	if strings.Compare("pull", os.Args[1]) == 0 {
		// for pull cmd test
		if len(os.Args) < 3 {
			panic("less num of args")
		}
		res, err := image_manager.PullImage(ctx, os.Args[2])
		if err != nil {
			fmt.Println("pull image err")
			panic(err)
		}
		fmt.Printf("pull image %s success\n", res)
	} else if strings.Compare("list", os.Args[1]) == 0 {
		res, err := runtime_manager.ListContainers(ctx)
		if err != nil {
			panic(err)
		}
		for _, container := range res {
			fmt.Printf("the container id is %s\n", container.ID())
		}
	} else if strings.Compare("get", os.Args[1]) == 0 {
		// the cmd used to get the image
		res, err := image_manager.GetImage(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		fmt.Printf("the image name is %s\n", res.Name())
	}
}
