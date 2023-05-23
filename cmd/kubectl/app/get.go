package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/podmanager"
)

var (
	getCmd = &cobra.Command{
		Use:     "get <resources>",
		Short:   "get a resource from minik8s",
		Long:    `get a resource from minik8s`,
		Example: "get pod",
		Run: func(cmd *cobra.Command, args []string) {
			if err := GetHandler(args[0]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func GetHandler(resourceKind string) error {
	fmt.Println("get handler")
	switch resourceKind{
	case "pod":
		fmt.Println("get pods")
		pods,err := podmanager.GetPods()
		if err!=nil{
			return err
		}
		output := "NAMESPACE\tNAME\tSTATUS\t\n"
		for _,p := range pods{
			output += "default\t\t" + p.Name + "\t" + string(p.Status.Phase) + "\n"
		}
		fmt.Println(output)
	}
	return nil
}

func GetHandlerWithName(resourceKind, resourceName string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
