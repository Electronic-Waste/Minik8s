package app

import (
	"fmt"
	"os"

	"minik8s.io/cmd/kubeadm/app/cmd"
)

func Run() error {
	fmt.Println("kubeadm app run")
	cmds := cmd.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr)

	return cmds.Execute()
}
