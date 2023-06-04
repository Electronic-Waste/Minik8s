package cadvisor

import (
	"fmt"
	"testing"
)

func TestGetMetrix(t *testing.T) {
	status, err := GetContainerMetric("go2")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(status)
}
