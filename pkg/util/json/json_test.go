package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	runtime "minik8s.io/cri-api/pkg/apis/runtime/v1"
)

const container_path = "./tests/container-config.json"
const sandbox_path = "./tests/sandbox-config.json"

func TestContainerConfig(t *testing.T) {
	res := new(runtime.ContainerConfig)
	buf := new(bytes.Buffer)
	ParseContainerConfig(container_path, res, buf)
	if !strings.Contains(buf.String(), "busybox") {
		fmt.Println(buf.String())
		t.Error("error parse")
	}
	fmt.Println(buf.String())
	fmt.Printf("ge the json is \n%s\n", res.String())
}

func TestSandBoxConfig(t *testing.T) {
	res := new(runtime.PodSandboxConfig)
	buf := new(bytes.Buffer)
	ParseSandBoxConfig(sandbox_path, res, buf)
	if !strings.Contains(buf.String(), "nginx-sandbox-too") {
		fmt.Println(buf.String())
		t.Error("error parse")
	}

	fmt.Printf("ge the json is \n%s\n", res.String())
}

func TestStructAndJSONStringConvert(t *testing.T) {
	res, err := StructToJSONString("3")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("result: ", res)
	var ret string
	var interface_val interface{}
	interface_val, err = JSONStringToStruct(res, ret)
	if err != nil {
		t.Error(err.Error())
	}
	ret = interface_val.(string)
	t.Log("result: ", ret)
	if ret != "3" {
		t.Error("test failed!")
	}

}
