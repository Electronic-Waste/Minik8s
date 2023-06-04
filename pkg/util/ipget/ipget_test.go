package ipget

import (
	"testing"
)

func TestIPGet(t *testing.T) {
	var hostIP, flannelIP string
	var err error
	hostIP, err = GetHostIP()
	if err != nil || hostIP != "192.168.1.7" {
		t.Errorf("Wrong hostIP: %s\n", hostIP)
	}
	flannelIP, err = GetFlannelIP()
	if err != nil || flannelIP != "10.0.6.0" {
		t.Errorf("Wrong flannelIP: %s\n", flannelIP)
	}
}