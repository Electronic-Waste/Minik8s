package controller

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	"minik8s.io/pkg/apiserver/etcd"
	util "minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/util/tools/queue"
	"strings"
	"time"
)

const (
	apply  int = 0
	modify int = 1
	delete int = 2
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
	watchres := etcd.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	dc.queue.Enqueue(watchres)
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
	_ = dc.syncDeployment(ctx, key.(etcd.WatchResult))
	return
}

func (dc *DeploymentController) syncDeployment(ctx context.Context, watchres etcd.WatchResult) error {
	var (
		err        error
		deployment core.Deployment
		actiontype int
		objecttype string
	)
	actiontype = watchres.ActionType
	objecttype = watchres.ObjectType
	switch objecttype {
	case "Deployment":
		deployment = core.Deployment{}
		err = json.Unmarshal(watchres.Payload, &deployment)
		if err != nil {
			return err
		}
		if deployment.Spec.Selector == "" {
			return nil
		}
		switch actiontype {
		case apply:
			uid := uuid.New()
			prefix := deployment.Metadata.Name + "-" + uid.String()
			replicas := deployment.Spec.Replicas
			label := map[string]string{}
			label["app"] = "test"
			for i := 0; i < replicas; i++ {
				pid := uuid.New()
				podname := prefix + "-" + pid.String()
				metadata := meta.ObjectMeta{
					Name:   podname,
					Labels: label,
				}
				pod := core.Pod{
					Name:   podname,
					Labels: label,
					Kind:   "Pod",
					Spec:   core.PodSpec{},
					Status: core.PodStatus{},
				}
				//client.addPod(pod)
			}
		case modify:
		case delete:
			//client.addPod(pod)

		}
	case "Pod":
		pod := core.Pod{}
		err = json.Unmarshal(watchres.Payload, &pod)
		if err != nil {
			return err
		}
		name := pod.Name
		namearr := strings.Split(name, "-")
		deploymentname := namearr[0] + "-" + namearr[1]
		//deployment = client.getDeployment(deploymentname)
		deployment = core.Deployment{}
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
	}
	//TODO: check the deployment status and do actions accordingly
	return nil
}

func (dc *DeploymentController) putDeployment(ctx context.Context) {

}
