package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
)

var (
	FormatNodes = []string{
		"Name", "MasterIp", "NodeIp", "NodeStatus",
	}
	FormatService = []string {
		"ServiceName", "ClusterIP", "PortName",
		"Port", "TargetPort", "ServiceStatus",
	}
	FormatDNS = []string {
		"DNSName", "Host", "SubPath", "ServiceName",
		"TargetPort", "DNSStatus", 
	}
)

var (
	getCmd = &cobra.Command{
		Use:     "get <resources> | (<resource> <resource-name>)",
		Short:   "get a resource from minik8s",
		Long:    `get a resource from minik8s`,
		Example: "apply ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get called")
			if len(args) == 1 {
				if err := GetHandler(args[0]); err != nil {
					fmt.Println(err.Error())
				}
			} else {
				if err := GetHandlerWithName(args[0], args[1]); err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
)

func GetHandler(resourceKind string) error {
	switch resourceKind {
	case "service":
		bytes, err := clientutil.HttpGetAll("Service")
		if err != nil {
			return fmt.Errorf("Error in getting service")
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		if err != nil {
			return err
		}
		serviceList := []core.Service{}
		for _, str := range strs {
			service := core.Service{}
			json.Unmarshal([]byte(str), &service)
			serviceList = append(serviceList, service)
		}
		FormatPrinting(FormatService, serviceList)
	case "dns":
		bytes, err := clientutil.HttpGetAll("DNS")
		if err != nil {
			return fmt.Errorf("Error in getting DNS")
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		if err != nil {
			return err
		}
		dnsList := []core.DNS{}
		for _, str := range strs {
			dns := core.DNS{}
			json.Unmarshal([]byte(str), &dns)
			dnsList = append(dnsList, dns)
		}
		FormatPrinting(FormatDNS, dnsList)
	}
		
	return nil
}

func FormatPrinting(formarStr []string, any interface{}) {
	for _, str := range formarStr {
		fmt.Printf("%s       ", str)
	}

	switch any.(type) {
	// case core.NodeList:
	// 	nodeList := any.(core.NodeList)
	// 	for _, node := range nodeList.NodeArray {
	// 		fmt.Printf("\n%s    %s     %s      %s", node.MetaData.Name, node.Spec.MasterIp, node.Spec.NodeIp, "Ready")
	// 	}
	case []core.Service:
		serviceList := any.([]core.Service)
		for _, service := range serviceList {
			for _, servicePort := range service.Spec.Ports {
				fmt.Printf("\n%s\t%s\t%s\t%d\t\t%d\t\t%s", 
					service.Name, service.Spec.ClusterIP, servicePort.Name, 
					servicePort.Port, servicePort.TargetPort, "READY")
			}
		}
	case []core.DNS:
		dnsList := any.([]core.DNS)
		for _, dns := range dnsList {
			for _, subpath := range dns.Spec.Subpaths {
				fmt.Printf("\n%s\t%s\t%s\t%s\t\t%d\t%s",
					dns.Name, dns.Spec.Host, subpath.Path, 
					subpath.Service, subpath.Port, "READY")
			}
		}
	}
	
	
	fmt.Println("")
}

func GetHandlerWithName(resourceKind, resourceName string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
