package kubeproxy

import (
	"fmt"
	"net"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/util/iptables"
	"minik8s.io/pkg/util/listwatch"
)

type Manager interface{
	// Create a service by adding chains and rules to iptables
	CreateService(
		serviceName string, clusterIP string, servicePorts []core.ServicePort,
		podNames []string, podIPs []string) error

	// Delete a service by deleting relevant chains and rule from iptables
	DelService(serviceName string) error

	// Processing request from redis: apply service
	HandleApplyService(msg *redis.Message)
	
	// Processing request from redis: delete service
	HandleDelService(msg *redis.Message)

	// Start Kubeproxy
	Run()
}

type KubeproxyManager struct {
	iptablesCli *iptables.IPTablesClient
}

func NewKubeProxy() (*KubeproxyManager, error) {
	cli, err := iptables.NewIPTablesClient("127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	return &KubeproxyManager{
		iptablesCli: cli,
	}, nil
}

func (manager *KubeproxyManager) Run() {
	// Bind list-watch function
	err := manager.iptablesCli.InitServiceIPTables()
	if err != nil {
		fmt.Printf("Error occured in init SerivceIPtables: %v", err)
	}
	go listwatch.Watch("/service/apply", manager.HandleApplyService)
	go listwatch.Watch("/service/del", manager.HandleDelService)
}

func (manager *KubeproxyManager) CreateService(
	serviceName string, 
	clusterIP string, 
	servicePorts []core.ServicePort,
	podNames []string, 
	podIPs []string) error {
	// Check whether clusterIP is valid or not
	if net.ParseIP(clusterIP) == nil {
		return fmt.Errorf("cluster IP %s is invalid", clusterIP)
	}
	// Check params
	if len(servicePorts) != len(podNames) ||
		len(podNames) != len(podIPs) {
		return fmt.Errorf("params' len mismatches!")
	}

	// Create service
	serviceChainName := manager.iptablesCli.CreateServiceChain()
	for i := range podNames {
		podChainName := manager.iptablesCli.CreatePodChain()
		if net.ParseIP(podIPs[i]) == nil {
			return fmt.Errorf("pod IP %s in invalid", podIPs[i])
		}
		err := manager.iptablesCli.ApplyPodChainRules(
			podChainName, 
			podIPs[i], 
			(uint16)(servicePorts[i].TargetPort),
		)
		if err != nil {
			return fmt.Errorf("Error in applying pod chain rules: %v", err)
		}
		err = manager.iptablesCli.ApplyPodChain(
			serviceName, 
			serviceChainName, 
			podNames[i], 
			podChainName, 
			i + 1,
		)
		if err != nil {
			return fmt.Errorf("Error in applying pod chain: %v", err)
		}
		err = manager.iptablesCli.ApplyServiceChain(
			serviceName, 
			clusterIP, 
			serviceChainName, 
			(uint16)(servicePorts[i].Port),
		)
		if err != nil {
			return fmt.Errorf("Error in applying service chain: %v", err)
		}
	}
	return nil
}

func (manager *KubeproxyManager) DelService(serviceName string) error {
	return nil
}

func (manager *KubeproxyManager) HandleApplyService(msg *redis.Message) {
	var params core.KubeproxyServiceParam
	json.Unmarshal([]byte(msg.Payload), &params)
	err := manager.CreateService(
		params.ServiceName, 
		params.ClusterIP,
		params.ServicePorts,
		params.PodNames,
		params.PodIPs,
	)
	if err != nil {
		fmt.Printf("Hanlde apply service error: %v", err)
	}
}

func (manager *KubeproxyManager) HandleDelService(msg *redis.Message) {
	serviceName := msg.Payload
	err := manager.DelService(serviceName)
	if err != nil {
		fmt.Printf("Handle delete service error: %v", err)
	}
}




