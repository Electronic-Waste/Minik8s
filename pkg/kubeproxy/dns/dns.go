package dns

import (
	"os"
	"net"
	"fmt"

	"minik8s.io/pkg/kubeproxy/path"
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

func NewDNSController() (*DNSController, error) {
	// Check if file exists
	_, err := os.Stat(path.HostsFile)
	if err != nil {
		return nil, err
	}
	return &DNSController{
		HostsFilePath : path.HostsFile,
	}, nil
}

func (controller *DNSController) CreateDNSRule(targetIP, hostname string) error {
	// Check whether the targetIP is valid or not
	if net.ParseIP(targetIP) != nil {
		return fmt.Errorf("Invalid targetIP: %s", targetIP)
	}
	// Open file with r/w privilege
	file, err := os.OpenFile(controller.HostsFilePath, os.O_RDWR, 777)
	if err != nil {
		return err
	}
	return nil
}

func (controller *DNSController) DelDNSRule(targetIP, hostname string) error {
	return nil
}