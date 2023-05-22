package controller

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/listwatch"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/util/tools/queue"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"strings"
	"time"
)

//const (
//	apply  int = 0
//	update int = 1
//	delete int = 2
//)

type DeploymentController struct {
	//Client
	//util

	// work queue
	queue   *queue.Queue
	nameMap map[interface{}]interface{}
	//channel chan struct{}
	//message *redis.Message
}

func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue:   new(queue.Queue),
		nameMap: make(map[interface{}]interface{}),
	}
	print("new deployment controller\n")
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	go dc.register()
	go dc.worker(ctx)
	print("deployment controller running\n")
	<-ctx.Done()
}

func (dc *DeploymentController) register() {
	print("register\n")
	listwatch.Watch(apiurl.DeploymentStatusApplyURL, dc.listener)
	listwatch.Watch(apiurl.DeploymentStatusUpdateURL, dc.listener)
	listwatch.Watch(apiurl.DeploymentStatusDelURL, dc.listener)
	//not reach here
	//print("registered\n")
}

func (dc *DeploymentController) listener(msg *redis.Message) {
	print("listening\n")
	bytes := []byte(msg.Payload)
	watchres := listwatch.WatchResult{}
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
	_ = dc.syncDeployment(ctx, key.(listwatch.WatchResult))
	return
}

func (dc *DeploymentController) syncDeployment(ctx context.Context, watchres listwatch.WatchResult) error {
	var (
		err        error
		deployment core.Deployment
		actiontype string
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
		case "apply":
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
				pod := deployment.Spec.Template
				pod.Name = podname
				AddPod(pod)
			}
			dc.nameMap[deployment.Metadata.Name] = nameSet
		case "update":
		case "delete":
			//client.addPod(pod)
			//var nameSet []string
			nameSet := dc.nameMap[deployment.Metadata.Name].([]string)
			for i := 0; i < deployment.Status.AvailableReplicas; i++ {
				podname := nameSet[i]
				fmt.Println(podname)
				DelPod(podname)
			}
		}
	case "Pod":
		pod := core.Pod{}
		err = json.Unmarshal(watchres.Payload, &pod)
		if err != nil {
			return err
		}
		switch actiontype{
		case "apply":
		case "update":
		case "delete":
			deploymentname := ""
			for k,v := range dc.nameMap{
				nameSet := v.([]string)
				for _,podname := range nameSet{
					if podname == pod.Name{
						deploymentname = k.(string)
					}
				}
			}
			if deploymentname != ""{
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
		}
		
		
	}
	//TODO: check the deployment status and do actions accordingly
	return nil
}

func (dc *DeploymentController) putDeployment(ctx context.Context) {

}

// just for test
func AddPod(pod core.Pod) {
	fmt.Printf("add pod %s\n",pod.Name)
	//podmanager.RunPod(&pod)
}

func DelPod(podname string) {
	//fmt.Printf("del pod %s\n",podname)
	podmanager.DelPod(podname)
}

func GetDeployment(name string) core.Deployment {
	return core.Deployment{}
}
