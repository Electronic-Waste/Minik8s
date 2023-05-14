package app

import (
	"github.com/spf13/cobra"
	"minik8s.io/pkg/kubelet"
	"minik8s.io/pkg/kubelet/config"
)

func NewKubeletCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "minikubelet",
		Short: "minikubelet is a tools for admin control plane",
		Long:  `for more detail, see git@gitee.com:jinglinwei/minik8s.git`,
		Run: func(cmd *cobra.Command, args []string) {
			// Run Kubelet
			Run()
		},
	}

	return rootCmd
}

func Run() error {
	return RunKubelet()
}

func RunKubelet() error {
	var podConfig *config.PodConfig
	klet, err := createAndInitKubelet(podConfig)
	if err != nil {
		return err
	}
	klet.Run(podConfig.Updates())
	return nil
}

func createAndInitKubelet(podConfig *config.PodConfig) (kubelet.Bootstrap, error) {
	// init a kubelet object
	klet, err := kubelet.NewMainKubelet(podConfig)
	if err != nil {
		return nil, err
	}
	return klet, nil
}
