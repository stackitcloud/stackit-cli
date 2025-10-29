package region

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "Provides functionality for regional configuration of STACKIT Network Area (SNA)",
		Long:  "Provides functionality for regional configuration of STACKIT Network Area (SNA).",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
