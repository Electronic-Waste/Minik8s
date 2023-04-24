package options

import (
	"fmt"
	"k8s/cmd/kube-controller-manager/app/config"
	opts "k8s/cmd/kube-controller-manager/app/controlleroptions"
)

type ControllerOptions interface {
	InitOptions()
}

type KubeControllerManagerOptions struct {
	// TODO : add more controllers here
	ReplicaSetController *opts.ReplicaSetControllerOptions
	DeploymentController *opts.DeploymentControllerOptions
}

/* it's same as KubeControllerManagerOptions, so it is useless
// KubeControllerManagerConfiguration contains elements describing kube-controller manager.
type KubeControllerManagerConfiguration struct {
	// DeploymentControllerConfiguration holds configuration for DeploymentController related features.
	DeploymentController DeploymentControllerConfiguration
	// ReplicaSetControllerConfiguration holds configuration for ReplicaSet related features.
	ReplicaSetController ReplicaSetControllerConfiguration
}
*/
// new options for controllers
func NewKubeControllerManagerOptions() (*KubeControllerManagerOptions, error) {
	//componentConfig, err := NewDefaultComponentConfig()
	s, err := NewDefaultComponentOptions()
	if err != nil {
		fmt.Printf("new options error\n")
	}
	return &s, nil
}

// get default controllers options
func NewDefaultComponentOptions() (KubeControllerManagerOptions, error) {
	versioned := KubeControllerManagerOptions{}
	versioned.DeploymentController = &opts.DeploymentControllerOptions{}
	versioned.ReplicaSetController = &opts.ReplicaSetControllerOptions{}

	InitOptions(versioned.DeploymentController)
	InitOptions(versioned.ReplicaSetController)
}

func InitOptions(controllerOptions *ControllerOptions) {
	controllerOptions.InitOptions()
}

// TODO: add more config and create a client
// get config of options
func (opts *KubeControllerManagerOptions) Config() *config.Config {
	// TODO : finish this function
	return &config.Config{
		ReplicaSetControllerOptions: opts.ReplicaSetController,
		DeploymentControllerOptions: opts.DeploymentController,
	}
}
