package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/podmanager"
)

var (
	FormatNodes = []string{
		"Name", "MasterIp", "NodeIp", "NodeStatus",
	}
	FormatJobs = []string{
		"Name", "Pod",
	}
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
		pods, _ := podmanager.GetPods()
		output := "NAMESPACE\tKIND\tNAME\t\t\t\t\tSTATUS\t\n"
		for _, p := range pods {
			output += "default\t\t" + "Pod\t" + p.Name + "\t" + string(p.Status.Phase) + "\n"
		}
		fmt.Println(output)
	case "deployment":
		//fmt.Println("get deployments")
		bytes, err := clientutil.HttpGetAll("Deployment")
		if err != nil {
			return err
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		fmt.Println("get deployment number: ", len(strs))
		var deployments []core.Deployment
		for _, s := range strs {
			if s == "" {
				continue
			}
			deployment := core.Deployment{}
			_ = json.Unmarshal([]byte(s), &deployment)
			deployments = append(deployments, deployment)
		}
		output := "NAMESPACE\tKIND\tNAME\tSTATUS\t\n"
		for _, d := range deployments {
			output += "default\t\t" + "Deployment\t" + d.Metadata.Name + "\t" + "Running" + "\n"
		}
		fmt.Println(output)
	case "autoscaler":
		fmt.Println("get autoscalers")
		bytes, err := clientutil.HttpGetAll("Autoscaler")
		if err != nil {
			return err
		}
		var strs []string
		err = json.Unmarshal(bytes, &strs)
		fmt.Println("get autoscaler number: ", len(strs))
		var autoscalers []core.Autoscaler
		for _, s := range strs {
			if s == "" {
				continue
			}
			autoscaler := core.Autoscaler{}
			_ = json.Unmarshal([]byte(s), &autoscaler)
			autoscalers = append(autoscalers, autoscaler)
		}
		output := "NAMESPACE\tKIND\t\tNAME\t\tSTATUS\t\n"
		for _, a := range autoscalers {
			output += "default\t\t" + "Autoscaler\t" + a.Metadata.Name + "\t" + "Running" + "\n"
		}
		fmt.Println(output)
	case "jobs":
		// deal with 'kubectl get nodes'
		bytes, err := clientutil.HttpGet("jobs", map[string]string{})
		if err != nil {
			return err
		}
		maps := core.JobMaps{}
		json.Unmarshal(bytes, &maps)
		FormatPrinting(FormatJobs, maps)
	}
	return nil
}

func FormatPrinting(formarStr []string, any interface{}) {
	for _, str := range formarStr {
		fmt.Printf("%s       ", str)
	}
	if nodeList, ok := any.(core.NodeList); ok {
		for _, node := range nodeList.NodeArray {
			fmt.Printf("\n%s    %s     %s      %s", node.MetaData.Name, node.Spec.MasterIp, node.Spec.NodeIp, "Ready")
		}
		fmt.Println("")
	} else if maps, ok := any.(core.JobMaps); ok {
		for _, Map := range maps.Maps {
			fmt.Printf("\n%s    %s", Map.JobName, Map.PodName)
		}
		fmt.Println("")
	}
}

func GetHandlerWithName(resourceKind, resourceName string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
