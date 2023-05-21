package ipgen

import (
	"net"
	"fmt"
)

const (
	clusterIPRange = "233.233.0.0/16"
)

type ClusterIPGenerator struct {
	currentIP string
}

type IPGen interface {
	NextClusterIP() string
}

func NewClusterIPGenerator() (*ClusterIPGenerator, error) {
	newIP, _, err := net.ParseCIDR(clusterIPRange)
	if err != nil {
		return nil, fmt.Errorf("Error %v in parsing clusterIPRange", err)
	}
	return &ClusterIPGenerator{
		currentIP: newIP.String(),
	}, nil
}

func (gen *ClusterIPGenerator) NextClusterIP() (string, error) {
	var newClusterIP net.IP
	// Parse current clusterIP to array
	prevClusterIP := net.ParseIP(gen.currentIP).To4()
	if prevClusterIP == nil {
		return "", fmt.Errorf("Error in parsing currentIP")
	}

	// Generate new clusterIP
	if prevClusterIP[3] + 1 != 0 {
		newClusterIP = net.IPv4(prevClusterIP[0], prevClusterIP[1], prevClusterIP[2], prevClusterIP[3] + 1)
	} else if prevClusterIP[2] + 1 != 0 {
		newClusterIP = net.IPv4(prevClusterIP[0], prevClusterIP[1], prevClusterIP[2] + 1, 0)
	} else {
		return "", fmt.Errorf("ClusterIp exceeds the valid range of ip")
	}
	gen.currentIP = newClusterIP.String()
	
	return newClusterIP.String(), nil
}