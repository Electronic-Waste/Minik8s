package app

import (
	"fmt"
	"github.com/spf13/cobra"
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
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
