package ipget

import (
	"fmt"
	"strings"
	"os/exec"
	"net"

	"minik8s.io/pkg/kubeproxy/path"
)

func GetHostIP() (string, error) {
	cmd := exec.Command("ip", "a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Execute `ip a` failed!")
	}
	lines := strings.Split(string(output), "\n")
	matchStr := fmt.Sprintf("%s:", path.ENS3)
	var targetLine string
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		blocks := strings.Split(line, " ")
		if blocks[1] == matchStr {
			targetLine = lines[i + 2]
			break;
		}
 	}
	blocks := strings.Split(targetLine, " ")
	var ipNet string
	for i, block := range blocks {
		if block == "inet" {
			ipNet = blocks[i + 1]
		}
	}
	ip, _, _ := net.ParseCIDR(ipNet)
	return ip.String(), nil
}

func GetFlannelIP() (string, error) {
	cmd := exec.Command("ip", "a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Execute `ip a` failed!")
	}
	lines := strings.Split(string(output), "\n")
	matchStr := fmt.Sprintf("%s:", path.Flannel)
	var targetLine string
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		blocks := strings.Split(line, " ")
		if blocks[1] == matchStr {
			targetLine = lines[i + 2]
			break;
		}
 	}
	 blocks := strings.Split(targetLine, " ")
	 var ipNet string
	 for i, block := range blocks {
		 if block == "inet" {
			 ipNet = blocks[i + 1]
		 }
	 }
	 ip, _, _ := net.ParseCIDR(ipNet)
	 return ip.String(), nil
}