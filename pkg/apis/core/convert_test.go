package core

import (
	"fmt"
	"minik8s.io/pkg/constant"
	"testing"
)

func TestPodParse(t *testing.T) {
	path := constant.SysPodDir + "/test.yaml"
	pod, err := ParsePod(path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*pod)
}

func TestNodeParse(t *testing.T) {
	path := constant.SysPodDir + "/../testcases/vmeet2.yaml"
	node, err := ParseNode(path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*node)
}
