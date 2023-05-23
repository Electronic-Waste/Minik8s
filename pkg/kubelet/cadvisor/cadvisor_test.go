package cadvisor

import "testing"

func TestGetMetrix(t *testing.T) {
	err := GetContainerMetric("go2")
	if err != nil {
		t.Error(err)
	}
}
