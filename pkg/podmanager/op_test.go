package podmanager

import (
	"fmt"
	"minik8s.io/pkg/apis/core"
	"testing"
)

func TestAllProcess(t *testing.T) {
	// construct a Pod Object
	var pod core.Pod
	pod.Name = "test"
	pod.Kind = "Pod"
	pod.Spec.Volumes = []core.Volume{{
		Name: "shared-data",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/root/test_vo",
			},
		},
	}}
	pod.Spec.Containers = []core.Container{
		{
			Name:  "go1",
			Image: "docker.io/library/golang:latest",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "shared-data",
					MountPath: "/mnt",
				},
			},
			Ports:   []core.ContainerPort{},
			Command: []string{"bash"},
		},
		{
			Name:  "go2",
			Image: "docker.io/library/golang:latest",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "shared-data",
					MountPath: "/go/src",
				},
			},
			Ports:   []core.ContainerPort{},
			Command: []string{"bash"},
		},
	}

	err := pod.ContainerConvert()
	if err != nil {
		t.Error(err)
	}
	err = RunPod(&pod)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("have start a new pod")
	err = DelPod(pod.Name)
	if err != nil {
		t.Error(err)
	}
}

func TestPodRunning(t *testing.T) {
	
}
