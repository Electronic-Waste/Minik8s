package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"io"
)

func NewServerCommand(in io.Reader, out, err io.Writer) *cobra.Command {

	cmds := &cobra.Command{
		Use:   "jobserver",
		Short: "jobserver: used to commit the job and receive the output",
		Long: dedent.Dedent(`
			nothing
		`),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewCmdLocal(out))
	cmds.AddCommand(NewCmdRemote(out))
	return cmds
}
