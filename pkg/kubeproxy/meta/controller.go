package meta

import (
	"minik8s.io/pkg/apis/core"
)

type Interface interface {
	AppendServiceChainNames(serviceName string, serviceChainNames []string)
	AppendClusterIP(serviceName, clusterIP string)
	AppendServicePorts(serviceName string, servicePorts []core.ServicePort)
	AppendPodNames(serviceName string, podNames []string)
	AppendPodChainName(podName, podChainName string)
	AppendPodIP(podName, podIP string)
	GetServiceChainName(serviceName string)
	DeleteServiceChainNames(serviceName string)
	DeleteClusterIP(serviceName string)
	DeleteServicePorts(serviceName string)
	DeletePodNames(serviceName string)
	DeletePodChainName(podName string)
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
	// podName -> podChainName (KUBE-SEP-)
	MapPodChainName			map[string]string
	// podName -> podIP
	MapPodIP				map[string]string
}

func NewMetaController() (*MetaController, error) {
	return &MetaController{
		MapServiceChainNames: map[string][]string{},
		MapClusterIP: map[string]string{},
		MapServicePorts: map[string][]core.ServicePort{},
		MapPodNames: map[string][]string{},
		MapPodChainName: map[string]string{},
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

func (controller *MetaController) AppendPodChainName(podName, podChainName string) {
	controller.MapPodChainName[podName] = podChainName
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

func (controller *MetaController) DeletePodChainName(podName string) {
	delete(controller.MapPodChainName, podName)
}

func (controller *MetaController) DeletePodIP(podName string) {
	delete(controller.MapPodIP, podName)
}

