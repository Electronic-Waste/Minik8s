package deployment

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

// Return certain deployment's status
// uri: /deployments/status/get?namespace=...&name=...
// @namespace: namespace requested; @name: deployment name
func HandleGetDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	deploymentName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || deploymentName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := path.Join(url.DeploymentStatus, namespace, deploymentName)
	DeploymentStatus, err := etcd.Get(etcdKey)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(DeploymentStatus))
	return
}

// Return all deployments' statuses
// uri: /deployments/status/getall
func HandleGetAllDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.DeploymentStatus
	var deploymentStatusArr []string
	deploymentStatusArr, err := etcd.GetWithPrefix(etcdPrefix)
	//err := etcd.DelWithPrefix(etcdPrefix)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	var jsonVal []byte
	jsonVal, err = json.Marshal(deploymentStatusArr)
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

// Apply a deployment's status in etcd
// uri: /deployments/status/apply?namespace=...&name=...
// @namespace: namespace requested; @name: deployment name
// body: core.Deployment in JSON form
func HandleApplyDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("receive http apply deployment request")
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	//deploymentName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)

	deployment := core.Deployment{}
	json.Unmarshal(body, &deployment)
	deploymentName := deployment.Metadata.Name
	//namespace := "default"

	// Param miss: return error to client
	if namespace == "" || deploymentName == "" {
		fmt.Println("deploymentName or namespace is missing")
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.DeploymentStatus, namespace, deploymentName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	fmt.Println("etcd apply deployment successfully")
	pubURL := path.Join(url.DeploymentStatus, "apply")
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "apply"
	watchres.ObjectType = "Deployment"
	watchres.Payload = body

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}

// Update a deployment's status in etcd
// uri: /deployments/status/update?namespace=...&name=...
// @namespace: namespace requested; @name: deployment name
// body: core.Deployment in JSON form
func HandleUpdateDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	//deploymentName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)

	deployment := core.Deployment{}
	json.Unmarshal(body, &deployment)
	deploymentName := deployment.Metadata.Name
	
	// Param miss: return error to client
	if namespace == "" || deploymentName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.DeploymentStatus, namespace, deploymentName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.DeploymentStatus, "update")
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "update"
	watchres.ObjectType = "Deployment"
	watchres.Payload = body

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}

// Delete a deployment's status in etcd
// uri: /deployments/status/del?namespace=...&name=...
// @namespace: namespace requested; @name: deployment name
func HandleDelDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	deploymentName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || deploymentName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.DeploymentStatus, namespace, deploymentName)
	err := etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.DeploymentStatus, "del", namespace, deploymentName)
	watchres := listwatch.WatchResult{}
	watchres.ActionType = "delete"
	watchres.ObjectType = "Deployment"
	watchres.Payload, _ = json.Marshal(deploymentName)

	bytes,_ := json.Marshal(watchres)
	listwatch.Publish(pubURL, bytes)
	resp.WriteHeader(http.StatusOK)
}
