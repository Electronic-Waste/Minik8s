package ipgen

import (
	"testing"
)

func TestIPGen(t *testing.T) {
	ipgen, err := NewClusterIPGenerator()
	if err != nil {
		t.Errorf("Error in func NewClusterIPGeneraotr: %v", err)
	}
	// Test1: currentIP is 233.233.0.0
	t.Logf("1. IP: %s", ipgen.currentIP)
	if ipgen.currentIP != "233.233.0.0" {
		t.Errorf("invalid IP %s, expected 233.233.0.0", ipgen.currentIP)
	}
	// Test2: nextIP is 233.233.0.1
	var IP string
	IP, err = ipgen.NextClusterIP()
	t.Logf("2. IP: %s", IP)
	if err != nil || IP != "233.233.0.1" {
		t.Errorf("error or invalid IP addr")
	}
	// Test3: nextIP is 233.233.1.0
	for i := 0; i < 255; i++ {
		IP, err = ipgen.NextClusterIP()
	}
	t.Logf("3. IP: %s", IP)
	if err != nil || IP != "233.233.1.0" {
		t.Errorf("invalid IP %s, expected 233.233.1.0", IP)
	}
}