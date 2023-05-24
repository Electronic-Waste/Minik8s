package main

import (
	"fmt"
	"minik8s.io/pkg/kubelet/cadvisor"
	"os"
	"time"
)

func main() {
	c := cadvisor.GetNewListener()
	// detect for go1
	err := c.RegisterContainer("go2")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for {
		status := c.GetStats("go2")
		fmt.Println(status)
		// every 3 second detect a cpu usage
		time.Sleep(3 * time.Second)
	}
}
