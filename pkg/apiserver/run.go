package apiserver

import(
	"net/http"
	"vmeet.io/minik8s/pkg/apiserver/util/url"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var postHandlerMap = map[string]HttpHandler{
	
}

var getHandlerMap = map[string]HttpHandler{

}

var deleteHandlerMap = map[string]HttpHandler{

}

func bindWatchHandler() {

}

func Run() {
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