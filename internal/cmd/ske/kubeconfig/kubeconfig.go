package kubeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig/login"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Provides functionality for SKE kubeconfig",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) kubeconfig.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(login.NewCmd(params))
}
