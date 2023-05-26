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
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/kubeproxy/dns"
	"minik8s.io/pkg/kubeproxy/nginx"
	"minik8s.io/pkg/kubeproxy/path"
)

type Manager interface{
	// Create a service by adding chains and rules to iptables
	CreateService(
		serviceName string, clusterIP string, servicePorts []core.ServicePort,
		podNames []string, podIPs []string) error

	// Delete a service by deleting relevant chains and rule from iptables
	DelService(serviceName string) error

	// Create a DNS rule to enable visiting service by domain name
	CreateDNS(hostName string, paths []core.DNSSubpath) error
	
	// Delete a DNS rule by Host path
	DelDNS(hostName string) error

	// Processing request from redis: apply service
	HandleApplyService(msg *redis.Message)
	
	// Processing request from redis: delete service
	HandleDelService(msg *redis.Message)

	// Processing request from redis: apply dns
	HandleApplyDNS(msg *redis.Message)

	// Processing request from redis: delete dns
	HandleDelDNS(msg *redis.Message)

	// Start Kubeproxy
	Run()
}

type KubeproxyManager struct {
	iptablesCli 	*iptables.IPTablesClient
	metaController 	*meta.MetaController
	dnsController	*dns.DNSController
	nginxController *nginx.NginxController
}

func NewKubeProxy() (*KubeproxyManager, error) {
	var cli *iptables.IPTablesClient
	var metaCtl *meta.MetaController
	var dnsCtl *dns.DNSController
	var nginxCtl *nginx.NginxController
	var err error
	cli, err = iptables.NewIPTablesClient("127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	metaCtl, err = meta.NewMetaController()
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	dnsCtl, err = dns.NewDNSController(path.HostsFile)
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	nginxCtl, err = nginx.NewNginxController(path.NginxConfFile)
	if err != nil {
		return nil, fmt.Errorf("Error occurred in creating new KubeProxy: %v", err)
	}
	return &KubeproxyManager{
		iptablesCli: 	cli,
		metaController:	metaCtl,
		dnsController: dnsCtl,
		nginxController: nginxCtl,
	}, nil
}

func (manager *KubeproxyManager) Run() {
	// Bind list-watch function
	err := manager.iptablesCli.InitServiceIPTables()
	if err != nil {
		fmt.Printf("Error occured in init SerivceIPtables: %v", err)
	}
	go listwatch.Watch(url.ServiceApplyURL, manager.HandleApplyService)
	go listwatch.Watch(url.ServiceDelURL, manager.HandleDelService)
	go listwatch.Watch(url.DNSApplyURL, manager.HandleApplyDNS)
	go listwatch.Watch(url.DNSDelURL, manager.HandleDelDNS)
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
		fmt.Printf("Hanlde apply service error: %v\n", err)
	}
}

func (manager *KubeproxyManager) HandleDelService(msg *redis.Message) {
	serviceName := msg.Payload
	err := manager.DelService(serviceName)
	if err != nil {
		fmt.Printf("Handle delete service error: %v\n", err)
	}
}

func (manager *KubeproxyManager) CreateDNS(hostName string, paths []core.DNSSubpath) error {
	err := manager.dnsController.CreateDNSRule(path.NginxIP, hostName)
	if err != nil {
		return err
	}
	err = manager.nginxController.ApplyNginxServer(hostName, paths)
	if err != nil {
		return err
	}
	return nil
}

func (manager *KubeproxyManager) DelDNS(hostName string) error {
	err := manager.dnsController.DelDNSRule(path.NginxIP, hostName)
	if err != nil {
		return err
	}
	err = manager.nginxController.DelNginxServer(hostName)
	if err != nil {
		return err
	}
	return nil
}

func (manager *KubeproxyManager) HandleApplyDNS(msg *redis.Message) {
	var dnsParams core.DNS
	json.Unmarshal([]byte(msg.Payload), &dnsParams)
	err := manager.CreateDNS(
		dnsParams.Spec.Host,
		dnsParams.Spec.Subpaths,
	)
	if err != nil {
		fmt.Printf("Handle apply dns error: %v\n", err)
	}
}

func (manager *KubeproxyManager) HandleDelDNS(msg *redis.Message) {
	hostName := msg.Payload
	err := manager.DelDNS(hostName)
	if err != nil {
		fmt.Printf("Handle delete dns error: %v\n", err)
	}
}

