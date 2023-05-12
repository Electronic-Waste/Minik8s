package controller

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apps"
	util "minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/util/tools/queue"
	"time"
)

type DeploymentController struct {
	//Client
	//listwatch

	// work queue
	queue   *queue.Queue
	channel <-chan *redis.Message
	message *redis.Message
}

func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue: new(queue.Queue),
	}
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	dc.register()
	go dc.worker(ctx)
}

func (dc *DeploymentController) register() {
	dc.channel = util.Subscribe("pod")
	util.Watch("pod change", dc.listener)
}

func (dc *DeploymentController) listener(msg *redis.Message) {
	bytes := []byte(msg.Payload)
	deployment := apps.Deployment{}
	err := json.Unmarshal(bytes, &deployment)
	if err != nil {
		return
	}
	dc.queue.Enqueue(deployment)
}

func (dc *DeploymentController) worker(ctx context.Context) {
	for {
		if !dc.queue.Empty() {
			dc.processNextWorkItem(ctx)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (dc *DeploymentController) processNextWorkItem(ctx context.Context) {
	key := dc.queue.Dequeue()
	dc.syncDeployment(ctx, key.(apps.Deployment))
	return
}

func (dc *DeploymentController) syncDeployment(ctx context.Context, deployment apps.Deployment) {
	if deployment.Spec.Selector == "" {
		return
	}

	if deployment.Status.AvailableReplicas < deployment.Spec.Replicas {
		bytes, _ := json.Marshal(deployment)
		msg := new(redis.Message)
		msg.Payload = string(bytes)
		util.Publish("add pod", msg)
	}

	if deployment.Status.AvailableReplicas > deployment.Spec.Replicas {
		bytes, _ := json.Marshal(deployment)
		msg := new(redis.Message)
		msg.Payload = string(bytes)
		util.Publish("delete pod", msg)
	}
	//TODO: check the deployment status and do actions accordingly
}

func (dc *DeploymentController) putDeployment(ctx context.Context) {

}
