package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/kubelet/config"
	"net/http"
	"path"
	// "github.com/go-redis/redis/v8"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/util/listwatch"
)

// Return certain pod's status
// uri: /pods/status/get?namespace=...&name=...
// @namespace: namespace requested; @name: pod name
func HandleGetPodStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	podName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || podName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := path.Join(url.PodStatus, namespace, podName)
	PodStatus, err := etcd.Get(etcdKey)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(PodStatus))
	return
}

// Return all pods' statuses
// uri: /pods/status/getall
// now get all living pods on all nodes
func HandleGetAllPodStatus(resp http.ResponseWriter, req *http.Request) {
	fmt.Println()
	bytes, err := clientutil.HttpGet("nodes", map[string]string{})
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	nodeList := core.NodeList{}
	json.Unmarshal(bytes, &nodeList)
	fmt.Println("nodelist:", nodeList)

	var allpods []core.Pod
	for _,n := range nodeList.NodeArray{
		geturl := url.HttpScheme + n.Spec.NodeIp + config.Port + config.GetAllPodUrl
		fmt.Println(geturl)
		bytes, err := clientutil.HttpGetPlus("Pod", geturl)
		//fmt.Println("get from nodes:",string(bytes))
		pods := []core.Pod{}
		err = json.Unmarshal(bytes, &pods)
		if err != nil {
			fmt.Println(err)
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(err.Error()))
			return
		}
		//fmt.Println("get pods from",n.Spec.NodeIp,"get",pods)
		fmt.Println("get pods from", n.Spec.NodeIp)
		allpods = append(allpods, pods...)
	}
	data,err := json.Marshal(allpods)
	if err != nil{
		fmt.Println(err)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
	/*
		etcdPrefix := url.PodStatus
		var podStatusArr []string
		podStatusArr, err := etcd.GetWithPrefix(etcdPrefix)
		// Error occur in etcd: return error to client
		if err != nil {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(err.Error()))
			return
		}
		var jsonVal []byte
		jsonVal, err = json.Marshal(podStatusArr)
		// Error occur in json parsing: return error to client
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte(err.Error()))
			return
		}
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(jsonVal)
		// return
	*/
}

// Return statuses of pods with prefix
// uri: /pods/status/getwithprefix
func HandleGetWithPrefixPodStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	prefix := vars.Get("prefix")
	etcdPrefix := path.Join(url.PodStatus, namespace, prefix)
	var podStatusArr []string
	podStatusArr, err := etcd.GetWithPrefix(etcdPrefix)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	var jsonVal []byte
	jsonVal, err = json.Marshal(podStatusArr)
	// Error occur in json parsing: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonVal)
	// return
}

// Apply a pod's status in etcd
// uri: /pods/status/apply?namespace=...&name=...
// @namespace: namespace requested; @name: pod name
// body: core.Pod in JSON form
func HandleApplyPodStatus(resp http.ResponseWriter, req *http.Request) {
	//vars := req.URL.Query()
	//namespace := vars.Get("namespace")
	//podName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)

	pod := core.Pod{}
	json.Unmarshal(body, &pod)
	podName := pod.Name
	namespace := "default"
	// Param miss: return error to client
	if namespace == "" || podName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.PodStatus, namespace, podName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.PodStatus, "apply")
	Param := core.ScheduleParam{
		RunPod: pod,
	}
	// get all node registered message
	nodeStr, err := etcd.GetWithPrefix(url.Node)
	for _, str := range nodeStr {
		node := core.Node{}
		json.Unmarshal([]byte(str), &node)
		Param.NodeList = append(Param.NodeList, node)
	}
	fmt.Println(Param)
	body, err = json.Marshal(Param)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	listwatch.Publish(pubURL, string(body))
	resp.WriteHeader(http.StatusOK)
}

// Update a pod's status in etcd
// uri: /pods/status/update?namespace=...&name=...
// @namespace: namespace requested; @name: pod name
// body: core.Pod in JSON form
func HandleUpdatePodStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	podName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
	// Param miss: return error to client
	if namespace == "" || podName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.PodStatus, namespace, podName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.PodStatus, "update")
	listwatch.Publish(pubURL, string(body))
	resp.WriteHeader(http.StatusOK)
}

// Delete a pod's status in etcd
// uri: /pods/status/del?namespace=...&name=...
// @namespace: namespace requested; @name: pod name
func HandleDelPodStatus(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("http del")
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	podName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || podName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.PodStatus, namespace, podName)
	data, err := etcd.Get(etcdURL)
	if err != nil {
		fmt.Println("HandleDelPodStatus etcd get fail")
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	pod := core.Pod{}
	json.Unmarshal([]byte(data), &pod)
	err = etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		fmt.Println("HandleDelPodStatus etcd delete fail")
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	fmt.Printf("del pod name is %s\n", pod.Name)
	clientutil.HttpPlus("Pod", pod, url.HttpScheme+pod.Spec.RunningNode.Spec.MasterIp+config.Port+config.DelPodRul)

	bytes, err := json.Marshal(podName)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	pubURL := url.PodStatusDelURL
	listwatch.Publish(pubURL, bytes)

	resp.WriteHeader(http.StatusOK)
	fmt.Println("http del success")
}

func HandleGetPodMetrics(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("handle get pod metrics")
	vars := req.URL.Query()
	podName := vars.Get("name")
	nodeip := vars.Get("nodeip")

	geturl := url.HttpScheme + nodeip + config.Port + config.PodMetricsUrl + "?name=" + podName

	bytes, err := clientutil.HttpGetPlus("Metrics", geturl)
	if err != nil {
		fmt.Println(err)
	}

	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(bytes)
}
