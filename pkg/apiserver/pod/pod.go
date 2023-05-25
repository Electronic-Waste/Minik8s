package pod

import (
	"encoding/json"
	"io/ioutil"
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
func HandleGetAllPodStatus(resp http.ResponseWriter, req *http.Request) {
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
	nodeStr, err := etcd.GetWithPrefix(url.NodeRergisterUrl)
	for _, str := range nodeStr {
		node := core.Node{}
		json.Unmarshal([]byte(str), &node)
		Param.NodeList = append(Param.NodeList, node)
	}
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
	err := etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.PodStatus, "del")
	listwatch.Publish(pubURL, podName)
	resp.WriteHeader(http.StatusOK)
}
