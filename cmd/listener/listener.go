package main

import (
	"fmt"
	"minik8s.io/pkg/kubelet/cadvisor"
	"os"
	"time"
)

func main() {
	c := cadvisor.GetCAdvisor()
	// detect for go1
	err := c.RegisterPod("test")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for {
		status, err := c.GetPodMetric("test")
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(status)
		// every 3 second detect a cpu usage
		time.Sleep(2 * time.Second)
	}
}
