package pod

import (
	"net/http"
	"github.com/go-redis/redis/v8"
	"encoding/json"

	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/etcdkeyprefix"
	"minik8s.io/pkg/apis/core"
)

func HandleGetPodStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	podName := vars.Get("name")
	// Param "name" miss: return error to client
	if podName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := etcdkeyprefix.POD + podName
	var pod core.Pod
	interface_value, err := etcd.Get(etcdKey, pod)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	pod = interface_value.(core.Pod)
	jsonVal, _ := json.Marshal(pod)
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(jsonVal))
	return
}

func HandleGetPodStatuses(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := etcdkeyprefix.POD
	var pods []core.Pod
	interface_values, err := etcd.GetWithPrefix(etcdPrefix, pods)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	
	// pods = interface_value.([]core.Pod)
	// jsonVal, _ := json.Marshal(pods)
	// resp.WriteHeader(http.StatusOK)
	// resp.Header().Set("Content-Type", "application/json")
	// resp.Write([]byte(jsonVal))
	// return
}

func HandleSetPodStatus(resp http.ResponseWriter, req *http.Request) {

}

func HandleDelPodStatus(resp http.ResponseWriter, req *http.Request) {

}

func HandleWatchPodStatus(msg *redis.Message) {
	
}