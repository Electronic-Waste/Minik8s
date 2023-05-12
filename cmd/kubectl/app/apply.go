package app

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	applyCmd = &cobra.Command{
		Use:     "apply <pathname>",
		Short:   "apply a yaml file to minik8s",
		Long:    `apply a yaml file to minik8s`,
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
	if strings.HasSuffix(path, ".yml") {
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
		//objectKind = viper.Get("objkind")
		//switch(objectKind)
		//case : handler...
		return nil
	}
	return errors.New("not a yaml file")
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
