package serverless

import (
	"fmt"
	"path"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	svlurl "minik8s.io/pkg/serverless/util/url"
)

// Knative handles function update
// uri: /func/update
// body: core.Function in JSON form
func (k *Knative) HandleFuncUpdate(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleFuncUpdate receive msg!")
	body, _ := ioutil.ReadAll(req.Body)

	function := core.Function{}
	json.Unmarshal(body, &function)
	funcName := function.Name
	funcFilePath := function.Path
	// Param miss: return error to client
	if funcName == "" || funcFilePath == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("funcName and filePath should not be empty string"))
		return
	}
	// 1. Update core.Function to apiserver /func/register
	err := clientutil.HttpApply("Function", function)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	// 2. Copy new file content to mounting file
	hostDestDir := path.Join(svlurl.HostMountDestPathPrefix, funcName)
	hostDestFuncPath := path.Join(hostDestDir, svlurl.FuncFile)
	err = CopyFile(hostDestFuncPath, funcFilePath)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to copy func file: %v", err)))
		return
	}

	// 3. Success
	resp.WriteHeader(http.StatusOK)
}