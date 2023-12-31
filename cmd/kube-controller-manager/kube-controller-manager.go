package main

import (
	"fmt"
	"minik8s.io/cmd/kube-controller-manager/app"
)

func main() {
	command := app.NewControllerManagerCommand()
	err := command.Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println("end of controllers")
}
