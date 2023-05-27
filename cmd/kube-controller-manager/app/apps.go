package app

import (
	"context"
	"fmt"
	"minik8s.io/pkg/controller"
)

func StartDeploymentController(ctx context.Context) error {
	fmt.Printf("start deployment controller\n")
	deploymentController, _ := controller.NewDeploymentController(ctx)
	go deploymentController.Run(ctx)
	return nil
}

func StartAutoSclaerController(ctx context.Context) error {
	fmt.Printf("start autoscaler controller\n")
	autoscalerController, _ := controller.NewAutoscalerController(ctx)
	go autoscalerController.Run(ctx)
	return nil
}

func StartJobController(ctx context.Context) error {
	fmt.Printf("start job controller\n")
	
	return nil
}
