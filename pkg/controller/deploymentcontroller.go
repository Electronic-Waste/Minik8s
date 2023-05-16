package controller

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apis/core"
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
	print("new deployment controller\n")
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	go dc.register()
	go dc.worker(ctx)
	print("deployment controller running\n")
}

func (dc *DeploymentController) register() {
	print("register\n")
	//dc.channel = util.Subscribe("/api/v1/deployment/status")
	util.Watch("/api/v1/deployment/status", dc.listener)
	//not reach here
	print("registered\n")
}

func (dc *DeploymentController) listener(msg *redis.Message) {
	print("listening\n")
	bytes := []byte(msg.Payload)
	watchres := etcd.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	dc.queue.Enqueue(watchres)
}

func (dc *DeploymentController) worker(ctx context.Context) {
	print("working\n")
	for {
		if !dc.queue.Empty() {
			print("receive msg!\n")
			//dc.queue.Dequeue()
			dc.processNextWorkItem(ctx)
		} else {
			//print("worker pending\n")
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
	//format: pod: deployment-rsuid-poduid
	//expample:	deployment-123456-789456

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
				print(podname)
				pod := core.Pod{
					Kind:   "Pod",
					Spec:   core.PodSpec{},
					Status: core.PodStatus{},
				}
				AddPod(pod)
			}
		case modify:
		case delete:
			//client.addPod(pod)

		}
	case "Pod":
		//seems only delete pod will invoke controller
		pod := core.Pod{}
		err = json.Unmarshal(watchres.Payload, &pod)
		if err != nil {
			return err
		}
		name := pod.Name
		namearr := strings.Split(name, "-")
		deploymentname := namearr[0] + "-" + namearr[1]
		deployment = GetDeployment(deploymentname)
		if deployment.Status.AvailableReplicas < deployment.Spec.Replicas {
			num := deployment.Spec.Replicas - deployment.Status.AvailableReplicas
			for i := 0; i < num; i++ {
				pod := core.Pod{}
				pod = deployment.Spec.Template
				AddPod(pod)
			}
		}

		if deployment.Status.AvailableReplicas > deployment.Spec.Replicas {
			bytes, _ := json.Marshal(deployment)
			msg := new(redis.Message)
			msg.Payload = string(bytes)
			//client
		}
	}
	//TODO: check the deployment status and do actions accordingly
	return nil
}

func (dc *DeploymentController) putDeployment(ctx context.Context) {

}

// just for test
func AddPod(pod core.Pod) {
	print("add pod\n")
}

func GetDeployment(name string) core.Deployment {
	return core.Deployment{}
}
