package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version of kubeadm",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(out, cmd)
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringP("output", "o", "", "the format of version need to be print")
	return cmd
}

func RunVersion(out io.Writer, cmd *cobra.Command) error {
	kind, _ := cmd.Flags().GetString("output")
	switch kind {
	case "":
		{
			fmt.Println("default version print")
			fmt.Fprint(out, "minik8s version 1.0")
		}
	case "yaml":
		{
			fmt.Println("yaml version print")
			fmt.Fprint(out, "minik8s version 1.0")
		}
	case "json":
		{
			fmt.Println("json version print")
			fmt.Fprint(out, "minik8s version 1.0")
		}
	}
}
