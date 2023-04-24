package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	jsonutil "minik8s.io/pkg/util/json"

	runtimeapi "minik8s.io/cri-api/pkg/apis/runtime/v1"
	"minik8s.io/pkg/cri/remote"
)

// just a simple command line tools used to test the runc function
// try to parse the cmd and get the container id
func main() {
	if len(os.Args) < 2 {
		err := errors.New("less os args")
		panic(err)
	}

	cli, err := remote.NewRemoteRuntimeService(remote.IdenticalErrorDelay)
	if err != nil {
		panic(err)
	}

	if strings.Compare("list", os.Args[1]) == 0 {
		fmt.Println("Hello, List all Containers")
		cli.ListContainers(context.TODO(), nil)
		fmt.Println("list all sandbox")
		cli.ListPodSandbox(context.TODO(), nil)
	} else if strings.Compare("pull", os.Args[1]) == 0 {
		if len(os.Args) != 3 {
			err := errors.New("err os args")
			panic(err)
		}
		remote.PullImage(os.Args[2])
	} else if strings.Compare("create", os.Args[1]) == 0 {
		if len(os.Args) != 5 {
			err := errors.New("err os args")
			panic(err)
		}
		var con_fig = new(runtimeapi.ContainerConfig)
		var sand_fig = new(runtimeapi.PodSandboxConfig)
		jsonutil.ParseContainerConfig(os.Args[3], con_fig, os.Stdout)
		jsonutil.ParseSandBoxConfig(os.Args[4], sand_fig, os.Stdout)
		cli.CreateContainer(context.TODO(), os.Args[2], con_fig, sand_fig)
	} else if strings.Compare("start", os.Args[1]) == 0 {
		fmt.Printf("Hello World and get Container id is %s\n", os.Args[2])
		cli.StartContainer(context.TODO(), os.Args[2])
	} else if strings.Compare("stop", os.Args[1]) == 0 {
		fmt.Printf("try to stop the container %s\n", os.Args[2])
		cli.StopContainer(context.TODO(), os.Args[2], 100)
	} else if strings.Compare("remove", os.Args[1]) == 0 {
		fmt.Printf("try to remove the container %s\n", os.Args[2])
		cli.RemoveContainer(context.TODO(), os.Args[2])
	} else if strings.Compare("runp", os.Args[1]) == 0 {
		fmt.Printf("try to run a new Pod in the machine\n")
		var sand_fig = new(runtimeapi.PodSandboxConfig)
		jsonutil.ParseSandBoxConfig(os.Args[2], sand_fig, os.Stdout)
		cli.RunPodSandbox(context.TODO(), sand_fig, "")
	}
}
