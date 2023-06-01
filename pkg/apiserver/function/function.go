package function

import (
	"path"
	"encoding/json"
	"net/http"
	"io/ioutil"

	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
)

// Register a function in etcd
// uri: /func/register
// body: core.Function in JSON form
func HandleRegisterFunction(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	function := core.Function{}
	json.Unmarshal(body, &function)
	functionName := function.Name
	// Params miss: return error to serverless manager
	if functionName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdURL := path.Join(url.Function, functionName)
	err := etcd.Put(etcdURL, string(body))
	// Error occurred in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// Success
	resp.WriteHeader(http.StatusOK)
}



// Invoked when there is no pod for func
// uri: /func/trigger
// body: core.Function in JSON form
// func HandleTriggerFunction(resp http.ResponseWriter, req *http.Request) {

// }

// Get all functions
// uri: /func/getall
func HandleGetAllFunction(resp http.ResponseWriter, req *http.Request) {
	
}

// Get a function with function name
// uri: /func/get?name=xxx
func HandleGetFunction(resp http.ResponseWriter, req *http.Request) {

}

// Update a function
// uri: /func/update
// body: core.Function in JSON form
func HandleUpdateFunction(resp http.ResponseWriter, req *http.Request) {

}

// Delete a function with function name
// uri: /func/del
func HandleDelFunction(resp http.ResponseWriter, req *http.Request) {
	
}