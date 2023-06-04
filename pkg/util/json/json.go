package json

// a package used to parse the json file to get the config message

import (
	"encoding/json"
	"fmt"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"minik8s.io/pkg/apis/core"
	runtime "minik8s.io/cri-api/pkg/apis/runtime/v1"
)

func ParseContainerConfig(file_path string, config *runtime.ContainerConfig, writer io.Writer) {
	jsonFile, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, config)
	fmt.Fprintf(writer, "get the config of container is \n%s\n", config.String())
	return
}

func ParseSandBoxConfig(file_path string, config *runtime.PodSandboxConfig, writer io.Writer) {
	jsonFile, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, config)
	fmt.Fprintf(writer, "get the config of container is \n%s\n", config.String())
	return
}

// Parse struct in types.go to string in JSON form
// Params: @struct_to_parse: Input a struct variable
// Usage: var pod core.Pod; val, err := StructToJSONString(pod);
func StructToJSONString(struct_to_parse interface{}) (string, error) {
	jsonVal, err := json.Marshal(struct_to_parse)
	if err != nil {
		return "", errors.New("parse struct to JSONString failed!")
	}
	return string(jsonVal), nil
}

// Parse a JSON form string to struct
// Params: @string_to_convert: the string to convert
// @val: a variable of the type you want to parse to
// Usage: var pod core.Pod; interface_val, err := JSONStringToStruct(str, pod); pod = interface_val.(core.Pod)
func JSONStringToStruct(string_to_convert string, val interface{}) (interface{}, error) {
	switch val.(type) {
	case core.Pod:
		var retVal core.Pod
		json.Unmarshal([]byte(string_to_convert), &retVal)
		return retVal, nil
	case string:
		var retVal string
		json.Unmarshal([]byte(string_to_convert), &retVal)
		return retVal, nil
	default:
		return nil, errors.New("JSONStringToStruct: Unsupported type!")
	}
}
