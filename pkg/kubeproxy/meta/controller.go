package meta

import (
	"minik8s.io/pkg/apis/core"
)

type Interface interface {
	AppendServiceChainNames(serviceName string, serviceChainNames []string)
	AppendClusterIP(serviceName, clusterIP string)
	AppendServicePorts(serviceName string, servicePorts []core.ServicePort)
	AppendPodNames(serviceName string, podNames []string)
	AppendPodChainNames(serviceChainName string, podChainNames []string)
	AppendPodChainNameToPodName(podChainName, podName string)
	AppendPodIP(podName, podIP string)
	GetServiceChainName(serviceName string)
	DeleteServiceChainNames(serviceName string)
	DeleteClusterIP(serviceName string)
	DeleteServicePorts(serviceName string)
	DeletePodNames(serviceName string)
	DeletePodChainNames(serviceChainName string)
	DeletePodChainNameToPodName(podChainName string)
	DeletePodIP(podName string)
}

type MetaController struct {
	// serviceName -> serviceChainNames (KUBE-SVC-)
	MapServiceChainNames	map[string][]string
	// serviceName -> clusterIP
	MapClusterIP			map[string]string
	// serviceName -> servicePorts
	MapServicePorts			map[string][]core.ServicePort
	// serviceName -> podNames
	MapPodNames				map[string][]string
	// serviceChainName (KUBE-SVC) -> podChainNames (KUBE-SEP-)
	MapPodChainNames			map[string][]string
	// podChainName (KUBE-SEP-) -> podName
	MapPodChainNameToPodName	map[string]string
	// podName -> podIP
	MapPodIP				map[string]string
}

func NewMetaController() (*MetaController, error) {
	return &MetaController{
		MapServiceChainNames: map[string][]string{},
		MapClusterIP: map[string]string{},
		MapServicePorts: map[string][]core.ServicePort{},
		MapPodNames: map[string][]string{},
		MapPodChainNames: map[string][]string{},
		MapPodChainNameToPodName:	map[string]string{},
		MapPodIP: map[string]string{},
	}, nil
}

func (controller *MetaController) AppendServiceChainNames(serviceName string, serviceChainNames []string) {
	controller.MapServiceChainNames[serviceName] = serviceChainNames
}

func (controller *MetaController) AppendClusterIP(serviceName, clusterIP string) {
	controller.MapClusterIP[serviceName] = clusterIP
}

func (controller *MetaController) AppendServicePorts(serviceName string, servicePorts []core.ServicePort) {
	controller.MapServicePorts[serviceName] = servicePorts
}

func (controller *MetaController) AppendPodNames(serviceName string, podNames []string) {
	controller.MapPodNames[serviceName] = podNames
}

func (controller *MetaController) AppendPodChainNames(serviceChainName string, podChainNames []string) {
	controller.MapPodChainNames[serviceChainName] = podChainNames
}

func (controller *MetaController) AppendPodChainNameToPodName(podChainName, podName string) {
	controller.MapPodChainNameToPodName[podChainName] = podName
}

func (controller *MetaController) AppendPodIP(podName, podIP string) {
	controller.MapPodIP[podName] = podIP
}

func (controller *MetaController) DeleteServiceChainNames(serviceName string) {
	delete(controller.MapServiceChainNames, serviceName)
}

func (controller *MetaController) DeleteClusterIP(serviceName string) {
	delete(controller.MapClusterIP, serviceName)
}

func (controller *MetaController) DeleteServicePorts(serviceName string) {
	delete(controller.MapServicePorts, serviceName)
}

func (controller *MetaController) DeletePodNames(serviceName string) {
	delete(controller.MapPodNames, serviceName)
}

func (controller *MetaController) DeletePodChainNames(serviceChainName string) {
	delete(controller.MapPodChainNames, serviceChainName)
}

func (controller *MetaController) DeletePodChainNameToPodName(podChainName string) {
	delete(controller.MapPodChainNameToPodName, podChainName)
}

func (controller *MetaController) DeletePodIP(podName string) {
	delete(controller.MapPodIP, podName)
}

