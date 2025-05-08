package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for SQLServer Flex instances",
		Long:  "Provides functionality for SQLServer Flex instances.",
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
}
