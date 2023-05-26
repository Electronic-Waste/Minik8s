package apiserver

import (
	"net/http"

	"minik8s.io/pkg/apiserver/deployment"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/pod"
	"minik8s.io/pkg/apiserver/service"
	"minik8s.io/pkg/apiserver/dns"
	"minik8s.io/pkg/apiserver/util/url"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var postHandlerMap = map[string]HttpHandler{
	url.PodStatusApplyURL:         	pod.HandleApplyPodStatus,
	url.PodStatusUpdateURL:        	pod.HandleUpdatePodStatus,
	url.DeploymentStatusApplyURL:  	deployment.HandleApplyDeploymentStatus,
	url.DeploymentStatusUpdateURL: 	deployment.HandleUpdateDeploymentStatus,
	url.ServiceApplyURL:		   	service.HandleApplyService,
	url.ServiceUpdateURL:			service.HandleUpdateService,
	url.DNSApplyURL:				dns.HandleApplyDNS,
}

var getHandlerMap = map[string]HttpHandler{
	url.PodStatusGetURL:           	pod.HandleGetPodStatus,
	url.PodStatusGetAllURL:        	pod.HandleGetAllPodStatus,
	url.DeploymentStatusGetURL:    	deployment.HandleGetDeploymentStatus,
	url.DeploymentStatusGetAllURL: 	deployment.HandleGetAllDeploymentStatus,
	url.ServiceGetURL:				service.HandleGetService,
	url.ServiceGetAllURL:			service.HandleGetAllServices,
	url.DNSGetURL:					dns.HandleGetDNS,
	url.DNSGetAllURL:				dns.HandleGetAllDNS,
}

var deleteHandlerMap = map[string]HttpHandler{
	url.PodStatusDelURL:        pod.HandleDelPodStatus,
	url.DeploymentStatusDelURL: deployment.HandleDelDeploymentStatus,
	url.ServiceDelURL:			service.HandleDelService,
	url.DNSDelURL:				dns.HandleDelDNS,
}

func bindWatchHandler() {

}

func Run() {
	// Initialize etcd client
	etcd.InitializeEtcdKVStore()
	
	// Init clusterIPGen in service module
	service.InitServiceController()

	// Bind POST request with handler
	for url, handler := range postHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Bind GET request with handler
	for url, handler := range getHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Bind DELETE request with handler
	for url, handler := range deleteHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Bind watch message with WatchHandler
	bindWatchHandler()

	// Start Server
	http.ListenAndServe(url.Port, nil)
}
