package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"minik8s.io/pkg/util/jobserver"
)

func NewCmdLocal(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "run in local",
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
			return js.Run(true, file, scripts, result)
		},
	}
	cmd.Flags().String("file", "", "the config file name")
	cmd.Flags().String("scripts", "", "the config scripts name")
	cmd.Flags().String("result", "", "the config result name")
	return cmd
}
