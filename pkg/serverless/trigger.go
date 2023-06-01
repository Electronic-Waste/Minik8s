package serverless

import(
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	svlurl "minik8s.io/pkg/serverless/util/url"
	//apiurl "minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/util/listwatch"
)

// Knative handle function trigger
// uri: /func/trigger?name=...
// body: params in JSON form
func (k *Knative) HandleFuncTrigger(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleFuncTrigger receive msg!")
	vars :=  req.URL.Query()
	funcName := vars.Get("name")
	params, _ := ioutil.ReadAll(req.Body)
	// Param miss: return error to client
	if funcName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	result, err := k.TriggerFunction(funcName, params)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(result))
}

func (k *Knative) TriggerFunction(funcName string, params []byte) (string, error) {
	//inform functioncontroller in the first place
	redismsg,err := json.Marshal(funcName)
	if err != nil{
		fmt.Println(err)
		return "",err
	}
	listwatch.Publish(svlurl.FunctionTriggerURL, redismsg)
	// 1. Query corresponding pod
	podNamePrefix := svlurl.DeploymentNamePrefix + funcName
	podParams := make(map[string]string)
	podParams["namespace"] = "default"
	podParams["prefix"] = podNamePrefix
	content, err := clientutil.HttpGetWithPrefix("Pod", podParams)
	if err != nil {
		return "", err
	}
	var podStrings []string
	err = json.Unmarshal(content, &podStrings)
	if err != nil {
		return "", err
	}
	var pods []core.Pod
	for _, podString := range podStrings {
		var pod core.Pod
		json.Unmarshal([]byte(podString), &pod)
		pods = append(pods, pod)
	}

	// 2. Judge whether has scalee-to-0 or not. If so, polling the apiserver
	if (len(pods) == 0) {
		// Publish to redis to inform watcher of missing pod
		// - topic: /func/trigger
		// - payload: function's name
		redismsg,err := json.Marshal(funcName)
		if err != nil{
			fmt.Println(err)
			return "",err
		}
		listwatch.Publish(svlurl.FunctionTriggerURL, redismsg)
		
		// Polling apiserver in 3s
		for triggerCount := 0; len(pods) == 0; triggerCount++ {
			fmt.Println("trigger time: ", triggerCount)
			time.Sleep(3 * time.Second)
			content, err := clientutil.HttpGetWithPrefix("Pod", podParams)
			if err != nil {
				return "", err
			}
			fmt.Printf("content: %s", string(content))
			err = json.Unmarshal(content, &podStrings)
			if err != nil {
				return "", err
			}
			for _, podString := range podStrings {
				var pod core.Pod
				json.Unmarshal([]byte(podString), &pod)
				pods = append(pods, pod)
			}
		}
	}

	// 3. Choose the serving pod with round-robin policy & Send request
	targetPod := pods[k.rrCount % len(pods)]
	k.rrCount++
	targetPodIP := strings.Replace(targetPod.Status.PodIp, "\"", "", -1)
	triggerURL := svlurl.HttpScheme + targetPodIP + ":8080"
	result, err := clientutil.HttpTrigger("Knative-Function", triggerURL, params)
	if err != nil {
		return "", err
	}
	return result, nil
}