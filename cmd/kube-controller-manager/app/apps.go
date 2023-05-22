package app

import (
	"fmt"
	"golang.org/x/net/context"
	"minik8s.io/pkg/controller"
)

func StartDeploymentController(ctx context.Context) error {
	fmt.Printf("start deployment controller\n")
	deploymentController, _ := controller.NewDeploymentController(ctx)
	go deploymentController.Run(ctx)
	<-ctx.Done()
	return nil
}
