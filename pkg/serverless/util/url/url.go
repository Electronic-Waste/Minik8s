package url

import (
	globalURL "minik8s.io/pkg/util/url"
)

// Define some url used for pod mounting and serverless manager
const (
	HttpScheme 					= "http://"
	ManagerURL 					= globalURL.MasterNodeIP
	ManagerPort 				= ":8081"

	ManagerPrefix 				= HttpScheme + ManagerURL + ManagerPort

	Function					= "/func"
	FunctionRegisterURL			= Function + "/register"
	FunctionTriggerURL			= Function + "/trigger"

	DeploymentNamePrefix		= "deployment-svl-"
	PodNamePrefix 				= "pod-svl-"
	VolumeNamePrefix			= "vol-svl-"
	ContainerNamePrefix			= "svl-"

	HostMountDestPathPrefix 	= "/root/func"
	HostMountSrcPathPrefix 		= "/root/minik8s/minik8s/pkg/serverless/app"
	ContainerMountPath 			= "/python/src"
	ContainerImage 				= "docker.io/library/python:3.8.10"
	ContainerExposedPort 		= 8080

	ScriptFile 					= "start.sh"
	PythonServerFile 			= "main.py"
	FuncFile 					= "func.py"
	RequirementsFile 			= "requirements.txt"

	HostSrcScriptPath 			= HostMountSrcPathPrefix + "/" + ScriptFile
	HostSrcPythonServerPath 	= HostMountSrcPathPrefix + "/" + PythonServerFile
	// We need to specify the source function file path and copy its content to destination function file
	HostSrcRequirementsPath 	= HostMountSrcPathPrefix + "/" + RequirementsFile

	ContainerScriptPath			= ContainerMountPath + "/" + ScriptFile
	ContainerPythonServerPath	= ContainerMountPath + "/" + PythonServerFile
	ContainerFuncPath			= ContainerMountPath + "/" + FuncFile
	ContainerRequirementsPath	= ContainerMountPath + "/" + RequirementsFile
)