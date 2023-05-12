package kube_controller_manager

import (
	"minik8s.io/cmd/kube-controller-manager/app"
	"os"
)

func main() {
	command := app.NewControllerManagerCommand()
	_ = command.Execute()
	os.Exit(0)
}
