package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/clientutil"
)

func NewCmdJoin(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "join a minik8s group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunJoin(out, cmd)
		},
	}
	cmd.Flags().String("config", "", "the path of node config file")
	return cmd
}

func RunJoin(out io.Writer, cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Fprintf(out, "err is %v\n", err)
		return err
	}
	fmt.Fprintf(out, "file path is %s\n", config)
	node, err := core.ParseNode(config)
	if err != nil {
		fmt.Fprintf(out, "err is %v\n", err)
		return err
	}
	err = RegisterNode(node)
	return err
}

func RegisterNode(node *core.Node) error {
	err := clientutil.HttpApply("Node", *node)
	return err
}
