package controller

import (
	//"encoding/json"
	//"context"
	//"minik8s.io/pkg/apis/core"
	util "minik8s.io/pkg/util/listwatch"
	//"time"
	//"github.com/go-redis/redis/v8"
	apiurl "minik8s.io/pkg/apiserver/util/url"

	"testing"
)
/*
func TestDeployment(t *testing.T) {
	ctx := context.Background()
	deploymentController, _ := NewDeploymentController(ctx)
	go deploymentController.Run(ctx)
	time.Sleep(time.Second)

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

	deployment := core.Deployment{}
	deployment.Spec.Replicas = 3
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Selector = "pod"
	deployment.Metadata.Name = "deployment"
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Template = pod

	watchres := listwatch.WatchResult{}
	watchres.ObjectType = "Deployment"
	watchres.ActionType = apply
	watchres.Payload, _ = json.Marshal(deployment)

	bytes, _ := json.Marshal(watchres)
	//payload := string(bytes)

	util.Publish("/api/v1/deployment/status", bytes)
	time.Sleep(time.Second * 3)
}

func TestReplicaset(t *testing.T) {
	ctx := context.Background()
	deploymentController, _ := NewDeploymentController(ctx)

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

	deployment := core.Deployment{}
	deployment.Spec.Replicas = 3
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Selector = "pod"
	deployment.Metadata.Name = "deployment"
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Template = pod

	watchres := listwatch.WatchResult{}
	watchres.ObjectType = "Deployment"
	watchres.ActionType = apply
	watchres.Payload, _ = json.Marshal(deployment)

	deploymentController.syncDeployment(ctx, watchres)

	watchres.ActionType = delete

	deploymentController.syncDeployment(ctx, watchres)
}
*/
func TestApply(t *testing.T) {
	bytes := []byte{}
	util.Publish(apiurl.DeploymentStatusApplyURL, bytes)
	util.Publish(apiurl.DeploymentStatusUpdateURL, bytes)
	util.Publish(apiurl.DeploymentStatusDelURL, bytes)
}