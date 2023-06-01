package main

import (
	"fmt"
	"minik8s.io/cmd/knative/app"
	"os"
)

func main() {
	cmd := app.NewKnativeCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}