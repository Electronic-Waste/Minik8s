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
		Name:     "shared-data",
		HostPath: "/root/test_vo",
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

func TestGetPod(t *testing.T) {
	podSet, err := GetPods()
	if err != nil {
		t.Error(err)
	}
	for _, pod := range podSet {
		fmt.Printf("the %s's status is %s\n", pod.Name, pod.Status.Phase)
	}
}

func TestPodRunning(t *testing.T) {
	// construct a Pod Object
	var pod core.Pod
	pod.Name = "test"
	pod.Kind = "Pod"
	pod.Spec.Volumes = []core.Volume{{
		Name:     "shared-data",
		HostPath: "/root/test_vo",
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
	is_exist := IsPodRunning("test")
	if !is_exist {
		t.Error("error start up a Pod")
	} else {
		fmt.Println("find test Pod")
	}
	err = DelPod(pod.Name)
	is_exist = IsPodRunning("test")
	if is_exist {
		t.Error("Pod delete error")
	} else {
		fmt.Println("delete Pod success")
	}
	if err != nil {
		t.Error(err)
	}
}
