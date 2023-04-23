package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const container_path = "./tests/container-config.json"
const sandbox_path = "./tests/sandbox-config.json"

func TestContainerConfig(t *testing.T) {
	var res *runtime.ContainerConfig
	buf := new(bytes.Buffer)
	ParseContainerConfig(container_path, res, buf)
	if !strings.Contains(buf.String(), "busybox") {
		fmt.Println(buf.String())
		t.Error("error parse")
	}
}

func TestSandBoxConfig(t *testing.T) {
	var res *runtime.PodSandboxConfig
	buf := new(bytes.Buffer)
	ParseSandBoxConfig(sandbox_path, res, buf)
	if !strings.Contains(buf.String(), "nginx-sandbox-too") {
		fmt.Println(buf.String())
		t.Error("error parse")
	}
}
