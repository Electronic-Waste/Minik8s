package deployment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	// "github.com/go-redis/redis/v8"

	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/util/listwatch"
)

// Return certain deployment's status
// uri: /api/v1/deployment/status/get?namespace=...&name=...
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
// uri: /api/v1/deployment/status/getall
func HandleGetAllDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.DeploymentStatus
	var deploymentStatusArr []string
	deploymentStatusArr, err := etcd.GetWithPrefix(etcdPrefix)
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

// Update a deployment's status in etcd
// uri: /api/v1/deployment/status/put?namespace=...&name=...
// @namespace: namespace requested; @name: deployment name
// body: core.Deployment in JSON form
func HandlePutDeploymentStatus(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	deploymentName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
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
	pubURL := path.Join(url.DeploymentStatus, "get", namespace, deploymentName)
	listwatch.Publish(pubURL, string(body))
	resp.WriteHeader(http.StatusOK)
}

// Delete a deployment's status in etcd
// uri: /api/v1/deployment/status/del?namespace=...&name=...
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
	listwatch.Publish(pubURL, "")
	resp.WriteHeader(http.StatusOK)
}
