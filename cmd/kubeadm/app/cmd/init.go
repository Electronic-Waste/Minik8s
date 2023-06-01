package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"minik8s.io/pkg/util/file"
)

func NewCmdInit(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init a minik8s group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit(out, cmd)
		},
	}
	return cmd
}

func RunInit(out io.Writer, cmd *cobra.Command) error {
	fmt.Println("Init k8s group")
	err := file.GenConfigFile()
	return err
}
