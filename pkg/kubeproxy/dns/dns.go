package dns

import (
	"os"
	"net"
	"fmt"
	"io/ioutil"
	"strings"
)

type DNSFuncInterface interface {
	// CreateDNSRule creates a <ip>:<hostname> mapping in hosts file
	CreateDNSRule(targetIP, hostname string) error

	// DelDNSRule deletes a <ip>:<hostname> mapping in hosts file
	DelDNSRule(targetIP, hostname string) error
}

type DNSController struct {
	HostsFilePath string
}

func NewDNSController(hostFilePath string) (*DNSController, error) {
	// Check if file exists
	_, err := os.Stat(hostFilePath)
	if err != nil {
		return nil, err
	}
	return &DNSController{
		HostsFilePath : hostFilePath,
	}, nil
}

func (controller *DNSController) CreateDNSRule(targetIP, hostname string) error {
	// Check whether the targetIP is valid or not
	if net.ParseIP(targetIP) == nil {
		return fmt.Errorf("Invalid targetIP: %s\n", targetIP)
	}
	// Open file with r/w privilege
	file, err := os.OpenFile(controller.HostsFilePath, os.O_RDWR | os.O_APPEND, 777)
	if err != nil {
		return fmt.Errorf("Error in CreateDNSRule: %v\n", err)
	}
	ipHostPair := fmt.Sprintf("\n%s %s", targetIP, hostname)
	var bytesWritten int
	bytesWritten, err = file.WriteString(ipHostPair)
	if err != nil {
		return fmt.Errorf("Error in CreateDNSRule: %v\n", err)
	}
	fmt.Printf("bytesWritten: %d", bytesWritten)
	return nil
}

func (controller *DNSController) DelDNSRule(targetIP, hostname string) error {
	content, err := ioutil.ReadFile(controller.HostsFilePath)
	if err != nil {
		return fmt.Errorf("Error in DelDNSRule: %v\n", err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}
	targetLine := fmt.Sprintf("%s %s", targetIP, hostname)
	for i := range lines {
		if lines[i] == targetLine {
			lines = append(lines[:i], lines[i + 1:]...)
			break
		}
	}
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(controller.HostsFilePath, []byte(output), 0644)
}