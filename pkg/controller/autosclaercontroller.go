package controller

import(
	"context"
)

type AutoscalerController struct {
	DeploymentName string
}

func NewAutoscalerController(ctx context.Context) (*AutoscalerController, error) {
	ac := &AutoscalerController{
		DeploymentName: "test",
	}
	return ac, nil
}

func (ac *AutoscalerController) Run (ctx context.Context) {
	return
}