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

// Knative handle workflow trigger
// uri: /workflow/trigger
// body: core.Workflow in JSON form
func (k *Knative) HandleWorkflowTrigger(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleWorkflowTrigger receive msg!")
	body, _ := ioutil.ReadAll(req.Body)

	// 1. Parse body to core.Workflow
	workflow := core.Workflow{}
	json.Unmarshal(body, &workflow)
	workflowNodes := workflow.Nodes
	startFuncName := workflow.StartAt
	params := workflow.Params
	if len(workflowNodes) == 0 || startFuncName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Invalid workflow!"))
		return
	}

	// 2. Execute with control workflow DAG
	var result string
	var resultParams map[string]int
	triggerFuncName := startFuncName
	triggerParams, err := json.Marshal(params)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	for {
		workflowNode := workflowNodes[triggerFuncName]
		if workflowNode.Type == "Task" {
			result, err = k.TriggerFunction(triggerFuncName, triggerParams)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(err.Error()))
				return
			}
			fmt.Printf("Workflow get result: %s\n", result)
			// If Next does not exists, stop workflow and return result
			// Else update triggerFuncName and triggerParams, trigger next function
			if (workflowNode.Next == "") {
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(result))
				return
			} else {
				triggerFuncName = workflowNode.Next
				triggerParams = []byte(result)
			}
		} else if workflowNode.Type == "Choice" {
			hasOneChoiceMatch := false
			json.Unmarshal(triggerParams, &resultParams)
			for _, workflowChoice := range workflowNode.Choices {
				isMatch := true
				for key, val := range resultParams {
					expectedVal, ok := workflowChoice.Conditions[key]
					// If do not have key of val != expectedVal, the match ends
					if !ok || val != expectedVal {
						isMatch = false
						break
					}
				}
				if isMatch {
					triggerFuncName = workflowChoice.Next
					hasOneChoiceMatch = true
					break
				}
 			}
			if !hasOneChoiceMatch {
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(fmt.Sprint("Workflow branch error: could not find target branch")))
				return
			}
		}
	}
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
	time.Sleep(time.Second * 20)

	// 3. Choose the serving pod with round-robin policy & Send request
	fmt.Println("trigger: pod len:",len(pods),"and rrcount:",k.rrCount)
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