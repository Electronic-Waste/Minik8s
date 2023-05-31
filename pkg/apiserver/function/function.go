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