package pod

import (
	"net/http"
	"encoding/json"
	"path"
	"io/ioutil"

	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
)

// Return certain pod's status
// uri: /api/v1/pod/status/get?namespace=...&name=...
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
	etcdKey := path.Join(url.PodStatusGetURL, namespace, podName)
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
// uri: /api/v1/pod/status/getall
func HandleGetAllPodStatus(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.PodURL
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

// Update a pod's status in etcd
// uri: /api/v1/pod/status/put?namespace=...&name=...
// @namespace: namespace requested; @name: pod name
// body: core.Pod in JSON form
func HandlePutPodStatus(resp http.ResponseWriter, req *http.Request) {
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
	etcdURL := path.Join(url.PodStatusPutURL, namespace, podName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	listwatch.Publish(etcdURL, string(body))	// TODO(shaowang): Normalize listwatch msg
	resp.WriteHeader(http.StatusOK)
}

// Delete a pod's status in etcd
// uri: /api/v1/pod/status/del?namespace=...&name=...
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
	etcdURL := path.Join(url.PodStatusPutURL, namespace, podName)
	err := etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	listwatch.Publish(etcdURL, "type: del")	// TODO(shaowang): Normalize listwatch msg
	resp.WriteHeader(http.StatusOK)
}

// func HandleWatchPodStatus(msg *redis.Message) {
	
// }