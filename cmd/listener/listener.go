package main

import (
	"fmt"
	"minik8s.io/pkg/kubelet/cadvisor"
	"time"
)

func main() {
	for {
		status, err := cadvisor.GetContainerMetric("go1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(status)
		// every 3 second detect a cpu usage
		time.Sleep(3 * time.Second)
	}
}
