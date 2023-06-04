package functioncontroller

import(
	"minik8s.io/pkg/apis/core"
	//"minik8s.io/pkg/apis/meta"
	"minik8s.io/pkg/controller"
	"context"
	"minik8s.io/pkg/util/listwatch"
	"sync"
	"time"
	"minik8s.io/pkg/clientutil"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	"fmt"
	"minik8s.io/pkg/serverless/util/url"
)

var(
	FunctionStatus	 string = url.Function
	FunctionRegisterUrl string = url.FunctionRegisterURL
	FunctionDelUrl string = url.FunctionDelURL
	FunctionTriggerUrl	string = url.FunctionTriggerURL
	defaultcountdown int = 60
)

type FunctionController struct{
	deploymentMap		map[string]string	//record function name to deployment name
	replicaMap			map[string]int
	countdownMap 		map[string]int
	requestMap			map[string]int
	mutex 				sync.Mutex
}

func NewFunctionController() (*FunctionController,error) {
	fc := &FunctionController{
		deploymentMap: 	make(map[string]string),
		countdownMap:	make(map[string]int),
		replicaMap:		make(map[string]int),
		requestMap:		make(map[string]int),
		mutex: 			sync.Mutex{},
	}
	return fc,nil
}

func (fc *FunctionController) Run (ctx context.Context) {
	fmt.Println("fc running")
	go fc.register()
	//go fc.scaler()
	go fc.countdown()
	<-ctx.Done()
}

func (fc *FunctionController) register () {
	//fmt.Println("fc register")
	go listwatch.Watch(FunctionRegisterUrl, fc.registerlistener)
	go listwatch.Watch(FunctionTriggerUrl, fc.triggerlistener)
	go listwatch.Watch(FunctionDelUrl, fc.deletelistener)
}

func (fc *FunctionController) registerlistener (msg *redis.Message) {
	print("fc listen register\n")
	bytes := []byte(msg.Payload)
	var functionname string
	err := json.Unmarshal(bytes, &functionname)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("fc register functionname:",functionname)
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.replicaMap[functionname] = 1
	fc.deploymentMap[functionname] = url.DeploymentNamePrefix + functionname
	fc.countdownMap[functionname] = defaultcountdown + 30
	fc.requestMap[functionname] = 0
}

func (fc *FunctionController) triggerlistener (msg *redis.Message) {
	print("fc listen trigger\n")
	bytes := []byte(msg.Payload)
	var functionname string
	err := json.Unmarshal(bytes, &functionname)
	if err != nil {
		fmt.Println(err)
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.countdownMap[functionname] = defaultcountdown
	fc.requestMap[functionname] += 1
	fmt.Println("request/s:",fc.requestMap[functionname])
	//fc.mutex.Unlock()
	//scale from 0
	if fc.replicaMap[functionname] == 0 {
		fmt.Println("scale from 0")
		fc.IncreaseReplica(functionname)
		//wait for pod to start
		time.Sleep(time.Second * 30)
	}
	//if too many requests, add replicas 
	if fc.requestMap[functionname] >= 2{
		fc.IncreaseReplica(functionname)
		time.Sleep(time.Second * 30)
	}
}

func (fc *FunctionController) deletelistener (msg *redis.Message) {
	print("fc listen delete\n")
	bytes := []byte(msg.Payload)
	var functionname string
	err := json.Unmarshal(bytes, &functionname)
	if err != nil {
		fmt.Println(err)
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	delete(fc.deploymentMap, functionname)
	delete(fc.replicaMap, functionname)
	delete(fc.countdownMap, functionname)
	delete(fc.requestMap, functionname)
}

func (fc *FunctionController) countdown () {
	for{
		time.Sleep(time.Second)
		//fmt.Println("counting")
		fc.mutex.Lock()
		for f,c := range fc.countdownMap{
			fc.countdownMap[f] = c - 1
			fc.requestMap[f] = 0
			fmt.Println("countdown",f, fc.countdownMap[f],"replicas:",fc.replicaMap[f])
			//scale to 0, delete function
			if fc.countdownMap[f] == 0 {
				fmt.Println("scale to 0")
				err := fc.ScaleTo0(f)
				if err != nil{
					fmt.Println(err)
					continue
				}
			}
		}
		fc.mutex.Unlock()
	}
}

func (fc *FunctionController) ScaleTo0 (functionname string) error {
	fmt.Println("ScaleTo0",functionname)
	delete(fc.countdownMap, functionname)
	fc.replicaMap[functionname] = 0
	//get deployment
	deployment, err := GetDeploymentByName(fc.deploymentMap[functionname])
	if err != nil{
		return err
	}
	//set replica to 0
	deployment.Spec.Replicas = 0
	fmt.Println("fc set deployment",deployment.Metadata.Name,"to 0")
	clientutil.HttpUpdate("Deployment",deployment)
	return nil
}

//attention: there is no 's'
func (fc *FunctionController)IncreaseReplica(functionname string) {
	deploymentname := fc.deploymentMap[functionname]
	fmt.Println("increase replica", functionname, deploymentname)
	fc.replicaMap[functionname] += 1
	//get deployment
	deployment, err := GetDeploymentByName(deploymentname)
	if err != nil{
		fmt.Println(err)
		return
	}
	//change deployment
	controller.IncreaseReplicas(deployment, -1, 1)
}

//attention: there is no 's'
func (fc *FunctionController)DecreaseReplica(functionname string) {
	deploymentname := fc.deploymentMap[functionname]
	fmt.Println("decrease replica", functionname, deploymentname)
	fc.replicaMap[functionname] -= 1
	//get deployment
	deployment, err := GetDeploymentByName(deploymentname)
	if err != nil{
		fmt.Println(err)
		return
	}
	//change deployment
	controller.DecreaseReplicas(deployment, -1, 1)
}

func GetDeploymentByName (deploymentname string) (core.Deployment,error) {
	//get deployment
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = deploymentname
	bytes, err := clientutil.HttpGet("Deployment",params)
	if err != nil{
		fmt.Println("get deployment fail")
		return core.Deployment{},err
	}
	deployment := core.Deployment{}
	err = json.Unmarshal(bytes, &deployment)
	if err != nil {
		return core.Deployment{},err
	}
	fmt.Printf("fc get deployment name: %s\n",deployment.Metadata.Name)
	return deployment, nil
}
/*
func (fc *FunctionController) scaler () {
	timeout := time.Second * 5
	for{
		time.Sleep(timeout)
		for _,d := fc.deploymentMap{
			metrics := GetCpuMetrics(d.Metadata.Name)
			if metrics > 50 {
				IncreaseReplica(d.Meatadata.Name)
				//autoscaler.IncreaseReplicas(d.Metadata.Name,-1,1)
			}
			if metrics {
				DecreaseReplica(d.Meatadata.Name)
				//autoscaler.DecreaseReplicas(d.Metadata.Name,-1,1)
			}
		}
	}
}
func GetMetrics () float64 {
	return 32.88
}
func FindFunctionByName (functionname string) Function {
	return Function{}
}
*/
/*
func GetFunctionPods (functionname string) ([]core.Pod,error) {
	params := make(map[string]string)
	params["namespace"] = "default"
	params["prefix"] = functionname
	bytes,err := clientutil.HttpGetWithPrefix("Pod",params)
	pods := make([]core.Pod)
	err = json.Unmarshal(bytes,&pods)
	if err != nil{
		return nil,err
	}
	return pods,nil
}
*/