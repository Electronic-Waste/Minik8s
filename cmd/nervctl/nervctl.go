package main

import (
	"fmt"
	"os"
	"strings"

	"minik8s.io/pkg/cli/remote_cli"
)

// !!!:this version of ctl tools is just for test used

func main() {
	if len(os.Args) < 2 {
		panic("error args num")
	}

	image_manager, err := remote_cli.NewRemoteImageService(remote_cli.IdenticalErrorDelay)

	if err != nil {
		panic(err)
	}
	if strings.Compare("pull", os.Args[1]) == 0 {
		// for pull cmd test
		if len(os.Args) < 3 {
			panic("less num of args")
		}
		res, err := image_manager.PullImage(os.Args[2])
		if err != nil {
			fmt.Println("pull image err")
			panic(err)
		}
		fmt.Printf("pull image %s success\n", res)
	}
}
