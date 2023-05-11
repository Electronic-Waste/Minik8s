package main

import (

	"minik8s.io/cmd/apiserver/app"
)

func main() {
	command := app.NewAPIServerCommand()
	command.Execute()
}