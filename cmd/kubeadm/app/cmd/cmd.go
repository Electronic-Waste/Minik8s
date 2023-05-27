package cmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewKubeadmCommand returns cobra.Command to run kubeadm command
func NewKubeadmCommand(in io.Reader, out, err io.Writer) *cobra.Command {

	cmds := &cobra.Command{
		Use:   "kubeadm",
		Short: "kubeadm: easily bootstrap a secure Kubernetes cluster",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ KUBEADM                                                  │
			    │ Easily bootstrap a secure Kubernetes cluster             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/kubernetes/kubeadm/issues             │
			    └──────────────────────────────────────────────────────────┘

			Example usage:

			    Create a two-machine cluster with one control-plane node
			    (which controls the cluster), and one worker node
			    (where your workloads, like Pods and Deployments run).

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the first machine:                                    │
			    ├──────────────────────────────────────────────────────────┤
			    │ control-plane# kubeadm init                              │
			    └──────────────────────────────────────────────────────────┘

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the second machine:                                   │
			    ├──────────────────────────────────────────────────────────┤
			    │ worker# kubeadm join <arguments-returned-from-init>      │
			    └──────────────────────────────────────────────────────────┘

			    You can then repeat the second step on as many other machines as you like.

		`),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewCmdVersion(out))
	cmds.AddCommand(NewCmdJoin(out))

	return cmds
}
