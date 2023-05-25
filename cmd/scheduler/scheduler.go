package main

import (
	"fmt"
	"minik8s.io/cmd/scheduler/app"
	"os"
)

func main() {
	cmd := app.NewSchedulerCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
