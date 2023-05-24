package app

import (
	"github.com/spf13/cobra"
	scheduler2 "minik8s.io/pkg/scheduler"
)

func NewKubeletCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "scheduler",
		Short: "scheduler is a tools for Pod scheduling",
		Long:  `for more detail, see git@gitee.com:jinglinwei/minik8s.git`,
		Run: func(cmd *cobra.Command, args []string) {
			// Run Kubelet
			Run()
		},
	}

	return rootCmd
}

func Run() error {
	scheduler := scheduler2.GetNewScheduler()
	scheduler.Run()
	return nil
}
