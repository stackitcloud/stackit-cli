package publicip

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/associate"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/disassociate"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/ranges"
	"github.com/stackitcloud/stackit-cli/internal/cmd/public-ip/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-ip",
		Short: "Provides functionality for public IPs",
		Long:  "Provides functionality for public IPs.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(associate.NewCmd(params))
	cmd.AddCommand(disassociate.NewCmd(params))
	cmd.AddCommand(ranges.NewCmd(params))
}
