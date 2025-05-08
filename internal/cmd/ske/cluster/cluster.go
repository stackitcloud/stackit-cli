package cluster

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Provides functionality for SKE cluster",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(generatepayload.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
