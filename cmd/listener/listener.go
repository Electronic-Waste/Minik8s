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
	t := time.NewTicker(2 * time.Second)
	stopper := make(chan bool)
	go func(ch chan bool) {
		fmt.Println("enter anything to end listener")
		var tmp int
		fmt.Scanf("%d", tmp)
		ch <- true
		fmt.Printf("input %d and listener is ending...", tmp)
	}(stopper)
	is_out := false
	for {
		select {
		case <-stopper:
			c.UnRegisterPod("test")
			is_out = true
		case <-t.C:
			status, err := c.GetPodMetric("test")
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			fmt.Println(status)
		}
		if is_out {
			break
		}
	}
	fmt.Println("out loop")
	time.Sleep(5 * time.Second)
	fmt.Println("listener finished")
}
