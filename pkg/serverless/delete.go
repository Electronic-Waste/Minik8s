package serverless

import (
	"fmt"
	"net/http"

	svlurl "minik8s.io/pkg/serverless/util/url"
	"minik8s.io/pkg/clientutil"
)

// Knative handles function deletion
// uri: /func/del?name=...
func (k *Knative) HandleFuncDel(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleFunDel receive msg!")
	vars := req.URL.Query()
	funcName := vars.Get("name")
	// Param miss: return error to client
	if funcName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	// 1. Delete corresponding core.Deployment
	targetDeploymentName := svlurl.DeploymentNamePrefix + funcName
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = targetDeploymentName
	err := clientutil.HttpDel("Deployment", params)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Fail to delete deployment"))
		return
	}

	// 2. Delete corresponding core.Function
	params["name"] = funcName
	err = clientutil.HttpDel("Knative-Function", params)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Fail to func"))
		return
	}

	// 3. Success
	resp.WriteHeader(http.StatusOK)
}