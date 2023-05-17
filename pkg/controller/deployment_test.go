package controller

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	util "minik8s.io/pkg/util/listwatch"
	"time"

	"testing"
)

func TestDeployment(t *testing.T) {
	ctx := context.Background()
	deploymentController, _ := NewDeploymentController(ctx)
	go deploymentController.Run(ctx)
	time.Sleep(1000)

	msg := redis.Message{}

	watchres := etcd.WatchResult{}
	watchres.ObjectType = "Deployment"
	watchres.ActionType = apply

	deployment := core.Deployment{}
	deployment.Spec.Replicas = 3
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Selector = "pod"
	deployment.Metadata.Name = "deployment"
	deployment.Status.AvailableReplicas = 3

	watchres.Payload, _ = json.Marshal(deployment)

	bytes, _ := json.Marshal(watchres)
	msg.Payload = string(bytes)

	util.Publish("/api/v1/deployment/status", msg)
}

func TestReplicaset(t *testing.T) {
	ctx := context.Background()
	deploymentController, _ := NewDeploymentController(ctx)

	watchres := etcd.WatchResult{}
	watchres.ObjectType = "Deployment"
	watchres.ActionType = apply

	deployment := core.Deployment{}
	deployment.Spec.Replicas = 3
	deployment.Status.AvailableReplicas = 3
	deployment.Spec.Selector = "pod"
	deployment.Metadata.Name = "deployment"
	deployment.Status.AvailableReplicas = 3

	watchres.Payload, _ = json.Marshal(deployment)

	deploymentController.syncDeployment(ctx, watchres)

	watchres.ActionType = delete

	deploymentController.syncDeployment(ctx, watchres)
}
