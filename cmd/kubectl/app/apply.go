package app

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"minik8s.io/pkg/apis/core"
	"os"
	"path/filepath"
	"strings"
)

var (
	applyCmd = &cobra.Command{
		Use:     "apply <pathname>",
		Short:   "apply a yaml or json file to minik8s",
		Long:    `apply a yaml or json file to minik8s`,
		Example: "apply ./cmd/config/xxx/yml",
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

		workDir, err := os.Getwd()
		if err != nil {
			return err
		}
		dir, file := filepath.Split(path)
		workDir += dir[1:len(dir)]
		//fmt.Println(workDir)
		viper.SetConfigFile(file)
		viper.AddConfigPath(workDir)
		err = viper.ReadInConfig()
		if err != nil {
			//fmt.Println("error reading file, please use relative path\n for example: apply ./cmd/config/xxx.yml")
			return err
		}
		//apply to k8s according to yaml
		//...
		objectKind := viper.GetString("kind")
		switch objectKind {
		case "Deployment":
			deployment := core.Deployment{}
			err := viper.Unmarshal(&deployment)
			if err != nil {
				return err
			}
			err = applyDeployment(deployment)
		}
		return nil
	}
	return errors.New("not a yaml file")
}

func applyDeployment(core.Deployment) error {

	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
