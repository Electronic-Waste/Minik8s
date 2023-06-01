package app

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"

	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	
)

var (
	registerCmd = &cobra.Command{
		Use:     "register <path-to-file>",
		Short:   "register a function to Knative",
		Long:    `register a function to Knative`,
		Example: "register ./func/Add.py",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("register called")
			//fmt.Println(args[0])
			if err := RegisterHandler(args[0]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

func RegisterHandler(filepath string) error {
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
	return clientutil.HttpRegister(function)
}