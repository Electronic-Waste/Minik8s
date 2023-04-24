package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"minik8s.io/pkg/remote"
)

// just a simple command line tools used to test the runc function
// try to parse the cmd and get the container id
func main() {
	if len(os.Args) < 2 {
		err := errors.New("less os args")
		panic(err)
	}

	if strings.Compare("list", os.Args[1]) == 0 {
		fmt.Println("Hello, List all Containers")
		remote.ListContainers()
		fmt.Println("list all sandbox")
		remote.ListPodSandbox()
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
		remote.CreateContainer(os.Args[2], os.Args[3], os.Args[4])
	} else if strings.Compare("start", os.Args[1]) == 0 {
		fmt.Printf("Hello World and get Container id is %s\n", os.Args[2])
		remote.StartContainer(os.Args[2])
	} else if strings.Compare("stop", os.Args[1]) == 0 {
		fmt.Printf("try to stop the container %s\n", os.Args[2])
		remote.StopContainer(os.Args[2])
	} else if strings.Compare("remove", os.Args[1]) == 0 {
		fmt.Printf("try to remove the container %s\n", os.Args[2])
		remote.RemoveContainer(os.Args[2])
	} else if strings.Compare("runp", os.Args[1]) == 0 {
		fmt.Printf("try to run a new Pod in the machine\n")

	}
}
