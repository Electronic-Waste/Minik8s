package main

import (
	"fmt"
	"minik8s.io/cmd/kubelet/app"
	"os"
)

func main() {
	cmd := app.NewKubeletCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
