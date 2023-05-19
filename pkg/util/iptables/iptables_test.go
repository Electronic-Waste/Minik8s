package iptables

import (
	"testing"
	"os/exec"
)

func TestInitAndDeinitIptables(t *testing.T) {
	var output []byte
	cli, err := NewIPTablesClient("127.0.0.1")
	if err != nil {
		t.Error("create iptables client error")
	}

	// 1. Test func InitServiceIPTables
	cli.InitServiceIPTables()
	cmd := exec.Command("iptables-save")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("1. output is: \n%s", string(output))
		t.Errorf("exec `iptables-save` error: %v", err)
	}
	t.Logf("1. output is: %s", string(output))

	// 2. Test func DeinitServiceIPTables
	cli.DeinitServiceIPTables()
	cmd = exec.Command("iptables-save")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("2. output is: \n%s", string(output))
		t.Errorf("exec `iptables-save` error: %v", err)
	}
	t.Logf("2. output is: %s", string(output))
}

