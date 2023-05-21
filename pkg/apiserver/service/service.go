package service

import (
	"net/http"
	"encoding/json"
	"path"
	"io/ioutil"

	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
)

// Return certain service's status
// uri: /service/get?namespace=...&name=...
// @namespace: namespace requested; @name: service name
func HandleGetService(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	serviceName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || serviceName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := path.Join(url.Service, namespace, serviceName)
	Service, err := etcd.Get(etcdKey)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(Service))
	return
}

// Return all services' statuses
// uri: /service/getall
func HandleGetAllServices(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.Service
	var serviceArr []string
	serviceArr, err := etcd.GetWithPrefix(etcdPrefix)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	var jsonVal []byte
	jsonVal, err = json.Marshal(serviceArr)
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


// Apply a service's status in etcd
// uri: /service/apply?namespace=...&name=...
// @namespace: namespace requested; @name: service name
// body: core.Service in JSON form
func HandleApplyService(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	serviceName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
	// Param miss: return error to client
	if namespace == "" || serviceName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.Service, namespace, serviceName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	// TODO(shaowang): Implement selector & Design redis message pattern
	pubURL := path.Join(url.Service, "apply", namespace, serviceName)
	listwatch.Publish(pubURL, string(body))	
	resp.WriteHeader(http.StatusOK)
}

// Update a service's status in etcd
// uri: /service/update?namespace=...&name=...
// @namespace: namespace requested; @name: service name
// body: core.Service in JSON form
func HandleUpdateService(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	serviceName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
	// Param miss: return error to client
	if namespace == "" || serviceName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.Service, namespace, serviceName)
	err := etcd.Put(etcdURL, string(body))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.Service, "update", namespace, serviceName)
	listwatch.Publish(pubURL, string(body))	
	resp.WriteHeader(http.StatusOK)
}

// Delete a service's status in etcd
// uri: /service/del?namespace=...&name=...
// @namespace: namespace requested; @name: service name
func HandleDelService(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	serviceName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || serviceName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.Service, namespace, serviceName)
	err := etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success!
	pubURL := path.Join(url.Service, "del", namespace, serviceName)
	listwatch.Publish(pubURL, "")	
	resp.WriteHeader(http.StatusOK)
}