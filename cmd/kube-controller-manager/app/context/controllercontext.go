package controllercontext

import (
	"k8s/cmd/kube-controller-manager/app/config"
)

type ControllerContext struct {
	Config *config.CompletedConfig
}
