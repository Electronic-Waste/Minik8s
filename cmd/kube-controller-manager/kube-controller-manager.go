package kube_controller_manager

import (
	"minik8s.io/cmd/kube-controller-manager/app"
)

func main() {
	command := app.NewControllerManagerCommand()
	err := command.Execute()
	panic(err)
}
