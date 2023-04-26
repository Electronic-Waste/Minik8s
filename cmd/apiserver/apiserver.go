package main

import (

	"vmeet.io/minik8s/cmd/apiserver/app"
)

func main() {
	command := app.NewAPIServerCommand()
	command.Execute()
}