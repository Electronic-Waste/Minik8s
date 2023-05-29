package controller

import (
	"encoding/json"
	"fmt"
	"context"
	"minik8s.io/pkg/apis/core"
	//"minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/listwatch"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/util/tools/queue"
	"minik8s.io/pkg/util/tools/uid"
	"minik8s.io/pkg/clientutil"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	svlurl "minik8s.io/pkg/serverless/util/url"
	"time"
	"strings"
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
	isApplying int
	//channel chan struct{}
	//message *redis.Message
}

func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue:   new(queue.Queue),
		d2pMap: make(map[interface{}]interface{}),	//deploymentname: []podname
		p2dMap:	make(map[interface{}]interface{}),	//podname:	deploymentname
	}
	print("new deployment controller\n")
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	dc.restart()			//restart from crash
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
	//go listwatch.Watch(apiurl.PodStatusApplyURL, dc.podApplyListener)
	//not reach here
	print("dc registered\n")
}
/*
func (dc *DeploymentController) podApplyListener (msg *redis.Message) {
	fmt.Println("dc informed pod apply")
	dc.isApplying = 1
}
*/

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
			print("dc receive wacth res!\n")
			//dc.queue.Dequeue()
			dc.processNextWorkItem(ctx)
		} else {
			//print("worker pending\n")
			dc.replicaWatcher()
		}
		time.Sleep(time.Second * 2)
	}
}

func (dc *DeploymentController) processNextWorkItem(ctx context.Context) {
	key := dc.queue.Dequeue()
	_ = dc.syncDeployment(ctx, key.(listwatch.WatchResult))
	//wait for pods to be truly applied
	//time.Sleep(time.Second * 30)
	return
}

func (dc *DeploymentController) syncDeployment(ctx context.Context, watchres listwatch.WatchResult) error {
	var (
		err        error
		deployment core.Deployment
		actiontype string
		objecttype string
	)
	//format: pod: deployment-poduid
	//expample:	deployment-789456
	fmt.Println("sync deployment")
	actiontype = watchres.ActionType
	objecttype = watchres.ObjectType
	fmt.Println("type:",actiontype,objecttype)
	switch objecttype {
	case "Deployment":
		var deploymentname string
		//TODO: add selector
		//if deployment.Spec.Selector == "" {
		//	return nil
		//}
		switch actiontype {
		case "apply":
			deployment = core.Deployment{}
			err = json.Unmarshal(watchres.Payload, &deployment)
			if err != nil {
				fmt.Println(err)
				return err
			}
			deploymentname = deployment.Metadata.Name
			fmt.Println("apply deployment pods")
			//test duplicate deployment !!!test
			_, isfound := dc.d2pMap[deployment.Metadata.Name]
			if isfound {
				fmt.Println("cannot apply duplicate deployment")
				return nil
			}
			
			replicas := deployment.Spec.Replicas
			//label := map[string]string{}
			//label["app"] = "test"
			var nameSet []string
			var containerNameSet []string
			pod := deployment.Spec.Template
			//!!!test
			prefix := deployment.Metadata.Name
			fmt.Println("deployment apply pod:",pod)
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
		case "update"://should only modify replicas
			deployment = core.Deployment{}
			err = json.Unmarshal(watchres.Payload, &deployment)
			if err != nil {
				fmt.Println(err)
				return err
			}
			deploymentname = deployment.Metadata.Name
			fmt.Println("update deployment")
			nameSet := dc.d2pMap[deployment.Metadata.Name].([]string)
			oldReplicas := len(nameSet)
			newReplicas := deployment.Spec.Replicas
			fmt.Println("dc update from",oldReplicas,"to",newReplicas)

			pod := deployment.Spec.Template
			//!!!test
			prefix := deployment.Metadata.Name
			var containerNameSet []string
			for _,c := range pod.Spec.Containers{
				containerNameSet = append(containerNameSet, c.Name)
			}
			if oldReplicas < newReplicas{
				num := newReplicas - oldReplicas
				
				for i := 0; i < num; i++{
					pid := uid.NewUid()
					podname := prefix + "-" + pid
					dc.p2dMap[podname] = deployment.Metadata.Name
					dc.d2pMap[deployment.Metadata.Name] = append(nameSet, podname)
					pod.Name = podname
					for i,_ := range pod.Spec.Containers{
						cid := uid.NewUid()
						pod.Spec.Containers[i].Name = containerNameSet[i] + "-" + cid
					}
					AddPod(pod)
					fmt.Println("deployment update add pod")
				}
			}else{
				num := oldReplicas - newReplicas
				for i := 0; i < num; i++{
					podname := nameSet[0]
					dc.d2pMap[deployment.Metadata.Name] = nameSet[1:]
					delete(dc.p2dMap, podname)
					DelPod(podname)
					fmt.Println("deployment update delete pod")
				}
			}
		case "delete":
			err = json.Unmarshal(watchres.Payload, &deploymentname)
			if err != nil {
				fmt.Println(err)
				return err
			}
			//client.addPod(pod)
			//var nameSet []string
			fmt.Println("dc delete deployment")
			nameSet := dc.d2pMap[deploymentname].([]string)
			for i := 0; i < len(nameSet); i++ {
				podname := nameSet[i]
				fmt.Println("dc delete pod:",podname)
				DelPod(podname)
				delete(dc.p2dMap,podname)
			}
			delete(dc.d2pMap,deploymentname)
		default:
			fmt.Println("unknown action type")
		}
		//wait for pod to be started
		//deploymentname := deployment.Metadata.Name
		
		if actiontype != "delete"{
			//long wait
			if strings.HasPrefix(deploymentname, svlurl.DeploymentNamePrefix){
				fmt.Println("long wait")
				time.Sleep(time.Second * 30)
			}else{	//short wait
				fmt.Println("short wait")
				time.Sleep(time.Second * 5)
			}
		}
	}
	//TODO: check the deployment status and do actions accordingly
	return nil
}

func (dc *DeploymentController) replicaWatcher() {
	fmt.Println("!!!watching replicas")
	pods,err := GetPods()
	fmt.Println("replica watcher get all pods:", pods)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	if len(pods) == 0{
		return
	}
	var strSet []string
	var deploymentSet []core.Deployment
	bytes,err := clientutil.HttpGetAll("Deployment")
	if err != nil{
		fmt.Println("get deployments fail")
		return
	}
	json.Unmarshal(bytes,&strSet)
	for _,s := range strSet{
		if s == ""{
			continue
		}
		deployment := core.Deployment{}
		json.Unmarshal([]byte(s),&deployment)
		deploymentSet = append(deploymentSet, deployment)
		fmt.Println("rw get", deployment.Metadata.Name)
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
				fmt.Println("pod:", pod.Name, "deployment:", deploymentname)
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
		}else{
			fmt.Println("pod", pod.Name,"is failed")
		}
	}
	for deploymentname,replica := range numMap {
		for _,deployment := range deploymentSet{
			if deployment.Metadata.Name == deploymentname{
				fmt.Println("replica watcher deployment:", deploymentname,"has",replica,"replicas","should be",deployment.Spec.Replicas)
				if replica < deployment.Spec.Replicas{
					fmt.Println("start adding replicas")
					//did := uid.NewUid()
					num := deployment.Spec.Replicas - replica
					var nameSet []string
					
					var containerNameSet []string
					pod := deployment.Spec.Template
					//!!!test
					prefix := deployment.Metadata.Name
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
					//wait for pod to be started
					deploymentname = deployment.Metadata.Name
					//long wait
					if strings.HasPrefix(deploymentname, svlurl.DeploymentNamePrefix){
						fmt.Println("long wait")
						time.Sleep(time.Second * 30)
					}else{	//short wait
						fmt.Println("short wait")
						time.Sleep(time.Second * 5)
					}
					dc.d2pMap[deployment.Metadata.Name] = nameSet
				}
			}
		}
	}
}

//this is an asynchronized method!
func AddPod(pod core.Pod) {
	fmt.Printf("add pod %s\n",pod.Name)
	err := clientutil.HttpApply("Pod",pod)
	if err != nil{
		fmt.Println(err)
	}
}
//this is a synchronized method!
func DelPod(podname string) {
	fmt.Printf("del pod %s\n",podname)
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = podname
	err := clientutil.HttpDel("Pod",params)
	if err != nil{
		fmt.Println("deployment del pod error")
		fmt.Println(err)
	}
}

func GetPods() ([]core.Pod,error) {
	fmt.Println("get all pods running")
	bytes,err := clientutil.HttpGetAll("Pod")
	var pods []core.Pod
	err = json.Unmarshal(bytes, &pods)
	if err != nil{
		return nil,err
	}
	return pods,nil
}

//!!!test
func (dc *DeploymentController) restart(){
	//get all deployments
	fmt.Println("dc restart")
	var strSet []string
	var deploymentSet []core.Deployment
	bytes,err := clientutil.HttpGetAll("Deployment")
	if err != nil{
		fmt.Println(err)
		fmt.Println("dc restart get deployments fail")
		return
	}
	err = json.Unmarshal(bytes, &strSet)
	if err != nil{
		fmt.Println(err)
		fmt.Println("dc restart unmarshal deployments fail")
		return
	}
	for _,s := range strSet{
		if s == ""{
			continue
		}
		deployment := core.Deployment{}
		json.Unmarshal([]byte(s),&deployment)
		deploymentSet = append(deploymentSet, deployment)
		fmt.Println(deployment.Metadata.Name)
	}
	//get all pods of one deployment and record into d2pmap and p2dmap
	for _,d := range deploymentSet{
		pods,err := GetReplicaPods(d.Metadata.Name)
		if err != nil{
			fmt.Println(err)
			fmt.Println("dc restart fail get replica pods of:",d.Metadata.Name)
			continue
		}
		/*
		params := make(map[string]string)
		params["namespace"] = "default"
		params["prefix"] = d.Metadata.Name
		bytes,err := clientutil.HttpGetWithPrefix("Pod",params)
		var pods []core.Pod
		err = json.Unmarshal(bytes, &pods)
		if err != nil{
			fmt.Println(err)
			fmt.Println("dc restart get replica pods fail")
		}
		*/
		podnameSet := make([]string,0)
		for _,p := range pods{
			podnameSet = append(podnameSet, p.Name)
			dc.p2dMap[p.Name] = d.Metadata.Name
			fmt.Println("dc add",p.Name,"to p2dMap")
		}
		dc.d2pMap[d.Metadata.Name] = podnameSet
		fmt.Println("dc add",podnameSet,"to d2pMap")
	}
}
