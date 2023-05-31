package register

import(
	"os"
	"fmt"
	"path"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	svlurl "minik8s.io/pkg/serverless/util/url"
	"minik8s.io/pkg/clientutil"
)

func HandleFuncRegister(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("HandleFuncRegister receive msg!")
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
	// 1. Send function param to apiserver /func/register
	err := clientutil.HttpApply("Function", function)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	// 2. Copy file to the host destination mount dir
	hostDestDir := path.Join(svlurl.HostMountDestPathPrefix, funcName)
	if err = os.MkdirAll(hostDestDir, 777); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to create mount dir: %v", err)))
		return
	}
	// Clear all existing files
	dir, err := ioutil.ReadDir(hostDestDir)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to Read dir content: %v", err)))
		return
	}
	for _, file := range dir {
		os.RemoveAll(path.Join([]string{hostDestDir, file.Name()}...))
	}
	// Copy script file & chmod +x
	hostDestScriptFilePath := path.Join(hostDestDir, svlurl.ScriptFile)
	err = CopyFile(hostDestScriptFilePath, svlurl.HostSrcScriptPath)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to copy script file: %v", err)))
		return
	}
	err = os.Chmod(hostDestScriptFilePath, 777)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to change script file's mode: %v", err)))
		return
	}
	// Copy main.py
	hostDestPythonServerPath := path.Join(hostDestDir, svlurl.PythonServerFile)
	err = CopyFile(hostDestPythonServerPath, svlurl.HostSrcPythonServerPath)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to copy python server file: %v", err)))
		return
	}
	// Copy requirements.txt
	hostDestRequirementsPath := path.Join(hostDestDir, svlurl.RequirementsFile)
	err = CopyFile(hostDestRequirementsPath, svlurl.HostSrcRequirementsPath)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to copy requirements file: %v", err)))
		return
	}
	// Copy <funcname>.py
	hostDestFuncPath := path.Join(hostDestDir, svlurl.FuncFile)
	err = CopyFile(hostDestFuncPath, funcFilePath)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("Failed to copy func file: %v", err)))
		return
	}

	// 3. Construct pod param & Send http request to apiserver /deployment/status/apply
	deployment := ConstructDeploymentWithFuncNameAndFilepath(funcName)
	err = clientutil.HttpApply("Deployment", deployment)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	// 4. Success
	resp.WriteHeader(http.StatusOK)
}

func ConstructDeploymentWithFuncNameAndFilepath(funcName string) core.Deployment {
	deployment := core.Deployment{
		Metadata: meta.ObjectMeta {
			Name: svlurl.DeploymentNamePrefix + funcName,
		},
		Spec: core.DeploymentSpec {
			Replicas: 1,
			Template: core.Pod {
				Kind: "Pod",
				Spec: core.PodSpec {
					Volumes: []core.Volume {
						{
							Name: svlurl.VolumeNamePrefix + funcName,
							HostPath: svlurl.HostMountDestPathPrefix + "/" + funcName,
						},
					},
					Containers: []core.Container {
						{
							Name: svlurl.ContainerNamePrefix + funcName,
							Image: svlurl.ContainerImage,
							VolumeMounts: []core.VolumeMount {
								{
									Name: svlurl.VolumeNamePrefix + funcName,
									MountPath: svlurl.ContainerMountPath,
								},
							},
							Ports: []core.ContainerPort {
								{
									ContainerPort: svlurl.ContainerExposedPort,
								},
							},
							Command: []string {
								svlurl.ContainerScriptPath,
								// "bash",
							},
						},
					},
				},
			},
		},
	}
	deployment.Spec.Template.Name = svlurl.PodNamePrefix + funcName
	deployment.Spec.Template.Labels = map[string]string{}
	deployment.Spec.Template.Labels["node"]	= "vmeet3"	// Should delete after test!!!
	return deployment
}

func CopyFile(destFilePath, srcFilePath string) error {
	// Check if srcFile exists
	_, err := os.Stat(srcFilePath)
	if err != nil {
		return err
	}

	// Create destFile and copy srcFile to it 
	file, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer file.Close()	
	content, err := ioutil.ReadFile(srcFilePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destFilePath, content, 644)
}