package app

import (
	"fmt"
	"github.com/spf13/cobra"

	"minik8s.io/pkg/clientutil"
	svlurl "minik8s.io/pkg/serverless/util/url"
)

var (
	triggerCmd = &cobra.Command{
		Use:     "trigger <funName> <json>",
		Short:   "Trigger a function in Knative",
		Long:    "Trigger a function in Knative",
		Example: `trigger Add '{"x":1,"y":2}'`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("trigger called")
			//fmt.Println(args[0])
			if err := TriggerHandler(args[0], args[1]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(triggerCmd)
}

func TriggerHandler(funcName, jsonVal string) error {
	triggerURL := svlurl.ManagerPrefix + svlurl.FunctionTriggerURL + "?name=" + funcName
	result, err := clientutil.HttpTrigger("Kubectl-Function", triggerURL, jsonVal)
	fmt.Printf("Trigger func: %s, and result is %s\n", funcName, result)
	return err
}