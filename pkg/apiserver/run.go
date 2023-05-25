package apiserver

import (
	"net/http"

	"minik8s.io/pkg/apiserver/deployment"
	"minik8s.io/pkg/apiserver/autoscaler"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/pod"
	"minik8s.io/pkg/apiserver/service"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/kubeproxy"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var postHandlerMap = map[string]HttpHandler{
	url.PodStatusApplyURL:         	pod.HandleApplyPodStatus,
	url.PodStatusUpdateURL:        	pod.HandleUpdatePodStatus,
	url.DeploymentStatusApplyURL:  	deployment.HandleApplyDeploymentStatus,
	url.DeploymentStatusUpdateURL: 	deployment.HandleUpdateDeploymentStatus,
	url.ServiceApplyURL:		   	service.HandleApplyService,
	url.ServiceUpdateURL:			service.HandleUpdateService,
	url.AutoscalerStatusApplyURL:  autoscaler.HandleApplyAutoscalerStatus,
	url.AutoscalerStatusUpdateURL: autoscaler.HandleUpdateAutoscalerStatus,
}

var getHandlerMap = map[string]HttpHandler{
	url.PodStatusGetURL:           	pod.HandleGetPodStatus,
	url.PodStatusGetAllURL:        	pod.HandleGetAllPodStatus,
	url.PodStatusGetWithPrefixURL:      pod.HandleGetWithPrefixPodStatus,
	url.DeploymentStatusGetURL:    	deployment.HandleGetDeploymentStatus,
	url.DeploymentStatusGetAllURL: 	deployment.HandleGetAllDeploymentStatus,
	url.ServiceGetURL:				service.HandleGetService,
	url.ServiceGetAllURL:			service.HandleGetAllServices,
	url.AutoscalerStatusGetURL:    		autoscaler.HandleGetAutoscalerStatus,
	url.AutoscalerStatusGetAllURL: 		autoscaler.HandleGetAllAutoscalerStatus,
}

var deleteHandlerMap = map[string]HttpHandler{
	url.PodStatusDelURL:        pod.HandleDelPodStatus,
	url.DeploymentStatusDelURL: deployment.HandleDelDeploymentStatus,
	url.ServiceDelURL:			service.HandleDelService,
	url.AutoscalerStatusDelURL: autoscaler.HandleDelAutoscalerStatus,
}

func bindWatchHandler() {

}

var kubeProxyManager *kubeproxy.KubeproxyManager

func Run() {
	// Initialize etcd client
	etcd.InitializeEtcdKVStore()
	
	// Init clusterIPGen in service
	service.InitServiceController()

	kubeProxyManager, _ = kubeproxy.NewKubeProxy()
	kubeProxyManager.Run()

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
