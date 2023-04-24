package app

import (
	"context"
	"fmt"
	"k8s/cmd/kube-controller-manager/controllercontext"
	"k8s/pkg/controller/deployment"
	"k8s/pkg/controller/replicaset"
)

// TODO: add informer
func StartReplicaSetController(ctx context.Context, controllerCtx controllercontext.ControllerContext) error {
	fmt.Printf("start running replicaset controller\n")
	replicasetController := replicaset.NewReplicaSetController(ctx, controllerCtx)
	go replicasetController.Run(ctx)
	return nil
}

func StartDeploymentController(ctx context.Context, controllerCtx controllercontext.ControllerContext) error {
	fmt.Printf("start running deployment controller\n")
	deploymentController := deployment.NewDeploymentController(ctx, controllerCtx)
	go deploymentController.Run(ctx)
	return nil
}
