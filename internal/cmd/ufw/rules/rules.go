package rules

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ufw/rules/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ufw/rules/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ufw/rules/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ufw/rules/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ufw/rules/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Provides functionality for UFW rules",
		Long:  "Provides functionality for STACKIT Unified Firewall (UFW) rules.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
