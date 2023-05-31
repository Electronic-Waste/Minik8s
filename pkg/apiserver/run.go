package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/job"
	"minik8s.io/pkg/apiserver/node"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/kubelet/config"
	"minik8s.io/pkg/util/listwatch"
	"net/http"
	"path"

	"minik8s.io/pkg/apiserver/autoscaler"
	"minik8s.io/pkg/apiserver/deployment"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/pod"
	"minik8s.io/pkg/apiserver/service"
	"minik8s.io/pkg/apiserver/dns"
	"minik8s.io/pkg/apiserver/util/url"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var postHandlerMap = map[string]HttpHandler{
	url.PodStatusApplyURL:         pod.HandleApplyPodStatus,
	url.PodStatusUpdateURL:        pod.HandleUpdatePodStatus,
	url.DeploymentStatusApplyURL:  deployment.HandleApplyDeploymentStatus,
	url.DeploymentStatusUpdateURL: deployment.HandleUpdateDeploymentStatus,
	url.AutoscalerStatusApplyURL:  autoscaler.HandleApplyAutoscalerStatus,
	url.AutoscalerStatusUpdateURL: autoscaler.HandleUpdateAutoscalerStatus,
	url.ServiceApplyURL:           service.HandleApplyService,
	url.ServiceUpdateURL:          service.HandleUpdateService,
	url.JobApplyUrl:               job.HandleApplyJob,
	url.JobMapUrl:                 job.HandleMapJob,
	url.DNSApplyURL:				dns.HandleApplyDNS,
}

var getHandlerMap = map[string]HttpHandler{
	url.PodStatusGetURL:           pod.HandleGetPodStatus,
	url.PodStatusGetAllURL:        pod.HandleGetAllPodStatus,
	url.PodStatusGetWithPrefixURL: 	pod.HandleGetWithPrefixPodStatus,
	url.DeploymentStatusGetURL:    deployment.HandleGetDeploymentStatus,
	url.DeploymentStatusGetAllURL: deployment.HandleGetAllDeploymentStatus,
	url.AutoscalerStatusGetURL:    autoscaler.HandleGetAutoscalerStatus,
	url.AutoscalerStatusGetAllURL: autoscaler.HandleGetAllAutoscalerStatus,
	url.ServiceGetURL:             service.HandleGetService,
	url.ServiceGetAllURL:          service.HandleGetAllServices,
	url.NodesGetUrl:               node.HandleGetNodes,
	url.JobGetUrl:                 job.HandleGetJob,
	url.DNSGetURL:					dns.HandleGetDNS,
	url.DNSGetAllURL:				dns.HandleGetAllDNS,
	url.MetricsGetUrl:				pod.HandleGetPodMetrics,
}

var deleteHandlerMap = map[string]HttpHandler{
	url.PodStatusDelURL:        pod.HandleDelPodStatus,
	url.DeploymentStatusDelURL: deployment.HandleDelDeploymentStatus,
	url.AutoscalerStatusDelURL: autoscaler.HandleDelAutoscalerStatus,
	url.ServiceDelURL:          service.HandleDelService,
	url.DNSDelURL:				dns.HandleDelDNS,
}

var nodeHandlerMap = map[string]HttpHandler{
	url.NodeRergisterUrl: node.HandleNodeRegister,
}

func bindWatchHandler() {
	go listwatch.Watch(url.SchedApplyURL, HitNode)
}

func HitNode(msg *redis.Message) {
	fmt.Printf("call HitNode\n")
	pod := core.Pod{}
	json.Unmarshal([]byte(msg.Payload), &pod)
	//fmt.Println(pod)
	podName := pod.Name
	namespace := "default"
	etcdURL := path.Join(url.PodStatus, namespace, podName)
	body, err := json.Marshal(pod)
	if err != nil {
		fmt.Println(err)
	}
	// inform core kubelet to apply the Pod
	fmt.Println("hit node: ",url.HttpScheme+pod.Spec.RunningNode.Spec.NodeIp+config.Port+config.RunPodUrl)
	err, str := clientutil.HttpPlus("Pod", pod, url.HttpScheme+pod.Spec.RunningNode.Spec.NodeIp+config.Port+config.RunPodUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("receive %s\n", str)
	pod.Status.PodIp = str
	body, err = json.Marshal(pod)
	if err != nil {
		fmt.Println(err)
	}
	err = etcd.Put(etcdURL, string(body))
	if err != nil {
		fmt.Println(err)
	}
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
	// Bind Node Method
	for url, handler := range nodeHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Bind watch message with WatchHandler
	bindWatchHandler()

	// Start Server
	http.ListenAndServe(url.Port, nil)
}
