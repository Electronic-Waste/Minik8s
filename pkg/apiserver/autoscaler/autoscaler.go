package autoscaler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"fmt"
	// "github.com/go-redis/redis/v8"

	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apis/core"
)

// Return certain autoscaler's status
// uri: /autoscalers/status/get?namespace=...&name=...
// @namespace: namespace requested; @name: autoscaler name
func HandleGetAutoscalerStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	autoscalerName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || autoscalerName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := path.Join(url.AutoscalerStatus, namespace, autoscalerName)
	AutoscalerStatus, err := etcd.Get(etcdKey)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(AutoscalerStatus))
	return
}

// Return all autoscalers' statuses
// uri: /autoscalers/status/getall
func HandleGetAllAutoscalerStatus(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.AutoscalerStatus
	var autoscalerStatusArr []string
	autoscalerStatusArr, err := etcd.GetWithPrefix(etcdPrefix)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	var jsonVal []byte
	jsonVal, err = json.Marshal(autoscalerStatusArr)
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

// Apply a autoscaler's status in etcd
// uri: /autoscalers/status/apply?namespace=...&name=...
// @namespace: namespace requested; @name: autoscaler name
// body: core.Autoscaler in JSON form
func HandleApplyAutoscalerStatus(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("receive http apply autoscaler request")
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	//autoscalerName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)

	autoscaler := core.Autoscaler{}
	json.Unmarshal(body, &autoscaler)
	autoscalerName := autoscaler.Metadata.Name
	//namespace := "default"

	// Param miss: return error to client
	if namespace == "" || autoscalerName == "" {
		fmt.Println("autoscalerName or namespace is missing")
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.AutoscalerStatus, namespace, autoscalerName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	fmt.Println("etcd apply autoscaler successfully")
	pubURL := path.Join(url.AutoscalerStatus, "apply")
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "apply"
	watchres.ObjectType = "Autoscaler"
	watchres.Payload = body

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}

// Update a autoscaler's status in etcd
// uri: /autoscalers/status/update?namespace=...&name=...
// @namespace: namespace requested; @name: autoscaler name
// body: core.Autoscaler in JSON form
func HandleUpdateAutoscalerStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	autoscalerName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
	// Param miss: return error to client
	if namespace == "" || autoscalerName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.AutoscalerStatus, namespace, autoscalerName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.AutoscalerStatus, "update")
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "update"
	watchres.ObjectType = "Autoscaler"
	watchres.Payload = body

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}

// Delete a autoscaler's status in etcd
// uri: /autoscalers/status/del?namespace=...&name=...
// @namespace: namespace requested; @name: autoscaler name
func HandleDelAutoscalerStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	autoscalerName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || autoscalerName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.AutoscalerStatus, namespace, autoscalerName)
	err := etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.AutoscalerStatus, "del", namespace, autoscalerName)
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "delete"
	watchres.ObjectType = "Autoscaler"
	watchres.Payload, _ = json.Marshal(autoscalerName)

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}
