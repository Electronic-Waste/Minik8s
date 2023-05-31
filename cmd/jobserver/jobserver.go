package main

import "minik8s.io/cmd/jobserver/app"

func main() {
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
