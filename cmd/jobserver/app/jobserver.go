package app

import (
	"fmt"
	"minik8s.io/cmd/jobserver/app/cmd"
	"os"
)

func Run() error {
	fmt.Println("kubeadm app run")
	cmds := cmd.NewServerCommand(os.Stdin, os.Stdout, os.Stderr)

	return cmds.Execute()
}
