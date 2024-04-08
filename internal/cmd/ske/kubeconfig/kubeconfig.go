package kubeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig/create"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Provides functionality for SKE kubeconfig",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) kubeconfig.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
}
