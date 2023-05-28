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

	// Read hosts file
	content, err := ioutil.ReadFile(controller.HostsFilePath)
	if err != nil {
		return fmt.Errorf("Error in DelDNSRule: %v\n", err)
	}

	// Parse file content to []string
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}

	// Combine them & write back to hosts file
	ipHostPair := fmt.Sprintf("%s %s", targetIP, hostname)
	lines = append(lines, ipHostPair)
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(controller.HostsFilePath, []byte(output), 0644)
}

func (controller *DNSController) DelDNSRule(targetIP, hostname string) error {
	// Read hosts file
	content, err := ioutil.ReadFile(controller.HostsFilePath)
	if err != nil {
		return fmt.Errorf("Error in DelDNSRule: %v\n", err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}

	// Find target line and delete it & Write back to file
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