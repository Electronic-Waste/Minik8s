package controller

import(
	"context"
	"time"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/kubelet/cadvisor"
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
	cadvisor		*cadvisor.CAdvisor
}

func NewAutoscalerController(ctx context.Context) (*AutoscalerController, error) {
	ac := &AutoscalerController{
		autoscalerList: []core.Autoscaler{},
		cadvisor: 		cadvisor.GetCAdvisor(),
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
	fmt.Println("ac apply new autoscaler")
	//apply
	for _,a := range ac.autoscalerList{
		if a.Metadata.Name == autoscaler.Metadata.Name{
			fmt.Println("cannot apply an autoscaler already exist")
			return
		}
	}
	//start supervise
	pods, err := GetReplicaPods(autoscaler.Spec.ScaleTargetRef.Name)
	if err != nil{
		fmt.Println(err)
		return
	}
	for _,pod := range pods{
		err = ac.cadvisor.RegisterPod(pod.Name)
		if err != nil {
			fmt.Println(err)
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
		timeout = time.Second * 5
	}
	for {
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
		//fmt.Printf("get deployment name: %s\n",deployment.Metadata.Name)
		//get pods
		pods, err := GetReplicaPods(deployment.Metadata.Name)
		if err != nil{
			continue
		}
		//fmt.Println("get pod status:")
		statusList := []stats.PodStats{}
		for _,pod := range pods{
			status, err := ac.cadvisor.GetPodMetric(pod.Name)
			if err != nil{
				fmt.Println(err)
				continue
			}
			statusList = append(statusList, status)
			//fmt.Println(status)
		}
		metricsMap = ac.calculateMetrics(autoscaler, statusList, deployment.Spec.Replicas)

		for metric,value := range metricsMap{
			if metric == "cpu"{
				auto
			}else if metric == "memory"{

			}
		}
		
		time.Sleep(timeout)
	}
}

type metricsCompare struct{
	Metrics		float64	//actual usage
	Utilization	int		//spec
}

//exapmle return: {"cpu":2,"memory":3} or {"cpu":6}
func (ac *AutoscalerController) calculateMetrics(autoscaler core.Autoscaler, status []stats.PodStats, currReplicas int) map[string]float64 {
	metrics := autoscaler.Spec.Metrics
	metricsMap := make(map[string]float64)
	for _,r := range metrics{
		switch r.Resource.Name{
		case "cpu":
			totalcpu := 0.0
			for _,s := range status{
				cpu := s.CPUPercentage
				totalcpu += cpu
			}
			totalcpu *= 100
			fmt.Println("cpu:", totalcpu)
			compare := metricsCompare{
				Metrics: totalcpu,
				Utilization: autoscaler.Spec.Metrics.Resource.Utilization
			}
			metricsMap["cpu"] = totalcpu
		case "memory":
			totalmemory := 0.0
			for _,s := range status{
				memory := s.MemoryPercentage
				totalmemory += memory
			}
			totalmemory *= 100
			//if memory is too small, it will be NAN
			if math.IsNaN(totalmemory){
				totalmemory = 0
			}
			fmt.Println("memory:", totalmemory)
			metricsMap["memory"] = totalmemory
		}
	}
	return metricsMap
}

func GetReplicaPods(deploymentname string) ([]core.Pod,error) {
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