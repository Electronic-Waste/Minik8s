package kubeproxy

import (
	"fmt"
	"net"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/util/iptables"
	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/kubeproxy/meta"
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
	iptablesCli 	*iptables.IPTablesClient
	metaController 	*meta.MetaController
}

func NewKubeProxy() (*KubeproxyManager, error) {
	var cli *iptables.IPTablesClient
	var controller *meta.MetaController
	var err error
	cli, err = iptables.NewIPTablesClient("127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	controller, err = meta.NewMetaController()
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	return &KubeproxyManager{
		iptablesCli: 	cli,
		metaController:	controller,
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
	if len(podNames) != len(podIPs) {
		return fmt.Errorf("params' len mismatches!")
	}

	// Create service
	serviceChainNames := []string{}
	manager.metaController.AppendClusterIP(serviceName, clusterIP)
	manager.metaController.AppendServicePorts(serviceName, servicePorts)
	manager.metaController.AppendPodNames(serviceName, podNames)
	for _, servicePort := range servicePorts {
		serviceChainName := manager.iptablesCli.CreateServiceChain()
		serviceChainNames = append(serviceChainNames, serviceChainName)
		podChainNames := []string{}
		for i := range podNames {
			podChainName := manager.iptablesCli.CreatePodChain()
			podChainNames = append(podChainNames, podChainName)
			if net.ParseIP(podIPs[i]) == nil {
				return fmt.Errorf("pod IP %s is invalid", podIPs[i])
			}
			manager.metaController.AppendPodChainNameToPodName(podChainName, podNames[i])
			manager.metaController.AppendPodIP(podNames[i], podIPs[i])
			// 1. Create KUBE-SEP- rule
			err := manager.iptablesCli.ApplyPodChainRules(
				podChainName, 
				podIPs[i], 
				(uint16)(servicePort.TargetPort),
			)
			if err != nil {
				return fmt.Errorf("Error in applying pod chain rules: %v", err)
			}
			// 2. Create KUBE-SVC- -> KUBE-SEP- rule
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
		}
		// Map serviceChainName -> podChainNames
		manager.metaController.AppendPodChainNames(serviceChainName, podChainNames)
		// Create KUBE-SERVICES -> KUBE-SVC- rule
		err := manager.iptablesCli.ApplyServiceChain(
			serviceName, 
			clusterIP, 
			serviceChainName, 
			(uint16)(servicePort.Port),
		)
		if err != nil {
			return fmt.Errorf("Error in applying service chain: %v", err)
		}

	}
	manager.metaController.AppendServiceChainNames(serviceName, serviceChainNames)
	return nil
}

func (manager *KubeproxyManager) DelService(serviceName string) error {
	// Get relevant information from metaController
	serviceChainNames := manager.metaController.MapServiceChainNames[serviceName]
	clusterIP := manager.metaController.MapClusterIP[serviceName]
	servicePorts := manager.metaController.MapServicePorts[serviceName]
	
	// Clear service
	for i, _ := range serviceChainNames {
		// Delete some rules and chains in service:
		// 1. The rule jumping from KUBE-SERVICES to KUBE-SVC-
		// 2. The chain KUBE-SVC-
		err := manager.iptablesCli.DeleteServiceChain(
			serviceName,
			clusterIP,
			serviceChainNames[i],
			(uint16)(servicePorts[i].Port),
		)
		if err != nil {
			return fmt.Errorf("Error in deleting serviceChain %s: %v", serviceChainNames[i], err)
		}

		podChainNames := manager.metaController.MapPodChainNames[serviceChainNames[i]]
		for _, podChainName := range podChainNames {
			// Delete some rule and chain in pod:
			// 1. The chain KUBE-SEP
			podName := manager.metaController.MapPodChainNameToPodName[podChainName]
			err = manager.iptablesCli.DeletePodChain(
				podName,
				podChainName,
			)
			if err != nil {
				return fmt.Errorf("Error in deleting podChain %s: %v", podChainName, err)
			}

			// Update map data in metaController
			manager.metaController.DeletePodChainNameToPodName(podChainName)
			manager.metaController.DeletePodIP(podName)
		}
		manager.metaController.DeletePodChainNames(serviceChainNames[i])
	}
	// Update map data in metaController
	manager.metaController.DeleteServiceChainNames(serviceName)
	manager.metaController.DeleteClusterIP(serviceName)
	manager.metaController.DeleteServicePorts(serviceName)
	manager.metaController.DeletePodNames(serviceName)
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

