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
			Image: "docker.io/library/pythonplus:latest",
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

func TestSysPod(t *testing.T) {
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
	err = RunSysPod(&pod)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("have start a new pod")
}

func TestJobServer(t *testing.T) {
	// construct a Pod Object
	var pod core.Pod
	pod.Name = "test"
	pod.Kind = "Pod"
	pod.Spec.Volumes = []core.Volume{{
		Name:     "shared-data",
		HostPath: "/root/minik8s/minik8s/scripts/data",
	},
		{
			Name:     "shared-scripts",
			HostPath: "/root/minik8s/minik8s/scripts/gpuscripts",
		}}
	pod.Spec.Containers = []core.Container{
		{
			Name:  "t1",
			Image: "docker.io/library/jobserver:latest",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "shared-data",
					MountPath: "/mnt/data",
				},
				{
					Name:      "shared-scripts",
					MountPath: "/mnt/scripts",
				},
			},
			Ports:   []core.ContainerPort{},
			Command: []string{"/mnt/scripts/jobserver"},
			Args:    []string{"remote", "--file=test.cu", "--scripts=test1.slurm", "--result=result"},
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
	fmt.Printf("have start a new pod and ip is %s\n", pod.Status.PodIp)
	is_exist := IsPodRunning("test")
	if !is_exist {
		t.Error("error start up a Pod")
	} else {
		fmt.Println("find test Pod")
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
	fmt.Printf("have start a new pod and ip is %s\n", pod.Status.PodIp)
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
