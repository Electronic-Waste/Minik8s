package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
)

var (
	deleteCmd = &cobra.Command{
		Use:     "delete <resource> <resource-name>",
		Short:   "delete a resource from minik8s",
		Long:    `delete a resource from minik8s`,
		Example: "delete pod go1",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
			if err := DeleteHandler(args[0], args[1]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func DeleteHandler(resourceKind, resourceName string) error {
	var err error
	switch resourceKind {
	case "service":
		err = deleteService(resourceName)
	case "dns":
		err = deleteDNS(resourceName)
	case "pod":
		err = deletePod(resourceName)
	case "deployment":
		err = deleteDeployment(resourceName)
	case "autoscaler":
		err = deleteAutoscaler(resourceName)
	case "node":
		err = deleteNode(resourceName)
	case "func":
		err = deleteFunc(resourceName)
	}
	return err
}

func deletePod(resourceName string) error {
	fmt.Println("delete pod")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = resourceName
	return clientutil.HttpDel("Pod", params)
}

func deleteNode(nodeName string) error {
	fmt.Printf("try to delete %s\n", nodeName)
	err, _ := clientutil.HttpPlus("Node", nodeName, url.Prefix+url.NodeDelUrl)
	return err
}

func deleteDeployment(resourceName string) error {
	fmt.Println("delete deployment")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = resourceName
	return clientutil.HttpDel("Deployment", params)
}

func deleteAutoscaler(resourceName string) error {
	fmt.Println("delete autoscaler")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = resourceName
	return clientutil.HttpDel("Autoscaler", params)
}

func deleteService(serviceName string) error {
	fmt.Println("del service")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = serviceName
	return clientutil.HttpDel("Service", params)
}

func deleteDNS(dnsName string) error {
	fmt.Println("del dns")
	params := make(map[string]string)
	params["namespace"] = "default"
	params["name"] = dnsName
	return clientutil.HttpDel("DNS", params)
}

func deleteFunc(funcName string) error {
	fmt.Println("del func")
	params := make(map[string]string)
	params["name"] = funcName
	return clientutil.HttpDel("Kubectl-Function", params)
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
