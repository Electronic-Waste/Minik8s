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
)

type AutoscalerController struct {
	autoscalerList	[]core.Autoscaler
	cadvisor		*cadvisor.CAdvisor
}

func NewAutoscalerController(ctx context.Context) (*AutoscalerController, error) {
	ac := &AutoscalerController{
		autoscalerList: make([]core.Autoscaler,0),
		cadvisor: 		cadvisor.GetCAdvisor(),
	}
	return ac, nil
}

func (ac *AutoscalerController) Run (ctx context.Context) {
	print("ac run\n")
	go ac.register()
	go ac.worker()
	<-ctx.Done()
}

func (ac *AutoscalerController) register() {
	print("ac register\n")
	go listwatch.Watch(apiurl.AutoscalerStatusApplyURL, ac.applylistener)
	go listwatch.Watch(apiurl.AutoscalerStatusUpdateURL, ac.updatelistener)
	go listwatch.Watch(apiurl.AutoscalerStatusDelURL, ac.deletelistener)
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
	ac.autoscalerList = append(ac.autoscalerList, autoscaler)
}

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

//polling to get deployments
func (ac *AutoscalerController) worker () {
	timeout := time.Second * 5
	for {
		for i,autoscaler := range ac.autoscalerList{
			fmt.Printf("process autoscaler: %s\n", autoscaler.Metadata.Name)
			//check validity of autoscaler first
			if autoscaler.Spec.ScaleTargetRef.Name == "" || autoscaler.Spec.ScaleTargetRef.Kind != "Deployment" {
				ac.autoscalerList = append(ac.autoscalerList[:i],ac.autoscalerList[i+1:]...)
				continue
			}
			//get deployment
			params := make(map[string]string)
			params["namespace"] = "default"
			params["name"] = autoscaler.Spec.ScaleTargetRef.Name
			bytes, err := clientutil.HttpGet("Deployment",params)
			if(err != nil){
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
			params = make(map[string]string)
			params["namespace"] = "default"
			params["prefix"] = deployment.Metadata.Name
			bytes,err = clientutil.HttpGetWithPrefix("Pod",params)
			var strs []string
			var pods []core.Pod
			err = json.Unmarshal(bytes, &strs)
			if err != nil {
				continue
			}
			for _,s := range strs{
				pod := core.Pod{}
				err = json.Unmarshal([]byte(s), &pod)
				if err != nil {
					continue
				}
				fmt.Printf("get pod name: %s\n", pod.Name)
				pods = append(pods, pod)
			}
		}
		time.Sleep(timeout)
	}
	
}