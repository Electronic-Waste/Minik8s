package config

import (
	"fmt"
	"minik8s.io/pkg/constant"
	"testing"
)

func TestNewFile(t *testing.T) {
	ch := make(chan interface{})
	NewSourceFile(ch)
	<-ch
}

func TestListFile(t *testing.T) {
	s, err := ListAllConfig(constant.SysPodDir)
	if err != nil {
		t.Error(err)
	}
	for _, ele := range s {
		fmt.Printf("file name is %s\n", ele)
	}
}
