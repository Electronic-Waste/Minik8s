package config

import "minik8s.io/pkg/kubelet/types"

type PodConfig struct {
	// that is the last channel used to hold all Pod message
	update chan types.PodUpdate
}

func (p *PodConfig) Updates() chan types.PodUpdate {
	return p.update
}

func NewPodConfig() *PodConfig {
	ch := make(chan types.PodUpdate, 50)
	return &PodConfig{
		update: ch,
	}
}
