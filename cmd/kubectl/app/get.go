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
)

var (
	getCmd = &cobra.Command{
		Use:     "get <resources> | (<resource> <resource-name>)",
		Short:   "get a resource from minik8s",
		Long:    `get a resource from minik8s`,
		Example: "apply ",
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
	case "nodes":
		{
			// deal with 'kubectl get nodes'
			bytes, err := clientutil.HttpGet("nodes", map[string]string{})
			if err != nil {
				return err
			}
			nodeList := core.NodeList{}
			json.Unmarshal(bytes, &nodeList)
			FormatPrinting(FormatNodes, nodeList)
		}
	}
	return nil
}

func FormatPrinting(formarStr []string, any interface{}) {
	for _, str := range formarStr {
		fmt.Printf("%s       ", str)
	}
	nodeList := any.(core.NodeList)
	fmt.Printf("the num of node is %d\n", len(nodeList.NodeArray))
	for _, node := range nodeList.NodeArray {
		fmt.Printf("\n%s    %s     %s      %s", node.MetaData.Name, node.Spec.MasterIp, node.Spec.NodeIp, "Ready")
	}
	fmt.Println("")
}

func GetHandlerWithName(resourceKind, resourceName string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
