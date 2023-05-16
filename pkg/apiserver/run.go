package apiserver

import(
	"net/http"

	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/apiserver/pod"
	"minik8s.io/pkg/apiserver/etcd"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var postHandlerMap = map[string]HttpHandler{
	url.PodStatusPutURL : pod.HandlePutPodStatus,
	
}

var getHandlerMap = map[string]HttpHandler{
	url.PodStatusGetURL : pod.HandleGetPodStatus,
	url.PodStatusGetAllURL : pod.HandleGetAllPodStatus,
}

var deleteHandlerMap = map[string]HttpHandler{
	url.PodStatusDelURL : pod.HandleDelPodStatus,
}

func bindWatchHandler() {

}

func Run() {
	// Initialize etcd client
	etcd.InitializeEtcdKVStore()

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