package functioncontroller

import(
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	//"minik8s.io/pkg/controller"
	"context"
	"minik8s.io/pkg/util/listwatch"
	"sync"
	"time"
	//"minik8s.io/pkg/clientutil"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	"fmt"
)

var(
	FunctionStatus	 string = "/function/status"
	FunctionApplyUrl string = FunctionStatus + "/apply".
	FunctionDelUrl string = FunctionStatus + "/delete"
	FunctionTriggerUrl	string = FunctionStatus + "/trigger"
	defaultcountdown int = 10
)

type Function struct{
	//mock
	Metadata meta.ObjectMeta
	//Pods []core.Pod
}

type FunctionController struct{
	deploymentMap		map[string]core.Deployment	//record function name to function
	replicaMap			map[string]int
	countdownMap 		map[string]int
	requestMap			map[string]int
	mutex 				sync.Mutex
}

func NewFunctionController() (*FunctionController,error) {
	fc := &FunctionController{
		deploymentMap: 	make(map[string]core.Deployment),
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
	go listwatch.Watch(FunctionApplyUrl, fc.applylistener)
	go listwatch.Watch(FunctionTriggerUrl, fc.triggerlistener)
	go listwatch.Watch(FunctionDelUrl, fc.deletelistener)
}

func (fc *FunctionController) applylistener (msg *redis.Message) {
	print("fc listen apply\n")
	bytes := []byte(msg.Payload)
	deployment := core.Deployment{}
	err := json.Unmarshal(bytes, &deployment)
	if err != nil {
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.replicaMap[deployment.Metadata.Name] = 1
	fc.countdownMap[deployment.Metadata.Name] = defaultcountdown
	fc.requestMap[deployment.Metadata.Name] = 0
}

func (fc *FunctionController) triggerlistener (msg *redis.Message) {
	print("fc listen trigger\n")
	bytes := []byte(msg.Payload)
	deployment := core.Deployment{}
	err := json.Unmarshal(bytes, &deployment)
	if err != nil {
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.countdownMap[deployment.Metadata.Name] = defaultcountdown
	//fc.IncreaseReplica(deployment.Metadata.Name)
	fc.requestMap[deployment.Metadata.Name] += 1
	fmt.Println("request/s:",fc.requestMap[deployment.Metadata.Name])
	//fc.mutex.Unlock()

	if fc.requestMap[deployment.Metadata.Name] >= 3{
		fc.IncreaseReplica(deployment.Metadata.Name)
	}
	if fc.requestMap[deployment.Metadata.Name] >= 6{
		fc.IncreaseReplica(deployment.Metadata.Name)
	}
}

func (fc *FunctionController) deletelistener (msg *redis.Message) {
	print("fc listen trigger\n")
	bytes := []byte(msg.Payload)
	deployment := core.Deployment{}
	err := json.Unmarshal(bytes, &deployment)
	if err != nil {
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	delete(fc.deploymentMap, deployment.Metadata.Name)
	delete(fc.replicaMap, deployment.Metadata.Name)
	delete(fc.countdownMap, deployment.Metadata.Name)
	delete(fc.requestMap, deployment.Metadata.Name)
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

func (fc *FunctionController) ScaleTo0 (deploymentname string) error {
	fmt.Println("ScaleTo0",deploymentname)
	delete(fc.countdownMap, deploymentname)
	fc.replicaMap[deploymentname] = 0
	return nil
}

func (fc *FunctionController)IncreaseReplica(deploymentname string) {
	fmt.Println("increase replica",deploymentname)
	fc.replicaMap[deploymentname] += 1
}

func (fc *FunctionController)DecreaseReplica(deploymentname string) {
	fmt.Println("decrease replica",deploymentname)
	fc.replicaMap[deploymentname] -= 1
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