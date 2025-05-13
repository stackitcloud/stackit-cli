package ske

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/kubeconfig"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/options"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ske",
		Short: "Provides functionality for SKE",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE).",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(kubeconfig.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
	cmd.AddCommand(cluster.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params))
}
