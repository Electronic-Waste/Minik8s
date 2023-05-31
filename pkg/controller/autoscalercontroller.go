package controller

import(
	"context"
	"time"
	"minik8s.io/pkg/clientutil"
	//"minik8s.io/pkg/kubelet/cadvisor"
	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apis/core"
	"github.com/go-redis/redis/v8"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"encoding/json"
	"fmt"
	"minik8s.io/pkg/kubelet/cadvisor/stats"
	"math"
)

type AutoscalerController struct {
	autoscalerList	[]core.Autoscaler
	//cadvisor		*cadvisor.CAdvisor
}

func NewAutoscalerController(ctx context.Context) (*AutoscalerController, error) {
	ac := &AutoscalerController{
		autoscalerList: []core.Autoscaler{},
		//cadvisor: 		cadvisor.GetCAdvisor(),
	}
	return ac, nil
}

func (ac *AutoscalerController) Run (ctx context.Context) {
	print("ac run\n")
	go ac.register()
	go ac.startworker()
	<-ctx.Done()
}

func (ac *AutoscalerController) register() {
	print("ac register\n")
	go listwatch.Watch(apiurl.AutoscalerStatusApplyURL, ac.applylistener)
	go listwatch.Watch(apiurl.AutoscalerStatusUpdateURL, ac.updatelistener)
	go listwatch.Watch(apiurl.AutoscalerStatusDelURL, ac.deletelistener)
	//not reach here
	print("ac registerd\n")
}

func (ac *AutoscalerController) applylistener (msg *redis.Message) {
	//get autoscaler from msg
	fmt.Println("ac apply new autoscaler")
	bytes := []byte(msg.Payload)
	watchres := listwatch.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	//if watchres.ActionType != "apply" || watchres.ObjectType != "Autoscaler"{
	//	return	//won`t happen
	//}
	autoscaler := core.Autoscaler{}
	err = json.Unmarshal(watchres.Payload, &autoscaler)
	if err != nil {
		return
	}
	//apply
	for _,a := range ac.autoscalerList{
		if a.Metadata.Name == autoscaler.Metadata.Name{
			fmt.Println("cannot apply an autoscaler already exist")
			return
		}
	}
	//add autoscaler
	ac.autoscalerList = append(ac.autoscalerList, autoscaler)
	fmt.Println("ac apply new autoscaler success")
}

//no update
func (ac *AutoscalerController) updatelistener (msg *redis.Message) {
	//get autoscaler from msg
	bytes := []byte(msg.Payload)
	watchres := listwatch.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	autoscaler := core.Autoscaler{}
	err = json.Unmarshal(watchres.Payload, &autoscaler)
	if err != nil {
		return
	}
	//update 
	for i,a := range ac.autoscalerList{
		if a.Metadata.Name == autoscaler.Metadata.Name{
			ac.autoscalerList[i] = autoscaler
		}
	}
	fmt.Println("autoscaler to update does not exist")
}

func (ac *AutoscalerController) deletelistener (msg *redis.Message) {
	//get autoscaler from msg
	bytes := []byte(msg.Payload)
	watchres := listwatch.WatchResult{}
	err := json.Unmarshal(bytes, &watchres)
	if err != nil {
		return
	}
	autoscaler := core.Autoscaler{}
	err = json.Unmarshal(watchres.Payload, &autoscaler)
	if err != nil {
		return
	}
	//delete
	for i,a := range ac.autoscalerList{
		if a.Metadata.Name == autoscaler.Metadata.Name{
			ac.autoscalerList = append(ac.autoscalerList[:i],ac.autoscalerList[i+1:]...)
		}
	}
	fmt.Println("autoscaler to delete does not exist")
}

func (ac *AutoscalerController) startworker () {
	for{
		for i,a := range ac.autoscalerList{
			go ac.worker(a)
			ac.autoscalerList = append(ac.autoscalerList[:i],ac.autoscalerList[i+1:]...)
		}
		time.Sleep(time.Second)
	}
}

//polling to get deployments
func (ac *AutoscalerController) worker (autoscaler core.Autoscaler) {
	timeout := time.Second * time.Duration(int64(autoscaler.Spec.ScaleInterval))
	if timeout == 0{
		timeout = time.Second * 15
	}
	for {
		time.Sleep(timeout)

		fmt.Println("ac working")
		//fmt.Println("ac list numbers: ", len(ac.autoscalerList))
		fmt.Printf("process autoscaler: %s\n", autoscaler.Metadata.Name)
		//check validity of autoscaler first
		//if autoscaler.Spec.ScaleTargetRef.Name == "" || autoscaler.Spec.ScaleTargetRef.Kind != "Deployment" {
		//	ac.autoscalerList = append(ac.autoscalerList[:i],ac.autoscalerList[i+1:]...)
		//	continue
		//}
		//get deployment
		params := make(map[string]string)
		params["namespace"] = "default"
		params["name"] = autoscaler.Spec.ScaleTargetRef.Name
		bytes, err := clientutil.HttpGet("Deployment",params)
		if err != nil{
			fmt.Println("get deployment fail")
			continue
		}
		deployment := core.Deployment{}
		err = json.Unmarshal(bytes, &deployment)
		if err != nil {
			continue
		}
		fmt.Printf("get deployment name: %s\n",deployment.Metadata.Name)
		//get pods
		pods, err := GetReplicaPods(deployment.Metadata.Name)
		if err != nil{
			continue
		}
		/*
		for _,pod := range pods{
			err = ac.cadvisor.RegisterPod(pod.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			time.Sleep(time.Second * 5)
		}
		*/
		fmt.Println("get pod status:")
		statusList := []stats.PodStats{}
		for _,pod := range pods{
			//status, err := ac.cadvisor.GetPodMetric(pod.Name)
			status, err := GetPodMetrics(pod.Name, pod.Spec.RunningNode.Spec.NodeIp)
			if err != nil{
				fmt.Println(err)
				continue
			}
			fmt.Println(status)
			statusList = append(statusList, status)
			//fmt.Println(status)
		}
		metricsMap := ac.calculateMetrics(autoscaler, statusList)
		//currentReplicas := deployment.Spec.Replicas
		minreplicas := autoscaler.Spec.MinReplicas
		maxreplicas := autoscaler.Spec.MaxReplicas

		metricsnum := len(metricsMap)
		if metricsnum == 1{
			for name,value := range metricsMap{
				if name == "cpu"{
					if value.Utilization < value.Metrics{
						fmt.Println("cpu increase replicas")
						IncreaseReplicas(deployment,maxreplicas,minreplicas)
					}else{
						fmt.Println("cpu decrease replicas")
						DecreaseReplicas(deployment,maxreplicas,minreplicas)
					}
				}else if name == "memory"{
					if value.Utilization < value.Metrics{
						fmt.Println("memory increase replicas")
						IncreaseReplicas(deployment,maxreplicas,minreplicas)
					}else{
						fmt.Println("memory decrease replicas")
						DecreaseReplicas(deployment,maxreplicas,minreplicas)
					}
				}else{
					fmt.Println("error: unknown metric")
				}
			}
		}else{
			cpuvalue := metricsMap["cpu"]
			memoryvalue := metricsMap["memory"]
			fmt.Println("metrics: ",cpuvalue.Metrics," ",memoryvalue.Metrics)
			fmt.Println("Utilization: ",cpuvalue.Utilization," ",memoryvalue.Utilization)
			flag1 := cpuvalue.Metrics > cpuvalue.Utilization
			flag2 := memoryvalue.Metrics > memoryvalue.Utilization
			if flag1 == flag2{
				if flag1 == true{
					fmt.Println("cpu and memory increase replicas")
					IncreaseReplicas(deployment,maxreplicas,minreplicas)
				}else{
					fmt.Println("cpu and memory decrease replicas")
					DecreaseReplicas(deployment,maxreplicas,minreplicas)
				}
			}else if flag1 == true{	//cpu increase but memory decrease
				fmt.Println("cpu (memory) decrease replicas")
				DecreaseReplicas(deployment,maxreplicas,minreplicas)
			}else{	//cpu decrease but memory increase
				if memoryvalue.Metrics > 80{
					DecreaseReplicas(deployment,maxreplicas,minreplicas)
				}else{
					IncreaseReplicas(deployment,maxreplicas,minreplicas)
				}
			}
		}
	}
}

type MetricsCompare struct{
	Metrics		float64		//actual usage
	Utilization	float64		//spec
}

//exapmle return: {"cpu":2,"memory":3} or {"cpu":6}
func (ac *AutoscalerController) calculateMetrics(autoscaler core.Autoscaler, status []stats.PodStats) map[string]MetricsCompare {
	metrics := autoscaler.Spec.Metrics
	metricsMap := make(map[string]MetricsCompare)
	for _,r := range metrics{
		switch r.Resource.Name{
		case "cpu":
			totalcpu := 0.0
			for _,s := range status{
				cpu := s.CPUPercentage
				totalcpu += cpu
			}
			fmt.Println("cpu:", totalcpu)
			compare := MetricsCompare{
				Metrics: totalcpu,
				Utilization: float64(r.Resource.Utilization),
			}
			metricsMap["cpu"] = compare
		case "memory":
			totalmemory := 0.0
			for _,s := range status{
				memory := s.MemoryPercentage
				totalmemory += memory
			}
			//if memory is too small, it will be NAN
			if math.IsNaN(totalmemory){
				totalmemory = 0
			}
			fmt.Println("memory:", totalmemory)
			compare := MetricsCompare{
				Metrics: totalmemory,
				Utilization: float64(r.Resource.Utilization),
			}
			metricsMap["memory"] = compare
		}
	}
	return metricsMap
}

func GetReplicaPods(deploymentname string) ([]core.Pod,error) {
	fmt.Println("GetReplicaPods")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["prefix"] = deploymentname
	bytes,err := clientutil.HttpGetWithPrefix("Pod",params)
	var strs []string
	var pods []core.Pod
	err = json.Unmarshal(bytes, &strs)
	if err != nil {
		return nil,err
	}
	for _,s := range strs{
		if s == ""{
			continue
		}
		pod := core.Pod{}
		err = json.Unmarshal([]byte(s), &pod)
		if err != nil {
			return nil,err
		}
		//fmt.Printf("get pod name: %s\n", pod.Name)
		pods = append(pods, pod)
	}
	return pods,nil
}
//maxreplicas = -1 representing infinity
func IncreaseReplicas(deployment core.Deployment, maxreplicas int, minreplicas int){
	if maxreplicas != -1 && deployment.Spec.Replicas == maxreplicas{
		fmt.Println("reach maxreplicas")
		return
	}
	deployment.Spec.Replicas = deployment.Spec.Replicas + 1
	fmt.Println("increase deployment ",deployment.Metadata.Name," to ",deployment.Spec.Replicas)
	clientutil.HttpUpdate("Deployment",deployment)
}

func DecreaseReplicas(deployment core.Deployment, maxreplicas int, minreplicas int){
	if deployment.Spec.Replicas == minreplicas{
		fmt.Println("reach minreplicas")
		return
	}
	deployment.Spec.Replicas = deployment.Spec.Replicas - 1
	fmt.Println("decrease deployment ",deployment.Metadata.Name," to ",deployment.Spec.Replicas)
	clientutil.HttpUpdate("Deployment",deployment)
}

func GetPodMetrics(podname string, nodeIP string) (stats.PodStats,error) {
	fmt.Println("GetPodMetrics ",podname," ",nodeIP)
	params := make(map[string]string)
	params["name"] = podname
	params["nodeip"] = nodeIP
	bytes,err := clientutil.HttpGet("metrics",params)
	if err != nil{
		return stats.PodStats{}, err
	}
	status := stats.PodStats{}
	err = json.Unmarshal(bytes, &status)
	if err != nil{
		return stats.PodStats{}, err
	}
	return status,nil
}