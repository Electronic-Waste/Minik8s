package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/podmanager"
	"github.com/liushuochen/gotable"
)

var (
	FormatNodes = []string{
		"Name", "MasterIp", "NodeIp", "NodeStatus",
	}

	FormatJobs = []string{
		"Name", "Pod",
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
		Use:     "get <resources>",
		Short:   "get a resource from minik8s",
		Long:    `get a resource from minik8s`,
		Example: "get pod",
		Run: func(cmd *cobra.Command, args []string) {
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
		//FormatPrinting(FormatService, serviceList)
		table,_ := gotable.Create("ServiceName","ClusterIP","PortName","Port","TargetPort","ServiceStatus")
		rows := make([]map[string]string, 0)
		for _, service := range serviceList {
			for _, servicePort := range service.Spec.Ports{
				row := make(map[string]string)
				row["ServiceName"] = service.Name
				row["ClusterIP"] = service.Spec.ClusterIP
				row["PortName"] = servicePort.Name
				row["Port"] = fmt.Sprintf("%d",servicePort.Port) 
				row["TargetPort"] = fmt.Sprintf("%d",servicePort.TargetPort)
				row["ServiceStatus"] = "READY"
				rows = append(rows,row)
			}
		}
		table.AddRows(rows)
		fmt.Println(table)
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
		//FormatPrinting(FormatDNS, dnsList)
		table,_ := gotable.Create("DNSName", "Host", "SubPath", "ServiceName","TargetPort", "DNSStatus")
		rows := make([]map[string]string, 0)
		for _, dns := range dnsList {
			for _, subpath := range dns.Spec.Subpaths{
				row := make(map[string]string)
				row["DNSName"] = dns.Name
				row["Host"] = dns.Spec.Host
				row["SubPath"] = subpath.Path
				row["ServiceName"] = subpath.Service
				row["TargetPort"] = fmt.Sprintf("%d",subpath.Port)
				row["DNSStatus"] = "READY"
				rows = append(rows,row)
			}
		}
		table.AddRows(rows)
		fmt.Println(table)
	case "nodes":
		{
			// deal with 'kubectl get nodes'
			bytes, err := clientutil.HttpGet("nodes", map[string]string{})
			if err != nil {
				return err
			}
			nodeList := core.NodeList{}
			json.Unmarshal(bytes, &nodeList)
			//FormatPrinting(FormatNodes, nodeList)
			table,_ := gotable.Create("Name", "MasterIp", "NodeIp", "NodeStatus")
			rows := make([]map[string]string, 0)
			for _, node := range nodeList.NodeArray {
				row := make(map[string]string)
				row["Name"] = node.MetaData.Name
				row["MasterIp"] = node.Spec.MasterIp
				row["NodeIp"] = node.Spec.NodeIp
				row["NodeStatus"] = "Ready"
				rows = append(rows,row)
			}
			table.AddRows(rows)
			fmt.Println(table)
		}
	case "pod":
		//fmt.Println("get pods")
		//bytes,err := clientutil.HttpGetAll("Pod")
		//if err!=nil{
		//	return err
		//}
		//var strs []string
		//err = json.Unmarshal(bytes, &strs)
		//fmt.Println("get pod number: ", len(strs))
		//var pods []core.Pod
		//for _,s := range strs{
		//	if s == ""{
		//		continue
		//	}
		//	pod := core.Pod{}
		//	_ = json.Unmarshal([]byte(s), &pod)
		//	pods = append(pods, pod
		//}

		pods, _ := podmanager.GetPods()
		table,_ := gotable.Create("NAMESPACE","KIND","NAME","STATUS")
		rows := make([]map[string]string, 0)
		for _, p := range pods {
			row := make(map[string]string)
			row["NAMESPACE"] = "default"
			row["KIND"] = "Pod"
			row["NAME"] = p.Name
			row["STATUS"] = string(p.Status.Phase)
			rows = append(rows,row)
		}
		table.AddRows(rows)
		fmt.Println(table)
	case "deployment":
		//fmt.Println("get deployments")
		bytes, err := clientutil.HttpGetAll("Deployment")
		if err != nil {
			return err
		}
		//fmt.Println("get deployment number: ", len(strs))
		var strs []string
		var deployments []core.Deployment
		err = json.Unmarshal(bytes, &strs)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, s := range strs {
			if s == "" {
				continue
			}
			deployment := core.Deployment{}
			_ = json.Unmarshal([]byte(s), &deployment)
			deployments = append(deployments, deployment)
		}

		table,_ := gotable.Create("NAMESPACE","KIND","NAME","REPLICAS")
		rows := make([]map[string]string, 0)

		for _, d := range deployments {
			row := make(map[string]string)
			row["NAMESPACE"] = "default"
			row["KIND"] = "Deployment"
			row["NAME"] = d.Metadata.Name
			row["REPLICAS"] = fmt.Sprintf("%d/%d",d.Spec.Replicas,d.Spec.Replicas)
			rows = append(rows,row)
		}
		table.AddRows(rows)
		fmt.Println(table)
	case "autoscaler":
		fmt.Println("get autoscalers")
		bytes, err := clientutil.HttpGetAll("Autoscaler")
		if err != nil {
			return err
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("get autoscaler number: ", len(strs))
		var autoscalers []core.Autoscaler
		for _, s := range strs {
			if s == "" {
				continue
			}
			autoscaler := core.Autoscaler{}
			_ = json.Unmarshal([]byte(s), &autoscaler)
			autoscalers = append(autoscalers, autoscaler)
		}

		table,_ := gotable.Create("NAMESPACE","KIND","NAME","TARGET")
		rows := make([]map[string]string, 0)
		for _, a := range autoscalers {
			row := make(map[string]string)
			row["NAMESPACE"] = "default"
			row["KIND"] = "Autoscaler"
			row["NAME"] = a.Metadata.Name
			row["TARGET"] = a.Spec.ScaleTargetRef.Name
			rows = append(rows,row)
		}
		table.AddRows(rows)
		fmt.Println(table)
	case "jobs":
		// deal with 'kubectl get nodes'
		bytes, err := clientutil.HttpGet("jobs", map[string]string{})
		if err != nil {
			return err
		}
		maps := core.JobMaps{}
		json.Unmarshal(bytes, &maps)
		//FormatPrinting(FormatJobs, maps)
		table,_ := gotable.Create("Name", "Pod")
		rows := make([]map[string]string, 0)
		for _, Map := range maps.Maps {
			row := make(map[string]string)
			row["Name"] = Map.JobName
			row["Pod"] = Map.PodName
			rows = append(rows,row)
		}
		table.AddRows(rows)
		fmt.Println(table)
	default:
		fmt.Println("unknown input type")
	}
	return nil
}

func FormatPrinting(formarStr []string, any interface{}) {
	for _, str := range formarStr {
		fmt.Printf("%s       ", str)
	}


	switch any.(type) {
	case core.NodeList:
		nodeList := any.(core.NodeList)
		for _, node := range nodeList.NodeArray {
			fmt.Printf("\n%s    %s     %s      %s", node.MetaData.Name, node.Spec.MasterIp, node.Spec.NodeIp, "Ready")
		}
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
	case core.JobMaps:
		maps := any.(core.JobMaps)
		for _, Map := range maps.Maps {
			fmt.Printf("\n%s    %s", Map.JobName, Map.PodName)
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
