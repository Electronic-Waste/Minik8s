package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/apis/core"
	"encoding/json"
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
		//fmt.Println("get pods")
		//bytes,err := clientutil.HttpGetAll("Pod")
		//if err!=nil{
		//	return err
		//}
		//var strs []string
		//err = json.Unmarshal(bytes, &strs)
		//fmt.Println("get pod number: ", len(strs))
		//var pods []core.Pod
		//for _,s := range strs{
		//	if s == ""{
		//		continue
		//	}
		//	pod := core.Pod{}
		//	_ = json.Unmarshal([]byte(s), &pod)
		//	pods = append(pods, pod
		//}
		pods,_ := podmanager.GetPods()
		output := "NAMESPACE\tKIND\tNAME\t\t\t\t\tSTATUS\t\n"
		for _,p := range pods{
			output += "default\t\t" + "Pod\t" + p.Name + "\t" + string(p.Status.Phase) + "\n"
		}
		fmt.Println(output)
	case "deployment":
		//fmt.Println("get deployments")
		bytes,err := clientutil.HttpGetAll("Deployment")
		if err!=nil{
			return err
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		fmt.Println("get deployment number: ", len(strs))
		var deployments []core.Deployment
		for _,s := range strs{
			if s == ""{
				continue
			}
			deployment := core.Deployment{}
			_ = json.Unmarshal([]byte(s), &deployment)
			deployments = append(deployments, deployment)
		}
		output := "NAMESPACE\tKIND\tNAME\tSTATUS\t\n"
		for _,d := range deployments{
			output += "default\t\t" + "Deployment\t" + d.Metadata.Name + "\t" + "Running" + "\n"
		}
		fmt.Println(output)
	case "autoscaler":
		fmt.Println("get autoscalers")
		bytes,err := clientutil.HttpGetAll("Autoscaler")
		if err!=nil{
			return err
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		fmt.Println("get autoscaler number: ", len(strs))
		var autoscalers []core.Autoscaler
		for _,s := range strs{
			if s == ""{
				continue
			}
			autoscaler := core.Autoscaler{}
			_ = json.Unmarshal([]byte(s), &autoscaler)
			autoscalers = append(autoscalers, autoscaler)
		}
		output := "NAMESPACE\tKIND\t\tNAME\t\tSTATUS\t\n"
		for _,a := range autoscalers{
			output += "default\t\t" + "Autoscaler\t" + a.Metadata.Name + "\t" + "Running" + "\n"
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
