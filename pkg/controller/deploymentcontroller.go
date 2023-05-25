package controller

import (
	"encoding/json"
	"fmt"
	"context"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/listwatch"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/util/tools/queue"
	"minik8s.io/pkg/util/tools/uid"
	"minik8s.io/pkg/clientutil"
	apiurl "minik8s.io/pkg/apiserver/util/url"
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
	d2pMap 	map[interface{}]interface{}
	p2dMap	map[interface{}]interface{}
	//channel chan struct{}
	//message *redis.Message
}

func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue:   new(queue.Queue),
		d2pMap: make(map[interface{}]interface{}),
		p2dMap:	make(map[interface{}]interface{}),
	}
	print("new deployment controller\n")
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	go dc.register()		//register list watch handler
	//go dc.replicaWatcher()	//supervise pod replica numbers
	go dc.worker(ctx)		//main thread processing messages
	print("deployment controller running\n")
	<-ctx.Done()
}

func (dc *DeploymentController) register() {
	print("dc register\n")
	go listwatch.Watch(apiurl.DeploymentStatusApplyURL, dc.listener)
	go listwatch.Watch(apiurl.DeploymentStatusUpdateURL, dc.listener)
	go listwatch.Watch(apiurl.DeploymentStatusDelURL, dc.listener)
	//not reach here
	print("dc registered\n")
}

func (dc *DeploymentController) listener(msg *redis.Message) {
	print("dc listening\n")
	bytes := []byte(msg.Payload)
	watchres := listwatch.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	dc.queue.Enqueue(watchres)
}

func (dc *DeploymentController) worker(ctx context.Context) {
	print("dc working\n")
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
			did := uid.NewUid()
			prefix := deployment.Metadata.Name + "-" + did
			replicas := deployment.Spec.Replicas
			//label := map[string]string{}
			//label["app"] = "test"
			var nameSet []string
			var containerNameSet []string
			pod := deployment.Spec.Template
			for _,c := range pod.Spec.Containers{
				containerNameSet = append(containerNameSet, c.Name)
			}
			for i := 0; i < replicas; i++ {
				//give pod names
				pid := uid.NewUid()
				podname := prefix + "-" + pid
				dc.p2dMap[podname] = deployment.Metadata.Name
				nameSet = append(nameSet, podname)
				fmt.Println("podname: " + podname)
				pod.Name = podname
				//give container names
				for i,_ := range pod.Spec.Containers{
					cid := uid.NewUid()
					pod.Spec.Containers[i].Name = containerNameSet[i] + "-" + cid
				}
				AddPod(pod)
			}
			dc.d2pMap[deployment.Metadata.Name] = nameSet
		case "update":
		case "delete":
			//client.addPod(pod)
			//var nameSet []string
			nameSet := dc.d2pMap[deployment.Metadata.Name].([]string)
			for i := 0; i < deployment.Status.AvailableReplicas; i++ {
				podname := nameSet[i]
				fmt.Println(podname)
				DelPod(podname)
				delete(dc.p2dMap,podname)
			}
			delete(dc.d2pMap,deployment.Metadata.Name)
		}
	case "Pod":
		pod := core.Pod{}
		err = json.Unmarshal(watchres.Payload, &pod)
		if err != nil {
			return err
		}
		switch actiontype{
		case "apply":
			fmt.Println("apply single pod")
			pid := uid.NewUid()
			podname := pod.Name + "-" + pid
			fmt.Println(podname)
			pod.Name = podname
			AddPod(pod)
		case "update":
		case "delete":
			_,ok := dc.p2dMap[pod.Name]
			//for k,v := range dc.d2pMap{
			//	nameSet := v.([]string)
			//	for _,podname := range nameSet{
			//		if podname == pod.Name{
			//			deploymentname = k.(string)
			//		}
			//	}
			//}
			if ok {
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

func (dc *DeploymentController) replicaWatcher() {
	timeout := time.Second * 3
	for {
		fmt.Println("!!!watching replicas")
		pods,err := podmanager.GetPods()
		if err!=nil{
			fmt.Println(err.Error())
			continue
		}
		var strSet []string
		var deploymentSet []core.Deployment
		bytes,err := clientutil.HttpGetAll("Deployment")
		if err != nil{
			fmt.Println("get deployments fail")
			continue
		}

		json.Unmarshal(bytes,&strSet)
		for _,s := range strSet{
			if s == ""{
				continue
			}
			//fmt.Println(i)
			//fmt.Println(s)
			deployment := core.Deployment{}
			json.Unmarshal([]byte(s),&deployment)
			deploymentSet = append(deploymentSet, deployment)
			fmt.Println(deployment.Metadata.Name)
		}
		//fmt.Println("map:")
		//for k,v := range dc.p2dMap{
		//	fmt.Println("podname:" + k.(string) + ", deployname: " + v.(string))
		//}
		//fmt.Println("pods:")

		numMap := make(map[string]int)
		dc.d2pMap = make(map[interface{}]interface{})
		for _,pod := range pods{
			if pod.Status.Phase != core.PodFailed{
				deploymentname,ok := dc.p2dMap[pod.Name]
				if ok == true{
					//fmt.Println("pod: " + pod.Name + ", deployment: deploymentname")
					replica,ok := numMap[deploymentname.(string)]
					if ok{
						replica++
						//fmt.Println("deployment recorded:")
						//fmt.Println(replica)
						numMap[deploymentname.(string)] = replica
						nameSet := dc.d2pMap[deploymentname.(string)].([]string)
						nameSet = append(nameSet, pod.Name)
						dc.d2pMap[deploymentname.(string)] = nameSet
					}else{
						//fmt.Println("deployment unrecorded:")
						//fmt.Println(1)
						numMap[deploymentname.(string)] = 1
						nameSet := make([]string,0)
						nameSet = append(nameSet, pod.Name)
						dc.d2pMap[deploymentname.(string)] = nameSet
					}
				}
			}
		}

		for deploymentname,replica := range numMap {
			for _,deployment := range deploymentSet{
				if deployment.Metadata.Name == deploymentname{
					if replica < deployment.Spec.Replicas{
						fmt.Println("start adding replicas")
						did := uid.NewUid()
						prefix := deployment.Metadata.Name + "-" + did
						num := deployment.Spec.Replicas - replica
						var nameSet []string
						
						var containerNameSet []string
						pod := deployment.Spec.Template
						for _,c := range pod.Spec.Containers{
							containerNameSet = append(containerNameSet, c.Name)
						}
						for i := 0; i < num; i++ {
							//give pod names
							pid := uid.NewUid()
							podname := prefix + "-" + pid
							dc.p2dMap[podname] = deployment.Metadata.Name
							nameSet = append(nameSet, podname)
							pod := deployment.Spec.Template
							pod.Name = podname
							//give container names
							for i,_ := range pod.Spec.Containers{
								cid := uid.NewUid()
								pod.Spec.Containers[i].Name = containerNameSet[i] + "-" + cid
							}
							AddPod(pod)
						}
						dc.d2pMap[deployment.Metadata.Name] = nameSet
					}
				}
			}
		}
		
		time.Sleep(timeout)
	}
	
}

// just for test
func AddPod(pod core.Pod) {
	fmt.Printf("add pod %s\n",pod.Name)
	//for _,c := range pod.Spec.Containers{
	//	fmt.Printf("pod container %s\n", c.Name)
	//}
	//podmanager.RunPod(&pod)
	clientutil.HttpApply("Pod",pod)
}

func DelPod(podname string) {
	fmt.Printf("del pod %s\n",podname)
	podmanager.DelPod(podname)
}

func GetDeployment(name string) core.Deployment {
	return core.Deployment{}
}