package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/podmanager"
	util "minik8s.io/pkg/util/listwatch"
	_map "minik8s.io/pkg/util/tools/map"
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
	nameMap *_map.ConcurrentMap
	//channel chan struct{}
	//message *redis.Message
}

func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue:   new(queue.Queue),
		nameMap: _map.NewConcurrentMap(),
	}
	print("new deployment controller\n")
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	go dc.register()
	go dc.worker(ctx)
	<-ctx.Done()
	print("deployment controller running\n")

}

func (dc *DeploymentController) register() {
	print("register\n")
	util.Watch("/api/v1/deployment/status", dc.listener)
	//not reach here
	//print("registered\n")
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
	fmt.Println("sync deployment")
	actiontype = watchres.ActionType
	objecttype = watchres.ObjectType
	fmt.Println("type: " + objecttype)
	switch objecttype {
	case "Deployment":
		deployment = core.Deployment{}
		err = json.Unmarshal(watchres.Payload, &deployment)
		if err != nil {
			return err
		}
		//TODO: add selector
		//if deployment.Spec.Selector == "" {
		//	return nil
		//}
		switch actiontype {
		case apply:
			fmt.Println("apply deployment pods")
			uid := uuid.New()
			uidstr := strings.Split(uid.String(), "-")[0]
			prefix := deployment.Metadata.Name + "-" + uidstr
			replicas := deployment.Spec.Replicas
			label := map[string]string{}
			label["app"] = "test"
			var nameSet []string
			for i := 0; i < replicas; i++ {
				pid := uuid.New()
				pidstr := strings.Split(pid.String(), "-")[0]
				podname := prefix + "-" + pidstr
				nameSet = append(nameSet, podname)
				fmt.Println(podname)
				pod := core.Pod{
					Kind:   "Pod",
					Spec:   core.PodSpec{},
					Status: core.PodStatus{},
				}
				AddPod(pod)
			}
			dc.nameMap.Put(deployment.Metadata.Name, nameSet)
		case modify:
		case delete:
			//client.addPod(pod)
			//var nameSet []string
			nameSet := dc.nameMap.Get(deployment.Metadata.Name).([]string)
			for i := 0; i < deployment.Status.AvailableReplicas; i++ {
				podname := nameSet[i]
				fmt.Println(podname)
				DelPod(podname)
			}
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
	//fmt.Println("add pod")
	podmanager.RunPod(&pod)
}

func DelPod(podname string) {
	podmanager.DelPod(podname)
}

func GetDeployment(name string) core.Deployment {
	return core.Deployment{}
}
