package app

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"minik8s.io/pkg/serverless"
)

func NewKnativeCommand() *cobra.Command {
	rootCmd := &cobra.Command {
		Use: "knative",
		Short: "Knatvie is a FaaS platform based on Minik8s",
		Long: "Knative is a FaaS platform based on Minik8s",
		Run: func (cmd *cobra.Command, args []string) {
			fmt.Println("Knative starts!")
			knative := serverless.NewKnative()
			knative.Run()
		},
	}
	return rootCmd
}
