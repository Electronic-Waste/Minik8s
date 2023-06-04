package clientutil

import (
	"fmt"
	"testing"
)

func TestGetMem(t *testing.T) {
	err, str := HttpPlus("Mem", "", "http://127.0.0.1:3000/Pod/memory")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("get output is %s\n", str)
}
