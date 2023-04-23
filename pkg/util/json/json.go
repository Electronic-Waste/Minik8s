package json

// a package used to parse the json file to get the config message

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func ParseContainerConfig(file_path string, config *runtime.ContainerConfig, writer io.Writer) {
	jsonFile, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var res runtime.ContainerConfig
	json.Unmarshal(byteValue, &res)
	fmt.Fprintf(writer, "get the config of container is \n%s\n", res.String())
	config = &res
	return
}

func ParseSandBoxConfig(file_path string, config *runtime.PodSandboxConfig, writer io.Writer) {
	jsonFile, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var res runtime.PodSandboxConfig
	json.Unmarshal(byteValue, &res)
	fmt.Fprintf(writer, "get the config of container is \n%s\n", res.String())
	config = &res
	return
}
