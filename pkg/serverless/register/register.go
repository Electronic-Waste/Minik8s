package register

import(
	"net/http"
	"io/ioutil"
	"encoding/json"

	"minik8s.io/pkg/apis/core"
)

func HandleFuncRegister(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	function := core.Function{}
	json.Unmarshal(body, &function)
	funcName := function.Name
	filePath := function.Path
	// Param miss: return error to client
	if funcName == "" || filePath == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("funcName and filePath should not be empty string"))
		return
	}
	// Send function param to apiserver /func/register
	

	// Construct pod param & Send http request to apiserver /deployment/status/apply
	deployment := ConstructDeploymentWithFuncNameAndFilepath(funcName, filePath)

}

func ConstructDeploymentWithFuncNameAndFilepath(funcName, filePath string) core.Deployment {

}