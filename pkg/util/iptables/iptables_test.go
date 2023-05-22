package iptables

import (
	"net/http"
	"testing"
	"os/exec"
	"fmt"
	"io/ioutil"

	"minik8s.io/pkg/util/ipgen"
)

func TestInitAndDeinitIptables(t *testing.T) {
	// var output []byte
	// cli, err := NewIPTablesClient("127.0.0.1")
	// if err != nil {
	// 	t.Error("create iptables client error")
	// }

	// // 1. Test func InitServiceIPTables
	// cli.InitServiceIPTables()
	// cmd := exec.Command("iptables-save")
	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	t.Logf("1. output is: \n%s", string(output))
	// 	t.Errorf("exec `iptables-save` error: %v", err)
	// }
	// t.Logf("1. output is: %s", string(output))

	// // 2. Test func DeinitServiceIPTables
	// cli.DeinitServiceIPTables()
	// cmd = exec.Command("iptables-save")
	// output, err = cmd.CombinedOutput()
	// if err != nil {
	// 	t.Logf("2. output is: \n%s", string(output))
	// 	t.Errorf("exec `iptables-save` error: %v", err)
	// }
	// t.Logf("2. output is: %s", string(output))
}

func TestDelete(t *testing.T) {
	cli, err := NewIPTablesClient("127.0.0.1")
	if err != nil {
		t.Error("create iptables client error")
	}
	cli.DeinitServiceIPTables()
}

func TestServiceChain(t *testing.T) {
	cli, err := NewIPTablesClient("127.0.0.1")
	if err != nil {
		t.Error("create iptables client error")
	}
	cli.InitServiceIPTables()

	// Create service chain
	var clusterIP string
	var gen *ipgen.ClusterIPGenerator
	gen, err = ipgen.NewClusterIPGenerator()
	if err != nil {
		t.Error("Error in creating clusterIP generator")
	}
	clusterIP, err = gen.NextClusterIP()
	if err != nil {
		t.Error("Error in generate clusterIP")
	}
	serviceChainName := cli.CreateServiceChain()
	serviceName := "service-test"
	err = cli.ApplyServiceChain(serviceName, clusterIP, serviceChainName, 23333)
	if err != nil {
		t.Error("Error in applying service chain")
	}

	// Create pod chain
	podChainName := cli.CreatePodChain()
	err = cli.ApplyPodChain(serviceName, serviceChainName, "pod-test", podChainName, 1)
	if err != nil {
		t.Errorf("Error in applying pod chain: %v", err)
	}
	err = cli.ApplyPodChainRules(podChainName, "127.0.0.1", 8080)
	if err != nil {
		t.Errorf("Error in applying pod chain rules: %v", err)
	}

	// Output result
	var output []byte
	cmd := exec.Command("iptables-save")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("output is: \n%s", string(output))
		t.Errorf("exec `iptables-save` error: %v", err)
	}
	t.Logf("output is: %s", string(output))

	// Test
	var response *http.Response
	response, err = http.Get(fmt.Sprintf("http://%s:%d/", clusterIP, 23333))
	if err != nil {
		t.Errorf("Error in http response: %v", err)
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		t.Logf("response: %s", string(body))
		if string(body) != "test" {
			t.Error("Test result error!")
		}
	}
	

	// Clean chains and rules
	// err = cli.DeletePodChain("pod-test", podChainName)
	// if err != nil {
	// 	t.Logf("Error in deleting pod chain: %v", err)
	// }
	// cli.DeleteServiceChain(serviceName, clusterIP, serviceChainName, 23333)
	// if err != nil {
	// 	t.Logf("Error in deleting service chain: %v", err)
	// }
	// cli.DeinitServiceIPTables()
	// if err != nil {
	// 	t.Logf("Error in deiniting pod chain: %v", err)
	// }
}

