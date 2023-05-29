package functioncontroller

import(
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	"minik8s.io/pkg/controller"
	"context"
	"minik8s.io/pkg/util/listwatch"
	"sync"
	"time"
	"minik8s.io/pkg/clientutil"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	"fmt"
)

var(
	FunctionStatus	 string = "/function/status"
	FunctionApplyUrl string = FunctionStatus + "/apply"
	FunctionTriggerUrl	string = FunctionStatus + "/trigger"
)

type Function struct{
	//mock
	Metadata meta.ObjectMeta
	//Pods []core.Pod
}

type FunctionController struct{
	deploymentMap		map[string]core.Deployment	//record function name to function
	countdownMap 		map[string]int
	mutex 				sync.Mutex
}

func NewFunctionController() (*FunctionController,error) {
	fc := &FunctionController{
		deploymentMap: 	make(map[string]core.Deployment),
		countdownMap:	make(map[string]int),
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
}

func (fc *FunctionController) applylistener (msg *redis.Message) {
	print("fc listen apply\n")
	bytes := []byte(msg.Payload)
	function := Function{}
	err := json.Unmarshal(bytes, &function)
	if err != nil {
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	//fc.replicaMap[function.Metadata.Name] = 1
	fc.countdownMap[function.Metadata.Name] = 10
}

func (fc *FunctionController) triggerlistener (msg *redis.Message) {
	print("fc listen trigger\n")
	bytes := []byte(msg.Payload)
	function := Function{}
	err := json.Unmarshal(bytes, &function)
	if err != nil {
		return
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.countdownMap[function.Metadata.Name] = 10
}

func (fc *FunctionController) countdown () {
	for{
		time.Sleep(time.Second)
		//fmt.Println("counting")
		fc.mutex.Lock()
		for f,c := range fc.countdownMap{
			fc.countdownMap[f] = c - 1
			fmt.Println("countdown",f, fc.countdownMap[f])
			//scale to 0, delete function
			if fc.countdownMap[f] == 0 {
				fmt.Println("scale to 0")
				err := fc.DeleteFunction(f)
				if err != nil{
					fmt.Println(err)
					continue
				}
			}
		}
		fc.mutex.Unlock()
	}
}

func (fc *FunctionController) DeleteFunction (functionname string) error {
	fmt.Println("deletefunction",functionname)
	delete(fc.countdownMap, functionname)
	return nil
	//funtcion := findFunctionByName(functionname)
	pods,_ := controller.GetReplicaPods(functionname)
	//pods := function.Pods
	for _,p := range pods{
		params := make(map[string]string)
		params["namespace"] = "default"
		params["name"] = p.Name
		err := clientutil.HttpDel("Pod", params)
		if err != nil{
			return err
		}
	}
	return nil
}

func (fc *FunctionController) scaler () {
	return
}

func FindFunctionByName (functionname string) Function {
	return Function{}
}
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