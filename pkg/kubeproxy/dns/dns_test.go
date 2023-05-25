package dns

import (
	// "fmt"
	"testing"
)

func TestCreatDNS(t *testing.T) {
	controller, err := NewDNSController("/root/minik8s/conf/test-hosts")
	if err != nil {
		t.Errorf("%v", err)
	}
	err = controller.CreateDNSRule("127.0.0.1", "test1.com")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestDelDNS(t *testing.T) {
	controller, err := NewDNSController("/root/minik8s/conf/test-hosts")
	if err != nil {
		t.Errorf("%v", err)
	}
	err = controller.DelDNSRule("127.0.0.1", "test1.com")
	if err != nil {
		t.Errorf("%v", err)
	}
}