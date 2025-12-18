package distribution

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCommand(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distribution",
		Short: "Manage CDN distributions",
		Long:  "Manage the lifecycle of CDN distributions.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
}
