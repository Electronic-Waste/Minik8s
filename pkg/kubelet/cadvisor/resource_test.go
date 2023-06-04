package cadvisor

import (
	"fmt"
	"testing"
)

func TestGetMem(t *testing.T) {
	num, err := GetFreeMem()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("get free mem is %d MB\n", num)
}
