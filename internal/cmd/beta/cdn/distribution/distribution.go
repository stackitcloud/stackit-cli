package distribution

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/cdn/distribution/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCommand(params *params.CmdParams) *cobra.Command {
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

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
}
