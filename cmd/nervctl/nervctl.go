package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"minik8s.io/pkg/apis/core"

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
	ctx := namespaces.WithNamespace(context.Background(), "default")

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
	} else if strings.Compare("run", os.Args[1]) == 0 {
		// finish port map first and port random assign after that
		// format : nervctl run image name  port-map path:path netnamespace command arg...
		//            arg0  arg1 arg2 arg3     arg4    arg5      arg6         arg7  arg...
		// construct the Container Object
		Container := core.Container{}
		if strings.Contains(os.Args[2], "registry.aliyuncs.com") || strings.Contains(os.Args[2], "docker.io/library") {
			Container.Image = os.Args[2]
		} else {
			Container.Image = "docker.io/library/" + os.Args[2]
		}
		Container.Name = os.Args[3]
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
				if i < 7 {
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
	}
}
