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
	case "pod":
		err = deletePod(resourceName)
	case "deployment":
		err = deleteDeployment(resourceName)
	case "autoscaler":
		err = deleteAutoscaler(resourceName)
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

func init() {
	rootCmd.AddCommand(deleteCmd)
}
