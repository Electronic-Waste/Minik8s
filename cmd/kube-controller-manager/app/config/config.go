package config

import (
	opts "k8s/cmd/kube-controller-manager/app/controlleroptions"
)

// TODO: add client
type Config struct {
	*opts.ReplicaSetControllerOptions
	*opts.DeploymentControllerOptions
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// TODO: authorize client
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}
	return &CompletedConfig{&cc}
}
