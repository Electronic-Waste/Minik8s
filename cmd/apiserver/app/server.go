package app

import (
	"fmt"
	"minik8s.io/pkg/apiserver"
	// "encoding/json"

	"github.com/spf13/cobra"
)

func NewAPIServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "minik8s-apiserver",
		Long: `The Kubernetes API server validates and configures data
for the api objects which include pods, services, replicationcontrollers, and
others. The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		Run: func(cmd *cobra.Command, args []string){
			fmt.Println("Minik8s's apiserver starts!")
			apiserver.Run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}
