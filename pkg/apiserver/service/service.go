package service

import (
	"fmt"
	"net/http"
	"encoding/json"
	"path"
	"strings"
	"io/ioutil"
	"os/exec"

	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/util/ipgen"
)

var clusterIPGen *ipgen.ClusterIPGenerator


func InitServiceController() {
	gen, err := ipgen.NewClusterIPGenerator()
	if err != nil {
		fmt.Printf("Error occurred in init cluster IP generator: %v", err)
	}
	clusterIPGen = gen
}

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
	fmt.Printf("HandleApplyService receive msg: namespace is %s, serviceName is %s, body is %s\n",
					namespace, serviceName, string(body))
	// Param miss: return error to client
	if namespace == "" || serviceName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	// 1. Allocate new clusterIP for service & Upate serviceSpec & Write serviceSpec to etcd
	var serviceSpec core.Service
	clusterIP, err := clusterIPGen.NextClusterIP()	// New clusterIP
	fmt.Printf("New clusterIP is %s\n", clusterIP)
	json.Unmarshal(body, &serviceSpec)
	serviceSpec.Spec.ClusterIP = clusterIP
	fmt.Printf("serviceSpec is: %v", serviceSpec)
	serviceJsonVal, _ := json.Marshal(serviceSpec)
	etcdURL := path.Join(url.Service, namespace, serviceName)
	err = etcd.Put(etcdURL, string(serviceJsonVal))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// 2. Match with selector
	appName := serviceSpec.Spec.Selector["app"]
	fmt.Printf("appName: %s\n", appName)
	var podStrings []string
	var podStatuses []core.Pod
	var podNames []string
	var podIPs []string
	podStrings, err = etcd.GetWithPrefix(url.PodStatus)
	// Error occured in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	for _, podString := range podStrings {
		// TODO(shaowang): Address empty string problem
		if podString == "" {
			continue
		}
		fmt.Printf("podString: %s\n", podString)
		var tmpPod core.Pod
		json.Unmarshal([]byte(podString), &tmpPod)
		podStatuses = append(podStatuses, tmpPod)
	}
	for _, podStatus := range podStatuses {
		if podStatus.Labels["app"] != appName {
			continue
		}
		podName := podStatus.Name
		podNames = append(podNames, podName)
		cmd := exec.Command("nerdctl", "inspect", "-f", "{{.NetworkSettings.IPAddress}}", podName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Could not find IP!")
		}
		podIP := strings.Replace(string(output), "\n", "", -1)
		fmt.Printf("podIP: %s", podIP)
		podIPs = append(podIPs, podIP)
	}
	// 3. Construct KubeproxyServiceParam & Publish
	serviceParam := core.KubeproxyServiceParam{
		ServiceName		: serviceSpec.Name,
		ClusterIP		: clusterIP,
		ServicePorts	: serviceSpec.Spec.Ports,
		PodNames		: podNames,
		PodIPs			: podIPs,
	}
	serviceParamJsonVal, _ := json.Marshal(serviceParam)
	fmt.Printf("Params to redis: %s\n", string(serviceParamJsonVal))
	// Publish to redis:
	// - topic: /service/apply
	// - payload: <KubeproxyServiceParam>
	listwatch.Publish(url.ServiceApplyURL, string(serviceParamJsonVal))	
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
	// Publish to redis:
	// 1. topic: /service/update
	// 2. payload: core.Service's JSON format
	pubURL := path.Join(url.Service, "update")
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
	fmt.Printf("HandleDelService receive msg: namespace is %s, serviceName is %s\n", namespace, serviceName)
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
	// Publish to redis:
	// 1. topic: /service/del
	// 2. payload: <serviceName>
	listwatch.Publish(url.ServiceDelURL, serviceName)
	resp.WriteHeader(http.StatusOK)
}