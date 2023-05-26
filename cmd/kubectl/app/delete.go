package app

import (
	"fmt"
	"github.com/spf13/cobra"

	"minik8s.io/pkg/clientutil"
)

var (
	deleteCmd = &cobra.Command{
		Use:     "delete <resource> <resource-name>",
		Short:   "delete a resource from minik8s",
		Long:    `delete a resource from minik8s`,
		Example: "apply ",
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
	case "replicaset":
		return nil
		//add more
		//...
	case "service":
		err = deleteService(resourceName)
	case "dns":
		err = deleteDNS(resourceName)
	}
	return err
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

func init() {
	rootCmd.AddCommand(deleteCmd)
}
