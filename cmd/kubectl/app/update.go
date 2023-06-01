package app

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"

	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/apis/core"
)

var (
	updateCmd = &cobra.Command{
		Use:     "update <path-to-file>",
		Short:   "update a function in Knative",
		Long:    "update a function in Knative",
		Example: `update ./func/Add.py`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("update called")
			//fmt.Println(args[0])
			if err := UpdateHandler(args[0]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

func UpdateHandler(filepath string) error {
	if !strings.HasSuffix(filepath, ".py") {
		return fmt.Errorf("File type error: only support python file!")
	}
	blocks := strings.Split(filepath, "/")
	fileName := blocks[len(blocks) - 1]
	funcName := strings.Split(fileName, ".")[0]
	function := core.Function {
		Name: funcName,
		Path: filepath,
	}
	fmt.Printf("function: %v\n", function)
	return clientutil.HttpUpdate("Function", function)
}