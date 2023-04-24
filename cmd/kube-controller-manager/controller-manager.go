package main

import (
	"k8s/cmd/kube-controller-manager/app"
	"os"
)

func main() {
	command := app.NewControllerManagerCommand()
	err := command.Execute()
	os.Exit(err)
}
