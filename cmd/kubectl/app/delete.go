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
	switch resourceKind {
	case "replicaset":
		return nil
		//add more
		//...
	case "service":
		fmt.Printf("Delete service name: %s\n", resourceName)
		err := deleteService(resourceName)
		if err != nil {
			fmt.Printf("Error in sending del service msg to apiserver: %v", err)
		}
	}

	return nil
}

func deleteService(serviceName string) error {
	fmt.Println("del service")
	params := map[string]string {
		"namespace"	: "default",
		"name"		: serviceName,
	}
	rawResp, err := clientutil.HttpDel("Service", params)
	if err != nil {
		return fmt.Errorf("Error in clientutil: %v", err)
	}
	fmt.Printf("deleteService response body: %s\n", string(rawResp))
	return nil;
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
