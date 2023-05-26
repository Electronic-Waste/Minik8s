package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
	"os"
	"strings"
	//"encoding/json"
	//"github.com/go-yaml/yaml"
)

var (
	applyCmd = &cobra.Command{
		Use:     "apply <pathname>",
		Short:   "apply a yaml or json file to minik8s",
		Long:    `apply a yaml or json file to minik8s`,
		Example: "apply ./cmd/config/test.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply called")
			//fmt.Println(args[0])
			if err := ApplyHandler(args[0]); err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

// ApplyHandler parse the filepath, read the file and analyze it
func ApplyHandler(path string) error {
	if strings.HasSuffix(path, ".yaml") {
		//get yaml file content
		fmt.Println("apply a yaml file")
		viper.SetConfigType("yaml")
	} else if strings.HasSuffix(path, ".json") {
		//get yaml file content
		fmt.Println("apply a json file")
		viper.SetConfigType("json")
	} else {
		return errors.New("not a yaml or json file")
	}
	file, err := os.ReadFile(path)
	err = viper.ReadConfig(bytes.NewReader(file))
	if err != nil {
		fmt.Println("error reading file, please use relative path\n for example: apply ./cmd/config/xxx.yml")
		return err
	}
	//apply to k8s according to yaml
	objectKind := viper.GetString("kind")
	fmt.Println(objectKind)
	switch objectKind {
	case "HorizontalPodAutoscaler":
		autoscaler := core.Autoscaler{}
		err := viper.Unmarshal(&autoscaler)
		if err != nil {
			return err
		}
		fmt.Printf("autoscaler: %s\n", autoscaler.Metadata.Name)
		//bytes,_ := json.Marshal(autoscaler)
		//fmt.Println(string(bytes))
		err = applyAutoscaler(autoscaler)
		if err != nil{
			return err
		}
	case "Deployment":
		deployment := core.Deployment{}
		err := viper.Unmarshal(&deployment)
		if err != nil {
			return err
		}
		fmt.Printf("deployment: %s\n", deployment.Metadata.Name)
		err = applyDeployment(deployment)
		if err != nil{
			return err
		}
		//TODO: add more case handlers
	case "Pod":
		pod := core.Pod{}
		//err := yaml.Unmarshal(file,&pod)
		err := viper.Unmarshal(&pod)
		if err != nil {
			return err
		}
		fmt.Printf("pod: %s\n", pod.Name)
		err = applyPod(pod)
		if err != nil{
			return err
		}
	case "Service":
		service := core.Service{}
		err := viper.Unmarshal(&service)
		if err != nil {
			fmt.Printf("Error in unmarshaling service yaml file: %v", err)
			return err
		}
		fmt.Printf("service name: %s\n", service.Name)
		err = applyService(service)
		if err != nil {
			fmt.Printf("Error in sending service to apiserver: %v", err)
			return err
		}
	}
	return nil
}

func applyPod(pod core.Pod) error {
	fmt.Println("apply pod")
	return clientutil.HttpApply("Pod", pod)
}

func applyDeployment(deployment core.Deployment) error {
	fmt.Println("apply deployment")
	return clientutil.HttpApply("Deployment", deployment)
}

func applyService(service core.Service) error {
	fmt.Println("apply service")
	return clientutil.HttpApply("Service", service)
}

func applyAutoscaler(autoscaler core.Autoscaler) error {
	fmt.Println("apply autoscaler")
	return clientutil.HttpApply("Autoscaler", autoscaler)
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
