package app

import (
	"fmt"
	"context"
	"minik8s.io/pkg/controller"
)

func StartDeploymentController(ctx context.Context) error {
	fmt.Printf("start deployment controller\n")
	deploymentController, _ := controller.NewDeploymentController(ctx)
	go deploymentController.Run(ctx)
	<-ctx.Done()
	return nil
}

func StartAutoSclaerController(ctx context.Context) error {
	fmt.Printf("start deployment controller\n")
	autoscalerController, _ := controller.NewAutoscalerController(ctx)
	go autoscalerController.Run(ctx)
	<-ctx.Done()
	return nil
}