package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:     "get <resources> | (<resource> <resource-name>)",
		Short:   "delete a resource from minik8s",
		Long:    `delete a resource from minik8s`,
		Example: "apply ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
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
	return nil
}

func GetHandlerWithName(resourceKind, resourceName string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
