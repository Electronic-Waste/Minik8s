package app

import (
	"fmt"
	"github.com/spf13/cobra"
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
