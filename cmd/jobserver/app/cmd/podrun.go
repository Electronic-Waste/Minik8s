package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"minik8s.io/pkg/util/jobserver"
)

func NewCmdRemote(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "tun in the pod",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			scripts, err := cmd.Flags().GetString("scripts")
			if err != nil {
				return err
			}
			result, err := cmd.Flags().GetString("result")
			if err != nil {
				return err
			}
			js := jobserver.NewJobServer()
			return js.Run(false, file, scripts, result)
		},
	}
	cmd.Flags().String("file", "", "the config file name")
	cmd.Flags().String("scripts", "", "the config scripts name")
	cmd.Flags().String("result", "", "the config result name")
	return cmd
}
